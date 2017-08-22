package patch

import "errors"

type RootToken struct{}

func (t RootToken) String() string {
	return ""
}

func (t RootToken) processDescent(ctx *tokenContext) (interface{}, error) {
	return nil, errors.New("not supported")
}
