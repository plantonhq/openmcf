package module

import (
	kuberneteskarpenternodepoolv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskarpenternodepool/v1"
	karpenterv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/karpenter/kubernetes/karpenter/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Proto-declared defaults for the nodeClassRef, applied module-side when the
// manifest leaves them unset: the CRD requires group and kind to be present
// and non-empty, so they can never be omitted from the rendered CR.
const (
	defaultNodeClassRefGroup = "karpenter.k8s.aws"
	defaultNodeClassRefKind  = "EC2NodeClass"
)

// buildTemplate maps the NodeClaim template — node metadata, the NodeClass
// reference, scheduling requirements, taints, and node lifetime — onto the
// typed crd2pulumi template args. template.metadata is only rendered when
// labels or annotations exist; expireAfter and terminationGracePeriod are
// only rendered when set, so the apiserver applies the CRD defaults
// (expireAfter: "720h") otherwise.
func buildTemplate(template *kuberneteskarpenternodepoolv1.KubernetesKarpenterNodePoolTemplate) karpenterv1.NodePoolSpecTemplateArgs {
	templateSpec := karpenterv1.NodePoolSpecTemplateSpecArgs{
		NodeClassRef: buildNodeClassRef(template.GetNodeClassRef()),
		Requirements: buildRequirements(template.GetRequirements()),
	}

	if taints := template.GetTaints(); len(taints) > 0 {
		templateSpec.Taints = buildTaints(taints)
	}

	if startupTaints := template.GetStartupTaints(); len(startupTaints) > 0 {
		templateSpec.StartupTaints = buildStartupTaints(startupTaints)
	}

	if expireAfter := template.GetExpireAfter(); expireAfter != "" {
		templateSpec.ExpireAfter = pulumi.String(expireAfter)
	}

	if terminationGracePeriod := template.GetTerminationGracePeriod(); terminationGracePeriod != "" {
		templateSpec.TerminationGracePeriod = pulumi.String(terminationGracePeriod)
	}

	args := karpenterv1.NodePoolSpecTemplateArgs{
		Spec: templateSpec,
	}

	if len(template.GetLabels()) > 0 || len(template.GetAnnotations()) > 0 {
		metadata := karpenterv1.NodePoolSpecTemplateMetadataArgs{}
		if labels := template.GetLabels(); len(labels) > 0 {
			metadata.Labels = pulumi.ToStringMap(labels)
		}
		if annotations := template.GetAnnotations(); len(annotations) > 0 {
			metadata.Annotations = pulumi.ToStringMap(annotations)
		}
		args.Metadata = metadata
	}

	return args
}

// buildNodeClassRef maps the provider-specific NodeClass reference. group
// and kind fall back to the proto-declared AWS defaults when unset — the
// CRD requires all three keys with non-empty values. name is a
// KubernetesKarpenterEc2NodeClass foreign key resolved to its literal value
// before the module runs, so GetValue() returns the final name.
func buildNodeClassRef(nodeClassRef *kuberneteskarpenternodepoolv1.KubernetesKarpenterNodePoolNodeClassRef) karpenterv1.NodePoolSpecTemplateSpecNodeClassRefArgs {
	group := nodeClassRef.GetGroup()
	if nodeClassRef.Group == nil {
		group = defaultNodeClassRefGroup
	}
	kind := nodeClassRef.GetKind()
	if nodeClassRef.Kind == nil {
		kind = defaultNodeClassRefKind
	}

	return karpenterv1.NodePoolSpecTemplateSpecNodeClassRefArgs{
		Group: pulumi.String(group),
		Kind:  pulumi.String(kind),
		Name:  pulumi.String(nodeClassRef.GetName().GetValue()),
	}
}

// buildRequirements maps the scheduling requirements. values is only
// rendered when non-empty (Exists/DoesNotExist require an absent or empty
// list) and minValues only when set — its presence activates the ALPHA
// instance-type-diversity check.
func buildRequirements(requirements []*kuberneteskarpenternodepoolv1.KubernetesKarpenterNodePoolRequirement) karpenterv1.NodePoolSpecTemplateSpecRequirementsArray {
	arr := karpenterv1.NodePoolSpecTemplateSpecRequirementsArray{}
	for _, requirement := range requirements {
		args := karpenterv1.NodePoolSpecTemplateSpecRequirementsArgs{
			Key:      pulumi.String(requirement.GetKey()),
			Operator: pulumi.String(requirement.GetOperator()),
		}
		if values := requirement.GetValues(); len(values) > 0 {
			args.Values = pulumi.ToStringArray(values)
		}
		if requirement.MinValues != nil {
			args.MinValues = pulumi.Int(int(requirement.GetMinValues()))
		}
		arr = append(arr, args)
	}
	return arr
}

// buildTaints maps the always-on node taints. value is optional in the CRD
// and only rendered when set; timeAdded is controller-owned and never set.
func buildTaints(taints []*kuberneteskarpenternodepoolv1.KubernetesKarpenterNodePoolTaint) karpenterv1.NodePoolSpecTemplateSpecTaintsArray {
	arr := karpenterv1.NodePoolSpecTemplateSpecTaintsArray{}
	for _, taint := range taints {
		args := karpenterv1.NodePoolSpecTemplateSpecTaintsArgs{
			Key:    pulumi.String(taint.GetKey()),
			Effect: pulumi.String(taint.GetEffect()),
		}
		if value := taint.GetValue(); value != "" {
			args.Value = pulumi.String(value)
		}
		arr = append(arr, args)
	}
	return arr
}

// buildStartupTaints maps the startup taints (same taint shape, distinct
// typed SDK array type).
func buildStartupTaints(startupTaints []*kuberneteskarpenternodepoolv1.KubernetesKarpenterNodePoolTaint) karpenterv1.NodePoolSpecTemplateSpecStartupTaintsArray {
	arr := karpenterv1.NodePoolSpecTemplateSpecStartupTaintsArray{}
	for _, taint := range startupTaints {
		args := karpenterv1.NodePoolSpecTemplateSpecStartupTaintsArgs{
			Key:    pulumi.String(taint.GetKey()),
			Effect: pulumi.String(taint.GetEffect()),
		}
		if value := taint.GetValue(); value != "" {
			args.Value = pulumi.String(value)
		}
		arr = append(arr, args)
	}
	return arr
}
