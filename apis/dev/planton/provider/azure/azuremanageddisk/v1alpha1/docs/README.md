# Azure Managed Disk Deployment: Storage That Outlives the Machine

## Introduction: The Disk Is Not Part of the VM

There is a mental model, inherited from physical servers, that a disk belongs to the machine it's plugged into. On that model, storage is a property of compute: you size disks when you size the VM, and when the VM goes away, its disks go with it.

Azure's own resource model says otherwise, and it matters. A **managed disk** is a first-class Azure Resource Manager resource with its own ARM ID, its own lifecycle, its own SKU, its own encryption configuration, and its own network posture. A virtual machine does not *contain* a data disk — it *attaches* one, at a logical unit number (LUN), with a caching mode. Detach the disk and the VM loses a mount point; the data is untouched. Delete the VM entirely and the disk sits in its resource group, waiting for the next machine.

This inversion is what makes stateful infrastructure survivable:

- **VM replacement is routine, data loss is not.** Rebuild a database VM from a fresh image, attach the same data disk, and the database comes back. The VM is cattle; the disk is the herd's memory.
- **A disk can outlive many machines.** Resize the VM, change its family, move it to a bigger SKU — the disk rides along, or waits.
- **One disk can serve several machines at once.** Azure's shared-disk capability (`maxShares`) lets multiple VMs attach the same disk simultaneously — the seam that Windows Server Failover Clustering, SQL Server FCI, and scale-out file systems build on. Only a standalone disk can express this; a disk defined inside one VM's spec structurally cannot.

This document explores how Azure Managed Disks actually work — the SKU spectrum, the create-option origin matrix, performance provisioning, encryption, and network posture — and how Planton's declarative API models each decision.

## The Seven SKUs: Choosing a Performance Model

The storage SKU (`storageAccountType`) is the disk's most fundamental choice, because it selects not just a speed but a *performance model*. There are two models in play.

### Fixed Per-Size Performance (the Classic Model)

For Standard HDD, Standard SSD, and classic Premium SSD, performance is a function of size. A 128 GiB Premium SSD is a P10 (500 IOPS, 100 MBps); a 1 TiB Premium SSD is a P30 (5,000 IOPS, 200 MBps). Need more IOPS? Buy a bigger disk, even if you don't need the space.

| SKU | Media | Redundancy | Best For |
|-----|-------|------------|----------|
| `STANDARD_LRS` | HDD | Local | Cold data, dev/test, infrequent access |
| `STANDARD_SSD_LRS` | SSD | Local | Light production, web servers, small databases |
| `STANDARD_SSD_ZRS` | SSD | Zone-redundant | Light workloads that must survive a zone outage |
| `PREMIUM_LRS` | Premium SSD | Local | The production default: OS disks, databases |
| `PREMIUM_ZRS` | Premium SSD | Zone-redundant | Premium performance that survives a zone outage |

The classic Premium tier has two escape hatches from the size-buys-performance coupling:

- **The `tier` field** lets a disk *buy a bigger size's performance tier*. A 256 GiB disk (natively P15) can be set to `tier: P40` and get a 2 TiB disk's IOPS and throughput without the capacity. This suits bursty workloads and pre-provisioned performance — but changing the tier on an attached disk briefly deallocates the VM.
- **Bursting.** Premium disks up to 512 GiB burst automatically on banked credits. Disks *larger* than 512 GiB can enable **on-demand bursting** (`onDemandBurstingEnabled`), which bursts beyond the provisioned tier whenever needed and bills per burst.

Both are classic-Premium-only concepts: they exist to work around the fixed-tier model, so the spec rejects them on any other SKU.

### Independently Dialed Performance (Premium SSD v2 and Ultra)

`PREMIUM_V2_LRS` and `ULTRA_SSD_LRS` abandon the tier ladder entirely. Capacity, IOPS, and throughput become three independent dials:

