package awselasticfilesystemv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestAwsElasticFileSystemSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsElasticFileSystemSpec Validation Suite")
}

// helper to create a StringValueOrRef with a literal value.
func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// helper to build a mount target with just a subnet.
func mountTarget(subnetId string) *AwsElasticFileSystemMountTarget {
	return &AwsElasticFileSystemMountTarget{SubnetId: strRef(subnetId)}
}

var _ = ginkgo.Describe("AwsElasticFileSystemSpec validations", func() {
	var spec *AwsElasticFileSystemSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: mount_targets required (min 1).
		spec = &AwsElasticFileSystemSpec{
			Region: "us-east-1",
			MountTargets: []*AwsElasticFileSystemMountTarget{
				mountTarget("subnet-abc123"),
			},
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal valid spec (one mount target)", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts encrypted with default KMS", func() {
		spec.Encrypted = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts encrypted with custom KMS key", func() {
		spec.Encrypted = true
		spec.KmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/test-key")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts One Zone storage (single AZ)", func() {
		spec.AvailabilityZoneName = "us-east-1a"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts provisioned throughput mode", func() {
		spec.ThroughputMode = "provisioned"
		spec.ProvisionedThroughputInMibps = 100.0
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts elastic throughput on generalPurpose", func() {
		spec.ThroughputMode = "elastic"
		spec.PerformanceMode = "generalPurpose"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts maxIO with bursting throughput", func() {
		spec.PerformanceMode = "maxIO"
		spec.ThroughputMode = "bursting"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts lifecycle policy (transition_to_ia only)", func() {
		spec.TransitionToIa = "AFTER_30_DAYS"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts lifecycle policy (IA + archive + primary storage class)", func() {
		spec.TransitionToIa = "AFTER_30_DAYS"
		spec.TransitionToArchive = "AFTER_90_DAYS"
		spec.TransitionToPrimaryStorageClass = "AFTER_1_ACCESS"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts backup enabled", func() {
		spec.BackupEnabled = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts replication overwrite protection DISABLED", func() {
		spec.ReplicationOverwriteProtection = "DISABLED"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts multiple mount targets with security groups", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			mountTarget("subnet-abc123"),
			mountTarget("subnet-def456"),
			mountTarget("subnet-ghi789"),
		}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{strRef("sg-123")}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a mount target with a static IPv4 address", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:  strRef("subnet-abc123"),
				IpAddress: "10.0.1.50",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a dual-stack mount target with a static IPv6 address", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:      strRef("subnet-abc123"),
				IpAddressType: "DUAL_STACK",
				Ipv6Address:   "2001:db8::1",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts an IPv6-only mount target", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:      strRef("subnet-abc123"),
				IpAddressType: "IPV6_ONLY",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a resource policy with bypass safety check", func() {
		policy, _ := structpb.NewStruct(map[string]interface{}{
			"Version": "2012-10-17",
		})
		spec.Policy = policy
		spec.BypassPolicyLockoutSafetyCheck = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts cross-region replication", func() {
		spec.Replication = &AwsElasticFileSystemReplication{
			DestinationRegion: "us-east-2",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts One Zone replication (AZ only)", func() {
		spec.Replication = &AwsElasticFileSystemReplication{
			DestinationAvailabilityZoneName: "us-east-2a",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts replication into an existing file system", func() {
		spec.Replication = &AwsElasticFileSystemReplication{
			DestinationRegion:       "us-east-2",
			DestinationFileSystemId: strRef("fs-0123456789abcdef0"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts production-ready configuration", func() {
		policy, _ := structpb.NewStruct(map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": []interface{}{
				map[string]interface{}{
					"Effect":    "Deny",
					"Principal": map[string]interface{}{"AWS": "*"},
					"Action":    "*",
					"Condition": map[string]interface{}{
						"Bool": map[string]interface{}{
							"aws:SecureTransport": "false",
						},
					},
				},
			},
		})
		spec.Encrypted = true
		spec.KmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/test-key")
		spec.PerformanceMode = "generalPurpose"
		spec.ThroughputMode = "elastic"
		spec.TransitionToIa = "AFTER_30_DAYS"
		spec.TransitionToArchive = "AFTER_90_DAYS"
		spec.BackupEnabled = true
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			mountTarget("subnet-abc123"),
			mountTarget("subnet-def456"),
		}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{strRef("sg-123")}
		spec.Policy = policy
		spec.Replication = &AwsElasticFileSystemReplication{
			DestinationRegion: "us-east-2",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Failure cases
	// -------------------------------------------------------------------------

	ginkgo.It("fails when mount_targets is missing", func() {
		spec.MountTargets = nil
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when mount_targets is empty", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a mount target has no subnet", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{{}}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when performance_mode is invalid", func() {
		spec.PerformanceMode = "invalid-mode"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when throughput_mode is invalid", func() {
		spec.ThroughputMode = "invalid-throughput"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when elastic throughput is combined with maxIO", func() {
		spec.PerformanceMode = "maxIO"
		spec.ThroughputMode = "elastic"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when provisioned throughput is set without provisioned mode", func() {
		spec.ThroughputMode = "bursting"
		spec.ProvisionedThroughputInMibps = 100.0
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when provisioned mode is set without throughput value", func() {
		spec.ThroughputMode = "provisioned"
		spec.ProvisionedThroughputInMibps = 0.0
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when transition_to_archive is set without transition_to_ia", func() {
		spec.TransitionToArchive = "AFTER_90_DAYS"
		spec.TransitionToIa = ""
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when kms_key_id is set without encrypted", func() {
		spec.Encrypted = false
		spec.KmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/test-key")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when transition_to_ia has invalid value", func() {
		spec.TransitionToIa = "AFTER_2_DAYS"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when transition_to_primary_storage_class has invalid value", func() {
		spec.TransitionToPrimaryStorageClass = "AFTER_2_ACCESS"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when replication_overwrite_protection has invalid value", func() {
		spec.ReplicationOverwriteProtection = "enabled"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when bypass_policy_lockout_safety_check is set without policy", func() {
		spec.BypassPolicyLockoutSafetyCheck = true
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when ip_address_type is invalid", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:      strRef("subnet-abc123"),
				IpAddressType: "ipv4",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a static IPv4 address is set on an IPv6-only mount target", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:      strRef("subnet-abc123"),
				IpAddressType: "IPV6_ONLY",
				IpAddress:     "10.0.1.50",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when a static IPv6 address is set without an IPv6-capable type", func() {
		spec.MountTargets = []*AwsElasticFileSystemMountTarget{
			{
				SubnetId:    strRef("subnet-abc123"),
				Ipv6Address: "2001:db8::1",
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})

	ginkgo.It("fails when replication has neither region nor AZ", func() {
		spec.Replication = &AwsElasticFileSystemReplication{
			DestinationKmsKeyId: strRef("arn:aws:kms:us-east-2:123456789012:key/test-key"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).NotTo(gomega.BeNil())
	})
})
