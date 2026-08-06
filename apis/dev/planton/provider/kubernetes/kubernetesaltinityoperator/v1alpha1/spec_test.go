package kubernetesaltinityoperatorv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesAltinityOperatorSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesAltinityOperatorSpec Validation Suite")
}

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

var _ = ginkgo.Describe("KubernetesAltinityOperatorSpec validations", func() {
	var spec *KubernetesAltinityOperatorSpec

	ginkgo.BeforeEach(func() {
		spec = &KubernetesAltinityOperatorSpec{
			Namespace: literal("clickhouse-operator"),
		}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.It("accepts a minimal spec (all chart defaults)", func() {
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts watch namespaces and namespace-scoped rbac", func() {
			spec.WatchNamespaces = []string{"analytics", "data-.*"}
			spec.NamespaceScopedRbac = true
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts operator credentials with a password", func() {
			spec.OperatorCredentials = &KubernetesAltinityOperatorCredentials{
				Password: literal("a-real-password"),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts crd hook and metrics tuning", func() {
			spec.CrdHook = &KubernetesAltinityOperatorCrdHook{
				Enabled: boolPtr(true),
			}
			spec.Metrics = &KubernetesAltinityOperatorMetrics{
				Enabled: boolPtr(false),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Describe("When invalid input is passed", func() {
		ginkgo.It("rejects a missing namespace", func() {
			spec.Namespace = nil
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects operator credentials without a password", func() {
			spec.OperatorCredentials = &KubernetesAltinityOperatorCredentials{
				Username: stringPtr("clickhouse_operator"),
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})
})

func boolPtr(v bool) *bool { return &v }

func stringPtr(v string) *string { return &v }
