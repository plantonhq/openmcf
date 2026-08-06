package gcpcloudrunv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestGcpCloudRunSpec(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GcpCloudRunSpec Validation Suite")
}

var _ = Describe("GcpCloudRunSpec validations", func() {

	strVal := func(v string) *foreignkeyv1.StringValueOrRef {
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
		}
	}

	strRef := func(kind, name string) *foreignkeyv1.StringValueOrRef {
		return &foreignkeyv1.StringValueOrRef{
			LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
				ValueFrom: &foreignkeyv1.ValueFromRef{Name: name},
			},
		}
	}

	makeValidSpec := func() *GcpCloudRunSpec {
		return &GcpCloudRunSpec{
			ProjectId: strVal("my-gcp-project"),
			Region:    "us-central1",
			Containers: []*GcpCloudRunContainer{
				{Image: "us-docker.pkg.dev/my-project/repo/app:v1.0.0"},
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

		It("rejects spec with no containers", func() {
			spec := makeValidSpec()
			spec.Containers = nil
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a container with no image", func() {
			spec := makeValidSpec()
			spec.Containers[0].Image = ""
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Region validation", func() {
		It("accepts a standard region", func() {
			spec := makeValidSpec()
			spec.Region = "europe-west1"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a multi-digit region", func() {
			spec := makeValidSpec()
			spec.Region = "europe-west12"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a zone as region", func() {
			spec := makeValidSpec()
			spec.Region = "us-central1-a"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects arbitrary text", func() {
			spec := makeValidSpec()
			spec.Region = "not a region"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Service name validation", func() {
		It("accepts a valid service name", func() {
			spec := makeValidSpec()
			spec.ServiceName = "my-api-service"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a name starting with a digit", func() {
			spec := makeValidSpec()
			spec.ServiceName = "1api"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a name with uppercase letters", func() {
			spec := makeValidSpec()
			spec.ServiceName = "MyApi"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a name over 63 characters", func() {
			spec := makeValidSpec()
			spec.ServiceName = "a" + string(make([]byte, 0)) + "bcdefghij-bcdefghij-bcdefghij-bcdefghij-bcdefghij-bcdefghij-bcde"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Auth posture", func() {
		It("accepts allow_unauthenticated alone", func() {
			spec := makeValidSpec()
			spec.AllowUnauthenticated = true
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts invoker_iam_disabled alone", func() {
			spec := makeValidSpec()
			spec.InvokerIamDisabled = true
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects allow_unauthenticated together with invoker_iam_disabled", func() {
			spec := makeValidSpec()
			spec.AllowUnauthenticated = true
			spec.InvokerIamDisabled = true
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Containers", func() {
		It("accepts a sidecar arrangement with names and depends_on", func() {
			spec := makeValidSpec()
			spec.Containers = []*GcpCloudRunContainer{
				{
					Name:      "app",
					Image:     "us-docker.pkg.dev/p/r/app:1",
					Ports:     &GcpCloudRunContainerPort{ContainerPort: proto.Int32(8080)},
					DependsOn: []string{"proxy"},
				},
				{
					Name:  "proxy",
					Image: "us-docker.pkg.dev/p/r/proxy:1",
				},
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an invalid container name", func() {
			spec := makeValidSpec()
			spec.Containers[0].Name = "App!"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects an out-of-range container port", func() {
			spec := makeValidSpec()
			spec.Containers[0].Ports = &GcpCloudRunContainerPort{ContainerPort: proto.Int32(70000)}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts the h2c port protocol", func() {
			spec := makeValidSpec()
			spec.Containers[0].Ports = &GcpCloudRunContainerPort{Name: "h2c"}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an unknown port protocol", func() {
			spec := makeValidSpec()
			spec.Containers[0].Ports = &GcpCloudRunContainerPort{Name: "tcp"}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Container resources", func() {
		It("accepts whole, fractional, and millicore CPU", func() {
			for _, cpu := range []string{"1", "4", "0.5", "500m"} {
				spec := makeValidSpec()
				spec.Containers[0].Resources = &GcpCloudRunContainerResources{Cpu: cpu, Memory: "512Mi"}
				Expect(protovalidate.Validate(spec)).To(BeNil(), "cpu=%s", cpu)
			}
		})

		It("rejects malformed CPU", func() {
			spec := makeValidSpec()
			spec.Containers[0].Resources = &GcpCloudRunContainerResources{Cpu: "one"}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts unit-suffixed memory", func() {
			for _, mem := range []string{"512Mi", "2Gi", "128M"} {
				spec := makeValidSpec()
				spec.Containers[0].Resources = &GcpCloudRunContainerResources{Memory: mem}
				Expect(protovalidate.Validate(spec)).To(BeNil(), "memory=%s", mem)
			}
		})

		It("rejects memory without a unit", func() {
			spec := makeValidSpec()
			spec.Containers[0].Resources = &GcpCloudRunContainerResources{Memory: "512"}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts the CPU allocation levers", func() {
			spec := makeValidSpec()
			spec.Containers[0].Resources = &GcpCloudRunContainerResources{
				CpuIdle:         proto.Bool(false),
				StartupCpuBoost: true,
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})

	Context("Environment variables", func() {
		It("accepts a literal value", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{Name: "LOG_LEVEL", Value: "info"}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a Secret Manager reference", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{
				Name:            "DB_PASSWORD",
				ValueFromSecret: &GcpCloudRunSecretEnvSource{Secret: "db-password", Version: "latest"},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects both a value and a secret reference", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{
				Name:            "DB_PASSWORD",
				Value:           "plain",
				ValueFromSecret: &GcpCloudRunSecretEnvSource{Secret: "db-password"},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a variable without a name", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{Value: "v"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a name starting with a digit", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{Name: "1BAD", Value: "v"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a secret reference without the secret name", func() {
			spec := makeValidSpec()
			spec.Containers[0].Env = []*GcpCloudRunEnvVar{{
				Name:            "DB_PASSWORD",
				ValueFromSecret: &GcpCloudRunSecretEnvSource{},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Probes", func() {
		It("accepts an HTTP startup probe", func() {
			spec := makeValidSpec()
			spec.Containers[0].StartupProbe = &GcpCloudRunProbe{
				InitialDelaySeconds: proto.Int32(0),
				PeriodSeconds:       proto.Int32(10),
				TimeoutSeconds:      proto.Int32(1),
				FailureThreshold:    proto.Int32(3),
				Handler:             &GcpCloudRunProbe_HttpGet{HttpGet: &GcpCloudRunHttpGetAction{Path: "/healthz"}},
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a TCP startup probe", func() {
			spec := makeValidSpec()
			spec.Containers[0].StartupProbe = &GcpCloudRunProbe{
				Handler: &GcpCloudRunProbe_TcpSocket{TcpSocket: &GcpCloudRunTcpSocketAction{Port: proto.Int32(8080)}},
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a gRPC liveness probe", func() {
			spec := makeValidSpec()
			spec.Containers[0].LivenessProbe = &GcpCloudRunProbe{
				Handler: &GcpCloudRunProbe_Grpc{Grpc: &GcpCloudRunGrpcAction{Service: "my.Service"}},
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a TCP liveness probe", func() {
			spec := makeValidSpec()
			spec.Containers[0].LivenessProbe = &GcpCloudRunProbe{
				Handler: &GcpCloudRunProbe_TcpSocket{TcpSocket: &GcpCloudRunTcpSocketAction{}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a probe without a handler", func() {
			spec := makeValidSpec()
			spec.Containers[0].StartupProbe = &GcpCloudRunProbe{PeriodSeconds: proto.Int32(10)}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects timeout exceeding period", func() {
			spec := makeValidSpec()
			spec.Containers[0].StartupProbe = &GcpCloudRunProbe{
				TimeoutSeconds: proto.Int32(20),
				PeriodSeconds:  proto.Int32(10),
				Handler:        &GcpCloudRunProbe_HttpGet{HttpGet: &GcpCloudRunHttpGetAction{}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a probe path not starting with /", func() {
			spec := makeValidSpec()
			spec.Containers[0].StartupProbe = &GcpCloudRunProbe{
				Handler: &GcpCloudRunProbe_HttpGet{HttpGet: &GcpCloudRunHttpGetAction{Path: "healthz"}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Volumes", func() {
		It("accepts a Cloud SQL volume with a reference", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name: "cloudsql",
				Source: &GcpCloudRunVolume_CloudSqlInstance{CloudSqlInstance: &GcpCloudRunVolumeCloudSql{
					Instances: []*foreignkeyv1.StringValueOrRef{strRef("GcpCloudSql", "my-db")},
				}},
			}}
			spec.Containers[0].VolumeMounts = []*GcpCloudRunVolumeMount{{Name: "cloudsql", MountPath: "/cloudsql"}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a Cloud SQL volume with no instances", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "cloudsql",
				Source: &GcpCloudRunVolume_CloudSqlInstance{CloudSqlInstance: &GcpCloudRunVolumeCloudSql{}},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts a secret volume with items", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name: "certs",
				Source: &GcpCloudRunVolume_Secret{Secret: &GcpCloudRunVolumeSecret{
					Secret:      "tls-cert",
					DefaultMode: proto.Int32(292),
					Items:       []*GcpCloudRunVolumeSecretItem{{Path: "tls.crt", Version: "latest"}},
				}},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a secret volume without the secret name", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "certs",
				Source: &GcpCloudRunVolume_Secret{Secret: &GcpCloudRunVolumeSecret{}},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts an in-memory empty_dir volume", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "scratch",
				Source: &GcpCloudRunVolume_EmptyDir{EmptyDir: &GcpCloudRunVolumeEmptyDir{Medium: "MEMORY", SizeLimit: "256Mi"}},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an unknown empty_dir medium", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "scratch",
				Source: &GcpCloudRunVolume_EmptyDir{EmptyDir: &GcpCloudRunVolumeEmptyDir{Medium: "SSD"}},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts a GCS volume with a bucket reference", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name: "assets",
				Source: &GcpCloudRunVolume_Gcs{Gcs: &GcpCloudRunVolumeGcs{
					Bucket:   strRef("GcpGcsBucket", "my-bucket"),
					ReadOnly: true,
				}},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts an NFS volume", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "shared",
				Source: &GcpCloudRunVolume_Nfs{Nfs: &GcpCloudRunVolumeNfs{Server: "10.0.0.5", Path: "/share1"}},
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an NFS volume with a relative path", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{
				Name:   "shared",
				Source: &GcpCloudRunVolume_Nfs{Nfs: &GcpCloudRunVolumeNfs{Server: "10.0.0.5", Path: "share1"}},
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a volume with no source", func() {
			spec := makeValidSpec()
			spec.Volumes = []*GcpCloudRunVolume{{Name: "empty"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a volume mount with a relative path", func() {
			spec := makeValidSpec()
			spec.Containers[0].VolumeMounts = []*GcpCloudRunVolumeMount{{Name: "v", MountPath: "data"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Scaling", func() {
		It("accepts revision scaling bounds", func() {
			spec := makeValidSpec()
			spec.Scaling = &GcpCloudRunRevisionScaling{
				MinInstanceCount: proto.Int32(1),
				MaxInstanceCount: proto.Int32(50),
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects min above max", func() {
			spec := makeValidSpec()
			spec.Scaling = &GcpCloudRunRevisionScaling{
				MinInstanceCount: proto.Int32(10),
				MaxInstanceCount: proto.Int32(5),
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts MANUAL service scaling with a count", func() {
			spec := makeValidSpec()
			spec.ServiceScaling = &GcpCloudRunServiceScaling{
				ScalingMode:         "MANUAL",
				ManualInstanceCount: proto.Int32(3),
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects manual_instance_count without MANUAL mode", func() {
			spec := makeValidSpec()
			spec.ServiceScaling = &GcpCloudRunServiceScaling{
				ManualInstanceCount: proto.Int32(3),
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects out-of-range concurrency", func() {
			spec := makeValidSpec()
			spec.MaxInstanceRequestConcurrency = proto.Int32(1001)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects out-of-range timeout", func() {
			spec := makeValidSpec()
			spec.TimeoutSeconds = proto.Int32(3601)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Traffic", func() {
		It("accepts a latest-revision target", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{{
				Type:    "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST",
				Percent: proto.Int32(100),
			}}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a canary split with a tagged revision", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{
				{Type: "TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION", Revision: "my-api-v41", Percent: proto.Int32(90)},
				{Type: "TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION", Revision: "my-api-v42", Percent: proto.Int32(10), Tag: "canary"},
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a REVISION target without a revision name", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{{
				Type:    "TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION",
				Percent: proto.Int32(100),
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a LATEST target naming a revision", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{{
				Type:     "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST",
				Revision: "my-api-v42",
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects an unknown traffic type", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{{Type: "SPLIT"}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a percent above 100", func() {
			spec := makeValidSpec()
			spec.Traffic = []*GcpCloudRunTrafficTarget{{
				Type:    "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST",
				Percent: proto.Int32(150),
			}}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("VPC access", func() {
		It("accepts direct VPC egress with references", func() {
			spec := makeValidSpec()
			spec.VpcAccess = &GcpCloudRunVpcAccess{
				NetworkInterfaces: []*GcpCloudRunNetworkInterface{{
					Network:    strRef("GcpVpcNetwork", "my-vpc"),
					Subnetwork: strRef("GcpSubnetwork", "my-subnet"),
					Tags:       []string{"cloud-run-egress"},
				}},
				Egress: "PRIVATE_RANGES_ONLY",
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a connector", func() {
			spec := makeValidSpec()
			spec.VpcAccess = &GcpCloudRunVpcAccess{
				Connector: strVal("projects/p/locations/us-central1/connectors/c"),
				Egress:    "ALL_TRAFFIC",
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects a connector combined with network_interfaces", func() {
			spec := makeValidSpec()
			spec.VpcAccess = &GcpCloudRunVpcAccess{
				Connector: strVal("projects/p/locations/us-central1/connectors/c"),
				NetworkInterfaces: []*GcpCloudRunNetworkInterface{{
					Subnetwork: strRef("GcpSubnetwork", "my-subnet"),
				}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects an unknown egress setting", func() {
			spec := makeValidSpec()
			spec.VpcAccess = &GcpCloudRunVpcAccess{Egress: "SOME_TRAFFIC"}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects an invalid network tag", func() {
			spec := makeValidSpec()
			spec.VpcAccess = &GcpCloudRunVpcAccess{
				NetworkInterfaces: []*GcpCloudRunNetworkInterface{{
					Subnetwork: strRef("GcpSubnetwork", "my-subnet"),
					Tags:       []string{"Bad_Tag"},
				}},
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("GPU", func() {
		It("accepts an accelerator with zonal redundancy opt-out", func() {
			spec := makeValidSpec()
			spec.NodeSelector = &GcpCloudRunNodeSelector{Accelerator: "nvidia-l4"}
			spec.GpuZonalRedundancyDisabled = true
			spec.Containers[0].Resources = &GcpCloudRunContainerResources{Cpu: "4", Memory: "16Gi"}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects zonal redundancy opt-out without an accelerator", func() {
			spec := makeValidSpec()
			spec.GpuZonalRedundancyDisabled = true
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("rejects a node selector without an accelerator", func() {
			spec := makeValidSpec()
			spec.NodeSelector = &GcpCloudRunNodeSelector{}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Binary authorization", func() {
		It("accepts the project default policy", func() {
			spec := makeValidSpec()
			spec.BinaryAuthorization = &GcpCloudRunBinaryAuthorization{UseDefault: true}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a named policy", func() {
			spec := makeValidSpec()
			spec.BinaryAuthorization = &GcpCloudRunBinaryAuthorization{
				Policy: "projects/p/platforms/cloudRun/policies/default",
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects default combined with a named policy", func() {
			spec := makeValidSpec()
			spec.BinaryAuthorization = &GcpCloudRunBinaryAuthorization{
				UseDefault: true,
				Policy:     "projects/p/platforms/cloudRun/policies/default",
			}
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Service-level fields", func() {
		It("accepts a launch stage", func() {
			spec := makeValidSpec()
			spec.LaunchStage = "BETA"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an unknown launch stage", func() {
			spec := makeValidSpec()
			spec.LaunchStage = "PREVIEW"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts custom audiences", func() {
			spec := makeValidSpec()
			spec.CustomAudiences = []string{"https://api.example.com"}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("accepts a pinned revision name", func() {
			spec := makeValidSpec()
			spec.Revision = "my-api-v42"
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an invalid revision name", func() {
			spec := makeValidSpec()
			spec.Revision = "My_Api"
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})

		It("accepts an explicit deletion_protection opt-out", func() {
			spec := makeValidSpec()
			spec.DeletionProtection = proto.Bool(false)
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})

		It("rejects an over-long description", func() {
			spec := makeValidSpec()
			long := make([]byte, 513)
			for i := range long {
				long[i] = 'a'
			}
			spec.Description = string(long)
			Expect(protovalidate.Validate(spec)).NotTo(BeNil())
		})
	})

	Context("Full production spec", func() {
		It("accepts a fully populated spec", func() {
			spec := &GcpCloudRunSpec{
				ProjectId:   strVal("my-gcp-project"),
				Region:      "us-central1",
				ServiceName: "my-api",
				Description: "Primary API service",
				Labels:      map[string]string{"team": "platform"},
				Containers: []*GcpCloudRunContainer{
					{
						Name:  "app",
						Image: "us-docker.pkg.dev/p/r/app:1.2.3",
						Env: []*GcpCloudRunEnvVar{
							{Name: "LOG_LEVEL", Value: "info"},
							{Name: "DB_PASSWORD", ValueFromSecret: &GcpCloudRunSecretEnvSource{Secret: "db-password", Version: "3"}},
						},
						Ports: &GcpCloudRunContainerPort{ContainerPort: proto.Int32(8080)},
						Resources: &GcpCloudRunContainerResources{
							Cpu:             "2",
							Memory:          "1Gi",
							CpuIdle:         proto.Bool(true),
							StartupCpuBoost: true,
						},
						VolumeMounts: []*GcpCloudRunVolumeMount{{Name: "cloudsql", MountPath: "/cloudsql"}},
						StartupProbe: &GcpCloudRunProbe{
							PeriodSeconds: proto.Int32(5),
							Handler:       &GcpCloudRunProbe_HttpGet{HttpGet: &GcpCloudRunHttpGetAction{Path: "/healthz"}},
						},
						LivenessProbe: &GcpCloudRunProbe{
							Handler: &GcpCloudRunProbe_HttpGet{HttpGet: &GcpCloudRunHttpGetAction{Path: "/livez"}},
						},
					},
				},
				Volumes: []*GcpCloudRunVolume{{
					Name: "cloudsql",
					Source: &GcpCloudRunVolume_CloudSqlInstance{CloudSqlInstance: &GcpCloudRunVolumeCloudSql{
						Instances: []*foreignkeyv1.StringValueOrRef{strRef("GcpCloudSql", "my-db")},
					}},
				}},
				ServiceAccount:                strRef("GcpServiceAccount", "api-runtime"),
				Scaling:                       &GcpCloudRunRevisionScaling{MinInstanceCount: proto.Int32(1), MaxInstanceCount: proto.Int32(100)},
				MaxInstanceRequestConcurrency: proto.Int32(80),
				TimeoutSeconds:                proto.Int32(300),
				ExecutionEnvironment:          GcpCloudRunExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN2,
				SessionAffinity:               true,
				VpcAccess: &GcpCloudRunVpcAccess{
					NetworkInterfaces: []*GcpCloudRunNetworkInterface{{
						Network:    strRef("GcpVpcNetwork", "my-vpc"),
						Subnetwork: strRef("GcpSubnetwork", "my-subnet"),
					}},
					Egress: "PRIVATE_RANGES_ONLY",
				},
				Ingress:              GcpCloudRunIngress_INGRESS_TRAFFIC_ALL,
				AllowUnauthenticated: true,
				Traffic: []*GcpCloudRunTrafficTarget{{
					Type:    "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST",
					Percent: proto.Int32(100),
				}},
				DeletionProtection: proto.Bool(false),
			}
			Expect(protovalidate.Validate(spec)).To(BeNil())
		})
	})
})
