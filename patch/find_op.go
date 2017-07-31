package patch

type FindOp struct {
	Path Pointer
}

func (op FindOp) Apply(doc interface{}) (interface{}, error) {
	tokens := op.Path.Tokens()

	if len(tokens) == 1 {
		return doc, nil
	}

	return (&tokenContext{
		Tokens:     tokens,
		TokenIndex: 0,
		Node:       doc,
		Method:     methodFind,
	}).Descend()
}
