# Scale-Out FSx ONTAP (Multiple HA Pairs)

SINGLE_AZ_2 scale-out file system with 4 HA pairs delivering 6 GB/s aggregate throughput (1536 MB/s per pair) and 16 TiB of SSD storage with 160,000 provisioned IOPS. The shape for high-throughput enterprise workloads that outgrow a single HA pair.

## When to Use

- High-throughput analytics, media rendering, or EDA workloads
- Large VMware Cloud on AWS datastores consolidating many VMs
- Database farms needing aggregate IOPS beyond a single pair's ceiling
- Petabyte-track datasets: scale-out is the only path past 512 TiB (up to 1 PiB with 12 pairs)

## What It Configures

- **SINGLE_AZ_2 with 4 HA pairs** — The only deployment type supporting scale-out (1-12 pairs); pairs can be ADDED in place later
- **1536 MB/s per HA pair** — Scale-out tier (valid per-pair tiers with multiple pairs: 1536, 3072, 6144); total throughput = per-pair × pairs
- **16384 GiB SSD** — 4 TiB per HA pair (minimum is 1024 GiB per pair, maximum 512 TiB per pair)
- **160,000 provisioned IOPS** — USER_PROVISIONED mode decouples IOPS from capacity (up to 2.4M)
- **Customer-managed KMS + 7-day backups** — Production posture

## What to Customize

- Replace placeholders: `name`, `id`, `org`, `env`, `<aws-region>`, `<subnet-id>`, `<security-group-id>`, and the KMS key ARN; `fsx_admin_password` ships as a `$secret/<slug>` managed-secret reference — create the referenced secret, never inline a plaintext password
- Scale `ha_pairs` (1-12) with `storage_capacity_gib` ≥ 1024 × pairs
- Raise `throughput_capacity_per_ha_pair` to 3072 or 6144 per pair for hotter workloads
- Drop `disk_iops_configuration` to fall back to AUTOMATIC (3 IOPS per GiB)
- Use `valueFrom` references to wire AwsSubnet, AwsSecurityGroup, and AwsKmsKey
