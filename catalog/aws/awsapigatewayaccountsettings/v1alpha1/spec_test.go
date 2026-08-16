package awsapigatewayaccountsettingsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsApiGatewayAccountSettingsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsApiGatewayAccountSettingsSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalSettings is the smallest valid instance: region only - the
// explicit "no logging role" posture.
func minimalSettings() *AwsApiGatewayAccountSettingsSpec {
	return &AwsApiGatewayAccountSettingsSpec{
		Region: "us-west-2",
	}
}

var _ = ginkgo.Describe("AwsApiGatewayAccountSettingsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with region only (the explicit no-role posture)", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalSettings())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a CloudWatch role", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalSettings()
				spec.CloudwatchRoleArn = svr("arn:aws:iam::123456789012:role/apigw-cloudwatch")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalSettings()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
