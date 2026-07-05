package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/spanner"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// spannerInstance provisions the Cloud Spanner instance — the unit of
// compute/storage allocation that pins a geographic topology (config) and a
// capacity envelope shared by every database on it. name, config, and
// project are immutable; capacity, edition, and autoscaling all update in
// place (online, no downtime).
func spannerInstance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSpannerInstance.Spec

	// Enable the Spanner API so a fresh project can host instances.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("spanner.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"spanner-spanner.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable spanner.googleapis.com api")
	}

	args := &spanner.InstanceArgs{
		Name:        pulumi.StringPtr(locals.InstanceName),
		Config:      pulumi.String(spec.Config),
		DisplayName: pulumi.String(spec.DisplayName),
		Labels:      pulumi.ToStringMap(locals.GcpLabels),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omit the arg entirely).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Exactly one capacity method is enforced by spec validation. When none
	// is set on a PROVISIONED instance, the API defaults to 1 node.
	if spec.NumNodes > 0 {
		args.NumNodes = pulumi.IntPtr(int(spec.NumNodes))
	}
	if spec.ProcessingUnits > 0 {
		args.ProcessingUnits = pulumi.IntPtr(int(spec.ProcessingUnits))
	}
	if spec.AutoscalingConfig != nil {
		autoscalingArgs := &spanner.InstanceAutoscalingConfigArgs{}

		if spec.AutoscalingConfig.AutoscalingLimits != nil {
			limits := spec.AutoscalingConfig.AutoscalingLimits
			limitsArgs := &spanner.InstanceAutoscalingConfigAutoscalingLimitsArgs{}

			if limits.MinNodes > 0 {
				limitsArgs.MinNodes = pulumi.IntPtr(int(limits.MinNodes))
			}
			if limits.MaxNodes > 0 {
				limitsArgs.MaxNodes = pulumi.IntPtr(int(limits.MaxNodes))
			}
			if limits.MinProcessingUnits > 0 {
				limitsArgs.MinProcessingUnits = pulumi.IntPtr(int(limits.MinProcessingUnits))
			}
			if limits.MaxProcessingUnits > 0 {
				limitsArgs.MaxProcessingUnits = pulumi.IntPtr(int(limits.MaxProcessingUnits))
			}
			autoscalingArgs.AutoscalingLimits = limitsArgs
		}

		if spec.AutoscalingConfig.AutoscalingTargets != nil {
			targets := spec.AutoscalingConfig.AutoscalingTargets
			targetsArgs := &spanner.InstanceAutoscalingConfigAutoscalingTargetsArgs{}

			if targets.HighPriorityCpuUtilizationPercent > 0 {
				targetsArgs.HighPriorityCpuUtilizationPercent = pulumi.IntPtr(int(targets.HighPriorityCpuUtilizationPercent))
			}
			if targets.StorageUtilizationPercent > 0 {
				targetsArgs.StorageUtilizationPercent = pulumi.IntPtr(int(targets.StorageUtilizationPercent))
			}
			autoscalingArgs.AutoscalingTargets = targetsArgs
		}

		// Per-replica-location node bounds for multi-region instances: a
		// read-heavy region scales independently instead of sizing every
		// region for the hottest one. The spec flattens the provider's
		// single-field replica_selection wrapper to replica_location.
		if len(spec.AutoscalingConfig.AsymmetricAutoscalingOptions) > 0 {
			asymmetricOptions := make(spanner.InstanceAutoscalingConfigAsymmetricAutoscalingOptionArray, 0,
				len(spec.AutoscalingConfig.AsymmetricAutoscalingOptions))
			for _, option := range spec.AutoscalingConfig.AsymmetricAutoscalingOptions {
				asymmetricOptions = append(asymmetricOptions,
					&spanner.InstanceAutoscalingConfigAsymmetricAutoscalingOptionArgs{
						ReplicaSelection: &spanner.InstanceAutoscalingConfigAsymmetricAutoscalingOptionReplicaSelectionArgs{
							Location: pulumi.String(option.ReplicaLocation),
						},
						Overrides: &spanner.InstanceAutoscalingConfigAsymmetricAutoscalingOptionOverridesArgs{
							AutoscalingLimits: &spanner.InstanceAutoscalingConfigAsymmetricAutoscalingOptionOverridesAutoscalingLimitsArgs{
								MinNodes: pulumi.Int(int(option.Overrides.MinNodes)),
								MaxNodes: pulumi.Int(int(option.Overrides.MaxNodes)),
							},
						},
					})
			}
			autoscalingArgs.AsymmetricAutoscalingOptions = asymmetricOptions
		}

		args.AutoscalingConfig = autoscalingArgs
	}

	if spec.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(spec.InstanceType)
	}
	if spec.Edition != "" {
		args.Edition = pulumi.StringPtr(spec.Edition)
	}
	if spec.DefaultBackupScheduleType != "" {
		args.DefaultBackupScheduleType = pulumi.StringPtr(spec.DefaultBackupScheduleType)
	}

	// When true, destroy deletes all backups held on the instance first;
	// when false (the safe default), destroy fails while any backup exists.
	if spec.ForceDestroy {
		args.ForceDestroy = pulumi.BoolPtr(true)
	}

	createdInstance, err := spanner.NewInstance(ctx, "spanner-instance", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create spanner instance")
	}

	// instance_id is built from the created resource's resolved project so
	// the output is correct under the ambient-project fallback (the spec
	// project may be empty).
	ctx.Export(OpInstanceId, pulumi.Sprintf(
		"projects/%s/instances/%s",
		createdInstance.Project,
		createdInstance.Name,
	))
	ctx.Export(OpInstanceName, createdInstance.Name)
	ctx.Export(OpState, createdInstance.State)

	// The bridged provider returns config as the fully qualified
	// projects/{p}/instanceConfigs/{name} path while the Terraform provider
	// stores the plain name the spec passed in. The output contract is the
	// plain config name (what spec authors and API callers use), so strip
	// the path prefix for cross-engine parity.
	ctx.Export(OpConfig, createdInstance.Config.ApplyT(func(config string) string {
		if idx := strings.LastIndex(config, "/"); idx >= 0 {
			return config[idx+1:]
		}
		return config
	}).(pulumi.StringOutput))

	return nil
}
