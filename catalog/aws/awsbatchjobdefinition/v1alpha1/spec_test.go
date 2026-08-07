package awsbatchjobdefinitionv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

func TestAwsBatchJobDefinitionSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "AwsBatchJobDefinitionSpec Validation Suite")
}

func svRef(val string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: val},
	}
}

func minimalEc2Spec() *AwsBatchJobDefinitionSpec {
	return &AwsBatchJobDefinitionSpec{
		Region: "us-west-2",
		Container: &AwsBatchJobDefinitionContainer{
			Image:     "public.ecr.aws/amazonlinux/amazonlinux:2023",
			Command:   []string{"echo", "hello"},
			Vcpus:     1,
			MemoryMib: 2048,
		},
	}
}

func minimalFargateSpec() *AwsBatchJobDefinitionSpec {
	return &AwsBatchJobDefinitionSpec{
		Region:               "us-west-2",
		PlatformCapabilities: []string{"FARGATE"},
		Container: &AwsBatchJobDefinitionContainer{
			Image:         "public.ecr.aws/amazonlinux/amazonlinux:2023",
			Command:       []string{"echo", "hello"},
			Vcpus:         0.25,
			MemoryMib:     512,
			ExecutionRole: svRef("arn:aws:iam::123456789012:role/batch-execution-role"),
		},
	}
}

