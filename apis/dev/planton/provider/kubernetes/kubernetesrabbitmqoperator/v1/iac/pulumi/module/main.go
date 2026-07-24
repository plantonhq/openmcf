package module

import (
	"github.com/pkg/errors"
	kubernetesrabbitmqoperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesrabbitmqoperator/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources installs the RabbitMQ Cluster Operator from its released
// single-file manifest — the operator's OFFICIAL distribution (it has no
// Helm chart). The manifest applies per document, with the spec's typed
// overrides patched onto the operator Deployment and every other document
// applied verbatim: the namespace `rabbitmq-system`, the RabbitmqCluster
// CRD, RBAC, the webhook + metrics Services, the cert-manager Issuer and
// Certificates, and the mutating/validating webhook configurations.
//
// CERT-MANAGER IS A HARD PREREQUISITE (a registry prerequisite of this
// kind): the webhook serving certificate is a cert-manager Certificate
// with CA injection — without a running cert-manager the certificate
// never issues and every RabbitmqCluster admission fails (the webhooks
// are failurePolicy: Fail).
//
// APPLY MODE: the shared provider helper enables server-side apply, which
// this manifest REQUIRES — the RabbitmqCluster CRD document (~342 KB)
// exceeds the client-side last-applied-configuration annotation cap
// (256 KB). The Terraform twin applies with server_side_apply = true for
// the same reason.
//
// DESTROY SEMANTICS: every document deletes with the resource, INCLUDING
// the CRD — which cascade-deletes every RabbitmqCluster on the cluster.
// The spec's CRD-lifecycle note carries the warning.
func Resources(ctx *pulumi.Context, stackInput *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorStackInput) error {
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
				deploymentTransformation(locals.Spec),
			},
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrap(err, "failed to apply the cluster-operator release manifest")
	}

	return exportOutputs(ctx, locals, manifest)
}
