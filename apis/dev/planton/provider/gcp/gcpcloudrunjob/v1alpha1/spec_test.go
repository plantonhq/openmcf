package gcpcloudrunjobv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestGcpCloudRunJobSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GcpCloudRunJobSpec Validation Suite")
}

var _ = Describe("GcpCloudRunJobSpec validations", func() {

	strVal := func(v string) *foreignkeyv1.StringValueOrRef {
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
		}
	}

	makeValidSpec := func() *GcpCloudRunJobSpec {
		return &GcpCloudRunJobSpec{
			ProjectId: strVal("my-gcp-project"),
			Region:    "us-central1",
			Template: &GcpCloudRunJobTemplate{
				Containers: []*GcpCloudRunJobContainer{
					{Image: "us-docker.pkg.dev/my-project/repo/worker:v1.0.0"},
				},
			},
		}
	}

	Context("Required fields", func() {
		It("accepts a minimal valid spec", func() {
			Expect(protovalidate.Validate(makeValidSpec())).To(BeNil())
		})

		It("accepts a spec without project_id (ambient project)", func() {
			spec := makeValidSpec()
			spec.ProjectId = nil
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects spec with missing region", func() {
			spec := makeValidSpec()
			spec.Region = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects spec with missing template", func() {
			spec := makeValidSpec()
			spec.Template = nil
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects template with no containers", func() {
			spec := makeValidSpec()
			spec.Template.Containers = nil
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a container with no image", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Image = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Region validation", func() {
		It("accepts a standard region", func() {
			spec := makeValidSpec()
			spec.Region = "europe-west1"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a zone as region", func() {
			spec := makeValidSpec()
			spec.Region = "us-central1-a"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Job name validation", func() {
		It("accepts a valid job name", func() {
			spec := makeValidSpec()
			spec.JobName = "my-batch-job"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an invalid job name", func() {
			spec := makeValidSpec()
			spec.JobName = "My_Job"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Task count and parallelism", func() {
		It("accepts valid task_count and parallelism", func() {
			spec := makeValidSpec()
			spec.TaskCount = proto.Int32(10)
			spec.Parallelism = proto.Int32(5)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects parallelism greater than task_count", func() {
			spec := makeValidSpec()
			spec.TaskCount = proto.Int32(3)
			spec.Parallelism = proto.Int32(5)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects task_count below 1", func() {
			spec := makeValidSpec()
			spec.TaskCount = proto.Int32(0)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("GPU zonal redundancy", func() {
		It("rejects gpu_zonal_redundancy_disabled without node_selector", func() {
			spec := makeValidSpec()
			spec.GpuZonalRedundancyDisabled = true
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts gpu_zonal_redundancy_disabled with node_selector", func() {
			spec := makeValidSpec()
			spec.GpuZonalRedundancyDisabled = true
			spec.Template.NodeSelector = &GcpCloudRunJobNodeSelector{Accelerator: "nvidia-l4"}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})

	Context("Environment variables", func() {
		It("accepts a literal env var", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Env = []*GcpCloudRunJobEnvVar{{Name: "BATCH_SIZE", Value: "100"}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a secret env var", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Env = []*GcpCloudRunJobEnvVar{{
				Name:            "DB_PASSWORD",
				ValueFromSecret: &GcpCloudRunJobSecretEnvSource{Secret: "db-password", Version: "latest"},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects env var with both value and secret", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Env = []*GcpCloudRunJobEnvVar{{
				Name:            "X",
				Value:           "v",
				ValueFromSecret: &GcpCloudRunJobSecretEnvSource{Secret: "s"},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects env var with invalid name", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Env = []*GcpCloudRunJobEnvVar{{Name: "1BAD", Value: "v"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Container resources", func() {
		It("accepts valid cpu and memory", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Resources = &GcpCloudRunJobContainerResources{Cpu: "2", Memory: "4Gi"}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects invalid memory format", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].Resources = &GcpCloudRunJobContainerResources{Memory: "512"}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Volumes", func() {
		It("accepts a Cloud SQL volume", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "cloudsql",
				Source: &GcpCloudRunJobVolume_CloudSqlInstance{
					CloudSqlInstance: &GcpCloudRunJobVolumeCloudSql{
						Instances: []*foreignkeyv1.StringValueOrRef{strVal("p:r:i")},
					},
				},
			}}
			spec.Template.Containers[0].VolumeMounts = []*GcpCloudRunJobVolumeMount{{Name: "cloudsql", MountPath: "/cloudsql"}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects Cloud SQL volume without instances", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "cloudsql",
				Source: &GcpCloudRunJobVolume_CloudSqlInstance{
					CloudSqlInstance: &GcpCloudRunJobVolumeCloudSql{},
				},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts a secret volume", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "certs",
				Source: &GcpCloudRunJobVolume_Secret{
					Secret: &GcpCloudRunJobVolumeSecret{
						Secret: "tls-cert",
						Items:  []*GcpCloudRunJobVolumeSecretItem{{Path: "tls.crt", Version: "latest"}},
					},
				},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts an empty_dir volume", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "scratch",
				Source: &GcpCloudRunJobVolume_EmptyDir{
					EmptyDir: &GcpCloudRunJobVolumeEmptyDir{Medium: "MEMORY", SizeLimit: "256Mi"},
				},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a GCS volume", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "data",
				Source: &GcpCloudRunJobVolume_Gcs{
					Gcs: &GcpCloudRunJobVolumeGcs{Bucket: strVal("my-bucket"), ReadOnly: true},
				},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts an NFS volume", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{
				Name: "share",
				Source: &GcpCloudRunJobVolume_Nfs{
					Nfs: &GcpCloudRunJobVolumeNfs{Server: "10.0.0.5", Path: "/share1"},
				},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects volume without source", func() {
			spec := makeValidSpec()
			spec.Template.Volumes = []*GcpCloudRunJobVolume{{Name: "empty"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects mount path not absolute", func() {
			spec := makeValidSpec()
			spec.Template.Containers[0].VolumeMounts = []*GcpCloudRunJobVolumeMount{{Name: "v", MountPath: "data"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("VPC access", func() {
		It("accepts connector-only vpc access", func() {
			spec := makeValidSpec()
			spec.Template.VpcAccess = &GcpCloudRunJobVpcAccess{Connector: strVal("my-connector")}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts direct VPC network_interfaces", func() {
			spec := makeValidSpec()
			spec.Template.VpcAccess = &GcpCloudRunJobVpcAccess{
				NetworkInterfaces: []*GcpCloudRunJobNetworkInterface{{
					Subnetwork: strVal("my-subnet"),
				}},
				Egress: "PRIVATE_RANGES_ONLY",
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects connector and network_interfaces together", func() {
			spec := makeValidSpec()
			spec.Template.VpcAccess = &GcpCloudRunJobVpcAccess{
				Connector: strVal("my-connector"),
				NetworkInterfaces: []*GcpCloudRunJobNetworkInterface{{
					Subnetwork: strVal("my-subnet"),
				}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Binary authorization", func() {
		It("accepts use_default", func() {
			spec := makeValidSpec()
			spec.BinaryAuthorization = &GcpCloudRunJobBinaryAuthorization{UseDefault: true}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects use_default and policy together", func() {
			spec := makeValidSpec()
			spec.BinaryAuthorization = &GcpCloudRunJobBinaryAuthorization{
				UseDefault: true,
				Policy:     "projects/p/platforms/cloudRun/policies/default",
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Timeout and retries", func() {
		It("accepts valid timeout_seconds", func() {
			spec := makeValidSpec()
			spec.Template.TimeoutSeconds = proto.Int32(3600)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects timeout_seconds above 86400", func() {
			spec := makeValidSpec()
			spec.Template.TimeoutSeconds = proto.Int32(90000)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts max_retries of zero", func() {
			spec := makeValidSpec()
			spec.Template.MaxRetries = proto.Int32(0)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})

	Context("Launch stage", func() {
		It("accepts BETA launch stage", func() {
			spec := makeValidSpec()
			spec.LaunchStage = "BETA"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects invalid launch stage", func() {
			spec := makeValidSpec()
			spec.LaunchStage = "PREVIEW"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})
})
