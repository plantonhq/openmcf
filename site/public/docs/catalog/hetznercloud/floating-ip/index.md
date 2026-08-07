---
title: "Floating IP"
description: "Floating IP deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudfloatingip"
---

# Hetzner Cloud Floating IP

Allocates a reassignable public IP address (IPv4 or IPv6) on Hetzner Cloud that can be moved between servers in the same location. Floating IPs persist independently of any server and can be reassigned at any time, making them ideal for failover scenarios where a stable endpoint must survive server replacement. Includes optional reverse DNS and initial server assignment.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Floating IP** -- an `hcloud_floating_ip` resource allocating a public IPv4 address or IPv6 /64 network block at the specified home location
- **Server Assignment** (optional) -- an `hcloud_floating_ip_assignment` resource attaching the IP to a server when `serverId` is specified
- **Reverse DNS** (optional) -- an `hcloud_rdns` resource mapping the allocated IP to a hostname when `dnsPtr` is specified

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **Location selection** -- choose a Hetzner Cloud location (fsn1, nbg1, hel1, ash, hil, sin) as the home location for the IP. The IP can only be assigned to servers in the same location.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Floating IP**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudFloatingIp
metadata:
  name: failover-ip
  org: acme-corp
  env: prod
spec:
  type: ipv4
  homeLocation: fsn1
  description: "Production failover IP"
  serverId:
    value: "12345678"
  dnsPtr: "app.example.com"
```

```shell
planton apply -f hetznercloud-floating-ip.yaml
```

This allocates an IPv4 Floating IP in Falkenstein, assigns it to the specified server, and configures reverse DNS. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of an HA cluster, use ValueFromRef to assign the Floating IP to a primary server:

```yaml
spec:
  serverId:
    valueFrom:
      kind: HetznerCloudServer
      name: primary-server
      fieldPath: status.outputs.server_id
```

The InfraPipeline resolves the dependency graph, creates the server first, then assigns the Floating IP to it.

## Key Configuration

These are the most important decisions when configuring a Floating IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type** -- The `type` field selects IPv4 (single address) or IPv6 (/64 network block). Changing this value forces replacement of the resource.

**Home location** -- The `homeLocation` field determines where the IP is homed (e.g., fsn1, nbg1, hel1, ash, hil, sin). The IP can only be assigned to servers in the same location. Changing this value forces replacement.

**Server assignment** -- The `serverId` field optionally attaches the Floating IP to a server at creation time. Accepts a literal server ID or a ValueFromRef reference to a HetznerCloudServer output. If omitted, the IP is created unassigned (reserved) and can be assigned later.

**Reverse DNS** -- The `dnsPtr` field sets an optional rDNS record mapping the IP to a hostname. Required for mail servers and services relying on reverse DNS consistency.

**Delete protection** -- The `deleteProtection` field prevents accidental deletion via the Hetzner Cloud API.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **HetznerCloudServer** (optional) | `serverId` | `status.outputs.server_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `floating_ip_id` | Hetzner Cloud numeric ID of the Floating IP | Monitoring, DNS record creation |
| `ip_address` | Allocated IP address (IPv4) or first address in /64 block (IPv6) | DNS record configuration, application endpoints |
| `ip_network` | IPv6 network in CIDR notation (empty for IPv4) | IPv6 network planning |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Failover IP** -- Allocate a Floating IP and assign it to a primary server. On failure, reassign the IP to a standby server -- the public endpoint stays constant while the underlying server changes.

**Reserved IP** -- Allocate a Floating IP without server assignment for future use. Useful for reserving a stable IP before the target server is provisioned.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- Floating IPs are assigned to servers for stable public addressing
- [**Hetzner Cloud DNS Zone**](/cloud-catalog/hetznercloud-dns-zone) -- DNS records can reference the allocated `ip_address` output
