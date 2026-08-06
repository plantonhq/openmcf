# AWS FSx for ONTAP: Architecture Reference

This document provides a deep technical reference for Amazon FSx for NetApp ONTAP as deployed via the AwsFsxOntapFileSystem API. It covers the file system hierarchy, deployment types, HA pairs, storage options, networking, endpoints, backup strategy, encryption, and cost optimization.

File system creation is a long-running operation — the Terraform provider's default create timeout is 60 minutes; budget deploy windows accordingly.

---

## 1. Architecture: File System → SVMs → Volumes

### Hierarchy

FSx for ONTAP follows a three-tier model:

```
File System (this component)
  └── Storage Virtual Machine (SVM) — AwsFsxOntapStorageVirtualMachine
        └── Volume — AwsFsxOntapVolume
```

- **File system:** The top-level ONTAP cluster. Defines deployment type, capacity, throughput, networking, and encryption. One file system can host multiple SVMs.
- **Storage Virtual Machine (SVM):** A logical storage server. Exposes NFS, SMB, and/or iSCSI protocols. Each SVM has its own namespace, LIFs (Logical Interfaces), and volumes.
- **Volume:** A logical storage container within an SVM. Volumes are what clients mount (NFS export path, SMB share, or iSCSI LUN).

**Key design note:** This Planton component manages only the file system. SVMs and volumes are separate resources with their own lifecycle. Create them after the file system is provisioned, referencing `status.outputs.file_system_id`.

### Data Flow

```
Client (EC2, ECS, VMware)
  → NFS/SMB/iSCSI
  → LIF (Logical Interface) on SVM
  → Volume
  → Aggregate (ONTAP storage pool)
  → SSD/HDD storage
```

---

## 2. Deployment Types

### SINGLE_AZ_1

- **Availability:** Single AZ. No cross-AZ failover.
- **HA pairs:** Fixed at 1. Capacity caps at 192 TiB.
- **Subnets:** Exactly one.
- **Scaling caveat:** Increasing per-pair throughput replaces the file system.
- **Status:** First generation. Prefer SINGLE_AZ_2 for new deployments.

### SINGLE_AZ_2

- **Availability:** Single AZ. Current generation.
- **HA pairs:** 1–12, and the ONLY deployment type supporting scale-out. HA
  pairs can be increased **without replacement**.
- **Subnets:** Exactly one.
- **Use case:** Most workloads. Recommended default.

### MULTI_AZ_1

- **Availability:** Multi-AZ with automatic failover across two AZs.
- **HA pairs:** Fixed at 1 (standby in second AZ). Capacity caps at 192 TiB.
- **Subnets:** Exactly two, in different AZs.
- **Requirements:** `preferred_subnet_id` (endpoint IP range and route tables
  optional — AWS picks defaults).
- **Scaling caveat:** Increasing per-pair throughput replaces the file system.
- **Status:** First generation. Prefer MULTI_AZ_2 for new multi-AZ deployments.

### MULTI_AZ_2

- **Availability:** Multi-AZ with automatic failover. Current generation.
- **HA pairs:** Fixed at 1.
- **Subnets:** Exactly two, in different AZs.
- **Requirements:** `preferred_subnet_id` (endpoint IP range and route tables
  optional — AWS picks defaults).
- **Use case:** Mission-critical workloads requiring high availability.

### Decision Matrix

| Criteria | SINGLE_AZ_1 | SINGLE_AZ_2 | MULTI_AZ_1 | MULTI_AZ_2 |
|----------|-------------|-------------|-----------|-----------|
| Generation | First | Current | First | Current |
| Cross-AZ failover | No | No | Yes | Yes |
| HA pairs | 1 | 1–12 (in-place add) | 1 | 1 |
| Max capacity | 192 TiB | 512 TiB × pairs (1 PiB max) | 192 TiB | 512 TiB |
| Subnets | 1 | 1 | 2 | 2 |
| Per-pair throughput tiers (MB/s) | 128–4096 | 384–6144 (1536+ for scale-out) | 128–4096 | 384–6144 |

---

## 3. HA Pairs Explained

### What Is an HA Pair?

An HA (High Availability) pair is a pair of file servers (nodes) that provide redundancy within a single AZ. Each pair contributes independent throughput and IOPS capacity.

### Scale-Out (SINGLE_AZ_2 Only)

- **SINGLE_AZ_2:** 1–12 HA pairs, added in place. Every other deployment type
  is fixed at 1 pair.
