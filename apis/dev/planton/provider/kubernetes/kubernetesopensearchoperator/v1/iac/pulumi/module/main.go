package module

import (
	"github.com/pkg/errors"
	kubernetesopensearchoperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesopensearchoperator/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the OpenSearch Kubernetes Operator from the official
// opensearch-operator Helm chart as a single Helm release named after
// metadata.name. The operator reconciles OpenSearchCluster custom
// resources (declared through KubernetesOpenSearch) into running search
// clusters with managed TLS, security bootstrap, rolling upgrades and
// Dashboards.
//
// CRD LIFECYCLE: the chart templates its ten CRDs as release-owned
// resources with NO keep-on-uninstall knob, so a Helm-owned install would
// cascade-delete every OpenSearchCluster (and its data) on uninstall. The
// module therefore OWNS the CRDs: installCRDs pins false unconditionally
// in the rendered values, and crds.go applies the staged CRD files with
// retainOnDelete — destroy drops them from state without deleting them
// from the cluster. The release depends on the CRDs so the operator never
// starts against an unregistered API group.
//
// The typed spec renders into chart values (values.go); the helm_values
// escape hatch merges last with Helm -f semantics — the exact semantic
// twin of the Terraform module's helm_release with
// values = [typed, helm_values].
func Resources(ctx *pulumi.Context, stackInput *kubernetesopensearchoperatorv1.KubernetesOpenSearchOperatorStackInput) error {
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

	// ------------------------------ CRDs ----------------------------------
	// Module-owned, retained on delete — see crds.go for the full posture.
	createdCrds, err := customResourceDefinitions(ctx, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to apply opensearch operator CRDs")
	}

	operatorDeps := append([]pulumi.Resource{}, createdCrds...)
	if createdNamespace != nil {
		operatorDeps = append(operatorDeps, createdNamespace)
	}

	// ------------------------------ operator release ----------------------
	mergedValues, err := buildHelmValues(locals)
	if err != nil {
		return errors.Wrap(err, "failed to build helm values")
	}

	_, err = helmv3.NewRelease(ctx, locals.ReleaseName, &helmv3.ReleaseArgs{
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
		// Wait for the operator to become Available — an operator that
		// never becomes ready (an unpullable image from a private mirror
		// is the classic case) should fail THIS deploy with a readiness
		// timeout, not surface later as OpenSearch clusters that
		// mysteriously never reconcile.
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(vars.HelmTimeoutSeconds),
	}, append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider)},
		dependsOn(operatorDeps)...)...)
	if err != nil {
		return errors.Wrap(err, "failed to install opensearch-operator helm release")
	}

	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpReleaseName, pulumi.String(locals.ReleaseName))
	ctx.Export(OpDeploymentName, pulumi.String(locals.DeploymentName))

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
