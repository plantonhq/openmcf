# VPC on DigitalOcean

Deploys a DigitalOcean Virtual Private Cloud with configurable CIDR range and region placement, providing an isolated private network for Droplets, Kubernetes clusters, databases, and load balancers. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean VPC** -- a regional private network with the configured IP range (or an auto-generated /20 CIDR if omitted), optional description, and default-for-region setting

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A target region** -- DigitalOcean VPCs are regional and cannot span regions. Choose the region where your resources will be deployed (e.g., `nyc1`, `sfo3`, `fra1`).
- **IP address planning** (optional) -- if you need specific, non-overlapping CIDR blocks across multiple VPCs, plan your ranges before deployment. Only /16, /20, and /24 blocks are supported.

## Deploy

### Console

Open the deployment store, find **VPC on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration with an explicit /16 CIDR block.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanVpc
metadata:
  name: app-network
  org: acme-corp
  env: prod
spec:
  region: nyc1
```

```shell
planton apply -f do-vpc.yaml
```

This creates a VPC in the NYC1 region with a DigitalOcean auto-generated /20 CIDR block. No explicit IP range or default-for-region setting is configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**IP range** -- The `ipRangeCidr` field accepts /16, /20, or /24 CIDR blocks (e.g., `"10.10.0.0/16"`). When omitted, DigitalOcean auto-generates a non-conflicting /20 block with 4,096 IPs. Specify an explicit range for production environments with IPAM requirements or when running multiple VPCs that must not overlap.

**Region** -- The `region` field determines where the VPC is created. All resources placed in this VPC (Droplets, databases, Kubernetes clusters) must be in the same region.

**Default for region** -- Set `isDefaultForRegion: true` to make this VPC the automatic default for new resources in the region. Only one VPC can be the default per region.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpc_id` | Unique VPC identifier (UUID) on DigitalOcean | Droplet VPC placement, Kubernetes cluster networking, database cluster VPC attachment, load balancer VPC placement |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard VPC** -- explicit /16 CIDR block providing 65,536 IPs, suitable for production environments with multiple Droplets, Kubernetes clusters, and databases. Adjust the second octet to avoid overlap when running multiple VPCs. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.