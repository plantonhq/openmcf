# AWS EBS gp3

This preset creates the standard general-purpose SSD class for EKS: encrypted gp3 volumes through the AWS EBS CSI driver, provisioned in the zone the consuming pod schedules into, and expandable after creation. gp3 is the current-generation general-purpose EBS type — it delivers a 3,000 IOPS / 125 MiB/s baseline at any size (unlike gp2, whose IOPS scale with size), and IOPS and throughput can be raised independently of capacity.

## When to Use

- The workhorse storage tier on EKS — databases, queues, anything that wants an SSD-backed ReadWriteOnce volume
- Replacing the cluster's built-in gp2 default with a gp3 tier (better baseline performance at lower cost per GB)
- As the candidate for the cluster's default class: this is the preset to promote with `is_default_class: true` once the existing default (e.g. EKS's built-in `gp2`) has been demoted. The preset ships with `is_default_class: false` so applying it never creates a second, dueling default

## Key Configuration Choices

- **`provisioner: ebs.csi.aws.com`** — the AWS EBS CSI driver; it must be installed in the cluster (an EKS add-on) for claims of this class to provision. IMMUTABLE after creation
- **`type: gp3`** — the parameters vocabulary belongs to the EBS CSI driver; other documented keys include `iops` (e.g. `"6000"`), `throughput` (MiB/s, e.g. `"250"`), and `kmsKeyId` for a customer-managed encryption key
- **`encrypted: "true"`** — encryption at rest for every volume of this class; costs nothing and removes a whole compliance conversation
- **`wait_for_first_consumer`** — EBS volumes are zonal; binding waits for the consuming pod so the volume provisions in that pod's zone. Claims of this class stay **Pending until a pod uses them — correct behavior, not an error**
- **`allow_volume_expansion: true`** — claims can grow later (never shrink); the EBS CSI driver supports online expansion

## Placeholders to Replace

This preset has no placeholders — it is deployable as-is on any EKS cluster with the EBS CSI driver installed. Add `iops`/`throughput` parameters to raise performance beyond the gp3 baseline.

## Related Presets

- **02-gcp-pd-ssd** — the equivalent tier on GKE
- **03-azure-premium** — the equivalent tier on AKS
