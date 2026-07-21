package module

import (
	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetescronjobv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetescronjob/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/workloadpod"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/batch/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cronJob renders the batch/v1 CronJob. Scheduling controls (schedule, time
// zone, concurrency, history retention) map at the top level; the Job stamped
// out at each run comes from spec.job_template, with its pod template
// assembled through the shared workloadpod builders so container and pod
// semantics match every other workload kind exactly.
//
// CronJobs are mutable — the controller only reads the template when stamping
// out a run — so unlike Jobs, no replace-on-change resource options are needed.
func cronJob(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, createdImagePullSecret *kubernetescorev1.Secret,
	namespaceDeps []pulumi.ResourceOption) (*batchv1.CronJob, error) {

	spec := locals.KubernetesCronJob.Spec

	cronJobSpecArgs := &batchv1.CronJobSpecArgs{
		Schedule: pulumi.String(spec.Schedule),
		JobTemplate: &batchv1.JobTemplateSpecArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Labels: pulumi.ToStringMap(locals.Labels),
			},
			Spec: buildJobSpec(spec.JobTemplate, locals, createdImagePullSecret),
		},
	}

	if spec.GetTimeZone() != "" {
		cronJobSpecArgs.TimeZone = pulumi.String(spec.GetTimeZone())
	}
	if spec.StartingDeadlineSeconds != nil {
		cronJobSpecArgs.StartingDeadlineSeconds = pulumi.Int(int(*spec.StartingDeadlineSeconds))
	}

	// Concurrency policy defaults to "Forbid" — deliberately safer than
	// upstream's "Allow". Overlapping cron runs are the classic
	// scheduled-workload incident (two backups writing the same target, two
	// migrations racing), so the spec documents Forbid as OUR default and the
	// module applies it explicitly rather than letting Kubernetes fall back to
	// Allow when the field is unset.
	concurrencyPolicy := spec.GetConcurrencyPolicy()
	if concurrencyPolicy == "" {
		concurrencyPolicy = "Forbid"
	}
	cronJobSpecArgs.ConcurrencyPolicy = pulumi.String(concurrencyPolicy)

	if spec.Suspend {
		cronJobSpecArgs.Suspend = pulumi.Bool(true)
	}
	if spec.SuccessfulJobsHistoryLimit != nil {
		cronJobSpecArgs.SuccessfulJobsHistoryLimit = pulumi.Int(int(*spec.SuccessfulJobsHistoryLimit))
	}
	if spec.FailedJobsHistoryLimit != nil {
		cronJobSpecArgs.FailedJobsHistoryLimit = pulumi.Int(int(*spec.FailedJobsHistoryLimit))
	}

	cronJobArgs := &batchv1.CronJobArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubernetesCronJob.Metadata.Name),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Spec: cronJobSpecArgs,
	}

	cronJobOpts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)
	createdCronJob, err := batchv1.NewCronJob(ctx,
		locals.KubernetesCronJob.Metadata.Name,
		cronJobArgs,
		cronJobOpts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create cron job")
	}

	return createdCronJob, nil
}

