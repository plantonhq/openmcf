package module

import (
	"github.com/pkg/errors"
	kubernetesgrafanav1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesgrafana/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs Grafana from the official Helm chart as a real Helm
// release. The typed spec renders into chart values (values.go); admin
// credentials stay chart-owned (generated once via the chart's lookup, or
// read from an existing Secret); database and datasource credentials ride
// environment variables sourced from Secrets so no credential ever lands in
// the chart's rendered configuration; the helm_values escape hatch merges
// last with Helm -f semantics — the exact semantic twin of the Terraform
// module's helm_release with values = [typed, helm_values].
func Resources(ctx *pulumi.Context, stackInput *kubernetesgrafanav1alpha1.KubernetesGrafanaStackInput) error {
	locals := initializeLocals(ctx, stackInput)

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
		// Wait for Grafana to become Ready — a UI that never starts (bad
		// plugin ID, unreachable database, unbindable volume) should fail
		// THIS deploy, not the first login attempt. Plugin downloads at
		// startup are what the generous budget covers. SkipAwait false is
		// Helm --wait, stated explicitly to mirror the Terraform twin's
		// `wait = true`.
		SkipAwait:     pulumi.Bool(false),
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(600),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, releaseDeps...)

	_, err = helmv3.NewRelease(ctx, locals.ReleaseName, releaseArgs, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to install grafana helm release")
	}

	exportOutputs(ctx, locals)
	return nil
}

// exportOutputs publishes the composition handles. The service name is the
// chart's ClusterIP Service — grafana.fullname, pinned to the resource name
// via fullnameOverride. The admin Secret name follows the credential arm:
// the chart-owned `<name>` Secret for the generate arm, the referenced
// Secret's own name for the existing arm.
func exportOutputs(ctx *pulumi.Context, locals *Locals) {
	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpReleaseName, pulumi.String(locals.ReleaseName))
	ctx.Export(OpService, pulumi.String(locals.ServiceName))
	ctx.Export(OpEndpoint, pulumi.String(locals.Endpoint))
	ctx.Export(OpAdminSecretName, pulumi.String(locals.AdminSecretName))
	ctx.Export(OpPortForwardCommand, pulumi.String(locals.PortForwardCommand))
}
