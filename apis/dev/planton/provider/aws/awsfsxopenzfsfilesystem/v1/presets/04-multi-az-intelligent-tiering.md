# Multi-AZ Intelligent-Tiering FSx OpenZFS

MULTI_AZ_1 file system on the INTELLIGENT_TIERING storage class — elastic, pay-for-what-you-store NFS capacity with 1280 MB/s of provisioned throughput, an AWS-sized SSD read cache, and automatic cross-AZ failover. No storage capacity to provision or forecast.

## When to Use

- Shared NFS datasets that grow unpredictably (content libraries, build artifact stores, research data)
- Cool datasets with hot working subsets where always-provisioned SSD is wasteful
- High-availability requirements (cross-AZ active/standby failover) combined with elastic economics

## What It Configures

- **MULTI_AZ_1** — Active file server in the preferred subnet, standby in the second AZ; floating endpoint IPs fail over automatically (routes managed in the given route tables)
- **INTELLIGENT_TIERING** — Capacity is elastic; data tiers automatically between the SSD read cache and cold storage
- **1280 MB/s provisioned throughput** — Gen-2 tier; scale up or down in place
- **Proportional read cache** — AWS sizes the SSD cache from the provisioned throughput
- **ZSTD compression** on the root volume, exported over NFS to `10.0.0.0/16`
- **Automatic backups** — 7-day retention

## What to Customize

- Replace placeholders: `<aws-region>`, subnet/route-table/security-group IDs
- Adjust `throughput_capacity` within the gen-2 set (160–10240 MB/s)
- Pin cache economics with `sizing_mode: USER_PROVISIONED` + `size_gib`
- Tighten the NFS export `clients` range to your compute subnets
- Create child volumes as separate resources against the `root_volume_id` output
