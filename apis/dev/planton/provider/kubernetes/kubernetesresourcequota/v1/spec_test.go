package kubernetesresourcequotav1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func TestKubernetesResourceQuotaSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesResourceQuotaSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesResourceQuotaSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a minimal compute quota", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "team-quota",
				Hard: map[string]string{"requests.cpu": "10", "requests.memory": "20Gi"},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a namespace provided as a resource reference", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "team-quota",
				Hard: map[string]string{"pods": "50"},
				Namespace: &foreignkeyv1.StringValueOrRef{
					LiteralOrRef: &foreignkeyv1.StringValueOrRef_ValueFrom{
						ValueFrom: &foreignkeyv1.ValueFromRef{Name: "team-namespace"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts object-count and storage caps", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "full-quota",
				Hard: map[string]string{
					"pods":                   "100",
					"services.loadbalancers": "2",
					"persistentvolumeclaims": "20",
					"requests.storage":       "500Gi",
					"count/deployments.apps": "30",
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a best_effort-scoped pods quota", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name:   "besteffort-cap",
				Hard:   map[string]string{"pods": "10"},
				Scopes: []KubernetesResourceQuotaSpec_KubernetesResourceQuotaScope{KubernetesResourceQuotaSpec_best_effort},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a priority-class scope selector with In values", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "critical-quota",
				Hard: map[string]string{"pods": "20"},
				ScopeSelector: []*KubernetesResourceQuotaScopeSelectorRequirement{{
					ScopeName: KubernetesResourceQuotaSpec_priority_class,
					Operator:  "In",
					Values:    []string{"critical"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts limit defaults for containers, pods, and claims", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "governed",
				Hard: map[string]string{"requests.cpu": "10", "limits.cpu": "20"},
				LimitDefaults: []*KubernetesResourceQuotaLimitDefaults{
					{
						Type:           KubernetesResourceQuotaLimitDefaults_container,
						DefaultRequest: map[string]string{"cpu": "100m", "memory": "128Mi"},
						DefaultLimit:   map[string]string{"cpu": "500m", "memory": "512Mi"},
						Max:            map[string]string{"cpu": "2"},
					},
					{
						Type: KubernetesResourceQuotaLimitDefaults_pod,
						Max:  map[string]string{"cpu": "4", "memory": "8Gi"},
					},
					{
						Type: KubernetesResourceQuotaLimitDefaults_persistent_volume_claim,
						Min:  map[string]string{"storage": "1Gi"},
						Max:  map[string]string{"storage": "100Gi"},
					},
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})

		ginkgo.It("accepts a max limit-request ratio bound", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "burst-bound",
				Hard: map[string]string{"pods": "50"},
				LimitDefaults: []*KubernetesResourceQuotaLimitDefaults{{
					Type:                 KubernetesResourceQuotaLimitDefaults_container,
					MaxLimitRequestRatio: map[string]string{"cpu": "4"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a missing name", func() {
			spec := &KubernetesResourceQuotaSpec{Hard: map[string]string{"pods": "10"}}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty hard map", func() {
			spec := &KubernetesResourceQuotaSpec{Name: "empty"}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects conflicting best_effort scopes", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "conflict",
				Hard: map[string]string{"pods": "10"},
				Scopes: []KubernetesResourceQuotaSpec_KubernetesResourceQuotaScope{
					KubernetesResourceQuotaSpec_best_effort,
					KubernetesResourceQuotaSpec_not_best_effort,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects conflicting terminating scopes", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "conflict",
				Hard: map[string]string{"pods": "10"},
				Scopes: []KubernetesResourceQuotaSpec_KubernetesResourceQuotaScope{
					KubernetesResourceQuotaSpec_terminating,
					KubernetesResourceQuotaSpec_not_terminating,
				},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a best_effort scope capping compute resources", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name:   "besteffort-compute",
				Hard:   map[string]string{"requests.cpu": "10"},
				Scopes: []KubernetesResourceQuotaSpec_KubernetesResourceQuotaScope{KubernetesResourceQuotaSpec_best_effort},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a pod-behavior scope selector with a non-Exists operator", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "bad-selector",
				Hard: map[string]string{"pods": "10"},
				ScopeSelector: []*KubernetesResourceQuotaScopeSelectorRequirement{{
					ScopeName: KubernetesResourceQuotaSpec_best_effort,
					Operator:  "In",
					Values:    []string{"x"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an In scope selector without values", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "bad-selector",
				Hard: map[string]string{"pods": "10"},
				ScopeSelector: []*KubernetesResourceQuotaScopeSelectorRequirement{{
					ScopeName: KubernetesResourceQuotaSpec_priority_class,
					Operator:  "In",
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an Exists scope selector carrying values", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "bad-selector",
				Hard: map[string]string{"pods": "10"},
				ScopeSelector: []*KubernetesResourceQuotaScopeSelectorRequirement{{
					ScopeName: KubernetesResourceQuotaSpec_priority_class,
					Operator:  "Exists",
					Values:    []string{"critical"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects defaults on a pod-type limit item", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "bad-defaults",
				Hard: map[string]string{"pods": "10"},
				LimitDefaults: []*KubernetesResourceQuotaLimitDefaults{{
					Type:         KubernetesResourceQuotaLimitDefaults_pod,
					DefaultLimit: map[string]string{"cpu": "1"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects an empty limit item", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "empty-item",
				Hard: map[string]string{"pods": "10"},
				LimitDefaults: []*KubernetesResourceQuotaLimitDefaults{{
					Type: KubernetesResourceQuotaLimitDefaults_container,
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})

		ginkgo.It("rejects a limit item without a type", func() {
			spec := &KubernetesResourceQuotaSpec{
				Name: "no-type",
				Hard: map[string]string{"pods": "10"},
				LimitDefaults: []*KubernetesResourceQuotaLimitDefaults{{
					Max: map[string]string{"cpu": "2"},
				}},
			}
			gomega.Expect(protovalidate.Validate(spec)).NotTo(gomega.BeNil())
		})
	})
})
