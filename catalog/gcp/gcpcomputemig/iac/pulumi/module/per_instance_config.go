package module

import (
	"github.com/pkg/errors"
	gcpcomputemigv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcomputemig/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// perInstanceConfigs creates one stateful per-instance config per spec
// entry — zonal or regional per the location selector — each wired to
// the group manager by NAME (the provider's expected reference form for
// this resource pair).
func perInstanceConfigs(
	ctx *pulumi.Context,
	locals *Locals,
	gcpProvider *gcp.Provider,
	groupManager *igmResult,
) error {

	spec := locals.GcpComputeMig.Spec

	for _, config := range spec.PerInstanceConfigs {
		if locals.IsRegional {
			args := &compute.RegionPerInstanceConfigArgs{
				Name:                       pulumi.String(config.ConfigName),
				Region:                     pulumi.StringPtr(spec.Region),
				RegionInstanceGroupManager: groupManager.Name,
			}
			if spec.ProjectId.GetValue() != "" {
				args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
			}
			if config.PreservedState != nil {
				args.PreservedState = regionalPreservedState(config.PreservedState)
			}
			if config.MinimalAction != "" {
				args.MinimalAction = pulumi.StringPtr(config.MinimalAction)
			}
			if config.MostDisruptiveAllowedAction != "" {
				args.MostDisruptiveAllowedAction = pulumi.StringPtr(config.MostDisruptiveAllowedAction)
			}
			if config.RemoveInstanceOnDestroy != nil {
				args.RemoveInstanceOnDestroy = pulumi.BoolPtr(config.GetRemoveInstanceOnDestroy())
			}
			if config.RemoveInstanceStateOnDestroy != nil {
				args.RemoveInstanceStateOnDestroy = pulumi.BoolPtr(config.GetRemoveInstanceStateOnDestroy())
			}
			if spec.DeletionPolicy != "" {
				args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
			}
			if _, err := compute.NewRegionPerInstanceConfig(ctx,
				config.ConfigName,
				args,
				pulumi.Provider(gcpProvider),
			); err != nil {
				return errors.Wrapf(err, "failed to create regional per-instance config %s", config.ConfigName)
			}
			continue
		}

		args := &compute.PerInstanceConfigArgs{
			Name:                 pulumi.String(config.ConfigName),
			Zone:                 pulumi.StringPtr(spec.Zone),
			InstanceGroupManager: groupManager.Name,
		}
		if spec.ProjectId.GetValue() != "" {
			args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}
		if config.PreservedState != nil {
			args.PreservedState = zonalPreservedState(config.PreservedState)
		}
		if config.MinimalAction != "" {
			args.MinimalAction = pulumi.StringPtr(config.MinimalAction)
		}
		if config.MostDisruptiveAllowedAction != "" {
			args.MostDisruptiveAllowedAction = pulumi.StringPtr(config.MostDisruptiveAllowedAction)
		}
		if config.RemoveInstanceOnDestroy != nil {
			args.RemoveInstanceOnDestroy = pulumi.BoolPtr(config.GetRemoveInstanceOnDestroy())
		}
		if config.RemoveInstanceStateOnDestroy != nil {
			args.RemoveInstanceStateOnDestroy = pulumi.BoolPtr(config.GetRemoveInstanceStateOnDestroy())
		}
		if spec.DeletionPolicy != "" {
			args.DeletionPolicy = pulumi.StringPtr(spec.DeletionPolicy)
		}
		if _, err := compute.NewPerInstanceConfig(ctx,
			config.ConfigName,
			args,
			pulumi.Provider(gcpProvider),
		); err != nil {
			return errors.Wrapf(err, "failed to create per-instance config %s", config.ConfigName)
		}
	}

	return nil
}

// zonalPreservedState maps the spec's preserved state to the zonal
// config's types.
func zonalPreservedState(preserved *gcpcomputemigv1alpha1.GcpComputeMigPreservedState) *compute.PerInstanceConfigPreservedStateArgs {
	args := &compute.PerInstanceConfigPreservedStateArgs{}
	if len(preserved.Metadata) > 0 {
		args.Metadata = pulumi.ToStringMap(preserved.Metadata)
	}
	if len(preserved.Disks) > 0 {
		disks := compute.PerInstanceConfigPreservedStateDiskArray{}
		for _, disk := range preserved.Disks {
			diskArgs := &compute.PerInstanceConfigPreservedStateDiskArgs{
				DeviceName: pulumi.String(disk.DeviceName),
				Source:     pulumi.String(disk.Source.GetValue()),
			}
			if disk.Mode != "" {
				diskArgs.Mode = pulumi.StringPtr(disk.Mode)
			}
			if disk.DeleteRule != "" {
				diskArgs.DeleteRule = pulumi.StringPtr(disk.DeleteRule)
			}
			disks = append(disks, diskArgs)
		}
		args.Disks = disks
	}
	if len(preserved.ExternalIps) > 0 {
		ips := compute.PerInstanceConfigPreservedStateExternalIpArray{}
		for _, ip := range preserved.ExternalIps {
			ipArgs := &compute.PerInstanceConfigPreservedStateExternalIpArgs{
				InterfaceName: pulumi.String(ip.InterfaceName),
			}
			if ip.Address.GetValue() != "" {
				ipArgs.IpAddress = &compute.PerInstanceConfigPreservedStateExternalIpIpAddressArgs{
					Address: pulumi.StringPtr(ip.Address.GetValue()),
				}
			}
			if ip.AutoDelete != "" {
				ipArgs.AutoDelete = pulumi.StringPtr(ip.AutoDelete)
			}
			ips = append(ips, ipArgs)
		}
		args.ExternalIps = ips
	}
	if len(preserved.InternalIps) > 0 {
		ips := compute.PerInstanceConfigPreservedStateInternalIpArray{}
		for _, ip := range preserved.InternalIps {
			ipArgs := &compute.PerInstanceConfigPreservedStateInternalIpArgs{
				InterfaceName: pulumi.String(ip.InterfaceName),
			}
			if ip.Address.GetValue() != "" {
				ipArgs.IpAddress = &compute.PerInstanceConfigPreservedStateInternalIpIpAddressArgs{
					Address: pulumi.StringPtr(ip.Address.GetValue()),
				}
			}
			if ip.AutoDelete != "" {
				ipArgs.AutoDelete = pulumi.StringPtr(ip.AutoDelete)
			}
			ips = append(ips, ipArgs)
		}
		args.InternalIps = ips
	}
	return args
}

// regionalPreservedState mirrors the zonal builder for the regional
// config's types.
func regionalPreservedState(preserved *gcpcomputemigv1alpha1.GcpComputeMigPreservedState) *compute.RegionPerInstanceConfigPreservedStateArgs {
	args := &compute.RegionPerInstanceConfigPreservedStateArgs{}
	if len(preserved.Metadata) > 0 {
		args.Metadata = pulumi.ToStringMap(preserved.Metadata)
	}
	if len(preserved.Disks) > 0 {
		disks := compute.RegionPerInstanceConfigPreservedStateDiskArray{}
		for _, disk := range preserved.Disks {
			diskArgs := &compute.RegionPerInstanceConfigPreservedStateDiskArgs{
				DeviceName: pulumi.String(disk.DeviceName),
				Source:     pulumi.String(disk.Source.GetValue()),
			}
			if disk.Mode != "" {
				diskArgs.Mode = pulumi.StringPtr(disk.Mode)
			}
			if disk.DeleteRule != "" {
				diskArgs.DeleteRule = pulumi.StringPtr(disk.DeleteRule)
			}
			disks = append(disks, diskArgs)
		}
		args.Disks = disks
	}
	if len(preserved.ExternalIps) > 0 {
		ips := compute.RegionPerInstanceConfigPreservedStateExternalIpArray{}
		for _, ip := range preserved.ExternalIps {
			ipArgs := &compute.RegionPerInstanceConfigPreservedStateExternalIpArgs{
				InterfaceName: pulumi.String(ip.InterfaceName),
			}
			if ip.Address.GetValue() != "" {
				ipArgs.IpAddress = &compute.RegionPerInstanceConfigPreservedStateExternalIpIpAddressArgs{
					Address: pulumi.StringPtr(ip.Address.GetValue()),
				}
			}
			if ip.AutoDelete != "" {
				ipArgs.AutoDelete = pulumi.StringPtr(ip.AutoDelete)
			}
			ips = append(ips, ipArgs)
		}
		args.ExternalIps = ips
	}
	if len(preserved.InternalIps) > 0 {
		ips := compute.RegionPerInstanceConfigPreservedStateInternalIpArray{}
		for _, ip := range preserved.InternalIps {
			ipArgs := &compute.RegionPerInstanceConfigPreservedStateInternalIpArgs{
				InterfaceName: pulumi.String(ip.InterfaceName),
			}
			if ip.Address.GetValue() != "" {
				ipArgs.IpAddress = &compute.RegionPerInstanceConfigPreservedStateInternalIpIpAddressArgs{
					Address: pulumi.StringPtr(ip.Address.GetValue()),
				}
			}
			if ip.AutoDelete != "" {
				ipArgs.AutoDelete = pulumi.StringPtr(ip.AutoDelete)
			}
			ips = append(ips, ipArgs)
		}
		args.InternalIps = ips
	}
	return args
}
