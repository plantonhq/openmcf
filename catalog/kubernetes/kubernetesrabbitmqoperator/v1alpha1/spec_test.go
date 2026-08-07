package kubernetesrabbitmqoperatorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	kubernetes "github.com/plantonhq/planton/catalog/kubernetes"
)

func TestKubernetesRabbitMqOperatorSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesRabbitMqOperatorSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesRabbitMqOperatorSpec validations", func() {
	var spec *KubernetesRabbitMqOperatorSpec

	ginkgo.BeforeEach(func() {
		spec = &KubernetesRabbitMqOperatorSpec{}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.It("accepts an empty spec (the release-manifest defaults)", func() {
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a fenced watch scope", func() {
			spec.WatchNamespaces = []string{"messaging", "team-a"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts air-gap image overrides", func() {
			spec.DefaultRabbitmqImage = "mirror.internal/rabbitmq:4.2.6-management"
			spec.DefaultUserUpdaterImage = "mirror.internal/default-user-credential-updater:1.0.14"
			spec.OperatorImage = &kubernetes.ContainerImage{
				Repo: "mirror.internal/rabbitmq/cluster-operator",
				Tag:  "2.22.3",
			}
			spec.ImagePullSecrets = []string{"mirror-pull"}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts sizing and placement", func() {
			spec.Resources = &kubernetes.ContainerResources{
				Requests: &kubernetes.CpuMemory{Cpu: "200m", Memory: "500Mi"},
				Limits:   &kubernetes.CpuMemory{Cpu: "500m", Memory: "1Gi"},
			}
			spec.NodeSelector = map[string]string{"workload": "system"}
			spec.Tolerations = []*kubernetes.WorkloadToleration{
				{Key: "system", Operator: "Exists", Effect: "NoSchedule"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("rejects watch namespaces that are not DNS-1123 names", func() {
			spec.WatchNamespaces = []string{"Team_A"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty watch namespace entry", func() {
			spec.WatchNamespaces = []string{""}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
