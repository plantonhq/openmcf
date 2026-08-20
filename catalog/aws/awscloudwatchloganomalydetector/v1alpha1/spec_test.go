package awscloudwatchloganomalydetectorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsCloudwatchLogAnomalyDetectorSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsCloudwatchLogAnomalyDetectorSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalDetector() *AwsCloudwatchLogAnomalyDetectorSpec {
	return &AwsCloudwatchLogAnomalyDetectorSpec{
		Region: "us-east-1",
		LogGroupArns: []*foreignkeyv1.StringValueOrRef{
			svr("arn:aws:logs:us-east-1:123456789012:log-group:/app/api"),
		},
		Enabled: true,
	}
}

var _ = ginkgo.Describe("AwsCloudwatchLogAnomalyDetectorSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal detector", func() {
			gomega.Expect(protovalidate.Validate(minimalDetector())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the visibility window boundaries", func() {
			for _, days := range []int64{7, 90} {
				spec := minimalDetector()
				spec.AnomalyVisibilityTime = proto.Int64(days)
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			}
		})

		ginkgo.It("accepts a paused detector", func() {
			spec := minimalDetector()
			spec.Enabled = false
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects an empty log group list", func() {
			spec := minimalDetector()
			spec.LogGroupArns = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a visibility window outside 7-90 days", func() {
			for _, days := range []int64{6, 91} {
				spec := minimalDetector()
				spec.AnomalyVisibilityTime = proto.Int64(days)
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			}
		})

		ginkgo.It("rejects an unknown evaluation frequency", func() {
			spec := minimalDetector()
			spec.EvaluationFrequency = "TWO_MIN"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
