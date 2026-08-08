package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/dataproc"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// autoscalingPolicy provisions the Dataproc autoscaling policy. A
// shareable resource: one policy can govern many clusters (each
// attaches it by reference), so scaling behavior is tuned in one place.
// Policy contents are mutable — updating re-tunes every attached
// cluster — but the API refuses to delete a policy while any cluster
// references it.
func autoscalingPolicy(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpDataprocAutoscalingPolicy.Spec

	// Enable the Dataproc API — the control plane that owns the policy.
	// disable_on_destroy stays false: tearing down one policy must never
	// disable the API for everything else in the project.
	dataprocApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("dataproc.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		dataprocApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdDataprocApi, err := projects.NewService(ctx,
		"dpasp-dataproc.googleapis.com", dataprocApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable dataproc.googleapis.com api")
	}

	// Primary workers are the stable base (they carry HDFS DataNodes);
	// min_instances 0 accepts the API's default of 2.
	workerArgs := &dataproc.AutoscalingPolicyWorkerConfigArgs{
		MaxInstances: pulumi.Int(int(spec.WorkerConfig.MaxInstances)),
	}
	if spec.WorkerConfig.MinInstances > 0 {
		workerArgs.MinInstances = pulumi.IntPtr(int(spec.WorkerConfig.MinInstances))
	}
	if spec.WorkerConfig.Weight > 0 {
		workerArgs.Weight = pulumi.IntPtr(int(spec.WorkerConfig.Weight))
	}

	// The YARN memory-based algorithm. The scale factors express what
	// fraction of pending/available memory the autoscaler acts on per
	// cooldown period — 1.0 is maximally aggressive, small values smooth
	// scaling at the cost of reaction speed.
	// The scale factors use getter access: they are explicit-presence
	// fields precisely so 0.0 (a legitimate value — 0.0 scale-down
	// disables shrinking) survives validation and serialization.
	yarnArgs := &dataproc.AutoscalingPolicyBasicAlgorithmYarnConfigArgs{
		GracefulDecommissionTimeout: pulumi.String(spec.BasicAlgorithm.YarnConfig.GracefulDecommissionTimeout),
		ScaleUpFactor:               pulumi.Float64(spec.BasicAlgorithm.YarnConfig.GetScaleUpFactor()),
		ScaleDownFactor:             pulumi.Float64(spec.BasicAlgorithm.YarnConfig.GetScaleDownFactor()),
	}
	if spec.BasicAlgorithm.YarnConfig.ScaleUpMinWorkerFraction > 0 {
		yarnArgs.ScaleUpMinWorkerFraction = pulumi.Float64Ptr(spec.BasicAlgorithm.YarnConfig.ScaleUpMinWorkerFraction)
	}
	if spec.BasicAlgorithm.YarnConfig.ScaleDownMinWorkerFraction > 0 {
		yarnArgs.ScaleDownMinWorkerFraction = pulumi.Float64Ptr(spec.BasicAlgorithm.YarnConfig.ScaleDownMinWorkerFraction)
	}

	basicAlgorithmArgs := &dataproc.AutoscalingPolicyBasicAlgorithmArgs{
		YarnConfig: yarnArgs,
	}
	if spec.BasicAlgorithm.CooldownPeriod != "" {
		basicAlgorithmArgs.CooldownPeriod = pulumi.StringPtr(spec.BasicAlgorithm.CooldownPeriod)
	}

	args := &dataproc.AutoscalingPolicyArgs{
		PolicyId:       pulumi.String(spec.PolicyId),
		Location:       pulumi.StringPtr(spec.Location),
		WorkerConfig:   workerArgs,
		BasicAlgorithm: basicAlgorithmArgs,
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// DELETE (default) deletes the policy; ABANDON removes it from IaC
	// management; PREVENT fails destroying previews. The API's own
	// referenced-by-a-cluster guard applies on top.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	// Secondary workers are the cost-optimized burst arm — no HDFS data,
	// so the autoscaler can add and remove them aggressively.
	if spec.SecondaryWorkerConfig != nil {
		secondaryArgs := &dataproc.AutoscalingPolicySecondaryWorkerConfigArgs{}
		if spec.SecondaryWorkerConfig.MaxInstances > 0 {
			secondaryArgs.MaxInstances = pulumi.IntPtr(int(spec.SecondaryWorkerConfig.MaxInstances))
		}
		if spec.SecondaryWorkerConfig.MinInstances > 0 {
			secondaryArgs.MinInstances = pulumi.IntPtr(int(spec.SecondaryWorkerConfig.MinInstances))
		}
		if spec.SecondaryWorkerConfig.Weight > 0 {
			secondaryArgs.Weight = pulumi.IntPtr(int(spec.SecondaryWorkerConfig.Weight))
		}
		args.SecondaryWorkerConfig = secondaryArgs
	}

	createdPolicy, err := dataproc.NewAutoscalingPolicy(ctx, "autoscaling-policy", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdDataprocApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create dataproc autoscaling policy")
	}

	// The computed name is the fully qualified policy resource name —
	// the handle a cluster's autoscaling_policy_uri reference resolves to.
	ctx.Export(OpName, createdPolicy.Name)
	ctx.Export(OpPolicyId, createdPolicy.PolicyId)
	// The plain spec region name (not a provider-computed attribute), so
	// API callers and verifiers can address the policy without parsing
	// paths — identical contract to the Terraform module.
	ctx.Export(OpLocation, pulumi.String(spec.Location))

	return nil
}
