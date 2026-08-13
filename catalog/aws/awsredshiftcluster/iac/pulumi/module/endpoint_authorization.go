package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// endpointAuthorizations grant OTHER AWS accounts permission to create
// managed VPC endpoints to this cluster -- the grantor side of
// cross-account access, living in this cluster's own credential domain.
// Each entry renders one authorization keyed by the grantee account
// (AWS keeps one authorization per account).
func endpointAuthorizations(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdCluster *redshift.Cluster,
) error {
	spec := locals.AwsRedshiftCluster.Spec

	for _, authorization := range spec.EndpointAuthorizations {
		args := &redshift.EndpointAuthorizationArgs{
			ClusterIdentifier: createdCluster.ClusterIdentifier,
			Account:           pulumi.String(authorization.Account),

			// When true, delete revokes the grant even while the grantee
			// still has live endpoints (deleting them too). AWS's
			// default refuses.
			ForceDelete: pulumi.Bool(authorization.ForceDelete),
		}

		// Empty authorizes ALL of the grantee account's VPCs.
		if len(authorization.VpcIds) > 0 {
			vpcIds := pulumi.StringArray{}
			for _, vpcId := range authorization.VpcIds {
				vpcIds = append(vpcIds, pulumi.String(vpcId.GetValue()))
			}
			args.VpcIds = vpcIds
		}

		if _, err := redshift.NewEndpointAuthorization(ctx, "endpoint-authorization-"+authorization.Account, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster})); err != nil {
			return errors.Wrapf(err, "failed to create endpoint authorization for account %s", authorization.Account)
		}
	}
	return nil
}
