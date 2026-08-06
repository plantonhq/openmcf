package module

import (
	"github.com/pkg/errors"
	gcpdataprocclusterv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpdataproccluster/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/dataproc"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// dataprocCluster provisions the Dataproc cluster. Two mutually
// exclusive arms mirror the API: cluster_config provisions dedicated
// Compute Engine VMs; virtual_cluster_config runs Dataproc as pods on
// an existing GKE cluster. Omitting both creates a default GCE cluster
// (2 workers, default machine types).
//
// Mutability: the cluster is create-mostly-immutable. The only in-place
// updates the API supports are labels, primary/secondary worker counts
// (manual scaling), min_num_instances, the autoscaling-policy
// attachment, and the lifecycle TTLs — everything else forces
// recreation. The virtual arm has no update paths at all.
func dataprocCluster(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpDataprocCluster.Spec

	// Enable the Dataproc API — the control plane that owns the cluster.
	// disable_on_destroy stays false: tearing down one cluster must never
	// disable the API for everything else in the project.
	dataprocApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("dataproc.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	if spec.ProjectId.GetValue() != "" {
		dataprocApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdDataprocApi, err := projects.NewService(ctx,
		"dpc-dataproc.googleapis.com", dataprocApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable dataproc.googleapis.com api")
	}

	args := &dataproc.ClusterArgs{
		Name:   pulumi.String(spec.ClusterName),
		Region: pulumi.StringPtr(spec.Region),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project (omitting the argument lets the provider
	// resolve its own project).
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// The Dataproc API rejects user labels on virtual (GKE-based)
	// clusters — labels (including the platform attribution set) are sent
	// only for the GCE arm, identically to the Terraform module.
	if spec.VirtualClusterConfig == nil {
		args.Labels = pulumi.ToStringMap(locals.GcpLabels)
	}

	// Applied when worker counts shrink during an update: YARN drains
	// running tasks for up to this window before nodes are removed.
	if spec.GracefulDecommissionTimeout != "" {
		args.GracefulDecommissionTimeout = pulumi.StringPtr(spec.GracefulDecommissionTimeout)
	}

	if spec.ClusterConfig != nil {
		args.ClusterConfig = buildClusterConfig(spec.ClusterConfig)
	}

	if spec.VirtualClusterConfig != nil {
		args.VirtualClusterConfig = buildVirtualClusterConfig(spec.VirtualClusterConfig)
	}

	createdCluster, err := dataproc.NewCluster(ctx, "dataproc-cluster", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdDataprocApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create dataproc cluster")
	}

	// The resource ID is the fully qualified cluster resource name
	// (projects/{p}/regions/{r}/clusters/{c}) — the bridged provider
	// inherits Terraform's ID format, so both engines export the exact
	// path downstream composition (e.g. spark_history_server_config)
	// consumes.
	ctx.Export(OpClusterId, createdCluster.ID())
	ctx.Export(OpClusterName, createdCluster.Name)

	// The staging bucket actually in use: the user-supplied bucket when
	// one was referenced, otherwise the bucket GCP auto-created. The
	// virtual arm reports its own staging bucket the same way.
	if spec.VirtualClusterConfig != nil {
		ctx.Export(OpStagingBucket, createdCluster.VirtualClusterConfig.StagingBucket().Elem())
	} else {
		ctx.Export(OpStagingBucket, createdCluster.ClusterConfig.Bucket().Elem())
	}

	return nil
}

// buildDiskConfig converts the shared disk-config message for any node
// group into the given block constructor via the supplied setter.
func diskFields(d *gcpdataprocclusterv1alpha1.GcpDataprocClusterDiskConfig) (bootDiskSizeGb *int, bootDiskType *string, numLocalSsds *int, localSsdInterface *string) {
	if d.BootDiskSizeGb > 0 {
		v := int(d.BootDiskSizeGb)
		bootDiskSizeGb = &v
	}
	if d.BootDiskType != "" {
		v := d.BootDiskType
		bootDiskType = &v
	}
	if d.NumLocalSsds > 0 {
		v := int(d.NumLocalSsds)
		numLocalSsds = &v
	}
	if d.LocalSsdInterface != "" {
		v := d.LocalSsdInterface
		localSsdInterface = &v
	}
	return
}

// buildClusterConfig assembles the standard Compute Engine arm.
func buildClusterConfig(cfg *gcpdataprocclusterv1alpha1.GcpDataprocClusterConfig) *dataproc.ClusterClusterConfigArgs {
	clusterConfig := &dataproc.ClusterClusterConfigArgs{}

	// Bucket names arrive as resolved references. GCP auto-creates
	// staging/temp buckets when these are unset.
	if cfg.StagingBucket.GetValue() != "" {
		clusterConfig.StagingBucket = pulumi.StringPtr(cfg.StagingBucket.GetValue())
	}
	if cfg.TempBucket.GetValue() != "" {
		clusterConfig.TempBucket = pulumi.StringPtr(cfg.TempBucket.GetValue())
	}
	if cfg.ClusterTier != "" {
		clusterConfig.ClusterTier = pulumi.StringPtr(cfg.ClusterTier)
	}

	// ── GCE environment: networking, identity, hardening, placement ──

	if cfg.GceConfig != nil {
		gce := cfg.GceConfig
		gceArgs := &dataproc.ClusterClusterConfigGceClusterConfigArgs{}

		if gce.Network.GetValue() != "" {
			gceArgs.Network = pulumi.StringPtr(gce.Network.GetValue())
		}
		if gce.Subnetwork.GetValue() != "" {
			gceArgs.Subnetwork = pulumi.StringPtr(gce.Subnetwork.GetValue())
		}
		if gce.ServiceAccount.GetValue() != "" {
			gceArgs.ServiceAccount = pulumi.StringPtr(gce.ServiceAccount.GetValue())
		}
		if len(gce.ServiceAccountScopes) > 0 {
			gceArgs.ServiceAccountScopes = pulumi.ToStringArray(gce.ServiceAccountScopes)
		}
		if gce.Zone != "" {
			gceArgs.Zone = pulumi.StringPtr(gce.Zone)
		}
		if gce.InternalIpOnly {
			gceArgs.InternalIpOnly = pulumi.BoolPtr(true)
		}
		if len(gce.Tags) > 0 {
			gceArgs.Tags = pulumi.ToStringArray(gce.Tags)
		}
		if len(gce.Metadata) > 0 {
			gceArgs.Metadata = pulumi.ToStringMap(gce.Metadata)
		}

		if gce.ShieldedInstanceConfig != nil {
			gceArgs.ShieldedInstanceConfig = &dataproc.ClusterClusterConfigGceClusterConfigShieldedInstanceConfigArgs{
				EnableSecureBoot:          pulumi.BoolPtr(gce.ShieldedInstanceConfig.EnableSecureBoot),
				EnableVtpm:                pulumi.BoolPtr(gce.ShieldedInstanceConfig.EnableVtpm),
				EnableIntegrityMonitoring: pulumi.BoolPtr(gce.ShieldedInstanceConfig.EnableIntegrityMonitoring),
			}
		}

		if gce.ReservationAffinity != nil {
			reservationArgs := &dataproc.ClusterClusterConfigGceClusterConfigReservationAffinityArgs{}
			if gce.ReservationAffinity.ConsumeReservationType != "" {
				reservationArgs.ConsumeReservationType = pulumi.StringPtr(gce.ReservationAffinity.ConsumeReservationType)
			}
			if gce.ReservationAffinity.Key != "" {
				reservationArgs.Key = pulumi.StringPtr(gce.ReservationAffinity.Key)
			}
			if len(gce.ReservationAffinity.Values) > 0 {
				reservationArgs.Values = pulumi.ToStringArray(gce.ReservationAffinity.Values)
			}
			gceArgs.ReservationAffinity = reservationArgs
		}

		if gce.NodeGroupAffinity != nil {
			gceArgs.NodeGroupAffinity = &dataproc.ClusterClusterConfigGceClusterConfigNodeGroupAffinityArgs{
				NodeGroupUri: pulumi.String(gce.NodeGroupAffinity.NodeGroupUri),
			}
		}

		if gce.ConfidentialInstanceConfig != nil {
			gceArgs.ConfidentialInstanceConfig = &dataproc.ClusterClusterConfigGceClusterConfigConfidentialInstanceConfigArgs{
				EnableConfidentialCompute: pulumi.BoolPtr(gce.ConfidentialInstanceConfig.EnableConfidentialCompute),
			}
		}

		clusterConfig.GceClusterConfig = gceArgs
	}

	// ── Master node group ──

	if cfg.MasterConfig != nil {
		m := cfg.MasterConfig
		masterArgs := &dataproc.ClusterClusterConfigMasterConfigArgs{}

		if m.NumInstances > 0 {
			masterArgs.NumInstances = pulumi.IntPtr(int(m.NumInstances))
		}
		if m.MachineType != "" {
			masterArgs.MachineType = pulumi.StringPtr(m.MachineType)
		}
		if m.MinCpuPlatform != "" {
			masterArgs.MinCpuPlatform = pulumi.StringPtr(m.MinCpuPlatform)
		}
		if m.ImageUri != "" {
			masterArgs.ImageUri = pulumi.StringPtr(m.ImageUri)
		}
		if m.DiskConfig != nil {
			sizeGb, diskType, ssds, ssdIface := diskFields(m.DiskConfig)
			diskArgs := &dataproc.ClusterClusterConfigMasterConfigDiskConfigArgs{}
			if sizeGb != nil {
				diskArgs.BootDiskSizeGb = pulumi.IntPtr(*sizeGb)
			}
			if diskType != nil {
				diskArgs.BootDiskType = pulumi.StringPtr(*diskType)
			}
			if ssds != nil {
				diskArgs.NumLocalSsds = pulumi.IntPtr(*ssds)
			}
			if ssdIface != nil {
				diskArgs.LocalSsdInterface = pulumi.StringPtr(*ssdIface)
			}
			masterArgs.DiskConfig = diskArgs
		}
		if len(m.Accelerators) > 0 {
			var accels dataproc.ClusterClusterConfigMasterConfigAcceleratorArray
			for _, a := range m.Accelerators {
				accels = append(accels, &dataproc.ClusterClusterConfigMasterConfigAcceleratorArgs{
					AcceleratorType:  pulumi.String(a.AcceleratorType),
					AcceleratorCount: pulumi.Int(int(a.AcceleratorCount)),
				})
			}
			masterArgs.Accelerators = accels
		}

		clusterConfig.MasterConfig = masterArgs
	}

	// ── Primary worker node group ──
	// num_instances / min_num_instances are the manual-scaling levers —
	// the only node counts that update in place (with
	// graceful_decommission_timeout honored on shrink).

	if cfg.WorkerConfig != nil {
		w := cfg.WorkerConfig
		workerArgs := &dataproc.ClusterClusterConfigWorkerConfigArgs{}

		if w.NumInstances > 0 {
			workerArgs.NumInstances = pulumi.IntPtr(int(w.NumInstances))
		}
		if w.MachineType != "" {
			workerArgs.MachineType = pulumi.StringPtr(w.MachineType)
		}
		if w.MinCpuPlatform != "" {
			workerArgs.MinCpuPlatform = pulumi.StringPtr(w.MinCpuPlatform)
		}
		if w.ImageUri != "" {
			workerArgs.ImageUri = pulumi.StringPtr(w.ImageUri)
		}
		if w.MinNumInstances > 0 {
			workerArgs.MinNumInstances = pulumi.IntPtr(int(w.MinNumInstances))
		}
		if w.DiskConfig != nil {
			sizeGb, diskType, ssds, ssdIface := diskFields(w.DiskConfig)
			diskArgs := &dataproc.ClusterClusterConfigWorkerConfigDiskConfigArgs{}
			if sizeGb != nil {
				diskArgs.BootDiskSizeGb = pulumi.IntPtr(*sizeGb)
			}
			if diskType != nil {
				diskArgs.BootDiskType = pulumi.StringPtr(*diskType)
			}
			if ssds != nil {
				diskArgs.NumLocalSsds = pulumi.IntPtr(*ssds)
			}
			if ssdIface != nil {
				diskArgs.LocalSsdInterface = pulumi.StringPtr(*ssdIface)
			}
			workerArgs.DiskConfig = diskArgs
		}
		if len(w.Accelerators) > 0 {
			var accels dataproc.ClusterClusterConfigWorkerConfigAcceleratorArray
			for _, a := range w.Accelerators {
				accels = append(accels, &dataproc.ClusterClusterConfigWorkerConfigAcceleratorArgs{
					AcceleratorType:  pulumi.String(a.AcceleratorType),
					AcceleratorCount: pulumi.Int(int(a.AcceleratorCount)),
				})
			}
			workerArgs.Accelerators = accels
		}

		clusterConfig.WorkerConfig = workerArgs
	}

	// ── Secondary (preemptible/spot) worker group ──
	// The count updates in place; preemptibility is immutable. Machine
	// shape is inherited from the primary workers unless the flexibility
	// policy overrides it.

	if cfg.SecondaryWorkerConfig != nil {
		s := cfg.SecondaryWorkerConfig
		secondaryArgs := &dataproc.ClusterClusterConfigPreemptibleWorkerConfigArgs{}

		if s.NumInstances > 0 {
			secondaryArgs.NumInstances = pulumi.IntPtr(int(s.NumInstances))
		}
		if s.Preemptibility != "" {
			secondaryArgs.Preemptibility = pulumi.StringPtr(s.Preemptibility)
		}
		if s.DiskConfig != nil {
			sizeGb, diskType, ssds, ssdIface := diskFields(s.DiskConfig)
			diskArgs := &dataproc.ClusterClusterConfigPreemptibleWorkerConfigDiskConfigArgs{}
			if sizeGb != nil {
				diskArgs.BootDiskSizeGb = pulumi.IntPtr(*sizeGb)
			}
			if diskType != nil {
				diskArgs.BootDiskType = pulumi.StringPtr(*diskType)
			}
			if ssds != nil {
				diskArgs.NumLocalSsds = pulumi.IntPtr(*ssds)
			}
			if ssdIface != nil {
				diskArgs.LocalSsdInterface = pulumi.StringPtr(*ssdIface)
			}
			secondaryArgs.DiskConfig = diskArgs
		}

		if s.InstanceFlexibilityPolicy != nil {
			flexArgs := &dataproc.ClusterClusterConfigPreemptibleWorkerConfigInstanceFlexibilityPolicyArgs{}

			if len(s.InstanceFlexibilityPolicy.InstanceSelectionList) > 0 {
				var selections dataproc.ClusterClusterConfigPreemptibleWorkerConfigInstanceFlexibilityPolicyInstanceSelectionListArray
				for _, sel := range s.InstanceFlexibilityPolicy.InstanceSelectionList {
					selections = append(selections, &dataproc.ClusterClusterConfigPreemptibleWorkerConfigInstanceFlexibilityPolicyInstanceSelectionListArgs{
						MachineTypes: pulumi.ToStringArray(sel.MachineTypes),
						Rank:         pulumi.IntPtr(int(sel.Rank)),
					})
				}
				flexArgs.InstanceSelectionLists = selections
			}

			if s.InstanceFlexibilityPolicy.ProvisioningModelMix != nil {
				flexArgs.ProvisioningModelMix = &dataproc.ClusterClusterConfigPreemptibleWorkerConfigInstanceFlexibilityPolicyProvisioningModelMixArgs{
					StandardCapacityBase:             pulumi.IntPtr(int(s.InstanceFlexibilityPolicy.ProvisioningModelMix.StandardCapacityBase)),
					StandardCapacityPercentAboveBase: pulumi.IntPtr(int(s.InstanceFlexibilityPolicy.ProvisioningModelMix.StandardCapacityPercentAboveBase)),
				}
			}

			secondaryArgs.InstanceFlexibilityPolicy = flexArgs
		}

		clusterConfig.PreemptibleWorkerConfig = secondaryArgs
	}

	// ── Software config ──
	// The spec's `properties` map feeds the provider's
	// override_properties — the API's writable surface (the provider's
	// `properties` attribute is the computed resolved set).

	if cfg.SoftwareConfig != nil {
		sw := cfg.SoftwareConfig
		softwareArgs := &dataproc.ClusterClusterConfigSoftwareConfigArgs{}

		if sw.ImageVersion != "" {
			softwareArgs.ImageVersion = pulumi.StringPtr(sw.ImageVersion)
		}
		if len(sw.OptionalComponents) > 0 {
			softwareArgs.OptionalComponents = pulumi.ToStringArray(sw.OptionalComponents)
		}
		if len(sw.Properties) > 0 {
			softwareArgs.OverrideProperties = pulumi.ToStringMap(sw.Properties)
		}

		clusterConfig.SoftwareConfig = softwareArgs
	}

	// ── Initialization actions ──

	if len(cfg.InitializationActions) > 0 {
		var initActions dataproc.ClusterClusterConfigInitializationActionArray
		for _, action := range cfg.InitializationActions {
			initArgs := &dataproc.ClusterClusterConfigInitializationActionArgs{
				Script: pulumi.String(action.Script),
			}
			if action.TimeoutSec > 0 {
				initArgs.TimeoutSec = pulumi.IntPtr(int(action.TimeoutSec))
			}
			initActions = append(initActions, initArgs)
		}
		clusterConfig.InitializationActions = initActions
	}

	// ── Autoscaling policy attachment ──
	// A first-class resource referenced by its full resource name;
	// attach/swap/detach updates in place.

	if cfg.AutoscalingPolicyUri.GetValue() != "" {
		clusterConfig.AutoscalingConfig = &dataproc.ClusterClusterConfigAutoscalingConfigArgs{
			PolicyUri: pulumi.String(cfg.AutoscalingPolicyUri.GetValue()),
		}
	}

	// ── CMEK for all persistent disks (key change forces recreation) ──

	if cfg.EncryptionKmsKeyName.GetValue() != "" {
		clusterConfig.EncryptionConfig = &dataproc.ClusterClusterConfigEncryptionConfigArgs{
			KmsKeyName: pulumi.String(cfg.EncryptionKmsKeyName.GetValue()),
		}
	}

	// ── Kerberos XOR personal-cluster identity mapping ──
	// Kerberos secret fields are GCS URIs of KMS-encrypted files — never
	// inline material (the API's own contract).

	if cfg.SecurityConfig != nil {
		securityArgs := &dataproc.ClusterClusterConfigSecurityConfigArgs{}

		if k := cfg.SecurityConfig.KerberosConfig; k != nil {
			kerberosArgs := &dataproc.ClusterClusterConfigSecurityConfigKerberosConfigArgs{
				RootPrincipalPasswordUri: pulumi.String(k.RootPrincipalPasswordUri),
				KmsKeyUri:                pulumi.String(k.KmsKeyUri.GetValue()),
			}
			if k.EnableKerberos {
				kerberosArgs.EnableKerberos = pulumi.BoolPtr(true)
			}
			if k.Realm != "" {
				kerberosArgs.Realm = pulumi.StringPtr(k.Realm)
			}
			if k.TgtLifetimeHours > 0 {
				kerberosArgs.TgtLifetimeHours = pulumi.IntPtr(int(k.TgtLifetimeHours))
			}
			if k.KdcDbKeyUri != "" {
				kerberosArgs.KdcDbKeyUri = pulumi.StringPtr(k.KdcDbKeyUri)
			}
			if k.KeystoreUri != "" {
				kerberosArgs.KeystoreUri = pulumi.StringPtr(k.KeystoreUri)
			}
			if k.KeystorePasswordUri != "" {
				kerberosArgs.KeystorePasswordUri = pulumi.StringPtr(k.KeystorePasswordUri)
			}
			if k.KeyPasswordUri != "" {
				kerberosArgs.KeyPasswordUri = pulumi.StringPtr(k.KeyPasswordUri)
			}
			if k.TruststoreUri != "" {
				kerberosArgs.TruststoreUri = pulumi.StringPtr(k.TruststoreUri)
			}
			if k.TruststorePasswordUri != "" {
				kerberosArgs.TruststorePasswordUri = pulumi.StringPtr(k.TruststorePasswordUri)
			}
			if k.CrossRealmTrustRealm != "" {
				kerberosArgs.CrossRealmTrustRealm = pulumi.StringPtr(k.CrossRealmTrustRealm)
			}
			if k.CrossRealmTrustKdc != "" {
				kerberosArgs.CrossRealmTrustKdc = pulumi.StringPtr(k.CrossRealmTrustKdc)
			}
			if k.CrossRealmTrustAdminServer != "" {
				kerberosArgs.CrossRealmTrustAdminServer = pulumi.StringPtr(k.CrossRealmTrustAdminServer)
			}
			if k.CrossRealmTrustSharedPasswordUri != "" {
				kerberosArgs.CrossRealmTrustSharedPasswordUri = pulumi.StringPtr(k.CrossRealmTrustSharedPasswordUri)
			}
			securityArgs.KerberosConfig = kerberosArgs
		}

		if id := cfg.SecurityConfig.IdentityConfig; id != nil {
			securityArgs.IdentityConfig = &dataproc.ClusterClusterConfigSecurityConfigIdentityConfigArgs{
				UserServiceAccountMapping: pulumi.ToStringMap(id.UserServiceAccountMapping),
			}
		}

		clusterConfig.SecurityConfig = securityArgs
	}

	// ── Component Gateway (authenticated web UIs) ──

	if cfg.EndpointConfig != nil {
		clusterConfig.EndpointConfig = &dataproc.ClusterClusterConfigEndpointConfigArgs{
			EnableHttpPortAccess: pulumi.Bool(cfg.EndpointConfig.EnableHttpPortAccess),
		}
	}

	// ── Cost-control TTLs — both update in place ──

	if cfg.LifecycleConfig != nil {
		lifecycleArgs := &dataproc.ClusterClusterConfigLifecycleConfigArgs{}
		if cfg.LifecycleConfig.IdleDeleteTtl != "" {
			lifecycleArgs.IdleDeleteTtl = pulumi.StringPtr(cfg.LifecycleConfig.IdleDeleteTtl)
		}
		if cfg.LifecycleConfig.AutoDeleteTime != "" {
			lifecycleArgs.AutoDeleteTime = pulumi.StringPtr(cfg.LifecycleConfig.AutoDeleteTime)
		}
		clusterConfig.LifecycleConfig = lifecycleArgs
	}

	// ── Persistent shared Hive metastore ──

	if cfg.MetastoreConfig != nil {
		clusterConfig.MetastoreConfig = &dataproc.ClusterClusterConfigMetastoreConfigArgs{
			DataprocMetastoreService: pulumi.String(cfg.MetastoreConfig.DataprocMetastoreService.GetValue()),
		}
	}

	// ── OSS metric collection into Cloud Monitoring ──

	if cfg.DataprocMetricConfig != nil {
		var metrics dataproc.ClusterClusterConfigDataprocMetricConfigMetricArray
		for _, m := range cfg.DataprocMetricConfig.Metrics {
			metricArgs := &dataproc.ClusterClusterConfigDataprocMetricConfigMetricArgs{
				MetricSource: pulumi.String(m.MetricSource),
			}
			if len(m.MetricOverrides) > 0 {
				metricArgs.MetricOverrides = pulumi.ToStringArray(m.MetricOverrides)
			}
			metrics = append(metrics, metricArgs)
		}
		clusterConfig.DataprocMetricConfig = &dataproc.ClusterClusterConfigDataprocMetricConfigArgs{
			Metrics: metrics,
		}
	}

	// ── Dedicated DRIVER node groups ──

	if len(cfg.AuxiliaryNodeGroups) > 0 {
		var groups dataproc.ClusterClusterConfigAuxiliaryNodeGroupArray
		for _, g := range cfg.AuxiliaryNodeGroups {
			nodeGroupArgs := &dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupArgs{
				Roles: pulumi.ToStringArray(g.Roles),
			}

			if g.NodeGroupConfig != nil {
				ngc := g.NodeGroupConfig
				configArgs := &dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupNodeGroupConfigArgs{}
				if ngc.NumInstances > 0 {
					configArgs.NumInstances = pulumi.IntPtr(int(ngc.NumInstances))
				}
				if ngc.MachineType != "" {
					configArgs.MachineType = pulumi.StringPtr(ngc.MachineType)
				}
				if ngc.MinCpuPlatform != "" {
					configArgs.MinCpuPlatform = pulumi.StringPtr(ngc.MinCpuPlatform)
				}
				if ngc.DiskConfig != nil {
					sizeGb, diskType, ssds, ssdIface := diskFields(ngc.DiskConfig)
					diskArgs := &dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupNodeGroupConfigDiskConfigArgs{}
					if sizeGb != nil {
						diskArgs.BootDiskSizeGb = pulumi.IntPtr(*sizeGb)
					}
					if diskType != nil {
						diskArgs.BootDiskType = pulumi.StringPtr(*diskType)
					}
					if ssds != nil {
						diskArgs.NumLocalSsds = pulumi.IntPtr(*ssds)
					}
					if ssdIface != nil {
						diskArgs.LocalSsdInterface = pulumi.StringPtr(*ssdIface)
					}
					configArgs.DiskConfig = diskArgs
				}
				if len(ngc.Accelerators) > 0 {
					var accels dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupNodeGroupConfigAcceleratorArray
					for _, a := range ngc.Accelerators {
						accels = append(accels, &dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupNodeGroupConfigAcceleratorArgs{
							AcceleratorType:  pulumi.String(a.AcceleratorType),
							AcceleratorCount: pulumi.Int(int(a.AcceleratorCount)),
						})
					}
					configArgs.Accelerators = accels
				}
				nodeGroupArgs.NodeGroupConfig = configArgs
			}

			groupArgs := &dataproc.ClusterClusterConfigAuxiliaryNodeGroupArgs{
				NodeGroups: dataproc.ClusterClusterConfigAuxiliaryNodeGroupNodeGroupArray{nodeGroupArgs},
			}
			if g.NodeGroupId != "" {
				groupArgs.NodeGroupId = pulumi.StringPtr(g.NodeGroupId)
			}
			groups = append(groups, groupArgs)
		}
		clusterConfig.AuxiliaryNodeGroups = groups
	}

	return clusterConfig
}

// buildVirtualClusterConfig assembles the Dataproc-on-GKE arm. All
// references arrive resolved to the fully qualified resource names the
// Dataproc API requires. The whole arm is immutable — changes replace
// the virtual cluster without touching the underlying GKE resources.
func buildVirtualClusterConfig(vcc *gcpdataprocclusterv1alpha1.GcpDataprocClusterVirtualClusterConfig) *dataproc.ClusterVirtualClusterConfigArgs {
	args := &dataproc.ClusterVirtualClusterConfigArgs{}

	if vcc.StagingBucket.GetValue() != "" {
		args.StagingBucket = pulumi.StringPtr(vcc.StagingBucket.GetValue())
	}

	kcc := vcc.KubernetesClusterConfig
	kccArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigArgs{}

	if kcc.KubernetesNamespace.GetValue() != "" {
		kccArgs.KubernetesNamespace = pulumi.StringPtr(kcc.KubernetesNamespace.GetValue())
	}

	gkeArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigArgs{
		GkeClusterTarget: pulumi.StringPtr(kcc.GkeClusterConfig.GkeClusterTarget.GetValue()),
	}

	if len(kcc.GkeClusterConfig.NodePoolTarget) > 0 {
		var targets dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetArray
		for _, t := range kcc.GkeClusterConfig.NodePoolTarget {
			targetArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetArgs{
				NodePool: pulumi.String(t.NodePool.GetValue()),
				Roles:    pulumi.ToStringArray(t.Roles),
			}

			if t.NodePoolConfig != nil {
				npc := t.NodePoolConfig
				npcArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigArgs{
					Locations: pulumi.ToStringArray(npc.Locations),
				}

				if npc.Autoscaling != nil {
					npcArgs.Autoscaling = &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigAutoscalingArgs{
						MinNodeCount: pulumi.IntPtr(int(npc.Autoscaling.MinNodeCount)),
						MaxNodeCount: pulumi.IntPtr(int(npc.Autoscaling.MaxNodeCount)),
					}
				}

				// Sizing for a Dataproc-created pool; ignored when the
				// referenced pool already exists.
				if npc.MachineType != "" || npc.LocalSsdCount > 0 || npc.MinCpuPlatform != "" || npc.Preemptible || npc.Spot {
					configArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigGkeClusterConfigNodePoolTargetNodePoolConfigConfigArgs{}
					if npc.MachineType != "" {
						configArgs.MachineType = pulumi.StringPtr(npc.MachineType)
					}
					if npc.LocalSsdCount > 0 {
						configArgs.LocalSsdCount = pulumi.IntPtr(int(npc.LocalSsdCount))
					}
					if npc.MinCpuPlatform != "" {
						configArgs.MinCpuPlatform = pulumi.StringPtr(npc.MinCpuPlatform)
					}
					if npc.Preemptible {
						configArgs.Preemptible = pulumi.BoolPtr(true)
					}
					if npc.Spot {
						configArgs.Spot = pulumi.BoolPtr(true)
					}
					npcArgs.Config = configArgs
				}

				targetArgs.NodePoolConfig = npcArgs
			}

			targets = append(targets, targetArgs)
		}
		gkeArgs.NodePoolTargets = targets
	}

	kccArgs.GkeClusterConfig = gkeArgs

	softwareArgs := &dataproc.ClusterVirtualClusterConfigKubernetesClusterConfigKubernetesSoftwareConfigArgs{
		ComponentVersion: pulumi.ToStringMap(kcc.KubernetesSoftwareConfig.ComponentVersion),
	}
	if len(kcc.KubernetesSoftwareConfig.Properties) > 0 {
		softwareArgs.Properties = pulumi.ToStringMap(kcc.KubernetesSoftwareConfig.Properties)
	}
	kccArgs.KubernetesSoftwareConfig = softwareArgs

	args.KubernetesClusterConfig = kccArgs

	if vcc.AuxiliaryServicesConfig != nil {
		auxArgs := &dataproc.ClusterVirtualClusterConfigAuxiliaryServicesConfigArgs{}

		if vcc.AuxiliaryServicesConfig.MetastoreConfig != nil {
			auxArgs.MetastoreConfig = &dataproc.ClusterVirtualClusterConfigAuxiliaryServicesConfigMetastoreConfigArgs{
				DataprocMetastoreService: pulumi.StringPtr(vcc.AuxiliaryServicesConfig.MetastoreConfig.DataprocMetastoreService.GetValue()),
			}
		}

		if shs := vcc.AuxiliaryServicesConfig.SparkHistoryServerConfig; shs != nil && shs.DataprocCluster.GetValue() != "" {
			auxArgs.SparkHistoryServerConfig = &dataproc.ClusterVirtualClusterConfigAuxiliaryServicesConfigSparkHistoryServerConfigArgs{
				DataprocCluster: pulumi.StringPtr(shs.DataprocCluster.GetValue()),
			}
		}

		args.AuxiliaryServicesConfig = auxArgs
	}

	return args
}
