package module

import (
	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesdaemonsetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesdaemonset/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/workloadpod"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// daemonSet renders the apps/v1 DaemonSet. The pod template is assembled
// through the shared workloadpod builders so container and pod semantics match
// every other workload kind exactly; only DaemonSet-specific mechanics
// (node-by-node rollout) live here. There is no replica count — the scheduler
// places one pod on every node matching the pod's scheduling rules.
func daemonSet(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, createdImagePullSecret *kubernetescorev1.Secret,
	namespaceDeps []pulumi.ResourceOption) error {

	spec := locals.KubernetesDaemonSet.Spec

	containers := workloadpod.BuildContainers(
		spec.Container.App, spec.Container.Sidecars,
		// Default name for the main container when the spec omits one. "app" matches
		// the selector-label convention and reads correctly in kubectl output.
		"app",
		locals.EnvSecretName,
	)

	// The pod volume list is the union of every container's mounts — app,
	// sidecars, and init containers all contribute.
	volumeSourceContainers := make([]*kubernetesv1.WorkloadContainer, 0, 2+len(spec.Container.Sidecars))
	volumeSourceContainers = append(volumeSourceContainers, spec.Container.App)
	volumeSourceContainers = append(volumeSourceContainers, spec.Container.Sidecars...)
	if spec.Pod != nil {
		volumeSourceContainers = append(volumeSourceContainers, spec.Pod.InitContainers...)
	}
	volumes := workloadpod.BuildVolumes(volumeSourceContainers...)

	// The module-created docker-config secret joins the pod's pull-secret list;
	// spec-listed pull secrets pass through by name.
	moduleCreatedPullSecret := ""
	if createdImagePullSecret != nil {
		moduleCreatedPullSecret = locals.ImagePullSecretName
	}

	podTemplate := workloadpod.BuildPodTemplateSpec(spec.Pod, workloadpod.PodTemplateInputs{
		Labels:               locals.Labels,
		Containers:           containers,
		Volumes:              volumes,
		ImagePullSecretNames: workloadpod.ResolveImagePullSecretNames(spec.Pod, moduleCreatedPullSecret),
	}, locals.EnvSecretName)

	daemonSetSpecArgs := &appsv1.DaemonSetSpecArgs{
		Selector: &metav1.LabelSelectorArgs{
			MatchLabels: pulumi.ToStringMap(locals.SelectorLabels),
		},
		Template:       podTemplate,
		UpdateStrategy: buildUpdateStrategy(spec.UpdateStrategy),
	}

	if spec.MinReadySeconds != nil {
		daemonSetSpecArgs.MinReadySeconds = pulumi.Int(int(*spec.MinReadySeconds))
	}
	if spec.RevisionHistoryLimit != nil {
		daemonSetSpecArgs.RevisionHistoryLimit = pulumi.Int(int(*spec.RevisionHistoryLimit))
	}

	daemonSetArgs := &appsv1.DaemonSetArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubernetesDaemonSet.Metadata.Name),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
			Annotations: pulumi.StringMap{
				// Force server-side apply to take ownership of fields on conflicts —
				// without it, a previously kubectl-managed daemon set blocks updates.
				"pulumi.com/patchForce": pulumi.String("true"),
			},
		},
		Spec: daemonSetSpecArgs,
	}

	dsOpts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)
	_, err := appsv1.NewDaemonSet(ctx,
		locals.KubernetesDaemonSet.Metadata.Name,
		daemonSetArgs,
		dsOpts...)
	if err != nil {
		return errors.Wrap(err, "failed to create daemon set")
	}

	return nil
}

// buildUpdateStrategy converts the update strategy spec into Pulumi args.
// Returns nil when unset so Kubernetes applies its defaults (RollingUpdate,
// at most 1 node's pod unavailable at a time).
func buildUpdateStrategy(protoStrategy *kubernetesdaemonsetv1.KubernetesDaemonSetUpdateStrategy) *appsv1.DaemonSetUpdateStrategyArgs {
	if protoStrategy == nil {
		return nil
	}

	strategyType := protoStrategy.Type
	if strategyType == "" {
		strategyType = "RollingUpdate"
	}

	strategy := &appsv1.DaemonSetUpdateStrategyArgs{
		Type: pulumi.String(strategyType),
	}

	// An OnDelete strategy must NOT carry a rollingUpdate block — the API server
	// rejects the combination.
	if strategyType == "RollingUpdate" {
		rollingUpdate := &appsv1.RollingUpdateDaemonSetArgs{}
		hasRollingConfig := false
		if protoStrategy.MaxUnavailable != "" {
			rollingUpdate.MaxUnavailable = parseIntOrString(protoStrategy.MaxUnavailable)
			hasRollingConfig = true
		}
		if protoStrategy.MaxSurge != "" {
			rollingUpdate.MaxSurge = parseIntOrString(protoStrategy.MaxSurge)
			hasRollingConfig = true
		}
		if hasRollingConfig {
			strategy.RollingUpdate = rollingUpdate
		}
	}

	return strategy
}
