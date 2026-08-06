package awssnssubscriptionv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsSnsSubscriptionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsSnsSubscriptionSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a Struct from a map, failing loudly on error.
func mustStruct(m map[string]interface{}) *structpb.Struct {
	s, err := structpb.NewStruct(m)
	if err != nil {
		panic(err)
	}
	return s
}

var _ = ginkgo.Describe("AwsSnsSubscriptionSpec validations", func() {
	var spec *AwsSnsSubscriptionSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: an SQS subscription.
		spec = &AwsSnsSubscriptionSpec{
			Region:   "us-west-2",
			TopicArn: strRef("arn:aws:sns:us-west-2:123456789012:order-events"),
			Protocol: "sqs",
			Endpoint: strRef("arn:aws:sqs:us-west-2:123456789012:fulfillment"),
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal SQS subscription", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts filtering, raw delivery, and a DLQ", func() {
		spec.FilterPolicy = mustStruct(map[string]interface{}{
			"event_type": []interface{}{"order_placed"},
		})
		spec.FilterPolicyScope = "MessageBody"
		spec.RawMessageDelivery = true
		spec.DeadLetterConfig = &AwsSnsSubscriptionDeadLetterConfig{
			DeadLetterTargetArn: strRef("arn:aws:sqs:us-west-2:123456789012:dlq"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a firehose subscription with a role", func() {
		spec.Protocol = "firehose"
		spec.Endpoint = strRef("arn:aws:firehose:us-west-2:123456789012:deliverystream/archive")
		spec.SubscriptionRoleArn = strRef("arn:aws:iam::123456789012:role/sns-firehose")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts an auto-confirming HTTPS subscription with a timeout", func() {
		spec.Protocol = "https"
		spec.Endpoint = strRef("https://hooks.example.com/sns")
		spec.EndpointAutoConfirms = true
		spec.ConfirmationTimeoutMinutes = 5
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a replay policy", func() {
		spec.ReplayPolicy = mustStruct(map[string]interface{}{
			"PointType": "Timestamp",
		})
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Required fields
	// -------------------------------------------------------------------------

	ginkgo.It("fails without a region", func() {
		spec.Region = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails without a topic_arn", func() {
		spec.TopicArn = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails without an endpoint", func() {
		spec.Endpoint = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// CEL: protocol and coupling rules
	// -------------------------------------------------------------------------

	ginkgo.It("fails on an unknown protocol", func() {
		spec.Protocol = "mqtt"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when filter_policy_scope is set without filter_policy", func() {
		spec.FilterPolicyScope = "MessageAttributes"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when filter_policy_scope has an invalid value", func() {
		spec.FilterPolicy = mustStruct(map[string]interface{}{"k": "v"})
		spec.FilterPolicyScope = "Everything"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when firehose protocol lacks subscription_role_arn", func() {
		spec.Protocol = "firehose"
		spec.Endpoint = strRef("arn:aws:firehose:us-west-2:123456789012:deliverystream/archive")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when endpoint_auto_confirms is set on a non-HTTP protocol", func() {
		spec.EndpointAutoConfirms = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when confirmation_timeout_minutes is set on a non-HTTP protocol", func() {
		spec.ConfirmationTimeoutMinutes = 5
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
