package awsrdsinstancev1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsRdsInstanceSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsRdsInstanceSpec Custom Validation Tests")
}

// twoSubnets returns the minimum valid subnet list (two literal subnet IDs).
func twoSubnets() []*foreignkeyv1.StringValueOrRef {
	return []*foreignkeyv1.StringValueOrRef{
		{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-12345678"}},
		{LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "subnet-87654321"}},
	}
}

// validInstance returns a minimal valid postgres instance that individual
// tests mutate into specific scenarios.
func validInstance() *AwsRdsInstance {
	return &AwsRdsInstance{
		ApiVersion: "aws.planton.dev/v1",
		Kind:       "AwsRdsInstance",
		Metadata: &shared.CloudResourceMetadata{
			Name: "test-rds-instance",
		},
		Spec: &AwsRdsInstanceSpec{
			Region:                   "us-west-2",
			SubnetIds:                twoSubnets(),
			Engine:                   "postgres",
			InstanceClass:            "db.t4g.micro",
			AllocatedStorageGb:       20,
			Username:                 "dbadmin",
			ManageMasterUserPassword: true,
			SkipFinalSnapshot:        true,
		},
	}
}

var _ = ginkgo.Describe("AwsRdsInstanceSpec Custom Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.It("accepts a minimal postgres instance with a managed password", func() {
			err := protovalidate.Validate(validInstance())
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a production-shaped multi-AZ instance", func() {
			input := validInstance()
			input.Spec.MultiAz = true
			input.Spec.StorageType = "gp3"
			input.Spec.MaxAllocatedStorageGb = 100
			input.Spec.BackupRetentionPeriod = 7
			input.Spec.DeletionProtection = true
			input.Spec.PerformanceInsightsEnabled = true
			input.Spec.EnabledCloudwatchLogsExports = []string{"postgresql", "upgrade"}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a read replica without engine, storage, or credentials", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.Username = ""
			input.Spec.ManageMasterUserPassword = false
			input.Spec.ReplicateSourceDb = "source-instance"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a mounted Oracle replica", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.Username = ""
			input.Spec.ManageMasterUserPassword = false
			input.Spec.ReplicateSourceDb = "arn:aws:rds:us-east-1:123456789012:db:oracle-source"
			input.Spec.ReplicaMode = "mounted"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a snapshot restore without engine or storage", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.SnapshotIdentifier = "my-snapshot"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a point-in-time restore", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.RestoreToPointInTime = &AwsRdsInstanceRestoreToPointInTime{
				SourceDbInstanceIdentifier: "source-instance",
				UseLatestRestorableTime:    true,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a direct password without the managed strategy", func() {
			input := validInstance()
			input.Spec.ManageMasterUserPassword = false
			input.Spec.Password = "super-secret-password"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts an AWS-managed Active Directory join", func() {
			input := validInstance()
			input.Spec.Engine = "sqlserver-se"
			input.Spec.LicenseModel = "license-included"
			input.Spec.ActiveDirectory = &AwsRdsInstanceActiveDirectory{
				Domain:            "d-1234567890",
				DomainIamRoleName: "rds-ad-role",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a self-managed Active Directory join", func() {
			input := validInstance()
			input.Spec.ActiveDirectory = &AwsRdsInstanceActiveDirectory{
				DomainFqdn:          "corp.example.com",
				DomainOu:            "OU=Databases,DC=corp,DC=example,DC=com",
				DomainAuthSecretArn: "arn:aws:secretsmanager:us-west-2:123456789012:secret:ad-join",
				DomainDnsIps:        []string{"10.0.0.10", "10.0.1.10"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("networking validations", func() {

		ginkgo.It("rejects a spec with neither subnets nor a subnet group", func() {
			input := validInstance()
			input.Spec.SubnetIds = nil
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("provide at least two subnet_ids (distinct AZs) or an"))
		})

		ginkgo.It("rejects an invalid network_type", func() {
			input := validInstance()
			input.Spec.NetworkType = "IPV6"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("network_type must be 'IPV4' or 'DUAL' when set"))
		})

		ginkgo.It("rejects an AZ pin on a multi-AZ instance", func() {
			input := validInstance()
			input.Spec.MultiAz = true
			input.Spec.AvailabilityZone = "us-west-2a"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("availability_zone cannot be pinned on a Multi-AZ instance"))
		})
	})

	ginkgo.Describe("engine and storage validations", func() {

		ginkgo.It("rejects a new instance without an engine", func() {
			input := validInstance()
			input.Spec.Engine = ""
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("engine is required unless the instance derives it from a"))
		})

		ginkgo.It("rejects a new instance without storage", func() {
			input := validInstance()
			input.Spec.AllocatedStorageGb = 0
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("allocated_storage_gb is required unless storage is"))
		})

		ginkgo.It("rejects an instance class without the db. prefix", func() {
			input := validInstance()
			input.Spec.InstanceClass = "t4g.micro"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown storage type", func() {
			input := validInstance()
			input.Spec.StorageType = "gp4"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an autoscaling ceiling at or below allocated storage", func() {
			input := validInstance()
			input.Spec.MaxAllocatedStorageGb = 20
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("max_allocated_storage_gb must exceed allocated_storage_gb"))
		})

		ginkgo.It("rejects iops on burst-credit storage", func() {
			input := validInstance()
			input.Spec.StorageType = "gp2"
			input.Spec.Iops = 3000
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("iops applies to io1, io2, or gp3 storage"))
		})

		ginkgo.It("rejects storage throughput outside gp3", func() {
			input := validInstance()
			input.Spec.StorageType = "io1"
			input.Spec.Iops = 1000
			input.Spec.StorageThroughput = 500
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("storage_throughput only applies to gp3 storage"))
		})
	})

	ginkgo.Describe("password strategy validations", func() {

		ginkgo.It("rejects a password alongside the managed password", func() {
			input := validInstance()
			input.Spec.Password = "also-a-password"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("password cannot be set when manage_master_user_password is"))
		})

		ginkgo.It("rejects a new instance without a username", func() {
			input := validInstance()
			input.Spec.Username = ""
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("username is required"))
		})

		ginkgo.It("accepts a blank username on a snapshot restore", func() {
			input := validInstance()
			input.Spec.Username = ""
			input.Spec.ManageMasterUserPassword = false
			input.Spec.SnapshotIdentifier = "my-snapshot"
			err := protovalidate.Validate(input)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("deletion safety validations", func() {

		ginkgo.It("rejects a missing final snapshot name while not skipping", func() {
			input := validInstance()
			input.Spec.SkipFinalSnapshot = false
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("final_snapshot_identifier is required when"))
		})
	})

	ginkgo.Describe("replica and restore validations", func() {

		ginkgo.It("rejects two create sources together", func() {
			input := validInstance()
			input.Spec.ReplicateSourceDb = "source-instance"
			input.Spec.SnapshotIdentifier = "my-snapshot"
			input.Spec.ManageMasterUserPassword = false
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("replicate_source_db, snapshot_identifier, and"))
		})

		ginkgo.It("rejects credentials on a read replica", func() {
			input := validInstance()
			input.Spec.ReplicateSourceDb = "source-instance"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("username, password, and manage_master_user_password cannot"))
		})

		ginkgo.It("rejects replica_mode without a replica source", func() {
			input := validInstance()
			input.Spec.ReplicaMode = "mounted"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("replica_mode only applies to a read replica"))
		})

		ginkgo.It("rejects an invalid replica_mode", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.ManageMasterUserPassword = false
			input.Spec.ReplicateSourceDb = "source-instance"
			input.Spec.ReplicaMode = "read-write"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("replica_mode must be 'open-read-only' or 'mounted' when set"))
		})

		ginkgo.It("rejects blue/green on a read replica", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.ManageMasterUserPassword = false
			input.Spec.ReplicateSourceDb = "source-instance"
			input.Spec.BlueGreenUpdateEnabled = true
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("blue_green_update_enabled cannot be combined with"))
		})

		ginkgo.It("rejects a point-in-time restore with two sources", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.RestoreToPointInTime = &AwsRdsInstanceRestoreToPointInTime{
				SourceDbInstanceIdentifier: "source-instance",
				SourceDbiResourceId:        "db-ABC123",
				UseLatestRestorableTime:    true,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of source_db_instance_identifier,"))
		})

		ginkgo.It("rejects a point-in-time restore with both time forms", func() {
			input := validInstance()
			input.Spec.Engine = ""
			input.Spec.AllocatedStorageGb = 0
			input.Spec.RestoreToPointInTime = &AwsRdsInstanceRestoreToPointInTime{
				SourceDbInstanceIdentifier: "source-instance",
				RestoreTime:                "2026-07-01T09:45:00Z",
				UseLatestRestorableTime:    true,
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("exactly one of restore_time or use_latest_restorable_time"))
		})
	})

	ginkgo.Describe("observability validations", func() {

		ginkgo.It("rejects an invalid performance insights retention", func() {
			input := validInstance()
			input.Spec.PerformanceInsightsRetentionPeriod = 30
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("performance_insights_retention_period must be 7, 731, or a"))
		})

		ginkgo.It("rejects an invalid monitoring interval", func() {
			input := validInstance()
			input.Spec.MonitoringInterval = 45
			input.Spec.MonitoringRoleArn = &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "arn:aws:iam::123456789012:role/rds-monitoring"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("monitoring_interval must be 0 (disabled), 1, 5, 10, 15, 30,"))
		})

		ginkgo.It("rejects enhanced monitoring without a role", func() {
			input := validInstance()
			input.Spec.MonitoringInterval = 30
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("monitoring_role_arn is required when monitoring_interval is"))
		})

		ginkgo.It("rejects an unknown log export type", func() {
			input := validInstance()
			input.Spec.EnabledCloudwatchLogsExports = []string{"binlog"}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an invalid database insights mode", func() {
			input := validInstance()
			input.Spec.DatabaseInsightsMode = "premium"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("database_insights_mode must be 'standard' or 'advanced'"))
		})
	})

	ginkgo.Describe("Active Directory validations", func() {

		ginkgo.It("rejects mixing the managed and self-managed shapes", func() {
			input := validInstance()
			input.Spec.ActiveDirectory = &AwsRdsInstanceActiveDirectory{
				Domain:            "d-1234567890",
				DomainIamRoleName: "rds-ad-role",
				DomainFqdn:        "corp.example.com",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("use either the AWS-managed shape (domain +"))
		})

		ginkgo.It("rejects a managed domain without its IAM role", func() {
			input := validInstance()
			input.Spec.ActiveDirectory = &AwsRdsInstanceActiveDirectory{
				Domain: "d-1234567890",
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("domain and domain_iam_role_name are required together for"))
		})

		ginkgo.It("rejects a self-managed AD without exactly two DNS IPs", func() {
			input := validInstance()
			input.Spec.ActiveDirectory = &AwsRdsInstanceActiveDirectory{
				DomainFqdn:          "corp.example.com",
				DomainOu:            "OU=Databases,DC=corp,DC=example,DC=com",
				DomainAuthSecretArn: "arn:aws:secretsmanager:us-west-2:123456789012:secret:ad-join",
				DomainDnsIps:        []string{"10.0.0.10"},
			}
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("self-managed AD requires exactly two domain_dns_ips"))
		})
	})

	ginkgo.Describe("misc validations", func() {

		ginkgo.It("rejects an invalid license model", func() {
			input := validInstance()
			input.Spec.LicenseModel = "site-license"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
			gomega.Expect(err.Error()).To(gomega.ContainSubstring("license_model must be 'license-included',"))
		})

		ginkgo.It("rejects a malformed backup window", func() {
			input := validInstance()
			input.Spec.BackupWindow = "5am-6am"
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a backup retention above 35 days", func() {
			input := validInstance()
			input.Spec.BackupRetentionPeriod = 36
			err := protovalidate.Validate(input)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})
