package awsbatchcomputeenvironmentv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBatchComputeEnvironmentSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBatchComputeEnvironmentSpec Validation Suite")
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func svRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalFargateSpec() *AwsBatchComputeEnvironmentSpec {
	return &AwsBatchComputeEnvironmentSpec{
		Region: "us-west-2",
		ComputeResources: &AwsBatchComputeResources{
			Type:             "FARGATE",
			MaxVcpus:         256,
			SubnetIds:        []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa"), svRef("subnet-bbb")},
			SecurityGroupIds: []*foreignkeyv1.StringValueOrRef{svRef("sg-111")},
		},
	}
}

func minimalEc2Spec() *AwsBatchComputeEnvironmentSpec {
	return &AwsBatchComputeEnvironmentSpec{
		Region: "us-west-2",
		ComputeResources: &AwsBatchComputeResources{
			Type:          "EC2",
			MaxVcpus:      256,
			SubnetIds:     []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
			InstanceTypes: []string{"optimal"},
			InstanceRole:  svRef("arn:aws:iam::123456789012:instance-profile/ecsInstanceRole"),
		},
	}
}

var _ = ginkgo.Describe("AwsBatchComputeEnvironmentSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with minimal Fargate configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalFargateSpec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with minimal EC2 configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalEc2Spec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a capacity-optimized SPOT configuration and no Spot fleet role", func() {
			ginkgo.It("should not return a validation error (Spot Fleet is only used by BEST_FIT)", func() {
				spec := &AwsBatchComputeEnvironmentSpec{
					Region: "us-west-2",
					ComputeResources: &AwsBatchComputeResources{
						Type:               "SPOT",
						MaxVcpus:           512,
						SubnetIds:          []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
						InstanceTypes:      []string{"m5.xlarge", "c5.xlarge"},
						InstanceRole:       svRef("arn:aws:iam::123456789012:instance-profile/ecsInstanceRole"),
						BidPercentage:      int32Ptr(60),
						AllocationStrategy: "SPOT_PRICE_CAPACITY_OPTIMIZED",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a BEST_FIT SPOT configuration carrying the Spot fleet role", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := &AwsBatchComputeEnvironmentSpec{
					Region: "us-west-2",
					ComputeResources: &AwsBatchComputeResources{
						Type:               "SPOT",
						MaxVcpus:           512,
						SubnetIds:          []*foreignkeyv1.StringValueOrRef{svRef("subnet-aaa")},
						InstanceTypes:      []string{"optimal"},
						InstanceRole:       svRef("arn:aws:iam::123456789012:instance-profile/ecsInstanceRole"),
						SpotIamFleetRole:   svRef("arn:aws:iam::123456789012:role/aws-ec2-spot-fleet-role"),
						AllocationStrategy: "BEST_FIT",
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with FARGATE_SPOT configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.Type = "FARGATE_SPOT"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with the ordered allocation strategies", func() {
			ginkgo.It("should accept BEST_FIT_PROGRESSIVE_ORDERED", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.AllocationStrategy = "BEST_FIT_PROGRESSIVE_ORDERED"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("should accept SPOT_CAPACITY_OPTIMIZED_PRIORITIZED on SPOT", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Type = "SPOT"
				spec.ComputeResources.AllocationStrategy = "SPOT_CAPACITY_OPTIMIZED_PRIORITIZED"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with update policy", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.UpdatePolicy = &AwsBatchUpdatePolicy{
					TerminateJobsOnUpdate:      true,
					JobExecutionTimeoutMinutes: int32Ptr(30),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with launch template reference and placement group", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.LaunchTemplate = &AwsBatchLaunchTemplate{
					LaunchTemplateId: svRef("lt-0123456789abcdef0"),
					Version:          "$Latest",
				}
				spec.ComputeResources.PlacementGroup = "hpc-cluster-pg"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with ec2_configurations including a Kubernetes version", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Ec2Configurations = []*AwsBatchEc2Configuration{
					{ImageType: "EKS_AL2023", ImageKubernetesVersion: "1.31"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with an EKS configuration", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.EksConfiguration = &AwsBatchEksConfiguration{
					EksClusterArn:       svRef("arn:aws:eks:us-west-2:123456789012:cluster/batch-eks"),
					KubernetesNamespace: "batch-jobs",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with state set to DISABLED", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalFargateSpec()
				spec.State = stringPtr("DISABLED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with resource tags on an EC2 environment", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.ResourceTags = map[string]string{"team": "data-eng"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no compute_resources", func() {
			ginkgo.It("should return a validation error", func() {
				spec := &AwsBatchComputeEnvironmentSpec{Region: "us-west-2"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid state", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.State = stringPtr("PAUSED")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a missing compute resource type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.Type = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid compute resource type", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.Type = "INVALID_TYPE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with max_vcpus less than 1", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.MaxVcpus = 0
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with no subnets", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.SubnetIds = []*foreignkeyv1.StringValueOrRef{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with EC2 type missing instance_role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.InstanceRole = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a BEST_FIT SPOT environment missing the Spot fleet role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Type = "SPOT"
				spec.ComputeResources.AllocationStrategy = "BEST_FIT"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a strategy-less SPOT environment missing the Spot fleet role", func() {
			ginkgo.It("should return a validation error (no strategy defaults to BEST_FIT)", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Type = "SPOT"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Fargate environment missing security groups", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.SecurityGroupIds = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with EC2-only fields on a Fargate environment", func() {
			ginkgo.It("should reject instance_types", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.InstanceTypes = []string{"optimal"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject a launch template", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.LaunchTemplate = &AwsBatchLaunchTemplate{
					LaunchTemplateId: svRef("lt-0123456789abcdef0"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject an allocation strategy", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.AllocationStrategy = "BEST_FIT_PROGRESSIVE"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject resource tags", func() {
				spec := minimalFargateSpec()
				spec.ComputeResources.ResourceTags = map[string]string{"team": "data-eng"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with Spot-only fields on an EC2 environment", func() {
			ginkgo.It("should reject bid_percentage", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.BidPercentage = int32Ptr(60)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject spot_iam_fleet_role", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.SpotIamFleetRole = svRef("arn:aws:iam::123456789012:role/spot-fleet-role")
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with invalid allocation_strategy", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.AllocationStrategy = "INVALID_STRATEGY"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with bid_percentage out of range", func() {
			ginkgo.It("should return a validation error when > 100", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Type = "SPOT"
				spec.ComputeResources.AllocationStrategy = "SPOT_CAPACITY_OPTIMIZED"
				spec.ComputeResources.BidPercentage = int32Ptr(150)
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with launch template missing its reference", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.LaunchTemplate = &AwsBatchLaunchTemplate{}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with ec2_configurations exceeding max 2", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.ComputeResources.Ec2Configurations = []*AwsBatchEc2Configuration{
					{ImageType: "ECS_AL2"},
					{ImageType: "ECS_AL2023"},
					{ImageType: "ECS_AL2_NVIDIA"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an EKS configuration missing its cluster reference", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.EksConfiguration = &AwsBatchEksConfiguration{
					KubernetesNamespace: "batch-jobs",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an EKS configuration missing its namespace", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.EksConfiguration = &AwsBatchEksConfiguration{
					EksClusterArn: svRef("arn:aws:eks:us-west-2:123456789012:cluster/batch-eks"),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with job_execution_timeout_minutes out of range", func() {
			ginkgo.It("should return a validation error when > 360", func() {
				spec := minimalEc2Spec()
				spec.UpdatePolicy = &AwsBatchUpdatePolicy{
					TerminateJobsOnUpdate:      true,
					JobExecutionTimeoutMinutes: int32Ptr(500),
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
