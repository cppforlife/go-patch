package patch

import (
	"fmt"
)

type RemoveOp struct {
	Path Pointer
}

func (op RemoveOp) Apply(doc interface{}) (interface{}, error) {
	tokens := op.Path.Tokens()

	if len(tokens) == 1 {
		return nil, fmt.Errorf("Cannot remove entire document")
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
				var newAry []interface{}
				newAry = append(newAry, typedObj[:idx]...)
				newAry = append(newAry, typedObj[idx+1:]...)
				ctx.PrevUpdate(newAry)
			} else {
				ctxStack = append(ctxStack, &mutationCtx{
					Obj:        typedObj[idx],
					PrevUpdate: func(newObj interface{}) { typedObj[idx] = newObj },
					I:          ctx.I + 1,
				})
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
				continue // don't exit early
			}

			if len(idxs) != 1 {
				return nil, opMultipleMatchingIndexErr{NewPointer(tokens[:ctx.I+2]), idxs}
			}

			idx := idxs[0]

			if isLast {
				var newAry []interface{}
				newAry = append(newAry, typedObj[:idx]...)
				newAry = append(newAry, typedObj[idx+1:]...)
				ctx.PrevUpdate(newAry)
			} else {
				ctxStack = append(ctxStack, &mutationCtx{
					Obj:        typedObj[idx],
					PrevUpdate: ctx.PrevUpdate, // no need to change prevUpdate since matching item can only be a map
					I:          ctx.I + 1,
				})
			}

		case KeyToken:
			typedObj, ok := ctx.Obj.(map[interface{}]interface{})
			if !ok {
				return nil, newOpMapMismatchTypeErr(tokens[:ctx.I+2], ctx.Obj)
			}

			o, found := typedObj[typedToken.Key]
			if !found {
				if typedToken.Optional {
					continue // don't return yet, as it may be present down alternate paths
				}

				return nil, opMissingMapKeyErr{typedToken.Key, NewPointer(tokens[:ctx.I+2]), typedObj}
			}

			if isLast {
				delete(typedObj, typedToken.Key)
			} else {
				ctxStack = append(ctxStack, &mutationCtx{
					Obj:        o,
					PrevUpdate: func(newObj interface{}) { typedObj[typedToken.Key] = newObj },
					I:          ctx.I + 1,
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
