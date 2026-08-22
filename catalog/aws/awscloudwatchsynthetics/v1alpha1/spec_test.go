package awscloudwatchsyntheticsv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsCloudwatchSyntheticsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchSyntheticsSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalCanary() *AwsSyntheticsCanary {
	return &AwsSyntheticsCanary{
		ArtifactBucket:   svr("my-canary-artifacts"),
		ExecutionRoleArn: svr("arn:aws:iam::123456789012:role/canary-exec"),
		Handler:          "index.handler",
		RuntimeVersion:   "syn-nodejs-puppeteer-9.1",
		Code: &AwsSyntheticsCanaryCode{
			S3Bucket: svr("my-canary-code"),
			S3Key:    "e2e/canary.zip",
		},
		Schedule: &AwsSyntheticsCanarySchedule{
			Expression: "rate(5 minutes)",
		},
	}
}

func canarySpec() *AwsCloudwatchSyntheticsSpec {
	return &AwsCloudwatchSyntheticsSpec{
		Region: "us-east-1",
		Canary: minimalCanary(),
	}
}

var _ = ginkgo.Describe("AwsCloudwatchSyntheticsSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal canary", func() {
			gomega.Expect(protovalidate.Validate(canarySpec())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a groups-only instance", func() {
			spec := &AwsCloudwatchSyntheticsSpec{
				Region: "us-east-1",
				Groups: []*AwsSyntheticsGroup{{Name: "checkout-canaries"}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a canary joining groups", func() {
			spec := canarySpec()
			spec.Groups = []*AwsSyntheticsGroup{{Name: "checkout-canaries"}}
			spec.GroupNames = []string{"checkout-canaries", "org-wide"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cron schedule and run config at the boundaries", func() {
			spec := canarySpec()
			spec.Canary.Schedule.Expression = "cron(0 8 * * ? *)"
			spec.Canary.Schedule.MaxRetries = proto.Int32(0)
			spec.Canary.RunConfig = &AwsSyntheticsCanaryRunConfig{
				MemoryInMb:       proto.Int32(960),
				EphemeralStorage: proto.Int32(1024),
				TimeoutInSeconds: proto.Int32(840),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts SSE_KMS artifact encryption with a key", func() {
			spec := canarySpec()
			spec.Canary.ArtifactEncryptionMode = "SSE_KMS"
			spec.Canary.ArtifactEncryptionKmsKeyArn = svr("arn:aws:kms:us-east-1:123456789012:key/abc")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an empty spec (no arm)", func() {
			spec := &AwsCloudwatchSyntheticsSpec{Region: "us-east-1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects group joins without the canary arm", func() {
			spec := &AwsCloudwatchSyntheticsSpec{
				Region:     "us-east-1",
				Groups:     []*AwsSyntheticsGroup{{Name: "checkout-canaries"}},
				GroupNames: []string{"checkout-canaries"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate owned group names", func() {
			spec := canarySpec()
			spec.Groups = []*AwsSyntheticsGroup{{Name: "dup"}, {Name: "dup"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects memory off Lambda's 64MB granularity or below 960", func() {
			for _, mem := range []int32{950, 970} {
				spec := canarySpec()
				spec.Canary.RunConfig = &AwsSyntheticsCanaryRunConfig{MemoryInMb: proto.Int32(mem)}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			}
		})

		ginkgo.It("rejects a KMS key without SSE_KMS mode", func() {
			spec := canarySpec()
			spec.Canary.ArtifactEncryptionKmsKeyArn = svr("arn:aws:kms:us-east-1:123456789012:key/abc")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a schedule that is neither rate nor cron", func() {
			spec := canarySpec()
			spec.Canary.Schedule.Expression = "every 5 minutes"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects retry counts above AWS's cap", func() {
			spec := canarySpec()
			spec.Canary.Schedule.MaxRetries = proto.Int32(3)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects retention outside 1-455 days", func() {
			spec := canarySpec()
			spec.Canary.FailureRetentionPeriod = proto.Int32(456)
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-Synthetics runtime string", func() {
			spec := canarySpec()
			spec.Canary.RuntimeVersion = "nodejs18.x"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
