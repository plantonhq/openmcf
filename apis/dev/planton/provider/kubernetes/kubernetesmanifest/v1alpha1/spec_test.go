package kubernetesmanifestv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesManifestSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesManifestSpec Validation Suite")
}

func literal(v string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: v},
	}
}

const singleDocYaml = `apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  key: value
`

const multiDocYaml = singleDocYaml + `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 1
`

var _ = ginkgo.Describe("KubernetesManifestSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal single-document manifest", func() {
			spec := &KubernetesManifestSpec{
				Namespace:    literal("apps"),
				ManifestYaml: singleDocYaml,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a multi-document manifest with namespace creation", func() {
			spec := &KubernetesManifestSpec{
				Namespace:       literal("apps"),
				CreateNamespace: true,
				ManifestYaml:    multiDocYaml,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace expressed as a foreign-key reference", func() {
			spec := &KubernetesManifestSpec{
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "apps-namespace"},
					},
				},
				ManifestYaml: singleDocYaml,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts skip_await", func() {
			spec := &KubernetesManifestSpec{
				Namespace:    literal("apps"),
				ManifestYaml: singleDocYaml,
				SkipAwait:    true,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing namespace", func() {
			spec := &KubernetesManifestSpec{
				ManifestYaml: singleDocYaml,
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a missing manifest body", func() {
			spec := &KubernetesManifestSpec{
				Namespace: literal("apps"),
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a whitespace-only manifest body", func() {
			spec := &KubernetesManifestSpec{
				Namespace:    literal("apps"),
				ManifestYaml: "   \n\t\n",
			}
			gomega.Expect(protovalidate.Validate(spec)).ToNot(gomega.BeNil())
		})
	})
})
