package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ghaRunnerScaleSetController deploys the GitHub Actions Runner Scale Set Controller
// using the official Helm chart.
func ghaRunnerScaleSetController(ctx *pulumi.Context, locals *Locals, k8sProvider *kubernetes.Provider) error {
	var dependencies []pulumi.Resource

	// Create namespace if requested
	if locals.CreateNamespace {
		ns, err := corev1.NewNamespace(ctx, locals.Namespace, &corev1.NamespaceArgs{
			Metadata: &metav1.ObjectMetaArgs{
				Name:   pulumi.String(locals.Namespace),
				Labels: pulumi.ToStringMap(locals.KubeLabels),
			},
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return errors.Wrap(err, "create namespace")
		}
		dependencies = append(dependencies, ns)
	}

	// Build Helm values
	helmValues := buildHelmValues(locals)

	// Deploy Helm chart
	// For OCI charts, the full URL must be passed as the Chart parameter
	// (RepositoryOpts.Repo doesn't work with OCI registries in Pulumi)
	_, err := helmv3.NewRelease(ctx, locals.ReleaseName, &helmv3.ReleaseArgs{
		Name:            pulumi.String(locals.ReleaseName),
		Namespace:       pulumi.String(locals.Namespace),
		CreateNamespace: pulumi.Bool(false), // We handle namespace creation ourselves
		Chart:           pulumi.String(vars.HelmChartOCI),
		Version:         pulumi.String(locals.ChartVersion),
		Values:          helmValues,
	}, pulumi.Provider(k8sProvider), pulumi.DependsOn(dependencies))
	if err != nil {
		return errors.Wrap(err, "deploy helm release")
	}

	return nil
}

// buildHelmValues constructs the Helm values map from the Locals struct.
func buildHelmValues(locals *Locals) pulumi.Map {
	values := pulumi.Map{
		"replicaCount": pulumi.Int(locals.ReplicaCount),
		"labels":       pulumi.ToStringMap(locals.KubeLabels),
	}

	// Container resources
	resources := pulumi.Map{}
	if locals.CpuRequests != "" || locals.MemoryRequests != "" {
		requests := pulumi.Map{}
		if locals.CpuRequests != "" {
			requests["cpu"] = pulumi.String(locals.CpuRequests)
		}
		if locals.MemoryRequests != "" {
			requests["memory"] = pulumi.String(locals.MemoryRequests)
		}
		resources["requests"] = requests
	}
	if locals.CpuLimits != "" || locals.MemoryLimits != "" {
		limits := pulumi.Map{}
		if locals.CpuLimits != "" {
			limits["cpu"] = pulumi.String(locals.CpuLimits)
		}
		if locals.MemoryLimits != "" {
			limits["memory"] = pulumi.String(locals.MemoryLimits)
		}
		resources["limits"] = limits
	}
	if len(resources) > 0 {
		values["resources"] = resources
	}

	// Custom image
	if locals.ImageRepository != "" || locals.ImageTag != "" || locals.ImagePullPolicy != "" {
		image := pulumi.Map{}
		if locals.ImageRepository != "" {
			image["repository"] = pulumi.String(locals.ImageRepository)
		}
		if locals.ImageTag != "" {
			image["tag"] = pulumi.String(locals.ImageTag)
		}
		if locals.ImagePullPolicy != "" {
			image["pullPolicy"] = pulumi.String(locals.ImagePullPolicy)
		}
		values["image"] = image
	}

	// Flags
	flags := pulumi.Map{}
	if locals.LogLevel != "" {
		flags["logLevel"] = pulumi.String(locals.LogLevel)
	}
	if locals.LogFormat != "" {
		flags["logFormat"] = pulumi.String(locals.LogFormat)
	}
	if locals.WatchSingleNamespace != "" {
		flags["watchSingleNamespace"] = pulumi.String(locals.WatchSingleNamespace)
	}
	if locals.RunnerMaxConcurrentReconciles > 0 {
		flags["runnerMaxConcurrentReconciles"] = pulumi.Int(locals.RunnerMaxConcurrentReconciles)
	}
	if locals.UpdateStrategy != "" {
		flags["updateStrategy"] = pulumi.String(locals.UpdateStrategy)
	}
	if len(locals.ExcludeLabelPropagationPrefixes) > 0 {
		flags["excludeLabelPropagationPrefixes"] = pulumi.ToStringArray(locals.ExcludeLabelPropagationPrefixes)
	}
	if locals.K8sClientRateLimiterQPS > 0 {
		flags["k8sClientRateLimiterQPS"] = pulumi.Int(locals.K8sClientRateLimiterQPS)
	}
	if locals.K8sClientRateLimiterBurst > 0 {
		flags["k8sClientRateLimiterBurst"] = pulumi.Int(locals.K8sClientRateLimiterBurst)
	}
	if len(flags) > 0 {
		values["flags"] = flags
	}

	// Metrics
	if locals.MetricsEnabled {
		values["metrics"] = pulumi.Map{
			"controllerManagerAddr": pulumi.String(locals.ControllerManagerAddr),
			"listenerAddr":          pulumi.String(locals.ListenerAddr),
			"listenerEndpoint":      pulumi.String(locals.ListenerEndpoint),
		}
	}

	// Image pull secrets
	if len(locals.ImagePullSecrets) > 0 {
		secrets := pulumi.Array{}
		for _, secret := range locals.ImagePullSecrets {
			secrets = append(secrets, pulumi.Map{"name": pulumi.String(secret)})
		}
		values["imagePullSecrets"] = secrets
	}

	// Priority class
	if locals.PriorityClassName != "" {
		values["priorityClassName"] = pulumi.String(locals.PriorityClassName)
	}

	return values
}
