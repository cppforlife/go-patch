package patch

import (
	"errors"
	"fmt"
)

type IndexToken struct {
	Index int
}

func (t IndexToken) String() string {
	return fmt.Sprintf("%d", t.Index)
}

func (t IndexToken) processDescent(ctx *tokenContext) (interface{}, error) {
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

	if t.Index >= len(typedObj) {
		return nil, opMissingIndexErr{NewPointer(ctx.Tokens[:ctx.TokenIndex+1]), t.Index, typedObj}
	}

	if !ctx.IsLast() {
		ctx.Node = typedObj[t.Index]
		ctx.Setter = func(newObj interface{}) { typedObj[t.Index] = newObj }
		return ctx.Descend()
	}

	switch ctx.Method {
	case methodFind:
		return typedObj[t.Index], nil

	case methodReplace:
		v, err := ctx.Value()
		if err != nil {
			return nil, err
		}
		typedObj[t.Index] = v
		return nil, nil

	case methodRemove:
		ctx.Setter(append(append(make([]interface{}, 0, len(typedObj)-1), typedObj[:t.Index]...), typedObj[t.Index+1:]...))
		return nil, nil

	default:
		return nil, errors.New("unsupported")
	}
}
