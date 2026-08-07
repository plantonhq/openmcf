---
title: "Server"
description: "Server deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudserver"
---

# Hetzner Cloud Server

Deploys a cloud server on Hetzner Cloud with configurable server type, OS image, location, SSH keys, firewall rules, placement group assignment, private network attachments, and public networking options. The server is the core compute resource in Hetzner Cloud, supporting both Intel/AMD (shared and dedicated) and ARM64 (Ampere) instance types. References to SSH keys, firewalls, placement groups, networks, and Primary IPs are wired via StringValueOrRef for infra-chart composability.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Server** -- an `hcloud_server` resource running the specified OS image on the chosen server type in the given location, with SSH keys, firewalls, network attachments, and public networking configured
- **Reverse DNS** (optional) -- an `hcloud_rdns` resource mapping the server's auto-assigned IPv4 address to a hostname when `dnsPtr` is specified

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **SSH key** -- at least one HetznerCloudSshKey registered in the same project for server access.
- **Location selection** -- choose a location (fsn1, nbg1, hel1, ash, hil, sin). Primary IPs and Floating IPs must be in the same location.
- **Server type selection** -- choose a type based on workload requirements (e.g., cx22 for 2 vCPU/4 GB, cpx11 for AMD, cax11 for ARM64).

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Server**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including server type, image, location, SSH keys, firewalls, and networking.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudServer
metadata:
  name: web-1
  org: acme-corp
  env: prod
spec:
  serverType: cx22
  image: ubuntu-24.04
  location: fsn1
  sshKeys:
    - value: "my-ssh-key"
  backups: true
```

```shell
planton apply -f hetznercloud-server.yaml
```

This creates a cx22 server running Ubuntu 24.04 in Falkenstein with SSH key access and daily backups. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a server environment, use ValueFromRef to wire SSH keys, firewalls, placement groups, and networks:

```yaml
spec:
  serverType: cx22
  image: ubuntu-24.04
  location: fsn1
  sshKeys:
    - valueFrom:
        kind: HetznerCloudSshKey
        name: deploy-key
        fieldPath: status.outputs.ssh_key_id
  firewallIds:
    - valueFrom:
        kind: HetznerCloudFirewall
        name: web-firewall
        fieldPath: status.outputs.firewall_id
  placementGroupId:
    valueFrom:
      kind: HetznerCloudPlacementGroup
      name: ha-group
      fieldPath: status.outputs.placement_group_id
  networks:
    - networkId:
        valueFrom:
          kind: HetznerCloudNetwork
          name: main-vpc
          fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph and provisions foundation resources before the server.

## Key Configuration

These are the most important decisions when configuring a server. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Server type** -- The `serverType` field selects the vCPU, RAM, and disk combination. Changing it triggers a server resize (temporary stop). Use `keepDisk: true` to prevent irreversible disk upgrades during downgrades.

**Image** -- The `image` field specifies the OS (e.g., "ubuntu-24.04", "debian-12", "rocky-9") or a snapshot ID. Changing this forces server replacement.

**Location** -- The `location` field determines the datacenter. Primary IPs, Floating IPs, and volumes must be co-located. Changing this forces replacement.

**SSH keys** -- The `sshKeys` field injects SSH keys at creation time. Accepts literal key names/IDs or ValueFromRef references to HetznerCloudSshKey outputs. Changing after creation forces replacement.

**Public networking** -- The `publicNet` block controls IPv4/IPv6 enablement and allows attaching existing Primary IPs instead of auto-assigned addresses. Omitting it uses default auto-assigned public IPs.

**Private networks** -- The `networks` field attaches the server to private networks for internal communication.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **HetznerCloudSshKey** | `sshKeys` | `status.outputs.ssh_key_id` |
| **HetznerCloudPlacementGroup** | `placementGroupId` | `status.outputs.placement_group_id` |
| **HetznerCloudFirewall** | `firewallIds` | `status.outputs.firewall_id` |
| **HetznerCloudNetwork** | `networks[].networkId` | `status.outputs.network_id` |
| **HetznerCloudPrimaryIp** | `publicNet.ipv4`, `publicNet.ipv6` | `status.outputs.primary_ip_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `server_id` | Hetzner Cloud numeric ID of the server | HetznerCloudVolume, HetznerCloudSnapshot, HetznerCloudFloatingIp, HetznerCloudLoadBalancer targets |
| `ipv4_address` | Public IPv4 address (empty if disabled) | DNS record configuration, application endpoints |
| `ipv6_address` | Public IPv6 address (empty if disabled) | DNS record configuration |
| `status` | Server status (running, off, etc.) | Monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web server** -- A cx22 or cx32 server with Ubuntu, SSH key, web firewall (ports 22/80/443), and optional private networking. The most common starting point.

**Private backend** -- A server with public networking disabled, attached to a private network, accessible only via a bastion host or load balancer.

## Works With

- [**Hetzner Cloud SSH Key**](/cloud-catalog/hetznercloud-ssh-key) -- provides SSH keys for server access
- [**Hetzner Cloud Placement Group**](/cloud-catalog/hetznercloud-placement-group) -- anti-affinity for HA deployments
- [**Hetzner Cloud Firewall**](/cloud-catalog/hetznercloud-firewall) -- network security rules
- [**Hetzner Cloud Network**](/cloud-catalog/hetznercloud-network) -- private networking
- [**Hetzner Cloud Primary IP**](/cloud-catalog/hetznercloud-primary-ip) -- stable public IP addresses
- [**Hetzner Cloud Volume**](/cloud-catalog/hetznercloud-volume) -- block storage attachment
- [**Hetzner Cloud Snapshot**](/cloud-catalog/hetznercloud-snapshot) -- disk image capture
- [**Hetzner Cloud Load Balancer**](/cloud-catalog/hetznercloud-load-balancer) -- traffic distribution target