var _ = ginkgo.Describe("AwsBatchJobDefinitionSpec validations", func() {

	ginkgo.Describe("When valid input is passed", func() {

		ginkgo.Context("with a minimal EC2 container job", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalEc2Spec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with a minimal Fargate container job", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(minimalFargateSpec())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with Fargate-only knobs on a Fargate job", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalFargateSpec()
				spec.Container.FargatePlatformVersion = "1.4.0"
				spec.Container.AssignPublicIp = true
				spec.Container.EphemeralStorageGib = 50
				spec.Container.RuntimePlatform = &AwsBatchJobDefinitionRuntimePlatform{
					CpuArchitecture:       "ARM64",
					OperatingSystemFamily: "LINUX",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with EC2-only knobs on an EC2 job", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.Gpus = 1
				spec.Container.Privileged = true
				spec.Container.Ulimits = []*AwsBatchJobDefinitionUlimit{
					{Name: "nofile", SoftLimit: 8192, HardLimit: 65535},
				}
				spec.Container.LinuxParameters = &AwsBatchJobDefinitionLinuxParameters{
					InitProcessEnabled:  true,
					SharedMemorySizeMib: 256,
					Tmpfs: []*AwsBatchJobDefinitionTmpfs{
						{ContainerPath: "/tmp/scratch", SizeMib: 128, MountOptions: []string{"noexec"}},
					},
					Devices: []*AwsBatchJobDefinitionDevice{
						{HostPath: "/dev/xvdf", Permissions: []string{"READ", "WRITE"}},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with identities, environment, secrets, and logging", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.JobRole = svRef("arn:aws:iam::123456789012:role/etl-job-role")
				spec.Container.ExecutionRole = svRef("arn:aws:iam::123456789012:role/batch-execution-role")
				spec.Container.Environment = map[string]string{"STAGE": "prod"}
				spec.Container.Secrets = map[string]string{
					"DB_PASSWORD": "arn:aws:secretsmanager:us-west-2:123456789012:secret:db-pass-AbCdEf",
				}
				spec.Container.LogConfiguration = &AwsBatchJobDefinitionLogConfiguration{
					LogDriver: "awslogs",
					Options:   map[string]string{"awslogs-group": "/custom/batch"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with EFS and host-path volumes and matching mount points", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.Volumes = []*AwsBatchJobDefinitionVolume{
					{
						Name: "shared-data",
						Efs: &AwsBatchJobDefinitionEfsVolume{
							FileSystemId:     svRef("fs-0123456789abcdef0"),
							AccessPointId:    svRef("fsap-0123456789abcdef0"),
							IamAuthorization: true,
						},
					},
					{Name: "scratch", HostPath: "/mnt/scratch"},
				}
				spec.Container.MountPoints = []*AwsBatchJobDefinitionMountPoint{
					{SourceVolume: "shared-data", ContainerPath: "/data"},
					{SourceVolume: "scratch", ContainerPath: "/scratch", ReadOnly: false},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("with retry strategy, timeout, and scheduling priority", func() {
			ginkgo.It("should not return a validation error", func() {
				spec := minimalEc2Spec()
				spec.RetryStrategy = &AwsBatchJobDefinitionRetryStrategy{
					Attempts: 3,
					EvaluateOnExit: []*AwsBatchJobDefinitionEvaluateOnExit{
						{Action: "RETRY", OnStatusReason: "Host EC2*"},
						{Action: "EXIT", OnExitCode: "1*"},
					},
				}
				spec.Timeout = &AwsBatchJobDefinitionTimeout{AttemptDurationSeconds: 3600}
				spec.SchedulingPriority = 500
				spec.PropagateTags = true
				spec.Parameters = map[string]string{"input_path": "s3://bucket/default"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {

		ginkgo.Context("with no container", func() {
			ginkgo.It("should return a validation error", func() {
				err := protovalidate.Validate(&AwsBatchJobDefinitionSpec{Region: "us-west-2"})
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with no image", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.Image = ""
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with zero vcpus or memory", func() {
			ginkgo.It("should reject vcpus 0", func() {
				spec := minimalEc2Spec()
				spec.Container.Vcpus = 0
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject memory below the 4 MiB floor", func() {
				spec := minimalEc2Spec()
				spec.Container.MemoryMib = 2
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid platform capability", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.PlatformCapabilities = []string{"LAMBDA"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a Fargate job missing its execution role", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalFargateSpec()
				spec.Container.ExecutionRole = nil
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with EC2-only fields on a Fargate job", func() {
			ginkgo.It("should reject gpus", func() {
				spec := minimalFargateSpec()
				spec.Container.Gpus = 1
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject privileged", func() {
				spec := minimalFargateSpec()
				spec.Container.Privileged = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject ulimits", func() {
				spec := minimalFargateSpec()
				spec.Container.Ulimits = []*AwsBatchJobDefinitionUlimit{{Name: "nofile", SoftLimit: 1024, HardLimit: 4096}}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject linux_parameters", func() {
				spec := minimalFargateSpec()
				spec.Container.LinuxParameters = &AwsBatchJobDefinitionLinuxParameters{InitProcessEnabled: true}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with Fargate-only fields on an EC2 job", func() {
			ginkgo.It("should reject fargate_platform_version", func() {
				spec := minimalEc2Spec()
				spec.Container.FargatePlatformVersion = "LATEST"
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject assign_public_ip", func() {
				spec := minimalEc2Spec()
				spec.Container.AssignPublicIp = true
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject ephemeral_storage_gib", func() {
				spec := minimalEc2Spec()
				spec.Container.EphemeralStorageGib = 50
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject runtime_platform", func() {
				spec := minimalEc2Spec()
				spec.Container.RuntimePlatform = &AwsBatchJobDefinitionRuntimePlatform{CpuArchitecture: "ARM64"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a reserved environment variable name", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.Environment = map[string]string{"AWS_BATCH_CUSTOM": "value"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a mount point referencing an undeclared volume", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.MountPoints = []*AwsBatchJobDefinitionMountPoint{
					{SourceVolume: "missing", ContainerPath: "/data"},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a volume backed by both EFS and a host path", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.Volumes = []*AwsBatchJobDefinitionVolume{
					{
						Name:     "conflicted",
						HostPath: "/mnt/x",
						Efs: &AwsBatchJobDefinitionEfsVolume{
							FileSystemId: svRef("fs-0123456789abcdef0"),
						},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid log driver", func() {
			ginkgo.It("should return a validation error (awsfirelens is not a Batch driver)", func() {
				spec := minimalEc2Spec()
				spec.Container.LogConfiguration = &AwsBatchJobDefinitionLogConfiguration{
					LogDriver: "awsfirelens",
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid retry strategy", func() {
			ginkgo.It("should reject attempts above 10", func() {
				spec := minimalEc2Spec()
				spec.RetryStrategy = &AwsBatchJobDefinitionRetryStrategy{Attempts: 11}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject more than five evaluate_on_exit conditions", func() {
				spec := minimalEc2Spec()
				conditions := make([]*AwsBatchJobDefinitionEvaluateOnExit, 6)
				for i := range conditions {
					conditions[i] = &AwsBatchJobDefinitionEvaluateOnExit{Action: "RETRY", OnExitCode: "1"}
				}
				spec.RetryStrategy = &AwsBatchJobDefinitionRetryStrategy{Attempts: 3, EvaluateOnExit: conditions}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject an exit-code glob with a leading wildcard", func() {
				spec := minimalEc2Spec()
				spec.RetryStrategy = &AwsBatchJobDefinitionRetryStrategy{
					Attempts:       3,
					EvaluateOnExit: []*AwsBatchJobDefinitionEvaluateOnExit{{Action: "RETRY", OnExitCode: "*1"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
			ginkgo.It("should reject an unknown action", func() {
				spec := minimalEc2Spec()
				spec.RetryStrategy = &AwsBatchJobDefinitionRetryStrategy{
					Attempts:       3,
					EvaluateOnExit: []*AwsBatchJobDefinitionEvaluateOnExit{{Action: "PAUSE", OnExitCode: "1"}},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with a timeout below 60 seconds", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Timeout = &AwsBatchJobDefinitionTimeout{AttemptDurationSeconds: 30}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with scheduling_priority above 9999", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.SchedulingPriority = 10000
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with ephemeral storage out of range on a Fargate job", func() {
			ginkgo.It("should return a validation error when below 21", func() {
				spec := minimalFargateSpec()
				spec.Container.EphemeralStorageGib = 10
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid runtime platform", func() {
			ginkgo.It("should reject an unknown CPU architecture", func() {
				spec := minimalFargateSpec()
				spec.Container.RuntimePlatform = &AwsBatchJobDefinitionRuntimePlatform{CpuArchitecture: "RISCV"}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("with an invalid device permission", func() {
			ginkgo.It("should return a validation error", func() {
				spec := minimalEc2Spec()
				spec.Container.LinuxParameters = &AwsBatchJobDefinitionLinuxParameters{
					Devices: []*AwsBatchJobDefinitionDevice{
						{HostPath: "/dev/xvdf", Permissions: []string{"EXECUTE"}},
					},
				}
				err := protovalidate.Validate(spec)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
