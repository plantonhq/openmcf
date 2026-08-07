package module

import (
	"github.com/pkg/errors"
	gcpgkenodepoolv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpgkenodepool/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// boolOr reads an optional proto bool with the field's documented default,
// so a manifest that omits the field and GKE's own behavior agree — the
// same contract the Terraform module expresses with optional(bool, X).
func boolOr(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

// nodePool provisions a GKE node pool — a group of identically configured
// Compute Engine VMs attached to a GcpGkeCluster. The cluster and location
// arrive as plain strings resolved from the cluster's outputs, so the pool
// addresses its parent exactly the way GKE named it — no lookup.
//
// Lifecycle notes the API enforces:
//   - name, location, initial_node_count, max_pods_per_node,
//     placement_policy, queued_provisioning, and nearly all of node_config
//     (machine type, disks, image, identity, accelerators, shielded/
//     confidential settings, local SSDs) are immutable — changing them
//     replaces the pool (GKE drains and recreates the nodes).
//   - node_count, autoscaling, management, upgrade_settings,
//     node_locations, labels, taints, tags, and resource_labels update in
//     place.
//   - For autoscaled pools the autoscaler owns node_count at runtime;
//     IgnoreChanges keeps the preview from fighting it.
func nodePool(ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider) error {

	spec := locals.GcpGkeNodePool.Spec

	// Enable the Kubernetes Engine API first so a fresh project works on
	// the first deploy. disable_on_destroy stays false: tearing down one
	// node pool must never disable the API for everything else in the
	// project (including its own cluster).
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("container.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"gkenp-container.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable container.googleapis.com api")
	}

	args := &container.NodePoolArgs{
		Name:     pulumi.String(locals.NodePoolName),
		Cluster:  pulumi.String(spec.ClusterName.GetValue()),
		Location: pulumi.StringPtr(spec.Location.GetValue()),
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	// Nodes may span fewer zones than the cluster; empty inherits the
	// cluster's node_locations.
	if len(spec.NodeLocations) > 0 {
		args.NodeLocations = pulumi.ToStringArray(spec.NodeLocations)
	}

	// Explicit version and auto_upgrade fight each other (the API
	// re-upgrades what the preview pins); the spec documents the trade and
	// the field passes through untouched.
	if spec.Version != "" {
		args.Version = pulumi.StringPtr(spec.Version)
	}

	if spec.MaxPodsPerNode != nil {
		args.MaxPodsPerNode = pulumi.IntPtr(int(*spec.MaxPodsPerNode))
	}
	if spec.InitialNodeCount != nil {
		args.InitialNodeCount = pulumi.IntPtr(int(*spec.InitialNodeCount))
	}

	// Fixed size XOR autoscaling — the proto oneof guarantees at most one
	// arrives. Neither means GKE's default size (3 nodes), unmanaged. The
	// effective min/max exported below mirror whichever mode is active so
	// downstream consumers read one honest pair.
	effectiveMin := 0
	effectiveMax := 0
	switch size := spec.NodePoolSize.(type) {
	case *gcpgkenodepoolv1alpha1.GcpGkeNodePoolSpec_NodeCount:
		args.NodeCount = pulumi.IntPtr(int(size.NodeCount))
		effectiveMin = int(size.NodeCount)
		effectiveMax = int(size.NodeCount)
	case *gcpgkenodepoolv1alpha1.GcpGkeNodePoolSpec_Autoscaling:
		autoscaling := size.Autoscaling
		autoscalingArgs := &container.NodePoolAutoscalingArgs{}
		// Per-zone bounds XOR total bounds (spec-level CEL enforces the
		// exclusivity); the unused arm stays out of the API payload.
		if autoscaling.MinNodes != nil {
			autoscalingArgs.MinNodeCount = pulumi.IntPtr(int(*autoscaling.MinNodes))
			effectiveMin = int(*autoscaling.MinNodes)
		}
		if autoscaling.MaxNodes != nil {
			autoscalingArgs.MaxNodeCount = pulumi.IntPtr(int(*autoscaling.MaxNodes))
			effectiveMax = int(*autoscaling.MaxNodes)
		}
		if autoscaling.TotalMinNodes != nil {
			autoscalingArgs.TotalMinNodeCount = pulumi.IntPtr(int(*autoscaling.TotalMinNodes))
			effectiveMin = int(*autoscaling.TotalMinNodes)
		}
		if autoscaling.TotalMaxNodes != nil {
			autoscalingArgs.TotalMaxNodeCount = pulumi.IntPtr(int(*autoscaling.TotalMaxNodes))
			effectiveMax = int(*autoscaling.TotalMaxNodes)
		}
		if autoscaling.LocationPolicy != "" {
			autoscalingArgs.LocationPolicy = pulumi.StringPtr(autoscaling.LocationPolicy)
		}
		args.Autoscaling = autoscalingArgs
	}

	// Auto-repair/auto-upgrade both default true — GKE's own defaults — so
	// an omitted management block and GKE's behavior agree.
	autoRepair, autoUpgrade := true, true
	if spec.Management != nil {
		autoRepair = boolOr(spec.Management.AutoRepair, true)
		autoUpgrade = boolOr(spec.Management.AutoUpgrade, true)
	}
	args.Management = &container.NodePoolManagementArgs{
		AutoRepair:  pulumi.Bool(autoRepair),
		AutoUpgrade: pulumi.Bool(autoUpgrade),
	}

	// Emitted only when the spec configures it: GKE's default surge
	// settings (max_surge=1, max_unavailable=0) apply otherwise.
	if upgrade := spec.UpgradeSettings; upgrade != nil {
		upgradeArgs := &container.NodePoolUpgradeSettingsArgs{}
		if upgrade.MaxSurge != nil {
			upgradeArgs.MaxSurge = pulumi.IntPtr(int(*upgrade.MaxSurge))
		}
		if upgrade.MaxUnavailable != nil {
			upgradeArgs.MaxUnavailable = pulumi.IntPtr(int(*upgrade.MaxUnavailable))
		}
		if upgrade.Strategy != "" {
			upgradeArgs.Strategy = pulumi.StringPtr(upgrade.Strategy)
		}
		if blueGreen := upgrade.BlueGreenSettings; blueGreen != nil {
			blueGreenArgs := &container.NodePoolUpgradeSettingsBlueGreenSettingsArgs{}
			if blueGreen.NodePoolSoakDuration != "" {
				blueGreenArgs.NodePoolSoakDuration = pulumi.StringPtr(blueGreen.NodePoolSoakDuration)
			}
			rollout := blueGreen.StandardRolloutPolicy
			rolloutArgs := &container.NodePoolUpgradeSettingsBlueGreenSettingsStandardRolloutPolicyArgs{}
			if rollout.BatchPercentage != nil {
				rolloutArgs.BatchPercentage = pulumi.Float64Ptr(float64(*rollout.BatchPercentage))
			}
			if rollout.BatchNodeCount != nil {
				rolloutArgs.BatchNodeCount = pulumi.IntPtr(int(*rollout.BatchNodeCount))
			}
			if rollout.BatchSoakDuration != "" {
				rolloutArgs.BatchSoakDuration = pulumi.StringPtr(rollout.BatchSoakDuration)
			}
			blueGreenArgs.StandardRolloutPolicy = rolloutArgs
			upgradeArgs.BlueGreenSettings = blueGreenArgs
		}
		args.UpgradeSettings = upgradeArgs
	}

	if placement := spec.PlacementPolicy; placement != nil {
		placementArgs := &container.NodePoolPlacementPolicyArgs{
			Type: pulumi.String(placement.Type),
		}
		if placement.PolicyName != "" {
			placementArgs.PolicyName = pulumi.StringPtr(placement.PolicyName)
		}
		if placement.TpuTopology != "" {
			placementArgs.TpuTopology = pulumi.StringPtr(placement.TpuTopology)
		}
		args.PlacementPolicy = placementArgs
	}

	if spec.QueuedProvisioningEnabled {
		args.QueuedProvisioning = &container.NodePoolQueuedProvisioningArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	if network := spec.NetworkConfig; network != nil {
		networkArgs := &container.NodePoolNetworkConfigArgs{}
		if network.CreatePodRange {
			networkArgs.CreatePodRange = pulumi.BoolPtr(true)
		}
		if network.PodRange != "" {
			networkArgs.PodRange = pulumi.StringPtr(network.PodRange)
		}
		if network.PodIpv4CidrBlock != "" {
			networkArgs.PodIpv4CidrBlock = pulumi.StringPtr(network.PodIpv4CidrBlock)
		}
		if network.EnablePrivateNodes != nil {
			networkArgs.EnablePrivateNodes = pulumi.BoolPtr(*network.EnablePrivateNodes)
		}
		if network.TotalEgressBandwidthTier != "" {
			networkArgs.NetworkPerformanceConfig = &container.NodePoolNetworkConfigNetworkPerformanceConfigArgs{
				TotalEgressBandwidthTier: pulumi.String(network.TotalEgressBandwidthTier),
			}
		}
		if network.PodCidrOverprovisionDisabled {
			networkArgs.PodCidrOverprovisionConfig = &container.NodePoolNetworkConfigPodCidrOverprovisionConfigArgs{
				Disabled: pulumi.Bool(true),
			}
		}
		args.NetworkConfig = networkArgs
	}

	args.NodeConfig = buildNodeConfig(spec.NodeConfig, locals)

	createdNodePool, err := container.NewNodePool(ctx,
		locals.NodePoolName,
		args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdProjectService}),
		// For autoscaled pools the cluster autoscaler owns the live node
		// count; without this, every preview after a scale event would try
		// to reset it.
		pulumi.IgnoreChanges([]string{"nodeCount"}))
	if err != nil {
		return errors.Wrap(err, "failed to create node pool")
	}

	ctx.Export(OpNodePoolName, createdNodePool.Name)
	ctx.Export(OpInstanceGroupUrls, createdNodePool.InstanceGroupUrls)
	ctx.Export(OpMinNodes, pulumi.Int(effectiveMin))
	ctx.Export(OpMaxNodes, pulumi.Int(effectiveMax))
	ctx.Export(OpCurrentNodeCount, createdNodePool.NodeCount)
	ctx.Export(OpNodePoolId, createdNodePool.ID())
	// The plain spec location (not a provider attribute) so the output is
	// the same human-readable region/zone name on both engines.
	ctx.Export(OpLocation, pulumi.String(spec.Location.GetValue()))
	ctx.Export(OpVersion, createdNodePool.Version)

	return nil
}

// buildNodeConfig translates the spec's node_config to provider args. The
// block is always emitted: the platform attribution resource_labels and the
// disable-legacy-endpoints metadata guard apply to every pool, spec block
// or not. Everything else honors "unset means GKE default" by omission.
func buildNodeConfig(nodeConfig *gcpgkenodepoolv1alpha1.GcpGkeNodePoolNodeConfig, locals *Locals) *container.NodePoolNodeConfigArgs {
	// GKE requires disable-legacy-endpoints=true on every node; it is
	// enforced beneath any user metadata so a spec entry can never weaken
	// it.
	metadata := map[string]string{}
	if nodeConfig != nil {
		for key, value := range nodeConfig.Metadata {
			metadata[key] = value
		}
	}
	metadata["disable-legacy-endpoints"] = "true"

	nodeConfigArgs := &container.NodePoolNodeConfigArgs{
		ResourceLabels: pulumi.ToStringMap(locals.GcpLabels),
		Metadata:       pulumi.ToStringMap(metadata),
	}

	// Defaults mirror the spec's field options (and the Terraform module's
	// optional() defaults) so both engines send identical payloads for an
	// omitted node_config.
	machineType := "e2-medium"
	imageType := "COS_CONTAINERD"

	if nodeConfig == nil {
		nodeConfigArgs.MachineType = pulumi.StringPtr(machineType)
		nodeConfigArgs.ImageType = pulumi.StringPtr(imageType)
		return nodeConfigArgs
	}

	if nodeConfig.GetMachineType() != "" {
		machineType = nodeConfig.GetMachineType()
	}
	if nodeConfig.GetImageType() != "" {
		imageType = nodeConfig.GetImageType()
	}
	nodeConfigArgs.MachineType = pulumi.StringPtr(machineType)
	nodeConfigArgs.ImageType = pulumi.StringPtr(imageType)

	if nodeConfig.DiskSizeGb != nil {
		nodeConfigArgs.DiskSizeGb = pulumi.IntPtr(int(*nodeConfig.DiskSizeGb))
	}
	if nodeConfig.DiskType != "" {
		nodeConfigArgs.DiskType = pulumi.StringPtr(nodeConfig.DiskType)
	}
	if nodeConfig.ServiceAccount.GetValue() != "" {
		nodeConfigArgs.ServiceAccount = pulumi.StringPtr(nodeConfig.ServiceAccount.GetValue())
	}

	// Empty applies GKE's default scopes; with Workload Identity, workload
	// permissions come from IAM on Kubernetes service accounts, not node
	// scopes, so the defaults are normally right.
	if len(nodeConfig.OauthScopes) > 0 {
		nodeConfigArgs.OauthScopes = pulumi.ToStringArray(nodeConfig.OauthScopes)
	}

	// Kubernetes node labels (nodeSelector targets) — distinct from the
	// GCE resource labels above, which carry the platform attribution.
	if len(nodeConfig.Labels) > 0 {
		nodeConfigArgs.Labels = pulumi.ToStringMap(nodeConfig.Labels)
	}
	if len(nodeConfig.Tags) > 0 {
		nodeConfigArgs.Tags = pulumi.ToStringArray(nodeConfig.Tags)
	}

	if len(nodeConfig.Taints) > 0 {
		taints := container.NodePoolNodeConfigTaintArray{}
		for _, taint := range nodeConfig.Taints {
			taints = append(taints, &container.NodePoolNodeConfigTaintArgs{
				Key:    pulumi.String(taint.Key),
				Value:  pulumi.String(taint.Value),
				Effect: pulumi.String(taint.Effect),
			})
		}
		nodeConfigArgs.Taints = taints
	}

	// spot supersedes preemptible (no 24h lifetime); spec-level CEL rejects
	// both together, so each passes through independently.
	if nodeConfig.Spot {
		nodeConfigArgs.Spot = pulumi.BoolPtr(true)
	}
	if nodeConfig.Preemptible {
		nodeConfigArgs.Preemptible = pulumi.BoolPtr(true)
	}

	if len(nodeConfig.GuestAccelerators) > 0 {
		accelerators := container.NodePoolNodeConfigGuestAcceleratorArray{}
		for _, accelerator := range nodeConfig.GuestAccelerators {
			acceleratorArgs := &container.NodePoolNodeConfigGuestAcceleratorArgs{
				Type:  pulumi.String(accelerator.Type),
				Count: pulumi.Int(int(accelerator.Count)),
			}
			if accelerator.GpuPartitionSize != "" {
				acceleratorArgs.GpuPartitionSize = pulumi.StringPtr(accelerator.GpuPartitionSize)
			}
			if accelerator.GpuDriverVersion != "" {
				acceleratorArgs.GpuDriverInstallationConfig = &container.NodePoolNodeConfigGuestAcceleratorGpuDriverInstallationConfigArgs{
					GpuDriverVersion: pulumi.String(accelerator.GpuDriverVersion),
				}
			}
			if sharing := accelerator.GpuSharingConfig; sharing != nil {
				acceleratorArgs.GpuSharingConfig = &container.NodePoolNodeConfigGuestAcceleratorGpuSharingConfigArgs{
					GpuSharingStrategy:     pulumi.String(sharing.GpuSharingStrategy),
					MaxSharedClientsPerGpu: pulumi.Int(int(sharing.MaxSharedClientsPerGpu)),
				}
			}
			accelerators = append(accelerators, acceleratorArgs)
		}
		nodeConfigArgs.GuestAccelerators = accelerators
	}

	if shielded := nodeConfig.ShieldedInstanceConfig; shielded != nil {
		nodeConfigArgs.ShieldedInstanceConfig = &container.NodePoolNodeConfigShieldedInstanceConfigArgs{
			EnableSecureBoot:          pulumi.BoolPtr(shielded.EnableSecureBoot),
			EnableIntegrityMonitoring: pulumi.BoolPtr(boolOr(shielded.EnableIntegrityMonitoring, true)),
		}
	}

	if confidential := nodeConfig.ConfidentialNodes; confidential != nil {
		confidentialArgs := &container.NodePoolNodeConfigConfidentialNodesArgs{
			Enabled: pulumi.Bool(confidential.Enabled),
		}
		if confidential.ConfidentialInstanceType != "" {
			confidentialArgs.ConfidentialInstanceType = pulumi.StringPtr(confidential.ConfidentialInstanceType)
		}
		nodeConfigArgs.ConfidentialNodes = confidentialArgs
	}

	if nodeConfig.MinCpuPlatform != "" {
		nodeConfigArgs.MinCpuPlatform = pulumi.StringPtr(nodeConfig.MinCpuPlatform)
	}
	if nodeConfig.LocalSsdCount != nil {
		nodeConfigArgs.LocalSsdCount = pulumi.IntPtr(int(*nodeConfig.LocalSsdCount))
	}

	if ephemeral := nodeConfig.EphemeralStorageLocalSsd; ephemeral != nil {
		ephemeralArgs := &container.NodePoolNodeConfigEphemeralStorageLocalSsdConfigArgs{
			LocalSsdCount: pulumi.Int(int(ephemeral.LocalSsdCount)),
		}
		if ephemeral.DataCacheCount != nil {
			ephemeralArgs.DataCacheCount = pulumi.IntPtr(int(*ephemeral.DataCacheCount))
		}
		nodeConfigArgs.EphemeralStorageLocalSsdConfig = ephemeralArgs
	}

	if nvmeBlock := nodeConfig.LocalNvmeSsdBlock; nvmeBlock != nil {
		nodeConfigArgs.LocalNvmeSsdBlockConfig = &container.NodePoolNodeConfigLocalNvmeSsdBlockConfigArgs{
			LocalSsdCount: pulumi.Int(int(nvmeBlock.LocalSsdCount)),
		}
	}

	if nodeConfig.GcfsEnabled {
		nodeConfigArgs.GcfsConfig = &container.NodePoolNodeConfigGcfsConfigArgs{
			Enabled: pulumi.Bool(true),
		}
	}
	if nodeConfig.GvnicEnabled {
		nodeConfigArgs.Gvnic = &container.NodePoolNodeConfigGvnicArgs{
			Enabled: pulumi.Bool(true),
		}
	}
	if nodeConfig.FastSocketEnabled {
		nodeConfigArgs.FastSocket = &container.NodePoolNodeConfigFastSocketArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	if nodeConfig.BootDiskKmsKey.GetValue() != "" {
		nodeConfigArgs.BootDiskKmsKey = pulumi.StringPtr(nodeConfig.BootDiskKmsKey.GetValue())
	}

	if nodeConfig.WorkloadMetadataMode != "" {
		nodeConfigArgs.WorkloadMetadataConfig = &container.NodePoolNodeConfigWorkloadMetadataConfigArgs{
			Mode: pulumi.String(nodeConfig.WorkloadMetadataMode),
		}
	}

	if reservation := nodeConfig.ReservationAffinity; reservation != nil {
		reservationArgs := &container.NodePoolNodeConfigReservationAffinityArgs{
			ConsumeReservationType: pulumi.String(reservation.ConsumeReservationType),
		}
		if reservation.Key != "" {
			reservationArgs.Key = pulumi.StringPtr(reservation.Key)
		}
		if len(reservation.Values) > 0 {
			reservationArgs.Values = pulumi.ToStringArray(reservation.Values)
		}
		nodeConfigArgs.ReservationAffinity = reservationArgs
	}

	if len(nodeConfig.SecondaryBootDisks) > 0 {
		disks := container.NodePoolNodeConfigSecondaryBootDiskArray{}
		for _, disk := range nodeConfig.SecondaryBootDisks {
			diskArgs := &container.NodePoolNodeConfigSecondaryBootDiskArgs{
				DiskImage: pulumi.String(disk.DiskImage),
			}
			if disk.Mode != "" {
				diskArgs.Mode = pulumi.StringPtr(disk.Mode)
			}
			disks = append(disks, diskArgs)
		}
		nodeConfigArgs.SecondaryBootDisks = disks
	}

	if kubelet := nodeConfig.KubeletConfig; kubelet != nil {
		kubeletArgs := &container.NodePoolNodeConfigKubeletConfigArgs{}
		if kubelet.CpuManagerPolicy != "" {
			kubeletArgs.CpuManagerPolicy = pulumi.StringPtr(kubelet.CpuManagerPolicy)
		}
		if kubelet.CpuCfsQuota != nil {
			kubeletArgs.CpuCfsQuota = pulumi.BoolPtr(*kubelet.CpuCfsQuota)
		}
		if kubelet.CpuCfsQuotaPeriod != "" {
			kubeletArgs.CpuCfsQuotaPeriod = pulumi.StringPtr(kubelet.CpuCfsQuotaPeriod)
		}
		if kubelet.PodPidsLimit != nil {
			kubeletArgs.PodPidsLimit = pulumi.IntPtr(int(*kubelet.PodPidsLimit))
		}
		if kubelet.InsecureKubeletReadonlyPortEnabled != "" {
			kubeletArgs.InsecureKubeletReadonlyPortEnabled = pulumi.StringPtr(kubelet.InsecureKubeletReadonlyPortEnabled)
		}
		if kubelet.MaxParallelImagePulls != nil {
			kubeletArgs.MaxParallelImagePulls = pulumi.IntPtr(int(*kubelet.MaxParallelImagePulls))
		}
		if kubelet.ContainerLogMaxSize != "" {
			kubeletArgs.ContainerLogMaxSize = pulumi.StringPtr(kubelet.ContainerLogMaxSize)
		}
		if kubelet.ContainerLogMaxFiles != nil {
			kubeletArgs.ContainerLogMaxFiles = pulumi.IntPtr(int(*kubelet.ContainerLogMaxFiles))
		}
		if kubelet.ImageGcLowThresholdPercent != nil {
			kubeletArgs.ImageGcLowThresholdPercent = pulumi.IntPtr(int(*kubelet.ImageGcLowThresholdPercent))
		}
		if kubelet.ImageGcHighThresholdPercent != nil {
			kubeletArgs.ImageGcHighThresholdPercent = pulumi.IntPtr(int(*kubelet.ImageGcHighThresholdPercent))
		}
		nodeConfigArgs.KubeletConfig = kubeletArgs
	}

	if linux := nodeConfig.LinuxNodeConfig; linux != nil {
		linuxArgs := &container.NodePoolNodeConfigLinuxNodeConfigArgs{}
		if len(linux.Sysctls) > 0 {
			linuxArgs.Sysctls = pulumi.ToStringMap(linux.Sysctls)
		}
		if linux.CgroupMode != "" {
			linuxArgs.CgroupMode = pulumi.StringPtr(linux.CgroupMode)
		}
		if hugepages := linux.HugepagesConfig; hugepages != nil {
			hugepagesArgs := &container.NodePoolNodeConfigLinuxNodeConfigHugepagesConfigArgs{}
			if hugepages.HugepageSize_2M != nil {
				hugepagesArgs.HugepageSize2m = pulumi.IntPtr(int(*hugepages.HugepageSize_2M))
			}
			if hugepages.HugepageSize_1G != nil {
				hugepagesArgs.HugepageSize1g = pulumi.IntPtr(int(*hugepages.HugepageSize_1G))
			}
			linuxArgs.HugepagesConfig = hugepagesArgs
		}
		nodeConfigArgs.LinuxNodeConfig = linuxArgs
	}

	if nodeConfig.LoggingVariant != "" {
		nodeConfigArgs.LoggingVariant = pulumi.StringPtr(nodeConfig.LoggingVariant)
	}
	if nodeConfig.FlexStart {
		nodeConfigArgs.FlexStart = pulumi.BoolPtr(true)
	}
	if nodeConfig.MaxRunDuration != "" {
		nodeConfigArgs.MaxRunDuration = pulumi.StringPtr(nodeConfig.MaxRunDuration)
	}

	return nodeConfigArgs
}
