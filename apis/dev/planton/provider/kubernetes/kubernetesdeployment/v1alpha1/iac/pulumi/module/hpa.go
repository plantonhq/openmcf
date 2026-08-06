package module

import (
	"fmt"

	"github.com/pkg/errors"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	autoscalingv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/autoscaling/v2"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// horizontalPodAutoscaler creates an autoscaling/v2 HPA targeting the Deployment
// when autoscaling is enabled. CPU targets are expressed as average Utilization
// (percentage of requests); memory targets as an absolute AverageValue per pod —
// matching how each metric is meaningfully compared across replicas.
func horizontalPodAutoscaler(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource, createdDeployment *appsv1.Deployment) error {

	availability := locals.KubernetesDeployment.Spec.Availability
	if availability == nil || availability.HorizontalPodAutoscaling == nil ||
		!availability.HorizontalPodAutoscaling.Enabled {
		return nil
	}
	hpaSpec := availability.HorizontalPodAutoscaling

	minReplicas := int32(1)
	if availability.Replicas != nil {
		minReplicas = *availability.Replicas
	}

	metrics := autoscalingv2.MetricSpecArray{}

	if hpaSpec.TargetCpuUtilizationPercent != nil {
		metrics = append(metrics, &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("Resource"),
			Resource: &autoscalingv2.ResourceMetricSourceArgs{
				Name: pulumi.String("cpu"),
				Target: &autoscalingv2.MetricTargetArgs{
					Type:               pulumi.String("Utilization"),
					AverageUtilization: pulumi.Int(int(*hpaSpec.TargetCpuUtilizationPercent)),
				},
			},
		})
	}

	if hpaSpec.TargetMemoryUtilization != "" {
		metrics = append(metrics, &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("Resource"),
			Resource: &autoscalingv2.ResourceMetricSourceArgs{
				Name: pulumi.String("memory"),
				Target: &autoscalingv2.MetricTargetArgs{
					Type:         pulumi.String("AverageValue"),
					AverageValue: pulumi.String(hpaSpec.TargetMemoryUtilization),
				},
			},
		})
	}

	hpaArgs := &autoscalingv2.HorizontalPodAutoscalerArgs{
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String(locals.KubernetesDeployment.Metadata.Name),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Spec: &autoscalingv2.HorizontalPodAutoscalerSpecArgs{
			ScaleTargetRef: &autoscalingv2.CrossVersionObjectReferenceArgs{
				ApiVersion: pulumi.String("apps/v1"),
				Kind:       pulumi.String("Deployment"),
				Name:       pulumi.String(locals.KubernetesDeployment.Metadata.Name),
			},
			MinReplicas: pulumi.Int(int(minReplicas)),
			MaxReplicas: pulumi.Int(int(hpaSpec.MaxReplicas)),
			Metrics:     metrics,
		},
	}

	hpaResourceName := fmt.Sprintf("%s-hpa", locals.KubernetesDeployment.Metadata.Name)
	_, err := autoscalingv2.NewHorizontalPodAutoscaler(ctx,
		hpaResourceName,
		hpaArgs,
		pulumi.Provider(kubernetesProvider),
		pulumi.DependsOn([]pulumi.Resource{createdDeployment}))
	if err != nil {
		return errors.Wrap(err, "failed to create horizontal pod autoscaler")
	}

	return nil
}
