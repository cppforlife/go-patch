package patch

type CopyOp struct {
	Path Pointer
	From Pointer
}

func (op CopyOp) Apply(doc interface{}) (interface{}, error) {
	value, err := FindOp{Path: op.From}.Apply(doc)
	if err != nil {
		return nil, err
	}

	return ReplaceOp{Path: op.Path, Value: value}.Apply(doc)
}
