package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/filestore"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// filestoreInstance enables the Filestore API and creates the instance.
//
// Sharp edges, all taught by the API rather than invented here:
//
//   - name, location, tier, protocol, network attachment, KMS key, and
//     replication are immutable — changing any of them replaces the
//     instance (and its data). File share capacity grows in place but
//     never shrinks.
//
//   - deletion_protection_enabled must be flipped false before a
//     protected instance can be destroyed.
//
//   - connect_mode defaults to DIRECT_PEERING; PRIVATE_SERVICE_ACCESS is
//     required for Shared VPC consumers and rides an existing
//     service-networking connection on the VPC.
func filestoreInstance(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpFilestoreInstance.Spec

	// Enable the Filestore API — the control plane that owns instances.
	// DisableOnDestroy stays false: tearing down one instance must never
	// disable the API for everything else in the project.
	fileApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("file.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		fileApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdFileApi, err := projects.NewService(ctx,
		"gcpnfs-file.googleapis.com", fileApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable file.googleapis.com api")
	}

	// Build NFS export options for the file share.
	var nfsExportOptions filestore.InstanceFileSharesNfsExportOptionArray
	for _, opt := range spec.FileShare.NfsExportOptions {
		exportOpt := &filestore.InstanceFileSharesNfsExportOptionArgs{}

		if len(opt.IpRanges) > 0 {
			exportOpt.IpRanges = pulumi.ToStringArray(opt.IpRanges)
		}
		if opt.AccessMode != "" {
			exportOpt.AccessMode = pulumi.StringPtr(opt.AccessMode)
		}
		if opt.SquashMode != "" {
			exportOpt.SquashMode = pulumi.StringPtr(opt.SquashMode)
		}
		if opt.AnonUid != nil {
			exportOpt.AnonUid = pulumi.IntPtr(int(*opt.AnonUid))
		}
		if opt.AnonGid != nil {
			exportOpt.AnonGid = pulumi.IntPtr(int(*opt.AnonGid))
		}
		// Source VPC network (name) for ip_ranges — GCP requires it on
		// PSC instances where client IPs aren't otherwise attributable.
		if opt.Network.GetValue() != "" {
			exportOpt.Network = pulumi.StringPtr(opt.Network.GetValue())
		}

		nfsExportOptions = append(nfsExportOptions, exportOpt)
	}

	// Build the file share configuration (singular — one per instance).
	fileShareArgs := &filestore.InstanceFileSharesArgs{
		Name:       pulumi.String(spec.FileShare.Name),
		CapacityGb: pulumi.Int(int(spec.FileShare.CapacityGb)),
	}
	if len(nfsExportOptions) > 0 {
		fileShareArgs.NfsExportOptions = nfsExportOptions
	}
	// Restore is create-time only and single-source (Filestore backup OR
	// Backup and DR backup — CEL-enforced); the share's capacity must
	// cover the backup's source capacity.
	if spec.FileShare.SourceBackup != "" {
		fileShareArgs.SourceBackup = pulumi.StringPtr(spec.FileShare.SourceBackup)
	}
	if spec.FileShare.SourceBackupdrBackup != "" {
		fileShareArgs.SourceBackupdrBackup = pulumi.StringPtr(spec.FileShare.SourceBackupdrBackup)
	}

	// Build the network configuration (singular — one per instance).
	// Empty modes follow the spec's documented default: IPv4 service.
	modes := pulumi.StringArray{}
	if len(spec.NetworkConfig.Modes) > 0 {
		modes = pulumi.ToStringArray(spec.NetworkConfig.Modes)
	} else {
		modes = pulumi.StringArray{pulumi.String("MODE_IPV4")}
	}
	networkArgs := &filestore.InstanceNetworkArgs{
		Network: pulumi.String(spec.NetworkConfig.Network.GetValue()),
		Modes:   modes,
	}
	if spec.NetworkConfig.ConnectMode != "" {
		networkArgs.ConnectMode = pulumi.StringPtr(spec.NetworkConfig.ConnectMode)
	}
	if spec.NetworkConfig.ReservedIpRange != "" {
		networkArgs.ReservedIpRange = pulumi.StringPtr(spec.NetworkConfig.ReservedIpRange)
	}
	// Consumer project hosting the PSC endpoint; PSC connect mode only
	// (CEL-enforced). Omitted means the instance's own project.
	if spec.NetworkConfig.PscEndpointProject.GetValue() != "" {
		networkArgs.PscConfig = &filestore.InstanceNetworkPscConfigArgs{
			EndpointProject: pulumi.StringPtr(spec.NetworkConfig.PscEndpointProject.GetValue()),
		}
	}

	args := &filestore.InstanceArgs{
		Name:       pulumi.String(locals.InstanceName),
		Location:   pulumi.StringPtr(spec.Location),
		Tier:       pulumi.String(spec.Tier),
		FileShares: fileShareArgs,
		Networks:   filestore.InstanceNetworkArray{networkArgs},
		Labels:     pulumi.ToStringMap(locals.GcpLabels),
	}

	// Client-side destroy behavior (DELETE deletes the share's data;
	// PREVENT refuses; ABANDON drops from state but keeps the instance
	// running), evaluated only after deletion_protection_enabled allows
	// the destroy. Empty follows the provider default (DELETE) —
	// mirrored zero-vs-omit with Terraform.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.Protocol != "" {
		args.Protocol = pulumi.StringPtr(spec.Protocol)
	}
	if spec.KmsKeyName.GetValue() != "" {
		args.KmsKeyName = pulumi.StringPtr(spec.KmsKeyName.GetValue())
	}
	args.DeletionProtectionEnabled = pulumi.BoolPtr(spec.DeletionProtectionEnabled)
	if spec.DeletionProtectionReason != "" {
		args.DeletionProtectionReason = pulumi.StringPtr(spec.DeletionProtectionReason)
	}
	if len(spec.Tags) > 0 {
		args.Tags = pulumi.ToStringMap(spec.Tags)
	}

	// IOPS tuning (ZONAL/REGIONAL/ENTERPRISE tiers).
	if spec.PerformanceConfig != nil {
		perfArgs := &filestore.InstancePerformanceConfigArgs{}
		if spec.PerformanceConfig.FixedIops != nil {
			perfArgs.FixedIops = &filestore.InstancePerformanceConfigFixedIopsArgs{
				MaxIops: pulumi.IntPtr(int(spec.PerformanceConfig.FixedIops.MaxIops)),
			}
		}
		if spec.PerformanceConfig.IopsPerTb != nil {
			perfArgs.IopsPerTb = &filestore.InstancePerformanceConfigIopsPerTbArgs{
				MaxIopsPerTb: pulumi.IntPtr(int(spec.PerformanceConfig.IopsPerTb.MaxIopsPerTb)),
			}
		}
		args.PerformanceConfig = perfArgs
	}

	// LDAP directory integration for NFSv4.1 identity mapping (protocol
	// NFS_V4_1 required — CEL-enforced pre-deploy).
	if spec.Ldap != nil {
		ldapArgs := &filestore.InstanceDirectoryServicesLdapArgs{
			Domain:  pulumi.String(spec.Ldap.Domain),
			Servers: pulumi.ToStringArray(spec.Ldap.Servers),
		}
		if spec.Ldap.GroupsOu != "" {
			ldapArgs.GroupsOu = pulumi.StringPtr(spec.Ldap.GroupsOu)
		}
		if spec.Ldap.UsersOu != "" {
			ldapArgs.UsersOu = pulumi.StringPtr(spec.Ldap.UsersOu)
		}
		args.DirectoryServices = &filestore.InstanceDirectoryServicesArgs{
			Ldap: ldapArgs,
		}
	}

	// Replica-relationship state: the provider pauses/resumes replication
	// to match. Sent explicitly from the spec (default READY) so both
	// engines and the provider agree on who chose the value — the Redis
	// deletion_protection posture.
	desiredReplicaState := "READY"
	if spec.DesiredReplicaState != nil && spec.GetDesiredReplicaState() != "" {
		desiredReplicaState = spec.GetDesiredReplicaState()
	}
	args.DesiredReplicaState = pulumi.StringPtr(desiredReplicaState)

	// Create-time replication: this instance joins as ACTIVE source or
	// STANDBY replica of the referenced peers. Backups cannot be taken
	// from a STANDBY replica.
	if spec.InitialReplication != nil {
		replicas := filestore.InstanceInitialReplicationReplicaArray{}
		for _, peer := range spec.InitialReplication.PeerInstances {
			replicas = append(replicas, &filestore.InstanceInitialReplicationReplicaArgs{
				PeerInstance: pulumi.String(peer.GetValue()),
			})
		}
		replicationArgs := &filestore.InstanceInitialReplicationArgs{
			Replicas: replicas,
		}
		if spec.InitialReplication.Role != "" {
			replicationArgs.Role = pulumi.StringPtr(spec.InitialReplication.Role)
		}
		args.InitialReplication = replicationArgs
	}

	createdInstance, err := filestore.NewInstance(ctx,
		"filestore-instance",
		args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdFileApi}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create filestore instance")
	}

	// Semantic outputs — names and shapes byte-identical to the Terraform
	// module's outputs.
	ctx.Export(OpInstanceId, createdInstance.ID())
	ctx.Export(OpInstanceName, createdInstance.Name)
	ctx.Export(OpFileShareName, pulumi.String(spec.FileShare.Name))
	ctx.Export(OpCreateTime, createdInstance.CreateTime)
	ctx.Export(OpEtag, createdInstance.Etag)

	// Extract addresses and the GCP-resolved reserved range from the
	// first (only) network.
	ipAddresses := createdInstance.Networks.ApplyT(func(networks []filestore.InstanceNetwork) []string {
		if len(networks) > 0 && len(networks[0].IpAddresses) > 0 {
			return networks[0].IpAddresses
		}
		return []string{}
	}).(pulumi.StringArrayOutput)
	ctx.Export(OpIpAddresses, ipAddresses)

	reservedIpRange := createdInstance.Networks.ApplyT(func(networks []filestore.InstanceNetwork) string {
		if len(networks) > 0 && networks[0].ReservedIpRange != nil {
			return *networks[0].ReservedIpRange
		}
		return ""
	}).(pulumi.StringOutput)
	ctx.Export(OpReservedIpRange, reservedIpRange)

	return nil
}
