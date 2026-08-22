# Standard VPC

This preset creates a DigitalOcean VPC with an explicit /16 CIDR block, providing a private isolated network for Droplets, Kubernetes clusters, databases, and load balancers within a single region. This is the most common production VPC configuration.

## When to Use

- Any new DigitalOcean environment that needs private networking
- Production workloads requiring predictable, non-overlapping IP ranges
- Environments where multiple resources (Droplets, DOKS, databases) must communicate privately

## Key Configuration Choices

- **Explicit /16 CIDR** (`ipRangeCidr: 10.10.0.0/16`) -- provides 65,536 IPs, sufficient for most production environments. Adjust the second octet (`10.10`, `10.20`, etc.) to avoid overlap when running multiple VPCs. The range is immutable -- changing it later replaces the VPC.
- **Region** (`region: nyc1`) -- placeholder; change to your target region. DigitalOcean VPCs are regional and cannot span regions.
- **Explicit membership** -- whether a VPC is the region's default is computed by DigitalOcean and cannot be set here; wire every resource's `vpc` reference to this VPC's `vpc_id` output instead of relying on regional defaults.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `nyc1` | Target DigitalOcean region slug | [DigitalOcean Regions API](https://docs.digitalocean.com/reference/api/api-reference/#tag/Regions) |
| `10.10.0.0/16` | VPC CIDR block (prefix /16 through /24) | Your IP address management plan |
