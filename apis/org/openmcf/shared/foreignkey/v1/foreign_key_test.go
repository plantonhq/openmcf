// filename: foreign_key_test.go
package foreignkeyv1

import (
	"testing"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"

	"buf.build/go/protovalidate"
)

func TestForeignKey(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "ForeignKey Suite")
}

var _ = ginkgo.Describe("ForeignKey Oneof Tests", func() {

	ginkgo.Describe("StringValueOrRef usage", func() {
		ginkgo.Context("when setting a literal value only", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &StringValueOrRef{
					LiteralOrRef: &StringValueOrRef_Value{
						Value: "my-string",
					},
				}
				// No custom rule; we expect no errors from protovalidate
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when setting a ValueFromRef only", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &StringValueOrRef{
					LiteralOrRef: &StringValueOrRef_ValueFrom{
						ValueFrom: &ValueFromRef{
							Kind:      cloudresourcekind.CloudResourceKind_TestCloudResourceOne,
							Env:       "dev",
							Name:      "my-cert",
							FieldPath: "status.outputs.cert_arn",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when setting both (proto oneof overwrites the first)", func() {
			ginkgo.It("should end up with the last field set and not produce a validation error", func() {
				// The 'value' gets overwritten by the oneof assignment to 'value_from'
				input := &StringValueOrRef{
					LiteralOrRef: &StringValueOrRef_Value{
						Value: "my-string",
					},
				}
				input.LiteralOrRef = &StringValueOrRef_ValueFrom{
					ValueFrom: &ValueFromRef{
						Kind: cloudresourcekind.CloudResourceKind_TestCloudResourceOne,
						Env:  "dev",
						Name: "overwrites-literal",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())

				// Because this is a oneof, the 'Value' is no longer set
				gomega.Expect(input.GetValue()).To(gomega.Equal(""))
				gomega.Expect(input.GetValueFrom().GetName()).To(gomega.Equal("overwrites-literal"))
			})
		})
	})

	ginkgo.Describe("Int32ValueOrRef usage", func() {
		ginkgo.Context("when setting an int32 literal only", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &Int32ValueOrRef{
					LiteralOrRef: &Int32ValueOrRef_Value{
						Value: 123,
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when setting a ValueFromRef only", func() {
			ginkgo.It("should not return a validation error", func() {
				input := &Int32ValueOrRef{
					LiteralOrRef: &Int32ValueOrRef_ValueFrom{
						ValueFrom: &ValueFromRef{
							Kind: cloudresourcekind.CloudResourceKind_TestCloudResourceOne,
							Env:  "dev",
							Name: "ref-int32",
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("when setting both fields in code sequentially", func() {
			ginkgo.It("should overwrite the first field with the second (no error)", func() {
				input := &Int32ValueOrRef{
					LiteralOrRef: &Int32ValueOrRef_Value{
						Value: 456,
					},
				}

				// Overwrite with ValueFromRef
				input.LiteralOrRef = &Int32ValueOrRef_ValueFrom{
					ValueFrom: &ValueFromRef{
						Kind: cloudresourcekind.CloudResourceKind_TestCloudResourceOne,
						Env:  "dev",
						Name: "ref-overwrites-int",
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())

				// Because it's a oneof, the integer literal is no longer set
				gomega.Expect(input.GetValue()).To(gomega.BeEquivalentTo(0))
				gomega.Expect(input.GetValueFrom().GetName()).To(gomega.Equal("ref-overwrites-int"))
			})
		})
	})
})