- `diskSizeGb` — pay for the capacity the data actually needs
- `diskIopsReadWrite` — provisioned IOPS, updated in place as needs change
- `diskMbpsReadWrite` — provisioned throughput in MBps, also updated in place

A 128 GiB Premium SSD v2 disk with 8,000 IOPS is a perfectly normal configuration — small disk, big performance, impossible on the classic tiers. Premium SSD v2's baseline (3,000 IOPS, 125 MBps) is free; provisioning beyond it bills per unit. Ultra Disk reaches the highest performance envelope Azure offers but demands more from its surroundings: zonal deployment and VM-family support.

The constraints worth knowing before choosing these SKUs:

- **Data disks only.** OS disks cannot use Premium SSD v2 or Ultra; they stay on the classic tiers.
- **Zonal only** (in regions with zones): plan the zone alignment with the attaching VM.
- **`logicalSectorSize`** (512 or 4096 bytes) exists only on these SKUs. Azure defaults to 4096; choose 512 only for legacy applications that require it — it's fixed at creation.

The spec gates all four performance dials and the sector size to exactly these two SKUs, so a `diskIopsReadWrite` on a classic Premium disk fails validation instead of failing the deploy.

### The ZRS Rule: Redundancy or Placement, Not Both

The ZRS SKUs (`STANDARD_SSD_ZRS`, `PREMIUM_ZRS`) synchronously replicate the disk across three availability zones — the data survives a zone outage. That property is incompatible with zone-pinning: a disk that lives in every zone cannot also be told to live in one. The spec enforces it: **a ZRS SKU must not set `zone`**.

Conversely, a zonal (LRS, zone-pinned) disk only attaches to VMs in the same zone. Zone alignment between disk and VM is a placement decision to make once, up front, because `zone` is fixed at creation.

## The Create Option: Every Disk Has an Origin Story

`createOption` answers the question "where does this disk's initial content come from?" — and it's fixed at creation, because a disk cannot retroactively change its origin. Each option requires specific source fields, and the spec enforces the pairings exactly as ARM does:

| Origin | Option | Source Fields | Typical Use |
|--------|--------|---------------|-------------|
| Nothing | `EMPTY` | `diskSizeGb` | New data volumes — the common case |
| A snapshot or disk | `COPY` | `sourceResourceId` | Restore, environment duplication, forensics |
| An image | `FROM_IMAGE` | exactly one of `imageReferenceId` / `galleryImageReferenceId` | Pre-building OS disks from platform or gallery images |
| A VHD blob | `IMPORT` | `sourceUri` + `storageAccountId` | Migrating existing VHDs into managed disks |
| A VHD blob, secured | `IMPORT_SECURE` | `sourceUri` + `storageAccountId`, `hyperVGeneration: V2` | Confidential-VM OS disk imports |
| A backup recovery point | `RESTORE` | `sourceResourceId` | Materializing Azure Backup recovery points |
| A future upload | `UPLOAD` | `uploadSizeBytes` | Streaming a VHD directly, no staging storage account |

A few of these deserve elaboration:

**COPY is the operational workhorse.** Cloning a snapshot produces a full, independent disk — writes to the clone never touch the source. Size is inherited when `diskSizeGb` is omitted (and surfaces in the `disk_size_gb` output, so downstream capacity planning reads the real value); the SKU is free to differ, so a Standard snapshot can restore onto Premium as a performance upgrade. Placement is *not* inherited — pin the clone's zone to match the VM that will attach it.

**UPLOAD is the odd one out.** It creates a disk in a state that *awaits* content: you get a writable SAS target and stream the VHD into it directly, skipping the staging storage account IMPORT requires. Its `uploadSizeBytes` must equal the source VHD's byte size *exactly*, footer included — which is why the field is in bytes, not GiB.

**IMPORT_SECURE is the confidential-VM path.** It requires UEFI boot (`hyperVGeneration: V2`) and is one of only two origins (with `FROM_IMAGE`) valid for customer-key confidential-VM encryption, covered below.

