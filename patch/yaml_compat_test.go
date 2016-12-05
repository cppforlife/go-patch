package patch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"

	. "github.com/cppforlife/go-patch/patch"
)

var _ = Describe("YAML compatibility", func() {
	Describe("empty string", func() {
		It("[WORKAROUND] works serializing empty strings", func() {
			str := `
- type: replace
  path: /instance_groups/name=cloud_controller/instances
  value: !!str ""
`

			var opDefs []OpDefinition

			err := yaml.Unmarshal([]byte(str), &opDefs)
			Expect(err).ToNot(HaveOccurred())

			val := opDefs[0].Value
			Expect((*val).(string)).To(Equal(""))
		})

		It("[PORBLEM] does not works serializing empty strings", func() {
			str := `
- type: replace
  path: /instance_groups/name=cloud_controller/instances
  value: ""
`

			var opDefs []OpDefinition

			err := yaml.Unmarshal([]byte(str), &opDefs)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot unmarshal !!str"))
		})
	})
})
