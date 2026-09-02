# Single-AZ Production FSx ONTAP

SINGLE_AZ_2 SSD file system with 2 TiB (2048 GiB) and 768 MB/s throughput. One HA pair. Customer-managed KMS encryption. 7-day automatic backups at 05:00 UTC. Production-ready configuration for non-HA workloads.

## When to Use

- Production workloads that can tolerate single-AZ availability
- Database storage (Oracle, SAP, SQL Server) on shared NFS/iSCSI
- VMware Cloud on AWS datastores
- Enterprise file shares with compliance requirements (encryption, backups)

## What It Configures

- **SINGLE_AZ_2** — Current-generation single-AZ deployment
- **2048 GiB SSD** — 2 TiB storage. Sub-millisecond latency
- **768 MB/s throughput** — Production-grade second-generation throughput tier
- **1 HA pair** — Standard redundancy within the AZ
- **Customer-managed KMS** — Encryption at rest with your key
- **ONTAP admin access** — `fsx_admin_password` enables the ONTAP CLI and REST API
- **7-day backups** — Daily automatic backups at 05:00 UTC
- **Weekly maintenance** — Sunday at 02:00 UTC

## What to Customize

- Replace placeholders: `name`, `id`, `org`, `env`, `<aws-region>`, `<subnet-id>`, `<security-group-id>`, and the KMS key ARN; `fsx_admin_password` ships as a `$secret/<slug>` managed-secret reference — create the referenced secret, never inline a plaintext password
- Increase `storage_capacity_gib` for larger datasets
- Increase `throughput_capacity_per_ha_pair` (1536, 3072, 6144) for higher I/O
- Add `ha_pairs: 2` or more for scale-out throughput (see preset 04)
- Adjust `automatic_backup_retention_days` (up to 90) for longer retention
- Use `valueFrom` references to wire AwsSubnet, AwsSecurityGroup, and AwsKmsKey
- Switch to preset 03 for multi-AZ high availability
