package awshttpapivpclinkv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsHttpApiVpcLinkSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsHttpApiVpcLinkSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

var _ = ginkgo.Describe("AwsHttpApiVpcLinkSpec validations", func() {
	var spec *AwsHttpApiVpcLinkSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: region + one subnet.
		spec = &AwsHttpApiVpcLinkSpec{
			Region: "us-west-2",
			SubnetIds: []*foreignkeyv1.StringValueOrRef{
				strRef("subnet-0abc123"),
			},
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal spec (region + one subnet)", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a two-AZ subnet spread with security groups", func() {
		spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{
			strRef("subnet-0abc123"),
			strRef("subnet-0def456"),
		}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{
			strRef("sg-0abc123"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a spec without security groups (AWS applies none)", func() {
		spec.SecurityGroupIds = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required field validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when subnet_ids is empty", func() {
		spec.SubnetIds = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
