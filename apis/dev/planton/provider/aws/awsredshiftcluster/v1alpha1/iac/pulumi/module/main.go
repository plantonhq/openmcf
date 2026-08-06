package module

import (
	"github.com/pkg/errors"
	awsredshiftclusterv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsredshiftcluster/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources provisions the Redshift cluster and its folded settings. The
// cluster composes onto its neighbors instead of embedding them: subnets,
// security groups, IAM roles, KMS keys, and the Elastic IP attach by
// reference, and warehouse ingress rules live on the referenced
// AwsSecurityGroup nodes -- this module never creates or mutates a
// resource that deserves to be its own node. Audit logging and
// cross-region snapshot copy are cluster settings keyed by the cluster
// itself, so they are managed here rather than modeled as standalone
// kinds.
func Resources(ctx *pulumi.Context, stackInput *awsredshiftclusterv1alpha1.AwsRedshiftClusterStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder,
	// which resolves the right credential mechanism (static keys, keyless
	// web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsRedshiftCluster.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	createdSubnetGroup, err := subnetGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create subnet group")
	}

	createdParameterGroup, err := clusterParameterGroup(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "failed to create parameter group")
	}

	createdCluster, err := redshiftCluster(ctx, locals, provider, createdSubnetGroup, createdParameterGroup)
	if err != nil {
		return errors.Wrap(err, "failed to create Redshift cluster")
	}

	if err := clusterLogging(ctx, locals, provider, createdCluster); err != nil {
		return errors.Wrap(err, "failed to configure audit logging")
	}

	if err := clusterSnapshotCopy(ctx, locals, provider, createdCluster); err != nil {
		return errors.Wrap(err, "failed to configure snapshot copy")
	}

	ctx.Export(OpClusterIdentifier, createdCluster.ClusterIdentifier)
	ctx.Export(OpClusterArn, createdCluster.Arn)
	ctx.Export(OpClusterNamespaceArn, createdCluster.ClusterNamespaceArn)
	ctx.Export(OpEndpoint, createdCluster.Endpoint)
	ctx.Export(OpDnsName, createdCluster.DnsName)
	ctx.Export(OpDatabaseName, createdCluster.DatabaseName)
	ctx.Export(OpPort, createdCluster.Port)

	// Group names come from the cluster's own attributes so the output
	// shape is identical whether the groups are managed here, referenced,
	// or left to the Redshift defaults.
	ctx.Export(OpSubnetGroupName, createdCluster.ClusterSubnetGroupName)
	ctx.Export(OpParameterGroupName, createdCluster.ClusterParameterGroupName)

	// The AWS-managed admin-password secret exists only when
	// manage_master_password is on; the attribute resolves to "" otherwise,
	// so the output shape is stable across both password strategies.
	ctx.Export(OpMasterPasswordSecretArn, createdCluster.MasterPasswordSecretArn)

	return nil
}
