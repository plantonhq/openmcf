package module

import (
	"github.com/pkg/errors"
	kubernetesrbacv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesrbac/v1alpha1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	rbacv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/rbac/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createRole creates the Role (namespace scope) or ClusterRole (cluster scope)
// when the spec defines one via create_role. Returns nil (and no error) when the
// grant binds to an existing role instead — nothing to create on the role side.
func createRole(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (pulumi.Resource, error) {
	createRoleSpec := locals.Spec.GetCreateRole()
	if createRoleSpec == nil {
		return nil, nil
	}

	metadata := &metav1.ObjectMetaArgs{
		Name:        pulumi.String(locals.RoleName),
		Labels:      pulumi.ToStringMap(locals.Labels),
		Annotations: pulumi.ToStringMap(locals.Annotations),
	}

	rules := buildPolicyRules(createRoleSpec.GetRules())

	if locals.IsNamespaceScoped {
		// Namespaced Role: rules only. The Kubernetes Role type has no
		// aggregationRule field — spec validation already rejects aggregation
		// in namespace scope, so it is simply never set here.
		metadata.Namespace = pulumi.String(locals.Namespace)
		role, err := rbacv1.NewRole(
			ctx,
			locals.RoleName,
			&rbacv1.RoleArgs{
				Metadata: metadata,
				Rules:    rules,
			},
			pulumi.Provider(provider),
		)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create role %s", locals.RoleName)
		}
		return role, nil
	}

	// ClusterRole: rules plus, when set, the aggregation rule. With aggregation
	// the controller continuously composes the rules from every ClusterRole
	// matching the selectors, and any directly listed rules are controller-managed.
	clusterRole, err := rbacv1.NewClusterRole(
		ctx,
		locals.RoleName,
		&rbacv1.ClusterRoleArgs{
			Metadata:        metadata,
			Rules:           rules,
			AggregationRule: buildAggregationRule(createRoleSpec.GetAggregationRule()),
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster role %s", locals.RoleName)
	}
	return clusterRole, nil
}

// buildPolicyRules maps spec policy rules onto Kubernetes rbac/v1 PolicyRules.
// The spec mirrors the Kubernetes type field-for-field, so this is a direct copy:
// each rule independently grants a set of verbs over resources (or non-resource
// URLs); there is no ordering and no deny semantics.
func buildPolicyRules(specRules []*kubernetesrbacv1alpha1.KubernetesRbacPolicyRule) rbacv1.PolicyRuleArray {
	rules := make(rbacv1.PolicyRuleArray, 0, len(specRules))
	for _, specRule := range specRules {
		rules = append(rules, rbacv1.PolicyRuleArgs{
			Verbs:           pulumi.ToStringArray(specRule.GetVerbs()),
			ApiGroups:       pulumi.ToStringArray(specRule.GetApiGroups()),
			Resources:       pulumi.ToStringArray(specRule.GetResources()),
			ResourceNames:   pulumi.ToStringArray(specRule.GetResourceNames()),
			NonResourceURLs: pulumi.ToStringArray(specRule.GetNonResourceUrls()),
		})
	}
	return rules
}

// buildAggregationRule maps the spec aggregation rule onto the Kubernetes
// AggregationRule: a list of label selectors, where a ClusterRole matching ANY
// selector contributes its rules to the aggregated ClusterRole.
func buildAggregationRule(specAggregation *kubernetesrbacv1alpha1.KubernetesRbacAggregationRule) rbacv1.AggregationRulePtrInput {
	if specAggregation == nil {
		return nil
	}

	selectors := make(metav1.LabelSelectorArray, 0, len(specAggregation.GetClusterRoleSelectors()))
	for _, specSelector := range specAggregation.GetClusterRoleSelectors() {
		// Within one selector, matchLabels and matchExpressions AND together;
		// an empty selector matches every ClusterRole.
		expressions := make(metav1.LabelSelectorRequirementArray, 0, len(specSelector.GetMatchExpressions()))
		for _, specExpression := range specSelector.GetMatchExpressions() {
			expressions = append(expressions, metav1.LabelSelectorRequirementArgs{
				Key:      pulumi.String(specExpression.GetKey()),
				Operator: pulumi.String(specExpression.GetOperator()),
				Values:   pulumi.ToStringArray(specExpression.GetValues()),
			})
		}
		selectors = append(selectors, metav1.LabelSelectorArgs{
			MatchLabels:      pulumi.ToStringMap(specSelector.GetMatchLabels()),
			MatchExpressions: expressions,
		})
	}

	return rbacv1.AggregationRuleArgs{
		ClusterRoleSelectors: selectors,
	}
}
