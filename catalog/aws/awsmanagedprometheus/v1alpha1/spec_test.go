package awsmanagedprometheusv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsManagedPrometheusSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsManagedPrometheusSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalWorkspace() *AwsManagedPrometheusSpec {
	return &AwsManagedPrometheusSpec{
		Region: "us-east-1",
	}
}

func sampleDetector() *AwsManagedPrometheusAnomalyDetector {
	return &AwsManagedPrometheusAnomalyDetector{
		Alias:             "request-rate",
		Query:             "sum(rate(http_requests_total[5m]))",
		MissingDataAction: "SKIP",
	}
}

var _ = ginkgo.Describe("AwsManagedPrometheusSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal workspace", func() {
			gomega.Expect(protovalidate.Validate(minimalWorkspace())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a fully loaded workspace", func() {
			spec := minimalWorkspace()
			spec.Alias = "platform-metrics"
			spec.Logging = &AwsManagedPrometheusLogging{
				LogGroupArn: svr("arn:aws:logs:us-east-1:123456789012:log-group:/amp/workspace"),
			}
			spec.Configuration = &AwsManagedPrometheusWorkspaceConfiguration{
				RetentionPeriodInDays: proto.Int32(90),
				LimitsPerLabelSet: []*AwsManagedPrometheusLabelSetLimit{
					{LabelSet: map[string]string{"team": "ingest"}, MaxSeries: 1000000},
				},
			}
			spec.AlertManagerDefinition = "alertmanager_config: |\n  route:\n    receiver: default\n"
			spec.RuleGroupNamespaces = []*AwsManagedPrometheusRuleGroupNamespace{
				{Name: "slo-rules", Data: "groups:\n  - name: slo\n    rules: []\n"},
			}
			spec.QueryLogging = &AwsManagedPrometheusQueryLogging{
				Destinations: []*AwsManagedPrometheusQueryLoggingDestination{
					{
						LogGroupArn:  svr("arn:aws:logs:us-east-1:123456789012:log-group:/amp/queries"),
						QspThreshold: 10000,
					},
				},
			}
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{sampleDetector()}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a detector suppression band with exactly one shape", func() {
			spec := minimalWorkspace()
			detector := sampleDetector()
			detector.IgnoreNearExpectedFromAbove = &AwsManagedPrometheusIgnoreNearExpected{Ratio: proto.Float64(0.1)}
			detector.IgnoreNearExpectedFromBelow = &AwsManagedPrometheusIgnoreNearExpected{Amount: proto.Float64(5)}
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{detector}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects duplicate rule namespace names", func() {
			spec := minimalWorkspace()
			spec.RuleGroupNamespaces = []*AwsManagedPrometheusRuleGroupNamespace{
				{Name: "dup", Data: "groups: []"},
				{Name: "dup", Data: "groups: []"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate detector aliases", func() {
			spec := minimalWorkspace()
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{sampleDetector(), sampleDetector()}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown missing-data action", func() {
			spec := minimalWorkspace()
			detector := sampleDetector()
			detector.MissingDataAction = "IGNORE"
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{detector}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a suppression band with both shapes", func() {
			spec := minimalWorkspace()
			detector := sampleDetector()
			detector.IgnoreNearExpectedFromAbove = &AwsManagedPrometheusIgnoreNearExpected{
				Amount: proto.Float64(5),
				Ratio:  proto.Float64(0.1),
			}
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{detector}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a suppression band with neither shape", func() {
			spec := minimalWorkspace()
			detector := sampleDetector()
			detector.IgnoreNearExpectedFromBelow = &AwsManagedPrometheusIgnoreNearExpected{}
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{detector}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an RCF sample size below AWS's floor", func() {
			spec := minimalWorkspace()
			detector := sampleDetector()
			detector.SampleSize = proto.Int32(255)
			spec.AnomalyDetectors = []*AwsManagedPrometheusAnomalyDetector{detector}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a label-set limit with an empty label set", func() {
			spec := minimalWorkspace()
			spec.Configuration = &AwsManagedPrometheusWorkspaceConfiguration{
				LimitsPerLabelSet: []*AwsManagedPrometheusLabelSetLimit{
					{MaxSeries: 100},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-order window above AWS's cap", func() {
			spec := minimalWorkspace()
			spec.Configuration = &AwsManagedPrometheusWorkspaceConfiguration{
				OutOfOrderTimeWindowInSeconds: proto.Int32(601),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty query-logging destination list", func() {
			spec := minimalWorkspace()
			spec.QueryLogging = &AwsManagedPrometheusQueryLogging{}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
