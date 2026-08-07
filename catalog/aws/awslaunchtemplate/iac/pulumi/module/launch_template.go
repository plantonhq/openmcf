package module

import (
	"encoding/base64"
	"fmt"

	"github.com/pkg/errors"
	awslaunchtemplatev1alpha1 "github.com/plantonhq/planton/catalog/aws/awslaunchtemplate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/datatypes/stringmaps/convertstringmaps"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// launchTemplate provisions the launch template. Only the template NAME is
// create-only in AWS: every other change creates a new immutable template
// VERSION, and update_default_version promotes it to the default -- so
// consumers following "$Default" (the common ASG and node-group setup) pick
// up the change on their next launch or instance refresh, while consumers
// pinned to a numeric version keep exactly what they tested.
func launchTemplate(ctx *pulumi.Context, locals *Locals, provider pulumi.ProviderResource) error {
	spec := locals.AwsLaunchTemplate.Spec

	// AWS limits launch template names to 125 characters; truncate
	// deterministically so the same manifest always yields the same name.
	launchTemplateName := truncateName(locals.AwsLaunchTemplate.Metadata.Name, 125)

	// Identity tags land in three places on purpose: on the template itself
	// (Tags), and via TagSpecifications on the instances and volumes each
	// launch creates -- a launch template's tags do NOT propagate to what it
	// launches, so untagged fleet instances would otherwise escape
	// cost-allocation and orphan-cleanup queries.
	launchTags := convertstringmaps.ConvertGoStringMapToPulumiStringMap(
		stringmaps.AddEntry(locals.AwsTags, "Name", launchTemplateName))

	args := &ec2.LaunchTemplateArgs{
		Name: pulumi.String(launchTemplateName),
		// Promote every new version to the template default. This is the
		// declarative-model contract: the spec describes ONE desired
		// configuration, so the newest version is always the intended one.
		UpdateDefaultVersion: pulumi.BoolPtr(true),
		Tags:                 launchTags,
		TagSpecifications: ec2.LaunchTemplateTagSpecificationArray{
			&ec2.LaunchTemplateTagSpecificationArgs{
				ResourceType: pulumi.StringPtr("instance"),
				Tags:         launchTags,
			},
			&ec2.LaunchTemplateTagSpecificationArgs{
				ResourceType: pulumi.StringPtr("volume"),
				Tags:         launchTags,
			},
		},
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.ImageId != "" {
		args.ImageId = pulumi.StringPtr(spec.ImageId)
	}
	if spec.InstanceType != "" {
		args.InstanceType = pulumi.StringPtr(spec.InstanceType)
	}
	if spec.InstanceRequirements != nil {
		args.InstanceRequirements = instanceRequirementsArgs(spec.InstanceRequirements)
	}
	if spec.KeyName != "" {
		args.KeyName = pulumi.StringPtr(spec.KeyName)
	}
	// The spec carries plain text so manifests stay readable; the EC2 API
	// requires base64 -- encode here, identically to the Terraform module.
	if spec.UserData != "" {
		args.UserData = pulumi.StringPtr(base64.StdEncoding.EncodeToString([]byte(spec.UserData)))
	}
	if spec.InstanceProfile.GetValue() != "" {
		args.IamInstanceProfile = &ec2.LaunchTemplateIamInstanceProfileArgs{
			Arn: pulumi.StringPtr(spec.InstanceProfile.GetValue()),
		}
	}
	if len(spec.SecurityGroupIds) > 0 {
		securityGroupIds := make(pulumi.StringArray, 0, len(spec.SecurityGroupIds))
		for _, securityGroupId := range spec.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
		args.VpcSecurityGroupIds = securityGroupIds
	}
	// ebs_optimized is a nullable tri-state at AWS (support varies by
	// instance type), so the provider takes a string; only an explicit true
	// is sent and unset keeps the type's own default.
	if spec.EbsOptimized {
		args.EbsOptimized = pulumi.StringPtr("true")
	}
	if len(spec.BlockDeviceMappings) > 0 {
		args.BlockDeviceMappings = blockDeviceMappingArgs(spec.BlockDeviceMappings)
	}
	if len(spec.NetworkInterfaces) > 0 {
		args.NetworkInterfaces = networkInterfaceArgs(spec.NetworkInterfaces)
	}
	if spec.MetadataOptions != nil {
		args.MetadataOptions = metadataOptionsArgs(spec.MetadataOptions)
	}
	if spec.DetailedMonitoring {
		args.Monitoring = &ec2.LaunchTemplateMonitoringArgs{Enabled: pulumi.BoolPtr(true)}
	}
	if spec.Placement != nil {
		args.Placement = placementArgs(spec.Placement)
	}
	if spec.CpuOptions != nil {
		args.CpuOptions = cpuOptionsArgs(spec.CpuOptions)
	}
	if spec.CpuCredits != "" {
		args.CreditSpecification = &ec2.LaunchTemplateCreditSpecificationArgs{
			CpuCredits: pulumi.StringPtr(spec.CpuCredits),
		}
	}
	// Configuring spot_options makes every launch a Spot request; the
	// market type is implied by the block's presence.
	if spec.SpotOptions != nil {
		args.InstanceMarketOptions = &ec2.LaunchTemplateInstanceMarketOptionsArgs{
			MarketType:  pulumi.StringPtr("spot"),
			SpotOptions: spotOptionsArgs(spec.SpotOptions),
		}
	}
	if spec.EnclaveEnabled {
		args.EnclaveOptions = &ec2.LaunchTemplateEnclaveOptionsArgs{Enabled: pulumi.BoolPtr(true)}
	}
	if spec.HibernationEnabled {
		args.HibernationOptions = &ec2.LaunchTemplateHibernationOptionsArgs{Configured: pulumi.Bool(true)}
	}
	if spec.AutoRecovery != "" {
		args.MaintenanceOptions = &ec2.LaunchTemplateMaintenanceOptionsArgs{
			AutoRecovery: pulumi.StringPtr(spec.AutoRecovery),
		}
	}
	if spec.PrivateDnsNameOptions != nil {
		args.PrivateDnsNameOptions = privateDnsNameOptionsArgs(spec.PrivateDnsNameOptions)
	}
	if spec.DisableApiStop {
		args.DisableApiStop = pulumi.BoolPtr(true)
	}
	if spec.DisableApiTermination {
		args.DisableApiTermination = pulumi.BoolPtr(true)
	}
	if spec.InstanceInitiatedShutdownBehavior != "" {
		args.InstanceInitiatedShutdownBehavior = pulumi.StringPtr(spec.InstanceInitiatedShutdownBehavior)
	}

	createdLaunchTemplate, err := ec2.NewLaunchTemplate(ctx, launchTemplateName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create launch template")
	}

	ctx.Export(OpLaunchTemplateId, createdLaunchTemplate.ID())
	ctx.Export(OpLaunchTemplateArn, createdLaunchTemplate.Arn)
	ctx.Export(OpLatestVersion, createdLaunchTemplate.LatestVersion)
	ctx.Export(OpDefaultVersion, createdLaunchTemplate.DefaultVersion)

	return nil
}

// blockDeviceMappingArgs maps the spec's block device mappings. Unset EBS
// fields are omitted so the AMI's own mapping (size, type) and the account
// default (encryption) keep deciding -- what makes a minimal root-volume
// override safe.
func blockDeviceMappingArgs(mappings []*awslaunchtemplatev1alpha1.AwsLaunchTemplateBlockDeviceMapping) ec2.LaunchTemplateBlockDeviceMappingArray {
	result := make(ec2.LaunchTemplateBlockDeviceMappingArray, 0, len(mappings))
	for _, mapping := range mappings {
		mappingArgs := &ec2.LaunchTemplateBlockDeviceMappingArgs{
			DeviceName: pulumi.StringPtr(mapping.DeviceName),
		}
		if mapping.VirtualName != "" {
			mappingArgs.VirtualName = pulumi.StringPtr(mapping.VirtualName)
		}
		// no_device is a presence-signal at AWS: an empty string suppresses
		// the AMI's device. The proto's bool maps onto that convention.
		if mapping.NoDevice {
			mappingArgs.NoDevice = pulumi.StringPtr("")
		}
		if mapping.Ebs != nil {
			ebsArgs := &ec2.LaunchTemplateBlockDeviceMappingEbsArgs{}
			if mapping.Ebs.VolumeSizeGb > 0 {
				ebsArgs.VolumeSize = pulumi.IntPtr(int(mapping.Ebs.VolumeSizeGb))
			}
			if mapping.Ebs.VolumeType != "" {
				ebsArgs.VolumeType = pulumi.StringPtr(mapping.Ebs.VolumeType)
			}
			if mapping.Ebs.Iops > 0 {
				ebsArgs.Iops = pulumi.IntPtr(int(mapping.Ebs.Iops))
			}
			if mapping.Ebs.ThroughputMibps > 0 {
				ebsArgs.Throughput = pulumi.IntPtr(int(mapping.Ebs.ThroughputMibps))
			}
			if mapping.Ebs.Encrypted {
				ebsArgs.Encrypted = pulumi.StringPtr("true")
			}
			if mapping.Ebs.KmsKeyId.GetValue() != "" {
				ebsArgs.KmsKeyId = pulumi.StringPtr(mapping.Ebs.KmsKeyId.GetValue())
			}
			if mapping.Ebs.SnapshotId != "" {
				ebsArgs.SnapshotId = pulumi.StringPtr(mapping.Ebs.SnapshotId)
			}
			// delete_on_termination is a nullable tri-state at AWS (the AMI
			// mapping decides the default), so the provider takes a string:
			// nil keeps the AMI default, an explicit value overrides it.
			if mapping.Ebs.DeleteOnTermination != nil {
				ebsArgs.DeleteOnTermination = pulumi.StringPtr(fmt.Sprintf("%t", mapping.Ebs.GetDeleteOnTermination()))
			}
			mappingArgs.Ebs = ebsArgs
		}
		result = append(result, mappingArgs)
	}
	return result
}

// networkInterfaceArgs maps the spec's explicit network interfaces.
func networkInterfaceArgs(interfaces []*awslaunchtemplatev1alpha1.AwsLaunchTemplateNetworkInterface) ec2.LaunchTemplateNetworkInterfaceArray {
	result := make(ec2.LaunchTemplateNetworkInterfaceArray, 0, len(interfaces))
	for _, networkInterface := range interfaces {
		interfaceArgs := &ec2.LaunchTemplateNetworkInterfaceArgs{
			DeviceIndex: pulumi.IntPtr(int(networkInterface.DeviceIndex)),
		}
		if networkInterface.NetworkCardIndex > 0 {
			interfaceArgs.NetworkCardIndex = pulumi.IntPtr(int(networkInterface.NetworkCardIndex))
		}
		if networkInterface.Description != "" {
			interfaceArgs.Description = pulumi.StringPtr(networkInterface.Description)
		}
		if networkInterface.InterfaceType != "" {
			interfaceArgs.InterfaceType = pulumi.StringPtr(networkInterface.InterfaceType)
		}
		if networkInterface.NetworkInterfaceId != "" {
			interfaceArgs.NetworkInterfaceId = pulumi.StringPtr(networkInterface.NetworkInterfaceId)
		}
		// associate_public_ip_address and delete_on_termination are nullable
		// tri-states at AWS (the subnet / interface origin decides the
		// default), so the provider takes strings: nil inherits, an explicit
		// value overrides.
		if networkInterface.AssociatePublicIpAddress != nil {
			interfaceArgs.AssociatePublicIpAddress = pulumi.StringPtr(fmt.Sprintf("%t", networkInterface.GetAssociatePublicIpAddress()))
		}
		if networkInterface.DeleteOnTermination != nil {
			interfaceArgs.DeleteOnTermination = pulumi.StringPtr(fmt.Sprintf("%t", networkInterface.GetDeleteOnTermination()))
		}
		if networkInterface.SubnetId.GetValue() != "" {
			interfaceArgs.SubnetId = pulumi.StringPtr(networkInterface.SubnetId.GetValue())
		}
		if len(networkInterface.SecurityGroupIds) > 0 {
			securityGroups := make(pulumi.StringArray, 0, len(networkInterface.SecurityGroupIds))
			for _, securityGroupId := range networkInterface.SecurityGroupIds {
				securityGroups = append(securityGroups, pulumi.String(securityGroupId.GetValue()))
			}
			interfaceArgs.SecurityGroups = securityGroups
		}
		if networkInterface.PrivateIpAddress != "" {
			interfaceArgs.PrivateIpAddress = pulumi.StringPtr(networkInterface.PrivateIpAddress)
		}
		if networkInterface.Ipv4AddressCount > 0 {
			interfaceArgs.Ipv4AddressCount = pulumi.IntPtr(int(networkInterface.Ipv4AddressCount))
		}
		if len(networkInterface.Ipv4Addresses) > 0 {
			interfaceArgs.Ipv4Addresses = pulumi.ToStringArray(networkInterface.Ipv4Addresses)
		}
		if networkInterface.Ipv6AddressCount > 0 {
			interfaceArgs.Ipv6AddressCount = pulumi.IntPtr(int(networkInterface.Ipv6AddressCount))
		}
		if len(networkInterface.Ipv6Addresses) > 0 {
			interfaceArgs.Ipv6Addresses = pulumi.ToStringArray(networkInterface.Ipv6Addresses)
		}
		if networkInterface.Ipv4PrefixCount > 0 {
			interfaceArgs.Ipv4PrefixCount = pulumi.IntPtr(int(networkInterface.Ipv4PrefixCount))
		}
		if len(networkInterface.Ipv4Prefixes) > 0 {
			interfaceArgs.Ipv4Prefixes = pulumi.ToStringArray(networkInterface.Ipv4Prefixes)
		}
		if networkInterface.Ipv6PrefixCount > 0 {
			interfaceArgs.Ipv6PrefixCount = pulumi.IntPtr(int(networkInterface.Ipv6PrefixCount))
		}
		if len(networkInterface.Ipv6Prefixes) > 0 {
			interfaceArgs.Ipv6Prefixes = pulumi.ToStringArray(networkInterface.Ipv6Prefixes)
		}
		result = append(result, interfaceArgs)
	}
	return result
}

// instanceRequirementsArgs maps attribute-based instance selection. Only
// set fields are sent so AWS's own defaults (e.g. bare metal excluded)
// keep applying.
func instanceRequirementsArgs(requirements *awslaunchtemplatev1alpha1.AwsLaunchTemplateInstanceRequirements) *ec2.LaunchTemplateInstanceRequirementsArgs {
	// memory_mib and vcpu_count are the two AWS-required dimensions; the
	// spec enforces their presence (and that min is set).
	memoryArgs := &ec2.LaunchTemplateInstanceRequirementsMemoryMibArgs{
		Min: pulumi.Int(int(requirements.MemoryMib.Min)),
	}
	if requirements.MemoryMib.Max > 0 {
		memoryArgs.Max = pulumi.IntPtr(int(requirements.MemoryMib.Max))
	}
	vcpuArgs := &ec2.LaunchTemplateInstanceRequirementsVcpuCountArgs{
		Min: pulumi.Int(int(requirements.VcpuCount.Min)),
	}
	if requirements.VcpuCount.Max > 0 {
		vcpuArgs.Max = pulumi.IntPtr(int(requirements.VcpuCount.Max))
	}
	args := &ec2.LaunchTemplateInstanceRequirementsArgs{
		MemoryMib: memoryArgs,
		VcpuCount: vcpuArgs,
	}

	if len(requirements.AllowedInstanceTypes) > 0 {
		args.AllowedInstanceTypes = pulumi.ToStringArray(requirements.AllowedInstanceTypes)
	}
	if len(requirements.ExcludedInstanceTypes) > 0 {
		args.ExcludedInstanceTypes = pulumi.ToStringArray(requirements.ExcludedInstanceTypes)
	}
	if len(requirements.InstanceGenerations) > 0 {
		args.InstanceGenerations = pulumi.ToStringArray(requirements.InstanceGenerations)
	}
	if len(requirements.CpuManufacturers) > 0 {
		args.CpuManufacturers = pulumi.ToStringArray(requirements.CpuManufacturers)
	}
	if requirements.BareMetal != "" {
		args.BareMetal = pulumi.StringPtr(requirements.BareMetal)
	}
	if requirements.BurstablePerformance != "" {
		args.BurstablePerformance = pulumi.StringPtr(requirements.BurstablePerformance)
	}
	if requirements.RequireHibernateSupport {
		args.RequireHibernateSupport = pulumi.BoolPtr(true)
	}
	if requirements.SpotMaxPricePercentageOverLowestPrice > 0 {
		args.SpotMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.SpotMaxPricePercentageOverLowestPrice))
	}
	if requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice > 0 {
		args.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice = pulumi.IntPtr(int(requirements.MaxSpotPriceAsPercentageOfOptimalOnDemandPrice))
	}
	if requirements.OnDemandMaxPricePercentageOverLowestPrice > 0 {
		args.OnDemandMaxPricePercentageOverLowestPrice = pulumi.IntPtr(int(requirements.OnDemandMaxPricePercentageOverLowestPrice))
	}
	if requirements.LocalStorage != "" {
		args.LocalStorage = pulumi.StringPtr(requirements.LocalStorage)
	}
	if len(requirements.LocalStorageTypes) > 0 {
		args.LocalStorageTypes = pulumi.ToStringArray(requirements.LocalStorageTypes)
	}
	if requirements.TotalLocalStorageGb != nil {
		storageArgs := &ec2.LaunchTemplateInstanceRequirementsTotalLocalStorageGbArgs{}
		if requirements.TotalLocalStorageGb.Min > 0 {
			storageArgs.Min = pulumi.Float64Ptr(requirements.TotalLocalStorageGb.Min)
		}
		if requirements.TotalLocalStorageGb.Max > 0 {
			storageArgs.Max = pulumi.Float64Ptr(requirements.TotalLocalStorageGb.Max)
		}
		args.TotalLocalStorageGb = storageArgs
	}
	if requirements.MemoryGibPerVcpu != nil {
		ratioArgs := &ec2.LaunchTemplateInstanceRequirementsMemoryGibPerVcpuArgs{}
		if requirements.MemoryGibPerVcpu.Min > 0 {
			ratioArgs.Min = pulumi.Float64Ptr(requirements.MemoryGibPerVcpu.Min)
		}
		if requirements.MemoryGibPerVcpu.Max > 0 {
			ratioArgs.Max = pulumi.Float64Ptr(requirements.MemoryGibPerVcpu.Max)
		}
		args.MemoryGibPerVcpu = ratioArgs
	}
	if requirements.NetworkInterfaceCount != nil {
		countArgs := &ec2.LaunchTemplateInstanceRequirementsNetworkInterfaceCountArgs{}
		if requirements.NetworkInterfaceCount.Min > 0 {
			countArgs.Min = pulumi.IntPtr(int(requirements.NetworkInterfaceCount.Min))
		}
		if requirements.NetworkInterfaceCount.Max > 0 {
			countArgs.Max = pulumi.IntPtr(int(requirements.NetworkInterfaceCount.Max))
		}
		args.NetworkInterfaceCount = countArgs
	}
	if requirements.NetworkBandwidthGbps != nil {
		bandwidthArgs := &ec2.LaunchTemplateInstanceRequirementsNetworkBandwidthGbpsArgs{}
		if requirements.NetworkBandwidthGbps.Min > 0 {
			bandwidthArgs.Min = pulumi.Float64Ptr(requirements.NetworkBandwidthGbps.Min)
		}
		if requirements.NetworkBandwidthGbps.Max > 0 {
			bandwidthArgs.Max = pulumi.Float64Ptr(requirements.NetworkBandwidthGbps.Max)
		}
		args.NetworkBandwidthGbps = bandwidthArgs
	}
	if requirements.BaselineEbsBandwidthMbps != nil {
		ebsBandwidthArgs := &ec2.LaunchTemplateInstanceRequirementsBaselineEbsBandwidthMbpsArgs{}
		if requirements.BaselineEbsBandwidthMbps.Min > 0 {
			ebsBandwidthArgs.Min = pulumi.IntPtr(int(requirements.BaselineEbsBandwidthMbps.Min))
		}
		if requirements.BaselineEbsBandwidthMbps.Max > 0 {
			ebsBandwidthArgs.Max = pulumi.IntPtr(int(requirements.BaselineEbsBandwidthMbps.Max))
		}
		args.BaselineEbsBandwidthMbps = ebsBandwidthArgs
	}
	if requirements.AcceleratorCount != nil {
		acceleratorArgs := &ec2.LaunchTemplateInstanceRequirementsAcceleratorCountArgs{}
		if requirements.AcceleratorCount.Min > 0 {
			acceleratorArgs.Min = pulumi.IntPtr(int(requirements.AcceleratorCount.Min))
		}
		if requirements.AcceleratorCount.Max > 0 {
			acceleratorArgs.Max = pulumi.IntPtr(int(requirements.AcceleratorCount.Max))
		}
		args.AcceleratorCount = acceleratorArgs
	}
	if len(requirements.AcceleratorManufacturers) > 0 {
		args.AcceleratorManufacturers = pulumi.ToStringArray(requirements.AcceleratorManufacturers)
	}
	if len(requirements.AcceleratorNames) > 0 {
		args.AcceleratorNames = pulumi.ToStringArray(requirements.AcceleratorNames)
	}
	if len(requirements.AcceleratorTypes) > 0 {
		args.AcceleratorTypes = pulumi.ToStringArray(requirements.AcceleratorTypes)
	}
	if requirements.AcceleratorTotalMemoryMib != nil {
		acceleratorMemoryArgs := &ec2.LaunchTemplateInstanceRequirementsAcceleratorTotalMemoryMibArgs{}
		if requirements.AcceleratorTotalMemoryMib.Min > 0 {
			acceleratorMemoryArgs.Min = pulumi.IntPtr(int(requirements.AcceleratorTotalMemoryMib.Min))
		}
		if requirements.AcceleratorTotalMemoryMib.Max > 0 {
			acceleratorMemoryArgs.Max = pulumi.IntPtr(int(requirements.AcceleratorTotalMemoryMib.Max))
		}
		args.AcceleratorTotalMemoryMib = acceleratorMemoryArgs
	}

	return args
}

