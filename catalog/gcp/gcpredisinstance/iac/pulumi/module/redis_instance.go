package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/redis"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// redisInstance provisions a Memorystore for Redis instance — the classic
// VPC-peered managed Redis. One resource is one instance: a standalone BASIC
// node, or a STANDARD_HA primary + failover replica (optionally with read
// replicas).
//
// Lifecycle notes the API enforces:
//   - name, tier, connect_mode, transit_encryption_mode, authorized_network,
//     reserved_ip_range, location_id, alternative_location_id, and the CMEK
//     key are immutable — changing any of them replaces the instance.
//   - memory_size_gb resizes in place; redis_version upgrades in place but a
//     downgrade replaces the instance.
//   - read_replicas_mode is set at creation; replica_count and
//     secondary_ip_range are the in-place scale-out levers afterwards.
//   - with connect_mode PRIVATE_SERVICE_ACCESS, the network must already
//     carry a service networking connection — the API rejects the create
//     otherwise.
func redisInstance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpRedisInstance.Spec

	// Enable the Memorystore for Redis API first so a fresh project works on
	// the first deploy. disable_on_destroy stays false: tearing down one
	// instance must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("redis.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"redis-redis.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable redis.googleapis.com api")
	}

	// Destroy guard, sent explicitly from the spec (default true) so
	// destroy behavior is identical on both engines — omitting it would
	// leave the decision to whatever happens to be in Terraform state.
	deletionProtection := true
	if spec.DeletionProtection != nil {
		deletionProtection = spec.GetDeletionProtection()
	}

	args := &redis.InstanceArgs{
		Name:         pulumi.String(spec.InstanceName),
		Region:       pulumi.StringPtr(spec.Region),
		Tier:         pulumi.StringPtr(spec.Tier),
		MemorySizeGb: pulumi.Int(int(spec.MemorySizeGb)),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),

		DeletionProtection: pulumi.BoolPtr(deletionProtection),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Engine version. Upgrades apply in place; a downgrade replaces the
	// instance.
	if spec.RedisVersion != "" {
		args.RedisVersion = pulumi.StringPtr(spec.RedisVersion)
	}

	if spec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}

	// Zone placement: primary zone, and (STANDARD_HA only) the replica zone.
	// Left unset, GCP spreads the nodes across zones automatically.
	if spec.LocationId != "" {
		args.LocationId = pulumi.StringPtr(spec.LocationId)
	}
	if spec.AlternativeLocationId != "" {
		args.AlternativeLocationId = pulumi.StringPtr(spec.AlternativeLocationId)
	}

	// Connectivity is fixed at creation: the peered network, how it peers
	// (direct peering vs the network's private services access connection),
	// and which internal range the nodes occupy.
	if spec.AuthorizedNetwork.GetValue() != "" {
		args.AuthorizedNetwork = pulumi.StringPtr(spec.AuthorizedNetwork.GetValue())
	}
	if spec.ConnectMode != "" {
		args.ConnectMode = pulumi.StringPtr(spec.ConnectMode)
	}
	if spec.ReservedIpRange != "" {
		args.ReservedIpRange = pulumi.StringPtr(spec.ReservedIpRange)
	}

	// The scale-out range: adding read replicas to an existing instance
	// needs more address space than the original /29 provides.
	if spec.SecondaryIpRange != "" {
		args.SecondaryIpRange = pulumi.StringPtr(spec.SecondaryIpRange)
	}

	// Redis AUTH: when true GCP generates and rotates the AUTH string
	// (exported as a secret output below).
	if spec.AuthEnabled {
		args.AuthEnabled = pulumi.BoolPtr(true)
	}

	if spec.TransitEncryptionMode != "" {
		args.TransitEncryptionMode = pulumi.StringPtr(spec.TransitEncryptionMode)
	}

	if len(spec.RedisConfigs) > 0 {
		args.RedisConfigs = pulumi.ToStringMap(spec.RedisConfigs)
	}

	// Weekly maintenance window (UTC, fixed 1-hour duration).
	if spec.MaintenanceWindow != nil {
		policyArgs := &redis.InstanceMaintenancePolicyArgs{
			WeeklyMaintenanceWindows: redis.InstanceMaintenancePolicyWeeklyMaintenanceWindowArray{
				&redis.InstanceMaintenancePolicyWeeklyMaintenanceWindowArgs{
					Day: pulumi.String(spec.MaintenanceWindow.Day),
					StartTime: &redis.InstanceMaintenancePolicyWeeklyMaintenanceWindowStartTimeArgs{
						Hours:   pulumi.IntPtr(int(spec.MaintenanceWindow.Hour)),
						Minutes: pulumi.IntPtr(int(spec.MaintenanceWindow.Minute)),
					},
				},
			},
		}
		if spec.MaintenanceWindow.Description != "" {
			policyArgs.Description = pulumi.StringPtr(spec.MaintenanceWindow.Description)
		}
		args.MaintenancePolicy = policyArgs
	}

	// Self-service maintenance: setting a newer available version applies
	// the update now instead of waiting for GCP's rollout.
	if spec.MaintenanceVersion != "" {
		args.MaintenanceVersion = pulumi.StringPtr(spec.MaintenanceVersion)
	}

	// Read replicas (STANDARD_HA only; mode is fixed at creation).
	if spec.ReadReplicasMode != "" {
		args.ReadReplicasMode = pulumi.StringPtr(spec.ReadReplicasMode)
	}
	if spec.ReplicaCount > 0 {
		args.ReplicaCount = pulumi.IntPtr(int(spec.ReplicaCount))
	}

	// RDB snapshot persistence.
	if spec.PersistenceConfig != nil {
		persistenceArgs := &redis.InstancePersistenceConfigArgs{
			PersistenceMode: pulumi.StringPtr(spec.PersistenceConfig.PersistenceMode),
		}
		if spec.PersistenceConfig.RdbSnapshotPeriod != "" {
			persistenceArgs.RdbSnapshotPeriod = pulumi.StringPtr(spec.PersistenceConfig.RdbSnapshotPeriod)
		}
		if spec.PersistenceConfig.RdbSnapshotStartTime != "" {
			persistenceArgs.RdbSnapshotStartTime = pulumi.StringPtr(spec.PersistenceConfig.RdbSnapshotStartTime)
		}
		args.PersistenceConfig = persistenceArgs
	}

	// CMEK: the key must live in the instance's region. Immutable.
	if spec.CustomerManagedKey.GetValue() != "" {
		args.CustomerManagedKey = pulumi.StringPtr(spec.CustomerManagedKey.GetValue())
	}

	// Client-side destroy behavior: DELETE (default), PREVENT, or ABANDON.
	// Sent only when set so the provider default stays in charge otherwise.
	// Evaluated only after deletion_protection allows the destroy at all.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdInstance, err := redis.NewInstance(ctx, "redis-instance", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create redis instance")
	}

	ctx.Export(OpHost, createdInstance.Host)
	ctx.Export(OpPort, createdInstance.Port)
	ctx.Export(OpReadEndpoint, createdInstance.ReadEndpoint)
	ctx.Export(OpReadEndpointPort, createdInstance.ReadEndpointPort)
	ctx.Export(OpCurrentLocationId, createdInstance.CurrentLocationId)
	// GCP-generated AUTH credential — kept encrypted in Pulumi state.
	ctx.Export(OpAuthString, pulumi.ToSecret(createdInstance.AuthString))
	// The PEM contents are what clients install as trust anchors; the other
	// certificate attributes (serial, fingerprint, expiry) stay internal.
	ctx.Export(OpServerCaCerts, createdInstance.ServerCaCerts.ApplyT(
		func(certs []redis.InstanceServerCaCert) []string {
			pems := make([]string, 0, len(certs))
			for _, cert := range certs {
				if cert.Cert != nil {
					pems = append(pems, *cert.Cert)
				}
			}
			return pems
		}).(pulumi.StringArrayOutput))
	ctx.Export(OpPersistenceIamIdentity, createdInstance.PersistenceIamIdentity)
	ctx.Export(OpEffectiveReservedIpRange, createdInstance.EffectiveReservedIpRange)
	ctx.Export(OpInstanceName, createdInstance.Name)
	// The plain spec region name (not the provider's computed attribute) so
	// both engines emit the identical value for API callers.
	ctx.Export(OpRegion, pulumi.String(spec.Region))

	return nil
}
