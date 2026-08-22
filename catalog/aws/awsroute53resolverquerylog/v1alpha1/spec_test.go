package awsroute53resolverquerylogv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsRoute53ResolverQueryLogSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRoute53ResolverQueryLogSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalConfig() *AwsRoute53ResolverQueryLogSpec {
	return &AwsRoute53ResolverQueryLogSpec{
		Region:         "us-west-2",
		DestinationArn: literal("arn:aws:logs:us-west-2:111111111111:log-group:resolver-queries"),
	}
}

var _ = ginkgo.Describe("AwsRoute53ResolverQueryLogSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal configuration", func() {
			gomega.Expect(protovalidate.Validate(minimalConfig())).To(gomega.BeNil())
		})

		ginkgo.It("accepts VPC associations", func() {
			spec := minimalConfig()
			spec.VpcIds = []*foreignkeyv1.StringValueOrRef{
				literal("vpc-0123456789abcdef0"),
				literal("vpc-0123456789abcdef1"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an S3 destination", func() {
			spec := minimalConfig()
			spec.DestinationArn = literal("arn:aws:s3:::resolver-query-logs/prefix")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing destination", func() {
			spec := minimalConfig()
			spec.DestinationArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing region", func() {
			spec := minimalConfig()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
