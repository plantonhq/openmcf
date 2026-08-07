package module

import (
	"github.com/pkg/errors"
	awsneptuneclusterv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsneptunecluster/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the Neptune cluster and its folded compute. The
// cluster is the shared-storage brain (endpoints, backups, encryption,
// engine lifecycle); the spec's instances list is the compute that serves
// queries, each entry managed as its own provider resource keyed by name.
// Subnets, security groups, KMS keys, and IAM roles compose by reference
// -- this module never creates or mutates a resource that deserves to be
// its own node.
func Resources(ctx *pulumi.Context, stackInput *awsneptuneclusterv1alpha1.AwsNeptuneClusterStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder,
	// which resolves the right credential mechanism (static keys, keyless
	// web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsNeptuneCluster.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdSubnetGroup, err := subnetGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create subnet group")
	}

	createdParameterGroup, err := clusterParameterGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster parameter group")
	}

	createdCluster, err := neptuneCluster(ctx, locals, provider, createdSubnetGroup, createdParameterGroup)
	if err != nil {
		return errors.Wrap(err, "failed to create Neptune cluster")
	}

	createdInstances, err := clusterInstances(ctx, locals, provider, createdCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster instances")
	}

	ctx.Export(OpClusterIdentifier, createdCluster.ClusterIdentifier)
	ctx.Export(OpArn, createdCluster.Arn)
	ctx.Export(OpClusterResourceId, createdCluster.ClusterResourceId)
	ctx.Export(OpEndpoint, createdCluster.Endpoint)
	ctx.Export(OpReaderEndpoint, createdCluster.ReaderEndpoint)
	ctx.Export(OpPort, createdCluster.Port)
	ctx.Export(OpHostedZoneId, createdCluster.HostedZoneId)
	ctx.Export(OpEngineVersionActual, createdCluster.EngineVersion)
	ctx.Export(OpNeptuneSubnetGroupName, createdCluster.NeptuneSubnetGroupName)
	ctx.Export(OpNeptuneClusterParameterGroupName, createdCluster.NeptuneClusterParameterGroupName)

	// Per-instance endpoints in spec order -- empty for headless shapes
	// (restores, replicas, and global-cluster members created without
	// instances).
	instanceEndpoints := make(pulumi.StringArray, 0, len(createdInstances))
	for _, createdInstance := range createdInstances {
		instanceEndpoints = append(instanceEndpoints, createdInstance.Endpoint)
	}
	ctx.Export(OpInstanceEndpoints, instanceEndpoints)

	return nil
}
