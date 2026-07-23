package kubernetesistiov1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
	"google.golang.org/protobuf/proto"
)

func TestKubernetesIstioSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesIstioSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesIstioSpec validations", func() {
	var spec *KubernetesIstioSpec

	ginkgo.BeforeEach(func() {
		spec = &KubernetesIstioSpec{
			Namespace: &foreignkeyv1.StringValueOrRef{
				LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
					Value: "istio-system",
				},
			},
		}
	})

	ginkgo.Describe("minimal spec", func() {
		ginkgo.It("accepts a namespace-only spec (every other field defaulted or optional)", func() {
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects a spec without the required namespace", func() {
			spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("version", func() {
		ginkgo.It("accepts an omitted version (default applies at manifest load)", func() {
			spec.Version = nil
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit full-patch version", func() {
			spec.Version = proto.String("1.30.3")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts pinning an older version (sequential-upgrade safety)", func() {
			spec.Version = proto.String("1.28.2")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects a bare minor version (full patch required)", func() {
			spec.Version = proto.String("1.30")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a non-semver version string", func() {
			spec.Version = proto.String("latest")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("revision", func() {
		ginkgo.It("accepts a DNS-1123 revision label", func() {
			spec.Revision = proto.String("1-30-3")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase revision", func() {
			spec.Revision = proto.String("Canary")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a revision with a trailing hyphen", func() {
			spec.Revision = proto.String("canary-")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("dataplane_mode", func() {
		ginkgo.It("accepts sidecar", func() {
			spec.DataplaneMode = proto.String("sidecar")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts ambient", func() {
			spec.DataplaneMode = proto.String("ambient")
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown mode", func() {
			spec.DataplaneMode = proto.String("hybrid")
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("mode-scoped blocks", func() {
		ginkgo.It("rejects ztunnel settings in (default) sidecar mode — the release is not installed", func() {
			spec.Ztunnel = &KubernetesIstioZtunnel{LogLevel: proto.String("info")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects ztunnel settings in explicit sidecar mode", func() {
			spec.DataplaneMode = proto.String("sidecar")
			spec.Ztunnel = &KubernetesIstioZtunnel{LogLevel: proto.String("info")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts ztunnel settings in ambient mode", func() {
			spec.DataplaneMode = proto.String("ambient")
			spec.Ztunnel = &KubernetesIstioZtunnel{LogLevel: proto.String("info")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects sidecar_injector settings in ambient mode — ambient workloads are not injected", func() {
			spec.DataplaneMode = proto.String("ambient")
			spec.SidecarInjector = &KubernetesIstioSidecarInjector{EnableNamespacesByDefault: true}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts sidecar_injector settings in sidecar mode", func() {
			spec.SidecarInjector = &KubernetesIstioSidecarInjector{EnableNamespacesByDefault: true}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("istiod sizing", func() {
		ginkgo.It("accepts a fixed replica count without autoscale", func() {
			spec.Istiod = &KubernetesIstioIstiod{Replicas: proto.Int32(2)}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an autoscale block without replicas", func() {
			spec.Istiod = &KubernetesIstioIstiod{
				Autoscale: &KubernetesIstioIstiodAutoscale{
					MinReplicas: proto.Int32(2),
					MaxReplicas: proto.Int32(5),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects replicas combined with an enabled autoscale block — the HPA owns the count", func() {
			spec.Istiod = &KubernetesIstioIstiod{
				Replicas:  proto.Int32(2),
				Autoscale: &KubernetesIstioIstiodAutoscale{MaxReplicas: proto.Int32(5)},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("accepts replicas with an explicitly disabled autoscale block", func() {
			spec.Istiod = &KubernetesIstioIstiod{
				Replicas:  proto.Int32(2),
				Autoscale: &KubernetesIstioIstiodAutoscale{Enabled: proto.Bool(false)},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects zero replicas", func() {
			spec.Istiod = &KubernetesIstioIstiod{Replicas: proto.Int32(0)}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects autoscale max below min", func() {
			spec.Istiod = &KubernetesIstioIstiod{
				Autoscale: &KubernetesIstioIstiodAutoscale{
					MinReplicas: proto.Int32(3),
					MaxReplicas: proto.Int32(2),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a target CPU utilization above 100", func() {
			spec.Istiod = &KubernetesIstioIstiod{
				Autoscale: &KubernetesIstioIstiodAutoscale{
					TargetCpuUtilizationPercent: proto.Int32(150),
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("mesh_config", func() {
		ginkgo.It("accepts REGISTRY_ONLY outbound traffic policy", func() {
			spec.MeshConfig = &KubernetesIstioMeshConfig{
				OutboundTrafficPolicyMode: proto.String("REGISTRY_ONLY"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown outbound traffic policy mode", func() {
			spec.MeshConfig = &KubernetesIstioMeshConfig{
				OutboundTrafficPolicyMode: proto.String("BLOCK_ALL"),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("proxy defaults", func() {
		ginkgo.It("accepts a known proxy log level", func() {
			spec.Proxy = &KubernetesIstioProxy{LogLevel: proto.String("debug")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown proxy log level", func() {
			spec.Proxy = &KubernetesIstioProxy{LogLevel: proto.String("verbose")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown auto_inject policy", func() {
			spec.Proxy = &KubernetesIstioProxy{AutoInject: proto.String("always")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("ztunnel tuning", func() {
		ginkgo.It("rejects an unknown ztunnel log level", func() {
			spec.DataplaneMode = proto.String("ambient")
			spec.Ztunnel = &KubernetesIstioZtunnel{LogLevel: proto.String("critical")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("gateway_defaults", func() {
		ginkgo.It("accepts NodePort", func() {
			spec.GatewayDefaults = &KubernetesIstioGatewayDefaults{ServiceType: proto.String("NodePort")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown service type", func() {
			spec.GatewayDefaults = &KubernetesIstioGatewayDefaults{ServiceType: proto.String("ExternalName")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("images", func() {
		ginkgo.It("accepts the distroless variant", func() {
			spec.Images = &KubernetesIstioImages{Variant: proto.String("distroless")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("rejects an unknown image variant", func() {
			spec.Images = &KubernetesIstioImages{Variant: proto.String("alpine")}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Describe("full ambient surface", func() {
		ginkgo.It("accepts the complete ambient configuration", func() {
			spec.Version = proto.String("1.30.3")
			spec.Revision = proto.String("prod")
			spec.DataplaneMode = proto.String("ambient")
			spec.Istiod = &KubernetesIstioIstiod{
				Autoscale: &KubernetesIstioIstiodAutoscale{
					Enabled:                     proto.Bool(true),
					MinReplicas:                 proto.Int32(2),
					MaxReplicas:                 proto.Int32(5),
					TargetCpuUtilizationPercent: proto.Int32(75),
				},
				LogLevel:            proto.String("default:info"),
				PodDisruptionBudget: proto.Bool(true),
				PriorityClassName:   "system-cluster-critical",
			}
			spec.MeshConfig = &KubernetesIstioMeshConfig{
				TrustDomain:               proto.String("mesh.example.internal"),
				OutboundTrafficPolicyMode: proto.String("REGISTRY_ONLY"),
				AccessLogFile:             "/dev/stdout",
			}
			spec.Cni = &KubernetesIstioCni{
				ExcludeNamespaces: []string{"kube-system"},
				Chained:           proto.Bool(true),
			}
			spec.Ztunnel = &KubernetesIstioZtunnel{LogLevel: proto.String("info")}
			spec.GatewayDefaults = &KubernetesIstioGatewayDefaults{ServiceType: proto.String("NodePort")}
			spec.Images = &KubernetesIstioImages{Variant: proto.String("distroless")}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})
})
