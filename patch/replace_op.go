package patch

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type ReplaceOp struct {
	Path  Pointer
	Value interface{} // will be cloned using yaml library
}

type mutationCtx struct {
	PrevUpdate func(interface{})
	I          int
	Obj        interface{}
}

func replaceOpCloneValueErr(err error) error {
	return fmt.Errorf("ReplaceOp cloning value: %s", err)
}

func (op ReplaceOp) Apply(doc interface{}) (interface{}, error) {
	tokens := op.Path.Tokens()

	if len(tokens) == 1 {
		// Ensure that value is not modified by future operations
		clonedValue, err := op.cloneValue(op.Value)
		if err != nil {
			return nil, replaceOpCloneValueErr(err)
		}
		return clonedValue, nil
	}

	ctxStack := []*mutationCtx{&mutationCtx{
		PrevUpdate: func(newObj interface{}) { doc = newObj },
		I:          0,
		Obj:        doc,
	}}
	for len(ctxStack) != 0 {
		// Pop the next context off the stack
		ctx := ctxStack[len(ctxStack)-1]
		ctxStack = ctxStack[:len(ctxStack)-1]

		// Terminate if done
		if ctx.I+1 >= len(tokens) {
			continue
		}

		token := tokens[ctx.I+1]
		isLast := ctx.I == len(tokens)-2

		switch typedToken := token.(type) {
		case IndexToken:
			idx := typedToken.Index

			typedObj, ok := ctx.Obj.([]interface{})
			if !ok {
				return nil, newOpArrayMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			if idx >= len(typedObj) {
				return nil, opMissingIndexErr{idx, typedObj}
			}

			if isLast {
				clonedValue, err := op.cloneValue(op.Value)
				if err != nil {
					return nil, replaceOpCloneValueErr(err)
				}
				typedObj[idx] = clonedValue
			} else {
				ctxStack = append(ctxStack, &mutationCtx{
					PrevUpdate: func(newObj interface{}) { typedObj[idx] = newObj },
					I:          ctx.I + 1,
					Obj:        typedObj[idx],
				})
			}

		case AfterLastIndexToken:
			typedObj, ok := ctx.Obj.([]interface{})
			if !ok {
				return nil, newOpArrayMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			if isLast {
				clonedValue, err := op.cloneValue(op.Value)
				if err != nil {
					return nil, replaceOpCloneValueErr(err)
				}
				ctx.PrevUpdate(append(typedObj, clonedValue))
			} else {
				return nil, fmt.Errorf("Expected after last index token to be last in path '%s'", op.Path)
			}

		case MatchingIndexToken:
			typedObj, ok := ctx.Obj.([]interface{})
			if !ok {
				return nil, newOpArrayMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			var idxs []int

			for itemIdx, item := range typedObj {
				typedItem, ok := item.(map[interface{}]interface{})
				if ok {
					if typedItem[typedToken.Key] == typedToken.Value {
						idxs = append(idxs, itemIdx)
					}
				}
			}

			if typedToken.Optional && len(idxs) == 0 {
				if isLast {
					clonedValue, err := op.cloneValue(op.Value)
					if err != nil {
						return nil, replaceOpCloneValueErr(err)
					}
					ctx.PrevUpdate(append(typedObj, clonedValue))
				} else {
					o := map[interface{}]interface{}{typedToken.Key: typedToken.Value}
					ctx.PrevUpdate(append(typedObj, o))
					ctxStack = append(ctxStack, &mutationCtx{
						PrevUpdate: ctx.PrevUpdate, // no need to change prevUpdate since matching item can only be a map
						I:          ctx.I + 1,
						Obj:        o,
					})
				}
			} else {
				if len(idxs) != 1 {
					return nil, opMultipleMatchingIndexErr{NewPointer(tokens[:ctx.I+2]), idxs}
				}

				idx := idxs[0]

				if isLast {
					clonedValue, err := op.cloneValue(op.Value)
					if err != nil {
						return nil, replaceOpCloneValueErr(err)
					}
					typedObj[idx] = clonedValue
				} else {
					// no need to change prevUpdate since matching item can only be a map
					ctxStack = append(ctxStack, &mutationCtx{
						PrevUpdate: ctx.PrevUpdate, // no need to change prevUpdate since matching item can only be a map
						I:          ctx.I + 1,
						Obj:        typedObj[idx],
					})
				}
			}

		case KeyToken:
			typedObj, ok := ctx.Obj.(map[interface{}]interface{})
			if !ok {
				return nil, newOpMapMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			o, found := typedObj[typedToken.Key]
			if !found && !typedToken.Optional {
				return nil, opMissingMapKeyErr{typedToken.Key, NewPointer(tokens[:ctx.I+2]), typedObj}
			}

			if isLast {
				clonedValue, err := op.cloneValue(op.Value)
				if err != nil {
					return nil, replaceOpCloneValueErr(err)
				}
				typedObj[typedToken.Key] = clonedValue
			} else {
				if !found {
					// Determine what type of value to create based on next token
					switch tokens[ctx.I+2].(type) {
					case AfterLastIndexToken:
						o = []interface{}{}
					case WildcardToken:
						o = []interface{}{}
					case MatchingIndexToken:
						o = []interface{}{}
					case KeyToken:
						o = map[interface{}]interface{}{}
					default:
						errMsg := "Expected to find key, matching index or after last index token at path '%s'"
						return nil, fmt.Errorf(errMsg, NewPointer(tokens[:ctx.I+3]))
					}

					typedObj[typedToken.Key] = o
				}

				ctxStack = append(ctxStack, &mutationCtx{
					PrevUpdate: func(newObj interface{}) { typedObj[typedToken.Key] = newObj },
					I:          ctx.I + 1,
					Obj:        o,
				})
			}

		case WildcardToken:
			if isLast {
				return nil, fmt.Errorf("Wildcard must not be the last token", NewPointer(tokens[:ctx.I+2]))
			}

			typedObj, ok := ctx.Obj.([]interface{})
			if !ok {
				return nil, newOpArrayMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			for idx, o := range typedObj {
				ctxStack = append(ctxStack, &mutationCtx{
					PrevUpdate: func(newObj interface{}) { typedObj[idx] = newObj },
					I:          ctx.I + 1,
					Obj:        o,
				})
			}

		default:
			return nil, opUnexpectedTokenErr{token, NewPointer(tokens[:ctx.I+2])}
		}
	}

	return doc, nil
}

func (ReplaceOp) cloneValue(in interface{}) (out interface{}, err error) {
	defer func() {
		if recoverVal := recover(); recoverVal != nil {
			err = fmt.Errorf("Recovered: %s", recoverVal)
		}
	}()

	bytes, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
