package awsbedrockinvocationloggingv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsBedrockInvocationLoggingSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBedrockInvocationLoggingSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalLogging is the smallest valid instance: region + the S3
// delivery arm.
func minimalLogging() *AwsBedrockInvocationLoggingSpec {
	return &AwsBedrockInvocationLoggingSpec{
		Region: "us-west-2",
		S3: &AwsBedrockInvocationLoggingS3{
			BucketName: svr("bedrock-invocation-logs"),
		},
	}
}

// cloudwatchLogging is the CloudWatch arm with the S3 spillover.
func cloudwatchLogging() *AwsBedrockInvocationLoggingSpec {
	return &AwsBedrockInvocationLoggingSpec{
		Region: "us-west-2",
		Cloudwatch: &AwsBedrockInvocationLoggingCloudwatch{
			LogGroupName: svr("/bedrock/invocations"),
			RoleArn:      svr("arn:aws:iam::123456789012:role/bedrock-logging"),
			LargeDataDeliveryS3: &AwsBedrockInvocationLoggingS3{
				BucketName: svr("bedrock-large-payloads"),
				KeyPrefix:  "spillover",
			},
		},
	}
}

var _ = ginkgo.Describe("AwsBedrockInvocationLoggingSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with the S3 arm only", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalLogging())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the CloudWatch arm and its spillover", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(cloudwatchLogging())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with both destinations and explicit data-type opt-outs", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := cloudwatchLogging()
				spec.S3 = &AwsBedrockInvocationLoggingS3{BucketName: svr("bedrock-logs")}
				spec.VideoDataDeliveryEnabled = proto.Bool(false)
				spec.EmbeddingDataDeliveryEnabled = proto.Bool(false)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a configuration with no delivery destination", func() {
			spec := minimalLogging()
			spec.S3 = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("delivery destination"))
		})

		ginkgo.It("rejects a CloudWatch arm without its role", func() {
			spec := cloudwatchLogging()
			spec.Cloudwatch.RoleArn = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a CloudWatch arm without its log group", func() {
			spec := cloudwatchLogging()
			spec.Cloudwatch.LogGroupName = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an S3 arm without its bucket", func() {
			spec := minimalLogging()
			spec.S3.BucketName = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing region", func() {
			spec := minimalLogging()
			spec.Region = ""
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
