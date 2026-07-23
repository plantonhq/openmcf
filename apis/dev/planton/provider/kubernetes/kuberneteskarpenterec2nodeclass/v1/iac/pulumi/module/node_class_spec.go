package module

import (
	kuberneteskarpenterec2nodeclassv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskarpenterec2nodeclass/v1"
	karpenterv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/karpenter/kubernetes/karpenter/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildNodeClassSpec maps the proto spec onto the typed crd2pulumi
// EC2NodeClass spec args, field by field. Every optional is set only when
// present in the manifest: the CRD carries server-side defaults (IMDS
// metadataOptions in particular) that a rendered zero value would override.
//
// Acronym-cased CR keys (associatePublicIPAddress, ownerID, kmsKeyID,
// snapshotID, clusterDNS, cpuCFSQuota, imageGC*ThresholdPercent,
// httpProtocolIPv6) are pinned via json_name in the proto and match the
// generated SDK's pulumi tags — verified against the CRD at the pinned
// karpenter release.
func buildNodeClassSpec(spec *kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassSpec) karpenterv1.EC2NodeClassSpecArgs {
	args := karpenterv1.EC2NodeClassSpecArgs{
		AmiSelectorTerms:           buildAmiSelectorTerms(spec.GetAmiSelectorTerms()),
		SubnetSelectorTerms:        buildSubnetSelectorTerms(spec.GetSubnetSelectorTerms()),
		SecurityGroupSelectorTerms: buildSecurityGroupSelectorTerms(spec.GetSecurityGroupSelectorTerms()),
	}

	if spec.AmiFamily != nil {
		args.AmiFamily = pulumi.String(spec.GetAmiFamily())
	}

	// role XOR instance_profile — protovalidate guarantees exactly one arm,
	// so exactly one lands in the rendered CR.
	if spec.Role != nil {
		args.Role = pulumi.String(spec.GetRole())
	}
	if spec.InstanceProfile != nil {
		args.InstanceProfile = pulumi.String(spec.GetInstanceProfile())
	}

	if blockDeviceMappings := spec.GetBlockDeviceMappings(); len(blockDeviceMappings) > 0 {
		args.BlockDeviceMappings = buildBlockDeviceMappings(blockDeviceMappings)
	}

	if capacityReservationTerms := spec.GetCapacityReservationSelectorTerms(); len(capacityReservationTerms) > 0 {
		args.CapacityReservationSelectorTerms = buildCapacityReservationSelectorTerms(capacityReservationTerms)
	}

	if spec.AssociatePublicIpAddress != nil {
		args.AssociatePublicIPAddress = pulumi.Bool(spec.GetAssociatePublicIpAddress())
	}

	if connectionTracking := spec.GetConnectionTracking(); connectionTracking != nil {
		args.ConnectionTracking = buildConnectionTracking(connectionTracking)
	}

	if spec.GetContext() != "" {
		args.Context = pulumi.String(spec.GetContext())
	}

	if cpuOptions := spec.GetCpuOptions(); cpuOptions != nil {
		cpuOptionsArgs := karpenterv1.EC2NodeClassSpecCpuOptionsArgs{}
		if cpuOptions.NestedVirtualization != nil {
			cpuOptionsArgs.NestedVirtualization = pulumi.String(cpuOptions.GetNestedVirtualization())
		}
		args.CpuOptions = cpuOptionsArgs
	}

	// plain proto3 bool: false is indistinguishable from unset, and false is
	// the AWS default anyway — only render the key when enabling.
	if spec.GetDetailedMonitoring() {
		args.DetailedMonitoring = pulumi.Bool(true)
	}

	if spec.InstanceStorePolicy != nil {
		args.InstanceStorePolicy = pulumi.String(spec.GetInstanceStorePolicy())
	}

	if spec.IpPrefixCount != nil {
		args.IpPrefixCount = pulumi.Int(int(spec.GetIpPrefixCount()))
	}

	if kubelet := spec.GetKubelet(); kubelet != nil {
		args.Kubelet = buildKubelet(kubelet)
	}

	if metadataOptions := spec.GetMetadataOptions(); metadataOptions != nil {
		args.MetadataOptions = buildMetadataOptions(metadataOptions)
	}

	if networkInterfaces := spec.GetNetworkInterfaces(); len(networkInterfaces) > 0 {
		args.NetworkInterfaces = buildNetworkInterfaces(networkInterfaces)
	}

	if placementGroupSelector := spec.GetPlacementGroupSelector(); placementGroupSelector != nil {
		placementArgs := karpenterv1.EC2NodeClassSpecPlacementGroupSelectorArgs{}
		if placementGroupSelector.Name != nil {
			placementArgs.Name = pulumi.String(placementGroupSelector.GetName())
		}
		if placementGroupSelector.Id != nil {
			placementArgs.Id = pulumi.String(placementGroupSelector.GetId())
		}
		args.PlacementGroupSelector = placementArgs
	}

	if tags := spec.GetTags(); len(tags) > 0 {
		args.Tags = pulumi.ToStringMap(tags)
	}

	if spec.GetUserData() != "" {
		args.UserData = pulumi.String(spec.GetUserData())
	}

	return args
}

func buildConnectionTracking(connectionTracking *kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassConnectionTracking) karpenterv1.EC2NodeClassSpecConnectionTrackingArgs {
	args := karpenterv1.EC2NodeClassSpecConnectionTrackingArgs{}
	if connectionTracking.TcpEstablishedTimeout != nil {
		args.TcpEstablishedTimeout = pulumi.Int(int(connectionTracking.GetTcpEstablishedTimeout()))
	}
	if connectionTracking.UdpStreamTimeout != nil {
		args.UdpStreamTimeout = pulumi.Int(int(connectionTracking.GetUdpStreamTimeout()))
	}
	if connectionTracking.UdpTimeout != nil {
		args.UdpTimeout = pulumi.Int(int(connectionTracking.GetUdpTimeout()))
	}
	return args
}

// buildMetadataOptions maps the IMDS block. Fields are set only when present:
// the CRD defaults this block to the EKS security best practice (endpoint
// enabled, IPv6 disabled, hop limit 1, tokens required) and an absent key
// lets those defaults apply.
func buildMetadataOptions(metadataOptions *kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassMetadataOptions) karpenterv1.EC2NodeClassSpecMetadataOptionsArgs {
	args := karpenterv1.EC2NodeClassSpecMetadataOptionsArgs{}
	if metadataOptions.HttpEndpoint != nil {
		args.HttpEndpoint = pulumi.String(metadataOptions.GetHttpEndpoint())
	}
	if metadataOptions.HttpProtocolIpv6 != nil {
		args.HttpProtocolIPv6 = pulumi.String(metadataOptions.GetHttpProtocolIpv6())
	}
	if metadataOptions.HttpPutResponseHopLimit != nil {
		args.HttpPutResponseHopLimit = pulumi.Int(int(metadataOptions.GetHttpPutResponseHopLimit()))
	}
	if metadataOptions.HttpTokens != nil {
		args.HttpTokens = pulumi.String(metadataOptions.GetHttpTokens())
	}
	return args
}

// buildNetworkInterfaces maps the explicit network-interface layout. Unlike
// the rest of the spec, all three fields are ALWAYS rendered: the CRD marks
// deviceIndex, interfaceType and networkCardIndex required on every item, so
// zero is a legitimate value (device 0 / card 0 is the primary), never an
// omission.
func buildNetworkInterfaces(networkInterfaces []*kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassNetworkInterface) karpenterv1.EC2NodeClassSpecNetworkInterfacesArray {
	arr := karpenterv1.EC2NodeClassSpecNetworkInterfacesArray{}
	for _, networkInterface := range networkInterfaces {
		arr = append(arr, karpenterv1.EC2NodeClassSpecNetworkInterfacesArgs{
			DeviceIndex:      pulumi.Int(int(networkInterface.GetDeviceIndex())),
			InterfaceType:    pulumi.String(networkInterface.GetInterfaceType()),
			NetworkCardIndex: pulumi.Int(int(networkInterface.GetNetworkCardIndex())),
		})
	}
	return arr
}
