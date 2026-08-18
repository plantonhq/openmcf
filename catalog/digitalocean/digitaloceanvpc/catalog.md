# VPC on DigitalOcean

Deploys a DigitalOcean Virtual Private Cloud with configurable CIDR range and region placement, providing an isolated private network for Droplets, Kubernetes clusters, databases, and load balancers. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean VPC** -- a regional private network with the configured IP range (or a DigitalOcean-assigned range if omitted) and optional description

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A target region** -- DigitalOcean VPCs are regional and cannot span regions. Choose the region where your resources will be deployed (e.g., `nyc1`, `sfo3`, `fra1`).
- **IP address planning** (optional) -- if you need specific, non-overlapping CIDR blocks across multiple VPCs, plan your ranges before deployment. DigitalOcean accepts prefix lengths from /16 through /24, and the range is immutable after creation.

## Deploy

### Console

Open the deployment store, find **VPC on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard** preset in the [Presets](#presets) tab to pre-populate a working configuration with an explicit /16 CIDR block.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
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

This creates a VPC in the NYC1 region with a DigitalOcean-assigned, non-conflicting IP range (reported through the `ip_range` output). A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**IP range** -- The `ipRangeCidr` field accepts CIDR blocks with prefix lengths /16 through /24 (e.g., `"10.10.0.0/16"`). When omitted, DigitalOcean assigns a non-conflicting range and reports it through the `ip_range` output. Specify an explicit range for production environments with IPAM requirements or when running multiple VPCs that must not overlap. The range is immutable -- changing it later replaces the VPC.

**Region** -- The `region` field determines where the VPC is created. All resources placed in this VPC (Droplets, databases, Kubernetes clusters) must be in the same region.

Whether a VPC is the region's DEFAULT is computed by DigitalOcean and cannot be set here -- wire every resource's `vpc` reference explicitly instead of relying on regional defaults.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpc_id` | Unique VPC identifier (UUID) on DigitalOcean | Droplet VPC placement, Kubernetes cluster networking, database cluster VPC attachment, load balancer VPC placement |
| `ip_range` | The VPC's CIDR range as DigitalOcean reports it (covers the auto-assigned case) | Firewall rules, peering and VPN planning |
| `urn` | The VPC's uniform resource name (`do:vpc:<uuid>`) | DigitalOcean project assignment, audit |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard VPC** -- explicit /16 CIDR block providing 65,536 IPs, suitable for production environments with multiple Droplets, Kubernetes clusters, and databases. Adjust the second octet to avoid overlap when running multiple VPCs. Start from the **Standard** preset.

## Works With

This component operates independently and does not reference other components.