## Sizing: The One-Way Ratchet

`diskSizeGb` can only ever increase. There is no shrink operation on an Azure managed disk — reducing capacity means creating a smaller disk and migrating data at the filesystem level. This makes initial sizing a real decision: start realistic rather than generous, because every GiB is billed whether used or not, and growth is always available.

Growing has its own operational texture:

- Growing an **attached** disk may briefly detach it or deallocate the VM, *except* where Azure supports live resize (data disks under specific conditions).
- Crossing the 4 TiB boundary on classic SKUs always requires a detach.
- On Premium SSD v2 and Ultra, size is decoupled from performance anyway — you grow for capacity, not for IOPS, which removes the most common reason classic-tier disks got resized.

## Shared Disks: The Clustering Seam

Most disks attach to exactly one VM. Setting `maxShares` (2–10) turns a disk into a **shared disk**: several VMs attach it simultaneously, coordinating access at the application layer via SCSI persistent reservations — the mechanism Windows Server Failover Clustering, SQL Server FCI, and clustered file systems already speak.

This is a capability only a standalone disk resource can express, and it's the clearest illustration of why the disk is not modeled inside the VM: a disk that belongs to one VM's spec cannot belong to three VMs at once.

Two dedicated dials exist for the shared case, on Premium SSD v2 and Ultra: `diskIopsReadOnly` and `diskMbpsReadOnly` budget the aggregate performance available to VMs that mount the disk *read-only* (typically the standby nodes), separate from the read/write budget the active node consumes. The spec requires `maxShares` before accepting either — a read-only budget for a single-attach disk is meaningless.

The attach limit itself depends on SKU and size (premium SKUs and larger disks support more shares), so treat 10 as a ceiling, not a promise.

## Encryption: Three Layers, One Choice That Matters

Every managed disk is encrypted at rest, always — server-side encryption with platform-managed keys is Azure's default and requires nothing from you. The decisions begin when compliance requires *your* keys or *confidential computing*:

**Customer-managed keys (`diskEncryptionSetId`).** A disk encryption set is an Azure resource that wraps a Key Vault key and grants the disk service access to it. Point the disk at the set's ARM ID and server-side encryption uses your key — revocable, rotatable, auditable. This field takes a plain ARM ID: disk encryption sets are not modeled as a Planton kind yet, so the reference crosses the model boundary as a literal.

**Confidential-VM profiles (`securityType`).** For OS disks of confidential VMs, the security type selects what gets encrypted with what: guest state only with a platform key, full disk with a platform key, or full disk with a *customer* key. The customer-key variant is the most constrained: it requires its own dedicated encryption set (`secureVmDiskEncryptionSetId` — required with that security type and only valid then), and the disk's origin must be `FROM_IMAGE` or `IMPORT_SECURE`, because a confidential OS disk must be born from a trusted source.

**Trusted launch (`trustedLaunchEnabled`).** Secure boot and vTPM support for regular (non-confidential) VMs — a middle security tier. A disk is trusted-launch *or* confidential, never both, and trusted launch requires an image-based origin (`FROM_IMAGE` or `IMPORT`).

The spec enforces every pairing in this section — the two encryption sets are mutually exclusive, the secure set pairs exactly with the customer-key security type, and trusted launch conflicts with any security type — so an inconsistent security posture fails validation, not deployment.

One deliberate omission: the legacy `encryption_settings` block (Azure Disk Encryption — the in-guest BitLocker/dm-crypt path) is not modeled. Server-side encryption with disk encryption sets is the modern grain: it encrypts below the OS, works with every OS image, composes with confidential computing, and is what Azure steers new deployments toward. Workloads still on the in-guest path should migrate rather than expect the API to carry the legacy surface forward.

## Network Posture: Yes, Disks Have One

