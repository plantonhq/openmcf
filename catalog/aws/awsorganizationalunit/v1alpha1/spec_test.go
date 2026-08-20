package awsorganizationalunitv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsOrganizationalUnitSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsOrganizationalUnitSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func ref(name string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
			ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
		},
	}
}

// minimalOu is the smallest valid instance: a first-level OU under a
// literal root.
func minimalOu() *AwsOrganizationalUnitSpec {
	return &AwsOrganizationalUnitSpec{
		Region:   "us-west-2",
		OuName:   "Core Services",
		ParentId: svr("r-abcd"),
	}
}

var _ = ginkgo.Describe("AwsOrganizationalUnitSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a first-level OU under a literal root", func() {
			gomega.Expect(protovalidate.Validate(minimalOu())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a nested OU under a literal parent OU", func() {
			spec := minimalOu()
			spec.ParentId = svr("ou-abcd-12345678")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a parent by reference", func() {
			spec := minimalOu()
			spec.ParentId = ref("my-organization")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an OU name with spaces (the explicit-name-field reason)", func() {
			spec := minimalOu()
			spec.OuName = "Production Workloads / Tier 1"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalOu()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing OU name", func() {
			spec := minimalOu()
			spec.OuName = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an OU name above 128 characters", func() {
			spec := minimalOu()
			name := ""
			for i := 0; i < 129; i++ {
				name += "a"
			}
			spec.OuName = name
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing parent", func() {
			spec := minimalOu()
			spec.ParentId = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a literal parent that is neither a root nor an OU", func() {
			spec := minimalOu()
			spec.ParentId = svr("111111111111")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
