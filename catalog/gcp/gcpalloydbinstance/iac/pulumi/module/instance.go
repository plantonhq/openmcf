package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/alloydb"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// instance provisions an AlloyDB instance (`google_alloydb_instance`) on an
// existing cluster — typically a READ_POOL for read scaling, but PRIMARY and
// SECONDARY types are supported for advanced topologies.
func instance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpAlloydbInstance.Spec

	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("alloydb.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"alloydb-alloydb.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable alloydb.googleapis.com api")
	}

	args := &alloydb.InstanceArgs{
		Cluster:      pulumi.String(spec.Cluster.GetValue()),
		InstanceId:   pulumi.String(spec.InstanceId),
		InstanceType: pulumi.String(spec.GetInstanceType()),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),
	}

	if spec.CpuCount > 0 || spec.MachineType != "" {
		machineConfig := &alloydb.InstanceMachineConfigArgs{}
		if spec.CpuCount > 0 {
			machineConfig.CpuCount = pulumi.IntPtr(int(spec.CpuCount))
		}
		if spec.MachineType != "" {
			machineConfig.MachineType = pulumi.StringPtr(spec.MachineType)
		}
		args.MachineConfig = machineConfig
	}

	if spec.ReadPoolConfig != nil && spec.ReadPoolConfig.NodeCount > 0 {
		args.ReadPoolConfig = &alloydb.InstanceReadPoolConfigArgs{
			NodeCount: pulumi.IntPtr(int(spec.ReadPoolConfig.NodeCount)),
		}
	}

	// Only ever non-empty for PRIMARY/SECONDARY instances (spec CEL): read
	// pools derive availability from node_count and the API drops a sent
	// value, which would refresh dirty forever.
	if spec.AvailabilityType != "" {
		args.AvailabilityType = pulumi.StringPtr(spec.AvailabilityType)
	}

	if len(spec.DatabaseFlags) > 0 {
		args.DatabaseFlags = pulumi.ToStringMap(spec.DatabaseFlags)
	}

	if spec.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}

	if spec.QueryInsightsConfig != nil {
		qi := spec.QueryInsightsConfig
		qiArgs := &alloydb.InstanceQueryInsightsConfigArgs{}
		if qi.QueryPlansPerMinute > 0 {
			qiArgs.QueryPlansPerMinute = pulumi.IntPtr(int(qi.QueryPlansPerMinute))
		}
		if qi.QueryStringLength > 0 {
			qiArgs.QueryStringLength = pulumi.IntPtr(int(qi.QueryStringLength))
		}
		if qi.RecordApplicationTags {
			qiArgs.RecordApplicationTags = pulumi.BoolPtr(true)
		}
		if qi.RecordClientAddress {
			qiArgs.RecordClientAddress = pulumi.BoolPtr(true)
		}
		args.QueryInsightsConfig = qiArgs
	}

	if spec.RequireConnectors || spec.SslMode != "" {
		clientConfig := &alloydb.InstanceClientConnectionConfigArgs{}
		if spec.RequireConnectors {
			clientConfig.RequireConnectors = pulumi.BoolPtr(true)
		}
		if spec.SslMode != "" {
			clientConfig.SslConfig = &alloydb.InstanceClientConnectionConfigSslConfigArgs{
				SslMode: pulumi.StringPtr(spec.SslMode),
			}
		}
		args.ClientConnectionConfig = clientConfig
	}

	if spec.ActivationPolicy != "" {
		args.ActivationPolicy = pulumi.StringPtr(spec.ActivationPolicy)
	}

	// Client tool metadata, paired with the computed effective_annotations.
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(spec.Annotations)
	}

	// ZONAL instances only — GCP rejects it on REGIONAL instances
	// (spec-enforced pairing). Changing it live-migrates the instance.
	if spec.GceZone != "" {
		args.GceZone = pulumi.StringPtr(spec.GceZone)
	}

	// AlloyDB managed connection pooling (built-in pooler). Flags only
	// apply while enabled is true.
	if spec.ConnectionPoolConfig != nil {
		poolArgs := &alloydb.InstanceConnectionPoolConfigArgs{
			Enabled: pulumi.Bool(spec.ConnectionPoolConfig.Enabled),
		}
		if len(spec.ConnectionPoolConfig.Flags) > 0 {
			poolArgs.Flags = pulumi.ToStringMap(spec.ConnectionPoolConfig.Flags)
		}
		args.ConnectionPoolConfig = poolArgs
	}

	if spec.EnablePublicIp || spec.EnableOutboundPublicIp || len(spec.AuthorizedExternalNetworks) > 0 || spec.AllocatedIpRangeOverride != "" {
		netConfig := &alloydb.InstanceNetworkConfigArgs{}
		if spec.EnablePublicIp {
			netConfig.EnablePublicIp = pulumi.BoolPtr(true)
		}
		if spec.EnableOutboundPublicIp {
			netConfig.EnableOutboundPublicIp = pulumi.BoolPtr(true)
		}
		// Immutable: a different PSA range recreates the instance.
		if spec.AllocatedIpRangeOverride != "" {
			netConfig.AllocatedIpRangeOverride = pulumi.StringPtr(spec.AllocatedIpRangeOverride)
		}
		if len(spec.AuthorizedExternalNetworks) > 0 {
			nets := make(alloydb.InstanceNetworkConfigAuthorizedExternalNetworkArray, len(spec.AuthorizedExternalNetworks))
			for i, n := range spec.AuthorizedExternalNetworks {
				nets[i] = &alloydb.InstanceNetworkConfigAuthorizedExternalNetworkArgs{
					CidrRange: pulumi.StringPtr(n.CidrRange),
				}
			}
			netConfig.AuthorizedExternalNetworks = nets
		}
		args.NetworkConfig = netConfig
	}

	if spec.PscInstanceConfig != nil {
		psc := spec.PscInstanceConfig
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

	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdInstance, err := alloydb.NewInstance(ctx, "alloydb-instance", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create alloydb instance")
	}

	ctx.Export(OpInstanceName, createdInstance.Name)
	ctx.Export(OpIpAddress, createdInstance.IpAddress)
	ctx.Export(OpState, createdInstance.State)

	return nil
}
