package awsmemorydbclusterv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsMemorydbClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsMemorydbClusterSpec Validation Suite")
}

func int32Ptr(i int32) *int32 {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func svr(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

// minimalCluster is the smallest valid manifest: engine, node type, and the
// mandatory (deliberately explicit) ACL attachment.
func minimalCluster() *AwsMemorydbClusterSpec {
	return &AwsMemorydbClusterSpec{
		Region:   "us-west-2",
		Engine:   "redis",
		NodeType: "db.t4g.small",
		AclName:  svr("open-access"),
	}
}

var _ = ginkgo.Describe("AwsMemorydbClusterSpec validations", func() {

	// -----------------------------------------------------------------
	// Valid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal required fields", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalCluster())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with full production configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsMemorydbClusterSpec{
					Region:              "us-west-2",
					Engine:              "redis",
					EngineVersion:       "7.1",
					Description:         "Production session store",
					NodeType:            "db.r7g.large",
					Port:                int32Ptr(6379),
					NumShards:           int32Ptr(2),
					NumReplicasPerShard: int32Ptr(2),
					AclName:             svr("my-prod-acl"),
					SubnetIds: []*foreignkeyv1.StringValueOrRef{
						svr("subnet-11111111"),
						svr("subnet-22222222"),
					},
					SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{
						svr("sg-12345678"),
					},
					TlsEnabled:             boolPtr(true),
					KmsKeyArn:              svr("arn:aws:kms:us-east-1:123456789012:key/abc-123"),
					MaintenanceWindow:      "sun:05:00-sun:06:00",
					SnapshotRetentionLimit: 7,
					SnapshotWindow:         "03:00-04:00",
					FinalSnapshotName:      "final-snap",
					ParameterGroupFamily:   "memorydb_redis7",
					Parameters: []*AwsMemorydbClusterParameter{
						{Name: "activedefrag", Value: "yes"},
					},
					SnsTopicArn:             svr("arn:aws:sns:us-east-1:123456789012:alerts"),
					AutoMinorVersionUpgrade: boolPtr(true),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with valkey engine", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.Engine = "valkey"
				spec.NodeType = "db.r7g.large"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an ACL reference instead of a literal", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.AclName = &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "payments-env-acl"},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an existing subnet group by name", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.SubnetGroupName = "shared-data-subnets"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an existing parameter group by name", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.ParameterGroupName = "shared-tuning"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with dual-stack networking and IPv6 discovery", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.NetworkType = "dual_stack"
				spec.IpDiscovery = "ipv6"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a multi-region cluster join", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.NodeType = "db.r7g.xlarge"
				spec.MultiRegionClusterName = "virxk-global-sessions"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with snapshot restore from ARNs", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotArns = []string{"arn:aws:s3:::my-bucket/snapshot.rdb"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with data tiering enabled", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalCluster()
				spec.NodeType = "db.r6gd.xlarge"
				spec.DataTiering = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	// -----------------------------------------------------------------
	// Invalid inputs
	// -----------------------------------------------------------------
	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with missing engine", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.Engine = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing node_type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.NodeType = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with missing acl_name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.AclName = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid engine value", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.Engine = "memcached"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with both subnet_ids and subnet_group_name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SubnetIds = []*foreignkeyv1.StringValueOrRef{svr("subnet-11111111")}
				spec.SubnetGroupName = "shared-data-subnets"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with both parameters and parameter_group_name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.ParameterGroupFamily = "memorydb_redis7"
				spec.Parameters = []*AwsMemorydbClusterParameter{
					{Name: "activedefrag", Value: "yes"},
				}
				spec.ParameterGroupName = "shared-tuning"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with parameters but no parameter_group_family", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.Parameters = []*AwsMemorydbClusterParameter{
					{Name: "activedefrag", Value: "yes"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with both snapshot_arns and snapshot_name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotArns = []string{"arn:aws:s3:::bucket/snap.rdb"}
				spec.SnapshotName = "my-snapshot"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a snapshot ARN that is not an S3 object ARN", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotArns = []string{"arn:aws:elasticache:us-west-2:123456789012:snapshot/my-snap"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a snapshot ARN containing a comma", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotArns = []string{"arn:aws:s3:::bucket/snap,shard2.rdb"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a final_snapshot_name containing consecutive hyphens", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.FinalSnapshotName = "final--snap"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a final_snapshot_name ending in a hyphen", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.FinalSnapshotName = "final-snap-"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a final_snapshot_name containing uppercase characters", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.FinalSnapshotName = "Final-Snap"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid network_type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.NetworkType = "dualstack"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with ip_discovery ipv6 on an ipv4 network", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.IpDiscovery = "ipv6"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid port (out of range)", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.Port = int32Ptr(0)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with num_replicas_per_shard out of range", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.NumReplicasPerShard = int32Ptr(6)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with snapshot_retention_limit out of range", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotRetentionLimit = 36
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid maintenance_window format", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.MaintenanceWindow = "invalid-format"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid snapshot_window format", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.SnapshotWindow = "bad-window"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with parameter missing name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalCluster()
				spec.ParameterGroupFamily = "memorydb_redis7"
				spec.Parameters = []*AwsMemorydbClusterParameter{
					{Value: "yes"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
