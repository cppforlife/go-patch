package patch

import "errors"

type WildcardToken struct{}

func (t WildcardToken) String() string {
	return "*"
}

func (t WildcardToken) processDescent(ctx *tokenContext) (interface{}, error) {
	if ctx.IsLast() {
		return nil, errors.New("wildcard can't be used for last element. Use - instead")
	}
	if ctx.Method != methodReplace && ctx.Method != methodRemove {
		return nil, errors.New("operation using wildcard")
	}

	typedArray, ok := ctx.Node.([]interface{})
	if !ok {
		return nil, newOpArrayMismatchTypeErr(ctx.Tokens[:ctx.TokenIndex+1], ctx.Node)
	}
	for idx, e := range typedArray {
		ctx.Node = e
		ctx.Setter = func(newObj interface{}) { typedArray[idx] = newObj }
		_, err := ctx.Descend()
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
