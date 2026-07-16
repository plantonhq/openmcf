package module

import (
	"strings"

	azureaksclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azureakscluster/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AzureAksCluster *azureaksclusterv1.AzureAksCluster

	// ResourceGroupName is a StringValueOrRef field; the platform middleware
	// resolves valueFrom references before IaC modules run, so GetValue()
	// always returns the resolved literal name.
	ResourceGroupName string

	// AzureTags is the metadata-derived tag map with the spec's user tags
	// merged over it (user tags win on key collision), mirroring the
	// Terraform module's merge order.
	AzureTags map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *azureaksclusterv1.AzureAksClusterStackInput) *Locals {
	locals := &Locals{}

	locals.AzureAksCluster = stackInput.Target
	target := stackInput.Target

	locals.ResourceGroupName = target.Spec.ResourceGroup.GetValue()

	// Metadata-derived tags first, then the user's spec tags merged over
	// them: user tags deliberately win so an org's governance conventions
	// (cost center, owner) can override the derived values where they
	// collide.
	locals.AzureTags = map[string]string{
		// PARITY-EXCEPTION: resource_kind here is the lowered
		// CloudResourceKind enum string and resource_id is omitted when
		// metadata.id is empty, while the Terraform module emits the
		// family-wide snake-case literal and falls back to metadata.name.
		// Output-neutral (tags never feed stack outputs); aligning the two
		// shapes is a family-wide convention change, not a per-kind fix.
		"resource":      "true",
		"resource_name": target.Metadata.Name,
		"resource_kind": strings.ToLower(cloudresourcekind.CloudResourceKind_AzureAksCluster.String()),
	}

	if target.Metadata.Id != "" {
		locals.AzureTags["resource_id"] = target.Metadata.Id
	}

	if target.Metadata.Org != "" {
		locals.AzureTags["organization"] = target.Metadata.Org
	}

	if target.Metadata.Env != "" {
		locals.AzureTags["environment"] = target.Metadata.Env
	}

	return locals
}

// enum-name -> ARM string maps shared by this module. Empty results mean
// "leave unset so Azure applies its default" -- keeping an unspecified
// spec and Azure's default identical on both engines.

func skuTierToArm(v azureaksclusterv1.AzureAksClusterSkuTier) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterSkuTier_FREE:
		return "Free"
	case azureaksclusterv1.AzureAksClusterSkuTier_STANDARD:
		return "Standard"
	case azureaksclusterv1.AzureAksClusterSkuTier_PREMIUM:
		return "Premium"
	}
	return ""
}

func supportPlanToArm(v azureaksclusterv1.AzureAksClusterSupportPlan) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterSupportPlan_KUBERNETES_OFFICIAL:
		return "KubernetesOfficial"
	case azureaksclusterv1.AzureAksClusterSupportPlan_AKS_LONG_TERM_SUPPORT:
		return "AKSLongTermSupport"
	}
	return ""
}

func upgradeChannelToArm(v azureaksclusterv1.AzureAksClusterUpgradeChannel) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterUpgradeChannel_PATCH:
		return "patch"
	case azureaksclusterv1.AzureAksClusterUpgradeChannel_STABLE:
		return "stable"
	case azureaksclusterv1.AzureAksClusterUpgradeChannel_RAPID:
		return "rapid"
	case azureaksclusterv1.AzureAksClusterUpgradeChannel_NODE_IMAGE:
		return "node-image"
	}
	return ""
}

func nodeOsUpgradeChannelToArm(v azureaksclusterv1.AzureAksClusterNodeOsUpgradeChannel) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterNodeOsUpgradeChannel_NODE_OS_NODE_IMAGE:
		return "NodeImage"
	case azureaksclusterv1.AzureAksClusterNodeOsUpgradeChannel_SECURITY_PATCH:
		return "SecurityPatch"
	case azureaksclusterv1.AzureAksClusterNodeOsUpgradeChannel_UNMANAGED:
		return "Unmanaged"
	case azureaksclusterv1.AzureAksClusterNodeOsUpgradeChannel_NODE_OS_NONE:
		return "None"
	}
	return ""
}

