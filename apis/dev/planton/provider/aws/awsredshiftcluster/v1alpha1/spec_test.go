package awsredshiftclusterv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsRedshiftClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRedshiftClusterSpec Custom Validation Tests")
}

// twoSubnets returns the minimum valid subnet list (two literal subnet IDs).
func twoSubnets() []*foreignkeyv1.StringValueOrRef {
	return []*foreignkeyv1.StringValueOrRef{
		{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-12345678"}},
		{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-87654321"}},
	}
}

// validCluster returns a minimal valid Redshift cluster that individual
// tests mutate into specific scenarios.
func validCluster() *AwsRedshiftCluster {
	return &AwsRedshiftCluster{
		ApiVersion: "aws.planton.dev/v1alpha1",
		Kind:       "AwsRedshiftCluster",
		Metadata: &shared.CloudResourceMetadata{
			Name: "test-redshift",
		},
		Spec: &AwsRedshiftClusterSpec{
			Region:               "us-west-2",
			NodeType:             "ra3.large",
			SubnetIds:            twoSubnets(),
			MasterUsername:       "rsadmin",
			ManageMasterPassword: true,
			SkipFinalSnapshot:    true,
		},
	}
}

var _ = ginkgo.Describe("AwsRedshiftClusterSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal managed-password cluster", func() {
			err := protovalidate.Validate(validCluster())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a multi-node production shape", func() {
			input := validCluster()
			input.Spec.NumberOfNodes = 4
			input.Spec.MultiAz = true
			input.Spec.EnhancedVpcRouting = true
			input.Spec.AutomatedSnapshotRetentionPeriod = int32Ptr(14)
			input.Spec.ManualSnapshotRetentionPeriod = 90
			input.Spec.Logging = &AwsRedshiftClusterLogging{
				LogDestinationType: "cloudwatch",
				LogExports:         []string{"connectionlog", "useractivitylog", "userlog"},
			}
			input.Spec.SnapshotCopy = &AwsRedshiftClusterSnapshotCopy{
				DestinationRegion: "us-east-1",
				RetentionPeriod:   7,
			}
			input.Spec.Parameters = []*AwsRedshiftClusterParameter{
				{Name: "require_ssl", Value: "true"},
			}
			input.Spec.ParameterGroupFamily = "redshift-2.0"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a plaintext-password cluster", func() {
			input := validCluster()
			input.Spec.ManageMasterPassword = false
			input.Spec.MasterPassword = "Correct1HorseBattery"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a headless snapshot restore without master_username", func() {
			input := validCluster()
			input.Spec.MasterUsername = ""
			input.Spec.ManageMasterPassword = false
			input.Spec.SnapshotIdentifier = "prod-final-2026-07-01"
			input.Spec.SnapshotClusterIdentifier = "prod-warehouse"
			input.Spec.OwnerAccount = "123456789012"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a snapshot-ARN restore", func() {
			input := validCluster()
			input.Spec.MasterUsername = ""
			input.Spec.SnapshotArn = "arn:aws:redshift:us-west-2:123456789012:snapshot:prod/prod-final"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a relocatable public cluster with an elastic IP", func() {
			input := validCluster()
			input.Spec.AvailabilityZoneRelocationEnabled = true
			input.Spec.PubliclyAccessible = true
			input.Spec.Port = 5439
			input.Spec.ElasticIp = &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "203.0.113.10"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.It("rejects missing subnets and subnet group (subnets_or_group)", func() {
			input := validCluster()
			input.Spec.SubnetIds = nil
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("provide at least two subnet_ids"))
		})

		ginkgo.It("rejects a port below the Redshift floor (port_range)", func() {
			input := validCluster()
			input.Spec.Port = 1024
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("port must be between 1115 and 65535"))
		})

		ginkgo.It("rejects an out-of-range node count (number_of_nodes_range)", func() {
			input := validCluster()
			input.Spec.NumberOfNodes = 129
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("number_of_nodes must be between 1 and 128"))
		})

		ginkgo.It("rejects both password strategies at once (password_xor_managed)", func() {
			input := validCluster()
			input.Spec.ManageMasterPassword = true
			input.Spec.MasterPassword = "Correct1HorseBattery"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("master_password cannot be set when manage_master_password is true"))
		})

		ginkgo.It("rejects a new cluster without master_username (master_username_required_unless_derived)", func() {
			input := validCluster()
			input.Spec.MasterUsername = ""
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("master_username is required for a new cluster"))
		})

		ginkgo.It("rejects relocation combined with Multi-AZ (relocation_xor_multi_az)", func() {
			input := validCluster()
			input.Spec.AvailabilityZoneRelocationEnabled = true
			input.Spec.MultiAz = true
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("availability_zone_relocation_enabled and multi_az are mutually exclusive"))
		})

		ginkgo.It("rejects an elastic IP on a private cluster (elastic_ip_requires_public)", func() {
			input := validCluster()
			input.Spec.ElasticIp = &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "203.0.113.10"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("elastic_ip requires publicly_accessible"))
		})

		ginkgo.It("rejects an out-of-range manual retention (manual_snapshot_retention_range)", func() {
			input := validCluster()
			input.Spec.ManualSnapshotRetentionPeriod = 4000
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("manual_snapshot_retention_period must be -1"))
		})

		ginkgo.It("rejects deletion without a final snapshot name (final_snapshot_id_required_when_not_skipping)", func() {
			input := validCluster()
			input.Spec.SkipFinalSnapshot = false
			input.Spec.FinalSnapshotIdentifier = ""
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("final_snapshot_identifier is required when skip_final_snapshot is false"))
		})

		ginkgo.It("rejects both restore sources at once (snapshot_id_xor_arn)", func() {
			input := validCluster()
			input.Spec.SnapshotIdentifier = "prod-final"
			input.Spec.SnapshotArn = "arn:aws:redshift:us-west-2:123456789012:snapshot:prod/prod-final"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("snapshot_identifier and snapshot_arn are mutually exclusive"))
		})

		ginkgo.It("rejects snapshot_cluster_identifier without a snapshot name (snapshot_cluster_requires_snapshot_id)", func() {
			input := validCluster()
			input.Spec.SnapshotClusterIdentifier = "prod-warehouse"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("snapshot_cluster_identifier is only meaningful alongside snapshot_identifier"))
		})

		ginkgo.It("rejects owner_account without a restore source (owner_account_requires_restore)", func() {
			input := validCluster()
			input.Spec.OwnerAccount = "123456789012"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("owner_account is only meaningful alongside a restore source"))
		})

		ginkgo.It("rejects inline parameters combined with an existing group (own_parameters_xor_existing_group)", func() {
			input := validCluster()
			input.Spec.ClusterParameterGroupName = "existing-group"
			input.Spec.Parameters = []*AwsRedshiftClusterParameter{
				{Name: "require_ssl", Value: "true"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("parameters and cluster_parameter_group_name are mutually exclusive"))
		})

		ginkgo.It("rejects a parameter-group family without inline parameters (family_requires_parameters)", func() {
			input := validCluster()
			input.Spec.ParameterGroupFamily = "redshift-2.0"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("parameter_group_family is only meaningful alongside inline parameters"))
		})

		ginkgo.It("rejects S3 logging without a bucket (s3_bucket_required)", func() {
			input := validCluster()
			input.Spec.Logging = &AwsRedshiftClusterLogging{
				LogDestinationType: "s3",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("s3_bucket_name is required when log_destination_type is 's3'"))
		})

		ginkgo.It("rejects CloudWatch logging without exports (cloudwatch_exports_required)", func() {
			input := validCluster()
			input.Spec.Logging = &AwsRedshiftClusterLogging{
				LogDestinationType: "cloudwatch",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("log_exports must have at least one entry"))
		})

		ginkgo.It("rejects an invalid log export type", func() {
			input := validCluster()
			input.Spec.Logging = &AwsRedshiftClusterLogging{
				LogDestinationType: "cloudwatch",
				LogExports:         []string{"queries"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a snapshot copy without a destination region", func() {
			input := validCluster()
			input.Spec.SnapshotCopy = &AwsRedshiftClusterSnapshotCopy{}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an out-of-range snapshot copy retention (copy_retention_range)", func() {
			input := validCluster()
			input.Spec.SnapshotCopy = &AwsRedshiftClusterSnapshotCopy{
				DestinationRegion: "us-east-1",
				RetentionPeriod:   36,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("retention_period must be between 1 and 35"))
		})

		ginkgo.It("rejects an out-of-range snapshot copy manual retention (copy_manual_retention_range)", func() {
			input := validCluster()
			input.Spec.SnapshotCopy = &AwsRedshiftClusterSnapshotCopy{
				DestinationRegion:             "us-east-1",
				ManualSnapshotRetentionPeriod: -2,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("manual_snapshot_retention_period must be -1"))
		})

		ginkgo.It("rejects an out-of-range automated snapshot retention", func() {
			input := validCluster()
			input.Spec.AutomatedSnapshotRetentionPeriod = int32Ptr(36)
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed maintenance window", func() {
			input := validCluster()
			input.Spec.PreferredMaintenanceWindow = "saturday:03:00-04:00"
			err := protovalidate.Validate(input)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})

// int32Ptr returns a pointer to the given int32 (proto3 optional scalars
// take pointers).
func int32Ptr(v int32) *int32 {
	return &v
}
