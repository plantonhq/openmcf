# Hetzner Cloud Primary IP

Allocates a managed public IP address (IPv4 or IPv6) on Hetzner Cloud that persists independently of any server. Primary IPs are created at a specific location and can be assigned to servers at creation time, providing stable public endpoints that survive server deletion and recreation. Includes optional reverse DNS (rDNS) configuration for mail servers and services requiring identity verification via reverse lookups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Primary IP** -- an `hcloud_primary_ip` resource allocating a public IPv4 address or IPv6 /64 network block at the specified location
- **Reverse DNS** (optional) -- an `hcloud_rdns` resource mapping the allocated IP to a hostname when `dnsPtr` is specified

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **Location selection** -- choose a Hetzner Cloud location (fsn1, nbg1, hel1, ash, hil, sin) where the IP will be allocated. The IP can only be assigned to servers in the same location.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Primary IP**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudPrimaryIp
metadata:
  name: web-ip
  org: acme-corp
  env: prod
spec:
  type: ipv4
  location: fsn1
  dnsPtr: "web.example.com"
```

```shell
planton apply -f hetznercloud-primary-ip.yaml
```

This allocates an IPv4 address in Falkenstein with a reverse DNS record pointing to web.example.com. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a server environment, use ValueFromRef to assign this IP to a server:

```yaml
# In the HetznerCloudServer manifest:
spec:
  publicNet:
    ipv4:
      valueFrom:
        kind: HetznerCloudPrimaryIp
        name: web-ip
        fieldPath: status.outputs.primary_ip_id
```

The InfraPipeline resolves the dependency graph, allocates the IP first, then provisions the server with the stable public address.

## Key Configuration

These are the most important decisions when configuring a Primary IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Type** -- The `type` field selects IPv4 (single address) or IPv6 (/64 network block). Changing this value forces replacement of the resource.

**Location** -- The `location` field determines where the IP is allocated (e.g., fsn1, nbg1, hel1, ash, hil, sin). The IP can only be assigned to servers in the same location. Changing this value forces replacement.

**Reverse DNS** -- The `dnsPtr` field sets an optional rDNS record mapping the IP to a hostname. Required for mail servers and services that rely on forward/reverse DNS consistency for identity verification.

**Delete protection** -- The `deleteProtection` field prevents accidental deletion via the Hetzner Cloud API.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `primary_ip_id` | Hetzner Cloud numeric ID of the Primary IP | HetznerCloudServer public network IP assignment |
| `ip_address` | Allocated IP address (IPv4) or first address in /64 block (IPv6) | DNS record configuration, application endpoints |
| `ip_network` | IPv6 network in CIDR notation (empty for IPv4) | IPv6 network planning |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Stable web endpoint** -- Allocate an IPv4 Primary IP and assign it to a web server. The IP survives server replacement, keeping DNS records and client connections stable.

**Mail server IP** -- Allocate an IPv4 Primary IP with rDNS configured to match the mail server hostname. Forward and reverse DNS consistency is required for SPF/DKIM verification.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers reference this Primary IP for stable public addressing
- [**Hetzner Cloud DNS Zone**](/cloud-catalog/hetznercloud-dns-zone) -- DNS records can reference the allocated `ip_address` output
