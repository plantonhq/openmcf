package module

import (
	kuberneteskarpenternodepoolv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kuberneteskarpenternodepool/v1alpha1"
	karpenterv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/karpenter/kubernetes/karpenter/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// defaultConsolidateAfter is the proto-declared default ("0s", immediate
// consolidation). The CRD marks consolidateAfter REQUIRED inside the
// disruption object, so whenever the disruption block is rendered the field
// must carry a value.
const defaultConsolidateAfter = "0s"

// buildDisruption maps the disruption policy — consolidation behavior and
// the disruption budgets — onto the typed crd2pulumi disruption args.
// consolidationPolicy is only rendered when set (the apiserver defaults it
// to "WhenEmptyOrUnderutilized"); consolidateAfter always renders because
// the CRD requires it whenever disruption is present.
func buildDisruption(disruption *kuberneteskarpenternodepoolv1alpha1.KubernetesKarpenterNodePoolDisruption) karpenterv1.NodePoolSpecDisruptionArgs {
	consolidateAfter := disruption.GetConsolidateAfter()
	if consolidateAfter == "" {
		consolidateAfter = defaultConsolidateAfter
	}

	args := karpenterv1.NodePoolSpecDisruptionArgs{
		ConsolidateAfter: pulumi.String(consolidateAfter),
	}

	if consolidationPolicy := disruption.GetConsolidationPolicy(); consolidationPolicy != "" {
		args.ConsolidationPolicy = pulumi.String(consolidationPolicy)
	}

	if budgets := disruption.GetBudgets(); len(budgets) > 0 {
		args.Budgets = buildBudgets(budgets)
	}

	return args
}

// buildBudgets maps the disruption budgets. nodes is required and always
// set; schedule and duration are a set-together pair (protovalidate CEL
// mirrors the CRD rule) rendered only when present; reasons is only
// rendered when non-empty — an absent list means all reasons.
func buildBudgets(budgets []*kuberneteskarpenternodepoolv1alpha1.KubernetesKarpenterNodePoolDisruptionBudget) karpenterv1.NodePoolSpecDisruptionBudgetsArray {
	arr := karpenterv1.NodePoolSpecDisruptionBudgetsArray{}
	for _, budget := range budgets {
		args := karpenterv1.NodePoolSpecDisruptionBudgetsArgs{
			Nodes: pulumi.String(budget.GetNodes()),
		}
		if budget.Schedule != nil {
			args.Schedule = pulumi.String(budget.GetSchedule())
		}
		if budget.Duration != nil {
			args.Duration = pulumi.String(budget.GetDuration())
		}
		if reasons := budget.GetReasons(); len(reasons) > 0 {
			args.Reasons = pulumi.ToStringArray(reasons)
		}
		arr = append(arr, args)
	}
	return arr
}
