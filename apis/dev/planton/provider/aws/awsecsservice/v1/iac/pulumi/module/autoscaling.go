package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/appautoscaling"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecs"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// AWS defaults for target-tracking cooldowns, applied when the spec's
// optional cooldown fields are unset.
const (
	defaultScaleInCooldownSeconds  = 300
	defaultScaleOutCooldownSeconds = 60
)

// autoscaling registers the service as an Application Auto Scaling target
// and attaches one target-tracking policy per configured metric. The
// scaler's identity is "service/<cluster>/<service>" -- it exists only for
// this service, which is why scaling is folded into this spec rather than
// being its own kind.
func autoscaling(ctx *pulumi.Context, locals *Locals, provider *aws.Provider, createdService *ecs.Service) error {
	spec := locals.AwsEcsService.Spec
	serviceName := locals.AwsEcsService.Metadata.Name
	scaling := spec.Autoscaling

	// The scalable target's resource id wants the cluster NAME; the spec
	// carries the cluster ARN (arn:aws:ecs:region:account:cluster/<name>).
	clusterArn := spec.ClusterArn.GetValue()
	clusterName := clusterArn
	if idx := strings.LastIndex(clusterArn, "/"); idx != -1 {
		clusterName = clusterArn[idx+1:]
	}

	target, err := appautoscaling.NewTarget(ctx,
		"autoscaling-target",
		&appautoscaling.TargetArgs{
			MaxCapacity:       pulumi.Int(int(scaling.MaxTasks)),
			MinCapacity:       pulumi.Int(int(scaling.MinTasks)),
			ResourceId:        pulumi.Sprintf("service/%s/%s", clusterName, serviceName),
			ScalableDimension: pulumi.String("ecs:service:DesiredCount"),
			ServiceNamespace:  pulumi.String("ecs"),
		},
		pulumi.Provider(provider),
		pulumi.DependsOn([]pulumi.Resource{createdService}))
	if err != nil {
		return errors.Wrap(err, "failed to create autoscaling target")
	}

	if scaling.Cpu != nil {
		if err := trackingPolicy(ctx, provider, target, "cpu-scaling-policy",
			fmt.Sprintf("%s-cpu-scaling", serviceName),
			"ECSServiceAverageCPUUtilization", "",
			float64(scaling.Cpu.TargetPercent),
			scaling.Cpu.ScaleInCooldownSeconds, scaling.Cpu.ScaleOutCooldownSeconds,
			scaling.Cpu.DisableScaleIn); err != nil {
			return err
		}
	}

	if scaling.Memory != nil {
		if err := trackingPolicy(ctx, provider, target, "memory-scaling-policy",
			fmt.Sprintf("%s-memory-scaling", serviceName),
			"ECSServiceAverageMemoryUtilization", "",
			float64(scaling.Memory.TargetPercent),
			scaling.Memory.ScaleInCooldownSeconds, scaling.Memory.ScaleOutCooldownSeconds,
			scaling.Memory.DisableScaleIn); err != nil {
			return err
		}
	}

	if scaling.RequestsPerTarget != nil {
		requests := scaling.RequestsPerTarget
		// ALBRequestCountPerTarget is scoped by the resource label
		// "<lb-arn-suffix>/<tg-arn-suffix>" -- both halves come from the
		// referenced load balancer's and target group's arn_suffix outputs.
		resourceLabel := fmt.Sprintf("%s/%s",
			requests.LoadBalancerArnSuffix.GetValue(),
			requests.TargetGroupArnSuffix.GetValue())
		if err := trackingPolicy(ctx, provider, target, "requests-scaling-policy",
			fmt.Sprintf("%s-requests-scaling", serviceName),
			"ALBRequestCountPerTarget", resourceLabel,
			requests.TargetRequestsPerTarget,
			requests.ScaleInCooldownSeconds, requests.ScaleOutCooldownSeconds,
			requests.DisableScaleIn); err != nil {
			return err
		}
	}

	return nil
}

// trackingPolicy creates one target-tracking policy against the service's
// scalable target.
func trackingPolicy(
	ctx *pulumi.Context,
	provider *aws.Provider,
	target *appautoscaling.Target,
	resourceName string,
	policyName string,
	predefinedMetricType string,
	resourceLabel string,
	targetValue float64,
	scaleInCooldown *int32,
	scaleOutCooldown *int32,
	disableScaleIn bool,
) error {
	metricSpecification := &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationPredefinedMetricSpecificationArgs{
		PredefinedMetricType: pulumi.String(predefinedMetricType),
	}
	if resourceLabel != "" {
		metricSpecification.ResourceLabel = pulumi.StringPtr(resourceLabel)
	}

	scaleIn := int32(defaultScaleInCooldownSeconds)
	if scaleInCooldown != nil {
		scaleIn = *scaleInCooldown
	}
	scaleOut := int32(defaultScaleOutCooldownSeconds)
	if scaleOutCooldown != nil {
		scaleOut = *scaleOutCooldown
	}

	configuration := &appautoscaling.PolicyTargetTrackingScalingPolicyConfigurationArgs{
		TargetValue:                   pulumi.Float64(targetValue),
		PredefinedMetricSpecification: metricSpecification,
		ScaleInCooldown:               pulumi.IntPtr(int(scaleIn)),
		ScaleOutCooldown:              pulumi.IntPtr(int(scaleOut)),
	}
	if disableScaleIn {
		configuration.DisableScaleIn = pulumi.BoolPtr(true)
	}

	_, err := appautoscaling.NewPolicy(ctx,
		resourceName,
		&appautoscaling.PolicyArgs{
			Name:                                     pulumi.String(policyName),
			PolicyType:                               pulumi.String("TargetTrackingScaling"),
			ResourceId:                               target.ResourceId,
			ScalableDimension:                        target.ScalableDimension,
			ServiceNamespace:                         target.ServiceNamespace,
			TargetTrackingScalingPolicyConfiguration: configuration,
		},
		pulumi.Provider(provider),
		pulumi.DependsOn([]pulumi.Resource{target}))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s policy", predefinedMetricType)
	}
	return nil
}
