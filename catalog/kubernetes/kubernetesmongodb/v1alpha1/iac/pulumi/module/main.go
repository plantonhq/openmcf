package module

import (
	"github.com/pkg/errors"
	kubernetesmongodbv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesmongodb/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/pulumikubernetesprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources deploys one Percona-operator-managed MongoDB cluster:
//
//  1. the namespace (optional),
//  2. declared-credential Secrets (user passwords, backup-storage keys)
//     — secrets always travel via secret references, never inline in a
//     custom resource,
//  3. the PerconaServerMongoDB CR itself (rendered untyped — see
//     cluster.go for why — and validated server-side by the operator's
//     CRD schema).
//
// Ordering matters only for the namespace (everything is namespaced) and
// for credential Secrets (the operator reads them at reconcile time).
func Resources(ctx *pulumi.Context, stackInput *kubernetesmongodbv1alpha1.KubernetesMongodbStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	kubernetesProvider, err := pulumikubernetesprovider.GetWithKubernetesProviderConfig(ctx,
		stackInput.ProviderConfig, "kubernetes")
	if err != nil {
		return errors.Wrap(err, "failed to set up kubernetes provider")
	}

	createdNamespace, err := namespace(ctx, stackInput, locals, kubernetesProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create namespace")
	}

	var namespaceDeps []pulumi.ResourceOption
	if createdNamespace != nil {
		namespaceDeps = append(namespaceDeps, pulumi.DependsOn([]pulumi.Resource{createdNamespace}))
	}

	credentialSecrets, err := createCredentialSecrets(ctx, locals, kubernetesProvider, namespaceDeps)
	if err != nil {
		return errors.Wrap(err, "failed to create credential secrets")
	}

	clusterDeps := namespaceDeps
	if len(credentialSecrets) > 0 {
		clusterDeps = append(clusterDeps, pulumi.DependsOn(credentialSecrets))
	}

	if _, err := createCluster(ctx, locals, kubernetesProvider, clusterDeps); err != nil {
		return errors.Wrap(err, "failed to create PerconaServerMongoDB cluster")
	}

	exportOutputs(ctx, locals)
	return nil
}

func exportOutputs(ctx *pulumi.Context, locals *Locals) {
	ctx.Export(OpNamespace, pulumi.String(locals.Namespace))
	ctx.Export(OpClusterName, pulumi.String(locals.ClusterName))
	ctx.Export(OpService, pulumi.String(locals.ServiceName))
	ctx.Export(OpKubeEndpoint, pulumi.String(locals.KubeEndpoint))
	ctx.Export(OpReplicaSet, pulumi.String(locals.ReplicaSetOutput))
	ctx.Export(OpPortForwardCommand, pulumi.String(locals.PortForwardCommand))
	ctx.Export(OpAdminPasswordSecret, pulumi.Map{
		"name": pulumi.String(locals.UsersSecretName),
		"key":  pulumi.String(vars.AdminPasswordKey),
	})
}
