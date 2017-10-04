package patch

import (
	"fmt"
)

type ArrayInsertion struct {
	Index     int
	Modifiers []Modifier
	Array     []interface{}
}

type ArrayInsertionIndex struct {
	Number int
	Insert bool
}

func (i ArrayInsertion) Concrete() (ArrayInsertionIndex, error) {
	var mods []Modifier

	before := false
	after := false

	for _, modifier := range i.Modifiers {
		if before {
			return ArrayInsertionIndex{}, fmt.Errorf(
				"Expected to not find any modifiers after 'before' modifier, but found modifier '%T'", modifier)
		}
		if after {
			return ArrayInsertionIndex{}, fmt.Errorf(
				"Expected to not find any modifiers after 'after' modifier, but found modifier '%T'", modifier)
		}

		switch modifier.(type) {
		case BeforeModifier:
			before = true
		case AfterModifier:
			after = true
		default:
			mods = append(mods, modifier)
		}
	}

	idx := ArrayIndex{Index: i.Index, Modifiers: mods, Array: i.Array}

	num, err := idx.Concrete()
	if err != nil {
		return ArrayInsertionIndex{}, err
	}

	if before {
		num -= 1
		if num < 0 {
			num = 0
		}
	}

	if after {
		num += 1
	}

	return ArrayInsertionIndex{num, before || after}, nil
}