- **Total throughput** = `throughput_capacity_per_ha_pair` × `ha_pairs`.
- **Example:** 4 HA pairs × 1536 MB/s = 6144 MB/s aggregate throughput.
- **Capacity scales too:** each pair carries 1024 GiB – 512 TiB, up to 1 PiB
  total at 12 pairs.

### Multi-AZ Fixed at 1 HA Pair

- **MULTI_AZ_1 / MULTI_AZ_2:** Always 1 HA pair.
- The pair spans two AZs: active node in preferred subnet, standby in the other.
- Failover is automatic; no scale-out via HA pairs.

### Sizing Throughput (exactly one arm)

Throughput is set through exactly one of two fields:

- **`throughput_capacity`** — whole-file-system sizing, the first-generation
  arm. Valid values: 128, 256, 512, 1024, 2048, 4096 MB/s.
- **`throughput_capacity_per_ha_pair`** — per-pair sizing. Valid values by
  deployment type:
  - SINGLE_AZ_1 / MULTI_AZ_1: 128, 256, 512, 1024, 2048, 4096
  - SINGLE_AZ_2 / MULTI_AZ_2 with 1 HA pair: 384, 768, 1536, 3072, 6144
  - SINGLE_AZ_2 with multiple HA pairs: 1536, 3072, 6144 per pair

Higher tiers cost more. Choose based on workload I/O profile. Scales in place
on SINGLE_AZ_2 and MULTI_AZ_2; on the first generation an increase replaces
the file system.

---

## 4. Storage

### SSD Primary Storage (the only option)

- **Latency:** Sub-millisecond.
- **Capacity range:** 1024 GiB per HA pair to 512 TiB per pair (192 TiB cap on
  the first generation; 1 PiB total at 12 SINGLE_AZ_2 pairs).
- **IOPS:** AUTOMATIC (3 IOPS/GiB) or USER_PROVISIONED (up to 2,400,000).

ONTAP file systems are SSD-only at the file-system level — the HDD and
Intelligent-Tiering storage classes belong to other FSx types (Windows/Lustre
and OpenZFS/Lustre respectively), and AWS rejects them for ONTAP at create
time. ONTAP's cost tiering happens INSIDE the file system instead: every
volume has a `tiering_policy` that moves cold data to the built-in elastic
capacity pool, which is why a file-system-level HDD option does not exist.

### Data Reduction

ONTAP provides built-in compression and deduplication. Typical effective capacity is 2–5× the provisioned capacity for many workloads.

---

## 5. Networking

### Single-AZ

- **Subnets:** Exactly one subnet.
- **ENIs:** One per HA pair (e.g., 4 HA pairs = 4 ENIs).
- **Placement:** All nodes in the same AZ.

### Multi-AZ

- **Subnets:** Exactly two subnets in different AZs.
- **ENIs:** Two (one per AZ).
- **Endpoint IP address range:** Floating-IP CIDR for the management,
  intercluster, and SVM data endpoints. Must NOT overlap with any subnet in
  the VPC — AWS recommends a range outside the VPC CIDR entirely (the
  198.19.0.0/16 block); clients reach it through the managed route tables.
  Omit to let AWS pick an unused range.
- **Route tables:** Optional. `route_table_ids` lists every route table
  associated with client subnets; AWS creates and repoints routes to the
  floating IPs on failover. Omitted → the VPC main route table.

### Security Groups

Must allow:

| Port | Protocol | Purpose |
|------|----------|---------|
| 111 | TCP | Portmapper (NFS) |
| 635 | TCP | mountd (NFS) |
| 2049 | TCP | NFS |
| 4045-4046 | TCP | NFS lock/status |
| 445 | TCP | SMB |
| 3260 | TCP | iSCSI |
| 443 | TCP | ONTAP REST API |

---

## 6. Endpoints

### Management Endpoint

- **Purpose:** ONTAP CLI (SSH) and REST API access.
- **Outputs:** `management_dns_name`, `management_ip_addresses`.
- **Access:** `ssh fsxadmin@<management_dns_name>` (requires `fsx_admin_password` in spec).
- **Use cases:** LIF management, SnapMirror configuration, aggregate monitoring, advanced administration.

### Intercluster Endpoint

- **Purpose:** NetApp SnapMirror replication between FSx for ONTAP file systems (same or cross-region).
- **Outputs:** `intercluster_dns_name`, `intercluster_ip_addresses`.
- **Use case:** Hybrid cloud replication to/from on-premises NetApp, or disaster recovery between AWS regions.

