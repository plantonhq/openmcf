# Intelligent-Tiering Elastic FSx Lustre

PERSISTENT_2 file system on the INTELLIGENT_TIERING storage class — elastic, pay-for-what-you-store capacity with 4000 MB/s of provisioned throughput and an AWS-sized SSD read cache. No storage capacity to provision, grow, or forecast.

## When to Use

- Large, cool datasets with hot working subsets (petabyte-scale data lakes, media archives feeding render farms)
- ML training corpora that grow continuously and are read in bursts
- Workloads where forecasting provisioned Lustre capacity is the operational pain
- Replacing S3-only pipelines that need POSIX semantics without paying for always-provisioned SSD

## What It Configures

- **INTELLIGENT_TIERING** — Capacity is elastic: you store what you store and pay for it; data tiers automatically between cache and cold storage
- **4000 MB/s provisioned throughput** — The baseline tier (`throughput_capacity` must be 4000 or a multiple of it)
- **Proportional read cache** — AWS sizes the SSD read cache from the provisioned throughput; all reads are served through it
- **AUTOMATIC metadata IOPS** — Part of the INTELLIGENT_TIERING contract; scales with usage
- **LZ4 compression** — Compressible data costs less and moves faster
- **Automatic backups** — 7-day retention

## What to Customize

- Replace placeholders: `<aws-region>`, `<subnet-id>`, `<security-group-id>`
- Raise `throughput_capacity` in 4000 MB/s steps as the hot working set grows
- Switch the read cache to `sizing_mode: USER_PROVISIONED` with `size_gib` (32-131072 GiB per 4000 MB/s of throughput) to pin cache economics
- Link S3 data with `AwsFsxDataRepositoryAssociation` resources referencing the `file_system_id` output (the legacy in-spec `import_path` arm does not apply to PERSISTENT_2)
