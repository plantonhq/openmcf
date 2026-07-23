package module

import (
	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetesjobv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesjob/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/workloadpod"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/batch/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// job renders the batch/v1 Job. The pod template is assembled through the
// shared workloadpod builders so container and pod semantics match every other
// workload kind exactly; only Job-specific mechanics (parallelism, retry
// budgets, completion tracking, failure/success policies) live here.
func job(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, createdImagePullSecret *kubernetescorev1.Secret,
	namespaceDeps []pulumi.ResourceOption) (*batchv1.Job, error) {

	spec := locals.KubernetesJob.Spec

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

	// Restart policy defaults to "Never": one pod per attempt, so failed pods
	// survive for post-mortem inspection — and pod_failure_policy requires it.
	restartPolicy := spec.GetRestartPolicy()
	if restartPolicy == "" {
		restartPolicy = "Never"
	}

	podTemplate := workloadpod.BuildPodTemplateSpec(spec.Pod, workloadpod.PodTemplateInputs{
		Labels:               locals.Labels,
		Containers:           containers,
		Volumes:              volumes,
		RestartPolicy:        restartPolicy,
		ImagePullSecretNames: workloadpod.ResolveImagePullSecretNames(spec.Pod, moduleCreatedPullSecret),
	}, locals.EnvSecretName)

	// Selector is deliberately NOT set: the Job controller generates its own
	// unique selector (controller-uid) and stamps it on the pod template.
	// Supplying one requires manualSelector and a non-unique choice makes the
	// controller adopt or fight over unrelated pods — our selector labels are
	// for humans and tooling, not for the controller.
	jobSpecArgs := &batchv1.JobSpecArgs{
		Template: podTemplate,
	}

	// Batch controls: the generated optional fields are pointers — apply the
	// value when set, otherwise omit the field and let Kubernetes apply its own
	// defaults (parallelism 1, completions 1, backoffLimit 6, NonIndexed).
	if spec.Parallelism != nil {
		jobSpecArgs.Parallelism = pulumi.Int(int(*spec.Parallelism))
	}
	if spec.Completions != nil {
		jobSpecArgs.Completions = pulumi.Int(int(*spec.Completions))
	}
	if spec.GetCompletionMode() != "" {
		jobSpecArgs.CompletionMode = pulumi.String(spec.GetCompletionMode())
	}
	if spec.BackoffLimit != nil {
		jobSpecArgs.BackoffLimit = pulumi.Int(int(*spec.BackoffLimit))
	}
	if spec.BackoffLimitPerIndex != nil {
		jobSpecArgs.BackoffLimitPerIndex = pulumi.Int(int(*spec.BackoffLimitPerIndex))
	}
	if spec.MaxFailedIndexes != nil {
		jobSpecArgs.MaxFailedIndexes = pulumi.Int(int(*spec.MaxFailedIndexes))
	}
	if spec.ActiveDeadlineSeconds != nil {
		jobSpecArgs.ActiveDeadlineSeconds = pulumi.Int(int(*spec.ActiveDeadlineSeconds))
	}
	if spec.TtlSecondsAfterFinished != nil {
		jobSpecArgs.TtlSecondsAfterFinished = pulumi.Int(int(*spec.TtlSecondsAfterFinished))
	}
	// PARITY-EXCEPTION: the Terraform kubernetes provider's job spec has no
	// `suspend` argument, so the Terraform module cannot render this field. A
	// spec with suspend=true deploys identically through Pulumi; on Terraform
	// the Job starts immediately and must be suspended out-of-band.
	if spec.Suspend {
		jobSpecArgs.Suspend = pulumi.Bool(true)
	}
	if spec.PodFailurePolicy != nil {
		jobSpecArgs.PodFailurePolicy = buildPodFailurePolicy(spec.PodFailurePolicy)
	}
	// PARITY-EXCEPTION: the Terraform kubernetes provider's job spec has no
	// `success_policy` block, so the Terraform module cannot render this field.
	// A spec with success_policy deploys identically through Pulumi; on
	// Terraform the Job falls back to the default success criterion (all
	// `completions` pods must succeed).
	if spec.SuccessPolicy != nil {
		jobSpecArgs.SuccessPolicy = buildSuccessPolicy(spec.SuccessPolicy)
	}

	jobArgs := &batchv1.JobArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubernetesJob.Metadata.Name),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Spec: jobSpecArgs,
	}

	jobOpts := append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider),
		// Job specs are immutable after creation (only parallelism and suspend are
		// mutable), so any spec change must replace the Job rather than patch it.
		pulumi.ReplaceOnChanges([]string{"spec"}),
		// The Job carries a fixed name; the old object must be deleted before the
		// replacement is created or the names collide.
		pulumi.DeleteBeforeReplace(true),
	}, namespaceDeps...)
	createdJob, err := batchv1.NewJob(ctx,
		locals.KubernetesJob.Metadata.Name,
		jobArgs,
		jobOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create job")
	}

	return createdJob, nil
}

