// Node-level configuration builders for the EC2NodeClass spec: block device
// mappings (EBS) and kubelet configuration.
package module

import (
	kuberneteskarpenterec2nodeclassv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kuberneteskarpenterec2nodeclass/v1alpha1"
	karpenterv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/karpenter/kubernetes/karpenter/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func buildBlockDeviceMappings(mappings []*kuberneteskarpenterec2nodeclassv1alpha1.KubernetesKarpenterEc2NodeClassBlockDeviceMapping) karpenterv1.EC2NodeClassSpecBlockDeviceMappingsArray {
	arr := karpenterv1.EC2NodeClassSpecBlockDeviceMappingsArray{}
	for _, mapping := range mappings {
		mappingArgs := karpenterv1.EC2NodeClassSpecBlockDeviceMappingsArgs{}
		if mapping.GetDeviceName() != "" {
			mappingArgs.DeviceName = pulumi.String(mapping.GetDeviceName())
		}
		if ebs := mapping.GetEbs(); ebs != nil {
			mappingArgs.Ebs = buildEbs(ebs)
		}
		if mapping.RootVolume != nil {
			mappingArgs.RootVolume = pulumi.Bool(mapping.GetRootVolume())
		}
		arr = append(arr, mappingArgs)
	}
	return arr
}

// buildEbs maps the EBS volume parameters. kms_key_id and snapshot_id render
// as the acronym-cased 'kmsKeyID' / 'snapshotID' CR keys — the SDK's
// KmsKeyID / SnapshotID args carry those pulumi tags, matching the CRD.
func buildEbs(ebs *kuberneteskarpenterec2nodeclassv1alpha1.KubernetesKarpenterEc2NodeClassEbs) karpenterv1.EC2NodeClassSpecBlockDeviceMappingsEbsArgs {
	args := karpenterv1.EC2NodeClassSpecBlockDeviceMappingsEbsArgs{}
	if ebs.DeleteOnTermination != nil {
		args.DeleteOnTermination = pulumi.Bool(ebs.GetDeleteOnTermination())
	}
	if ebs.Encrypted != nil {
		args.Encrypted = pulumi.Bool(ebs.GetEncrypted())
	}
	if ebs.Iops != nil {
		args.Iops = pulumi.Int(int(ebs.GetIops()))
	}
	if ebs.KmsKeyId != nil {
		args.KmsKeyID = pulumi.String(ebs.GetKmsKeyId())
	}
	if ebs.SnapshotId != nil {
		args.SnapshotID = pulumi.String(ebs.GetSnapshotId())
	}
	if ebs.Throughput != nil {
		args.Throughput = pulumi.Int(int(ebs.GetThroughput()))
	}
	if ebs.VolumeInitializationRate != nil {
		args.VolumeInitializationRate = pulumi.Int(int(ebs.GetVolumeInitializationRate()))
	}
	if ebs.VolumeSize != nil {
		args.VolumeSize = pulumi.String(ebs.GetVolumeSize())
	}
	if ebs.VolumeType != nil {
		args.VolumeType = pulumi.String(ebs.GetVolumeType())
	}
	return args
}

// buildKubelet maps all 12 kubelet fields. cluster_dns and cpu_cfs_quota
// render as the acronym-cased 'clusterDNS' / 'cpuCFSQuota' CR keys, and the
// image GC thresholds as 'imageGCHighThresholdPercent' /
// 'imageGCLowThresholdPercent' — the SDK's ClusterDNS / CpuCFSQuota /
// ImageGC*ThresholdPercent args carry those pulumi tags, matching the CRD.
func buildKubelet(kubelet *kuberneteskarpenterec2nodeclassv1alpha1.KubernetesKarpenterEc2NodeClassKubelet) karpenterv1.EC2NodeClassSpecKubeletArgs {
	args := karpenterv1.EC2NodeClassSpecKubeletArgs{}
	if clusterDns := kubelet.GetClusterDns(); len(clusterDns) > 0 {
		args.ClusterDNS = pulumi.ToStringArray(clusterDns)
	}
	if kubelet.CpuCfsQuota != nil {
		args.CpuCFSQuota = pulumi.Bool(kubelet.GetCpuCfsQuota())
	}
	if evictionHard := kubelet.GetEvictionHard(); len(evictionHard) > 0 {
		args.EvictionHard = pulumi.ToStringMap(evictionHard)
	}
	if kubelet.EvictionMaxPodGracePeriod != nil {
		args.EvictionMaxPodGracePeriod = pulumi.Int(int(kubelet.GetEvictionMaxPodGracePeriod()))
	}
	if evictionSoft := kubelet.GetEvictionSoft(); len(evictionSoft) > 0 {
		args.EvictionSoft = pulumi.ToStringMap(evictionSoft)
	}
	if evictionSoftGracePeriod := kubelet.GetEvictionSoftGracePeriod(); len(evictionSoftGracePeriod) > 0 {
		args.EvictionSoftGracePeriod = pulumi.ToStringMap(evictionSoftGracePeriod)
	}
	if kubelet.ImageGcHighThresholdPercent != nil {
		args.ImageGCHighThresholdPercent = pulumi.Int(int(kubelet.GetImageGcHighThresholdPercent()))
	}
	if kubelet.ImageGcLowThresholdPercent != nil {
		args.ImageGCLowThresholdPercent = pulumi.Int(int(kubelet.GetImageGcLowThresholdPercent()))
	}
	if kubeReserved := kubelet.GetKubeReserved(); len(kubeReserved) > 0 {
		args.KubeReserved = pulumi.ToStringMap(kubeReserved)
	}
	if kubelet.MaxPods != nil {
		args.MaxPods = pulumi.Int(int(kubelet.GetMaxPods()))
	}
	if kubelet.PodsPerCore != nil {
		args.PodsPerCore = pulumi.Int(int(kubelet.GetPodsPerCore()))
	}
	if systemReserved := kubelet.GetSystemReserved(); len(systemReserved) > 0 {
		args.SystemReserved = pulumi.ToStringMap(systemReserved)
	}
	return args
}
