package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// metricsConfigMap renders the module-owned JMX Prometheus Exporter rules
// ConfigMap when spec.metrics.enabled — the canonical Strimzi rule set
// (metrics_rules.go) under the key the Kafka CR's metricsConfig points at.
// Returns nil when metrics are disabled. Terraform equivalent:
// kubernetes_config_map_v1 with count.
func metricsConfigMap(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.Spec.GetMetrics().GetEnabled() {
		return nil, nil
	}

	createdConfigMap, err := kubernetescorev1.NewConfigMap(ctx, locals.MetricsConfigMapName,
		&kubernetescorev1.ConfigMapArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(
				&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(locals.MetricsConfigMapName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
			Data: pulumi.StringMap{
				vars.MetricsConfigKey: pulumi.String(kafkaMetricsRules),
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kafka metrics ConfigMap")
	}

	return createdConfigMap, nil
}
