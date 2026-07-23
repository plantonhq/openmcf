package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createResourceQuota creates the core/v1 ResourceQuota resource carrying the
// aggregate caps, scopes, and scope selector.
func createResourceQuota(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.ResourceQuota, error) {
	spec := locals.Spec

	quotaSpecArgs := &kubernetescorev1.ResourceQuotaSpecArgs{
		Hard: pulumi.ToStringMap(spec.GetHard()),
	}

	if len(spec.GetScopes()) > 0 {
		scopes := make([]string, 0, len(spec.GetScopes()))
		for _, s := range spec.GetScopes() {
			scopes = append(scopes, scopeApiString(s))
		}
		quotaSpecArgs.Scopes = pulumi.ToStringArray(scopes)
	}

	if len(spec.GetScopeSelector()) > 0 {
		var matchExpressions kubernetescorev1.ScopedResourceSelectorRequirementArray
		for _, req := range spec.GetScopeSelector() {
			exprArgs := &kubernetescorev1.ScopedResourceSelectorRequirementArgs{
				ScopeName: pulumi.String(scopeApiString(req.GetScopeName())),
				Operator:  pulumi.String(req.GetOperator()),
			}
			if len(req.GetValues()) > 0 {
				exprArgs.Values = pulumi.ToStringArray(req.GetValues())
			}
			matchExpressions = append(matchExpressions, exprArgs)
		}
		quotaSpecArgs.ScopeSelector = &kubernetescorev1.ScopeSelectorArgs{
			MatchExpressions: matchExpressions,
		}
	}

	resourceQuota, err := kubernetescorev1.NewResourceQuota(
		ctx,
		locals.Name,
		&kubernetescorev1.ResourceQuotaArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: quotaSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create resource quota %s/%s", locals.Namespace, locals.Name)
	}

	return resourceQuota, nil
}

// createLimitRange creates the companion core/v1 LimitRange when the spec
// carries limit_defaults. It shares the quota's name and namespace — one
// governance pair, one identity — which is also what keeps a compute quota
// livable (workloads that omit requests/limits inherit the defaults instead
// of being rejected by the quota's admission check).
func createLimitRange(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetescorev1.LimitRange, error) {
	spec := locals.Spec

	var limitItems kubernetescorev1.LimitRangeItemArray
	for _, item := range spec.GetLimitDefaults() {
		itemArgs := &kubernetescorev1.LimitRangeItemArgs{
			Type: pulumi.String(limitTypeApiString(item.GetType())),
		}
		if len(item.GetMax()) > 0 {
			itemArgs.Max = pulumi.ToStringMap(item.GetMax())
		}
		if len(item.GetMin()) > 0 {
			itemArgs.Min = pulumi.ToStringMap(item.GetMin())
		}
		if len(item.GetDefaultLimit()) > 0 {
			itemArgs.Default = pulumi.ToStringMap(item.GetDefaultLimit())
		}
		if len(item.GetDefaultRequest()) > 0 {
			itemArgs.DefaultRequest = pulumi.ToStringMap(item.GetDefaultRequest())
		}
		if len(item.GetMaxLimitRequestRatio()) > 0 {
			itemArgs.MaxLimitRequestRatio = pulumi.ToStringMap(item.GetMaxLimitRequestRatio())
		}
		limitItems = append(limitItems, itemArgs)
	}

	limitRange, err := kubernetescorev1.NewLimitRange(
		ctx,
		locals.LimitRangeName,
		&kubernetescorev1.LimitRangeArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.LimitRangeName),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: &kubernetescorev1.LimitRangeSpecArgs{
				Limits: limitItems,
			},
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create limit range %s/%s", locals.Namespace, locals.LimitRangeName)
	}

	return limitRange, nil
}
