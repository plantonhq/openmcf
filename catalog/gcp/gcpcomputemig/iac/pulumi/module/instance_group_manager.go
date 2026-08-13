package module

import (
	"github.com/pkg/errors"
	gcpcomputemigv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcomputemig/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// igmResult carries what downstream resources (autoscaler, per-instance
// configs, resize requests) need from the created group manager.
type igmResult struct {
	Resource      pulumi.Resource
	Name          pulumi.StringOutput
	SelfLink      pulumi.StringOutput
	InstanceGroup pulumi.StringOutput
}

// versionTemplateRef resolves one version's template reference: an empty
// template_self_link means "this kind's own template" (the default), a
// non-empty value pins an external template URL — the canary escape
// hatch.
func versionTemplateRef(version *gcpcomputemigv1alpha1.GcpComputeMigVersion, ownTemplate pulumi.StringOutput) pulumi.StringInput {
	if version.TemplateSelfLink != "" {
		return pulumi.String(version.TemplateSelfLink)
	}
	return ownTemplate
}

// instanceGroupManager creates the group manager — zonal or regional per
// the spec's location selector — wired to the created template. When the
// spec declares no versions, the group runs one default version on the
// kind's own template.
func instanceGroupManager(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
	template *templateResult,
) (*igmResult, error) {

	spec := locals.GcpComputeMig.Spec

	if locals.IsRegional {
		versions := compute.RegionInstanceGroupManagerVersionArray{}
		if len(spec.Versions) == 0 {
			versions = append(versions, &compute.RegionInstanceGroupManagerVersionArgs{
				InstanceTemplate: template.TemplateRef,
			})
		}
		for _, version := range spec.Versions {
			versionArgs := &compute.RegionInstanceGroupManagerVersionArgs{
				InstanceTemplate: versionTemplateRef(version, template.TemplateRef),
			}
			if version.VersionName != "" {
				versionArgs.Name = pulumi.StringPtr(version.VersionName)
			}
			if version.TargetSizeFixed != nil || version.TargetSizePercent != nil {
				targetSize := &compute.RegionInstanceGroupManagerVersionTargetSizeArgs{}
				if version.TargetSizeFixed != nil {
					targetSize.Fixed = pulumi.IntPtr(int(version.GetTargetSizeFixed()))
				}
				if version.TargetSizePercent != nil {
					targetSize.Percent = pulumi.IntPtr(int(version.GetTargetSizePercent()))
				}
				versionArgs.TargetSize = targetSize
			}
			versions = append(versions, versionArgs)
		}

		args := &compute.RegionInstanceGroupManagerArgs{
			Name:             pulumi.StringPtr(locals.MigName),
			BaseInstanceName: pulumi.String(locals.BaseInstanceName),
			Region:           pulumi.String(spec.Region),
			Versions:         versions,
		}
		if spec.ProjectId.GetValue() != "" {
			args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}
		if spec.Description != "" {
			args.Description = pulumi.StringPtr(spec.Description)
		}
		if spec.TargetSize != nil {
			args.TargetSize = pulumi.IntPtr(int(spec.GetTargetSize()))
		}
		if len(spec.NamedPorts) > 0 {
			ports := compute.RegionInstanceGroupManagerNamedPortArray{}
			for _, port := range spec.NamedPorts {
				ports = append(ports, &compute.RegionInstanceGroupManagerNamedPortArgs{
					Name: pulumi.String(port.Name),
					Port: pulumi.Int(int(port.Port)),
				})
			}
			args.NamedPorts = ports
		}
		if updatePolicy := spec.UpdatePolicy; updatePolicy != nil {
			upArgs := &compute.RegionInstanceGroupManagerUpdatePolicyArgs{
				MinimalAction: pulumi.String(updatePolicy.MinimalAction),
				Type:          pulumi.String(updatePolicy.Type),
			}
			if updatePolicy.MostDisruptiveAllowedAction != "" {
				upArgs.MostDisruptiveAllowedAction = pulumi.StringPtr(updatePolicy.MostDisruptiveAllowedAction)
			}
			if updatePolicy.ReplacementMethod != "" {
				upArgs.ReplacementMethod = pulumi.StringPtr(updatePolicy.ReplacementMethod)
			}
			if updatePolicy.MaxSurgeFixed != nil {
				upArgs.MaxSurgeFixed = pulumi.IntPtr(int(updatePolicy.GetMaxSurgeFixed()))
			}
			if updatePolicy.MaxSurgePercent != nil {
				upArgs.MaxSurgePercent = pulumi.IntPtr(int(updatePolicy.GetMaxSurgePercent()))
			}
			if updatePolicy.MaxUnavailableFixed != nil {
				upArgs.MaxUnavailableFixed = pulumi.IntPtr(int(updatePolicy.GetMaxUnavailableFixed()))
			}
			if updatePolicy.MaxUnavailablePercent != nil {
				upArgs.MaxUnavailablePercent = pulumi.IntPtr(int(updatePolicy.GetMaxUnavailablePercent()))
			}
			// Regional-only: PROACTIVE (default) rebalances zones;
			// NONE is required for stateful regional groups.
			if updatePolicy.InstanceRedistributionType != "" {
				upArgs.InstanceRedistributionType = pulumi.StringPtr(updatePolicy.InstanceRedistributionType)
			}
			args.UpdatePolicy = upArgs
		}
		if spec.AutoHealing != nil {
			args.AutoHealingPolicies = &compute.RegionInstanceGroupManagerAutoHealingPoliciesArgs{
				HealthCheck:     pulumi.String(spec.AutoHealing.HealthCheck.GetValue()),
				InitialDelaySec: pulumi.Int(int(spec.AutoHealing.InitialDelaySec)),
			}
		}
		if spec.StandbyPolicy != nil {
			standby := &compute.RegionInstanceGroupManagerStandbyPolicyArgs{}
			if spec.StandbyPolicy.InitialDelaySec != nil {
				standby.InitialDelaySec = pulumi.IntPtr(int(spec.StandbyPolicy.GetInitialDelaySec()))
			}
			if spec.StandbyPolicy.Mode != "" {
				standby.Mode = pulumi.StringPtr(spec.StandbyPolicy.Mode)
			}
			args.StandbyPolicy = standby
		}
		if spec.TargetSuspendedSize != nil {
			args.TargetSuspendedSize = pulumi.IntPtr(int(spec.GetTargetSuspendedSize()))
		}
		if spec.TargetStoppedSize != nil {
			args.TargetStoppedSize = pulumi.IntPtr(int(spec.GetTargetStoppedSize()))
		}
		if len(spec.StatefulDisks) > 0 {
			disks := compute.RegionInstanceGroupManagerStatefulDiskArray{}
			for _, disk := range spec.StatefulDisks {
				diskArgs := &compute.RegionInstanceGroupManagerStatefulDiskArgs{
					DeviceName: pulumi.String(disk.DeviceName),
				}
				if disk.DeleteRule != "" {
					diskArgs.DeleteRule = pulumi.StringPtr(disk.DeleteRule)
				}
				disks = append(disks, diskArgs)
			}
			args.StatefulDisks = disks
		}
		if len(spec.StatefulExternalIps) > 0 {
			ips := compute.RegionInstanceGroupManagerStatefulExternalIpArray{}
			for _, ip := range spec.StatefulExternalIps {
				ipArgs := &compute.RegionInstanceGroupManagerStatefulExternalIpArgs{}
				if ip.InterfaceName != "" {
					ipArgs.InterfaceName = pulumi.StringPtr(ip.InterfaceName)
				}
				if ip.DeleteRule != "" {
					ipArgs.DeleteRule = pulumi.StringPtr(ip.DeleteRule)
				}
				ips = append(ips, ipArgs)
			}
			args.StatefulExternalIps = ips
		}
		if len(spec.StatefulInternalIps) > 0 {
			ips := compute.RegionInstanceGroupManagerStatefulInternalIpArray{}
			for _, ip := range spec.StatefulInternalIps {
				ipArgs := &compute.RegionInstanceGroupManagerStatefulInternalIpArgs{}
				if ip.InterfaceName != "" {
					ipArgs.InterfaceName = pulumi.StringPtr(ip.InterfaceName)
				}
				if ip.DeleteRule != "" {
					ipArgs.DeleteRule = pulumi.StringPtr(ip.DeleteRule)
				}
				ips = append(ips, ipArgs)
			}
			args.StatefulInternalIps = ips
		}
		if lifecycle := spec.InstanceLifecyclePolicy; lifecycle != nil {
			lifecycleArgs := &compute.RegionInstanceGroupManagerInstanceLifecyclePolicyArgs{}
			if lifecycle.DefaultActionOnFailure != "" {
				lifecycleArgs.DefaultActionOnFailure = pulumi.StringPtr(lifecycle.DefaultActionOnFailure)
			}
			if lifecycle.ForceUpdateOnRepair != "" {
				lifecycleArgs.ForceUpdateOnRepair = pulumi.StringPtr(lifecycle.ForceUpdateOnRepair)
			}
			if lifecycle.OnFailedHealthCheck != "" {
				lifecycleArgs.OnFailedHealthCheck = pulumi.StringPtr(lifecycle.OnFailedHealthCheck)
			}
			if lifecycle.OnRepairAllowChangingZone != "" {
				lifecycleArgs.OnRepair = &compute.RegionInstanceGroupManagerInstanceLifecyclePolicyOnRepairArgs{
					AllowChangingZone: pulumi.StringPtr(lifecycle.OnRepairAllowChangingZone),
				}
			}
			args.InstanceLifecyclePolicy = lifecycleArgs
		}
		if spec.AllInstancesConfig != nil {
			aic := &compute.RegionInstanceGroupManagerAllInstancesConfigArgs{}
			if len(spec.AllInstancesConfig.Labels) > 0 {
				aic.Labels = pulumi.ToStringMap(spec.AllInstancesConfig.Labels)
			}
			if len(spec.AllInstancesConfig.Metadata) > 0 {
				aic.Metadata = pulumi.ToStringMap(spec.AllInstancesConfig.Metadata)
			}
			args.AllInstancesConfig = aic
		}
		if spec.ListManagedInstancesResults != "" {
			args.ListManagedInstancesResults = pulumi.StringPtr(spec.ListManagedInstancesResults)
		}
		if spec.WorkloadPolicy != "" {
			args.ResourcePolicies = &compute.RegionInstanceGroupManagerResourcePoliciesArgs{
				WorkloadPolicy: pulumi.StringPtr(spec.WorkloadPolicy),
			}
		}
		if len(spec.TargetPools) > 0 {
			args.TargetPools = pulumi.ToStringArray(spec.TargetPools)
		}
		if spec.WaitForInstances != nil {
			args.WaitForInstances = pulumi.BoolPtr(spec.GetWaitForInstances())
		}
		if spec.WaitForInstancesStatus != "" {
			args.WaitForInstancesStatus = pulumi.StringPtr(spec.WaitForInstancesStatus)
		}
		if spec.DistributionPolicy != nil {
			if len(spec.DistributionPolicy.Zones) > 0 {
				args.DistributionPolicyZones = pulumi.ToStringArray(spec.DistributionPolicy.Zones)
			}
			if spec.DistributionPolicy.TargetShape != "" {
				args.DistributionPolicyTargetShape = pulumi.StringPtr(spec.DistributionPolicy.TargetShape)
			}
		}
		if spec.InstanceFlexibilityPolicy != nil {
			selections := compute.RegionInstanceGroupManagerInstanceFlexibilityPolicyInstanceSelectionArray{}
			for _, selection := range spec.InstanceFlexibilityPolicy.InstanceSelections {
				selectionArgs := &compute.RegionInstanceGroupManagerInstanceFlexibilityPolicyInstanceSelectionArgs{
					Name:         pulumi.String(selection.Name),
					MachineTypes: pulumi.ToStringArray(selection.MachineTypes),
				}
				if selection.Rank != nil {
					selectionArgs.Rank = pulumi.IntPtr(int(selection.GetRank()))
				}
				selections = append(selections, selectionArgs)
			}
			args.InstanceFlexibilityPolicy = &compute.RegionInstanceGroupManagerInstanceFlexibilityPolicyArgs{
				InstanceSelections: selections,
			}
		}
		if spec.TargetSizePolicyMode != "" {
			args.TargetSizePolicies = compute.RegionInstanceGroupManagerTargetSizePolicyArray{
				&compute.RegionInstanceGroupManagerTargetSizePolicyArgs{
					Mode: pulumi.String(spec.TargetSizePolicyMode),
				},
			}
		}
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}

		createdManager, err := compute.NewRegionInstanceGroupManager(ctx,
			locals.MigName,
			args,
			pulumi.Provider(gcpProvider),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create regional instance group manager")
		}
		return &igmResult{
			Resource:      createdManager,
			Name:          createdManager.Name,
			SelfLink:      createdManager.SelfLink,
			InstanceGroup: createdManager.InstanceGroup,
		}, nil
	}

	versions := compute.InstanceGroupManagerVersionArray{}
	if len(spec.Versions) == 0 {
		versions = append(versions, &compute.InstanceGroupManagerVersionArgs{
			InstanceTemplate: template.TemplateRef,
		})
	}
	for _, version := range spec.Versions {
		versionArgs := &compute.InstanceGroupManagerVersionArgs{
			InstanceTemplate: versionTemplateRef(version, template.TemplateRef),
		}
		if version.VersionName != "" {
			versionArgs.Name = pulumi.StringPtr(version.VersionName)
		}
		if version.TargetSizeFixed != nil || version.TargetSizePercent != nil {
			targetSize := &compute.InstanceGroupManagerVersionTargetSizeArgs{}
			if version.TargetSizeFixed != nil {
				targetSize.Fixed = pulumi.IntPtr(int(version.GetTargetSizeFixed()))
			}
			if version.TargetSizePercent != nil {
				targetSize.Percent = pulumi.IntPtr(int(version.GetTargetSizePercent()))
			}
			versionArgs.TargetSize = targetSize
		}
		versions = append(versions, versionArgs)
	}

	args := &compute.InstanceGroupManagerArgs{
		Name:             pulumi.StringPtr(locals.MigName),
		BaseInstanceName: pulumi.String(locals.BaseInstanceName),
		Zone:             pulumi.StringPtr(spec.Zone),
		Versions:         versions,
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}
	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.TargetSize != nil {
		args.TargetSize = pulumi.IntPtr(int(spec.GetTargetSize()))
	}
	if len(spec.NamedPorts) > 0 {
		ports := compute.InstanceGroupManagerNamedPortArray{}
		for _, port := range spec.NamedPorts {
			ports = append(ports, &compute.InstanceGroupManagerNamedPortArgs{
				Name: pulumi.String(port.Name),
				Port: pulumi.Int(int(port.Port)),
			})
		}
		args.NamedPorts = ports
	}
	if updatePolicy := spec.UpdatePolicy; updatePolicy != nil {
		upArgs := &compute.InstanceGroupManagerUpdatePolicyArgs{
			MinimalAction: pulumi.String(updatePolicy.MinimalAction),
			Type:          pulumi.String(updatePolicy.Type),
		}
		if updatePolicy.MostDisruptiveAllowedAction != "" {
			upArgs.MostDisruptiveAllowedAction = pulumi.StringPtr(updatePolicy.MostDisruptiveAllowedAction)
		}
		if updatePolicy.ReplacementMethod != "" {
			upArgs.ReplacementMethod = pulumi.StringPtr(updatePolicy.ReplacementMethod)
		}
		if updatePolicy.MaxSurgeFixed != nil {
			upArgs.MaxSurgeFixed = pulumi.IntPtr(int(updatePolicy.GetMaxSurgeFixed()))
		}
		if updatePolicy.MaxSurgePercent != nil {
			upArgs.MaxSurgePercent = pulumi.IntPtr(int(updatePolicy.GetMaxSurgePercent()))
		}
		if updatePolicy.MaxUnavailableFixed != nil {
			upArgs.MaxUnavailableFixed = pulumi.IntPtr(int(updatePolicy.GetMaxUnavailableFixed()))
		}
		if updatePolicy.MaxUnavailablePercent != nil {
			upArgs.MaxUnavailablePercent = pulumi.IntPtr(int(updatePolicy.GetMaxUnavailablePercent()))
		}
		args.UpdatePolicy = upArgs
	}
	if spec.AutoHealing != nil {
		args.AutoHealingPolicies = &compute.InstanceGroupManagerAutoHealingPoliciesArgs{
			HealthCheck:     pulumi.String(spec.AutoHealing.HealthCheck.GetValue()),
			InitialDelaySec: pulumi.Int(int(spec.AutoHealing.InitialDelaySec)),
		}
	}
	if spec.StandbyPolicy != nil {
		standby := &compute.InstanceGroupManagerStandbyPolicyArgs{}
		if spec.StandbyPolicy.InitialDelaySec != nil {
			standby.InitialDelaySec = pulumi.IntPtr(int(spec.StandbyPolicy.GetInitialDelaySec()))
		}
		if spec.StandbyPolicy.Mode != "" {
			standby.Mode = pulumi.StringPtr(spec.StandbyPolicy.Mode)
		}
		args.StandbyPolicy = standby
	}
	if spec.TargetSuspendedSize != nil {
		args.TargetSuspendedSize = pulumi.IntPtr(int(spec.GetTargetSuspendedSize()))
	}
	if spec.TargetStoppedSize != nil {
		args.TargetStoppedSize = pulumi.IntPtr(int(spec.GetTargetStoppedSize()))
	}
	if len(spec.StatefulDisks) > 0 {
		disks := compute.InstanceGroupManagerStatefulDiskArray{}
		for _, disk := range spec.StatefulDisks {
			diskArgs := &compute.InstanceGroupManagerStatefulDiskArgs{
				DeviceName: pulumi.String(disk.DeviceName),
			}
			if disk.DeleteRule != "" {
				diskArgs.DeleteRule = pulumi.StringPtr(disk.DeleteRule)
			}
			disks = append(disks, diskArgs)
		}
		args.StatefulDisks = disks
	}
	if len(spec.StatefulExternalIps) > 0 {
		ips := compute.InstanceGroupManagerStatefulExternalIpArray{}
		for _, ip := range spec.StatefulExternalIps {
			ipArgs := &compute.InstanceGroupManagerStatefulExternalIpArgs{}
			if ip.InterfaceName != "" {
				ipArgs.InterfaceName = pulumi.StringPtr(ip.InterfaceName)
			}
			if ip.DeleteRule != "" {
				ipArgs.DeleteRule = pulumi.StringPtr(ip.DeleteRule)
			}
			ips = append(ips, ipArgs)
		}
		args.StatefulExternalIps = ips
	}
	if len(spec.StatefulInternalIps) > 0 {
		ips := compute.InstanceGroupManagerStatefulInternalIpArray{}
		for _, ip := range spec.StatefulInternalIps {
			ipArgs := &compute.InstanceGroupManagerStatefulInternalIpArgs{}
			if ip.InterfaceName != "" {
				ipArgs.InterfaceName = pulumi.StringPtr(ip.InterfaceName)
			}
			if ip.DeleteRule != "" {
				ipArgs.DeleteRule = pulumi.StringPtr(ip.DeleteRule)
			}
			ips = append(ips, ipArgs)
		}
		args.StatefulInternalIps = ips
	}
	if lifecycle := spec.InstanceLifecyclePolicy; lifecycle != nil {
		lifecycleArgs := &compute.InstanceGroupManagerInstanceLifecyclePolicyArgs{}
		if lifecycle.DefaultActionOnFailure != "" {
			lifecycleArgs.DefaultActionOnFailure = pulumi.StringPtr(lifecycle.DefaultActionOnFailure)
		}
		if lifecycle.ForceUpdateOnRepair != "" {
			lifecycleArgs.ForceUpdateOnRepair = pulumi.StringPtr(lifecycle.ForceUpdateOnRepair)
		}
		if lifecycle.OnFailedHealthCheck != "" {
			lifecycleArgs.OnFailedHealthCheck = pulumi.StringPtr(lifecycle.OnFailedHealthCheck)
		}
		if lifecycle.OnRepairAllowChangingZone != "" {
			lifecycleArgs.OnRepair = &compute.InstanceGroupManagerInstanceLifecyclePolicyOnRepairArgs{
				AllowChangingZone: pulumi.StringPtr(lifecycle.OnRepairAllowChangingZone),
			}
		}
		args.InstanceLifecyclePolicy = lifecycleArgs
	}
	if spec.AllInstancesConfig != nil {
		aic := &compute.InstanceGroupManagerAllInstancesConfigArgs{}
		if len(spec.AllInstancesConfig.Labels) > 0 {
			aic.Labels = pulumi.ToStringMap(spec.AllInstancesConfig.Labels)
		}
		if len(spec.AllInstancesConfig.Metadata) > 0 {
			aic.Metadata = pulumi.ToStringMap(spec.AllInstancesConfig.Metadata)
		}
		args.AllInstancesConfig = aic
	}
	if spec.ListManagedInstancesResults != "" {
		args.ListManagedInstancesResults = pulumi.StringPtr(spec.ListManagedInstancesResults)
	}
	if spec.WorkloadPolicy != "" {
		args.ResourcePolicies = &compute.InstanceGroupManagerResourcePoliciesArgs{
			WorkloadPolicy: pulumi.StringPtr(spec.WorkloadPolicy),
		}
	}
	if len(spec.TargetPools) > 0 {
		args.TargetPools = pulumi.ToStringArray(spec.TargetPools)
	}
	if spec.WaitForInstances != nil {
		args.WaitForInstances = pulumi.BoolPtr(spec.GetWaitForInstances())
	}
	if spec.WaitForInstancesStatus != "" {
		args.WaitForInstancesStatus = pulumi.StringPtr(spec.WaitForInstancesStatus)
	}
	if spec.TargetSizePolicyMode != "" {
		args.TargetSizePolicies = compute.InstanceGroupManagerTargetSizePolicyArray{
			&compute.InstanceGroupManagerTargetSizePolicyArgs{
				Mode: pulumi.String(spec.TargetSizePolicyMode),
			},
		}
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
	}

	createdManager, err := compute.NewInstanceGroupManager(ctx,
		locals.MigName,
		args,
		pulumi.Provider(gcpProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create instance group manager")
	}
	return &igmResult{
		Resource:      createdManager,
		Name:          createdManager.Name,
		SelfLink:      createdManager.SelfLink,
		InstanceGroup: createdManager.InstanceGroup,
	}, nil
}