func osSkuToArm(v azureaksclusterv1.AzureAksClusterOsSku) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterOsSku_UBUNTU:
		return "Ubuntu"
	case azureaksclusterv1.AzureAksClusterOsSku_UBUNTU_2204:
		return "Ubuntu2204"
	case azureaksclusterv1.AzureAksClusterOsSku_UBUNTU_2404:
		return "Ubuntu2404"
	case azureaksclusterv1.AzureAksClusterOsSku_AZURE_LINUX:
		return "AzureLinux"
	case azureaksclusterv1.AzureAksClusterOsSku_AZURE_LINUX_3:
		return "AzureLinux3"
	case azureaksclusterv1.AzureAksClusterOsSku_WINDOWS_2019:
		return "Windows2019"
	case azureaksclusterv1.AzureAksClusterOsSku_WINDOWS_2022:
		return "Windows2022"
	}
	return ""
}

func osDiskTypeToArm(v azureaksclusterv1.AzureAksClusterOsDiskType) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterOsDiskType_MANAGED:
		return "Managed"
	case azureaksclusterv1.AzureAksClusterOsDiskType_EPHEMERAL:
		return "Ephemeral"
	}
	return ""
}

func kubeletDiskTypeToArm(v azureaksclusterv1.AzureAksClusterKubeletDiskType) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterKubeletDiskType_OS:
		return "OS"
	case azureaksclusterv1.AzureAksClusterKubeletDiskType_TEMPORARY:
		return "Temporary"
	}
	return ""
}

func gpuInstanceToArm(v azureaksclusterv1.AzureAksClusterGpuInstance) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterGpuInstance_MIG1G:
		return "MIG1g"
	case azureaksclusterv1.AzureAksClusterGpuInstance_MIG2G:
		return "MIG2g"
	case azureaksclusterv1.AzureAksClusterGpuInstance_MIG3G:
		return "MIG3g"
	case azureaksclusterv1.AzureAksClusterGpuInstance_MIG4G:
		return "MIG4g"
	case azureaksclusterv1.AzureAksClusterGpuInstance_MIG7G:
		return "MIG7g"
	}
	return ""
}

func gpuDriverToArm(v azureaksclusterv1.AzureAksClusterGpuDriver) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterGpuDriver_INSTALL:
		return "Install"
	case azureaksclusterv1.AzureAksClusterGpuDriver_NONE:
		return "None"
	}
	return ""
}

func scaleDownModeToArm(v azureaksclusterv1.AzureAksClusterScaleDownMode) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterScaleDownMode_DELETE:
		return "Delete"
	case azureaksclusterv1.AzureAksClusterScaleDownMode_DEALLOCATE:
		return "Deallocate"
	}
	return ""
}

func workloadRuntimeToArm(v azureaksclusterv1.AzureAksClusterWorkloadRuntime) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterWorkloadRuntime_OCI_CONTAINER:
		return "OCIContainer"
	case azureaksclusterv1.AzureAksClusterWorkloadRuntime_KATA_MSHV_VM_ISOLATION:
		return "KataMshvVmIsolation"
	}
	return ""
}

func cpuManagerPolicyToArm(v azureaksclusterv1.AzureAksClusterCpuManagerPolicy) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterCpuManagerPolicy_CPU_MANAGER_NONE:
		return "none"
	case azureaksclusterv1.AzureAksClusterCpuManagerPolicy_CPU_MANAGER_STATIC:
		return "static"
	}
	return ""
}

func topologyManagerPolicyToArm(v azureaksclusterv1.AzureAksClusterTopologyManagerPolicy) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterTopologyManagerPolicy_TOPOLOGY_NONE:
		return "none"
	case azureaksclusterv1.AzureAksClusterTopologyManagerPolicy_BEST_EFFORT:
		return "best-effort"
	case azureaksclusterv1.AzureAksClusterTopologyManagerPolicy_RESTRICTED:
		return "restricted"
	case azureaksclusterv1.AzureAksClusterTopologyManagerPolicy_SINGLE_NUMA_NODE:
		return "single-numa-node"
	}
	return ""
}

