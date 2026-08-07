package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster provisions the aws_elasticache_cluster itself — a distributed
// Memcached cache with sub-millisecond latency. The cluster composes onto
// its neighbors instead of embedding them: subnets, security groups, and
// parameter groups attach by reference, and client ingress rules live on
// the referenced AwsSecurityGroup nodes — this module never creates or
// mutates a resource that deserves to be its own node.
//
// Create-only in AWS: the cluster identifier, engine, port, subnet group,
// network_type, and node type (vertical scaling forces recreation).
// Everything else updates in place (immediately or at the next maintenance
// window, per apply_immediately).
func cluster(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdSubnetGroup *elasticache.SubnetGroup,
	createdParamGroup *elasticache.ParameterGroup,
) error {
	spec := locals.Spec

	args := &elasticache.ClusterArgs{
		ClusterId:        pulumi.String(locals.ClusterIdentifier),
		Engine:           pulumi.String("memcached"),
		NodeType:         pulumi.String(spec.NodeType),
		NumCacheNodes:    pulumi.Int(spec.NumCacheNodes),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
		ApplyImmediately: pulumi.Bool(spec.ApplyImmediately),
	}

	// Empty pins nothing: AWS picks the engine's current default version,
	// so an unpinned manifest never goes stale.
	if spec.EngineVersion != "" {
		args.EngineVersion = pulumi.String(spec.EngineVersion)
	}

	if spec.Port != nil {
		args.Port = pulumi.Int(*spec.Port)
	}

	if spec.AzMode != "" {
		args.AzMode = pulumi.String(spec.AzMode)
	}

	if spec.TransitEncryptionEnabled {
		args.TransitEncryptionEnabled = pulumi.Bool(true)
	}

	// -------------------------------------------------------------------
	// Networking: the subnet group managed here (or referenced), security
	// groups by reference, and optional dual-stack addressing.
	// -------------------------------------------------------------------

	if createdSubnetGroup != nil {
		args.SubnetGroupName = createdSubnetGroup.Name
	} else if spec.SubnetGroupName != "" {
		args.SubnetGroupName = pulumi.String(spec.SubnetGroupName)
	}

	securityGroupIds := pulumi.StringArray{}
	for _, securityGroupId := range spec.SecurityGroupIds {
		if securityGroupId.GetValue() != "" {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
	}
	if len(securityGroupIds) > 0 {
		args.SecurityGroupIds = securityGroupIds
	}

	if spec.NetworkType != "" {
		args.NetworkType = pulumi.String(spec.NetworkType)
	}
	if spec.IpDiscovery != "" {
		args.IpDiscovery = pulumi.String(spec.IpDiscovery)
	}

	// -------------------------------------------------------------------
	// Parameter groups: the managed inline group, an existing referenced
	// group, or the engine default.
	// -------------------------------------------------------------------

	if createdParamGroup != nil {
		args.ParameterGroupName = createdParamGroup.Name
	} else if spec.ParameterGroupName != "" {
		args.ParameterGroupName = pulumi.String(spec.ParameterGroupName)
	}

	// -------------------------------------------------------------------
	// Maintenance
	// -------------------------------------------------------------------

	if spec.MaintenanceWindow != "" {
		args.MaintenanceWindow = pulumi.String(spec.MaintenanceWindow)
	}
	if spec.AutoMinorVersionUpgrade {
		args.AutoMinorVersionUpgrade = pulumi.String("true")
	}

	// -------------------------------------------------------------------
	// Notifications
	// -------------------------------------------------------------------

	if spec.NotificationTopicArn.GetValue() != "" {
		args.NotificationTopicArn = pulumi.String(spec.NotificationTopicArn.GetValue())
	}

	// -------------------------------------------------------------------
	// Node placement
	// -------------------------------------------------------------------

	if len(spec.PreferredAvailabilityZones) > 0 {
		args.PreferredAvailabilityZones = pulumi.ToStringArray(spec.PreferredAvailabilityZones)
	}

	createdCluster, err := elasticache.NewCluster(ctx, "cluster", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create memcached cluster")
	}

	ctx.Export(OpClusterId, createdCluster.ClusterId)
	ctx.Export(OpClusterAddress, createdCluster.ClusterAddress)
	ctx.Export(OpConfigEndpoint, createdCluster.ConfigurationEndpoint)
	ctx.Export(OpArn, createdCluster.Arn)
	ctx.Export(OpPort, createdCluster.Port)

	return nil
}
