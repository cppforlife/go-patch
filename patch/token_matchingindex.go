package patch

import (
	"errors"
	"fmt"
)

type MatchingIndexToken struct {
	Key   string
	Value string

	Optional bool
}

func (t MatchingIndexToken) String() string {
	key := rfc6901Encoder.Replace(t.Key)
	val := rfc6901Encoder.Replace(t.Value)

	if t.Optional {
		val += "?"
	}

	return fmt.Sprintf("%s=%s", key, val)
}

func (t MatchingIndexToken) processDescent(ctx *tokenContext) (interface{}, error) {
	if ctx.Node == nil {
		ctx.Node = make([]interface{}, 0)
		if ctx.Method == methodReplace {
			ctx.Setter(ctx.Node)
		}
	}

	typedObj, ok := ctx.Node.([]interface{})
	if !ok {
		return nil, newOpArrayMismatchTypeErr(ctx.Tokens[:ctx.TokenIndex+1], ctx.Node)
	}

	var idxs []int
	for itemIdx, item := range typedObj {
		typedItem, ok := item.(map[interface{}]interface{})
		if ok {
			if typedItem[t.Key] == t.Value {
				idxs = append(idxs, itemIdx)
			}
		}
	}

	switch len(idxs) {
	case 0:
		if !t.Optional {
			return nil, fmt.Errorf("Expected to find exactly one matching array item for path '%s' but found 0", NewPointer(ctx.Tokens[:ctx.TokenIndex+1]))
		}
		// We know the type here - it must be a map (else we couldn't key match it)
		idxs = []int{len(typedObj)}
		typedObj = append(append(make([]interface{}, 0, len(typedObj)+1), typedObj...), map[interface{}]interface{}{
			t.Key: t.Value,
		})
		if ctx.Method == methodReplace {
			ctx.Setter(typedObj)
		}
	case 1:
		// good, proceed as normal
	default:
		return nil, opMultipleMatchingIndexErr{NewPointer(ctx.Tokens[:ctx.TokenIndex+1]), idxs}
	}

	if !ctx.IsLast() {
		ctx.Node = typedObj[idxs[0]]
		ctx.Setter = func(newObj interface{}) { typedObj[idxs[0]] = newObj }
		return ctx.Descend()
	}

	switch ctx.Method {
	case methodFind:
		return typedObj[idxs[0]], nil

	case methodReplace:
		v, err := ctx.Value()
		if err != nil {
			return nil, err
		}
		typedObj[idxs[0]] = v
		return nil, nil

	case methodRemove:
		ctx.Setter(append(append(make([]interface{}, 0, len(typedObj)-1), typedObj[:idxs[0]]...), typedObj[idxs[0]+1:]...))
		return nil, nil

	default:
		return nil, errors.New("unsupported")
	}
}
