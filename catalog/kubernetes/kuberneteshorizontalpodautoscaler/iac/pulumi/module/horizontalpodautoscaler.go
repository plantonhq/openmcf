package module

import (
	"github.com/pkg/errors"
	kuberneteshorizontalpodautoscalerv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kuberneteshorizontalpodautoscaler/v1alpha1"
	autoscalingv2 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/autoscaling/v2"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createHorizontalPodAutoscaler creates the autoscaling/v2 HPA.
//
// minReplicas is ALWAYS sent explicitly (Kubernetes default 1 applied
// module-side) so both engines submit identical objects. When the spec lists
// no metrics, the metrics field is OMITTED — the API server then applies its
// own default (80% average CPU utilization), and sending an empty list would
// instead disable metric-driven scaling.
func createHorizontalPodAutoscaler(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) (*autoscalingv2.HorizontalPodAutoscaler, error) {
	spec := locals.Spec

	hpaSpecArgs := &autoscalingv2.HorizontalPodAutoscalerSpecArgs{
		ScaleTargetRef: &autoscalingv2.CrossVersionObjectReferenceArgs{
			ApiVersion: pulumi.String(locals.TargetApiVersion),
			Kind:       pulumi.String(locals.TargetKind),
			Name:       pulumi.String(locals.TargetName),
		},
		MinReplicas: pulumi.Int(int(locals.MinReplicas)),
		MaxReplicas: pulumi.Int(int(spec.GetMaxReplicas())),
	}

	if len(spec.GetMetrics()) > 0 {
		metricArray := autoscalingv2.MetricSpecArray{}
		for _, m := range spec.GetMetrics() {
			metricArray = append(metricArray, buildMetric(m))
		}
		hpaSpecArgs.Metrics = metricArray
	}

	if behavior := spec.GetBehavior(); behavior != nil {
		behaviorArgs := &autoscalingv2.HorizontalPodAutoscalerBehaviorArgs{}
		if up := behavior.GetScaleUp(); up != nil {
			behaviorArgs.ScaleUp = buildScalingRules(up)
		}
		if down := behavior.GetScaleDown(); down != nil {
			behaviorArgs.ScaleDown = buildScalingRules(down)
		}
		hpaSpecArgs.Behavior = behaviorArgs
	}

	horizontalPodAutoscaler, err := autoscalingv2.NewHorizontalPodAutoscaler(
		ctx,
		locals.Name,
		&autoscalingv2.HorizontalPodAutoscalerArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:        pulumi.String(locals.Name),
				Namespace:   pulumi.String(locals.Namespace),
				Labels:      pulumi.ToStringMap(locals.Labels),
				Annotations: pulumi.ToStringMap(locals.Annotations),
			},
			Spec: hpaSpecArgs,
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create horizontal pod autoscaler %s/%s", locals.Namespace, locals.Name)
	}

	return horizontalPodAutoscaler, nil
}

// buildMetric converts one proto metric onto the API MetricSpec. The spec's
// CEL rules guarantee exactly the source matching the type is present.
func buildMetric(m *kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetric) *autoscalingv2.MetricSpecArgs {
	switch m.GetType() {
	case kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricType_container_resource:
		cr := m.GetContainerResource()
		return &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("ContainerResource"),
			ContainerResource: &autoscalingv2.ContainerResourceMetricSourceArgs{
				Name:      pulumi.String(cr.GetName()),
				Container: pulumi.String(cr.GetContainer()),
				Target:    buildMetricTarget(cr.GetTarget()),
			},
		}
	case kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricType_pods:
		p := m.GetPods()
		return &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("Pods"),
			Pods: &autoscalingv2.PodsMetricSourceArgs{
				Metric: buildMetricIdentifier(p.GetMetric()),
				Target: buildMetricTarget(p.GetTarget()),
			},
		}
	case kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricType_object:
		o := m.GetObject()
		return &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("Object"),
			Object: &autoscalingv2.ObjectMetricSourceArgs{
				DescribedObject: &autoscalingv2.CrossVersionObjectReferenceArgs{
					ApiVersion: pulumi.String(o.GetDescribedObject().GetApiVersion()),
					Kind:       pulumi.String(o.GetDescribedObject().GetKind()),
					Name:       pulumi.String(o.GetDescribedObject().GetName()),
				},
				Metric: buildMetricIdentifier(o.GetMetric()),
				Target: buildMetricTarget(o.GetTarget()),
			},
		}
	case kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricType_external:
		e := m.GetExternal()
		return &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("External"),
			External: &autoscalingv2.ExternalMetricSourceArgs{
				Metric: buildMetricIdentifier(e.GetMetric()),
				Target: buildMetricTarget(e.GetTarget()),
			},
		}
	default: // resource
		r := m.GetResource()
		return &autoscalingv2.MetricSpecArgs{
			Type: pulumi.String("Resource"),
			Resource: &autoscalingv2.ResourceMetricSourceArgs{
				Name:   pulumi.String(r.GetName()),
				Target: buildMetricTarget(r.GetTarget()),
			},
		}
	}
}

// buildMetricTarget converts the target with exactly the value form the
// spec's CEL rules guarantee is present for the declared type.
func buildMetricTarget(t *kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricTarget) *autoscalingv2.MetricTargetArgs {
	args := &autoscalingv2.MetricTargetArgs{
		Type: pulumi.String(metricTargetTypeApiString(t.GetType())),
	}
	if t.AverageUtilization != nil {
		args.AverageUtilization = pulumi.Int(int(t.GetAverageUtilization()))
	}
	if t.GetValue() != "" {
		args.Value = pulumi.String(t.GetValue())
	}
	if t.GetAverageValue() != "" {
		args.AverageValue = pulumi.String(t.GetAverageValue())
	}
	return args
}

// buildMetricIdentifier converts the metric identity; match_labels scope the
// series the metrics adapter reads.
func buildMetricIdentifier(id *kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerMetricIdentifier) *autoscalingv2.MetricIdentifierArgs {
	args := &autoscalingv2.MetricIdentifierArgs{
		Name: pulumi.String(id.GetName()),
	}
	if len(id.GetMatchLabels()) > 0 {
		args.Selector = &metav1.LabelSelectorArgs{
			MatchLabels: pulumi.ToStringMap(id.GetMatchLabels()),
		}
	}
	return args
}

// buildScalingRules converts one direction's tuning. selectPolicy is always
// sent (with the API default Max applied module-side) so both engines submit
// identical behavior blocks.
func buildScalingRules(rules *kuberneteshorizontalpodautoscalerv1alpha1.KubernetesHorizontalPodAutoscalerScalingRules) *autoscalingv2.HPAScalingRulesArgs {
	args := &autoscalingv2.HPAScalingRulesArgs{
		SelectPolicy: pulumi.String(selectPolicyApiString(rules)),
	}
	if rules.StabilizationWindowSeconds != nil {
		args.StabilizationWindowSeconds = pulumi.Int(int(rules.GetStabilizationWindowSeconds()))
	}
	if len(rules.GetPolicies()) > 0 {
		policyArray := autoscalingv2.HPAScalingPolicyArray{}
		for _, p := range rules.GetPolicies() {
			policyArray = append(policyArray, &autoscalingv2.HPAScalingPolicyArgs{
				Type:          pulumi.String(scalingPolicyTypeApiString(p.GetType())),
				Value:         pulumi.Int(int(p.GetValue())),
				PeriodSeconds: pulumi.Int(int(p.GetPeriodSeconds())),
			})
		}
		args.Policies = policyArray
	}
	return args
}
