package patch

import "errors"

type KeyToken struct {
	Key string

	Optional bool
}

func (t KeyToken) String() string {
	str := rfc6901Encoder.Replace(t.Key)

	if t.Optional { // /key?/key2/key3
		str += "?"
	}

	return str
}

func (t KeyToken) processDescent(ctx *tokenContext) (interface{}, error) {
	if ctx.Node == nil {
		ctx.Node = make(map[interface{}]interface{})
		if ctx.Method == methodReplace {
			ctx.Setter(ctx.Node)
		}
	}

	typedObj, ok := ctx.Node.(map[interface{}]interface{})
	if !ok {
		return nil, newOpMapMismatchTypeErr(ctx.Tokens[:ctx.TokenIndex+1], ctx.Node)
	}

	var found bool
	ctx.Node, found = typedObj[t.Key]
	if !found {
		if !t.Optional {
			return nil, opMissingMapKeyErr{t.Key, NewPointer(ctx.Tokens[:ctx.TokenIndex+1]), typedObj}
		}
		ctx.Node = nil // up to next to create thyself
	}

	if !ctx.IsLast() {
		ctx.Setter = func(newObj interface{}) { typedObj[t.Key] = newObj }
		return ctx.Descend()
	}

	switch ctx.Method {
	case methodFind:
		return typedObj[t.Key], nil

	case methodReplace:
		v, err := ctx.Value()
		if err != nil {
			return nil, err
		}
		typedObj[t.Key] = v
		return nil, nil

	case methodRemove:
		delete(typedObj, t.Key)
		return nil, nil

	default:
		return nil, errors.New("unsupported")
	}
}
