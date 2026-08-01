package module

import (
	"github.com/pkg/errors"
	kubernetessolroperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessolroperator/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the Apache Solr Operator from the official
// solr-operator Helm chart as a single Helm release named after
// metadata.name. The operator reconciles SolrCloud custom resources
// (declared through KubernetesSolr) into running Solr clusters, plus
// SolrBackup and SolrPrometheusExporter resources.
//
// CRD LIFECYCLE: unlike most operator charts, the solr-operator chart
// ships NO CRDs — they are separate release artifacts. The module OWNS
// them: the four staged files at ../crds (the three solr.apache.org CRDs
// plus the ZookeeperCluster CRD of the bundled zookeeper-operator
// dependency) are applied before the release, each keyed by its own
// metadata.name and carrying retainOnDelete — destroying the stack
// removes the operator but NEVER the CRDs, so SolrCloud resources are
// never cascade-deleted. The bundled subchart's own CRD switch is pinned
// off (`zookeeper-operator.crd.create: false`) so the ZookeeperCluster
// CRD never falls under Helm's delete-on-uninstall lifecycle.
//
// The typed spec renders into chart values (values.go); the helm_values
// escape hatch merges last with Helm -f semantics — the exact semantic
// twin of the Terraform module's helm_release with
// values = [typed, helm_values].
func Resources(ctx *pulumi.Context, stackInput *kubernetessolroperatorv1.KubernetesSolrOperatorStackInput) error {
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

	var operatorDeps []pulumi.Resource
	if createdNamespace != nil {
		operatorDeps = append(operatorDeps, createdNamespace)
	}

	// ------------------------------ CRDs ----------------------------------
	// Applied before the release: the operator's controllers start
	// watching these types immediately, and the bundled
	// zookeeper-operator refuses to start without the ZookeeperCluster
	// CRD present.
	//
	// The CRDs ride a DEDICATED upsert provider: retained-on-destroy
	// resources are, by design, already on the cluster the next time
	// this module installs, and a plain create fails AlreadyExists (the
	// provider adopts only with upsertExistingObjects — verified in the
	// pinned provider source). The upsert scope stays CRD-only; the
	// release keeps the plain provider's create-conflict semantics.
	upsertProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfigUpsert(ctx,
		stackInput.ProviderConfig, "kubernetes-crd-upsert")
	if err != nil {
		return errors.Wrap(err, "failed to create the CRD upsert kubernetes provider")
	}
	createdCrds, err := applyCrds(ctx, upsertProvider)
	if err != nil {
		return errors.Wrap(err, "failed to apply solr-operator crds")
	}
	operatorDeps = append(operatorDeps, createdCrds...)

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
		// timeout, not surface later as SolrCloud resources that
		// mysteriously never reconcile.
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(vars.HelmTimeoutSeconds),
	}, append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider)},
		dependsOn(operatorDeps)...)...)
	if err != nil {
		return errors.Wrap(err, "failed to install solr-operator helm release")
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
