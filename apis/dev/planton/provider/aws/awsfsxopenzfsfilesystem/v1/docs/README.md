# AwsFsxOpenzfsFileSystem — Technical Reference

## Service Overview

Amazon FSx for OpenZFS is a fully managed file storage service built on the OpenZFS file system. It provides sub-millisecond latency, up to 10 GB/s throughput, and over 1 million IOPS through standard NFS protocols (NFSv3, NFSv4.0, NFSv4.1, NFSv4.2). Key OpenZFS features include instant snapshots, data cloning, ZSTD/LZ4 compression, and per-user/group quotas.

## Terraform Resource Mapping

| Planton Concept | Terraform Resource | Notes |
|-----------------|-------------------|-------|
| File system | `aws_fsx_openzfs_file_system` | Primary resource |
| Root volume | Inline `root_volume_configuration` block | Configured on the file system resource |
| Child volumes | `aws_fsx_openzfs_volume` | NOT managed by this component (independent lifecycle) |
| Snapshots | `aws_fsx_openzfs_snapshot` | NOT managed by this component |

## Deployment Types

### SINGLE_AZ_1

First-generation single-AZ deployment. Lower throughput ceiling (max 4,096 MB/s). Uses one subnet and one ENI. Lowest cost option. Best for development and non-critical workloads.

### SINGLE_AZ_2

Latest single-AZ deployment (recommended for new workloads). Higher throughput ceiling (max 10,240 MB/s). More features than SINGLE_AZ_1. One subnet, one ENI.

### SINGLE_AZ_HA_1 / SINGLE_AZ_HA_2

HA variants of the two single-AZ generations: an active/standby file-server pair inside one availability zone. Automatic failover without cross-AZ data transfer charges — the middle ground between plain single-AZ and MULTI_AZ_1. One subnet; throughput follows the parent generation's value set (validated by AWS at create time).

### MULTI_AZ_1

Multi-AZ deployment with automatic failover. Data is synchronously replicated across two AZs. Requires two subnets in different AZs, a preferred subnet (active node), route table IDs for automatic route management, and optionally an endpoint IP address range for floating IPs. Same throughput range as SINGLE_AZ_2.

Failover is transparent to NFS clients — the DNS name resolves to a floating IP that follows the active file server.

## Storage Types

### SSD

Solid-state drives. Sub-millisecond latency. Available for all deployment types. Storage capacity: 64–524,288 GiB, provisioned upfront and grown in place.

### INTELLIGENT_TIERING

Elastic, pay-for-what-you-store capacity with automatic tiering between a provisioned SSD read cache and cold storage. MULTI_AZ_1 only; forbids `storage_capacity_gib` and requires `read_cache_configuration` (`PROPORTIONAL_TO_THROUGHPUT_CAPACITY` for AWS-sized caching, `USER_PROVISIONED` with `size_gib` to pin economics, or `NO_CACHE` for purely archival patterns).

## Disk IOPS

- **AUTOMATIC** (default): 3 IOPS per GiB of storage, up to deployment type limit
- **USER_PROVISIONED**: Explicit IOPS independent of storage. SINGLE_AZ_1: up to 160,000. SINGLE_AZ_2/MULTI_AZ_1: up to 400,000.

## Root Volume Configuration

The root volume is automatically created with the file system. Key settings:

- **data_compression_type**: NONE, ZSTD (best ratio, ~2-3x), LZ4 (fastest, ~1.5-2x)
- **record_size_kib**: 4–1024 KiB. Default 128. Smaller for random I/O (databases), larger for sequential (analytics).
- **nfs_exports**: Client configurations with IP/CIDR/wildcard + mount options (rw, ro, crossmnt, root_squash, etc.)
- **user_and_group_quotas**: Per-UID/GID storage limits (AWS validates counts service-side).
- **read_only**: Makes the entire root volume read-only.

## Networking

- **Ports required**: TCP 111 (portmapper), TCP 2049 (NFS), TCP 20001-20003 (NFS mount)
- **SINGLE_AZ**: 1 ENI in the specified subnet
- **MULTI_AZ**: 2 ENIs (one per subnet), floating IP for failover

## Backup and Restore

- Automatic backups: 0–90 day retention (0 = disabled)
- Daily backup window: HH:MM UTC format
- Final backup on deletion: controlled by `skip_final_backup`; `final_backup_tags` tag the backup when one is taken
- Restore: `backup_id` creates the file system from an existing FSx backup (capacity comes from the backup)
- Tag propagation: `copy_tags_to_backups`, `copy_tags_to_volumes`

## Deletion Behavior

By default, deleting a file system fails while child volumes or snapshots exist — a guard rail against cascading data loss. `delete_options: [DELETE_CHILD_VOLUMES_AND_SNAPSHOTS]` opts into the cascade deliberately.

## ForceNew Attributes

Changing these requires replacing the file system (destructive):

- `deployment_type`
- `storage_type`
- `subnet_ids`
- `security_group_ids`
- `kms_key_id`
- `backup_id`
- `preferred_subnet_id`
- `endpoint_ip_address_range`
- `root_volume_configuration.copy_tags_to_snapshots`

## Scope

### Included
- File system creation with all five deployment types and both storage classes
- Root volume configuration (compression, NFS, quotas, record size)
- Disk IOPS configuration
- Backup, restore-from-backup, and deletion behavior
- Customer-managed KMS encryption
- Multi-AZ networking (preferred subnet, route tables, endpoint IP range)

### Separate resources (referenced via outputs, never embedded)
- Child volumes (`aws_fsx_openzfs_volume`) — join via the `root_volume_id` output
- Snapshots (`aws_fsx_openzfs_snapshot`) — point-in-time operations
