package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// replicationGroup provisions the aws_elasticache_replication_group
// itself -- the in-memory data plane: endpoints, topology, encryption,
// authentication, and engine lifecycle. The group composes onto its
// neighbors instead of embedding them: subnets, security groups, KMS
// keys, and RBAC user groups attach by reference, and network ingress
// rules live on the referenced AwsSecurityGroup nodes -- this module
// never creates or mutates a resource that deserves to be its own node.
//
// Create-only in AWS: the replication_group_id, port, network_type,
// at-rest encryption + KMS key, restore sources (snapshot_arns/
// snapshot_name), global datastore membership, durability, and explicit
// shard placement (node_group_configurations). Everything else updates
// in place (immediately or at the next maintenance window, per
// apply_immediately).
func replicationGroup(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdSubnetGroup *elasticache.SubnetGroup,
	createdParameterGroup *elasticache.ParameterGroup,
) (*elasticache.ReplicationGroup, error) {
	spec := locals.AwsRedisElasticache.Spec

	args := &elasticache.ReplicationGroupArgs{
		ReplicationGroupId: pulumi.String(locals.ReplicationGroupId),
		Description:        pulumi.String(spec.Description),
		Tags:               pulumi.ToStringMap(locals.AwsTags),
		ApplyImmediately:   pulumi.Bool(spec.ApplyImmediately),
	}

	// Engine settings are inherited from the global primary when joining
	// a global datastore -- forward only when the spec carries them.
	if spec.Engine != "" {
		args.Engine = pulumi.String(spec.Engine)
	}
	if spec.EngineVersion != "" {
		args.EngineVersion = pulumi.String(spec.EngineVersion)
	}
	if spec.NodeType != "" {
		args.NodeType = pulumi.String(spec.NodeType)
	}
	if spec.Port != nil {
		args.Port = pulumi.Int(int(*spec.Port))
	}

	// Topology: exactly one of num_cache_clusters (non-clustered) or
	// num_node_groups (clustered), CEL-enforced. A global-datastore
	// secondary may set only num_cache_clusters.
	if spec.NumCacheClusters > 0 {
		args.NumCacheClusters = pulumi.Int(int(spec.NumCacheClusters))
	} else if spec.NumNodeGroups > 0 {
		args.NumNodeGroups = pulumi.Int(int(spec.NumNodeGroups))
		if spec.ReplicasPerNodeGroup > 0 {
			args.ReplicasPerNodeGroup = pulumi.Int(int(spec.ReplicasPerNodeGroup))
		}
	}

	if len(spec.PreferredCacheClusterAzs) > 0 {
		args.PreferredCacheClusterAzs = pulumi.ToStringArray(spec.PreferredCacheClusterAzs)
	}

	if len(spec.NodeGroupConfigurations) > 0 {
		nodeGroupConfigs := elasticache.ReplicationGroupNodeGroupConfigurationArray{}
		for _, nodeGroupConfig := range spec.NodeGroupConfigurations {
			configArgs := &elasticache.ReplicationGroupNodeGroupConfigurationArgs{}
			if nodeGroupConfig.NodeGroupId != "" {
				configArgs.NodeGroupId = pulumi.String(nodeGroupConfig.NodeGroupId)
			}
			if nodeGroupConfig.PrimaryAvailabilityZone != "" {
				configArgs.PrimaryAvailabilityZone = pulumi.String(nodeGroupConfig.PrimaryAvailabilityZone)
			}
			if len(nodeGroupConfig.ReplicaAvailabilityZones) > 0 {
				configArgs.ReplicaAvailabilityZones = pulumi.ToStringArray(nodeGroupConfig.ReplicaAvailabilityZones)
			}
			if nodeGroupConfig.ReplicaCount > 0 {
				configArgs.ReplicaCount = pulumi.Int(int(nodeGroupConfig.ReplicaCount))
			}
			if nodeGroupConfig.Slots != "" {
				configArgs.Slots = pulumi.String(nodeGroupConfig.Slots)
			}
			nodeGroupConfigs = append(nodeGroupConfigs, configArgs)
		}
		args.NodeGroupConfigurations = nodeGroupConfigs
	}

	args.AutomaticFailoverEnabled = pulumi.Bool(spec.AutomaticFailoverEnabled)
	args.MultiAzEnabled = pulumi.Bool(spec.MultiAzEnabled)

	if spec.Durability != "" {
		args.Durability = pulumi.String(spec.Durability)
	}

	if spec.GlobalReplicationGroupId != "" {
		args.GlobalReplicationGroupId = pulumi.String(spec.GlobalReplicationGroupId)
	}

	// Networking: the subnet group managed here (or referenced), security
	// groups for node-level access, and optional dual-stack settings.
	if createdSubnetGroup != nil {
		args.SubnetGroupName = createdSubnetGroup.Name
	} else if spec.SubnetGroupName != "" {
		args.SubnetGroupName = pulumi.String(spec.SubnetGroupName)
	}
	if len(spec.SecurityGroupIds) > 0 {
		securityGroupIds := pulumi.StringArray{}
		for _, securityGroupId := range spec.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
		args.SecurityGroupIds = securityGroupIds
	}
	if spec.NetworkType != "" {
		args.NetworkType = pulumi.String(spec.NetworkType)
	}
	if spec.IpDiscovery != "" {
		args.IpDiscovery = pulumi.String(spec.IpDiscovery)
	}

	// Presence-typed enable flags: forward only when the manifest set them.
	// Unset must be OMITTED entirely — AWS applies its engine default, and
	// a global-datastore secondary MUST omit them (the provider conflicts
	// their presence with global_replication_group_id).
	if spec.AtRestEncryptionEnabled != nil {
		args.AtRestEncryptionEnabled = pulumi.BoolPtr(*spec.AtRestEncryptionEnabled)
	}
	if spec.TransitEncryptionEnabled != nil {
		args.TransitEncryptionEnabled = pulumi.BoolPtr(*spec.TransitEncryptionEnabled)
	}
	if spec.TransitEncryptionMode != "" {
		args.TransitEncryptionMode = pulumi.String(spec.TransitEncryptionMode)
	}
	if spec.KmsKeyId != nil && spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}

	// Authentication: legacy AUTH token or RBAC user groups (mutually
	// exclusive, CEL-enforced). auth_token is a presence field -- forward
	// only when the manifest explicitly sets it.
	if spec.AuthToken != nil {
		args.AuthToken = pulumi.String(spec.AuthToken.GetValue())
	}
	if spec.AuthTokenUpdateStrategy != "" {
		args.AuthTokenUpdateStrategy = pulumi.String(spec.AuthTokenUpdateStrategy)
	}
	if len(spec.UserGroupIds) > 0 {
		userGroupIds := pulumi.StringArray{}
		for _, userGroupId := range spec.UserGroupIds {
			userGroupIds = append(userGroupIds, pulumi.String(userGroupId.GetValue()))
		}
		args.UserGroupIds = userGroupIds
	}

	// Create-time restore sources (mutually exclusive, CEL-enforced):
	// from S3 RDB files or from an ElastiCache snapshot.
	if len(spec.SnapshotArns) > 0 {
		args.SnapshotArns = pulumi.ToStringArray(spec.SnapshotArns)
	}
	if spec.SnapshotName != "" {
		args.SnapshotName = pulumi.String(spec.SnapshotName)
	}

	if spec.MaintenanceWindow != "" {
		args.MaintenanceWindow = pulumi.String(spec.MaintenanceWindow)
	}
	if spec.SnapshotRetentionLimit > 0 {
		args.SnapshotRetentionLimit = pulumi.Int(int(spec.SnapshotRetentionLimit))
	}
	if spec.SnapshotWindow != "" {
		args.SnapshotWindow = pulumi.String(spec.SnapshotWindow)
	}
	if spec.FinalSnapshotIdentifier != "" {
		args.FinalSnapshotIdentifier = pulumi.String(spec.FinalSnapshotIdentifier)
	}

	// Parameter groups: the managed inline group, an existing referenced
	// group, or the engine default.
	if createdParameterGroup != nil {
		args.ParameterGroupName = createdParameterGroup.Name
	} else if spec.ParameterGroupName != "" {
		args.ParameterGroupName = pulumi.String(spec.ParameterGroupName)
	}

	if len(spec.LogDeliveryConfigurations) > 0 {
		logConfigs := elasticache.ReplicationGroupLogDeliveryConfigurationArray{}
		for _, logConfig := range spec.LogDeliveryConfigurations {
			logConfigs = append(logConfigs, &elasticache.ReplicationGroupLogDeliveryConfigurationArgs{
				DestinationType: pulumi.String(logConfig.DestinationType),
				Destination:     pulumi.String(logConfig.Destination.GetValue()),
				LogFormat:       pulumi.String(logConfig.LogFormat),
				LogType:         pulumi.String(logConfig.LogType),
			})
		}
		args.LogDeliveryConfigurations = logConfigs
	}

	if spec.NotificationTopicArn != nil && spec.NotificationTopicArn.GetValue() != "" {
		args.NotificationTopicArn = pulumi.String(spec.NotificationTopicArn.GetValue())
	}
	// auto_minor_version_upgrade is presence-typed: AWS enables minor
	// upgrades by default, so unset is omitted (AWS decides) and an
	// explicit false is forwarded — never conflated.
	if spec.AutoMinorVersionUpgrade != nil {
		args.AutoMinorVersionUpgrade = pulumi.BoolPtr(*spec.AutoMinorVersionUpgrade)
	}
	if spec.DataTieringEnabled {
		args.DataTieringEnabled = pulumi.Bool(true)
	}
	if spec.ClusterMode != "" {
		args.ClusterMode = pulumi.String(spec.ClusterMode)
	}

	createdReplicationGroup, err := elasticache.NewReplicationGroup(ctx, "replication-group", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ElastiCache replication group")
	}
	return createdReplicationGroup, nil
}