// metadataOptionsArgs maps the IMDS posture. Only explicitly set fields are
// sent, so AWS keeps its own defaults for the rest.
func metadataOptionsArgs(options *awslaunchtemplatev1alpha1.AwsLaunchTemplateMetadataOptions) *ec2.LaunchTemplateMetadataOptionsArgs {
	args := &ec2.LaunchTemplateMetadataOptionsArgs{}
	if options.HttpEndpoint != "" {
		args.HttpEndpoint = pulumi.StringPtr(options.HttpEndpoint)
	}
	if options.HttpTokens != "" {
		args.HttpTokens = pulumi.StringPtr(options.HttpTokens)
	}
	if options.HttpPutResponseHopLimit > 0 {
		args.HttpPutResponseHopLimit = pulumi.IntPtr(int(options.HttpPutResponseHopLimit))
	}
	if options.HttpProtocolIpv6 != "" {
		args.HttpProtocolIpv6 = pulumi.StringPtr(options.HttpProtocolIpv6)
	}
	if options.InstanceMetadataTags != "" {
		args.InstanceMetadataTags = pulumi.StringPtr(options.InstanceMetadataTags)
	}
	return args
}

// placementArgs maps instance placement.
func placementArgs(placement *awslaunchtemplatev1alpha1.AwsLaunchTemplatePlacement) *ec2.LaunchTemplatePlacementArgs {
	args := &ec2.LaunchTemplatePlacementArgs{}
	if placement.AvailabilityZone != "" {
		args.AvailabilityZone = pulumi.StringPtr(placement.AvailabilityZone)
	}
	if placement.GroupName != "" {
		args.GroupName = pulumi.StringPtr(placement.GroupName)
	}
	if placement.PartitionNumber > 0 {
		args.PartitionNumber = pulumi.IntPtr(int(placement.PartitionNumber))
	}
	if placement.Tenancy != "" {
		args.Tenancy = pulumi.StringPtr(placement.Tenancy)
	}
	return args
}

