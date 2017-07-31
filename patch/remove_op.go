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

	_, err := (&tokenContext{
		Tokens:     tokens,
		TokenIndex: 0,
		Node:       doc,
		Setter:     func(newObj interface{}) { doc = newObj },
		Method:     methodRemove,
	}).Descend()
	if err != nil {
		return nil, err
	}

	return doc, nil
}