func transparentHugePageToArm(v azureaksclusterv1.AzureAksClusterTransparentHugePage) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterTransparentHugePage_THP_ALWAYS:
		return "always"
	case azureaksclusterv1.AzureAksClusterTransparentHugePage_THP_MADVISE:
		return "madvise"
	case azureaksclusterv1.AzureAksClusterTransparentHugePage_THP_NEVER:
		return "never"
	}
	return ""
}

func transparentHugePageDefragToArm(v azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag_DEFRAG_ALWAYS:
		return "always"
	case azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag_DEFRAG_DEFER:
		return "defer"
	case azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag_DEFRAG_DEFER_MADVISE:
		return "defer+madvise"
	case azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag_DEFRAG_MADVISE:
		return "madvise"
	case azureaksclusterv1.AzureAksClusterTransparentHugePageDefrag_DEFRAG_NEVER:
		return "never"
	}
	return ""
}

func undrainableNodeBehaviorToArm(v azureaksclusterv1.AzureAksClusterUndrainableNodeBehavior) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterUndrainableNodeBehavior_CORDON:
		return "Cordon"
	case azureaksclusterv1.AzureAksClusterUndrainableNodeBehavior_SCHEDULE:
		return "Schedule"
	}
	return ""
}

func hostPortProtocolToArm(v azureaksclusterv1.AzureAksClusterHostPortProtocol) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterHostPortProtocol_TCP:
		return "TCP"
	case azureaksclusterv1.AzureAksClusterHostPortProtocol_UDP:
		return "UDP"
	}
	return ""
}

func weekDayToArm(v azureaksclusterv1.AzureAksClusterWeekDay) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterWeekDay_SUNDAY:
		return "Sunday"
	case azureaksclusterv1.AzureAksClusterWeekDay_MONDAY:
		return "Monday"
	case azureaksclusterv1.AzureAksClusterWeekDay_TUESDAY:
		return "Tuesday"
	case azureaksclusterv1.AzureAksClusterWeekDay_WEDNESDAY:
		return "Wednesday"
	case azureaksclusterv1.AzureAksClusterWeekDay_THURSDAY:
		return "Thursday"
	case azureaksclusterv1.AzureAksClusterWeekDay_FRIDAY:
		return "Friday"
	case azureaksclusterv1.AzureAksClusterWeekDay_SATURDAY:
		return "Saturday"
	}
	return ""
}

func frequencyToArm(v azureaksclusterv1.AzureAksClusterMaintenanceFrequency) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterMaintenanceFrequency_DAILY:
		return "Daily"
	case azureaksclusterv1.AzureAksClusterMaintenanceFrequency_WEEKLY:
		return "Weekly"
	case azureaksclusterv1.AzureAksClusterMaintenanceFrequency_RELATIVE_MONTHLY:
		return "RelativeMonthly"
	case azureaksclusterv1.AzureAksClusterMaintenanceFrequency_ABSOLUTE_MONTHLY:
		return "AbsoluteMonthly"
	}
	return ""
}

func weekIndexToArm(v azureaksclusterv1.AzureAksClusterWeekIndex) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterWeekIndex_FIRST:
		return "First"
	case azureaksclusterv1.AzureAksClusterWeekIndex_SECOND:
		return "Second"
	case azureaksclusterv1.AzureAksClusterWeekIndex_THIRD:
		return "Third"
	case azureaksclusterv1.AzureAksClusterWeekIndex_FOURTH:
		return "Fourth"
	case azureaksclusterv1.AzureAksClusterWeekIndex_LAST:
		return "Last"
	}
	return ""
}

func nginxControllerToArm(v azureaksclusterv1.AzureAksClusterNginxDefaultController) string {
	switch v {
	case azureaksclusterv1.AzureAksClusterNginxDefaultController_ANNOTATION_CONTROLLED:
		return "AnnotationControlled"
	case azureaksclusterv1.AzureAksClusterNginxDefaultController_INTERNAL:
		return "Internal"
	case azureaksclusterv1.AzureAksClusterNginxDefaultController_EXTERNAL:
		return "External"
	case azureaksclusterv1.AzureAksClusterNginxDefaultController_NGINX_NONE:
		return "None"
	}
	return ""
}
