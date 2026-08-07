package module

import (
	"github.com/pkg/errors"
	kubernetesoteloperatorv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesoteloperator/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the OpenTelemetry Operator from the official
// opentelemetry-operator Helm chart as a single Helm release named after
// metadata.name. The operator reconciles OpenTelemetryCollector custom
// resources (declared through KubernetesOtelCollector) into running
// collector fleets, defaulting and validating them through its admission
// webhooks.
//
// CRD LIFECYCLE: the chart templates its opentelemetry.io CRDs as
// release-owned resources, so a Helm-owned install would cascade-delete
// every collector declaration on uninstall. The module therefore OWNS the
// CRDs: crds.create pins false unconditionally in the rendered values,
// and crds.go applies the staged, tokenized CRD files with retainOnDelete
// — destroy drops them from state without deleting them from the cluster.
// The release depends on the CRDs so the operator never starts against an
// unregistered API group.
//
// THE CONVERSION-TRUST COUPLING (why cert-manager is required): the
// collector CRD carries a version-conversion webhook and the
// cert-manager.io/inject-ca-from annotation. Because the CRDs are
// retained past the release's lifetime, their conversion trust must be
// kept current by a RUNNING reconciler — cert-manager's CA injector —
// not by a certificate embedded once at install time.
//
// The typed spec renders into chart values (values.go); the helm_values
// escape hatch merges last with Helm -f semantics — the exact semantic
// twin of the Terraform module's helm_release values documents.
func Resources(ctx *pulumi.Context, stackInput *kubernetesoteloperatorv1alpha1.KubernetesOtelOperatorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Fail-loud name budget: 63-char Kubernetes name limit minus the
	// chart's longest derived suffix, "-controller-manager-service-cert"
	// (33 chars). The module pins fullnameOverride to metadata.name, so
	// the budget is exact. Twin of the Terraform module's precondition.
	if len(locals.ReleaseName) > 30 {
		return errors.Errorf("metadata.name %q is %d characters — must be 30 or fewer: the chart derives \"<name>-controller-manager-service-cert\" (33-char suffix) and Kubernetes caps names at 63",
			locals.ReleaseName, len(locals.ReleaseName))
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

	// ------------------------------ CRDs ----------------------------------
	// Module-owned, retained on delete — see crds.go for the full posture.
	// skip_crds is the bring-your-own-CRDs arm: the CRDs are owned
	// elsewhere (a GitOps-managed bundle) and this module must not touch
	// them.
	//
	// The CRDs ride a DEDICATED upsert provider: retained-on-destroy
	// resources are, by design, already on the cluster the next time
	// this module installs, and a plain create fails AlreadyExists (the
	// provider adopts only with upsertExistingObjects — verified in the
	// pinned provider source). The upsert scope stays CRD-only; the
	// release keeps the plain provider's create-conflict semantics.
	var operatorDeps []pulumi.Resource
	if !stackInput.Target.Spec.SkipCrds {
		upsertProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfigUpsert(ctx,
			stackInput.ProviderConfig, "kubernetes-crd-upsert")
		if err != nil {
			return errors.Wrap(err, "failed to create the CRD upsert kubernetes provider")
		}
		createdCrds, err := customResourceDefinitions(ctx, locals, upsertProvider)
		if err != nil {
			return errors.Wrap(err, "failed to apply opentelemetry operator CRDs")
		}
		operatorDeps = append(operatorDeps, createdCrds...)
	}
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
		// Wait for the operator to become Available — the manager pod
		// mounts the cert-manager-issued webhook Secret, so an install
		// without a working cert-manager (or with an unpullable image)
		// should fail THIS deploy with a readiness timeout, not surface
		// later as collectors that mysteriously never reconcile.
		Atomic:        pulumi.Bool(true),
		CleanupOnFail: pulumi.Bool(true),
		Timeout:       pulumi.Int(vars.HelmTimeoutSeconds),
	}, append([]pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider)},
		dependsOn(operatorDeps)...)...)
	if err != nil {
		return errors.Wrap(err, "failed to install opentelemetry-operator helm release")
	}

	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpReleaseName, pulumi.String(locals.ReleaseName))
	ctx.Export(OpWebhookService, pulumi.String(locals.WebhookService))
	ctx.Export(OpWebhookCertSecretName, pulumi.String(locals.WebhookCertSecretName))

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
