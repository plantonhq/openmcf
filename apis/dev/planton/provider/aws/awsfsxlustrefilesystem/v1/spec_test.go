package awsfsxlustrefilesystemv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsFsxLustreFileSystemSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsFsxLustreFileSystemSpec Validation Suite")
}

func strRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func stringPtr(s string) *string {
	return &s
}

func int32Ptr(i int32) *int32 {
	return &i
}

var _ = ginkgo.Describe("AwsFsxLustreFileSystemSpec validations", func() {
	var spec *AwsFsxLustreFileSystemSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: subnet_id required, storage_capacity_gib >= 1200.
		spec = &AwsFsxLustreFileSystemSpec{
			Region:             "us-west-2",
			SubnetId:           strRef("subnet-abc123"),
			StorageCapacityGib: int32Ptr(1200),
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal valid spec", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts SCRATCH_1 deployment type", func() {
		spec.DeploymentType = stringPtr("SCRATCH_1")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts SCRATCH_2 deployment type", func() {
		spec.DeploymentType = stringPtr("SCRATCH_2")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts PERSISTENT_1 deployment type with throughput", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_1")
		spec.PerUnitStorageThroughput = int32Ptr(50)
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts PERSISTENT_2 deployment type with throughput", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.PerUnitStorageThroughput = int32Ptr(125)
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts HDD storage type with PERSISTENT_1 and a drive cache decision", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_1")
		spec.StorageType = stringPtr("HDD")
		spec.StorageCapacityGib = int32Ptr(6000)
		spec.PerUnitStorageThroughput = int32Ptr(12)
		spec.DriveCacheType = "READ"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts HDD with drive_cache_type NONE", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_1")
		spec.StorageType = stringPtr("HDD")
		spec.StorageCapacityGib = int32Ptr(6000)
		spec.PerUnitStorageThroughput = int32Ptr(40)
		spec.DriveCacheType = "NONE"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts SSD storage type", func() {
		spec.StorageType = stringPtr("SSD")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts the full INTELLIGENT_TIERING companion set", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.StorageType = stringPtr("INTELLIGENT_TIERING")
		spec.StorageCapacityGib = nil
		spec.ThroughputCapacity = int32Ptr(4000)
		spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
			SizingMode: "PROPORTIONAL_TO_THROUGHPUT_CAPACITY",
		}
		spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
			Mode: stringPtr("AUTOMATIC"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a user-provisioned read cache within bounds", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.StorageType = stringPtr("INTELLIGENT_TIERING")
		spec.StorageCapacityGib = nil
		spec.ThroughputCapacity = int32Ptr(8000)
		spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
			SizingMode: "USER_PROVISIONED",
			SizeGib:    int32Ptr(64),
		}
		spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts efa_enabled on PERSISTENT_2 with metadata configuration", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.PerUnitStorageThroughput = int32Ptr(500)
		spec.EfaEnabled = true
		spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{strRef("sg-efa123")}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a backup restore without storage capacity", func() {
		spec.StorageCapacityGib = nil
		spec.BackupId = "backup-0123456789abcdef0"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts LZ4 data compression type", func() {
		spec.DataCompressionType = stringPtr("LZ4")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts import_path on SCRATCH_2", func() {
		spec.DeploymentType = stringPtr("SCRATCH_2")
		spec.ImportPath = "s3://my-bucket/data"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts the full legacy S3 link arm on PERSISTENT_1", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_1")
		spec.PerUnitStorageThroughput = int32Ptr(100)
		spec.ImportPath = "s3://my-bucket/input"
		spec.ExportPath = "s3://my-bucket/output"
		spec.AutoImportPolicy = "NEW_CHANGED_DELETED"
		spec.ImportedFileChunkSize = int32Ptr(2048)
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a root squash configuration", func() {
		spec.RootSquashConfiguration = &AwsFsxLustreFileSystemRootSquashConfiguration{
			RootSquash:   "65534:65534",
			NoSquashNids: []string{"10.0.1.6@tcp", "10.0.[2-10].[1-255]@tcp"},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts metadata_configuration with USER_PROVISIONED mode and iops", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.PerUnitStorageThroughput = int32Ptr(125)
		spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
			Mode: stringPtr("USER_PROVISIONED"),
			Iops: int32Ptr(3000),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts log_configuration with valid level", func() {
		spec.LogConfiguration = &AwsFsxLustreFileSystemLogConfiguration{
			Destination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/fsx/lustre"),
			Level:       stringPtr("WARN_ERROR"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a production-ready PERSISTENT_2 configuration", func() {
		spec.DeploymentType = stringPtr("PERSISTENT_2")
		spec.StorageCapacityGib = int32Ptr(2400)
		spec.StorageType = stringPtr("SSD")
		spec.PerUnitStorageThroughput = int32Ptr(250)
		spec.DataCompressionType = stringPtr("LZ4")
		spec.SecurityGroupIds = []*foreignkeyv1.StringValueOrRef{strRef("sg-abc123")}
		spec.KmsKeyId = strRef("arn:aws:kms:us-east-1:123456789012:key/test-key")
		spec.AutomaticBackupRetentionDays = int32Ptr(7)
		spec.DailyAutomaticBackupStartTime = "05:00"
		spec.CopyTagsToBackups = true
		spec.SkipFinalBackup = boolPtr(false)
		spec.FinalBackupTags = map[string]string{"retention": "legal-hold"}
		spec.WeeklyMaintenanceStartTime = "1:05:00"
		spec.RootSquashConfiguration = &AwsFsxLustreFileSystemRootSquashConfiguration{
			RootSquash: "65534:65534",
		}
		spec.LogConfiguration = &AwsFsxLustreFileSystemLogConfiguration{
			Destination: strRef("arn:aws:logs:us-east-1:123456789012:log-group:/aws/fsx/lustre"),
			Level:       stringPtr("WARN_ERROR"),
		}
		spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
			Mode: stringPtr("AUTOMATIC"),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Field-level validation failures
	// -------------------------------------------------------------------------

	ginkgo.Context("field-level validations", func() {
		ginkgo.It("fails when subnet_id is missing", func() {
			spec.SubnetId = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when storage_capacity_gib is below minimum", func() {
			spec.StorageCapacityGib = int32Ptr(100)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when per_unit_storage_throughput is not a valid tier", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.PerUnitStorageThroughput = int32Ptr(300)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when file_system_type_version is malformed", func() {
			spec.FileSystemTypeVersion = "v2.15"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when daily_automatic_backup_start_time is malformed", func() {
			spec.DailyAutomaticBackupStartTime = "5am"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when weekly_maintenance_start_time is malformed", func() {
			spec.WeeklyMaintenanceStartTime = "8:05:00"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when import_path is not an s3 URI", func() {
			spec.ImportPath = "https://my-bucket/data"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when more than 50 security groups are supplied", func() {
			for i := 0; i < 51; i++ {
				spec.SecurityGroupIds = append(spec.SecurityGroupIds, strRef("sg-x"))
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: deployment_type_valid / storage_type_valid
	// -------------------------------------------------------------------------

	ginkgo.Context("type enums", func() {
		ginkgo.It("fails when deployment_type is invalid", func() {
			spec.DeploymentType = stringPtr("INVALID_TYPE")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when storage_type is invalid", func() {
			spec.StorageType = stringPtr("NVME")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when data_compression_type is invalid", func() {
			spec.DataCompressionType = stringPtr("GZIP")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: storage_capacity_required / backup_restore_excludes_capacity
	// -------------------------------------------------------------------------

	ginkgo.Context("storage capacity presence", func() {
		ginkgo.It("fails when storage_capacity_gib is omitted without backup or intelligent tiering", func() {
			spec.StorageCapacityGib = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a backup restore also provisions capacity", func() {
			spec.BackupId = "backup-0123456789abcdef0"
			spec.StorageCapacityGib = int32Ptr(1200)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: hdd_requires_persistent_1 / drive_cache_type_hdd_contract
	// -------------------------------------------------------------------------

	ginkgo.Context("HDD contract", func() {
		ginkgo.It("fails when HDD is used with SCRATCH_2", func() {
			spec.DeploymentType = stringPtr("SCRATCH_2")
			spec.StorageType = stringPtr("HDD")
			spec.DriveCacheType = "READ"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when HDD omits the drive cache decision", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_1")
			spec.StorageType = stringPtr("HDD")
			spec.PerUnitStorageThroughput = int32Ptr(12)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when drive_cache_type is set on SSD storage", func() {
			spec.DriveCacheType = "READ"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when drive_cache_type has an invalid value", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_1")
			spec.StorageType = stringPtr("HDD")
			spec.PerUnitStorageThroughput = int32Ptr(12)
			spec.DriveCacheType = "WRITE"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: intelligent_tiering_contract and companions
	// -------------------------------------------------------------------------

	ginkgo.Context("INTELLIGENT_TIERING contract", func() {
		ginkgo.It("fails without the companion set", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on a non-PERSISTENT_2 deployment", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_1")
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ThroughputCapacity = int32Ptr(4000)
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "NO_CACHE",
			}
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when throughput_capacity is not a multiple of 4000", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ThroughputCapacity = int32Ptr(6000)
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "PROPORTIONAL_TO_THROUGHPUT_CAPACITY",
			}
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when throughput_capacity is set on provisioned storage", func() {
			spec.ThroughputCapacity = int32Ptr(4000)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when data_read_cache_configuration is set on SSD storage", func() {
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "NO_CACHE",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when per_unit_storage_throughput rides intelligent tiering", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ThroughputCapacity = int32Ptr(4000)
			spec.PerUnitStorageThroughput = int32Ptr(125)
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "PROPORTIONAL_TO_THROUGHPUT_CAPACITY",
			}
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: efa_requires_persistent_2_metadata
	// -------------------------------------------------------------------------

	ginkgo.Context("efa_requires_persistent_2_metadata", func() {
		ginkgo.It("fails when efa_enabled rides a SCRATCH deployment", func() {
			spec.EfaEnabled = true
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when efa_enabled lacks metadata configuration", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.PerUnitStorageThroughput = int32Ptr(500)
			spec.EfaEnabled = true
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: per_unit_throughput_requires_persistent
	// -------------------------------------------------------------------------

	ginkgo.Context("per_unit_throughput_requires_persistent", func() {
		ginkgo.It("fails when per_unit_storage_throughput is set with SCRATCH_2", func() {
			spec.DeploymentType = stringPtr("SCRATCH_2")
			spec.PerUnitStorageThroughput = int32Ptr(50)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when per_unit_storage_throughput is set without a deployment type (default SCRATCH_2)", func() {
			spec.PerUnitStorageThroughput = int32Ptr(50)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: legacy S3 link arm
	// -------------------------------------------------------------------------

	ginkgo.Context("legacy S3 link arm", func() {
		ginkgo.It("fails when import_path is set on PERSISTENT_2", func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.PerUnitStorageThroughput = int32Ptr(125)
			spec.ImportPath = "s3://my-bucket/data"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when export_path is set without import_path", func() {
			spec.DeploymentType = stringPtr("SCRATCH_2")
			spec.ExportPath = "s3://my-bucket/output"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when auto_import_policy is set without import_path", func() {
			spec.AutoImportPolicy = "NEW"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when auto_import_policy has an invalid value", func() {
			spec.ImportPath = "s3://my-bucket/data"
			spec.AutoImportPolicy = "ALL"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when imported_file_chunk_size is set without import_path", func() {
			spec.ImportedFileChunkSize = int32Ptr(1024)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when imported_file_chunk_size exceeds the maximum", func() {
			spec.ImportPath = "s3://my-bucket/data"
			spec.ImportedFileChunkSize = int32Ptr(512001)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: root squash formats
	// -------------------------------------------------------------------------

	ginkgo.Context("root squash formats", func() {
		ginkgo.It("fails when root_squash is not UID:GID", func() {
			spec.RootSquashConfiguration = &AwsFsxLustreFileSystemRootSquashConfiguration{
				RootSquash: "nobody",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a no-squash NID is malformed", func() {
			spec.RootSquashConfiguration = &AwsFsxLustreFileSystemRootSquashConfiguration{
				RootSquash:   "65534:65534",
				NoSquashNids: []string{"10.0.1.6"},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: data read cache nested rules
	// -------------------------------------------------------------------------

	ginkgo.Context("data read cache nested rules", func() {
		ginkgo.BeforeEach(func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ThroughputCapacity = int32Ptr(4000)
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
		})

		ginkgo.It("fails when sizing_mode is invalid", func() {
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "AUTOMATIC",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when size_gib rides a non-user-provisioned mode", func() {
			spec.DataReadCacheConfiguration = &AwsFsxLustreFileSystemDataReadCacheConfiguration{
				SizingMode: "NO_CACHE",
				SizeGib:    int32Ptr(64),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: metadata configuration nested rules
	// -------------------------------------------------------------------------

	ginkgo.Context("metadata configuration nested rules", func() {
		ginkgo.BeforeEach(func() {
			spec.DeploymentType = stringPtr("PERSISTENT_2")
			spec.PerUnitStorageThroughput = int32Ptr(125)
		})

		ginkgo.It("fails when metadata_configuration is set on SCRATCH_2", func() {
			spec.DeploymentType = stringPtr("SCRATCH_2")
			spec.PerUnitStorageThroughput = nil
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when metadata mode is invalid", func() {
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
				Mode: stringPtr("MANUAL"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when iops is set with AUTOMATIC mode", func() {
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
				Mode: stringPtr("AUTOMATIC"),
				Iops: int32Ptr(3000),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when iops is not a valid tier", func() {
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
				Mode: stringPtr("USER_PROVISIONED"),
				Iops: int32Ptr(2000),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts iops with USER_PROVISIONED mode", func() {
			spec.MetadataConfiguration = &AwsFsxLustreFileSystemMetadataConfiguration{
				Mode: stringPtr("USER_PROVISIONED"),
				Iops: int32Ptr(6000),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: log_level_valid (nested message)
	// -------------------------------------------------------------------------

	ginkgo.Context("log_level_valid", func() {
		ginkgo.It("fails when log level is invalid", func() {
			spec.LogConfiguration = &AwsFsxLustreFileSystemLogConfiguration{
				Level: stringPtr("DEBUG"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})

func boolPtr(b bool) *bool {
	return &b
}
