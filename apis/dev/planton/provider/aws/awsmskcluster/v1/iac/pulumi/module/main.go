package module

import (
	"github.com/pkg/errors"
	awsmskclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskcluster/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources orchestrates creation of AWS MSK Cluster related resources and exports outputs.
func Resources(ctx *pulumi.Context, stackInput *awsmskclusterv1.AwsMskClusterStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	// Build the AWS provider from the stack input via the shared builder, which resolves
	// the right credential mechanism (static keys, keyless web identity, or ambient chain).
	provider, err := pulumiawsprovider.Get(ctx, stackInput.ProviderConfig, locals.AwsMskCluster.Spec.Region)
	if err != nil {
		return errors.Wrap(err, "failed to create AWS provider")
	}

	// Module-managed MSK Configuration (when server_properties provided)
	createdConfig, err := configuration(ctx, locals, provider)
	if err != nil {
		return errors.Wrap(err, "msk configuration")
	}

	// MSK Cluster
	createdCluster, err := cluster(ctx, locals, provider, createdConfig)
	if err != nil {
		return errors.Wrap(err, "msk cluster")
	}

	// Cluster-scoped satellites (SCRAM secret associations, resource policy)
	if err := satellites(ctx, locals, provider, createdCluster); err != nil {
		return errors.Wrap(err, "msk cluster satellites")
	}

	// Export outputs
	ctx.Export(OpClusterArn, createdCluster.Arn)
	ctx.Export(OpClusterName, createdCluster.ClusterName)
	ctx.Export(OpClusterUuid, createdCluster.ClusterUuid)
	ctx.Export(OpCurrentVersion, createdCluster.CurrentVersion)
	ctx.Export(OpBootstrapBrokers, createdCluster.BootstrapBrokers)
	ctx.Export(OpBootstrapBrokersTls, createdCluster.BootstrapBrokersTls)
	ctx.Export(OpBootstrapBrokersSaslIam, createdCluster.BootstrapBrokersSaslIam)
	ctx.Export(OpBootstrapBrokersSaslScram, createdCluster.BootstrapBrokersSaslScram)
	ctx.Export(OpBootstrapBrokersPublicTls, createdCluster.BootstrapBrokersPublicTls)
	ctx.Export(OpBootstrapBrokersPublicSaslIam, createdCluster.BootstrapBrokersPublicSaslIam)
	ctx.Export(OpBootstrapBrokersPublicSaslScram, createdCluster.BootstrapBrokersPublicSaslScram)
	ctx.Export(OpBootstrapBrokersVpcConnectivityTls, createdCluster.BootstrapBrokersVpcConnectivityTls)
	ctx.Export(OpBootstrapBrokersVpcConnectivitySaslIam, createdCluster.BootstrapBrokersVpcConnectivitySaslIam)
	ctx.Export(OpBootstrapBrokersVpcConnectivitySaslScram, createdCluster.BootstrapBrokersVpcConnectivitySaslScram)
	ctx.Export(OpZookeeperConnectString, createdCluster.ZookeeperConnectString)
	ctx.Export(OpZookeeperConnectStringTls, createdCluster.ZookeeperConnectStringTls)
	if createdConfig != nil {
		ctx.Export(OpConfigurationArn, createdConfig.Arn)
	}

	return nil
}
