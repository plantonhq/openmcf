package module

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// kubernetesElasticOperator installs the Elastic Cloud on Kubernetes operator.
func kubernetesElasticOperator(ctx *pulumi.Context, locals *Locals,
	k8sProvider *pulumikubernetes.Provider) error {

	// --------------------------------------------------------------------
	// 1. Namespace - conditionally create based on create_namespace flag
	// --------------------------------------------------------------------
	if locals.KubernetesElasticOperator.Spec.CreateNamespace {
		_, err := corev1.NewNamespace(ctx, locals.Namespace, &corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(&metav1.ObjectMetaArgs{
				Name:   pulumi.String(locals.Namespace),
				Labels: pulumi.ToStringMap(locals.KubeLabels),
			}),
		}, pulumi.Provider(k8sProvider))
		if err != nil {
			return errors.Wrap(err, "create namespace")
		}
	}

	// --------------------------------------------------------------------
	// 2. Helm values – propagate Planton labels + optional resources.
	// --------------------------------------------------------------------
	values := pulumi.Map{
		"configKubernetes": pulumi.Map{
			"inherited_labels": pulumi.ToStringArray([]string{
				kuberneteslabelkeys.Resource,
				kuberneteslabelkeys.Organization,
				kuberneteslabelkeys.Environment,
				kuberneteslabelkeys.ResourceKind,
				kuberneteslabelkeys.ResourceId,
			}),
		},
	}

	if cr := locals.KubernetesElasticOperator.Spec.GetContainer().GetResources(); cr != nil {
		res := pulumi.Map{}
		if lim := cr.GetLimits(); lim != nil &&
			(lim.Cpu != "" || lim.Memory != "") {
			res["limits"] = pulumi.StringMap{
				"cpu":    pulumi.String(lim.Cpu),
				"memory": pulumi.String(lim.Memory),
			}
		}
		if req := cr.GetRequests(); req != nil &&
			(req.Cpu != "" || req.Memory != "") {
			res["requests"] = pulumi.StringMap{
				"cpu":    pulumi.String(req.Cpu),
				"memory": pulumi.String(req.Memory),
			}
		}
		if len(res) > 0 {
			values["resources"] = res
		}
	}

	// --------------------------------------------------------------------
	// 3. Helm release
	// --------------------------------------------------------------------
	helmReleaseOpts := []pulumi.ResourceOption{
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}),
		pulumi.Provider(k8sProvider),
	}

	// Use computed HelmReleaseName to avoid conflicts when multiple instances share a namespace
	_, err := helm.NewRelease(ctx, locals.HelmReleaseName, &helm.ReleaseArgs{
		Name:            pulumi.String(locals.HelmReleaseName),
		Namespace:       pulumi.String(locals.Namespace),
		Chart:           pulumi.String(vars.HelmChartName),
		Version:         pulumi.String(vars.HelmChartVersion),
		RepositoryOpts:  helm.RepositoryOptsArgs{Repo: pulumi.String(vars.HelmChartRepo)},
		CreateNamespace: pulumi.Bool(false),
		Atomic:          pulumi.Bool(false),
		CleanupOnFail:   pulumi.Bool(true),
		WaitForJobs:     pulumi.Bool(true),
		Timeout:         pulumi.Int(180),
		Values:          values,
	}, helmReleaseOpts...)
	if err != nil {
		return errors.Wrap(err, "install helm chart")
	}

	return nil
}
