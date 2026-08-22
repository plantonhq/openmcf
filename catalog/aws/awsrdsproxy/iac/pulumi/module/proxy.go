package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// proxy creates the RDS proxy, its default-target-group pool tuning,
// additional endpoints, and the database target, then exports outputs.
//
// Lifecycle facts the render below depends on:
//   - engine_family, vpc_subnet_ids, and both network-type dials are
//     ForceNew on the proxy; everything else updates in place;
//   - the default target group is a PATCH satellite: it always exists
//     on the proxy, its provider delete is a no-op, and managing it
//     here just tunes the pool (name is always "default");
//   - the target registration waits out a database still in CREATING
//     (the provider retries for 5 minutes) - the module still orders
//     it after the target group so plans read cleanly;
//   - a proxy fronts exactly ONE database (instance XOR cluster - the
//     spec's CEL wall mirrors AWS's contract).
func proxy(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	subnetIds := make([]string, 0, len(spec.VpcSubnetIds))
	for _, subnetId := range spec.VpcSubnetIds {
		subnetIds = append(subnetIds, subnetId.GetValue())
	}
	securityGroupIds := make([]string, 0, len(spec.VpcSecurityGroupIds))
	for _, securityGroupId := range spec.VpcSecurityGroupIds {
		securityGroupIds = append(securityGroupIds, securityGroupId.GetValue())
	}

	auths := rds.ProxyAuthArray{}
	for _, auth := range spec.Auth {
		authArgs := &rds.ProxyAuthArgs{
			// SECRETS is the only auth scheme AWS supports - pinned
			// here, never spec surface.
			AuthScheme: pulumi.String("SECRETS"),
			SecretArn:  pulumi.String(auth.SecretArn.GetValue()),
		}
		if auth.Description != "" {
			authArgs.Description = pulumi.String(auth.Description)
		}
		if auth.IamAuth != "" {
			authArgs.IamAuth = pulumi.String(auth.IamAuth)
		}
		if auth.ClientPasswordAuthType != "" {
			authArgs.ClientPasswordAuthType = pulumi.String(auth.ClientPasswordAuthType)
		}
		if auth.Username != "" {
			authArgs.Username = pulumi.String(auth.Username)
		}
		auths = append(auths, authArgs)
	}

	proxyArgs := &rds.ProxyArgs{
		Name:         pulumi.String(locals.Target.Metadata.Name),
		EngineFamily: pulumi.String(spec.EngineFamily),
		RoleArn:      pulumi.String(spec.RoleArn.GetValue()),
		VpcSubnetIds: pulumi.ToStringArray(subnetIds),
		Auths:        auths,
		Tags:         pulumi.ToStringMap(locals.AwsTags),
	}
	if len(securityGroupIds) > 0 {
		proxyArgs.VpcSecurityGroupIds = pulumi.ToStringArray(securityGroupIds)
	}
	if spec.RequireTls {
		proxyArgs.RequireTls = pulumi.Bool(true)
	}
	if spec.IdleClientTimeout > 0 {
		proxyArgs.IdleClientTimeout = pulumi.Int(int(spec.IdleClientTimeout))
	}
	if spec.DebugLogging {
		proxyArgs.DebugLogging = pulumi.Bool(true)
	}
	if spec.DefaultAuthScheme != "" {
		proxyArgs.DefaultAuthScheme = pulumi.String(spec.DefaultAuthScheme)
	}
	if spec.EndpointNetworkType != "" {
		proxyArgs.EndpointNetworkType = pulumi.String(spec.EndpointNetworkType)
	}
	if spec.TargetConnectionNetworkType != "" {
		proxyArgs.TargetConnectionNetworkType = pulumi.String(spec.TargetConnectionNetworkType)
	}

	createdProxy, err := rds.NewProxy(ctx, "proxy", proxyArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create proxy")
	}

	// Always managed: with no pool tuning configured AWS keeps its
	// defaults, and managing the group gives the target a concrete
	// ordering point plus the ARN output either way.
	targetGroupArgs := &rds.ProxyDefaultTargetGroupArgs{
		DbProxyName: createdProxy.Name,
	}
	if pool := spec.ConnectionPool; pool != nil {
		poolArgs := &rds.ProxyDefaultTargetGroupConnectionPoolConfigArgs{}
		if pool.ConnectionBorrowTimeout != nil {
			poolArgs.ConnectionBorrowTimeout = pulumi.Int(int(*pool.ConnectionBorrowTimeout))
		}
		if pool.InitQuery != "" {
			poolArgs.InitQuery = pulumi.String(pool.InitQuery)
		}
		if pool.MaxConnectionsPercent != nil {
			poolArgs.MaxConnectionsPercent = pulumi.Int(int(*pool.MaxConnectionsPercent))
		}
		if pool.MaxIdleConnectionsPercent != nil {
			poolArgs.MaxIdleConnectionsPercent = pulumi.Int(int(*pool.MaxIdleConnectionsPercent))
		}
		if len(pool.SessionPinningFilters) > 0 {
			poolArgs.SessionPinningFilters = pulumi.ToStringArray(pool.SessionPinningFilters)
		}
		targetGroupArgs.ConnectionPoolConfig = poolArgs
	}

	createdTargetGroup, err := rds.NewProxyDefaultTargetGroup(ctx, "default-target-group",
		targetGroupArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "configure default target group")
	}

	endpointAddresses := pulumi.StringMap{}
	endpointArns := pulumi.StringMap{}
	for _, endpoint := range spec.Endpoints {
		endpointSubnetIds := make([]string, 0, len(endpoint.VpcSubnetIds))
		for _, subnetId := range endpoint.VpcSubnetIds {
			endpointSubnetIds = append(endpointSubnetIds, subnetId.GetValue())
		}
		endpointSecurityGroupIds := make([]string, 0, len(endpoint.VpcSecurityGroupIds))
		for _, securityGroupId := range endpoint.VpcSecurityGroupIds {
			endpointSecurityGroupIds = append(endpointSecurityGroupIds, securityGroupId.GetValue())
		}
		endpointArgs := &rds.ProxyEndpointArgs{
			DbProxyName:         createdProxy.Name,
			DbProxyEndpointName: pulumi.String(endpoint.Name),
			VpcSubnetIds:        pulumi.ToStringArray(endpointSubnetIds),
			Tags:                pulumi.ToStringMap(locals.AwsTags),
		}
		if endpoint.TargetRole != "" {
			endpointArgs.TargetRole = pulumi.String(endpoint.TargetRole)
		}
		if len(endpointSecurityGroupIds) > 0 {
			endpointArgs.VpcSecurityGroupIds = pulumi.ToStringArray(endpointSecurityGroupIds)
		}
		createdEndpoint, err := rds.NewProxyEndpoint(ctx,
			fmt.Sprintf("endpoint-%s", endpoint.Name),
			endpointArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create endpoint %s", endpoint.Name)
		}
		endpointAddresses[endpoint.Name] = createdEndpoint.Endpoint
		endpointArns[endpoint.Name] = createdEndpoint.Arn
	}

	targetType := pulumi.String("").ToStringOutput()
	targetRdsResourceId := pulumi.String("").ToStringOutput()
	if target := spec.Target; target != nil {
		targetArgs := &rds.ProxyTargetArgs{
			DbProxyName:     createdProxy.Name,
			TargetGroupName: createdTargetGroup.Name,
		}
		if target.DbInstanceIdentifier != nil {
			targetArgs.DbInstanceIdentifier = pulumi.String(target.DbInstanceIdentifier.GetValue())
		}
		if target.DbClusterIdentifier != nil {
			targetArgs.DbClusterIdentifier = pulumi.String(target.DbClusterIdentifier.GetValue())
		}
		createdTarget, err := rds.NewProxyTarget(ctx, "target", targetArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "register target")
		}
		targetType = createdTarget.Type
		targetRdsResourceId = createdTarget.RdsResourceId
	}

	ctx.Export(OpProxyName, createdProxy.Name)
	ctx.Export(OpProxyArn, createdProxy.Arn)
	ctx.Export(OpEndpoint, createdProxy.Endpoint)
	ctx.Export(OpDefaultTargetGroupArn, createdTargetGroup.Arn)
	ctx.Export(OpDefaultTargetGroupName, createdTargetGroup.Name)
	ctx.Export(OpEndpointAddresses, endpointAddresses)
	ctx.Export(OpEndpointArns, endpointArns)
	ctx.Export(OpTargetType, targetType)
	ctx.Export(OpTargetRdsResourceId, targetRdsResourceId)
	return nil
}