// buildPodFailurePolicy converts the spec's pod failure policy into Pulumi
// args. Rule order is preserved — Kubernetes evaluates rules in order and the
// first match wins, so order is semantics, not cosmetics.
func buildPodFailurePolicy(policy *kubernetesjobv1.KubernetesJobPodFailurePolicy) *batchv1.PodFailurePolicyArgs {
	rules := make(batchv1.PodFailurePolicyRuleArray, 0, len(policy.Rules))
	for _, rule := range policy.Rules {
		ruleArgs := &batchv1.PodFailurePolicyRuleArgs{
			Action: pulumi.String(rule.Action),
		}

		if rule.OnExitCodes != nil {
			onExitCodes := &batchv1.PodFailurePolicyOnExitCodesRequirementArgs{
				Operator: pulumi.String(rule.OnExitCodes.Operator),
				Values:   pulumi.ToIntArray(toIntSlice(rule.OnExitCodes.Values)),
			}
			if rule.OnExitCodes.ContainerName != "" {
				onExitCodes.ContainerName = pulumi.String(rule.OnExitCodes.ContainerName)
			}
			ruleArgs.OnExitCodes = onExitCodes
		}

		if len(rule.OnPodConditions) > 0 {
			conditions := make(batchv1.PodFailurePolicyOnPodConditionsPatternArray, 0, len(rule.OnPodConditions))
			for _, condition := range rule.OnPodConditions {
				// Status defaults to "True" — the API requires the field, and "True"
				// is the documented default in the spec.
				status := condition.GetStatus()
				if status == "" {
					status = "True"
				}
				conditions = append(conditions, &batchv1.PodFailurePolicyOnPodConditionsPatternArgs{
					Type:   pulumi.String(condition.Type),
					Status: pulumi.String(status),
				})
			}
			ruleArgs.OnPodConditions = conditions
		}

		rules = append(rules, ruleArgs)
	}

	return &batchv1.PodFailurePolicyArgs{Rules: rules}
}

// buildSuccessPolicy converts the spec's success policy into Pulumi args. The
// Job succeeds as soon as any rule is satisfied.
func buildSuccessPolicy(policy *kubernetesjobv1.KubernetesJobSuccessPolicy) *batchv1.SuccessPolicyArgs {
	rules := make(batchv1.SuccessPolicyRuleArray, 0, len(policy.Rules))
	for _, rule := range policy.Rules {
		ruleArgs := &batchv1.SuccessPolicyRuleArgs{}
		if rule.SucceededIndexes != "" {
			ruleArgs.SucceededIndexes = pulumi.String(rule.SucceededIndexes)
		}
		if rule.SucceededCount != nil {
			ruleArgs.SucceededCount = pulumi.Int(int(*rule.SucceededCount))
		}
		rules = append(rules, ruleArgs)
	}

	return &batchv1.SuccessPolicyArgs{Rules: rules}
}

// toIntSlice widens proto int32 exit codes to the int slice Pulumi's array
// helper accepts.
func toIntSlice(values []int32) []int {
	result := make([]int, 0, len(values))
	for _, v := range values {
		result = append(result, int(v))
	}
	return result
}
