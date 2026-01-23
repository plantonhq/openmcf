package module

import (
	"github.com/pkg/errors"
	kubernetesrookcephclusterv1 "github.com/plantonhq/project-planton/apis/org/project_planton/provider/kubernetes/kubernetesrookcephcluster/v1"
	"github.com/plantonhq/project-planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources creates all Pulumi resources for the Rook Ceph Cluster Kubernetes deployment.
func Resources(ctx *pulumi.Context, stackInput *kubernetesrookcephclusterv1.KubernetesRookCephClusterStackInput) error {
	// Initialize locals with computed values
	locals := initializeLocals(ctx, stackInput)

	// Set up kubernetes provider from the supplied cluster credential
	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx, stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to set up kubernetes provider")
	}

	// --------------------------------------------------------------------
	// 1. Namespace - conditionally create based on create_namespace flag
	// --------------------------------------------------------------------
	if stackInput.Target.Spec.GetCreateNamespace() {
		_, err := corev1.NewNamespace(ctx, locals.Namespace,
			&corev1.NamespaceArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Name:   pulumi.String(locals.Namespace),
					Labels: pulumi.ToStringMap(locals.Labels),
					// CRITICAL: Background Deletion Propagation Policy
					//
					// This annotation prevents namespace deletion from timing out during `pulumi destroy`.
					//
					// Problem: By default, Pulumi uses "Foreground" cascading deletion for namespaces.
					// Kubernetes adds a `foregroundDeletion` finalizer and waits for all resources inside
					// the namespace to be deleted before removing the namespace itself. However, if the
					// Helm release or CRDs are being deleted concurrently, there can be race conditions
					// where finalizers on child resources (like operator-managed CRs) prevent timely cleanup.
					//
					// Solution: Using "background" propagation policy causes Kubernetes to delete the
					// namespace object immediately. The namespace controller then asynchronously cleans up
					// all resources within the namespace. This avoids blocking on child resource finalizers.
					Annotations: pulumi.StringMap{
						"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
					},
				},
			},
			pulumi.Provider(kubernetesProvider))
		if err != nil {
			return errors.Wrap(err, "failed to create namespace")
		}
	}

	// --------------------------------------------------------------------
	// 2. Deploy the Rook Ceph Cluster via Helm
	// --------------------------------------------------------------------
	_, err = helm.NewRelease(ctx, locals.HelmReleaseName,
		&helm.ReleaseArgs{
			Name:            pulumi.String(locals.HelmReleaseName),
			Namespace:       pulumi.String(locals.Namespace),
			Chart:           pulumi.String(vars.HelmChartName),
			Version:         pulumi.String(locals.ChartVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(true),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(600), // Ceph cluster deployment can take longer
			Values:          pulumi.ToMap(locals.HelmValues),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.HelmChartRepo),
			},
		},
		pulumi.Provider(kubernetesProvider),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to install rook-ceph-cluster helm release")
	}

	return nil
}