// cpuOptionsArgs maps CPU topology and security features.
func cpuOptionsArgs(options *awslaunchtemplatev1alpha1.AwsLaunchTemplateCpuOptions) *ec2.LaunchTemplateCpuOptionsArgs {
	args := &ec2.LaunchTemplateCpuOptionsArgs{}
	if options.CoreCount > 0 {
		args.CoreCount = pulumi.IntPtr(int(options.CoreCount))
	}
	if options.ThreadsPerCore > 0 {
		args.ThreadsPerCore = pulumi.IntPtr(int(options.ThreadsPerCore))
	}
	if options.AmdSevSnp != "" {
		args.AmdSevSnp = pulumi.StringPtr(options.AmdSevSnp)
	}
	return args
}

// spotOptionsArgs maps the Spot purchase options.
func spotOptionsArgs(options *awslaunchtemplatev1alpha1.AwsLaunchTemplateSpotOptions) *ec2.LaunchTemplateInstanceMarketOptionsSpotOptionsArgs {
	args := &ec2.LaunchTemplateInstanceMarketOptionsSpotOptionsArgs{}
	if options.MaxPrice != "" {
		args.MaxPrice = pulumi.StringPtr(options.MaxPrice)
	}
	if options.SpotInstanceType != "" {
		args.SpotInstanceType = pulumi.StringPtr(options.SpotInstanceType)
	}
	if options.InstanceInterruptionBehavior != "" {
		args.InstanceInterruptionBehavior = pulumi.StringPtr(options.InstanceInterruptionBehavior)
	}
	if options.ValidUntil != "" {
		args.ValidUntil = pulumi.StringPtr(options.ValidUntil)
	}
	return args
}

// privateDnsNameOptionsArgs maps the private DNS hostname behavior.
func privateDnsNameOptionsArgs(options *awslaunchtemplatev1alpha1.AwsLaunchTemplatePrivateDnsNameOptions) *ec2.LaunchTemplatePrivateDnsNameOptionsArgs {
	args := &ec2.LaunchTemplatePrivateDnsNameOptionsArgs{}
	if options.HostnameType != "" {
		args.HostnameType = pulumi.StringPtr(options.HostnameType)
	}
	if options.EnableResourceNameDnsARecord {
		args.EnableResourceNameDnsARecord = pulumi.BoolPtr(true)
	}
	if options.EnableResourceNameDnsAaaaRecord {
		args.EnableResourceNameDnsAaaaRecord = pulumi.BoolPtr(true)
	}
	return args
}

// truncateName enforces AWS's 125-character launch template name limit.
func truncateName(name string, maxLen int) string {
	if len(name) <= maxLen {
		return name
	}
	return name[:maxLen]
}
