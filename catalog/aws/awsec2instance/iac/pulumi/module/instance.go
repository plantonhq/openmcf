package module

import (
	"github.com/pkg/errors"
	awsec2instancev1alpha1 "github.com/plantonhq/planton/catalog/aws/awsec2instance/v1alpha1"
	"github.com/plantonhq/planton/internal/valuefrom"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// instance provisions the EC2 instance. Launch identity (AMI, type,
// subnet, key pair, placement, CPU topology, purchase option) is
// create-time in EC2 -- changing those fields replaces or restarts the
// instance -- while the operational posture (security groups, IAM
// profile, metadata options, protections, monitoring) edits in place.
//
// Field-presence discipline: optional scalars are passed through only
// when set, so an omitted field keeps the AWS (or launch-template)
// default instead of forcing a zero value -- the same tri-state contract
// the Terraform module's generated variables carry.
func instance(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*ec2.Instance, error) {
	spec := locals.AwsEc2Instance.Spec

	args := &ec2.InstanceArgs{
		// EC2 has no name argument -- the Name tag IS the instance's
		// display identity, carried inside the identity tag set built
		// from metadata (both engines emit the same keys and values).
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	// --- Launch source ------------------------------------------------
	// AMI and type are each optional here when a launch template
	// supplies them (CEL enforces that at validate time); set only what
	// the spec carries so the template's values win for the rest.
	if spec.Ami != "" {
		args.Ami = pulumi.String(spec.Ami)
	}
	if spec.InstanceType != "" {
		args.InstanceType = pulumi.String(spec.InstanceType)
	}
	if spec.LaunchTemplate != nil {
		templateArgs := &ec2.InstanceLaunchTemplateArgs{}
		if spec.LaunchTemplate.Id.GetValue() != "" {
			templateArgs.Id = pulumi.String(spec.LaunchTemplate.Id.GetValue())
		}
		if spec.LaunchTemplate.Name != "" {
			templateArgs.Name = pulumi.String(spec.LaunchTemplate.Name)
		}
		if spec.LaunchTemplate.Version != "" {
			templateArgs.Version = pulumi.String(spec.LaunchTemplate.Version)
		}
		args.LaunchTemplate = templateArgs
	}

	// --- IAM / access --------------------------------------------------
	// The EC2 API takes the instance profile by NAME (launch templates
	// accept an ARN; instances do not) -- the spec's ref resolves the
	// profile's name output for exactly this reason.
	if spec.InstanceProfile.GetValue() != "" {
		args.IamInstanceProfile = pulumi.String(spec.InstanceProfile.GetValue())
	}
	if spec.KeyName != "" {
		args.KeyName = pulumi.String(spec.KeyName)
	}

	// --- Networking -----------------------------------------------------
	// Either a pre-provisioned primary ENI carries the network identity,
	// or the inline fields shape a new primary interface (CEL guarantees
	// the two never mix).
	if spec.PrimaryNetworkInterfaceId != "" {
		args.PrimaryNetworkInterface = &ec2.InstancePrimaryNetworkInterfaceArgs{
			NetworkInterfaceId: pulumi.String(spec.PrimaryNetworkInterfaceId),
		}
	}
	if spec.SubnetId.GetValue() != "" {
		args.SubnetId = pulumi.String(spec.SubnetId.GetValue())
	}
	if len(spec.SecurityGroupIds) > 0 {
		args.VpcSecurityGroupIds = pulumi.ToStringArray(valuefrom.ToStringArray(spec.SecurityGroupIds))
	}
	if spec.PrivateIp != "" {
		args.PrivateIp = pulumi.String(spec.PrivateIp)
	}
	if len(spec.SecondaryPrivateIps) > 0 {
		args.SecondaryPrivateIps = pulumi.ToStringArray(spec.SecondaryPrivateIps)
	}
	// Tri-state: unset inherits the subnet's map-public-IP setting.
	if spec.AssociatePublicIpAddress != nil {
		args.AssociatePublicIpAddress = pulumi.Bool(spec.GetAssociatePublicIpAddress())
	}
	// Platform middleware defaults this to true (the AWS default); an
	// explicit false is the NAT/router posture.
	if spec.SourceDestCheck != nil {
		args.SourceDestCheck = pulumi.Bool(spec.GetSourceDestCheck())
	}
	if spec.Ipv6AddressCount != 0 {
		args.Ipv6AddressCount = pulumi.Int(int(spec.Ipv6AddressCount))
	}
	if len(spec.Ipv6Addresses) > 0 {
		args.Ipv6Addresses = pulumi.ToStringArray(spec.Ipv6Addresses)
	}
	if spec.EnablePrimaryIpv6 != nil {
		args.EnablePrimaryIpv6 = pulumi.Bool(spec.GetEnablePrimaryIpv6())
	}
	if spec.PrivateDnsNameOptions != nil {
		dnsArgs := &ec2.InstancePrivateDnsNameOptionsArgs{
			EnableResourceNameDnsARecord:    pulumi.Bool(spec.PrivateDnsNameOptions.EnableResourceNameDnsARecord),
			EnableResourceNameDnsAaaaRecord: pulumi.Bool(spec.PrivateDnsNameOptions.EnableResourceNameDnsAaaaRecord),
		}
		if spec.PrivateDnsNameOptions.HostnameType != "" {
			dnsArgs.HostnameType = pulumi.String(spec.PrivateDnsNameOptions.HostnameType)
		}
		args.PrivateDnsNameOptions = dnsArgs
	}
	if len(spec.SecondaryNetworkInterfaces) > 0 {
		interfaces := ec2.InstanceSecondaryNetworkInterfaceArray{}
		for _, secondaryInterface := range spec.SecondaryNetworkInterfaces {
			interfaceArgs := ec2.InstanceSecondaryNetworkInterfaceArgs{
				NetworkCardIndex:  pulumi.Int(int(secondaryInterface.NetworkCardIndex)),
				SecondarySubnetId: pulumi.String(secondaryInterface.SubnetId.GetValue()),
			}
			if secondaryInterface.DeviceIndex != 0 {
				interfaceArgs.DeviceIndex = pulumi.Int(int(secondaryInterface.DeviceIndex))
			}
			if secondaryInterface.PrivateIpAddressCount != 0 {
				interfaceArgs.PrivateIpAddressCount = pulumi.Int(int(secondaryInterface.PrivateIpAddressCount))
			}
			if secondaryInterface.DeleteOnTermination != nil {
				interfaceArgs.DeleteOnTermination = pulumi.Bool(secondaryInterface.GetDeleteOnTermination())
			}
			interfaces = append(interfaces, interfaceArgs)
		}
		args.SecondaryNetworkInterfaces = interfaces
	}

	// --- Storage ---------------------------------------------------------
	if spec.RootBlockDevice != nil {
		args.RootBlockDevice = rootBlockDeviceArgs(spec.RootBlockDevice)
	}
	if len(spec.EbsBlockDevices) > 0 {
		devices := ec2.InstanceEbsBlockDeviceArray{}
		for _, device := range spec.EbsBlockDevices {
			devices = append(devices, ebsBlockDeviceArgs(device))
		}
		args.EbsBlockDevices = devices
	}
	if len(spec.EphemeralBlockDevices) > 0 {
		devices := ec2.InstanceEphemeralBlockDeviceArray{}
		for _, device := range spec.EphemeralBlockDevices {
			deviceArgs := ec2.InstanceEphemeralBlockDeviceArgs{
				DeviceName: pulumi.String(device.DeviceName),
			}
			if device.VirtualName != "" {
				deviceArgs.VirtualName = pulumi.String(device.VirtualName)
			}
			if device.NoDevice {
				deviceArgs.NoDevice = pulumi.Bool(true)
			}
			devices = append(devices, deviceArgs)
		}
		args.EphemeralBlockDevices = devices
	}
	// Only forced on: most current-generation types are EBS-optimized by
	// default at no charge, so an unset field keeps the type's default
	// rather than pinning false.
	if spec.EbsOptimized {
		args.EbsOptimized = pulumi.Bool(true)
	}

	// --- Instance posture -------------------------------------------------
	if spec.MetadataOptions != nil {
		metadataArgs := &ec2.InstanceMetadataOptionsArgs{}
		if spec.MetadataOptions.HttpEndpoint != "" {
			metadataArgs.HttpEndpoint = pulumi.String(spec.MetadataOptions.HttpEndpoint)
		}
		if spec.MetadataOptions.HttpTokens != "" {
			metadataArgs.HttpTokens = pulumi.String(spec.MetadataOptions.HttpTokens)
		}
		if spec.MetadataOptions.HttpPutResponseHopLimit != 0 {
			metadataArgs.HttpPutResponseHopLimit = pulumi.Int(int(spec.MetadataOptions.HttpPutResponseHopLimit))
		}
		if spec.MetadataOptions.HttpProtocolIpv6 != "" {
			metadataArgs.HttpProtocolIpv6 = pulumi.String(spec.MetadataOptions.HttpProtocolIpv6)
		}
		if spec.MetadataOptions.InstanceMetadataTags != "" {
			metadataArgs.InstanceMetadataTags = pulumi.String(spec.MetadataOptions.InstanceMetadataTags)
		}
		args.MetadataOptions = metadataArgs
	}
	if spec.DetailedMonitoring {
		args.Monitoring = pulumi.Bool(true)
	}
	if spec.CpuOptions != nil {
		cpuArgs := &ec2.InstanceCpuOptionsArgs{}
		if spec.CpuOptions.CoreCount != 0 {
			cpuArgs.CoreCount = pulumi.Int(int(spec.CpuOptions.CoreCount))
		}
		if spec.CpuOptions.ThreadsPerCore != 0 {
			cpuArgs.ThreadsPerCore = pulumi.Int(int(spec.CpuOptions.ThreadsPerCore))
		}
		if spec.CpuOptions.AmdSevSnp != "" {
			cpuArgs.AmdSevSnp = pulumi.String(spec.CpuOptions.AmdSevSnp)
		}
		if spec.CpuOptions.NestedVirtualization != "" {
			cpuArgs.NestedVirtualization = pulumi.String(spec.CpuOptions.NestedVirtualization)
		}
		args.CpuOptions = cpuArgs
	}
	if spec.CpuCredits != "" {
		args.CreditSpecification = &ec2.InstanceCreditSpecificationArgs{
			CpuCredits: pulumi.String(spec.CpuCredits),
		}
	}

	// --- Purchase options ---------------------------------------------------
	// The purchase market: an explicit market_type, or presence of
	// spot_options implying "spot" (the classic shape). On-Demand needs no
	// market options at all. Reservation-backed markets (capacity-block,
	// interruptible-capacity-reservation) carry no spot_options -- CEL
	// enforces that pairing and the required reservation target.
	if spec.MarketType != "" || spec.SpotOptions != nil {
		marketArgs := &ec2.InstanceInstanceMarketOptionsArgs{
			MarketType: pulumi.String("spot"),
		}
		if spec.MarketType != "" {
			marketArgs.MarketType = pulumi.String(spec.MarketType)
		}
		if spec.SpotOptions != nil {
			spotArgs := &ec2.InstanceInstanceMarketOptionsSpotOptionsArgs{}
			if spec.SpotOptions.MaxPrice != "" {
				spotArgs.MaxPrice = pulumi.String(spec.SpotOptions.MaxPrice)
			}
			if spec.SpotOptions.SpotInstanceType != "" {
				spotArgs.SpotInstanceType = pulumi.String(spec.SpotOptions.SpotInstanceType)
			}
			if spec.SpotOptions.InstanceInterruptionBehavior != "" {
				spotArgs.InstanceInterruptionBehavior = pulumi.String(spec.SpotOptions.InstanceInterruptionBehavior)
			}
			if spec.SpotOptions.ValidUntil != "" {
				spotArgs.ValidUntil = pulumi.String(spec.SpotOptions.ValidUntil)
			}
			marketArgs.SpotOptions = spotArgs
		}
		args.InstanceMarketOptions = marketArgs
	}
	if spec.CapacityReservation != nil {
		reservationArgs := &ec2.InstanceCapacityReservationSpecificationArgs{}
		if spec.CapacityReservation.Preference != "" {
			reservationArgs.CapacityReservationPreference = pulumi.String(spec.CapacityReservation.Preference)
		}
		if spec.CapacityReservation.CapacityReservationId != "" || spec.CapacityReservation.CapacityReservationResourceGroupArn != "" {
			targetArgs := &ec2.InstanceCapacityReservationSpecificationCapacityReservationTargetArgs{}
			if spec.CapacityReservation.CapacityReservationId != "" {
				targetArgs.CapacityReservationId = pulumi.String(spec.CapacityReservation.CapacityReservationId)
			}
			if spec.CapacityReservation.CapacityReservationResourceGroupArn != "" {
				targetArgs.CapacityReservationResourceGroupArn = pulumi.String(spec.CapacityReservation.CapacityReservationResourceGroupArn)
			}
			reservationArgs.CapacityReservationTarget = targetArgs
		}
		args.CapacityReservationSpecification = reservationArgs
	}

	// --- Placement -------------------------------------------------------------
	if spec.Placement != nil {
		if spec.Placement.AvailabilityZone != "" {
			args.AvailabilityZone = pulumi.String(spec.Placement.AvailabilityZone)
		}
		if spec.Placement.GroupName != "" {
			args.PlacementGroup = pulumi.String(spec.Placement.GroupName)
		}
		if spec.Placement.GroupId != "" {
			args.PlacementGroupId = pulumi.String(spec.Placement.GroupId)
		}
		if spec.Placement.PartitionNumber != 0 {
			args.PlacementPartitionNumber = pulumi.Int(int(spec.Placement.PartitionNumber))
		}
		if spec.Placement.Tenancy != "" {
			args.Tenancy = pulumi.String(spec.Placement.Tenancy)
		}
		if spec.Placement.HostId != "" {
			args.HostId = pulumi.String(spec.Placement.HostId)
		}
		if spec.Placement.HostResourceGroupArn != "" {
			args.HostResourceGroupArn = pulumi.String(spec.Placement.HostResourceGroupArn)
		}
	}

	// --- Lifecycle protections and recovery ---------------------------------
	if spec.EnclaveEnabled {
		args.EnclaveOptions = &ec2.InstanceEnclaveOptionsArgs{Enabled: pulumi.Bool(true)}
	}
	if spec.HibernationEnabled {
		args.Hibernation = pulumi.Bool(true)
	}
	if spec.AutoRecovery != "" {
		args.MaintenanceOptions = &ec2.InstanceMaintenanceOptionsArgs{
			AutoRecovery: pulumi.String(spec.AutoRecovery),
		}
	}
	if spec.InstanceInitiatedShutdownBehavior != "" {
		args.InstanceInitiatedShutdownBehavior = pulumi.String(spec.InstanceInitiatedShutdownBehavior)
	}
	if spec.DisableApiStop != nil {
		args.DisableApiStop = pulumi.Bool(spec.GetDisableApiStop())
	}
	if spec.DisableApiTermination != nil {
		args.DisableApiTermination = pulumi.Bool(spec.GetDisableApiTermination())
	}

	// --- User data -----------------------------------------------------------
	// Plain text passes through untouched (the provider hashes it for
	// state); template-looking content like ${HOME} needs no escaping on
	// this engine.
	if spec.UserData != "" {
		args.UserData = pulumi.String(spec.UserData)
	}
	if spec.UserDataBase64 != "" {
		args.UserDataBase64 = pulumi.String(spec.UserDataBase64)
	}
	if spec.UserDataReplaceOnChange {
		args.UserDataReplaceOnChange = pulumi.Bool(true)
	}

	// Uniform at-creation tags for EVERY volume (incl. AMI-inherited
	// mappings) -- the ABAC-compliant arm; mutually exclusive with the
	// per-device tags (CEL enforces it).
	if len(spec.VolumeTags) > 0 {
		args.VolumeTags = pulumi.ToStringMap(spec.VolumeTags)
	}

	// Destroy-time escape hatch: lift stop/termination protection before
	// terminating instead of failing the destroy.
	if spec.ForceDestroy {
		args.ForceDestroy = pulumi.Bool(true)
	}

	// Stable Pulumi resource name; the cloud identity travels in the
	// Name tag (EC2's only name), never in Pulumi auto-naming.
	createdInstance, err := ec2.NewInstance(ctx, "instance", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create EC2 instance")
	}

	return createdInstance, nil
}

// rootBlockDeviceArgs lowers the root-volume override. Unset fields stay
// nil so the AMI's block-device mapping keeps deciding them.
func rootBlockDeviceArgs(device *awsec2instancev1alpha1.AwsEc2InstanceRootBlockDevice) *ec2.InstanceRootBlockDeviceArgs {
	rootArgs := &ec2.InstanceRootBlockDeviceArgs{}
	if device.VolumeSizeGb != 0 {
		rootArgs.VolumeSize = pulumi.Int(int(device.VolumeSizeGb))
	}
	if device.VolumeType != "" {
		rootArgs.VolumeType = pulumi.String(device.VolumeType)
	}
	if device.Iops != 0 {
		rootArgs.Iops = pulumi.Int(int(device.Iops))
	}
	if device.ThroughputMibps != 0 {
		rootArgs.Throughput = pulumi.Int(int(device.ThroughputMibps))
	}
	if device.Encrypted {
		rootArgs.Encrypted = pulumi.Bool(true)
	}
	if device.KmsKeyId.GetValue() != "" {
		rootArgs.KmsKeyId = pulumi.String(device.KmsKeyId.GetValue())
	}
	if device.DeleteOnTermination != nil {
		rootArgs.DeleteOnTermination = pulumi.Bool(device.GetDeleteOnTermination())
	}
	// Per-volume tags apply post-creation (see the spec-level volume_tags
	// for the at-creation alternative; the provider forbids mixing them).
	if len(device.Tags) > 0 {
		rootArgs.Tags = pulumi.ToStringMap(device.Tags)
	}
	return rootArgs
}

// ebsBlockDeviceArgs lowers one additional data volume mapping.
func ebsBlockDeviceArgs(device *awsec2instancev1alpha1.AwsEc2InstanceEbsBlockDevice) ec2.InstanceEbsBlockDeviceArgs {
	deviceArgs := ec2.InstanceEbsBlockDeviceArgs{
		DeviceName: pulumi.String(device.DeviceName),
	}
	if device.VolumeSizeGb != 0 {
		deviceArgs.VolumeSize = pulumi.Int(int(device.VolumeSizeGb))
	}
	if device.VolumeType != "" {
		deviceArgs.VolumeType = pulumi.String(device.VolumeType)
	}
	if device.Iops != 0 {
		deviceArgs.Iops = pulumi.Int(int(device.Iops))
	}
	if device.ThroughputMibps != 0 {
		deviceArgs.Throughput = pulumi.Int(int(device.ThroughputMibps))
	}
	if device.Encrypted {
		deviceArgs.Encrypted = pulumi.Bool(true)
	}
	if device.KmsKeyId.GetValue() != "" {
		deviceArgs.KmsKeyId = pulumi.String(device.KmsKeyId.GetValue())
	}
	if device.SnapshotId != "" {
		deviceArgs.SnapshotId = pulumi.String(device.SnapshotId)
	}
	if device.DeleteOnTermination != nil {
		deviceArgs.DeleteOnTermination = pulumi.Bool(device.GetDeleteOnTermination())
	}
	if len(device.Tags) > 0 {
		deviceArgs.Tags = pulumi.ToStringMap(device.Tags)
	}
	return deviceArgs
}
