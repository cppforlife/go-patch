package patch

import (
	"fmt"
	"reflect"
)

type TestOp struct {
	Path  Pointer
	Value interface{}
}

func (op TestOp) Apply(doc interface{}) (interface{}, error) {
	foundVal, err := FindOp{Path: op.Path}.Apply(doc)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(foundVal, op.Value) {
		return nil, fmt.Errorf("Found value does not match expected value")
	}

	// Return same input document
	return doc, nil
}
