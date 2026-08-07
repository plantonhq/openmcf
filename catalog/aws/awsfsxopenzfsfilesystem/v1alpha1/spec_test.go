package awsfsxopenzfsfilesystemv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsFsxOpenzfsFileSystemSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsFsxOpenzfsFileSystemSpec Validation Suite")
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

// multiAzBase returns a valid MULTI_AZ_1 baseline: two subnets, a preferred
// subnet, and a gen-2 throughput value.
func multiAzBase() *AwsFsxOpenzfsFileSystemSpec {
	return &AwsFsxOpenzfsFileSystemSpec{
		Region:         "us-west-2",
		DeploymentType: stringPtr("MULTI_AZ_1"),
		SubnetIds: []*foreignkeyv1.StringValueOrRef{
			strRef("subnet-az1"),
			strRef("subnet-az2"),
		},
		PreferredSubnetId:  strRef("subnet-az1"),
		StorageCapacityGib: int32Ptr(256),
		ThroughputCapacity: 160,
	}
}

var _ = ginkgo.Describe("AwsFsxOpenzfsFileSystemSpec validations", func() {
	var spec *AwsFsxOpenzfsFileSystemSpec

	ginkgo.BeforeEach(func() {
		// Minimal valid spec: one subnet (SINGLE_AZ_2 default), provisioned SSD
		// capacity, gen-2 throughput.
		spec = &AwsFsxOpenzfsFileSystemSpec{
			Region:             "us-west-2",
			SubnetIds:          []*foreignkeyv1.StringValueOrRef{strRef("subnet-abc123")},
			StorageCapacityGib: int32Ptr(256),
			ThroughputCapacity: 160,
		}
	})

	// -------------------------------------------------------------------------
	// Happy path
	// -------------------------------------------------------------------------

	ginkgo.It("accepts a minimal valid spec", func() {
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts SINGLE_AZ_1 with a gen-1 throughput value", func() {
		spec.DeploymentType = stringPtr("SINGLE_AZ_1")
		spec.ThroughputCapacity = 64
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts the HA deployment variants", func() {
		spec.DeploymentType = stringPtr("SINGLE_AZ_HA_2")
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())

		spec.DeploymentType = stringPtr("SINGLE_AZ_HA_1")
		spec.ThroughputCapacity = 64
		err = protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full MULTI_AZ_1 shape", func() {
		spec = multiAzBase()
		spec.EndpointIpAddressRange = "10.0.255.0/24"
		spec.RouteTableIds = []*foreignkeyv1.StringValueOrRef{strRef("rtb-abc123")}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts INTELLIGENT_TIERING on MULTI_AZ_1 with a read cache", func() {
		spec = multiAzBase()
		spec.StorageType = stringPtr("INTELLIGENT_TIERING")
		spec.StorageCapacityGib = nil
		spec.ReadCacheConfiguration = &AwsFsxOpenzfsFileSystemReadCacheConfiguration{
			SizingMode: "PROPORTIONAL_TO_THROUGHPUT_CAPACITY",
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a backup restore without storage capacity", func() {
		spec.StorageCapacityGib = nil
		spec.BackupId = "backup-0123456789abcdef0"
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts a full root volume configuration", func() {
		spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
			DataCompressionType: stringPtr("ZSTD"),
			RecordSizeKib:       int32Ptr(1024),
			ReadOnly:            false,
			CopyTagsToSnapshots: true,
			NfsExports: &AwsFsxOpenzfsFileSystemNfsExports{
				ClientConfigurations: []*AwsFsxOpenzfsFileSystemNfsClientConfiguration{
					{Clients: "10.0.0.0/16", Options: []string{"rw", "crossmnt"}},
				},
			},
			UserAndGroupQuotas: []*AwsFsxOpenzfsFileSystemUserAndGroupQuota{
				{Id: 1000, StorageCapacityQuotaGib: 100, Type: "USER"},
				{Id: 2000, StorageCapacityQuotaGib: 500, Type: "GROUP"},
			},
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts user-provisioned disk IOPS", func() {
		spec.DiskIopsConfiguration = &AwsFsxOpenzfsFileSystemDiskIopsConfiguration{
			Mode: stringPtr("USER_PROVISIONED"),
			Iops: int32Ptr(50000),
		}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.It("accepts delete options and final backup tags", func() {
		spec.DeleteOptions = []string{"DELETE_CHILD_VOLUMES_AND_SNAPSHOTS"}
		spec.SkipFinalBackup = boolPtr(false)
		spec.FinalBackupTags = map[string]string{"retention": "audit"}
		err := protovalidate.Validate(spec)
		gomega.Expect(err).To(gomega.BeNil())
	})

	// -------------------------------------------------------------------------
	// Field-level validation failures
	// -------------------------------------------------------------------------

	ginkgo.Context("field-level validations", func() {
		ginkgo.It("fails when subnet_ids is empty", func() {
			spec.SubnetIds = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when storage_capacity_gib is below minimum", func() {
			spec.StorageCapacityGib = int32Ptr(32)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when storage_capacity_gib exceeds the maximum", func() {
			spec.StorageCapacityGib = int32Ptr(524289)
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when automatic_backup_retention_days exceeds 90", func() {
			spec.AutomaticBackupRetentionDays = int32Ptr(91)
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

		ginkgo.It("fails when more than 50 security groups are supplied", func() {
			for i := 0; i < 51; i++ {
				spec.SecurityGroupIds = append(spec.SecurityGroupIds, strRef("sg-x"))
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: deployment/storage types
	// -------------------------------------------------------------------------

	ginkgo.Context("type enums", func() {
		ginkgo.It("fails when deployment_type is invalid", func() {
			spec.DeploymentType = stringPtr("MULTI_AZ_2")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when storage_type is invalid", func() {
			spec.StorageType = stringPtr("HDD")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: throughput value sets
	// -------------------------------------------------------------------------

	ginkgo.Context("throughput_capacity value sets", func() {
		ginkgo.It("fails when a gen-1 value rides the default SINGLE_AZ_2", func() {
			spec.ThroughputCapacity = 64
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a gen-2 value rides SINGLE_AZ_1", func() {
			spec.DeploymentType = stringPtr("SINGLE_AZ_1")
			spec.ThroughputCapacity = 160
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when an arbitrary value rides MULTI_AZ_1", func() {
			spec = multiAzBase()
			spec.ThroughputCapacity = 200
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: INTELLIGENT_TIERING contract
	// -------------------------------------------------------------------------

	ginkgo.Context("INTELLIGENT_TIERING contract", func() {
		ginkgo.It("fails on a single-AZ deployment", func() {
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ReadCacheConfiguration = &AwsFsxOpenzfsFileSystemReadCacheConfiguration{
				SizingMode: "NO_CACHE",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails without a read cache decision", func() {
			spec = multiAzBase()
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when provisioned capacity rides intelligent tiering", func() {
			spec = multiAzBase()
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.ReadCacheConfiguration = &AwsFsxOpenzfsFileSystemReadCacheConfiguration{
				SizingMode: "PROPORTIONAL_TO_THROUGHPUT_CAPACITY",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a read cache rides SSD storage", func() {
			spec.ReadCacheConfiguration = &AwsFsxOpenzfsFileSystemReadCacheConfiguration{
				SizingMode: "NO_CACHE",
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when the read cache size rides a non-user-provisioned mode", func() {
			spec = multiAzBase()
			spec.StorageType = stringPtr("INTELLIGENT_TIERING")
			spec.StorageCapacityGib = nil
			spec.ReadCacheConfiguration = &AwsFsxOpenzfsFileSystemReadCacheConfiguration{
				SizingMode: "NO_CACHE",
				SizeGib:    int32Ptr(64),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: storage capacity presence
	// -------------------------------------------------------------------------

	ginkgo.Context("storage capacity presence", func() {
		ginkgo.It("fails when SSD capacity is omitted without a backup", func() {
			spec.StorageCapacityGib = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a backup restore also provisions capacity", func() {
			spec.BackupId = "backup-0123456789abcdef0"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: MULTI_AZ_1 networking contract
	// -------------------------------------------------------------------------

	ginkgo.Context("MULTI_AZ_1 networking contract", func() {
		ginkgo.It("fails when MULTI_AZ_1 omits preferred_subnet_id", func() {
			spec = multiAzBase()
			spec.PreferredSubnetId = nil
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when preferred_subnet_id rides a single-AZ deployment", func() {
			spec.PreferredSubnetId = strRef("subnet-abc123")
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when MULTI_AZ_1 has only one subnet", func() {
			spec = multiAzBase()
			spec.SubnetIds = spec.SubnetIds[:1]
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a single-AZ deployment has two subnets", func() {
			spec.SubnetIds = append(spec.SubnetIds, strRef("subnet-def456"))
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when endpoint_ip_address_range rides a single-AZ deployment", func() {
			spec.EndpointIpAddressRange = "10.0.255.0/24"
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when route_table_ids ride a single-AZ deployment", func() {
			spec.RouteTableIds = []*foreignkeyv1.StringValueOrRef{strRef("rtb-abc123")}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: delete options
	// -------------------------------------------------------------------------

	ginkgo.Context("delete_options", func() {
		ginkgo.It("fails on an unknown delete option", func() {
			spec.DeleteOptions = []string{"DELETE_EVERYTHING"}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: disk IOPS
	// -------------------------------------------------------------------------

	ginkgo.Context("disk IOPS", func() {
		ginkgo.It("fails when iops rides AUTOMATIC mode", func() {
			spec.DiskIopsConfiguration = &AwsFsxOpenzfsFileSystemDiskIopsConfiguration{
				Mode: stringPtr("AUTOMATIC"),
				Iops: int32Ptr(50000),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when iops exceeds the SINGLE_AZ_1 ceiling", func() {
			spec.DeploymentType = stringPtr("SINGLE_AZ_1")
			spec.ThroughputCapacity = 64
			spec.DiskIopsConfiguration = &AwsFsxOpenzfsFileSystemDiskIopsConfiguration{
				Mode: stringPtr("USER_PROVISIONED"),
				Iops: int32Ptr(160001),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when iops exceeds the SINGLE_AZ_2 ceiling", func() {
			spec.DiskIopsConfiguration = &AwsFsxOpenzfsFileSystemDiskIopsConfiguration{
				Mode: stringPtr("USER_PROVISIONED"),
				Iops: int32Ptr(400001),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	// -------------------------------------------------------------------------
	// CEL: root volume nested rules
	// -------------------------------------------------------------------------

	ginkgo.Context("root volume nested rules", func() {
		ginkgo.It("fails when data_compression_type is invalid", func() {
			spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
				DataCompressionType: stringPtr("GZIP"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when record_size_kib is invalid", func() {
			spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
				RecordSizeKib: int32Ptr(100),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when nfs_exports has no client configurations", func() {
			spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
				NfsExports: &AwsFsxOpenzfsFileSystemNfsExports{},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when an NFS client configuration has no options", func() {
			spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
				NfsExports: &AwsFsxOpenzfsFileSystemNfsExports{
					ClientConfigurations: []*AwsFsxOpenzfsFileSystemNfsClientConfiguration{
						{Clients: "*"},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a quota has an invalid type", func() {
			spec.RootVolumeConfiguration = &AwsFsxOpenzfsFileSystemRootVolumeConfiguration{
				UserAndGroupQuotas: []*AwsFsxOpenzfsFileSystemUserAndGroupQuota{
					{Id: 1000, StorageCapacityQuotaGib: 100, Type: "ROLE"},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})
})

func boolPtr(b bool) *bool {
	return &b
}