// buildJobSpec renders the JobSpec stamped out at each scheduled run: batch
// controls from the job template plus the pod template built through the
// shared workloadpod builders — the same mapping the KubernetesJob module
// applies to its flat spec.
func buildJobSpec(jobTemplate *kubernetescronjobv1.KubernetesCronJobJobTemplate, locals *Locals,
	createdImagePullSecret *kubernetescorev1.Secret) *batchv1.JobSpecArgs {

	containers := workloadpod.BuildContainers(
		jobTemplate.Container.App, jobTemplate.Container.Sidecars,
		// Default name for the main container when the spec omits one. "app" matches
		// the labeling convention and reads correctly in kubectl output.
		"app",
		locals.EnvSecretName,
	)

	// The pod volume list is the union of every container's mounts — app,
	// sidecars, and init containers all contribute.
	volumeSourceContainers := make([]*kubernetesv1.WorkloadContainer, 0, 2+len(jobTemplate.Container.Sidecars))
	volumeSourceContainers = append(volumeSourceContainers, jobTemplate.Container.App)
	volumeSourceContainers = append(volumeSourceContainers, jobTemplate.Container.Sidecars...)
	if jobTemplate.Pod != nil {
		volumeSourceContainers = append(volumeSourceContainers, jobTemplate.Pod.InitContainers...)
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
	restartPolicy := jobTemplate.GetRestartPolicy()
	if restartPolicy == "" {
		restartPolicy = "Never"
	}

	podTemplate := workloadpod.BuildPodTemplateSpec(jobTemplate.Pod, workloadpod.PodTemplateInputs{
		Labels:               locals.Labels,
		Containers:           containers,
		Volumes:              volumes,
		RestartPolicy:        restartPolicy,
		ImagePullSecretNames: workloadpod.ResolveImagePullSecretNames(jobTemplate.Pod, moduleCreatedPullSecret),
	}, locals.EnvSecretName)

	// Selector is deliberately NOT set: the Job controller generates its own
	// unique selector (controller-uid) for each stamped-out Job and manages pod
	// ownership through it.
	jobSpecArgs := &batchv1.JobSpecArgs{
		Template: podTemplate,
	}

	// Batch controls: the generated optional fields are pointers — apply the
	// value when set, otherwise omit the field and let Kubernetes apply its own
	// defaults (parallelism 1, completions 1, backoffLimit 6, NonIndexed).
	if jobTemplate.Parallelism != nil {
		jobSpecArgs.Parallelism = pulumi.Int(int(*jobTemplate.Parallelism))
	}
	if jobTemplate.Completions != nil {
		jobSpecArgs.Completions = pulumi.Int(int(*jobTemplate.Completions))
	}
	if jobTemplate.GetCompletionMode() != "" {
		jobSpecArgs.CompletionMode = pulumi.String(jobTemplate.GetCompletionMode())
	}
	if jobTemplate.BackoffLimit != nil {
		jobSpecArgs.BackoffLimit = pulumi.Int(int(*jobTemplate.BackoffLimit))
	}
	if jobTemplate.BackoffLimitPerIndex != nil {
		jobSpecArgs.BackoffLimitPerIndex = pulumi.Int(int(*jobTemplate.BackoffLimitPerIndex))
	}
	if jobTemplate.MaxFailedIndexes != nil {
		jobSpecArgs.MaxFailedIndexes = pulumi.Int(int(*jobTemplate.MaxFailedIndexes))
	}
	if jobTemplate.ActiveDeadlineSeconds != nil {
		jobSpecArgs.ActiveDeadlineSeconds = pulumi.Int(int(*jobTemplate.ActiveDeadlineSeconds))
	}
	if jobTemplate.TtlSecondsAfterFinished != nil {
		jobSpecArgs.TtlSecondsAfterFinished = pulumi.Int(int(*jobTemplate.TtlSecondsAfterFinished))
	}
	if jobTemplate.PodFailurePolicy != nil {
		jobSpecArgs.PodFailurePolicy = buildPodFailurePolicy(jobTemplate.PodFailurePolicy)
	}
	// PARITY-EXCEPTION: the Terraform kubernetes provider's job spec has no
	// `success_policy` block, so the Terraform module cannot render this field.
	// A spec with success_policy deploys identically through Pulumi; on
	// Terraform each run's Job falls back to the default success criterion
	// (all `completions` pods must succeed).
	if jobTemplate.SuccessPolicy != nil {
		jobSpecArgs.SuccessPolicy = buildSuccessPolicy(jobTemplate.SuccessPolicy)
	}

	return jobSpecArgs
}

// buildPodFailurePolicy converts the job template's pod failure policy into
// Pulumi args. Rule order is preserved — Kubernetes evaluates rules in order
// and the first match wins, so order is semantics, not cosmetics.
func buildPodFailurePolicy(policy *kubernetescronjobv1.KubernetesCronJobPodFailurePolicy) *batchv1.PodFailurePolicyArgs {
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

// buildSuccessPolicy converts the job template's success policy into Pulumi
// args. The run's Job succeeds as soon as any rule is satisfied.
func buildSuccessPolicy(policy *kubernetescronjobv1.KubernetesCronJobSuccessPolicy) *batchv1.SuccessPolicyArgs {
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
