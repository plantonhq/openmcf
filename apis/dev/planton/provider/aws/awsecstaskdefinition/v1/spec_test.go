package awsecstaskdefinitionv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/planton/apis/dev/planton/shared"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestAwsEcsTaskDefinitionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsEcsTaskDefinitionSpec Validation Tests")
}

func literalRef(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func boolPtr(b bool) *bool { return &b }

// minimalValidTaskDefinition is the common case: one Fargate container with
// task-level sizing and the default logging wiring.
func minimalValidTaskDefinition() *AwsEcsTaskDefinition {
	return &AwsEcsTaskDefinition{
		ApiVersion: "aws.planton.dev/v1",
		Kind:       "AwsEcsTaskDefinition",
		Metadata: &shared.CloudResourceMetadata{
			Name: "api",
		},
		Spec: &AwsEcsTaskDefinitionSpec{
			Region:        "us-west-2",
			Cpu:           256,
			Memory:        512,
			ExecutionRole: literalRef("arn:aws:iam::123456789012:role/ecs-execution"),
			Containers: []*AwsEcsTaskDefinitionContainer{
				{
					Name:  "api",
					Image: "public.ecr.aws/nginx/nginx:stable",
				},
			},
		},
	}
}

var _ = ginkgo.Describe("AwsEcsTaskDefinitionSpec Validation Tests", func() {

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("aws_ecs_task_definition", func() {

			ginkgo.It("should not return a validation error for a minimal Fargate task", func() {
				err := protovalidate.Validate(minimalValidTaskDefinition())
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept an EC2-only task without task-level sizing", func() {
				input := minimalValidTaskDefinition()
				input.Spec.RequiresCompatibilities = []string{"EC2"}
				input.Spec.Cpu = 0
				input.Spec.Memory = 0
				input.Spec.NetworkMode = "bridge"
				input.Spec.Containers[0].Memory = 512
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept a sidecar with startup ordering and health checks", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].HealthCheck = &AwsEcsTaskDefinitionHealthCheck{
					Command: []string{"CMD-SHELL", "curl -f http://localhost:8080/healthz || exit 1"},
				}
				input.Spec.Containers = append(input.Spec.Containers, &AwsEcsTaskDefinitionContainer{
					Name:      "otel-collector",
					Image:     "public.ecr.aws/aws-observability/aws-otel-collector:latest",
					Essential: boolPtr(false),
					DependsOn: []*AwsEcsTaskDefinitionContainerDependency{
						{ContainerName: "api", Condition: "HEALTHY"},
					},
				})
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept named ports with app protocols", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].PortMappings = []*AwsEcsTaskDefinitionPortMapping{
					{ContainerPort: 8080, Name: "http", AppProtocol: "http"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept a Fargate task with logging disabled and no execution role", func() {
				input := minimalValidTaskDefinition()
				input.Spec.ExecutionRole = nil
				input.Spec.Logging = &AwsEcsTaskDefinitionLogging{Disabled: true}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept roles, ARM64, ephemeral storage, and an EFS volume", func() {
				input := minimalValidTaskDefinition()
				input.Spec.TaskRole = literalRef("arn:aws:iam::123456789012:role/api-task")
				input.Spec.RuntimePlatform = &AwsEcsTaskDefinitionRuntimePlatform{
					CpuArchitecture: "ARM64",
				}
				input.Spec.EphemeralStorageGib = 50
				input.Spec.Volumes = []*AwsEcsTaskDefinitionVolume{
					{
						Name: "shared-data",
						Efs: &AwsEcsTaskDefinitionEfsVolume{
							FileSystemId:  literalRef("fs-0123456789abcdef0"),
							AccessPointId: literalRef("fsap-0123456789abcdef0"),
						},
					},
				}
				input.Spec.Containers[0].MountPoints = []*AwsEcsTaskDefinitionMountPoint{
					{SourceVolume: "shared-data", ContainerPath: "/var/data"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept secrets, environment files, and a FireLens router", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].Secrets = map[string]string{
					"DB_PASSWORD": "arn:aws:secretsmanager:us-west-2:123456789012:secret:db-abc123",
				}
				input.Spec.Containers[0].EnvironmentFiles = []string{
					"arn:aws:s3:::my-config-bucket/api.env",
				}
				input.Spec.Containers[0].LogConfiguration = &AwsEcsTaskDefinitionLogConfiguration{
					LogDriver: "awsfirelens",
					Options:   map[string]string{"Name": "cloudwatch"},
				}
				input.Spec.Containers = append(input.Spec.Containers, &AwsEcsTaskDefinitionContainer{
					Name:      "log-router",
					Image:     "public.ecr.aws/aws-observability/aws-for-fluent-bit:stable",
					Essential: boolPtr(false),
					FirelensConfiguration: &AwsEcsTaskDefinitionFirelens{
						Type: "fluentbit",
					},
				})
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})

			ginkgo.It("should accept a restart policy and stop timeout", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].StopTimeoutSeconds = 90
				input.Spec.Containers[0].RestartPolicy = &AwsEcsTaskDefinitionRestartPolicy{
					Enabled:                     true,
					RestartAttemptPeriodSeconds: 120,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("aws_ecs_task_definition", func() {

			ginkgo.It("should return an error when region is empty", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Region = ""
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should return an error when no containers are declared", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should return an error when a container has no name", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].Name = ""
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should return an error when a container has no image", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].Image = ""
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a Fargate task without task-level sizing", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Cpu = 0
				input.Spec.Memory = 0
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a Fargate task on the awslogs default without an execution role", func() {
				input := minimalValidTaskDefinition()
				input.Spec.ExecutionRole = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject explicit awslogs without an execution role even when the default is disabled", func() {
				input := minimalValidTaskDefinition()
				input.Spec.ExecutionRole = nil
				input.Spec.Logging = &AwsEcsTaskDefinitionLogging{Disabled: true}
				input.Spec.Containers[0].LogConfiguration = &AwsEcsTaskDefinitionLogConfiguration{
					LogDriver: "awslogs",
					Options:   map[string]string{"awslogs-group": "/ecs/api"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a Fargate task with bridge networking", func() {
				input := minimalValidTaskDefinition()
				input.Spec.NetworkMode = "bridge"
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown network mode", func() {
				input := minimalValidTaskDefinition()
				input.Spec.NetworkMode = "overlay"
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown launch-type compatibility", func() {
				input := minimalValidTaskDefinition()
				input.Spec.RequiresCompatibilities = []string{"LAMBDA"}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject ephemeral storage below the 21 GiB floor", func() {
				input := minimalValidTaskDefinition()
				input.Spec.EphemeralStorageGib = 10
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a task where every container is non-essential", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].Essential = boolPtr(false)
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a memory reservation above the hard limit", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].Memory = 512
				input.Spec.Containers[0].MemoryReservation = 1024
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an out-of-range container port", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].PortMappings = []*AwsEcsTaskDefinitionPortMapping{
					{ContainerPort: 70000},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an app protocol on an unnamed port", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].PortMappings = []*AwsEcsTaskDefinitionPortMapping{
					{ContainerPort: 8080, AppProtocol: "grpc"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown port protocol", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].PortMappings = []*AwsEcsTaskDefinitionPortMapping{
					{ContainerPort: 8080, Protocol: "sctp"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a health check without a command", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].HealthCheck = &AwsEcsTaskDefinitionHealthCheck{}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown dependency condition", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].DependsOn = []*AwsEcsTaskDefinitionContainerDependency{
					{ContainerName: "api", Condition: "READY"},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a volume backed by both EFS and a host path", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Volumes = []*AwsEcsTaskDefinitionVolume{
					{
						Name:     "conflicted",
						HostPath: "/mnt/data",
						Efs: &AwsEcsTaskDefinitionEfsVolume{
							FileSystemId: literalRef("fs-0123456789abcdef0"),
						},
					},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an EFS volume without a file system id", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Volumes = []*AwsEcsTaskDefinitionVolume{
					{Name: "data", Efs: &AwsEcsTaskDefinitionEfsVolume{}},
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown CPU architecture", func() {
				input := minimalValidTaskDefinition()
				input.Spec.RuntimePlatform = &AwsEcsTaskDefinitionRuntimePlatform{
					CpuArchitecture: "RISCV64",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown OS family", func() {
				input := minimalValidTaskDefinition()
				input.Spec.RuntimePlatform = &AwsEcsTaskDefinitionRuntimePlatform{
					OperatingSystemFamily: "DARWIN",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown log driver", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].LogConfiguration = &AwsEcsTaskDefinitionLogConfiguration{
					LogDriver: "logstash",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject an unknown FireLens router type", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].FirelensConfiguration = &AwsEcsTaskDefinitionFirelens{
					Type: "vector",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a stop timeout beyond the Fargate cap", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].StopTimeoutSeconds = 121
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})

			ginkgo.It("should reject a restart-attempt period below 60 seconds", func() {
				input := minimalValidTaskDefinition()
				input.Spec.Containers[0].RestartPolicy = &AwsEcsTaskDefinitionRestartPolicy{
					Enabled:                     true,
					RestartAttemptPeriodSeconds: 30,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).NotTo(gomega.BeNil())
			})
		})
	})
})
