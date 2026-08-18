package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/dsql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster creates the DSQL cluster (and its multi-region pairing when
// configured), then exports outputs.
//
// Lifecycle facts the render below depends on:
//   - only multi_region_properties.witness_region replaces the
//     cluster; deletion protection and the KMS key update in place
//     (a key change re-encrypts, no replacement);
//   - the PEERING resource is a disguised UpdateCluster call that AWS
//     accepts only while the cluster is in PENDING_SETUP - creating
//     it right after the cluster (as here) is the one valid order;
//   - the peering has NO update path at the provider (changes error
//     at apply, they do not replace) and a no-op delete - changing
//     peers means recreating the CLUSTER;
//   - force_destroy makes the provider disable deletion protection
//     before deleting.
func cluster(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	clusterArgs := &dsql.ClusterArgs{
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.DeletionProtectionEnabled {
		clusterArgs.DeletionProtectionEnabled = pulumi.Bool(true)
	}
	if spec.ForceDestroy {
		clusterArgs.ForceDestroy = pulumi.Bool(true)
	}
	if spec.KmsEncryptionKey != nil && spec.KmsEncryptionKey.GetValue() != "" {
		clusterArgs.KmsEncryptionKey = pulumi.String(spec.KmsEncryptionKey.GetValue())
	}

	peerArns := []string{}
	if multiRegion := spec.MultiRegion; multiRegion != nil {
		for _, peerArn := range multiRegion.PeerClusterArns {
			peerArns = append(peerArns, peerArn.GetValue())
		}
		// The witness makes the CLUSTER multi-region at create; the
		// peer list lands via the peering resource below.
		clusterArgs.MultiRegionProperties = &dsql.ClusterMultiRegionPropertiesArgs{
			WitnessRegion: pulumi.String(multiRegion.WitnessRegion),
		}
	}

	createdCluster, err := dsql.NewCluster(ctx, "cluster", clusterArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create cluster")
	}

	if multiRegion := spec.MultiRegion; multiRegion != nil {
		if _, err := dsql.NewClusterPeering(ctx, "peering", &dsql.ClusterPeeringArgs{
			Identifier:    createdCluster.Identifier,
			Clusters:      pulumi.ToStringArray(peerArns),
			WitnessRegion: pulumi.String(multiRegion.WitnessRegion),
		}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
			return errors.Wrap(err, "peer cluster")
		}
	}

	ctx.Export(OpIdentifier, createdCluster.Identifier)
	ctx.Export(OpClusterArn, createdCluster.Arn)
	// AWS exposes no endpoint attribute - the documented DNS shape is
	// derived here (the Terraform module derives the identical value).
	ctx.Export(OpEndpoint, pulumi.Sprintf("%s.dsql.%s.on.aws", createdCluster.Identifier, spec.Region))
	ctx.Export(OpVpcEndpointServiceName, createdCluster.VpcEndpointServiceName)
	ctx.Export(OpEncryptionType, createdCluster.EncryptionDetails.ApplyT(func(details []dsql.ClusterEncryptionDetail) string {
		if len(details) == 0 {
			return ""
		}
		return details[0].EncryptionType
	}).(pulumi.StringOutput))
	return nil
}
