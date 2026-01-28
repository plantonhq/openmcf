package kubernetescertmanagerv1

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

func TestKubernetesCertManager(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesCertManager Suite")
}

var _ = ginkgo.Describe("KubernetesCertManager Custom Validation Tests", func() {
	var input *KubernetesCertManager

	ginkgo.BeforeEach(func() {
		input = &KubernetesCertManager{
			ApiVersion: "kubernetes.openmcf.org/v1",
			Kind:       "KubernetesCertManager",
			Metadata: &shared.CloudResourceMetadata{
				Name: "test-kubernetes-cert-manager",
			},
			Spec: &KubernetesCertManagerSpec{
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
				Acme: &AcmeConfig{
					Email: "admin@example.com",
				},
				DnsProviders: []*DnsProviderConfig{
					{
						Name:     "cloudflare-test",
						DnsZones: []string{"example.com"},
						Provider: &DnsProviderConfig_Cloudflare{
							Cloudflare: &CloudflareProvider{
								ApiToken: "test-token",
							},
						},
					},
				},
			},
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("kubernetes_cert_manager", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(input)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
