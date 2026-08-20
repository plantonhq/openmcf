package awsebssnapshotv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsEbsSnapshotSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEbsSnapshotSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func volumeArmSnapshot() *AwsEbsSnapshotSpec {
	return &AwsEbsSnapshotSpec{
		Region:   "us-west-2",
		VolumeId: literal("vol-0123456789abcdef0"),
	}
}

var _ = ginkgo.Describe("AwsEbsSnapshotSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts the volume arm", func() {
			gomega.Expect(protovalidate.Validate(volumeArmSnapshot())).To(gomega.BeNil())
		})

		ginkgo.It("accepts the copy arm with re-encryption", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				CopyFrom: &AwsEbsSnapshotCopyFrom{
					SourceSnapshotId: literal("snap-0123456789abcdef0"),
					SourceRegion:     "us-west-2",
					Encrypted:        true,
					KmsKeyId:         literal("arn:aws:kms:us-west-2:111122223333:key/abc"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the import arm with an S3 disk container", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				ImportFrom: &AwsEbsSnapshotImportFrom{
					DiskContainer: &AwsEbsSnapshotDiskContainer{
						Format:   "RAW",
						S3Bucket: "images",
						S3Key:    "disk.raw",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts archive tiering with a temporary restore", func() {
			spec := volumeArmSnapshot()
			spec.StorageTier = "archive"
			spec.TemporaryRestoreDays = 30
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts fast restore zones and share grants", func() {
			spec := volumeArmSnapshot()
			spec.FastRestoreAvailabilityZones = []string{"us-west-2a", "us-west-2b"}
			spec.ShareWithAccountIds = []string{"111122223333"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a 15-multiple completion duration", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				CopyFrom: &AwsEbsSnapshotCopyFrom{
					SourceSnapshotId:          literal("snap-0123456789abcdef0"),
					SourceRegion:              "us-west-2",
					CompletionDurationMinutes: 45,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects zero source arms", func() {
			spec := &AwsEbsSnapshotSpec{Region: "us-west-2"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects two source arms", func() {
			spec := volumeArmSnapshot()
			spec.CopyFrom = &AwsEbsSnapshotCopyFrom{
				SourceSnapshotId: literal("snap-0123456789abcdef0"),
				SourceRegion:     "us-west-2",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a copy KMS key without encrypted", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				CopyFrom: &AwsEbsSnapshotCopyFrom{
					SourceSnapshotId: literal("snap-0123456789abcdef0"),
					SourceRegion:     "us-west-2",
					KmsKeyId:         literal("arn:aws:kms:us-west-2:111122223333:key/abc"),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a completion duration off the 15-minute grid", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				CopyFrom: &AwsEbsSnapshotCopyFrom{
					SourceSnapshotId:          literal("snap-0123456789abcdef0"),
					SourceRegion:              "us-west-2",
					CompletionDurationMinutes: 40,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a disk container with both url and s3 bucket", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				ImportFrom: &AwsEbsSnapshotImportFrom{
					DiskContainer: &AwsEbsSnapshotDiskContainer{
						Format:   "VMDK",
						Url:      "https://example.com/disk.vmdk",
						S3Bucket: "images",
						S3Key:    "disk.vmdk",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an s3 bucket without a key", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				ImportFrom: &AwsEbsSnapshotImportFrom{
					DiskContainer: &AwsEbsSnapshotDiskContainer{
						Format:   "RAW",
						S3Bucket: "images",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects restore dials on the standard tier", func() {
			spec := volumeArmSnapshot()
			spec.PermanentRestore = true
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a temporary restore alongside a permanent one", func() {
			spec := volumeArmSnapshot()
			spec.StorageTier = "archive"
			spec.PermanentRestore = true
			spec.TemporaryRestoreDays = 30
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects duplicate fast restore zones", func() {
			spec := volumeArmSnapshot()
			spec.FastRestoreAvailabilityZones = []string{"us-west-2a", "us-west-2a"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed share account id", func() {
			spec := volumeArmSnapshot()
			spec.ShareWithAccountIds = []string{"not-an-account"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown disk image format", func() {
			spec := &AwsEbsSnapshotSpec{
				Region: "us-west-2",
				ImportFrom: &AwsEbsSnapshotImportFrom{
					DiskContainer: &AwsEbsSnapshotDiskContainer{
						Format:   "QCOW2",
						S3Bucket: "images",
						S3Key:    "disk.qcow2",
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
