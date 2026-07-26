package module

import (
	"github.com/pkg/errors"
	kubernetestektonoperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetestektonoperator/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the Tekton Operator from its released single-file
// manifest — the operator's OFFICIAL distribution (the in-repo Helm chart
// is unpublished, version "devel"). The manifest applies per document:
// the namespace `tekton-operator`, the 14 operator.tekton.dev CRDs, the
// operator and webhook Deployments (with the spec's typed overrides
// patched on), ConfigMaps, RBAC, Services and the webhook cert Secret.
//
// AUTO-INSTALL IS ALWAYS DISABLED: the tekton-config-defaults ConfigMap's
// AUTOINSTALL_COMPONENTS key is patched from the release's "true" to
// "false" so the operator never creates its own TektonConfig — the
// KubernetesTekton declaration kind is the single owner of the cluster's
// Tekton configuration (see autoInstallTransformation). Installing the
// operator alone deploys no Tekton components.
//
// APPLY MODE: the shared provider helper enables server-side apply — the
// TektonConfig CRD document (~70 KB) plus its siblings stay comfortably
// managed, and SSA keeps re-installs tolerant of the operator's own
// field management. The Terraform twin applies with
// server_side_apply = true.
//
// DESTROY SEMANTICS: every document deletes with the resource, INCLUDING
// the CRDs — which cascade-deletes any TektonConfig on the cluster.
// Always destroy the KubernetesTekton resource FIRST while the operator
// still runs (TektonInstallerSet finalizers are operator-processed; see
// the spec's destroy note).
func Resources(ctx *pulumi.Context, stackInput *kubernetestektonoperatorv1.KubernetesTektonOperatorStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(
		ctx, stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to set up kubernetes provider")
	}

	manifest, err := pulumiyaml.NewConfigFile(ctx, locals.ResourceName,
		&pulumiyaml.ConfigFileArgs{
			File: ManifestURL(),
			Transformations: []pulumiyaml.Transformation{
				autoInstallTransformation(),
				deploymentTransformation(locals.Spec),
			},
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrap(err, "failed to apply the tekton-operator release manifest")
	}

	return exportOutputs(ctx, locals, manifest)
}
