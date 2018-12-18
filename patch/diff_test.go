package patch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/go-patch/patch"
)

var _ = Describe("Diff.Calculate", func() {
	testDiff := func(left, right interface{}, expectedOps []Op) {
		Expect(Diff{Left: left, Right: right}.Calculate()).To(Equal(Ops(expectedOps)))

		result, err := Ops(expectedOps).Apply(left)
		Expect(err).ToNot(HaveOccurred())

		if right == nil { // gomega does not allow nil==nil comparison
			Expect(result).To(BeNil())
		} else {
			Expect(result).To(Equal(right))
		}
	}

	It("returns no ops if both docs are same", func() {
		testDiff(nil, nil, []Op{})

		testDiff(
			map[interface{}]interface{}{"a": 124},
			map[interface{}]interface{}{"a": 124},
			[]Op{},
		)

		testDiff(
			[]interface{}{"a", 124},
			[]interface{}{"a", 124},
			[]Op{},
		)
	})

	It("can replace doc root", func() {
		testDiff(nil, "a", []Op{ReplaceOp{Path: MustNewPointerFromString(""), Value: "a"}})
	})

	It("can replace doc root with nil", func() {
		testDiff("a", nil, []Op{ReplaceOp{Path: MustNewPointerFromString(""), Value: nil}})
	})

	It("can diff maps", func() {
		testDiff(
			map[interface{}]interface{}{"a": 123},
			map[interface{}]interface{}{"a": 124},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/a"), Value: 124},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123, "b": 456},
			map[interface{}]interface{}{"a": 124, "c": 456},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/a"), Value: 124},
				RemoveOp{Path: MustNewPointerFromString("/b")},
				ReplaceOp{Path: MustNewPointerFromString("/c?"), Value: 456},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123, "b": 456},
			map[interface{}]interface{}{"a": 124},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/a"), Value: 124},
				RemoveOp{Path: MustNewPointerFromString("/b")},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123, "b": 456},
			map[interface{}]interface{}{},
			[]Op{
				RemoveOp{Path: MustNewPointerFromString("/a")},
				RemoveOp{Path: MustNewPointerFromString("/b")},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123},
			map[interface{}]interface{}{"a": nil},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/a"), Value: nil},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123, "b": map[interface{}]interface{}{"a": 1024, "b": 4056}},
			map[interface{}]interface{}{"a": 124, "b": map[interface{}]interface{}{"a": 1024, "c": 4056}},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/a"), Value: 124},
				RemoveOp{Path: MustNewPointerFromString("/b/b")},
				ReplaceOp{Path: MustNewPointerFromString("/b/c?"), Value: 4056},
			},
		)

		testDiff(
			map[interface{}]interface{}{"a": 123},
			"a",
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString(""), Value: "a"},
			},
		)

		testDiff(
			"a",
			map[interface{}]interface{}{"a": 123},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString(""), Value: map[interface{}]interface{}{"a": 123}},
			},
		)
	})

	It("can diff arrays", func() {
		testDiff(
			[]interface{}{"a", 123},
			[]interface{}{"b", 123},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/0"), Value: "b"},
			},
		)

		testDiff(
			[]interface{}{"a"},
			[]interface{}{"b", 123, 456},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/0"), Value: "b"},
				ReplaceOp{Path: MustNewPointerFromString("/-"), Value: 123},
				ReplaceOp{Path: MustNewPointerFromString("/-"), Value: 456},
			},
		)

		testDiff(
			[]interface{}{"a", 123, 456},
			[]interface{}{"b"},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/0"), Value: "b"},
				RemoveOp{Path: MustNewPointerFromString("/1")},
				RemoveOp{Path: MustNewPointerFromString("/1")},
			},
		)

		testDiff(
			[]interface{}{123, 456},
			[]interface{}{},
			[]Op{
				RemoveOp{Path: MustNewPointerFromString("/0")},
				RemoveOp{Path: MustNewPointerFromString("/0")},
			},
		)

		testDiff(
			[]interface{}{123, 456},
			[]interface{}{123, "a", 456}, // TODO unoptimized insertion
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/1"), Value: "a"},
				ReplaceOp{Path: MustNewPointerFromString("/-"), Value: 456},
			},
		)

		testDiff(
			[]interface{}{[]interface{}{456, 789}},
			[]interface{}{[]interface{}{789}},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString("/0/0"), Value: 789},
				RemoveOp{Path: MustNewPointerFromString("/0/1")},
			},
		)

		testDiff(
			[]interface{}{"a", 123},
			"a",
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString(""), Value: "a"},
			},
		)

		testDiff(
			"a",
			[]interface{}{"a", 123},
			[]Op{
				ReplaceOp{Path: MustNewPointerFromString(""), Value: []interface{}{"a", 123}},
			},
		)
	})
})
