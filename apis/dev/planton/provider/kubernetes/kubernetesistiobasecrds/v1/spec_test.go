package kubernetesistiobasecrdsv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestKubernetesIstioBaseCrdsSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesIstioBaseCrdsSpec Validation Suite")
}

// The spec is intentionally minimal: the CRD version is pinned to the typed SDK (no
// user version field), so there are no value-validation rules to exercise here. The
// case below asserts the empty spec validates cleanly.
var _ = ginkgo.Describe("KubernetesIstioBaseCrdsSpec validations", func() {
	var spec *KubernetesIstioBaseCrdsSpec

	ginkgo.BeforeEach(func() {
		spec = &KubernetesIstioBaseCrdsSpec{}
	})

	ginkgo.Describe("When valid input is passed", func() {
		ginkgo.Context("with an empty spec", func() {
			ginkgo.It("should not return a validation error", func() {
				err := protovalidate.Validate(spec)
				gomega.Expect(err).To(gomega.BeNil())
			})
		})
	})
})
