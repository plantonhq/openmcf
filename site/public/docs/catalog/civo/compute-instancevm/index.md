---
title: "Compute Instance/VM"
description: "Compute Instance/VM deployment documentation"
icon: "package"
order: 100
componentName: "civocomputeinstance"
---

# Compute Instance/VM on Civo

Deploys a virtual machine on Civo Cloud with configurable instance sizing, OS image selection, VPC networking, firewall attachment, volume mounting, and cloud-init scripting. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for wiring VPC, firewall, volume, and reserved IP dependencies.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Instance** -- a virtual machine in the specified region with the configured size, OS image, and hostname, placed on the referenced VPC network
- **Network Attachment** -- the instance is connected to the specified VPC network for private connectivity
- **Firewall Binding** -- created only when `firewallIds` is provided, attaching the instance to an existing firewall for traffic control
- **Reserved IP Assignment** -- created only when `reservedIpId` is provided, binding a static public IPv4 address to the instance
- **Volume Attachment** -- created only when `volumeIds` is provided, mounting existing storage volumes to the instance
- **Civo Tags** -- metadata tags applied to the instance for organizational tracking and tag-based firewall assignment

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo VPC network** in the target region. Provide the network ID directly or reference a CivoVpc Cloud Resource via ValueFromRef.
- **An OS image slug** for the instance (e.g., `ubuntu-jammy` for Ubuntu 22.04 LTS). List available images via the Civo CLI (`civo diskimage ls`) or dashboard.
- **SSH key IDs** (optional) -- pre-registered SSH public keys on your Civo account for passwordless login. If omitted, Civo sets a root password.

## Deploy

### Console

Open the deployment store, find **Compute Instance/VM on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production Web** preset in the [Presets](#presets) tab for a production-grade instance with firewall and cloud-init.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoComputeInstance
metadata:
  name: web-server
  org: acme-corp
  env: prod
spec:
  instanceName: web-server
  region: lon1
  size: g3.medium
  image: ubuntu-jammy
  network:
    value: "abc12345-6789-def0-1234-567890abcdef"
```

```shell
planton apply -f civo-instance.yaml
```

This creates a medium-sized Ubuntu 22.04 instance on the specified VPC network in Civo's London region. No firewall, volumes, or reserved IP are attached. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a VPC and firewall deployed in the same InfraPipeline:

```yaml
spec:
  network:
    valueFrom:
      kind: CivoVpc
      name: app-network
      fieldPath: status.outputs.network_id
  firewallIds:
    - valueFrom:
        kind: CivoFirewall
        name: web-firewall
        fieldPath: status.outputs.firewall_id
  reservedIpId:
    valueFrom:
      kind: CivoIpAddress
      name: web-ip
      fieldPath: status.outputs.reserved_ip_id
  volumeIds:
    - valueFrom:
        kind: CivoVolume
        name: data-vol
        fieldPath: status.outputs.volume_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC, firewall, IP address, and volume first, then provisions the instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring a Civo compute instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Instance size** -- The `size` field sets the instance flavor (e.g., `g3.small` for development, `g3.medium` for production web workloads). This determines CPU, memory, and local disk allocation. Check available sizes via the Civo CLI (`civo instance size`) or dashboard.

**OS image** -- The `image` field selects the base operating system (e.g., `ubuntu-jammy` for Ubuntu 22.04 LTS, `debian-11` for Debian). Choose an LTS image for production to align with your patching lifecycle.

**Cloud-init script** -- The `userData` field accepts a cloud-init script (up to 32 KiB) that runs on first boot. Use it to install packages, apply security updates, and configure the instance without manual SSH access.

**Firewall and tagging** -- Attach firewalls via `firewallIds` for explicit assignment, or use `tags` to enable tag-based firewall auto-assignment. Tag-based assignment simplifies multi-instance deployments where all instances share the same security posture.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoVpc** | `network` | `status.outputs.network_id` |
| **CivoFirewall** (optional) | `firewallIds` | `status.outputs.firewall_id` |
| **CivoVolume** (optional) | `volumeIds` | `status.outputs.volume_id` |
| **CivoIpAddress** (optional) | `reservedIpId` | `status.outputs.reserved_ip_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | Unique identifier of the instance on Civo | Civo API operations, monitoring dashboards |
| `public_ipv4` | Public IPv4 address assigned to the instance | DNS records, external access configuration |
| `private_ipv4` | Private IPv4 address within the VPC network | Service discovery, internal load balancing |
| `status` | Current instance status (e.g., ACTIVE, BUILDING) | Health monitoring, deployment verification |
| `created_at_rfc3339` | Instance creation timestamp in RFC 3339 format | Audit logs, lifecycle tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Production web server** -- a medium-sized Ubuntu instance with VPC networking, firewall protection, and cloud-init for automatic security updates on first boot. Suitable for web servers, API backends, and reverse proxies. Start from the **Production Web** preset.

**Development instance** -- a small, minimal instance with VPC networking but no firewall or cloud-init. Lowest cost for development, CI/CD agents, and temporary workloads. Start from the **Development** preset.

## Works With

- [**Civo VPC**](/cloud-catalog/civo-vpc) -- provides the VPC network for instance connectivity
- [**Civo Firewall**](/cloud-catalog/civo-firewall) -- controls inbound and outbound traffic to the instance
- [**Civo Volume**](/cloud-catalog/civo-volume) -- provides persistent block storage volumes for the instance
- [**Civo IP Address**](/cloud-catalog/civo-ip-address) -- provides a static public IPv4 address for the instance