package awsecsclusterv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	fk "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsEcsClusterSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEcsClusterSpec Validation Suite")
}

func literal(v string) *fk.StringValueOrRef {
	return &fk.StringValueOrRef{LiteralOrRef: &fk.StringValueOrRef_Value{Value: v}}
}

var _ = ginkgo.Describe("AwsEcsClusterSpec validations", func() {
	var spec *AwsEcsClusterSpec

	ginkgo.BeforeEach(func() {
		spec = &AwsEcsClusterSpec{
			Region:            "us-west-2",
			ContainerInsights: "enabled",
			CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
			DefaultCapacityProviderStrategy: []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "FARGATE", Base: 1, Weight: 1},
				{CapacityProvider: "FARGATE_SPOT", Weight: 4},
			},
		}
	})

	ginkgo.It("accepts a Fargate cost-optimized cluster", func() {
		gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
	})

	ginkgo.It("fails when region is empty", func() {
		spec.Region = ""
		gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
	})

	ginkgo.Context("container insights", func() {
		ginkgo.It("accepts enhanced observability", func() {
			spec.ContainerInsights = "enhanced"
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid value", func() {
			spec.ContainerInsights = "on"
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts unset (account default)", func() {
			spec.ContainerInsights = ""
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("capacity providers", func() {
		ginkgo.It("fails on a custom name in the built-in list", func() {
			spec.CapacityProviders = []string{"FARGATE", "my-asg-provider"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on duplicate built-ins", func() {
			spec.CapacityProviders = []string{"FARGATE", "FARGATE"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts an EC2 capacity provider with managed scaling", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "general-purpose",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
				ManagedScaling: &AwsEcsClusterManagedScaling{
					Status:         "ENABLED",
					TargetCapacity: 80,
				},
				ManagedDraining: "ENABLED",
			}}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when an EC2 provider has no auto-scaling group", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name: "general-purpose",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on a reserved provider name prefix", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "ecs-workers",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on duplicate EC2 provider names", func() {
			provider := func(name string) *AwsEcsClusterEc2CapacityProvider {
				return &AwsEcsClusterEc2CapacityProvider{
					Name:                name,
					AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
				}
			}
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{provider("workers"), provider("workers")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid managed_termination_protection value", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                         "workers",
				AutoScalingGroupArn:          literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
				ManagedTerminationProtection: "on",
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an out-of-range target capacity", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "workers",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
				ManagedScaling:      &AwsEcsClusterManagedScaling{TargetCapacity: 150},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when scaling step sizes are inverted", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "workers",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
				ManagedScaling: &AwsEcsClusterManagedScaling{
					MinimumScalingStepSize: 10,
					MaximumScalingStepSize: 5,
				},
			}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("default capacity provider strategy", func() {
		ginkgo.It("accepts a strategy naming a folded EC2 provider", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "workers",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
			}}
			spec.DefaultCapacityProviderStrategy = []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "workers", Weight: 1},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails when a strategy names an unassociated provider", func() {
			spec.DefaultCapacityProviderStrategy = []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "ghost-provider", Weight: 1},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when two entries set a non-zero base", func() {
			spec.DefaultCapacityProviderStrategy = []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "FARGATE", Base: 1, Weight: 1},
				{CapacityProvider: "FARGATE_SPOT", Base: 1, Weight: 4},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on a negative base", func() {
			spec.DefaultCapacityProviderStrategy = []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "FARGATE", Base: -1},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("execute command configuration", func() {
		ginkgo.It("accepts DEFAULT logging with a KMS key", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{
				Logging:  "DEFAULT",
				KmsKeyId: literal("arn:aws:kms:us-west-2:123456789012:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid logging value", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{Logging: "ALL"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when OVERRIDE logging has no log configuration", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{Logging: "OVERRIDE"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a log configuration is set without OVERRIDE", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{
				Logging: "DEFAULT",
				LogConfiguration: &AwsEcsClusterExecuteCommandLogConfiguration{
					CloudWatchLogGroupName: "/ecs/exec-audit",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when an OVERRIDE log configuration has no destination", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{
				Logging:          "OVERRIDE",
				LogConfiguration: &AwsEcsClusterExecuteCommandLogConfiguration{},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts OVERRIDE logging to CloudWatch and S3", func() {
			spec.ExecuteCommandConfiguration = &AwsEcsClusterExecuteCommandConfiguration{
				Logging: "OVERRIDE",
				LogConfiguration: &AwsEcsClusterExecuteCommandLogConfiguration{
					CloudWatchLogGroupName:      "/ecs/exec-audit",
					CloudWatchEncryptionEnabled: true,
					S3BucketName:                "exec-audit-archive",
					S3KeyPrefix:                 "prod/",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("managed storage", func() {
		ginkgo.It("accepts Fargate ephemeral storage encryption", func() {
			spec.ManagedStorageConfiguration = &AwsEcsClusterManagedStorageConfiguration{
				FargateEphemeralStorageKmsKeyId: literal("arn:aws:kms:us-west-2:123456789012:key/abc"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("managed instances capacity providers", func() {
		// A minimal valid managed-instances provider: role, profile,
		// network, and the two required requirement dimensions.
		minimalManagedInstances := func(name string) *AwsEcsClusterManagedInstancesCapacityProvider {
			return &AwsEcsClusterManagedInstancesCapacityProvider{
				Name:                  name,
				InfrastructureRoleArn: literal("arn:aws:iam::123456789012:role/ecs-mi-infrastructure"),
				InstanceLaunchTemplate: &AwsEcsClusterManagedInstancesLaunchTemplate{
					Ec2InstanceProfileArn: literal("arn:aws:iam::123456789012:instance-profile/ecs-mi-instances"),
					NetworkConfiguration: &AwsEcsClusterManagedInstancesNetworkConfiguration{
						Subnets: []*fk.StringValueOrRef{literal("subnet-0a1b2c3d4e5f60789")},
					},
					InstanceRequirements: &AwsEcsClusterManagedInstancesRequirements{
						MemoryMib: &AwsEcsClusterIntRange{Min: 2048},
						VcpuCount: &AwsEcsClusterIntRange{Min: 1},
					},
				},
			}
		}

		ginkgo.It("accepts a minimal managed-instances provider", func() {
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{
				minimalManagedInstances("mi-general"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("fails without the infrastructure role", func() {
			p := minimalManagedInstances("mi-general")
			p.InfrastructureRoleArn = nil
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails without a network configuration", func() {
			p := minimalManagedInstances("mi-general")
			p.InstanceLaunchTemplate.NetworkConfiguration = nil
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails without subnets", func() {
			p := minimalManagedInstances("mi-general")
			p.InstanceLaunchTemplate.NetworkConfiguration.Subnets = nil
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on a reserved name prefix", func() {
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{
				minimalManagedInstances("ecs-managed"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on duplicate managed-instances names", func() {
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{
				minimalManagedInstances("mi-general"),
				minimalManagedInstances("mi-general"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails when a managed-instances name collides with an ec2 provider name", func() {
			spec.Ec2CapacityProviders = []*AwsEcsClusterEc2CapacityProvider{{
				Name:                "general-purpose",
				AutoScalingGroupArn: literal("arn:aws:autoscaling:us-west-2:123456789012:autoScalingGroup:uuid:autoScalingGroupName/workers"),
			}}
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{
				minimalManagedInstances("general-purpose"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("lets the default strategy name a managed-instances provider", func() {
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{
				minimalManagedInstances("mi-general"),
			}
			spec.DefaultCapacityProviderStrategy = []*AwsEcsClusterCapacityProviderStrategy{
				{CapacityProvider: "mi-general", Weight: 1},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.Context("capacity reservations", func() {
			ginkgo.It("fails when RESERVED has no capacity_reservations", func() {
				p := minimalManagedInstances("mi-reserved")
				p.InstanceLaunchTemplate.CapacityOptionType = "RESERVED"
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("fails when capacity_reservations is set without RESERVED", func() {
				p := minimalManagedInstances("mi-spot")
				p.InstanceLaunchTemplate.CapacityOptionType = "SPOT"
				p.InstanceLaunchTemplate.CapacityReservations = &AwsEcsClusterManagedInstancesCapacityReservations{
					ReservationPreference: "RESERVATIONS_FIRST",
				}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("accepts RESERVED with reservations and requirements", func() {
				p := minimalManagedInstances("mi-reserved")
				p.InstanceLaunchTemplate.CapacityOptionType = "RESERVED"
				p.InstanceLaunchTemplate.CapacityReservations = &AwsEcsClusterManagedInstancesCapacityReservations{
					ReservationPreference: "RESERVATIONS_FIRST",
				}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
			})

			ginkgo.It("fails when reservation_group_arn is set without RESERVATIONS_ONLY", func() {
				p := minimalManagedInstances("mi-reserved")
				p.InstanceLaunchTemplate.CapacityOptionType = "RESERVED"
				p.InstanceLaunchTemplate.CapacityReservations = &AwsEcsClusterManagedInstancesCapacityReservations{
					ReservationPreference: "RESERVATIONS_FIRST",
					ReservationGroupArn:   "arn:aws:resource-groups:us-west-2:123456789012:group/cr-group",
				}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("fails when a reservations preference has no instance_requirements", func() {
				p := minimalManagedInstances("mi-reserved")
				p.InstanceLaunchTemplate.CapacityOptionType = "RESERVED"
				p.InstanceLaunchTemplate.InstanceRequirements = nil
				p.InstanceLaunchTemplate.CapacityReservations = &AwsEcsClusterManagedInstancesCapacityReservations{
					ReservationPreference: "RESERVATIONS_ONLY",
				}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.Context("instance requirements", func() {
			ginkgo.It("fails without the required memory and vcpu dimensions", func() {
				p := minimalManagedInstances("mi-general")
				p.InstanceLaunchTemplate.InstanceRequirements = &AwsEcsClusterManagedInstancesRequirements{}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("fails on an inverted range", func() {
				p := minimalManagedInstances("mi-general")
				p.InstanceLaunchTemplate.InstanceRequirements.MemoryMib = &AwsEcsClusterIntRange{Min: 4096, Max: 2048}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})

			ginkgo.It("fails when allow and deny lists are both set", func() {
				p := minimalManagedInstances("mi-general")
				p.InstanceLaunchTemplate.InstanceRequirements.AllowedInstanceTypes = []string{"m5.*"}
				p.InstanceLaunchTemplate.InstanceRequirements.ExcludedInstanceTypes = []string{"t2.*"}
				spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
				gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
			})
		})

		ginkgo.It("fails on an out-of-range scale_in_after_seconds", func() {
			p := minimalManagedInstances("mi-general")
			scaleIn := int32(7200)
			p.ScaleInAfterSeconds = &scaleIn
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("fails on an invalid monitoring value", func() {
			p := minimalManagedInstances("mi-general")
			p.InstanceLaunchTemplate.Monitoring = "VERBOSE"
			spec.ManagedInstancesCapacityProviders = []*AwsEcsClusterManagedInstancesCapacityProvider{p}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
