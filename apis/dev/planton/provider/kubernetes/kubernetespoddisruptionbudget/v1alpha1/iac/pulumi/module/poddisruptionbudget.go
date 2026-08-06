package module

import (
	"strconv"

	"github.com/pkg/errors"
	kubernetespoddisruptionbudgetv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetespoddisruptionbudget/v1alpha1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	kubernetespolicyv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/policy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createPodDisruptionBudget creates the policy/v1 PodDisruptionBudget.
//
// The selector is always sent (it is required in the spec): an empty selector
// block is the "all pods in the namespace" form — deliberately explicit,
// because a policy/v1 budget with a NULL selector matches no pods at all.
// unhealthyPodEvictionPolicy is always sent with the server default applied
// module-side, keeping both engines' submitted objects identical.
func createPodDisruptionBudget(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*kubernetespolicyv1.PodDisruptionBudget, error) {
	spec := locals.Spec

	// PARITY-EXCEPTION: the terraform kubernetes provider's PDB resource
	// cannot express spec.unhealthyPodEvictionPolicy at all; its module fails
	// the plan with a precondition when the spec asks for always_allow. This
	// engine sends the field explicitly (server default IfHealthyBudget when
	// the spec omits it), so the deployed object is identical across engines
	// for every spec the terraform module accepts.
	budgetSpecArgs := &kubernetespolicyv1.PodDisruptionBudgetSpecArgs{
		Selector:                   buildLabelSelector(spec.GetSelector()),
		UnhealthyPodEvictionPolicy: pulumi.String(locals.UnhealthyPodEvictionPolicy),
	}

	// min_available/max_unavailable are IntOrString upstream: a numeric
	// string is an absolute count, a "%"-suffixed one a percentage. The spec
	// enforces exactly one is set.
	if spec.GetMinAvailable() != "" {
		budgetSpecArgs.MinAvailable = intOrStringValue(spec.GetMinAvailable())
	}
	if spec.GetMaxUnavailable() != "" {
		budgetSpecArgs.MaxUnavailable = intOrStringValue(spec.GetMaxUnavailable())
	}

	podDisruptionBudget, err := kubernetespolicyv1.NewPodDisruptionBudget(
		ctx,
		locals.Name,
		&kubernetespolicyv1.PodDisruptionBudgetArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: budgetSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create pod disruption budget %s/%s", locals.Namespace, locals.Name)
	}

	return podDisruptionBudget, nil
}

// intOrStringValue renders the IntOrString wire form: numeric strings as
// integers, percentages as strings.
func intOrStringValue(v string) pulumi.Input {
	if num, err := strconv.Atoi(v); err == nil {
		return pulumi.Int(num)
	}
	return pulumi.String(v)
}

// buildLabelSelector converts the proto label selector into Pulumi args. An
// empty selector renders as the EMPTY selector — "all pods in the namespace".
func buildLabelSelector(s *kubernetespoddisruptionbudgetv1alpha1.KubernetesPodDisruptionBudgetLabelSelector) *metav1.LabelSelectorArgs {
	selectorArgs := &metav1.LabelSelectorArgs{}
	if s == nil {
		return selectorArgs
	}
	if len(s.GetMatchLabels()) > 0 {
		selectorArgs.MatchLabels = pulumi.ToStringMap(s.GetMatchLabels())
	}
	if len(s.GetMatchExpressions()) > 0 {
		var exprArray metav1.LabelSelectorRequirementArray
		for _, e := range s.GetMatchExpressions() {
			exprArgs := &metav1.LabelSelectorRequirementArgs{
				Key:      pulumi.String(e.GetKey()),
				Operator: pulumi.String(e.GetOperator()),
			}
			if len(e.GetValues()) > 0 {
				exprArgs.Values = pulumi.ToStringArray(e.GetValues())
			}
			exprArray = append(exprArray, exprArgs)
		}
		selectorArgs.MatchExpressions = exprArray
	}
	return selectorArgs
}