### Data Endpoints

Data access (NFS, SMB, iSCSI) is provided by **SVMs**, not the file system directly. Create an SVM and volumes to expose data endpoints.

---

## 7. ONTAP CLI and REST API Access

### Enabling Access

Set `fsx_admin_password` in the spec (8–50 characters). This enables:

- **SSH:** `ssh fsxadmin@<management_dns_name>`
- **REST API:** `https://<management_dns_name>/api`

### Security

- The password is sensitive and is not returned in read operations.
- Omit `fsx_admin_password` if ONTAP CLI access is not needed (reduces attack surface).

### Common CLI Operations

- **System health:** `system health status show`
- **Aggregate info:** `storage aggregate show`
- **SVM list:** `vserver show`
- **Volume list:** `volume show`
- **SnapMirror:** Configure via CLI or REST API for replication.

---

## 8. Backup Strategy

### FSx Automatic Backups

- **Retention:** `automatic_backup_retention_days` (0–90). Set 0 to disable.
- **Window:** `daily_automatic_backup_start_time` in `HH:MM` UTC.
- **Mechanics:** Incremental. Stored in AWS-managed location. Restore creates a new file system.

### ONTAP Snapshots

- **Independent:** ONTAP's built-in snapshots are separate from FSx backups.
- **Configured on volumes:** Create snapshots via SVM/volume configuration or ONTAP CLI.
- **Use case:** Point-in-time recovery, cloning, SnapMirror source.

### Recommendation

- **Development:** Disable FSx backups (`automatic_backup_retention_days: 0`). Use ONTAP snapshots if needed.
- **Production:** Enable FSx backups (7–30 days) plus ONTAP snapshots for layered protection.

---

## 9. Encryption at Rest

- **Default:** All FSx for ONTAP file systems are encrypted at rest with an AWS-managed key.
- **Customer-managed key:** Set `kms_key_id` to use a customer-managed KMS key. **ForceNew** — cannot be changed after creation.
- **Key requirements:** The KMS key must be in the same region as the file system and allow FSx to use it.

---

## 10. Cost Optimization Tips

1. **Right-size throughput:** Start with the lowest tier that meets your needs. Throughput can be scaled up after creation for SINGLE_AZ_2 and MULTI_AZ_2.
2. **Tier cold data per volume:** Set each volume's `tiering_policy` (AUTO or SNAPSHOT_ONLY) to move cold data to the elastic capacity pool — ONTAP's equivalent of an HDD tier, without provisioning it.
3. **Leverage data reduction:** ONTAP compression and deduplication reduce effective storage cost (2–5× typical).
4. **Manage backup retention:** Set `automatic_backup_retention_days` to the minimum required by recovery objectives.
5. **Single-AZ for non-critical workloads:** MULTI_AZ adds cost for the standby server and cross-AZ replication.
6. **Scale-out vs. scale-up:** For single-AZ, adding HA pairs increases throughput. Ensure workload can utilize the additional capacity before scaling.
7. **AUTOMATIC IOPS:** Use `disk_iops_configuration.mode: AUTOMATIC` unless you have a demonstrated need for USER_PROVISIONED IOPS.

---

## 11. Summary

| Topic | Key Takeaway |
|-------|--------------|
| **Hierarchy** | File system → SVM → Volume. This component manages file system only. |
| **Deployment types** | SINGLE_AZ_2 for most; MULTI_AZ_2 for HA. |
| **HA pairs** | SINGLE_AZ_2: 1–12 (in-place scale-out). All others: fixed at 1. |
| **Throughput** | Exactly one arm: whole-system `throughput_capacity` or per-pair `throughput_capacity_per_ha_pair` (gen2 tiers 384–6144 MB/s). |
| **Storage** | SSD only, 1024 GiB per pair to 1 PiB total. Cost tiering is per-volume (tiering policies). |
| **Networking** | Single-AZ: 1 subnet. Multi-AZ: 2 subnets + preferred subnet (+ optional floating-IP range and route tables). |
| **Endpoints** | Management (CLI/API), intercluster (SnapMirror). Data via SVMs. |
| **Backups** | FSx backups (0–90 days) and ONTAP snapshots; final-backup decisions live on volumes. |
| **Encryption** | Always on. Optional customer-managed KMS. |

For API reference and examples, see the parent [README.md](../README.md) and [catalog-page.md](../catalog-page.md).
