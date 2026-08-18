package awscloudwatchdashboardv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsCloudwatchDashboardSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchDashboardSpec Validation Suite")
}

func sampleBody() *structpb.Struct {
	body, err := structpb.NewStruct(map[string]any{
		"widgets": []any{
			map[string]any{
				"type":   "text",
				"x":      0.0,
				"y":      0.0,
				"width":  6.0,
				"height": 3.0,
				"properties": map[string]any{
					"markdown": "# Service health",
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return body
}

func minimalDashboard() *AwsCloudwatchDashboardSpec {
	return &AwsCloudwatchDashboardSpec{
		Region:        "us-east-1",
		DashboardName: "ServiceHealth",
		DashboardBody: sampleBody(),
	}
}

var _ = ginkgo.Describe("AwsCloudwatchDashboardSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal dashboard", func() {
			gomega.Expect(protovalidate.Validate(minimalDashboard())).To(gomega.BeNil())
		})

		ginkgo.It("accepts underscores and hyphens in the name", func() {
			spec := minimalDashboard()
			spec.DashboardName = "service_health-prod"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalDashboard()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a name with characters outside the dashboard charset", func() {
			spec := minimalDashboard()
			spec.DashboardName = "service health"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing dashboard body", func() {
			spec := minimalDashboard()
			spec.DashboardBody = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
