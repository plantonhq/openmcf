# GCP PD SSD

This preset creates the performance SSD class for GKE: SSD persistent disks through the GCE Persistent Disk CSI driver, provisioned in the zone the consuming pod schedules into, and expandable after creation. `pd-ssd` is the consistent-latency SSD tier — the choice for databases and latency-sensitive workloads over the default `pd-balanced`/`pd-standard` tiers GKE ships as built-in classes.

## When to Use

- Databases and latency-sensitive workloads on GKE that need SSD IOPS beyond the built-in `standard-rwo` (pd-balanced) class
- Pinning a workload's storage tier explicitly instead of inheriting the cluster default
- Any GKE claim where consistent single-digit-millisecond latency matters

## Key Configuration Choices

- **`provisioner: pd.csi.storage.gke.io`** — the GCE PD CSI driver, enabled by default on current GKE clusters. IMMUTABLE after creation
- **`type: pd-ssd`** — the parameters vocabulary belongs to the PD CSI driver; other documented values include `pd-balanced`, `pd-standard`, and `pd-extreme`. GKE persistent disks are encrypted at rest by default; a customer-managed key can be named with the `disk-encryption-kms-key` parameter
- **`wait_for_first_consumer`** — persistent disks are zonal; binding waits for the consuming pod so the disk provisions in that pod's zone. Claims of this class stay **Pending until a pod uses them — correct behavior, not an error**
- **`allow_volume_expansion: true`** — claims can grow later (never shrink); the PD CSI driver supports online expansion
- **`is_default_class: false`** — GKE ships its own default class (`standard-rwo`); promoting this preset to default would create a dueling pair unless the built-in default is demoted first

## Placeholders to Replace

This preset has no placeholders — it is deployable as-is on any GKE cluster with the PD CSI driver enabled (the default on current versions).

## Related Presets

- **01-aws-ebs-gp3** — the equivalent tier on EKS, and the preset whose discussion covers promoting a class to cluster default
- **03-azure-premium** — the equivalent tier on AKS
