package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshiftserverless"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// endpointAccesses expose the workgroup inside other subnets via VPC
// endpoints -- same-account cross-VPC access without peering. Each entry
// renders one endpoint keyed by endpoint name; the returned map exports
// each endpoint's private DNS address. The cross-account grantee side
// (owner_account) is deliberately not modeled -- it lives in the
// grantee's credential domain.
func endpointAccesses(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdWorkgroup *redshiftserverless.Workgroup,
) (pulumi.StringMap, error) {
	spec := locals.AwsRedshiftServerlessWorkgroup.Spec

	addresses := pulumi.StringMap{}
	for _, endpoint := range spec.EndpointAccesses {
		// An entry without its own subnets reuses the workgroup's (the
		// spec CEL guarantees the fallback exists).
		sourceSubnets := endpoint.SubnetIds
		if len(sourceSubnets) == 0 {
			sourceSubnets = spec.SubnetIds
		}
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range sourceSubnets {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}

		args := &redshiftserverless.EndpointAccessArgs{
			WorkgroupName: createdWorkgroup.WorkgroupName,
			EndpointName:  pulumi.String(endpoint.EndpointName),
			SubnetIds:     subnetIds,
		}

		// Empty uses the VPC's default security group (the AWS default).
		if len(endpoint.VpcSecurityGroupIds) > 0 {
			securityGroupIds := pulumi.StringArray{}
			for _, securityGroupId := range endpoint.VpcSecurityGroupIds {
				securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
			}
			args.VpcSecurityGroupIds = securityGroupIds
		}

		createdEndpoint, err := redshiftserverless.NewEndpointAccess(ctx, "endpoint-access-"+endpoint.EndpointName, args,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdWorkgroup}))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create endpoint access %s", endpoint.EndpointName)
		}
		addresses[endpoint.EndpointName] = createdEndpoint.Address
	}
	return addresses, nil
}
