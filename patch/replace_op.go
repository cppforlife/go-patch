package patch

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type ReplaceOp struct {
	Path  Pointer
	Value interface{} // will be cloned using yaml library
}

func (op ReplaceOp) Apply(doc interface{}) (interface{}, error) {
	// Ensure that value is not modified by future operations
	clonedValue, err := op.cloneValue(op.Value)
	if err != nil {
		return nil, fmt.Errorf("ReplaceOp cloning value: %s", err)
	}

	tokens := op.Path.Tokens()
	if len(tokens) == 1 {
		return clonedValue, nil
	}

	_, err = (&tokenContext{
		Tokens:     tokens,
		TokenIndex: 0,
		Node:       doc,
		Setter:     func(newObj interface{}) { doc = newObj },
		Value:      func() (interface{}, error) { return op.cloneValue(clonedValue) },
		Method:     methodReplace,
	}).Descend()
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func (ReplaceOp) cloneValue(in interface{}) (out interface{}, err error) {
	defer func() {
		if recoverVal := recover(); recoverVal != nil {
			err = fmt.Errorf("Recovered: %s", recoverVal)
		}
	}()

	bytes, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
