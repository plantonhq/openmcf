package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/alloydb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func primaryInstance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, createdCluster *alloydb.Cluster) error {
	spec := locals.GcpAlloydbCluster.Spec
	instanceSpec := spec.PrimaryInstance

	args := &alloydb.InstanceArgs{
		Cluster:      createdCluster.Name,
		InstanceId:   pulumi.String(instanceSpec.InstanceId),
		InstanceType: pulumi.String("PRIMARY"),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),
	}

	// Machine configuration: cpu_count or machine_type (mutually exclusive, validated by proto).
	if instanceSpec.CpuCount > 0 || instanceSpec.MachineType != "" {
		machineConfig := &alloydb.InstanceMachineConfigArgs{}
		if instanceSpec.CpuCount > 0 {
			machineConfig.CpuCount = pulumi.IntPtr(int(instanceSpec.CpuCount))
		}
		if instanceSpec.MachineType != "" {
			machineConfig.MachineType = pulumi.StringPtr(instanceSpec.MachineType)
		}
		args.MachineConfig = machineConfig
	}

	// Availability type.
	if instanceSpec.AvailabilityType != "" {
		args.AvailabilityType = pulumi.StringPtr(instanceSpec.AvailabilityType)
	}

	// Database flags.
	if len(instanceSpec.DatabaseFlags) > 0 {
		args.DatabaseFlags = pulumi.ToStringMap(instanceSpec.DatabaseFlags)
	}

	// Display name.
	if instanceSpec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(instanceSpec.DisplayName)
	}

	// Query insights configuration.
	if instanceSpec.QueryInsightsConfig != nil {
		qiConfig := instanceSpec.QueryInsightsConfig
		qiArgs := &alloydb.InstanceQueryInsightsConfigArgs{}

		if qiConfig.QueryPlansPerMinute > 0 {
			qiArgs.QueryPlansPerMinute = pulumi.IntPtr(int(qiConfig.QueryPlansPerMinute))
		}
		if qiConfig.QueryStringLength > 0 {
			qiArgs.QueryStringLength = pulumi.IntPtr(int(qiConfig.QueryStringLength))
		}
		if qiConfig.RecordApplicationTags {
			qiArgs.RecordApplicationTags = pulumi.BoolPtr(true)
		}
		if qiConfig.RecordClientAddress {
			qiArgs.RecordClientAddress = pulumi.BoolPtr(true)
		}
		args.QueryInsightsConfig = qiArgs
	}

	// Client connection configuration (require_connectors, ssl_mode).
	if instanceSpec.RequireConnectors || instanceSpec.SslMode != "" {
		clientConfig := &alloydb.InstanceClientConnectionConfigArgs{}
		if instanceSpec.RequireConnectors {
			clientConfig.RequireConnectors = pulumi.BoolPtr(true)
		}
		if instanceSpec.SslMode != "" {
			clientConfig.SslConfig = &alloydb.InstanceClientConnectionConfigSslConfigArgs{
				SslMode: pulumi.StringPtr(instanceSpec.SslMode),
			}
		}
		args.ClientConnectionConfig = clientConfig
	}

	// Stop/start lever: NEVER stops the primary's compute (storage and
	// configuration survive); ALWAYS restarts it.
	if instanceSpec.ActivationPolicy != "" {
		args.ActivationPolicy = pulumi.StringPtr(instanceSpec.ActivationPolicy)
	}

	// Client tool metadata, paired with the computed effective_annotations.
	if len(instanceSpec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(instanceSpec.Annotations)
	}

	// ZONAL primaries only — GCP rejects it on REGIONAL instances
	// (spec-enforced pairing). Changing it live-migrates the primary.
	if instanceSpec.GceZone != "" {
		args.GceZone = pulumi.StringPtr(instanceSpec.GceZone)
	}

	// AlloyDB managed connection pooling (built-in pooler).
	if instanceSpec.ConnectionPoolConfig != nil {
		poolArgs := &alloydb.InstanceConnectionPoolConfigArgs{
			Enabled: pulumi.Bool(instanceSpec.ConnectionPoolConfig.Enabled),
		}
		if len(instanceSpec.ConnectionPoolConfig.Flags) > 0 {
			poolArgs.Flags = pulumi.ToStringMap(instanceSpec.ConnectionPoolConfig.Flags)
		}
		args.ConnectionPoolConfig = poolArgs
	}

	// Public-IP / PSA-range surface on the bundled primary — the same
	// contract the standalone GcpAlloydbInstance kind models.
	if instanceSpec.EnablePublicIp || instanceSpec.EnableOutboundPublicIp || len(instanceSpec.AuthorizedExternalNetworks) > 0 || instanceSpec.AllocatedIpRangeOverride != "" {
		netConfig := &alloydb.InstanceNetworkConfigArgs{}
		if instanceSpec.EnablePublicIp {
			netConfig.EnablePublicIp = pulumi.BoolPtr(true)
		}
		if instanceSpec.EnableOutboundPublicIp {
			netConfig.EnableOutboundPublicIp = pulumi.BoolPtr(true)
		}
		// Immutable: a different PSA range recreates the primary.
		if instanceSpec.AllocatedIpRangeOverride != "" {
			netConfig.AllocatedIpRangeOverride = pulumi.StringPtr(instanceSpec.AllocatedIpRangeOverride)
		}
		if len(instanceSpec.AuthorizedExternalNetworks) > 0 {
			nets := make(alloydb.InstanceNetworkConfigAuthorizedExternalNetworkArray, len(instanceSpec.AuthorizedExternalNetworks))
			for i, n := range instanceSpec.AuthorizedExternalNetworks {
				nets[i] = &alloydb.InstanceNetworkConfigAuthorizedExternalNetworkArgs{
					CidrRange: pulumi.StringPtr(n.CidrRange),
				}
			}
			netConfig.AuthorizedExternalNetworks = nets
		}
		args.NetworkConfig = netConfig
	}

	// PSC on the bundled primary (meaningful only on PSC clusters).
	if instanceSpec.PscInstanceConfig != nil {
		psc := instanceSpec.PscInstanceConfig
		pscArgs := &alloydb.InstancePscInstanceConfigArgs{}
		if len(psc.AllowedConsumerProjects) > 0 {
			pscArgs.AllowedConsumerProjects = pulumi.ToStringArray(psc.AllowedConsumerProjects)
		}
		if len(psc.PscAutoConnections) > 0 {
			conns := make(alloydb.InstancePscInstanceConfigPscAutoConnectionArray, len(psc.PscAutoConnections))
			for i, c := range psc.PscAutoConnections {
				connArgs := &alloydb.InstancePscInstanceConfigPscAutoConnectionArgs{}
				if c.ConsumerNetwork != "" {
					connArgs.ConsumerNetwork = pulumi.StringPtr(c.ConsumerNetwork)
				}
				if c.ConsumerProject != "" {
					connArgs.ConsumerProject = pulumi.StringPtr(c.ConsumerProject)
				}
				conns[i] = connArgs
			}
			pscArgs.PscAutoConnections = conns
		}
		if len(psc.PscInterfaceConfigs) > 0 {
			ifaces := make(alloydb.InstancePscInstanceConfigPscInterfaceConfigArray, len(psc.PscInterfaceConfigs))
			for i, iface := range psc.PscInterfaceConfigs {
				ifaces[i] = &alloydb.InstancePscInstanceConfigPscInterfaceConfigArgs{
					NetworkAttachmentResource: pulumi.StringPtr(iface.NetworkAttachmentResource),
				}
			}
			pscArgs.PscInterfaceConfigs = ifaces
		}
		args.PscInstanceConfig = pscArgs
	}

	// The PRIMARY INSTANCE's own destroy behavior (the cluster's
	// deletion_policy is separate — AlloyDB gives each resource its own).
	if instanceSpec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(instanceSpec.DeletionPolicy)
	}

	createdInstance, err := alloydb.NewInstance(ctx, "alloydb-primary-instance", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdCluster}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create alloydb primary instance")
	}

	// Export primary instance outputs.
	ctx.Export(OpPrimaryInstanceIp, createdInstance.IpAddress)
	ctx.Export(OpPrimaryInstanceName, createdInstance.Name)

	return nil
}
