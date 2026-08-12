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
//
// Endpoint-access CREATE and DELETE both answer 400 ConflictException
// ("An operation is running on the serverless workgroup") unless the
// workgroup is idle -- the workgroup stays busy for ~15-30s AFTER a
// usage-limit call returns -- and the provider carries no
// ConflictException retry on either (only on the workgroup's own
// delete/update). Endpoint accesses therefore apply straight after the
// workgroup (idle from the provider's own wait-for-available; extraDeps
// carries the custom domain, whose window is unproven -- its live arm
// is deferred), with the conflict-immune usage limits chained behind
// the returned resources.
//
// PARITY-EXCEPTION: destroy MECHANICS differ from the Terraform module
// by design. Each endpoint is DeletedWith(workgroup): a full destroy
// skips the conflict-prone DeleteEndpointAccess call entirely and the
// workgroup's own delete -- which the provider retries on conflict and
// AWS cascades over live endpoint accesses (live-probed) -- removes it.
// Terraform cannot skip a tracked resource's delete, so it crosses the
// same window behind a time_sleep destroy settle instead
// (iac/tf/satellite_settle.tf). End state and outputs are identical.
// Entries WITHIN this group still dispatch concurrently -- a durable
// fix for many-entry specs is a provider-side retry, recorded upstream.
func endpointAccesses(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdWorkgroup *redshiftserverless.Workgroup,
	extraDeps []pulumi.Resource,
) (pulumi.StringMap, []pulumi.Resource, error) {
	spec := locals.AwsRedshiftServerlessWorkgroup.Spec

	dependencies := append([]pulumi.Resource{createdWorkgroup}, extraDeps...)
	createdEndpoints := []pulumi.Resource{}
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
			pulumi.Provider(provider), pulumi.DependsOn(dependencies), pulumi.DeletedWith(createdWorkgroup))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "failed to create endpoint access %s", endpoint.EndpointName)
		}
		createdEndpoints = append(createdEndpoints, createdEndpoint)
		addresses[endpoint.EndpointName] = createdEndpoint.Address
	}
	return addresses, createdEndpoints, nil
}
