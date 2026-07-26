package module

import (
	"github.com/pkg/errors"
	kuberneteskyvernov1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskyverno/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the Kyverno policy engine from the official chart as
// a real Helm release. The typed spec renders into chart values
// (values.go); the helm_values escape hatch merges last with Helm -f
// semantics — the exact semantic twin of the Terraform module's
// helm_release with values = [typed, helm_values].
//
// WEBHOOK LIFECYCLE: the chart templates NO webhook configurations — the
// admission controller REGISTERS them at runtime and the chart's
// pre-delete cleanup hook removes them at uninstall (webhooksCleanup,
// rendered explicitly). A release force-deleted without the hook strands
// kyverno-* webhook configurations whose backing service is gone; the
// spec's top-level comment carries the manual unstick command.
//
// CRD LIFECYCLE: the policy CRDs are chart-TEMPLATED via the crds
// subchart — installed and DELETED with the release unless
// crds.keep_on_uninstall injects the resource-policy annotation.
// Destroying the engine cascade-deletes every policy on the cluster in
// the default posture (the spec's CRD field comments carry the warning).
func Resources(ctx *pulumi.Context, stackInput *kuberneteskyvernov1.KubernetesKyvernoStackInput) error {
	locals, err := initializeLocals(ctx, stackInput)
	if err != nil {
		return errors.Wrap(err, "failed to initialize locals")
	}

	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	// ------------------------------ namespace ----------------------------
	createdNamespace, err := namespace(ctx, stackInput, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	var releaseDeps []pulumi.ResourceOption
	if createdNamespace != nil {
		releaseDeps = append(releaseDeps, pulumi.DependsOn([]pulumi.Resource{createdNamespace}))
	}

	// ------------------------------ helm release --------------------------
	mergedValues, err := buildHelmValues(locals)
	if err != nil {
		return errors.Wrap(err, "failed to build helm values")
	}

	releaseArgs := &helmv3.ReleaseArgs{
		Name:      pulumi.String(locals.ReleaseName),
		Namespace: pulumi.String(locals.Namespace),
		Chart:     pulumi.String(vars.HelmChartName),
		Version:   pulumi.String(locals.ChartVersion),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String(vars.HelmChartRepo),
		},
		Values: pulumi.ToMap(mergedValues),
		// The module owns namespace creation (create_namespace flag).
		CreateNamespace: pulumi.Bool(false),
		// Wait for the controllers to become Ready — an admission engine
		// that never starts should fail THIS deploy, not the first policy
		// apply. SkipAwait false is Helm --wait, stated explicitly to
		// mirror the Terraform twin's `wait = true`. The timeout covers
		// the four controller rollouts plus the runtime webhook
		// registration on a cold image pull.
		SkipAwait:     pulumi.Bool(false),
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(600),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, releaseDeps...)

	_, err = helmv3.NewRelease(ctx, locals.ReleaseName, releaseArgs, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to install kyverno helm release")
	}

	exportOutputs(ctx, locals)
	return nil
}
