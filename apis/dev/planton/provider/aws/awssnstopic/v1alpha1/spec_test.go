package awssnstopicv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsSnsTopicSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSnsTopicSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a Struct from a map, failing the spec on error.
func mustStruct(m map[string]interface{}) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return s
}

var _ = ginkgo.Describe("AwsSnsTopicSpec validations", func() {
	var spec *AwsSnsTopicSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: a standard topic with all AWS defaults.
		spec = &AwsSnsTopicSpec{
			Region: "us-west-2",
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal standard topic (all defaults)", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a FIFO topic with content-based deduplication and throughput scope", func() {
		spec.FifoTopic = true
		spec.ContentBasedDeduplication = true
		spec.FifoThroughputScope = "MessageGroup"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a FIFO topic with an archive policy", func() {
		spec.FifoTopic = true
		spec.ArchivePolicy = mustStruct(map[string]interface{}{
			"MessageRetentionPeriod": 30,
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a standard topic with KMS encryption and access policy", func() {
		spec.KmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/mrk-abc123")
		spec.Policy = mustStruct(map[string]interface{}{
			"Version": "2012-10-17",
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a standard topic with a data protection policy", func() {
		spec.DataProtectionPolicy = mustStruct(map[string]interface{}{
			"Name":    "pii-guard",
			"Version": "2021-06-01",
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts delivery feedback with roles and sample rate", func() {
		spec.DeliveryFeedback = &AwsSnsTopicDeliveryFeedback{
			Sqs: &AwsSnsTopicProtocolFeedback{
				SuccessFeedbackRole:       strRef("arn:aws:iam::123456789012:role/sns-feedback"),
				FailureFeedbackRole:       strRef("arn:aws:iam::123456789012:role/sns-feedback"),
				SuccessFeedbackSampleRate: 50,
			},
			Lambda: &AwsSnsTopicProtocolFeedback{
				FailureFeedbackRole: strRef("arn:aws:iam::123456789012:role/sns-feedback"),
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts observability settings", func() {
		spec.TracingConfig = "Active"
		spec.SignatureVersion = 2
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: FIFO-only fields on Standard topics
	// -------------------------------------------------------------------------

	ginkgo.It("fails when content_based_deduplication is set on a standard topic", func() {
		spec.FifoTopic = false
		spec.ContentBasedDeduplication = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when fifo_throughput_scope is set on a standard topic", func() {
		spec.FifoTopic = false
		spec.FifoThroughputScope = "Topic"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when fifo_throughput_scope has an invalid value", func() {
		spec.FifoTopic = true
		spec.FifoThroughputScope = "invalid"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when archive_policy is set on a standard topic", func() {
		spec.FifoTopic = false
		spec.ArchivePolicy = mustStruct(map[string]interface{}{
			"MessageRetentionPeriod": 30,
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: data protection policy is standard-only
	// -------------------------------------------------------------------------

	ginkgo.It("fails when data_protection_policy is set on a FIFO topic", func() {
		spec.FifoTopic = true
		spec.DataProtectionPolicy = mustStruct(map[string]interface{}{
			"Name": "pii-guard",
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: observability values
	// -------------------------------------------------------------------------

	ginkgo.It("fails when signature_version is invalid", func() {
		spec.SignatureVersion = 3
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when tracing_config is invalid", func() {
		spec.TracingConfig = "Sampled"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Delivery feedback validations
	// -------------------------------------------------------------------------

	ginkgo.It("fails when success_feedback_sample_rate exceeds 100", func() {
		spec.DeliveryFeedback = &AwsSnsTopicDeliveryFeedback{
			Http: &AwsSnsTopicProtocolFeedback{
				SuccessFeedbackRole:       strRef("arn:aws:iam::123456789012:role/sns-feedback"),
				SuccessFeedbackSampleRate: 101,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a sample rate is set without a success role", func() {
		spec.DeliveryFeedback = &AwsSnsTopicDeliveryFeedback{
			Firehose: &AwsSnsTopicProtocolFeedback{
				SuccessFeedbackSampleRate: 25,
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("accepts a failure-only feedback block", func() {
		spec.DeliveryFeedback = &AwsSnsTopicDeliveryFeedback{
			Application: &AwsSnsTopicProtocolFeedback{
				FailureFeedbackRole: strRef("arn:aws:iam::123456789012:role/sns-feedback"),
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})
})
