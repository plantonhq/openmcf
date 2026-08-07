---
title: "Droplet/VM"
description: "Droplet/VM deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceandroplet"
---

# Droplet/VM on DigitalOcean

Deploys a DigitalOcean Droplet with configurable sizing, VPC placement, optional backups, IPv6 networking, block storage volume attachments, and cloud-init user data. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for VPC and volume dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Droplet** -- a virtual machine in the specified region and VPC, using the configured size slug, base image, and optional cloud-init user data script
- **Block Storage Volume Attachments** -- created only when `volumeIds` are provided; attaches existing volumes to the Droplet (volumes must reside in the same region)
- **DigitalOcean Monitoring Agent** -- enabled by default unless `disableMonitoring` is set to true
- **DigitalOcean Tags** -- user-supplied tags applied to the Droplet for firewall targeting and resource organization

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **A VPC network** in the target region. Provide the VPC UUID directly or reference a DigitalOceanVpc Cloud Resource via ValueFromRef.
- **A valid Droplet size slug** (e.g., `"s-2vcpu-4gb"`) -- check available sizes via the DigitalOcean API or CLI (`doctl compute size list`).
- **A valid image slug** (e.g., `"ubuntu-24-04-x64"`) -- check available images via `doctl compute image list --public`.
- **Block storage volumes** (optional) -- must be pre-created in the same region if you plan to attach them via `volumeIds`.

## Deploy

### Console

Open the deployment store, find **Droplet/VM on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production** preset in the [Presets](#presets) tab to pre-populate a working configuration with backups and VPC isolation.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanDroplet
metadata:
  name: app-server
  org: acme-corp
  env: prod
spec:
  dropletName: app-server
  region: nyc1
  size: s-2vcpu-4gb
  image: ubuntu-24-04-x64
  vpc:
    value: "abc12345-6789-def0-1234-567890abcdef"
```

```shell
planton apply -f do-droplet.yaml
```

This creates a Droplet with the specified size and image in the NYC1 region. No backups, IPv6, or volume attachments are configured. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the Droplet to a VPC deployed in the same InfraPipeline:

```yaml
spec:
  vpc:
    valueFrom:
      kind: DigitalOceanVpc
      name: app-network
      fieldPath: status.outputs.vpc_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the Droplet within it.

## Key Configuration

These are the most important decisions when configuring a Droplet. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Sizing** -- The `size` field sets the Droplet's CPU and memory (e.g., `"s-1vcpu-1gb"` for development, `"s-2vcpu-4gb"` for production web servers, `"c-4vcpu-8gb"` for CPU-intensive workloads). Droplets can be resized after creation, but disk-only downgrades are not supported.

**Backups** -- Set `enableBackups: true` to enable DigitalOcean's automated weekly snapshots with 4-week retention. Recommended for production. Disabled by default to reduce cost in development.

**Cloud-init user data** -- The `userData` field accepts a cloud-init script (up to 32 KiB) for bootstrapping the Droplet on first boot. Use it to install packages, configure services, or pull application code.

**Volume attachments** -- The `volumeIds` field accepts a list of block storage volume UUIDs or ValueFromRef references. Volumes must be in the same region as the Droplet and are attached at creation time.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **DigitalOceanVpc** | `vpc` | `status.outputs.vpc_id` |
| **DigitalOceanVolume** (optional) | `volumeIds` | `status.outputs.volume_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `droplet_id` | Unique Droplet identifier on DigitalOcean | Load balancer backend attachment, firewall rules |
| `ipv4_address` | Primary IPv4 address (public or private) | DNS records, application configuration |
| `ipv6_address` | IPv6 address (if IPv6 was enabled) | IPv6 DNS records |
| `image_id` | Image ID of the Droplet's base image | Audit logs, image tracking |
| `vpc_uuid` | VPC network UUID the Droplet resides in | Verifying network placement |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production server** -- 2 vCPU / 4 GB instance with automated backups, VPC isolation, and tags for firewall targeting. Suitable for web servers, API backends, and application hosts. Start from the **Production** preset.

**Development instance** -- 1 vCPU / 1 GB instance with no backups, minimal cost for dev/test workloads, CI/CD build agents, and experimentation. Start from the **Development** preset.

## Works With

- [**DigitalOcean VPC**](/cloud-catalog/digital-ocean-vpc) -- provides the private network for Droplet placement
- [**DigitalOcean Volume**](/cloud-catalog/digital-ocean-volume) -- provides block storage volumes attached to the Droplet