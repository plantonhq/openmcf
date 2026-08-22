package module

import (
	"github.com/pkg/errors"
	kubernetesplantonoperatorv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesplantonoperator/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the Planton operator from the official
// planton-operator Helm chart (OCI, ghcr.io/plantonhq/charts) as ONE real
// Helm release plus the module-owned PlantonPlatform CRD:
//
//  1. The plantonplatforms.planton.ai CRD, applied from the staged copy at
//     ../crds (extracted from the published chart at the pinned default
//     version) with keep-on-uninstall semantics — destroying the operator
//     never cascade-deletes PlantonPlatform declarations. Applied BEFORE
//     the release through the UPSERT provider, so a reinstall after a
//     destroy adopts the retained CRD instead of failing AlreadyExists.
//  2. The "planton-operator" release, installed with SkipCrds — the module
//     owns the CRD lifecycle (see crds.go for why), so the chart's own
//     crds/ install never runs.
//
// The release name is FIXED: the operator enforces one installation per
// cluster itself at startup (a label-matched Deployment scan), so a second
// release could only crash-loop — the fixed name makes the collision
// impossible to express instead of merely refused. Installing more
// PLATFORMS needs no second operator: one operator watches all namespaces.
//
// The typed spec renders into chart values (values.go); the helm_values
// escape hatch merges last with Helm -f semantics — the exact semantic
// twin of the Terraform module's helm_release with
// values = [typed, helm_values].
func Resources(ctx *pulumi.Context, stackInput *kubernetesplantonoperatorv1alpha1.KubernetesPlantonOperatorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// FAIL LOUDLY below the schema-contract floor: charts below the floor
	// ship operators whose reconcilers predate the PlantonPlatform schema
	// the staged CRD advertises — the API server would ACCEPT fields the
	// running operator silently ignores. Twin: the Terraform module's
	// lifecycle precondition.
	ok, err := chartVersionAtLeast(locals.ChartVersion, vars.MinChartVersion)
	if err != nil {
		return errors.Wrapf(err, "failed to parse chart version %q", locals.ChartVersion)
	}
	if !ok {
		return errors.Errorf(
			"chart version %q predates the PlantonPlatform schema this catalog models — use %s or newer (older operators silently ignore fields the staged CRD accepts)",
			locals.ChartVersion, vars.MinChartVersion)
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

	// ------------------------------ CRD (module-owned) --------------------
	// Skipped only when something else manages the CRD (spec.skip_crds).
	// Routed through the UPSERT provider — the CRD is retained on destroy
	// by design, so the next install must ADOPT it (see crds.go).
	var releaseDeps []pulumi.Resource
	if createdNamespace != nil {
		releaseDeps = append(releaseDeps, createdNamespace)
	}
	if !locals.Spec.SkipCrds {
		upsertProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfigUpsert(ctx,
			stackInput.ProviderConfig, "kubernetes-upsert")
		if err != nil {
			return errors.Wrap(err, "failed to create kubernetes upsert provider")
		}
		crds, err := customResourceDefinitions(ctx, locals, upsertProvider)
		if err != nil {
			return errors.Wrap(err, "failed to apply module-owned CRDs")
		}
		releaseDeps = append(releaseDeps, crds...)
	}

	// ------------------------------ operator release ----------------------
	mergedValues, err := buildHelmValues(locals)
	if err != nil {
		return errors.Wrap(err, "failed to build helm values")
	}

	_, err = helmv3.NewRelease(ctx, vars.ReleaseName, &helmv3.ReleaseArgs{
		Name:      pulumi.String(vars.ReleaseName),
		Namespace: pulumi.String(locals.Namespace),
		// OCI chart reference — joined string, no RepositoryOpts (Pulumi's
		// helm.v3.Release does not resolve oci:// through RepositoryOpts
		// the way the Terraform provider does).
		Chart:   pulumi.String(vars.HelmOciRepo + "/" + vars.HelmChartName),
		Version: pulumi.String(locals.ChartVersion),
		Values:  pulumi.ToMap(mergedValues),
		// The module owns namespace creation (create_namespace flag) and
		// the CRD lifecycle (crds.go) — the chart's own crds/ install
		// never runs.
		CreateNamespace: pulumi.Bool(false),
		SkipCrds:        pulumi.Bool(true),
		// Wait for the operator to become Available — a manager that never
		// becomes ready (the one-per-cluster startup guard refusing beside
		// a sibling operator is THE classic case) should fail THIS deploy
		// with a readiness timeout, not surface later as PlantonPlatform
		// resources that mysteriously never reconcile.
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(vars.HelmTimeoutSeconds),
	}, append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider)},
		dependsOn(releaseDeps)...)...)
	if err != nil {
		return errors.Wrap(err, "failed to install planton-operator helm release")
	}

	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpReleaseName, pulumi.String(vars.ReleaseName))

	return nil
}

// dependsOn wraps a possibly-empty dependency list into resource options
// (an empty DependsOn is a valid no-op).
func dependsOn(deps []pulumi.Resource) []pulumi.ResourceOption {
	if len(deps) == 0 {
		return nil
	}
	return []pulumi.ResourceOption{pulumi.DependsOn(deps)}
}
