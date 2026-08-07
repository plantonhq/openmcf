# Azure Premium

This preset creates the premium SSD class for AKS: Premium SSD managed disks (locally redundant) through the Azure Disk CSI driver, provisioned in the zone the consuming pod schedules into, and expandable after creation. `Premium_LRS` is the production SSD tier — guaranteed IOPS and throughput per disk size, over the burstable `StandardSSD_LRS` tier AKS defaults to.

## When to Use

- Production databases and IO-bound workloads on AKS that need guaranteed disk performance
- Pinning a workload's storage tier explicitly instead of inheriting the cluster default
- Note: Premium SSD disks require a VM series that supports premium storage (most current AKS node sizes do — the `s`-suffixed series)

## Key Configuration Choices

- **`provisioner: disk.csi.azure.com`** — the Azure Disk CSI driver, built into current AKS clusters. IMMUTABLE after creation
- **`skuName: Premium_LRS`** — the parameters vocabulary belongs to the Azure Disk CSI driver; other documented values include `StandardSSD_LRS`, `Premium_ZRS` (zone-redundant), `PremiumV2_LRS`, and `UltraSSD_LRS`. Azure managed disks are encrypted at rest by default; a customer-managed key set can be named with the `diskEncryptionSetID` parameter
- **`wait_for_first_consumer`** — managed disks are zonal on zone-enabled clusters; binding waits for the consuming pod so the disk provisions in that pod's zone. Claims of this class stay **Pending until a pod uses them — correct behavior, not an error**
- **`allow_volume_expansion: true`** — claims can grow later (never shrink); the Azure Disk CSI driver supports expansion
- **`is_default_class: false`** — AKS ships its own default class (`default`, StandardSSD); promoting this preset to default would create a dueling pair unless the built-in default is demoted first

## Placeholders to Replace

This preset has no placeholders — it is deployable as-is on any AKS cluster whose nodes support premium storage.

## Related Presets

- **01-aws-ebs-gp3** — the equivalent tier on EKS, and the preset whose discussion covers promoting a class to cluster default
- **02-gcp-pd-ssd** — the equivalent tier on GKE
