package awsconfigrecorderv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestAwsConfigRecorderSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsConfigRecorderSpec Validation Suite")
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalRecorder is the smallest valid RUNNING posture: region, role,
// and a delivery channel.
func minimalRecorder() *AwsConfigRecorderSpec {
	return &AwsConfigRecorderSpec{
		Region:  "us-west-2",
		RoleArn: svr("arn:aws:iam::123456789012:role/config-recorder"),
		DeliveryChannel: &AwsConfigRecorderDeliveryChannel{
			S3BucketName: svr("my-config-bucket"),
		},
	}
}

var _ = ginkgo.Describe("AwsConfigRecorderSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal running recorder", func() {
			gomega.Expect(protovalidate.Validate(minimalRecorder())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a stopped recorder without a delivery channel", func() {
			spec := minimalRecorder()
			spec.DeliveryChannel = nil
			spec.RecordingEnabled = proto.Bool(false)
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the inclusion recording shape", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				AllSupported:      proto.Bool(false),
				ResourceTypes:     []string{"AWS::EC2::Instance", "AWS::S3::Bucket"},
				RecordingStrategy: "INCLUSION_BY_RESOURCE_TYPES",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the exclusion recording shape", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				AllSupported:             proto.Bool(false),
				ExclusionByResourceTypes: []string{"AWS::CloudFront::Distribution"},
				RecordingStrategy:        "EXCLUSION_BY_RESOURCE_TYPES",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a daily recording mode with an override", func() {
			spec := minimalRecorder()
			spec.RecordingMode = &AwsConfigRecorderRecordingMode{
				RecordingFrequency: "DAILY",
				Override: &AwsConfigRecorderRecordingModeOverride{
					RecordingFrequency: "CONTINUOUS",
					ResourceTypes:      []string{"AWS::IAM::Role"},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a retention window", func() {
			spec := minimalRecorder()
			spec.RetentionPeriodInDays = 365
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a missing region", func() {
			spec := minimalRecorder()
			spec.Region = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing role", func() {
			spec := minimalRecorder()
			spec.RoleArn = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a running recorder without a delivery channel", func() {
			spec := minimalRecorder()
			spec.DeliveryChannel = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects inclusion types while all_supported stays on", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				ResourceTypes: []string{"AWS::EC2::Instance"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects mixing inclusion and exclusion lists", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				AllSupported:             proto.Bool(false),
				ResourceTypes:            []string{"AWS::EC2::Instance"},
				ExclusionByResourceTypes: []string{"AWS::S3::Bucket"},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the ALL_SUPPORTED strategy with all_supported off", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				AllSupported:      proto.Bool(false),
				RecordingStrategy: "ALL_SUPPORTED_RESOURCE_TYPES",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an exclusion list under the inclusion strategy", func() {
			spec := minimalRecorder()
			spec.RecordingGroup = &AwsConfigRecorderRecordingGroup{
				AllSupported:             proto.Bool(false),
				ExclusionByResourceTypes: []string{"AWS::S3::Bucket"},
				RecordingStrategy:        "INCLUSION_BY_RESOURCE_TYPES",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an override without resource types", func() {
			spec := minimalRecorder()
			spec.RecordingMode = &AwsConfigRecorderRecordingMode{
				Override: &AwsConfigRecorderRecordingModeOverride{
					RecordingFrequency: "DAILY",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a retention window below 30 days", func() {
			spec := minimalRecorder()
			spec.RetentionPeriodInDays = 7
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown snapshot delivery frequency", func() {
			spec := minimalRecorder()
			spec.DeliveryChannel.SnapshotDeliveryFrequency = "Weekly"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
