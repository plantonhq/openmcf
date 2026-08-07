# Multi-AZ High Availability FSx ONTAP

MULTI_AZ_2 deployment with automatic failover across two availability zones. 2 TiB SSD, 768 MB/s throughput. 7-day backups, customer-managed KMS encryption. Mission-critical configuration for workloads requiring high availability.

## When to Use

- Mission-critical workloads requiring automatic failover
- Production databases that cannot tolerate AZ downtime
- VMware Cloud on AWS datastores with HA requirements
- Compliance-sensitive workloads requiring multi-AZ resilience

## What It Configures

- **MULTI_AZ_2** — Current-generation multi-AZ deployment with automatic failover
- **2048 GiB SSD** — 2 TiB storage. Sub-millisecond latency
- **768 MB/s throughput** — Production-grade second-generation throughput tier
- **1 HA pair** — Fixed for multi-AZ (active in preferred subnet, standby in second)
- **Two subnets** — One per AZ. Must be in different availability zones
- **Preferred subnet** — Active file server placement (required for multi-AZ)
- **Endpoint IP range** — Floating-IP CIDR outside the VPC range (seamless failover)
- **Managed route tables** — AWS repoints routes to the floating IPs on failover
- **Customer-managed KMS** — Encryption at rest
- **7-day backups** — Daily automatic backups at 05:00 UTC

## What to Customize

- Replace placeholders: `name`, `id`, `org`, `env`, `<aws-region>`, both subnet IDs, the route table ID, `<security-group-id>`, and the KMS key ARN
- Keep `endpoint_ip_address_range` OUTSIDE the VPC CIDR (AWS recommends the 198.19.0.0/16 block) and non-overlapping with any subnet; omit it to let AWS pick an unused range
- List every route table associated with your clients' subnets in `route_table_ids`
- Use `valueFrom` references to wire AwsSubnet, AwsSecurityGroup, and AwsKmsKey
- Increase `storage_capacity_gib` or `throughput_capacity_per_ha_pair` (1536, 3072, 6144) as needed
