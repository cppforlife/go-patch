package patch

import "fmt"

type BeforeFirstIndexToken struct{}

func (t BeforeFirstIndexToken) String() string {
	return "+"
}

func (t BeforeFirstIndexToken) processDescent(ctx *tokenContext) (interface{}, error) {
	if ctx.Node == nil {
		ctx.Node = make([]interface{}, 0)
		if ctx.Method == methodReplace {
			ctx.Setter(ctx.Node)
		}
	}

	if ctx.Method != methodReplace {
		if ctx.Method == methodFind {
			errMsg := "Expected after last index token to be last in path '%s' (not supported in find operations)"
			return nil, fmt.Errorf(errMsg, NewPointer(ctx.Tokens))
		}
		return nil, opUnexpectedTokenErr{t, NewPointer(ctx.Tokens[:ctx.TokenIndex+1])}
	}

	typedObj, ok := ctx.Node.([]interface{})
	if !ok {
		return nil, newOpArrayMismatchTypeErr(ctx.Tokens[:ctx.TokenIndex+1], ctx.Node)
	}

	if !ctx.IsLast() {
		return nil, fmt.Errorf("Expected before first index token to be last in path '%s'", Pointer{tokens: ctx.Tokens})
	}

	v, err := ctx.Value()
	if err != nil {
		return nil, err
	}
	ctx.Setter(append([]interface{}{v}, typedObj...))
	return nil, nil
}
