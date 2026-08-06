package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/elasticache"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serverlessCache provisions the aws_elasticache_serverless_cache itself —
// a fully managed, auto-scaling cache with no node types or replica counts.
// The cache composes onto its neighbors instead of embedding them: subnets,
// security groups, and KMS keys attach by reference, and client ingress
// rules live on the referenced AwsSecurityGroup nodes — this module never
// creates or mutates a resource that deserves to be its own node.
//
// Create-only in AWS: the cache name, engine (when switching to/from
// memcached), subnet IDs, KMS key, network_type, and snapshot restore
// sources. Everything else updates in place.
func serverlessCache(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &elasticache.ServerlessCacheArgs{
		Engine: pulumi.String(spec.Engine),
		Name:   pulumi.String(locals.CacheName),
		Tags:   pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	if spec.MajorEngineVersion != "" {
		args.MajorEngineVersion = pulumi.String(spec.MajorEngineVersion)
	}

	// -------------------------------------------------------------------
	// Cache usage limits (flattened in spec, nested in AWS)
	// -------------------------------------------------------------------

	hasDataStorage := spec.DataStorageMinGb > 0 || spec.DataStorageMaxGb > 0
	hasEcpu := spec.EcpuMin > 0 || spec.EcpuMax > 0

	if hasDataStorage || hasEcpu {
		limits := &elasticache.ServerlessCacheCacheUsageLimitsArgs{}

		if hasDataStorage {
			ds := &elasticache.ServerlessCacheCacheUsageLimitsDataStorageArgs{
				Unit: pulumi.String("GB"),
			}
			if spec.DataStorageMinGb > 0 {
				ds.Minimum = pulumi.Int(int(spec.DataStorageMinGb))
			}
			if spec.DataStorageMaxGb > 0 {
				ds.Maximum = pulumi.Int(int(spec.DataStorageMaxGb))
			}
			limits.DataStorage = ds
		}

		if hasEcpu {
			ecpu := &elasticache.ServerlessCacheCacheUsageLimitsEcpuPerSecondArgs{}
			if spec.EcpuMin > 0 {
				ecpu.Minimum = pulumi.Int(int(spec.EcpuMin))
			}
			if spec.EcpuMax > 0 {
				ecpu.Maximum = pulumi.Int(int(spec.EcpuMax))
			}
			limits.EcpuPerSeconds = elasticache.ServerlessCacheCacheUsageLimitsEcpuPerSecondArray{ecpu}
		}

		args.CacheUsageLimits = limits
	}

	// -------------------------------------------------------------------
	// Networking
	// -------------------------------------------------------------------

	subnetIds := pulumi.StringArray{}
	for _, subnetId := range spec.SubnetIds {
		if subnetId.GetValue() != "" {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
	}
	if len(subnetIds) > 0 {
		args.SubnetIds = subnetIds
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

	// -------------------------------------------------------------------
	// Encryption
	// -------------------------------------------------------------------

	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}

	// -------------------------------------------------------------------
	// Snapshots (Redis/Valkey only — CEL guards prevent Memcached usage)
	// -------------------------------------------------------------------

	if spec.DailySnapshotTime != "" {
		args.DailySnapshotTime = pulumi.String(spec.DailySnapshotTime)
	}

	if spec.SnapshotRetentionLimit > 0 {
		args.SnapshotRetentionLimit = pulumi.Int(int(spec.SnapshotRetentionLimit))
	}

	if len(spec.SnapshotArnsToRestore) > 0 {
		args.SnapshotArnsToRestores = pulumi.ToStringArray(spec.SnapshotArnsToRestore)
	}

	// -------------------------------------------------------------------
	// Authentication (Redis/Valkey only — CEL guards prevent Memcached usage)
	// -------------------------------------------------------------------

	if spec.UserGroupId.GetValue() != "" {
		args.UserGroupId = pulumi.String(spec.UserGroupId.GetValue())
	}

	createdCache, err := elasticache.NewServerlessCache(ctx, "serverless-cache", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create serverless cache")
	}

	ctx.Export(OpArn, createdCache.Arn)
	ctx.Export(OpFullEngineVersion, createdCache.FullEngineVersion)
	ctx.Export(OpName, createdCache.Name)

	ctx.Export(OpEndpointAddress, createdCache.Endpoints.Index(pulumi.Int(0)).Address())
	ctx.Export(OpEndpointPort, createdCache.Endpoints.Index(pulumi.Int(0)).Port())
	ctx.Export(OpReaderEndpointAddress, createdCache.ReaderEndpoints.Index(pulumi.Int(0)).Address())
	ctx.Export(OpReaderEndpointPort, createdCache.ReaderEndpoints.Index(pulumi.Int(0)).Port())

	return nil
}
