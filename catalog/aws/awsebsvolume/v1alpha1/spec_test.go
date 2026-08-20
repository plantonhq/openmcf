package awsebsvolumev1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsEbsVolumeSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEbsVolumeSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func minimalVolume() *AwsEbsVolumeSpec {
	return &AwsEbsVolumeSpec{
		Region:           "us-west-2",
		AvailabilityZone: "us-west-2a",
		Type:             "gp3",
		SizeGb:           20,
	}
}

var _ = ginkgo.Describe("AwsEbsVolumeSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the minimal fresh gp3 volume", func() {
			gomega.Expect(protovalidate.Validate(minimalVolume())).To(gomega.BeNil())
		})

		ginkgo.It("accepts a snapshot-restored volume without a size", func() {
			spec := minimalVolume()
			spec.SizeGb = 0
			spec.SnapshotId = literal("snap-0123456789abcdef0")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a copy without zone or size", func() {
			spec := &AwsEbsVolumeSpec{
				Region:   "us-west-2",
				CopyFrom: &AwsEbsVolumeCopyFrom{SourceVolumeId: literal("vol-0123456789abcdef0")},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an io2 multi-attach volume with two attachments", func() {
			spec := minimalVolume()
			spec.Type = "io2"
			spec.Iops = 1000
			spec.MultiAttachEnabled = true
			spec.Attachments = []*AwsEbsVolumeAttachment{
				{DeviceName: "/dev/sdf", InstanceId: literal("i-0123456789abcdef0")},
				{DeviceName: "/dev/sdf", InstanceId: literal("i-0123456789abcdef1")},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an initialization rate on a snapshot restore", func() {
			spec := minimalVolume()
			spec.SnapshotId = literal("snap-0123456789abcdef0")
			spec.VolumeInitializationRate = 200
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects a fresh volume without an availability zone", func() {
			spec := minimalVolume()
			spec.AvailabilityZone = ""
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a fresh volume with neither size nor snapshot", func() {
			spec := minimalVolume()
			spec.SizeGb = 0
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a copy that also sets create-arm fields", func() {
			spec := minimalVolume()
			spec.CopyFrom = &AwsEbsVolumeCopyFrom{SourceVolumeId: literal("vol-0123456789abcdef0")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects iops on a gp2 volume", func() {
			spec := minimalVolume()
			spec.Type = "gp2"
			spec.Iops = 3000
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an io1 volume without iops", func() {
			spec := minimalVolume()
			spec.Type = "io1"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects throughput outside gp3", func() {
			spec := minimalVolume()
			spec.Type = "io1"
			spec.Iops = 1000
			spec.ThroughputMibps = 250
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a throughput below the gp3 floor", func() {
			spec := minimalVolume()
			spec.ThroughputMibps = 50
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects multi-attach on gp3", func() {
			spec := minimalVolume()
			spec.MultiAttachEnabled = true
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two attachments without multi-attach", func() {
			spec := minimalVolume()
			spec.Attachments = []*AwsEbsVolumeAttachment{
				{DeviceName: "/dev/sdf", InstanceId: literal("i-0123456789abcdef0")},
				{DeviceName: "/dev/sdg", InstanceId: literal("i-0123456789abcdef1")},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an initialization rate without a snapshot", func() {
			spec := minimalVolume()
			spec.VolumeInitializationRate = 200
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an initialization rate outside 100-300", func() {
			spec := minimalVolume()
			spec.SnapshotId = literal("snap-0123456789abcdef0")
			spec.VolumeInitializationRate = 50
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown volume type", func() {
			spec := minimalVolume()
			spec.Type = "gp4"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a device name outside /dev/", func() {
			spec := minimalVolume()
			spec.Attachments = []*AwsEbsVolumeAttachment{
				{DeviceName: "sdf", InstanceId: literal("i-0123456789abcdef0")},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
