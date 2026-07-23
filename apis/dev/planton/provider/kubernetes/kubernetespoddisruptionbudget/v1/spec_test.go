package kubernetespoddisruptionbudgetv1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesPodDisruptionBudgetSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesPodDisruptionBudgetSpec Validation Suite")
}

func appSelector() *KubernetesPodDisruptionBudgetLabelSelector {
	return &KubernetesPodDisruptionBudgetLabelSelector{
		MatchLabels: map[string]string{"app": "checkout"},
	}
}

var _ = ginkgo.Describe("KubernetesPodDisruptionBudgetSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a min_available absolute budget", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:         "checkout-pdb",
				Selector:     appSelector(),
				MinAvailable: "2",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a max_unavailable percentage budget", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:           "checkout-pdb",
				Selector:       appSelector(),
				MaxUnavailable: "25%",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a resource reference", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:         "checkout-pdb",
				Selector:     appSelector(),
				MinAvailable: "50%",
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts an explicit protect-everything selector", func() {
			policy := KubernetesPodDisruptionBudgetSpec_always_allow
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:                       "namespace-wide",
				Selector:                   &KubernetesPodDisruptionBudgetLabelSelector{},
				MaxUnavailable:             "1",
				UnhealthyPodEvictionPolicy: &policy,
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts set-based selector expressions", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name: "tier-pdb",
				Selector: &KubernetesPodDisruptionBudgetLabelSelector{
					MatchExpressions: []*KubernetesPodDisruptionBudgetLabelSelectorRequirement{{
						Key:      "tier",
						Operator: "In",
						Values:   []string{"web", "api"},
					}},
				},
				MinAvailable: "1",
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing selector", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{Name: "no-selector", MinAvailable: "1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects both bounds set", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:           "both",
				Selector:       appSelector(),
				MinAvailable:   "1",
				MaxUnavailable: "1",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects neither bound set", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:     "neither",
				Selector: appSelector(),
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a percentage above 100", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:         "over",
				Selector:     appSelector(),
				MinAvailable: "150%",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a malformed bound", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name:           "bad",
				Selector:       appSelector(),
				MaxUnavailable: "one",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects selector In requirements without values", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{
				Name: "bad-selector",
				Selector: &KubernetesPodDisruptionBudgetLabelSelector{
					MatchExpressions: []*KubernetesPodDisruptionBudgetLabelSelectorRequirement{{
						Key:      "tier",
						Operator: "In",
					}},
				},
				MinAvailable: "1",
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a missing name", func() {
			spec := &KubernetesPodDisruptionBudgetSpec{Selector: appSelector(), MinAvailable: "1"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
