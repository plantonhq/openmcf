package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// cluster provisions a GKE cluster — the managed Kubernetes control plane
// plus cluster-wide configuration. Node pools are separate GcpGkeNodePool
// resources: the default node pool is always removed at create time on
// Standard clusters, so every pool is an explicitly managed, first-class
// node. Autopilot clusters manage nodes themselves and take no node pools.
//
// Lifecycle notes the API enforces:
//   - name, location, description, network, subnetwork, the whole
//     ip_allocation_policy, datapath_provider, default_max_pods_per_node,
//     confidential_nodes, enable_autopilot, and the private control-plane
//     placement fields are immutable — changing any of them replaces the
//     cluster (and everything running on it).
//   - enable_l4_ilb_subsetting is one-way: it can be turned on in place but
//     never off.
//   - deletion_protection is an engine-side guard: while true, a destroy
//     preview fails before touching the cluster.
func cluster(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpGkeCluster.Spec

	// Enable the Kubernetes Engine API first so a fresh project works on the
	// first deploy. disable_on_destroy stays false: tearing down one cluster
	// must never disable the API for everything else in the project.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("container.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"gke-container.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable container.googleapis.com api")
	}

	// The Workload Identity pool name is fixed by the API to
	// PROJECT_ID.svc.id.goog; when the spec omits project_id the ambient
	// provider project fills it in (the ambient-project contract).
	workloadPoolProject := spec.ProjectId.GetValue()
	if workloadPoolProject == "" {
		clientConfig, err := organizations.GetClientConfig(ctx, pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to resolve ambient project from provider client config")
		}
		workloadPoolProject = clientConfig.Project
	}
	workloadPool := fmt.Sprintf("%s.svc.id.goog", workloadPoolProject)

	args := &container.ClusterArgs{
		Name:     pulumi.String(locals.ClusterName),
		Location: pulumi.StringPtr(spec.Location),

		Network:    pulumi.StringPtr(spec.Network.GetValue()),
		Subnetwork: pulumi.StringPtr(spec.Subnetwork.GetValue()),

		DeletionProtection: pulumi.BoolPtr(spec.GetDeletionProtection()),

		ResourceLabels: pulumi.ToStringMap(locals.GcpLabels),

		ReleaseChannel: &container.ClusterReleaseChannelArgs{
			Channel: pulumi.String(locals.ReleaseChannel),
		},
	}

	// An empty project falls back to the provider's default project — the
	// ambient-project contract every GCP kind honors.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// Nodes may span fewer (regional) or more (zonal) zones than the control
	// plane; empty defers to GKE's location defaults.
	if len(spec.NodeLocations) > 0 {
		args.NodeLocations = pulumi.ToStringArray(spec.NodeLocations)
	}

	// Clusters are always VPC-native (alias IP). The ip_allocation_policy
	// block is emitted even when the spec omits ip_allocation: an empty
	// block tells GKE to create and manage the pod/service secondary ranges
	// itself, while named ranges pin allocation to planned subnetwork
	// ranges.
	args.NetworkingMode = pulumi.StringPtr("VPC_NATIVE")
	ipAllocationArgs := &container.ClusterIpAllocationPolicyArgs{}
	if spec.IpAllocation != nil {
		if spec.IpAllocation.ClusterSecondaryRangeName.GetValue() != "" {
			ipAllocationArgs.ClusterSecondaryRangeName = pulumi.StringPtr(spec.IpAllocation.ClusterSecondaryRangeName.GetValue())
		}
		if spec.IpAllocation.ServicesSecondaryRangeName.GetValue() != "" {
			ipAllocationArgs.ServicesSecondaryRangeName = pulumi.StringPtr(spec.IpAllocation.ServicesSecondaryRangeName.GetValue())
		}
		if spec.IpAllocation.ClusterIpv4CidrBlock != "" {
			ipAllocationArgs.ClusterIpv4CidrBlock = pulumi.StringPtr(spec.IpAllocation.ClusterIpv4CidrBlock)
		}
		if spec.IpAllocation.ServicesIpv4CidrBlock != "" {
			ipAllocationArgs.ServicesIpv4CidrBlock = pulumi.StringPtr(spec.IpAllocation.ServicesIpv4CidrBlock)
		}
		ipAllocationArgs.StackType = pulumi.StringPtr(spec.IpAllocation.GetStackType())
		if len(spec.IpAllocation.AdditionalPodRangeNames) > 0 {
			ipAllocationArgs.AdditionalPodRangesConfig = &container.ClusterIpAllocationPolicyAdditionalPodRangesConfigArgs{
				PodRangeNames: pulumi.ToStringArray(spec.IpAllocation.AdditionalPodRangeNames),
			}
		}
		if spec.IpAllocation.PodCidrOverprovisionDisabled {
			ipAllocationArgs.PodCidrOverprovisionConfig = &container.ClusterIpAllocationPolicyPodCidrOverprovisionConfigArgs{
				Disabled: pulumi.Bool(true),
			}
		}
	}
	args.IpAllocationPolicy = ipAllocationArgs

	// Standard clusters: drop the API-mandated default node pool
	// immediately — node pools are composed as GcpGkeNodePool resources.
	// Autopilot rejects both fields (GKE owns node management there).
	if spec.EnableAutopilot {
		args.EnableAutopilot = pulumi.BoolPtr(true)
		if spec.AllowNetAdmin {
			args.AllowNetAdmin = pulumi.BoolPtr(true)
		}
	} else {
		args.RemoveDefaultNodePool = pulumi.BoolPtr(true)
		args.InitialNodeCount = pulumi.IntPtr(1)
	}

	if spec.DatapathProvider != "" {
		args.DatapathProvider = pulumi.StringPtr(spec.DatapathProvider)
	}
	if spec.DefaultMaxPodsPerNode != nil {
		args.DefaultMaxPodsPerNode = pulumi.IntPtr(int(spec.GetDefaultMaxPodsPerNode()))
	}
	if spec.EnableIntranodeVisibility {
		args.EnableIntranodeVisibility = pulumi.BoolPtr(true)
	}
	if spec.EnableL4IlbSubsetting {
		args.EnableL4IlbSubsetting = pulumi.BoolPtr(true)
	}
	if spec.EnableFqdnNetworkPolicy {
		args.EnableFqdnNetworkPolicy = pulumi.BoolPtr(true)
	}
	if spec.EnableCiliumClusterwideNetworkPolicy {
		args.EnableCiliumClusterwideNetworkPolicy = pulumi.BoolPtr(true)
	}
	if spec.EnableMultiNetworking {
		args.EnableMultiNetworking = pulumi.BoolPtr(true)
	}
	if spec.PrivateIpv6GoogleAccess != "" {
		args.PrivateIpv6GoogleAccess = pulumi.StringPtr(spec.PrivateIpv6GoogleAccess)
	}
	if spec.InTransitEncryption != "" {
		args.InTransitEncryptionConfig = pulumi.StringPtr(spec.InTransitEncryption)
	}
	if spec.DisableL4LbFirewallReconciliation {
		args.DisableL4LbFirewallReconciliation = pulumi.BoolPtr(true)
	}

	// Calico NetworkPolicy enforcement is two coupled settings: the
	// enforcement block here and the addon below. Both follow the single
	// spec toggle so they can never drift apart. Omitted entirely on
	// Autopilot (which enforces NetworkPolicy natively via Dataplane V2).
	if spec.EnableNetworkPolicy {
		args.NetworkPolicy = &container.ClusterNetworkPolicyArgs{
			Enabled:  pulumi.Bool(true),
			Provider: pulumi.StringPtr("CALICO"),
		}
	}

	if spec.DisableDefaultSnat {
		args.DefaultSnatStatus = &container.ClusterDefaultSnatStatusArgs{
			Disabled: pulumi.Bool(true),
		}
	}

	if spec.TotalEgressBandwidthTier != "" {
		args.NetworkPerformanceConfig = &container.ClusterNetworkPerformanceConfigArgs{
			TotalEgressBandwidthTier: pulumi.String(spec.TotalEgressBandwidthTier),
		}
	}

	if spec.DnsConfig != nil {
		dnsArgs := &container.ClusterDnsConfigArgs{}
		if spec.DnsConfig.ClusterDns != "" {
			dnsArgs.ClusterDns = pulumi.StringPtr(spec.DnsConfig.ClusterDns)
		}
		if spec.DnsConfig.ClusterDnsScope != "" {
			dnsArgs.ClusterDnsScope = pulumi.StringPtr(spec.DnsConfig.ClusterDnsScope)
		}
		if spec.DnsConfig.ClusterDnsDomain != "" {
			dnsArgs.ClusterDnsDomain = pulumi.StringPtr(spec.DnsConfig.ClusterDnsDomain)
		}
		if spec.DnsConfig.AdditiveVpcScopeDnsDomain != "" {
			dnsArgs.AdditiveVpcScopeDnsDomain = pulumi.StringPtr(spec.DnsConfig.AdditiveVpcScopeDnsDomain)
		}
		args.DnsConfig = dnsArgs
	}

	if spec.GatewayApiChannel != "" {
		args.GatewayApiConfig = &container.ClusterGatewayApiConfigArgs{
			Channel: pulumi.String(spec.GatewayApiChannel),
		}
	}

	if spec.EnableServiceExternalIps {
		args.ServiceExternalIpsConfig = &container.ClusterServiceExternalIpsConfigArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	// Private topology: private nodes need Cloud NAT (a GcpRouterNat on the
	// network) for outbound internet; a private-only endpoint removes
	// public kubectl access entirely.
	if spec.PrivateCluster != nil {
		privateArgs := &container.ClusterPrivateClusterConfigArgs{
			EnablePrivateNodes:    pulumi.BoolPtr(spec.PrivateCluster.EnablePrivateNodes),
			EnablePrivateEndpoint: pulumi.BoolPtr(spec.PrivateCluster.EnablePrivateEndpoint),
		}
		if spec.PrivateCluster.MasterIpv4CidrBlock != "" {
			privateArgs.MasterIpv4CidrBlock = pulumi.StringPtr(spec.PrivateCluster.MasterIpv4CidrBlock)
		}
		if spec.PrivateCluster.PrivateEndpointSubnetwork.GetValue() != "" {
			privateArgs.PrivateEndpointSubnetwork = pulumi.StringPtr(spec.PrivateCluster.PrivateEndpointSubnetwork.GetValue())
		}
		if spec.PrivateCluster.EnableMasterGlobalAccess {
			privateArgs.MasterGlobalAccessConfig = &container.ClusterPrivateClusterConfigMasterGlobalAccessConfigArgs{
				Enabled: pulumi.Bool(true),
			}
		}
		args.PrivateClusterConfig = privateArgs
	}

	if spec.MasterAuthorizedNetworks != nil {
		manArgs := &container.ClusterMasterAuthorizedNetworksConfigArgs{}
		if spec.MasterAuthorizedNetworks.GcpPublicCidrsAccessEnabled != nil {
			manArgs.GcpPublicCidrsAccessEnabled = pulumi.BoolPtr(spec.MasterAuthorizedNetworks.GetGcpPublicCidrsAccessEnabled())
		}
		if spec.MasterAuthorizedNetworks.PrivateEndpointEnforcementEnabled != nil {
			manArgs.PrivateEndpointEnforcementEnabled = pulumi.BoolPtr(spec.MasterAuthorizedNetworks.GetPrivateEndpointEnforcementEnabled())
		}
		cidrBlocks := container.ClusterMasterAuthorizedNetworksConfigCidrBlockArray{}
		for _, cidr := range spec.MasterAuthorizedNetworks.CidrBlocks {
			cidrArgs := &container.ClusterMasterAuthorizedNetworksConfigCidrBlockArgs{
				CidrBlock: pulumi.String(cidr.CidrBlock),
			}
			if cidr.DisplayName != "" {
				cidrArgs.DisplayName = pulumi.StringPtr(cidr.DisplayName)
			}
			cidrBlocks = append(cidrBlocks, cidrArgs)
		}
		manArgs.CidrBlocks = cidrBlocks
		args.MasterAuthorizedNetworksConfig = manArgs
	}

	if spec.ControlPlaneEndpoints != nil {
		args.ControlPlaneEndpointsConfig = &container.ClusterControlPlaneEndpointsConfigArgs{
			DnsEndpointConfig: &container.ClusterControlPlaneEndpointsConfigDnsEndpointConfigArgs{
				AllowExternalTraffic: pulumi.BoolPtr(spec.ControlPlaneEndpoints.DnsEndpointAllowExternalTraffic),
			},
			IpEndpointsConfig: &container.ClusterControlPlaneEndpointsConfigIpEndpointsConfigArgs{
				Enabled: pulumi.BoolPtr(spec.ControlPlaneEndpoints.GetIpEndpointsEnabled()),
			},
		}
	}

	if spec.MinMasterVersion != "" {
		args.MinMasterVersion = pulumi.StringPtr(spec.MinMasterVersion)
	}

	if spec.MaintenancePolicy != nil {
		maintenanceArgs := &container.ClusterMaintenancePolicyArgs{}
		if spec.MaintenancePolicy.DailyWindow != nil {
			maintenanceArgs.DailyMaintenanceWindow = &container.ClusterMaintenancePolicyDailyMaintenanceWindowArgs{
				StartTime: pulumi.String(spec.MaintenancePolicy.DailyWindow.StartTime),
			}
		}
		if spec.MaintenancePolicy.RecurringWindow != nil {
			maintenanceArgs.RecurringWindow = &container.ClusterMaintenancePolicyRecurringWindowArgs{
				StartTime:  pulumi.String(spec.MaintenancePolicy.RecurringWindow.StartTime),
				EndTime:    pulumi.String(spec.MaintenancePolicy.RecurringWindow.EndTime),
				Recurrence: pulumi.String(spec.MaintenancePolicy.RecurringWindow.Recurrence),
			}
		}
		if len(spec.MaintenancePolicy.Exclusions) > 0 {
			exclusions := container.ClusterMaintenancePolicyMaintenanceExclusionArray{}
			for _, exclusion := range spec.MaintenancePolicy.Exclusions {
				exclusionArgs := &container.ClusterMaintenancePolicyMaintenanceExclusionArgs{
					ExclusionName: pulumi.String(exclusion.ExclusionName),
					StartTime:     pulumi.String(exclusion.StartTime),
					EndTime:       pulumi.String(exclusion.EndTime),
				}
				if exclusion.Scope != "" {
					exclusionArgs.ExclusionOptions = &container.ClusterMaintenancePolicyMaintenanceExclusionExclusionOptionsArgs{
						Scope: pulumi.String(exclusion.Scope),
					}
				}
				exclusions = append(exclusions, exclusionArgs)
			}
			maintenanceArgs.MaintenanceExclusions = exclusions
		}
		args.MaintenancePolicy = maintenanceArgs
	}

	// Node auto-provisioning: GKE creates/deletes node pools within the
	// resource limits — bounded by spec-level validation so an enabled NAP
	// always carries limits (an unbounded NAP is an unbounded bill).
	if spec.ClusterAutoscaling != nil {
		autoscalingArgs := &container.ClusterClusterAutoscalingArgs{
			Enabled:            pulumi.BoolPtr(spec.ClusterAutoscaling.Enabled),
			AutoscalingProfile: pulumi.StringPtr(spec.ClusterAutoscaling.GetAutoscalingProfile()),
		}
		if len(spec.ClusterAutoscaling.AutoProvisioningLocations) > 0 {
			autoscalingArgs.AutoProvisioningLocations = pulumi.ToStringArray(spec.ClusterAutoscaling.AutoProvisioningLocations)
		}
		if len(spec.ClusterAutoscaling.ResourceLimits) > 0 {
			limits := container.ClusterClusterAutoscalingResourceLimitArray{}
			for _, limit := range spec.ClusterAutoscaling.ResourceLimits {
				limits = append(limits, &container.ClusterClusterAutoscalingResourceLimitArgs{
					ResourceType: pulumi.String(limit.ResourceType),
					Minimum:      pulumi.IntPtr(int(limit.Minimum)),
					Maximum:      pulumi.Int(int(limit.Maximum)),
				})
			}
			autoscalingArgs.ResourceLimits = limits
		}
		if defaults := spec.ClusterAutoscaling.AutoProvisioningDefaults; defaults != nil {
			defaultsArgs := &container.ClusterClusterAutoscalingAutoProvisioningDefaultsArgs{
				ShieldedInstanceConfig: &container.ClusterClusterAutoscalingAutoProvisioningDefaultsShieldedInstanceConfigArgs{
					EnableSecureBoot:          pulumi.BoolPtr(defaults.EnableSecureBoot),
					EnableIntegrityMonitoring: pulumi.BoolPtr(defaults.GetEnableIntegrityMonitoring()),
				},
				Management: &container.ClusterClusterAutoscalingAutoProvisioningDefaultsManagementArgs{
					AutoUpgrade: pulumi.BoolPtr(defaults.GetAutoUpgrade()),
					AutoRepair:  pulumi.BoolPtr(defaults.GetAutoRepair()),
				},
			}
			if defaults.ServiceAccount.GetValue() != "" {
				defaultsArgs.ServiceAccount = pulumi.StringPtr(defaults.ServiceAccount.GetValue())
			}
			if len(defaults.OauthScopes) > 0 {
				defaultsArgs.OauthScopes = pulumi.ToStringArray(defaults.OauthScopes)
			}
			if defaults.DiskSizeGb != nil {
				defaultsArgs.DiskSize = pulumi.IntPtr(int(defaults.GetDiskSizeGb()))
			}
			if defaults.DiskType != "" {
				defaultsArgs.DiskType = pulumi.StringPtr(defaults.DiskType)
			}
			if defaults.ImageType != "" {
				defaultsArgs.ImageType = pulumi.StringPtr(defaults.ImageType)
			}
			if defaults.MinCpuPlatform != "" {
				defaultsArgs.MinCpuPlatform = pulumi.StringPtr(defaults.MinCpuPlatform)
			}
			if defaults.BootDiskKmsKey.GetValue() != "" {
				defaultsArgs.BootDiskKmsKey = pulumi.StringPtr(defaults.BootDiskKmsKey.GetValue())
			}
			autoscalingArgs.AutoProvisioningDefaults = defaultsArgs
		}
		args.ClusterAutoscaling = autoscalingArgs
	}

	if spec.EnableVerticalPodAutoscaling {
		args.VerticalPodAutoscaling = &container.ClusterVerticalPodAutoscalingArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	if spec.HpaProfile != "" {
		args.PodAutoscaling = &container.ClusterPodAutoscalingArgs{
			HpaProfile: pulumi.String(spec.HpaProfile),
		}
	}

	// Workload Identity Federation for GKE: the pool name is fixed by the
	// API to PROJECT_ID.svc.id.goog. Autopilot clusters have it always on —
	// the block is suppressed there to avoid a permanent diff.
	workloadIdentityEnabled := spec.GetWorkloadIdentityEnabled()
	if workloadIdentityEnabled && !spec.EnableAutopilot {
		args.WorkloadIdentityConfig = &container.ClusterWorkloadIdentityConfigArgs{
			WorkloadPool: pulumi.StringPtr(workloadPool),
		}
	}

	// Shielded GKE nodes: only forwarded when the spec sets it (GCP default
	// is true; Autopilot rejects the field and spec validation blocks it
	// there).
	if spec.EnableShieldedNodes != nil {
		args.EnableShieldedNodes = pulumi.BoolPtr(spec.GetEnableShieldedNodes())
	}

	if spec.DatabaseEncryption != nil {
		encryptionArgs := &container.ClusterDatabaseEncryptionArgs{
			State: pulumi.String(spec.DatabaseEncryption.State),
		}
		if spec.DatabaseEncryption.KeyName.GetValue() != "" {
			encryptionArgs.KeyName = pulumi.StringPtr(spec.DatabaseEncryption.KeyName.GetValue())
		}
		args.DatabaseEncryption = encryptionArgs
	}

	if spec.BinaryAuthorizationEvaluationMode != "" {
		args.BinaryAuthorization = &container.ClusterBinaryAuthorizationArgs{
			EvaluationMode: pulumi.StringPtr(spec.BinaryAuthorizationEvaluationMode),
		}
	}

	if spec.SecurityPosture != nil {
		postureArgs := &container.ClusterSecurityPostureConfigArgs{}
		if spec.SecurityPosture.Mode != "" {
			postureArgs.Mode = pulumi.StringPtr(spec.SecurityPosture.Mode)
		}
		if spec.SecurityPosture.VulnerabilityMode != "" {
			postureArgs.VulnerabilityMode = pulumi.StringPtr(spec.SecurityPosture.VulnerabilityMode)
		}
		args.SecurityPostureConfig = postureArgs
	}

	if spec.AuthenticatorSecurityGroup != "" {
		args.AuthenticatorGroupsConfig = &container.ClusterAuthenticatorGroupsConfigArgs{
			SecurityGroup: pulumi.String(spec.AuthenticatorSecurityGroup),
		}
	}

	if spec.EnableLegacyAbac {
		args.EnableLegacyAbac = pulumi.BoolPtr(true)
	}

	if spec.EnableMeshCertificates {
		args.MeshCertificates = &container.ClusterMeshCertificatesArgs{
			EnableCertificates: pulumi.Bool(true),
		}
	}

	if spec.EnableSecretManagerCsi {
		args.SecretManagerConfig = &container.ClusterSecretManagerConfigArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	if spec.ConfidentialNodes != nil {
		confidentialArgs := &container.ClusterConfidentialNodesArgs{
			Enabled: pulumi.Bool(spec.ConfidentialNodes.Enabled),
		}
		if spec.ConfidentialNodes.ConfidentialInstanceType != "" {
			confidentialArgs.ConfidentialInstanceType = pulumi.StringPtr(spec.ConfidentialNodes.ConfidentialInstanceType)
		}
		args.ConfidentialNodes = confidentialArgs
	}

	if spec.AnonymousAuthenticationMode != "" {
		args.AnonymousAuthenticationConfig = &container.ClusterAnonymousAuthenticationConfigArgs{
			Mode: pulumi.String(spec.AnonymousAuthenticationMode),
		}
	}

	if spec.EnableIdentityService {
		args.IdentityServiceConfig = &container.ClusterIdentityServiceConfigArgs{
			Enabled: pulumi.BoolPtr(true),
		}
	}

	// An explicit empty components list is meaningful: it disables the
	// Cloud Logging/Monitoring integration outright, so the spec message's
	// presence (not emptiness) drives whether the block is emitted.
	if spec.Logging != nil {
		args.LoggingConfig = &container.ClusterLoggingConfigArgs{
			EnableComponents: pulumi.ToStringArray(spec.Logging.Components),
		}
	}

	if spec.Monitoring != nil {
		monitoringArgs := &container.ClusterMonitoringConfigArgs{}
		if len(spec.Monitoring.Components) > 0 {
			monitoringArgs.EnableComponents = pulumi.ToStringArray(spec.Monitoring.Components)
		}
		prometheusArgs := &container.ClusterMonitoringConfigManagedPrometheusArgs{
			Enabled: pulumi.Bool(spec.Monitoring.GetManagedPrometheusEnabled()),
		}
		if spec.Monitoring.AutoMonitoringScope != "" {
			prometheusArgs.AutoMonitoringConfig = &container.ClusterMonitoringConfigManagedPrometheusAutoMonitoringConfigArgs{
				Scope: pulumi.String(spec.Monitoring.AutoMonitoringScope),
			}
		}
		monitoringArgs.ManagedPrometheus = prometheusArgs
		if spec.Monitoring.AdvancedDatapathMetricsEnabled || spec.Monitoring.AdvancedDatapathRelayEnabled {
			monitoringArgs.AdvancedDatapathObservabilityConfig = &container.ClusterMonitoringConfigAdvancedDatapathObservabilityConfigArgs{
				EnableMetrics: pulumi.Bool(spec.Monitoring.AdvancedDatapathMetricsEnabled),
				EnableRelay:   pulumi.Bool(spec.Monitoring.AdvancedDatapathRelayEnabled),
			}
		}
		args.MonitoringConfig = monitoringArgs
	}

	if spec.NotificationPubsub != nil {
		pubsubArgs := &container.ClusterNotificationConfigPubsubArgs{
			Enabled: pulumi.Bool(spec.NotificationPubsub.Enabled),
		}
		if spec.NotificationPubsub.Topic.GetValue() != "" {
			pubsubArgs.Topic = pulumi.StringPtr(spec.NotificationPubsub.Topic.GetValue())
		}
		if len(spec.NotificationPubsub.EventTypes) > 0 {
			pubsubArgs.Filter = &container.ClusterNotificationConfigPubsubFilterArgs{
				EventTypes: pulumi.ToStringArray(spec.NotificationPubsub.EventTypes),
			}
		}
		args.NotificationConfig = &container.ClusterNotificationConfigArgs{
			Pubsub: pubsubArgs,
		}
	}

	if spec.EnableCostManagement {
		args.CostManagementConfig = &container.ClusterCostManagementConfigArgs{
			Enabled: pulumi.Bool(true),
		}
	}

	if spec.ResourceUsageExport != nil {
		args.ResourceUsageExportConfig = &container.ClusterResourceUsageExportConfigArgs{
			EnableNetworkEgressMetering:       pulumi.BoolPtr(spec.ResourceUsageExport.EnableNetworkEgressMetering),
			EnableResourceConsumptionMetering: pulumi.BoolPtr(spec.ResourceUsageExport.GetEnableResourceConsumptionMetering()),
			BigqueryDestination: &container.ClusterResourceUsageExportConfigBigqueryDestinationArgs{
				DatasetId: pulumi.String(spec.ResourceUsageExport.BigqueryDatasetId.GetValue()),
			},
		}
	}

	// Addons: emitted when the spec configures them or when Calico network
	// policy needs its companion addon. The network-policy addon toggle
	// always mirrors enable_network_policy (never set on Autopilot).
	if spec.Addons != nil || (spec.EnableNetworkPolicy && !spec.EnableAutopilot) {
		addonsArgs := &container.ClusterAddonsConfigArgs{}
		if !spec.EnableAutopilot {
			addonsArgs.NetworkPolicyConfig = &container.ClusterAddonsConfigNetworkPolicyConfigArgs{
				Disabled: pulumi.Bool(!spec.EnableNetworkPolicy),
			}
		}
		httpLoadBalancingEnabled := true
		horizontalPodAutoscalingEnabled := true
		pdCsiEnabled := true
		if spec.Addons != nil {
			httpLoadBalancingEnabled = spec.Addons.GetHttpLoadBalancingEnabled()
			horizontalPodAutoscalingEnabled = spec.Addons.GetHorizontalPodAutoscalingEnabled()
			pdCsiEnabled = spec.Addons.GetGcePersistentDiskCsiDriverEnabled()
		}
		addonsArgs.HttpLoadBalancing = &container.ClusterAddonsConfigHttpLoadBalancingArgs{
			Disabled: pulumi.Bool(!httpLoadBalancingEnabled),
		}
		addonsArgs.HorizontalPodAutoscaling = &container.ClusterAddonsConfigHorizontalPodAutoscalingArgs{
			Disabled: pulumi.Bool(!horizontalPodAutoscalingEnabled),
		}
		addonsArgs.GcePersistentDiskCsiDriverConfig = &container.ClusterAddonsConfigGcePersistentDiskCsiDriverConfigArgs{
			Enabled: pulumi.Bool(pdCsiEnabled),
		}
		if spec.Addons != nil {
			if spec.Addons.GcpFilestoreCsiDriverEnabled {
				addonsArgs.GcpFilestoreCsiDriverConfig = &container.ClusterAddonsConfigGcpFilestoreCsiDriverConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.GcsFuseCsiDriverEnabled {
				addonsArgs.GcsFuseCsiDriverConfig = &container.ClusterAddonsConfigGcsFuseCsiDriverConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.GkeBackupAgentEnabled {
				addonsArgs.GkeBackupAgentConfig = &container.ClusterAddonsConfigGkeBackupAgentConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.DnsCacheEnabled {
				addonsArgs.DnsCacheConfig = &container.ClusterAddonsConfigDnsCacheConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.ConfigConnectorEnabled {
				addonsArgs.ConfigConnectorConfig = &container.ClusterAddonsConfigConfigConnectorConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.StatefulHaEnabled {
				addonsArgs.StatefulHaConfig = &container.ClusterAddonsConfigStatefulHaConfigArgs{
					Enabled: pulumi.Bool(true),
				}
			}
			if spec.Addons.RayOperatorEnabled {
				addonsArgs.RayOperatorConfigs = container.ClusterAddonsConfigRayOperatorConfigArray{
					&container.ClusterAddonsConfigRayOperatorConfigArgs{
						Enabled: pulumi.Bool(true),
					},
				}
			}
		}
		args.AddonsConfig = addonsArgs
	}

	if spec.FleetProject != "" {
		args.Fleet = &container.ClusterFleetArgs{
			Project: pulumi.StringPtr(spec.FleetProject),
		}
	}

	createdCluster, err := container.NewCluster(ctx, "cluster", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create container cluster")
	}

	ctx.Export(OpEndpoint, createdCluster.Endpoint)
	// The CA certificate is public trust material (clients install it as a
	// trust anchor), not a secret.
	ctx.Export(OpClusterCaCertificate, createdCluster.MasterAuth.ClusterCaCertificate())
	if workloadIdentityEnabled || spec.EnableAutopilot {
		ctx.Export(OpWorkloadIdentityPool, pulumi.String(workloadPool))
	} else {
		ctx.Export(OpWorkloadIdentityPool, pulumi.String(""))
	}
	ctx.Export(OpClusterId, createdCluster.ID())
	ctx.Export(OpName, createdCluster.Name)
	// The plain spec location (not the provider's computed attribute) so
	// both engines emit the identical value for API callers.
	ctx.Export(OpLocation, pulumi.String(spec.Location))
	ctx.Export(OpSelfLink, createdCluster.SelfLink)
	ctx.Export(OpMasterVersion, createdCluster.MasterVersion)

	return nil
}
