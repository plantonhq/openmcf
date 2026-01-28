package kubernetesopenbaov1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/plantonhq/openmcf/apis/org/openmcf/provider/kubernetes"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared"
	"github.com/plantonhq/openmcf/apis/org/openmcf/shared/cloudresourcekind"
	foreignkeyv1 "github.com/plantonhq/openmcf/apis/org/openmcf/shared/foreignkey/v1"
)

func TestKubernetesOpenBao(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesOpenBao Suite")
}

var _ = ginkgo.Describe("KubernetesOpenBao Custom Validation Tests", func() {
	var input *KubernetesOpenBao

	ginkgo.BeforeEach(func() {
		input = &KubernetesOpenBao{
			ApiVersion: "kubernetes.openmcf.org/v1",
			Kind:       "KubernetesOpenBao",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-openbao",
			},
			Spec: &KubernetesOpenBaoSpec{
				TargetCluster: &kubernetes.KubernetesClusterSelector{
					ClusterKind: cloudresourcekind.CloudResourceKind_GcpGkeCluster,
					ClusterName: "test-cluster",
				},
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{
						Value: "test-namespace",
					},
				},
				CreateNamespace: true,
				ServerContainer: &KubernetesOpenBaoServerContainer{
					Replicas:        1,
					DataStorageSize: "10Gi",
					Resources: &kubernetes.ContainerResources{
						Limits: &kubernetes.CpuMemory{
							Cpu:    "500m",
							Memory: "256Mi",
						},
						Requests: &kubernetes.CpuMemory{
							Cpu:    "100m",
							Memory: "128Mi",
						},
					},
				},
			},
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("standalone openbao deployment", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("standalone openbao with ingress enabled", func() {
			ginkgo.It("should not return a validation error", func() {
				input.Spec.Ingress = &KubernetesOpenBaoIngress{
					Enabled:  true,
					Hostname: "openbao.example.com",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("ha mode openbao deployment", func() {
			ginkgo.It("should not return a validation error", func() {
				haReplicas := int32(3)
				input.Spec.HighAvailability = &KubernetesOpenBaoHighAvailability{
					Enabled:  true,
					Replicas: &haReplicas,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("openbao with injector enabled", func() {
			ginkgo.It("should not return a validation error", func() {
				injectorReplicas := int32(2)
				input.Spec.Injector = &KubernetesOpenBaoInjector{
					Enabled:  true,
					Replicas: &injectorReplicas,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.Context("missing namespace", func() {
			ginkgo.It("should return a validation error", func() {
				input.Spec.Namespace = nil
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("invalid data storage size format", func() {
			ginkgo.It("should return a validation error", func() {
				input.Spec.ServerContainer.DataStorageSize = "invalid"
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("server replicas below minimum", func() {
			ginkgo.It("should return a validation error", func() {
				input.Spec.ServerContainer.Replicas = 0
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("server replicas above maximum", func() {
			ginkgo.It("should return a validation error", func() {
				input.Spec.ServerContainer.Replicas = 11
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("ha replicas below minimum", func() {
			ginkgo.It("should return a validation error", func() {
				haReplicas := int32(1)
				input.Spec.HighAvailability = &KubernetesOpenBaoHighAvailability{
					Enabled:  true,
					Replicas: &haReplicas,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("ingress enabled without hostname", func() {
			ginkgo.It("should return a validation error", func() {
				input.Spec.Ingress = &KubernetesOpenBaoIngress{
					Enabled:  true,
					Hostname: "",
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})

		ginkgo.Context("injector replicas above maximum", func() {
			ginkgo.It("should return a validation error", func() {
				injectorReplicas := int32(6)
				input.Spec.Injector = &KubernetesOpenBaoInjector{
					Enabled:  true,
					Replicas: &injectorReplicas,
				}
				err := protovalidate.Validate(input)
				gomega.Expect(err).ToNot(gomega.BeNil())
			})
		})
	})
})
