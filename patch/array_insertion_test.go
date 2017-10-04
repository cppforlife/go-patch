package patch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/go-patch/patch"
)

var _ = Describe("ArrayInsertion", func() {
	Describe("Concrete", func() {
		It("returns specified index when not using any modifiers", func() {
			idx := ArrayInsertion{Index: 1, Modifiers: []Modifier{}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{1, false}))
		})

		It("returns index adjusted for previous and next modifiers", func() {
			p := PrevModifier{}
			n := NextModifier{}

			idx := ArrayInsertion{Index: 1, Modifiers: []Modifier{p, n, n}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{2, false}))
		})

		It("returns error if both after and before are used", func() {
			idx := ArrayInsertion{Index: 0, Modifiers: []Modifier{BeforeModifier{}, AfterModifier{}}, Array: []interface{}{}}
			_, err := idx.Concrete()
			Expect(err.Error()).To(Equal("Expected to not find any modifiers after 'before' modifier, but found modifier 'patch.AfterModifier'"))

			idx = ArrayInsertion{Index: 0, Modifiers: []Modifier{AfterModifier{}, BeforeModifier{}}, Array: []interface{}{}}
			_, err = idx.Concrete()
			Expect(err.Error()).To(Equal("Expected to not find any modifiers after 'after' modifier, but found modifier 'patch.BeforeModifier'"))

			idx = ArrayInsertion{Index: 0, Modifiers: []Modifier{AfterModifier{}, PrevModifier{}}, Array: []interface{}{}}
			_, err = idx.Concrete()
			Expect(err.Error()).To(Equal("Expected to not find any modifiers after 'after' modifier, but found modifier 'patch.PrevModifier'"))
		})

		It("returns (0, true) when inserting in the beginning", func() {
			idx := ArrayInsertion{Index: 0, Modifiers: []Modifier{BeforeModifier{}}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{0, true}))
		})

		It("returns (last+1, true) when inserting in the end", func() {
			idx := ArrayInsertion{Index: 2, Modifiers: []Modifier{AfterModifier{}}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{3, true}))

			idx = ArrayInsertion{Index: -1, Modifiers: []Modifier{AfterModifier{}}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{3, true}))
		})

		It("returns (mid+1, true) when inserting in the middle", func() {
			idx := ArrayInsertion{Index: 1, Modifiers: []Modifier{AfterModifier{}}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{2, true}))
		})

		It("returns index adjusted for previous, next modifiers and before modifier", func() {
			p := PrevModifier{}
			n := NextModifier{}
			b := BeforeModifier{}

			idx := ArrayInsertion{Index: 1, Modifiers: []Modifier{p, n, n, b}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{1, true}))
		})

		It("returns index adjusted for previous, next modifiers and after modifier", func() {
			p := PrevModifier{}
			n := NextModifier{}
			a := AfterModifier{}

			idx := ArrayInsertion{Index: 1, Modifiers: []Modifier{p, n, n, a}, Array: []interface{}{1, 2, 3}}
			Expect(idx.Concrete()).To(Equal(ArrayInsertionIndex{3, true}))
		})
	})
})
