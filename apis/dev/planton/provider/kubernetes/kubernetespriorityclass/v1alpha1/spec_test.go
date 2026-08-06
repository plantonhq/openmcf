package kubernetespriorityclassv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestKubernetesPriorityClassSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesPriorityClassSpec Validation Suite")
}

func preemption(p KubernetesPriorityClassSpec_KubernetesPriorityClassPreemptionPolicy) *KubernetesPriorityClassSpec_KubernetesPriorityClassPreemptionPolicy {
	return &p
}

var _ = ginkgo.Describe("KubernetesPriorityClassSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal spec", func() {
			spec := &KubernetesPriorityClassSpec{Name: "critical", Value: 1000000}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts the maximum user-definable value", func() {
			spec := &KubernetesPriorityClassSpec{Name: "top", Value: 1000000000}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a negative value for always-preemptable tiers", func() {
			spec := &KubernetesPriorityClassSpec{Name: "batch", Value: -100}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a non-preempting global default", func() {
			spec := &KubernetesPriorityClassSpec{
				Name:             "standard",
				Value:            1000,
				GlobalDefault:    true,
				Description:      "Default priority for services",
				PreemptionPolicy: preemption(KubernetesPriorityClassSpec_never),
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing name", func() {
			spec := &KubernetesPriorityClassSpec{Value: 100}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a value above the user-definable ceiling", func() {
			spec := &KubernetesPriorityClassSpec{Name: "too-high", Value: 1000000001}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects the reserved system- name prefix", func() {
			spec := &KubernetesPriorityClassSpec{Name: "system-critical", Value: 100}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an uppercase name", func() {
			spec := &KubernetesPriorityClassSpec{Name: "Critical", Value: 100}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an undefined preemption policy enum value", func() {
			bad := KubernetesPriorityClassSpec_KubernetesPriorityClassPreemptionPolicy(99)
			spec := &KubernetesPriorityClassSpec{Name: "critical", Value: 100, PreemptionPolicy: &bad}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
