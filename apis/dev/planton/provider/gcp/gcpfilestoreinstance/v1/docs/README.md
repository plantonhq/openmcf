# GcpFilestoreInstance — Deep Dive

## The problem this resource solves

Filestore fills the gap between block storage (Persistent Disks, which attach to one VM) and object storage (Cloud Storage, which speaks HTTP): a managed POSIX filesystem shared across many clients over NFS. This kind models the instance as a first-class node so the network attachment, the encryption key, the replication peer, and the target project are all explicit references to reviewable resources — and so the dangerous decisions (deletion protection, immutable replacements, capacity growth) are deliberate spec choices rather than tool defaults.

## Where it sits in the composition

Outbound references the instance makes: `projectId` → GcpProject, `networkConfig.network` → GcpVpcNetwork, `kmsKeyName` → GcpKmsKey, `initialReplication.peerInstances[]` → other GcpFilestoreInstance resources (via `status.outputs.instance_id`).

Inbound, consumers use `ip_addresses` and `file_share_name` to construct the NFS mount path: GKE workloads through the Filestore CSI driver (ReadWriteMany PersistentVolumes), Compute Engine VMs through a plain `mount` command. The instance sits above the network foundation and below the workloads that mount it.

## Lifecycle contract

| Property | Behavior |
|---|---|
| `instanceName`, `location`, `tier`, `protocol`, `networkConfig` (all of it), `kmsKeyName`, `initialReplication` | Immutable (ForceNew) — replacement destroys the instance and its data |
| `fileShare.name` | Immutable — the export path is fixed at creation |
| `fileShare.capacityGb` | Grows in place, never shrinks |
| `fileShare.sourceBackup`, `tags` | Create-time only |
| `fileShare.nfsExportOptions`, `description`, `labels`, deletion protection, `performanceConfig` | Mutable in place |
| Deletion | Blocked while `deletionProtectionEnabled: true` — the flag must be flipped false first |

`instanceName` is optional: when empty, `metadata.name` becomes the cloud-side name. Both IaC engines derive the identical name via the same explicit conditional.

## One file share, one network

GCP Filestore supports exactly one file share and one VPC attachment per instance. The spec models both as singular sub-messages rather than repeated fields — a repeated field would promise a flexibility the API does not have.

## Tiers

| Tier | Storage | Availability | Min capacity | IOPS tuning | Location | NFSv4.1 |
|------|---------|--------------|--------------|-------------|----------|---------|
| STANDARD / BASIC_HDD | HDD | Single-zone | 1 TiB | No | Zone | No |
| PREMIUM / BASIC_SSD | SSD | Single-zone | 2.5 TiB | No | Zone | No |
| HIGH_SCALE_SSD | SSD | Single-zone | 10 TiB | No | Zone | Yes |
| ZONAL | SSD | Single-zone | 1 TiB | Yes | Zone | Yes |
| REGIONAL | SSD | Multi-zone | 1 TiB | Yes | Region | Yes |
| ENTERPRISE | SSD | Multi-zone | 1 TiB | Yes | Region | Yes |

STANDARD/PREMIUM are the rebranded names for BASIC_HDD/BASIC_SSD; the API accepts both. ZONAL is the modern replacement for BASIC_SSD and HIGH_SCALE_SSD. `location` is a zone for the zonal tiers and a region for REGIONAL/ENTERPRISE — one unified field, because the underlying `zone` attribute is deprecated in both providers.

## Connectivity: the three connect modes

- **DIRECT_PEERING** (default) — VPC peering. The simplest setup; right for a standalone VPC.
- **PRIVATE_SERVICE_ACCESS** — rides an existing service-networking connection on the VPC. Required for Shared VPC consumers and common in enterprise network topologies. The connection must already exist; this module does not create it.
- **PRIVATE_SERVICE_CONNECT** — PSC endpoints.

`reservedIpRange` pins the instance's `/29` block when address planning matters; left empty, GCP picks an unused range — and either way the resolved block surfaces in the `reserved_ip_range` output. `modes` selects IP versions; empty means `["MODE_IPV4"]`, the standard NFS posture, and both engines send that default explicitly so the realized instance is identical.

## Replication (`initialReplication`)

ENTERPRISE and REGIONAL instances support cross-instance replication, fixed at create time. The common DR posture: the ACTIVE (source) instance already exists, and this spec creates the STANDBY replica pointing at it — `role` empty or `STANDBY`, with the peer referenced as a `GcpFilestoreInstance` (its `status.outputs.instance_id` is exactly the full resource path the API wants) or as a literal path. Backups cannot be taken from a STANDBY replica, and the peer relationship cannot be changed later without replacing the instance.

## The safety posture

- **`deletionProtectionEnabled` is the destroy guard.** A protected instance cannot be deleted — by Planton or anyone else — until the flag is flipped false in a deliberate, reviewable change. Production presets set it.
- **Capacity is a ratchet.** Growth applies in place; shrinking is impossible. Start at the tier minimum and grow on demand.
- **The API enablement never disables on destroy** — tearing down one instance must not break every other Filestore user in the project.

## Labels and tags

User `labels` merge beneath the platform attribution labels (`planton-ai_*`), which win on key conflicts — identically in both engines. Resource Manager `tags` are a different mechanism entirely: org-policy and IAM-condition bindings in the form `tagKeys/{id}` → `tagValues/{id}`, applied at create time only.

## Deliberately not modeled

All of the following were verified absent from the released `google` `~> 6.x` provider line (schema-probe verified); revisit when the floor moves:

- **`psc_config.endpoint_project`** — PSC endpoint project selection.
- **`nfs_export_options.network`** — per-export source network for PSC.
- **LDAP `directory_services`** — NFSv4.1 LDAP integration.
- **`source_backupdr_backup`** — restore from a Backup and DR Service backup (`fileShare.sourceBackup` covers Filestore-native backups).
- **`deletion_policy`** — a client-side flag, not a server field; destroy semantics ride on `deletionProtectionEnabled` identically on both engines (the Pulumi module pins the bridged provider's extra flag to `DELETE` for parity).

Also not modeled as kinds: **`GcpFilestoreBackup`** and **`GcpFilestoreSnapshot`**. Both are one-shot point-in-time resources — a poor fit for continuously-reconciled IaC, where every re-apply would either do nothing or silently redefine what "the backup" means — and GCP offers no schedule resource for Filestore that would make them declarative. Revisit on concrete pull.
