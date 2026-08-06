package module

import (
	"strings"

	"github.com/pkg/errors"
	kuberneteshelmreleasev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteshelmrelease/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the main entry point for the Pulumi module.
// It installs the chart as a REAL Helm release (helm.v3.Release): hooks run,
// the release secret is written, and `helm list` sees the release — the
// exact semantic twin of the Terraform module's helm_release resource.
// (The render-only helm.v3.Chart resource is deliberately NOT used: it
// template-renders client-side without creating a release, which silently
// skips hooks and leaves nothing for Helm tooling to manage.)
func Resources(ctx *pulumi.Context, stackInput *kuberneteshelmreleasev1alpha1.KubernetesHelmReleaseStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Create kubernetes provider from the credential in the stack-input
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

	// Build conditional namespace dependency (Pulumi equivalent of Terraform depends_on).
	var namespaceDeps []pulumi.ResourceOption
	if createdNamespace != nil {
		namespaceDeps = append(namespaceDeps, pulumi.DependsOn([]pulumi.Resource{createdNamespace}))
	}

	// ------------------------------ helm release --------------------------
	createdRelease, err := helmRelease(ctx, locals, kubernetesProvider, namespaceDeps)
	if err != nil {
		return errors.Wrap(err, "failed to create helm release")
	}

	exportOutputs(ctx, locals, createdRelease)
	return nil
}

// helmRelease installs the chart with the merged values and the spec's
// lifecycle knobs. Every knob maps 1:1 onto a helm_release argument in the
// Terraform module — keep the two in lockstep.
func helmRelease(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	namespaceDeps []pulumi.ResourceOption) (*helmv3.Release, error) {

	spec := locals.Spec

	// PARITY-EXCEPTION: take_ownership is Pulumi-inexpressible at the pinned
	// SDK (helm.v3.ReleaseArgs gained TakeOwnership only in a later
	// pulumi-kubernetes release). The Terraform module implements it. A set
	// field must never be silently dropped, so fail the deploy loudly and
	// route the user to the working engine. Dissolves at the next
	// pulumi-kubernetes SDK upgrade.
	if spec.GetTakeOwnership() {
		return nil, errors.New("take_ownership is not supported by the pulumi engine at the current pulumi-kubernetes SDK version; " +
			"deploy this release with the terraform provisioner, or drop take_ownership")
	}

	mergedValues, err := buildHelmValues(spec)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build helm values")
	}

	valuesInput := pulumi.MapInput(pulumi.ToMap(mergedValues))
	if len(spec.GetSetSensitive()) > 0 {
		// The merged map carries secret entries: mark the WHOLE values map
		// secret in Pulumi state. Coarser than Terraform's per-entry
		// set_sensitive masking, but errs on the safe side — documented
		// behavioral difference, not a parity exception (the release Helm
		// installs is identical).
		valuesInput = pulumi.ToSecret(pulumi.ToMap(mergedValues)).(pulumi.MapOutput)
	}

	// Chart resolution mirrors the Terraform provider exactly:
	// - OCI registry: the chart reference is <repo>/<chart>, no repository opts URL.
	// - HTTP(S) repository: chart name stays bare; the repo URL goes in
	//   repository opts.
	// Private-repo credentials ride repository opts in both forms.
	chartRef := spec.GetChart()
	repositoryOpts := &helmv3.RepositoryOptsArgs{}
	if strings.HasPrefix(spec.GetRepo(), "oci://") {
		chartRef = strings.TrimSuffix(spec.GetRepo(), "/") + "/" + spec.GetChart()
	} else {
		repositoryOpts.Repo = pulumi.String(spec.GetRepo())
	}
	if spec.GetRepositoryUsername() != "" {
		repositoryOpts.Username = pulumi.String(spec.GetRepositoryUsername())
		repositoryOpts.Password = pulumi.ToSecret(pulumi.String(spec.GetRepositoryPassword())).(pulumi.StringOutput)
	}

	releaseArgs := &helmv3.ReleaseArgs{
		Name:            pulumi.String(locals.ReleaseName),
		Namespace:       pulumi.String(locals.Namespace),
		Chart:           pulumi.String(chartRef),
		Version:         pulumi.String(spec.GetVersion()),
		RepositoryOpts:  repositoryOpts,
		Values:          valuesInput,
		CreateNamespace: pulumi.Bool(false), // the module owns namespace creation (create_namespace flag)

		// Lifecycle knobs — 1:1 with the Terraform module's helm_release
		// arguments. SkipAwait inverts the TF `wait` flag: both engines
		// default to awaiting readiness.
		Atomic:                   pulumi.Bool(spec.GetAtomic()),
		CleanupOnFail:            pulumi.Bool(spec.GetCleanupOnFail()),
		SkipAwait:                pulumi.Bool(spec.GetSkipAwait()),
		WaitForJobs:              pulumi.Bool(spec.GetWaitForJobs()),
		Timeout:                  pulumi.Int(locals.TimeoutSeconds),
		SkipCrds:                 pulumi.Bool(spec.GetSkipCrds()),
		DependencyUpdate:         pulumi.Bool(spec.GetDependencyUpdate()),
		MaxHistory:               pulumi.Int(locals.MaxHistory),
		Replace:                  pulumi.Bool(spec.GetReplace()),
		ForceUpdate:              pulumi.Bool(spec.GetForceUpdate()),
		ReuseValues:              pulumi.Bool(spec.GetReuseValues()),
		ResetValues:              pulumi.Bool(spec.GetResetValues()),
		DisableWebhooks:          pulumi.Bool(spec.GetDisableWebhooks()),
		DisableOpenapiValidation: pulumi.Bool(spec.GetDisableOpenapiValidation()),
	}
	if spec.GetDescription() != "" {
		releaseArgs.Description = pulumi.String(spec.GetDescription())
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, namespaceDeps...)

	createdRelease, err := helmv3.NewRelease(ctx, locals.ReleaseName, releaseArgs, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to install helm release")
	}

	return createdRelease, nil
}
