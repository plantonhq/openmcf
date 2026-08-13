package module

import (
	"github.com/pkg/errors"
	gcpcomputemigv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcomputemig/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// autoscaler creates the group's autoscaler — zonal or regional per the
// spec's location selector — targeted at the created group manager by
// reference (the provider expects the manager's URL; wiring the resource
// output keeps dependency order and the ambient-project case correct on
// both engines).
func autoscaler(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
	groupManager *igmResult,
) error {

	spec := locals.GcpComputeMig.Spec
	autoscalerSpec := spec.Autoscaler

	if locals.IsRegional {
		args := &compute.RegionAutoscalerArgs{
			Name:              pulumi.StringPtr(locals.AutoscalerName),
			Region:            pulumi.StringPtr(spec.Region),
			Target:            groupManager.SelfLink,
			AutoscalingPolicy: regionalAutoscalingPolicy(autoscalerSpec),
		}
		if spec.ProjectId.GetValue() != "" {
			args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}
		if autoscalerSpec.Description != "" {
			args.Description = pulumi.StringPtr(autoscalerSpec.Description)
		}
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}
		if _, err := compute.NewRegionAutoscaler(ctx,
			locals.AutoscalerName,
			args,
			pulumi.Provider(gcpProvider),
		); err != nil {
			return errors.Wrap(err, "failed to create regional autoscaler")
		}
		return nil
	}

	args := &compute.AutoscalerArgs{
		Name:              pulumi.StringPtr(locals.AutoscalerName),
		Zone:              pulumi.StringPtr(spec.Zone),
		Target:            groupManager.SelfLink,
		AutoscalingPolicy: zonalAutoscalingPolicy(autoscalerSpec),
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if autoscalerSpec.Description != "" {
		args.Description = pulumi.StringPtr(autoscalerSpec.Description)
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}
	if _, err := compute.NewAutoscaler(ctx,
		locals.AutoscalerName,
		args,
		pulumi.Provider(gcpProvider),
	); err != nil {
		return errors.Wrap(err, "failed to create autoscaler")
	}
	return nil
}

// zonalAutoscalingPolicy maps the spec's autoscaler surface to the zonal
// autoscaler's policy type.
func zonalAutoscalingPolicy(autoscalerSpec *gcpcomputemigv1alpha1.GcpComputeMigAutoscaler) *compute.AutoscalerAutoscalingPolicyArgs {
	policy := &compute.AutoscalerAutoscalingPolicyArgs{
		MinReplicas: pulumi.Int(int(autoscalerSpec.MinReplicas)),
		MaxReplicas: pulumi.Int(int(autoscalerSpec.MaxReplicas)),
	}
	if autoscalerSpec.CooldownPeriod != nil {
		policy.CooldownPeriod = pulumi.IntPtr(int(autoscalerSpec.GetCooldownPeriod()))
	}
	if autoscalerSpec.Mode != "" {
		policy.Mode = pulumi.StringPtr(autoscalerSpec.Mode)
	}
	if autoscalerSpec.CpuTarget != nil || autoscalerSpec.CpuPredictiveMethod != "" {
		cpu := &compute.AutoscalerAutoscalingPolicyCpuUtilizationArgs{}
		if autoscalerSpec.CpuTarget != nil {
			cpu.Target = pulumi.Float64(autoscalerSpec.GetCpuTarget())
		}
		if autoscalerSpec.CpuPredictiveMethod != "" {
			cpu.PredictiveMethod = pulumi.StringPtr(autoscalerSpec.CpuPredictiveMethod)
		}
		policy.CpuUtilization = cpu
	}
	if autoscalerSpec.LoadBalancingTarget != nil {
		policy.LoadBalancingUtilization = &compute.AutoscalerAutoscalingPolicyLoadBalancingUtilizationArgs{
			Target: pulumi.Float64(autoscalerSpec.GetLoadBalancingTarget()),
		}
	}
	if len(autoscalerSpec.Metrics) > 0 {
		metrics := compute.AutoscalerAutoscalingPolicyMetricArray{}
		for _, metric := range autoscalerSpec.Metrics {
			metricArgs := &compute.AutoscalerAutoscalingPolicyMetricArgs{
				Name: pulumi.String(metric.Name),
			}
			if metric.Target != nil {
				metricArgs.Target = pulumi.Float64Ptr(metric.GetTarget())
			}
			if metric.Type != "" {
				metricArgs.Type = pulumi.StringPtr(metric.Type)
			}
			if metric.Filter != "" {
				metricArgs.Filter = pulumi.StringPtr(metric.Filter)
			}
			if metric.SingleInstanceAssignment != nil {
				metricArgs.SingleInstanceAssignment = pulumi.Float64Ptr(metric.GetSingleInstanceAssignment())
			}
			metrics = append(metrics, metricArgs)
		}
		policy.Metrics = metrics
	}
	if scaleIn := autoscalerSpec.ScaleInControl; scaleIn != nil {
		scaleInArgs := &compute.AutoscalerAutoscalingPolicyScaleInControlArgs{}
		if scaleIn.MaxScaledInReplicasFixed != nil || scaleIn.MaxScaledInReplicasPercent != nil {
			replicas := &compute.AutoscalerAutoscalingPolicyScaleInControlMaxScaledInReplicasArgs{}
			if scaleIn.MaxScaledInReplicasFixed != nil {
				replicas.Fixed = pulumi.IntPtr(int(scaleIn.GetMaxScaledInReplicasFixed()))
			}
			if scaleIn.MaxScaledInReplicasPercent != nil {
				replicas.Percent = pulumi.IntPtr(int(scaleIn.GetMaxScaledInReplicasPercent()))
			}
			scaleInArgs.MaxScaledInReplicas = replicas
		}
		if scaleIn.TimeWindowSec != nil {
			scaleInArgs.TimeWindowSec = pulumi.IntPtr(int(scaleIn.GetTimeWindowSec()))
		}
		policy.ScaleInControl = scaleInArgs
	}
	if len(autoscalerSpec.Schedules) > 0 {
		schedules := compute.AutoscalerAutoscalingPolicyScalingScheduleArray{}
		for _, schedule := range autoscalerSpec.Schedules {
			scheduleArgs := &compute.AutoscalerAutoscalingPolicyScalingScheduleArgs{
				Name:                pulumi.String(schedule.ScheduleName),
				Schedule:            pulumi.String(schedule.Schedule),
				DurationSec:         pulumi.Int(int(schedule.DurationSec)),
				MinRequiredReplicas: pulumi.Int(int(schedule.MinRequiredReplicas)),
			}
			if schedule.Disabled != nil {
				scheduleArgs.Disabled = pulumi.BoolPtr(schedule.GetDisabled())
			}
			if schedule.TimeZone != "" {
				scheduleArgs.TimeZone = pulumi.StringPtr(schedule.TimeZone)
			}
			if schedule.Description != "" {
				scheduleArgs.Description = pulumi.StringPtr(schedule.Description)
			}
			schedules = append(schedules, scheduleArgs)
		}
		policy.ScalingSchedules = schedules
	}
	if autoscalerSpec.StabilizationPeriod != nil {
		policy.StabilizationPeriod = pulumi.IntPtr(int(autoscalerSpec.GetStabilizationPeriod()))
	}
	return policy
}

// regionalAutoscalingPolicy mirrors the zonal builder for the regional
// autoscaler's types.
func regionalAutoscalingPolicy(autoscalerSpec *gcpcomputemigv1alpha1.GcpComputeMigAutoscaler) *compute.RegionAutoscalerAutoscalingPolicyArgs {
	policy := &compute.RegionAutoscalerAutoscalingPolicyArgs{
		MinReplicas: pulumi.Int(int(autoscalerSpec.MinReplicas)),
		MaxReplicas: pulumi.Int(int(autoscalerSpec.MaxReplicas)),
	}
	if autoscalerSpec.CooldownPeriod != nil {
		policy.CooldownPeriod = pulumi.IntPtr(int(autoscalerSpec.GetCooldownPeriod()))
	}
	if autoscalerSpec.Mode != "" {
		policy.Mode = pulumi.StringPtr(autoscalerSpec.Mode)
	}
	if autoscalerSpec.CpuTarget != nil || autoscalerSpec.CpuPredictiveMethod != "" {
		cpu := &compute.RegionAutoscalerAutoscalingPolicyCpuUtilizationArgs{}
		if autoscalerSpec.CpuTarget != nil {
			cpu.Target = pulumi.Float64(autoscalerSpec.GetCpuTarget())
		}
		if autoscalerSpec.CpuPredictiveMethod != "" {
			cpu.PredictiveMethod = pulumi.StringPtr(autoscalerSpec.CpuPredictiveMethod)
		}
		policy.CpuUtilization = cpu
	}
	if autoscalerSpec.LoadBalancingTarget != nil {
		policy.LoadBalancingUtilization = &compute.RegionAutoscalerAutoscalingPolicyLoadBalancingUtilizationArgs{
			Target: pulumi.Float64(autoscalerSpec.GetLoadBalancingTarget()),
		}
	}
	if len(autoscalerSpec.Metrics) > 0 {
		metrics := compute.RegionAutoscalerAutoscalingPolicyMetricArray{}
		for _, metric := range autoscalerSpec.Metrics {
			metricArgs := &compute.RegionAutoscalerAutoscalingPolicyMetricArgs{
				Name: pulumi.String(metric.Name),
			}
			if metric.Target != nil {
				metricArgs.Target = pulumi.Float64Ptr(metric.GetTarget())
			}
			if metric.Type != "" {
				metricArgs.Type = pulumi.StringPtr(metric.Type)
			}
			if metric.Filter != "" {
				metricArgs.Filter = pulumi.StringPtr(metric.Filter)
			}
			if metric.SingleInstanceAssignment != nil {
				metricArgs.SingleInstanceAssignment = pulumi.Float64Ptr(metric.GetSingleInstanceAssignment())
			}
			metrics = append(metrics, metricArgs)
		}
		policy.Metrics = metrics
	}
	if scaleIn := autoscalerSpec.ScaleInControl; scaleIn != nil {
		scaleInArgs := &compute.RegionAutoscalerAutoscalingPolicyScaleInControlArgs{}
		if scaleIn.MaxScaledInReplicasFixed != nil || scaleIn.MaxScaledInReplicasPercent != nil {
			replicas := &compute.RegionAutoscalerAutoscalingPolicyScaleInControlMaxScaledInReplicasArgs{}
			if scaleIn.MaxScaledInReplicasFixed != nil {
				replicas.Fixed = pulumi.IntPtr(int(scaleIn.GetMaxScaledInReplicasFixed()))
			}
			if scaleIn.MaxScaledInReplicasPercent != nil {
				replicas.Percent = pulumi.IntPtr(int(scaleIn.GetMaxScaledInReplicasPercent()))
			}
			scaleInArgs.MaxScaledInReplicas = replicas
		}
		if scaleIn.TimeWindowSec != nil {
			scaleInArgs.TimeWindowSec = pulumi.IntPtr(int(scaleIn.GetTimeWindowSec()))
		}
		policy.ScaleInControl = scaleInArgs
	}
	if len(autoscalerSpec.Schedules) > 0 {
		schedules := compute.RegionAutoscalerAutoscalingPolicyScalingScheduleArray{}
		for _, schedule := range autoscalerSpec.Schedules {
			scheduleArgs := &compute.RegionAutoscalerAutoscalingPolicyScalingScheduleArgs{
				Name:                pulumi.String(schedule.ScheduleName),
				Schedule:            pulumi.String(schedule.Schedule),
				DurationSec:         pulumi.Int(int(schedule.DurationSec)),
				MinRequiredReplicas: pulumi.Int(int(schedule.MinRequiredReplicas)),
			}
			if schedule.Disabled != nil {
				scheduleArgs.Disabled = pulumi.BoolPtr(schedule.GetDisabled())
			}
			if schedule.TimeZone != "" {
				scheduleArgs.TimeZone = pulumi.StringPtr(schedule.TimeZone)
			}
			if schedule.Description != "" {
				scheduleArgs.Description = pulumi.StringPtr(schedule.Description)
			}
			schedules = append(schedules, scheduleArgs)
		}
		policy.ScalingSchedules = schedules
	}
	if autoscalerSpec.StabilizationPeriod != nil {
		policy.StabilizationPeriod = pulumi.IntPtr(int(autoscalerSpec.GetStabilizationPeriod()))
	}
	return policy
}
