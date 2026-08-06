package kubernetesrbacv1alpha1

import (
	"testing"

	"buf.build/go/protovalidate"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	foreignkeyv1 "github.com/plantonhq/planton/apis/dev/planton/shared/foreignkey/v1"
)

func literal(value string) *foreignkeyv1.StringValueOrRef {
	return &foreignkeyv1.StringValueOrRef{
		LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: value},
	}
}

func namespaceScope(namespace string) *KubernetesRbacSpec_NamespaceScope {
	return &KubernetesRbacSpec_NamespaceScope{
		NamespaceScope: &KubernetesRbacNamespaceScope{
			Namespace: literal(namespace),
		},
	}
}

func clusterScope() *KubernetesRbacSpec_ClusterScope {
	return &KubernetesRbacSpec_ClusterScope{
		ClusterScope: &KubernetesRbacClusterScope{},
	}
}

func podReaderRule() *KubernetesRbacPolicyRule {
	return &KubernetesRbacPolicyRule{
		Verbs:     []string{"get", "list", "watch"},
		ApiGroups: []string{""},
		Resources: []string{"pods"},
	}
}

func TestKubernetesRbacSpec(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "KubernetesRbacSpec Validation Suite")
}

var _ = ginkgo.Describe("KubernetesRbacSpec validations", func() {

	ginkgo.Context("When valid specs are provided", func() {

		ginkgo.It("accepts a namespace-scoped grant creating a role for a service account", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_ServiceAccount{
							ServiceAccount: &KubernetesRbacServiceAccountSubject{
								Name: literal("app-identity"),
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster-scoped binding of cluster-admin to a user", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name: "cluster-admin",
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_User{
							User: "jane@example.com",
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a group subject", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name:          "view",
						IsClusterRole: true,
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_Group{
							Group: "system:authenticated",
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a create_role grant with no subjects (role definition only)", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster-scoped aggregated ClusterRole with empty rules", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchLabels: map[string]string{
										"rbac.example.com/aggregate-to-monitoring": "true",
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts match expressions honoring the operator/values contract", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchExpressions: []*KubernetesRbacLabelSelectorRequirement{
										{
											Key:      "tier",
											Operator: "In",
											Values:   []string{"monitoring", "logging"},
										},
										{
											Key:      "rbac.example.com/aggregate",
											Operator: "Exists",
										},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster-scoped rule with non_resource_urls", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								Verbs:           []string{"get"},
								NonResourceUrls: []string{"/healthz", "/metrics"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a cluster-scoped service account subject with an explicit namespace", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name: "view",
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_ServiceAccount{
							ServiceAccount: &KubernetesRbacServiceAccountSubject{
								Name:      literal("app-identity"),
								Namespace: literal("prod"),
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a rule with resource_names restricting objects", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								Verbs:         []string{"get"},
								ApiGroups:     []string{""},
								Resources:     []string{"configmaps"},
								ResourceNames: []string{"app-config"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})

		ginkgo.It("accepts a create_role with an explicit role name", func() {
			roleName := "pod-reader"
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Name:  &roleName,
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).To(gomega.BeNil())
		})
	})

	ginkgo.Context("When invalid specs are provided", func() {

		ginkgo.It("rejects a grant with no scope", func() {
			spec := &KubernetesRbacSpec{
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a grant with no role", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an existing_role grant with no subjects", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name: "view",
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an aggregation_rule in namespace scope", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchLabels: map[string]string{"aggregate": "true"},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects non_resource_urls in a namespace-scoped rule", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								Verbs:           []string{"get"},
								NonResourceUrls: []string{"/healthz"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a cluster-scoped service account subject without a namespace", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name: "view",
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_ServiceAccount{
							ServiceAccount: &KubernetesRbacServiceAccountSubject{
								Name: literal("app-identity"),
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a rule with no verbs", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								ApiGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a rule granting both resources and non_resource_urls", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								Verbs:           []string{"get"},
								ApiGroups:       []string{""},
								Resources:       []string{"pods"},
								NonResourceUrls: []string{"/healthz"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a rule granting neither resources nor non_resource_urls", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{
							{
								Verbs: []string{"get"},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a create_role with no rules and no aggregation_rule", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an aggregation_rule with no selectors", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an In match expression with empty values", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchExpressions: []*KubernetesRbacLabelSelectorRequirement{
										{
											Key:      "tier",
											Operator: "In",
										},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an Exists match expression with values", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchExpressions: []*KubernetesRbacLabelSelectorRequirement{
										{
											Key:      "tier",
											Operator: "Exists",
											Values:   []string{"monitoring"},
										},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a match expression with an unknown operator", func() {
			spec := &KubernetesRbacSpec{
				Scope: clusterScope(),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						AggregationRule: &KubernetesRbacAggregationRule{
							ClusterRoleSelectors: []*KubernetesRbacLabelSelector{
								{
									MatchExpressions: []*KubernetesRbacLabelSelectorRequirement{
										{
											Key:      "tier",
											Operator: "Equals",
											Values:   []string{"monitoring"},
										},
									},
								},
							},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an empty subject message", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a service account subject without a name", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_ServiceAccount{
							ServiceAccount: &KubernetesRbacServiceAccountSubject{},
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects an existing_role with an empty name", func() {
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_ExistingRole{
					ExistingRole: &KubernetesRbacExistingRole{
						Name: "",
					},
				},
				Subjects: []*KubernetesRbacSubject{
					{
						Subject: &KubernetesRbacSubject_User{
							User: "jane@example.com",
						},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})

		ginkgo.It("rejects a create_role name that is not a DNS subdomain", func() {
			roleName := "Bad_Role_Name"
			spec := &KubernetesRbacSpec{
				Scope: namespaceScope("prod"),
				Role: &KubernetesRbacSpec_CreateRole{
					CreateRole: &KubernetesRbacRoleDefinition{
						Name:  &roleName,
						Rules: []*KubernetesRbacPolicyRule{podReaderRule()},
					},
				},
			}
			err := protovalidate.Validate(spec)
			gomega.Expect(err).ToNot(gomega.BeNil())
		})
	})
})