It surprises people that a *disk* has network configuration. The reason: managed disks support **export** — generating a SAS URL to download the disk's content, which is how VHDs leave Azure and how some backup tooling operates. That endpoint is an attack surface, and Azure gives it a policy:

- **`ALLOW_ALL`** (Azure's default): the export endpoint is reachable over the network with proper authorization.
- **`ALLOW_PRIVATE`**: export only flows through the private endpoints of a **disk-access resource** (`diskAccessId` — required with this policy, and only valid then). The disk-access resource is referenced by plain ARM ID; it is not modeled as a Planton kind.
- **`DENY_ALL`**: network export disabled entirely. The lockdown posture — right for any disk that will never need SAS-based export, which is most production data disks.

`publicNetworkAccessEnabled` (default true) is the coarser companion switch: set it false alongside `ALLOW_PRIVATE` for a fully private posture, or alongside `DENY_ALL` for belt-and-suspenders.

VM disk *traffic* — the reads and writes the attached VM performs — never touches any of this. It flows over Azure's storage fabric regardless of the export posture.

## What Changes Safely, and What Replaces the Disk

For a resource whose entire purpose is holding data, knowing which changes are destructive is not optional:

**Identity — changing any of these replaces the disk and its data:**
name, region, zone, create option and its source fields, logical sector size, security profile (security type, trusted launch, secure encryption set), performance-plus, edge zone.

**Update in place, no interruption:**
size (increase only, with the caveats above), the PremiumV2/Ultra performance dials, network posture fields, `publicNetworkAccessEnabled`, tags.

**Update in place, with a VM interruption:**
SKU changes and `tier` changes on an attached disk deallocate the VM for the change and restart it after.

The practical rule: decide the identity fields once, at design time. Everything you'll want to tune over the disk's life — size, performance, network posture — is deliberately in the updatable set.

## The Planton Approach

Planton provides a declarative, protobuf-based API for Azure Managed Disks. Three design decisions define its shape.

### The Disk Knows Nothing About Its Consumers

The most important structural choice: **the attachment lives on the VM, not the disk.** An `AzureVirtualMachine`'s `dataDiskAttachments` lists each disk it mounts — referencing this disk's `disk_id` output — along with the LUN and caching mode, which are properties of *that VM's use of the disk*, not of the disk itself.

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedDisk
metadata:
  name: orders-db-data
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: orders-db-data
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 512
  zone: "1"
---
apiVersion: azure.planton.dev/v1alpha1
kind: AzureVirtualMachine
metadata:
  name: orders-db
spec:
  # ... VM configuration ...
  dataDiskAttachments:
    - managedDiskId:
        valueFrom:
          name: orders-db-data
      lun: 0
      caching: READ_ONLY
```

This mirrors Azure's own model and buys three things: replacing the VM is a pure VM-side change that never touches the disk resource; a shared disk can appear in several VMs' attachment lists without the disk model contorting; and the disk's spec stays a pure description of storage. Name the disk after the data it carries (`orders-db-data`), not the VM it happens to attach to — the VM is the transient party in this relationship.

### Validation That Mirrors ARM

The spec encodes the full conditional-validation matrix Azure Resource Manager enforces at deploy time — every create-option source pairing, every SKU gate on the performance dials, the shared-disk requirement for the read-only budgets, the encryption exclusions and pairings, the network-policy pairing with `diskAccessId`, and the ZRS/zone conflict. A manifest that would fail in ARM fails at validation instead, with a message naming the actual rule, before any deployment starts.

### References Where Kinds Exist, ARM IDs Where They Don't

The resource group is a first-class reference (an `AzureResourceGroup`'s output, or a literal). Disk encryption sets and disk-access resources are plain ARM IDs, because those kinds are not modeled yet — a DiskEncryptionSet kind is under evaluation alongside a Key Vault key kind, since the two compose (a set wraps a key) and modeling one without the other would leave the reference chain dangling. The API takes the honest position: a typed string field today, an upgradeable reference tomorrow.

### Example: Dialed-Performance Database Volume

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedDisk
metadata:
  name: ledger-data
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: ledger-data
  storageAccountType: PREMIUM_V2_LRS
  createOption: EMPTY
  diskSizeGb: 128
  zone: "1"
  diskIopsReadWrite: 8000
  diskMbpsReadWrite: 300
  tags:
    workload: ledger
    cost_center: payments
```

128 GiB of capacity with 8,000 IOPS — sized for the data, dialed for the workload, both dials adjustable in place as the workload evolves.

### Example: Locked-Down Regulated Data

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureManagedDisk
metadata:
  name: pii-vault-data
spec:
  region: eastus
  resourceGroup:
    valueFrom:
      name: prod-rg
  name: pii-vault-data
  storageAccountType: PREMIUM_LRS
  createOption: EMPTY
  diskSizeGb: 1024
  zone: "2"
  diskEncryptionSetId:
    value: /subscriptions/xxx/resourceGroups/security-rg/providers/Microsoft.Compute/diskEncryptionSets/prod-cmk
  networkAccessPolicy: DENY_ALL
  publicNetworkAccessEnabled: false
  tags:
    data-classification: restricted
```

Customer-managed keys, no network export, no public reachability — the full lockdown posture in four fields.

## Common Anti-Patterns to Avoid

**❌ Anti-Pattern 1: Naming Disks After VMs**

`web-vm-01-disk-2` tells you nothing when `web-vm-01` is replaced by `web-vm-07` and the disk lives on.

**✅ Solution:** Name the disk after the data it carries — `orders-db-data`, `ledger-wal`. The name is identity (changing it replaces the disk), so make it describe the durable thing.

---

**❌ Anti-Pattern 2: Oversizing Classic Premium for IOPS**

Buying a 2 TiB P40 to get 7,500 IOPS for 200 GiB of data — paying for 1.8 TiB of empty space as a performance tax.

**✅ Solution:** Use `PREMIUM_V2_LRS` and dial IOPS independently, or set `tier` on a right-sized classic Premium disk.

---

**❌ Anti-Pattern 3: Generous Initial Sizing "To Be Safe"**

Size can only increase. A 4 TiB disk provisioned for a 300 GiB dataset bills for 4 TiB forever — there is no path back.

**✅ Solution:** Start realistic. Growth is always available; shrinkage never is.

---

**❌ Anti-Pattern 4: Ignoring Zone Alignment**

A zonal disk in zone 1 cannot attach to a VM in zone 2. Discovering this at attach time means recreating the disk (zone is fixed at creation) or the VM.

**✅ Solution:** Decide zone placement for the disk and its VM together, up front — or use a ZRS SKU (without a zone) when the data must survive a zone outage.

---

**❌ Anti-Pattern 5: Leaving Export Open on Sensitive Data Disks**

Azure's default posture (`ALLOW_ALL`) leaves the SAS-export endpoint reachable. Most production data disks never need export at all.

**✅ Solution:** Set `networkAccessPolicy: DENY_ALL` (and `publicNetworkAccessEnabled: false`) on disks that never export; use `ALLOW_PRIVATE` with a disk-access resource when export must exist but privately.

## Conclusion: Storage as a First-Class Decision

Treating the disk as an appendage of the VM is how data ends up coupled to the machine that happens to hold it — and how VM replacement becomes a data-migration project. Azure's model, and this API's, inverts that: the disk is the durable party, the VM is the transient one, and the attachment between them is the VM's declaration, not the disk's.

Get four decisions right at creation — SKU, origin, size, and placement — and everything else about the disk stays adjustable for its whole life. The spec's validation matrix ensures the decisions you write down are ones Azure will actually accept, and the `disk_id` output is the single join point everything downstream builds on.

Model the data first. Attach the machines to it.
