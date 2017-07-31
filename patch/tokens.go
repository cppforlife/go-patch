package patch

const (
	methodFind    = 0
	methodReplace = 1
	methodRemove  = 2
)

type tokenContext struct {
	Tokens     []Token
	TokenIndex int

	Node interface{}

	Setter func(newObj interface{})
	Value  func() (interface{}, error)

	Method int
}

func (rc *tokenContext) IsLast() bool {
	return (rc.TokenIndex + 1) == len(rc.Tokens)
}

func (rc *tokenContext) Descend() (interface{}, error) {
	// Clone our context so values can be safely overridden
	nc := *rc
	nc.TokenIndex++

	return nc.Tokens[nc.TokenIndex].processDescent(&nc)
}

type Token interface {
	processDescent(ctx *tokenContext) (interface{}, error)
	String() string
}
