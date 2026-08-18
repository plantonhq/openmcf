package awscostanomalymonitorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsCostAnomalyMonitorSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCostAnomalyMonitorSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalMonitor is the smallest valid instance: the recommended
// by-service dimensional monitor.
func minimalMonitor() *AwsCostAnomalyMonitorSpec {
	return &AwsCostAnomalyMonitorSpec{
		Region:           "us-east-1",
		MonitorName:      "Service Spend Monitor",
		MonitorType:      "DIMENSIONAL",
		MonitorDimension: "SERVICE",
	}
}

func validSubscription() *AwsCostAnomalyMonitorSubscription {
	return &AwsCostAnomalyMonitorSubscription{
		Name:      "finops-daily",
		Frequency: "DAILY",
		Subscribers: []*AwsCostAnomalyMonitorSubscriber{{
			Address: svr("finops@example.com"),
			Type:    "EMAIL",
		}},
	}
}

var _ = ginkgo.Describe("AwsCostAnomalyMonitorSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal dimensional monitor", func() {
			gomega.Expect(protovalidate.Validate(minimalMonitor())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a custom monitor with a specification", func() {
			specification, err := structpb.NewStruct(map[string]any{
				"Dimensions": map[string]any{
					"Key":    "LINKED_ACCOUNT",
					"Values": []any{"123456789012"},
				},
			})
			gomega.Expect(err).To(gomega.BeNil())
			spec := minimalMonitor()
			spec.MonitorType = "CUSTOM"
			spec.MonitorDimension = ""
			spec.MonitorSpecification = specification
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a daily email subscription", func() {
			spec := minimalMonitor()
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{validSubscription()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an immediate SNS subscription with a threshold", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.Frequency = "IMMEDIATE"
			sub.Subscribers = []*AwsCostAnomalyMonitorSubscriber{{
				Address: svr("arn:aws:sns:us-east-1:123456789012:cost-alerts"),
				Type:    "SNS",
			}}
			sub.ThresholdExpression = &AwsCostAnomalyMonitorExpression{
				Dimension: &AwsCostAnomalyMonitorExpressionDimension{
					Key:          "ANOMALY_TOTAL_IMPACT_ABSOLUTE",
					MatchOptions: []string{"GREATER_THAN_OR_EQUAL"},
					Values:       []string{"100"},
				},
			}
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a composed AND threshold", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.ThresholdExpression = &AwsCostAnomalyMonitorExpression{
				And: []*AwsCostAnomalyMonitorExpressionLeaf{
					{Dimension: &AwsCostAnomalyMonitorExpressionDimension{
						Key:          "ANOMALY_TOTAL_IMPACT_ABSOLUTE",
						MatchOptions: []string{"GREATER_THAN_OR_EQUAL"},
						Values:       []string{"100"},
					}},
					{Dimension: &AwsCostAnomalyMonitorExpressionDimension{
						Key:          "ANOMALY_TOTAL_IMPACT_PERCENTAGE",
						MatchOptions: []string{"GREATER_THAN_OR_EQUAL"},
						Values:       []string{"10"},
					}},
				},
			}
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalMonitor()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a dimensional monitor without its dimension", func() {
			spec := minimalMonitor()
			spec.MonitorDimension = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a custom monitor without its specification", func() {
			spec := minimalMonitor()
			spec.MonitorType = "CUSTOM"
			spec.MonitorDimension = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a dimensional monitor carrying a specification", func() {
			specification, err := structpb.NewStruct(map[string]any{"Dimensions": map[string]any{}})
			gomega.Expect(err).To(gomega.BeNil())
			spec := minimalMonitor()
			spec.MonitorSpecification = specification
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown monitor dimension", func() {
			spec := minimalMonitor()
			spec.MonitorDimension = "USAGE_TYPE"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate subscription names", func() {
			spec := minimalMonitor()
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{validSubscription(), validSubscription()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a subscriber-less subscription", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.Subscribers = nil
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an immediate subscription with an email subscriber", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.Frequency = "IMMEDIATE"
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a daily subscription with an SNS subscriber", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.Subscribers = []*AwsCostAnomalyMonitorSubscriber{{
				Address: svr("arn:aws:sns:us-east-1:123456789012:cost-alerts"),
				Type:    "SNS",
			}}
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown threshold dimension key", func() {
			spec := minimalMonitor()
			sub := validSubscription()
			sub.ThresholdExpression = &AwsCostAnomalyMonitorExpression{
				Dimension: &AwsCostAnomalyMonitorExpressionDimension{
					Key:    "TOTAL_IMPACT",
					Values: []string{"100"},
				},
			}
			spec.Subscriptions = []*AwsCostAnomalyMonitorSubscription{sub}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
