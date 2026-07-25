package module

import (
	"github.com/pkg/errors"
	kubernetessignozv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessignoz/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs SigNoz from the official Helm chart as a real Helm
// release. The typed spec renders into chart values (values.go); exactly
// one database arm renders (bundled ClickHouse XOR external); the bundled
// ClickHouse admin password is module-GENERATED, injected as a Pulumi
// secret Output after the values merge (the set_sensitive twin) and
// exported through the module-owned "<name>-clickhouse-auth" Secret; the
// helm_values escape hatch merges last with Helm -f semantics — the exact
// semantic twin of the Terraform module's helm_release with
// values = [typed, helm_values, re-pin] + set_sensitive.
func Resources(ctx *pulumi.Context, stackInput *kubernetessignozv1.KubernetesSignozStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// The chart wraps the ClickHouseInstallation's operator-generated
	// StatefulSet names in ~27 characters of scaffolding within
	// Kubernetes' 63-character cap — an over-long resource name corrupts
	// the naming contract the outputs promise. Fail THIS deploy loudly
	// instead (twin: the Terraform module's plan-time precondition).
	if len(locals.ReleaseName) > vars.MaxNameLength {
		return errors.Errorf(
			"metadata.name %q is %d characters — the signoz chart's child-name budget allows at most %d "+
				"(the bundled ClickHouse composes names like chi-<name>-clickhouse-cluster-0-0 within Kubernetes' 63-character cap)",
			locals.ReleaseName, len(locals.ReleaseName), vars.MaxNameLength)
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

	// ------------------------------ helm values ---------------------------
	mergedValues, err := buildHelmValues(locals)
	if err != nil {
		return errors.Wrap(err, "failed to build helm values")
	}

	// ---------------- bundled clickhouse credential -----------------------
	// Generated and injected AFTER the values merge so it never rides a
	// plain values map, the escape hatch can never override it, and the
	// preview shows it as a secret.
	if !locals.IsExternal {
		password, createdSecret, err := clickHouseAuthSecret(ctx, locals, kubernetesProvider, releaseDeps)
		if err != nil {
			return errors.Wrap(err, "failed to create clickhouse auth secret")
		}
		clickhouse, ok := mergedValues["clickhouse"].(map[string]interface{})
		if !ok {
			return errors.New("internal: bundled-arm values missing the clickhouse block")
		}
		clickhouse["password"] = pulumi.ToSecret(password)
		releaseDeps = append(releaseDeps, pulumi.DependsOn([]pulumi.Resource{createdSecret}))
	}

	// ------------------------------ helm release --------------------------
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
		// Wait for the whole platform to become Ready — a SigNoz whose
		// ClickHouse never binds storage or whose schema migration fails
		// should fail THIS deploy, not the first trace query. SkipAwait
		// false is Helm --wait, stated explicitly to mirror the Terraform
		// twin's `wait = true`; the budget covers the ClickHouse operator
		// reconcile + schema migration on a cold cluster.
		SkipAwait:     pulumi.Bool(false),
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(1200),
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, releaseDeps...)

	_, err = helmv3.NewRelease(ctx, locals.ReleaseName, releaseArgs, opts...)
	if err != nil {
		return errors.Wrap(err, "failed to install signoz helm release")
	}

	exportOutputs(ctx, locals)
	return nil
}

// exportOutputs publishes the composition handles. Every child name
// derives from the two fullnames pinned via fullnameOverride (the release
// name and `<name>-clickhouse`).
func exportOutputs(ctx *pulumi.Context, locals *Locals) {
	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpService, pulumi.String(locals.Service))
	ctx.Export(OpKubeEndpoint, pulumi.String(locals.KubeEndpoint))
	ctx.Export(OpPortForwardCommand, pulumi.String(locals.PortForwardCommand))
	ctx.Export(OpOtelCollectorService, pulumi.String(locals.OtelCollectorService))
	ctx.Export(OpOtlpGrpcEndpoint, pulumi.String(locals.OtlpGrpcEndpoint))
	ctx.Export(OpOtlpHttpEndpoint, pulumi.String(locals.OtlpHttpEndpoint))
	ctx.Export(OpClickHouseEndpoint, pulumi.String(locals.ClickHouseEndpoint))
	ctx.Export(OpClickHouseUsername, pulumi.String(locals.ClickHouseUsername))
	ctx.Export(OpClickHousePasswordSecretName, pulumi.String(locals.ClickHousePasswordSecretName))
	ctx.Export(OpClickHousePasswordSecretKey, pulumi.String(locals.ClickHousePasswordSecretKey))
}
