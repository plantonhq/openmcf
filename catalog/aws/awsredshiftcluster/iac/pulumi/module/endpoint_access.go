package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshift"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// endpointAccesses expose the cluster inside other subnet groups via
// Redshift-managed VPC endpoints -- same-account cross-VPC access
// without peering (RA3 only). Each entry renders one endpoint keyed by
// endpoint name; the returned map exports each endpoint's private DNS
// address. The cross-account grantee side (resource_owner) is
// deliberately not modeled -- a grantee creates its endpoint in its own
// account against an authorization this cluster grants (see
// endpoint_authorization.go).
func endpointAccesses(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdCluster *redshift.Cluster,
) (pulumi.StringMap, error) {
	spec := locals.AwsRedshiftCluster.Spec

	addresses := pulumi.StringMap{}
	for _, endpoint := range spec.EndpointAccesses {
		args := &redshift.EndpointAccessArgs{
			ClusterIdentifier: createdCluster.ClusterIdentifier,
			EndpointName:      pulumi.String(endpoint.EndpointName),
		}

		// An entry without its own group reuses the cluster's subnet
		// group (managed or referenced), yielding an extra endpoint in
		// the cluster's own VPC. The cluster attribute carries the
		// effective group either way.
		if endpoint.SubnetGroupName != "" {
			args.SubnetGroupName = pulumi.String(endpoint.SubnetGroupName)
		} else {
			args.SubnetGroupName = createdCluster.ClusterSubnetGroupName
		}

		// Empty uses the VPC's default security group (the AWS default).
		if len(endpoint.VpcSecurityGroupIds) > 0 {
			securityGroupIds := pulumi.StringArray{}
			for _, securityGroupId := range endpoint.VpcSecurityGroupIds {
				securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
			}
			args.VpcSecurityGroupIds = securityGroupIds
		}

		createdEndpoint, err := redshift.NewEndpointAccess(ctx, "endpoint-access-"+endpoint.EndpointName, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdCluster}))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create endpoint access %s", endpoint.EndpointName)
		}
		addresses[endpoint.EndpointName] = createdEndpoint.Address
	}
	return addresses, nil
}
