# Firewall on DigitalOcean

Deploys a DigitalOcean Cloud Firewall with configurable inbound and outbound rules, applied to Droplets by ID or tag. Firewalls filter traffic at the network edge before it reaches the Droplet, providing stateful packet inspection without consuming Droplet resources. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud Firewall** -- a network firewall with the configured inbound and outbound rules, applied to the specified Droplets or Droplet tags
- **Inbound Rules** -- traffic filtering rules that control which sources can reach Droplets on specific protocols and port ranges
- **Outbound Rules** -- traffic filtering rules that control which destinations Droplets can reach on specific protocols and port ranges

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **Existing Droplets or Droplet tags** -- the firewall targets Droplets by ID (via `dropletIds`) or by tag name (via `tags`). Tag-based targeting is recommended for dynamic environments where Droplets are created and destroyed frequently.
- **Network planning** -- determine which protocols, ports, and source/destination addresses your application requires. Firewalls deny all traffic by default; only explicitly allowed traffic passes through.

## Deploy

### Console

Open the deployment store, find **Firewall on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Tier** preset in the [Presets](#presets) tab for a web-facing firewall with HTTP/HTTPS inbound and restricted SSH.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digitalocean.planton.dev/v1
kind: DigitalOceanFirewall
metadata:
  name: web-firewall
  org: acme-corp
  env: prod
spec:
  name: web-firewall
  tags:
    - web
  inboundRules:
    - protocol: tcp
      portRange: "443"
      sourceAddresses:
        - "0.0.0.0/0"
        - "::/0"
    - protocol: tcp
      portRange: "80"
      sourceAddresses:
        - "0.0.0.0/0"
        - "::/0"
  outboundRules:
    - protocol: tcp
      portRange: "1-65535"
      destinationAddresses:
        - "0.0.0.0/0"
        - "::/0"
```

```shell
planton apply -f do-firewall.yaml
```

This creates a firewall allowing HTTPS and HTTP inbound from all addresses and all TCP outbound, applied to Droplets tagged `web`. SSH access is not included. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Cloud Firewall. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Targeting strategy** -- Use `tags` for dynamic targeting (any Droplet with the tag inherits the firewall) or `dropletIds` for explicit targeting. Tag-based targeting scales better in environments where Droplets are created and destroyed by automation. A maximum of 5 tags and 10 Droplet IDs can be specified per firewall.

**Inbound rules** -- Each rule specifies a `protocol` (tcp, udp, or icmp), a `portRange` (e.g., `"443"`, `"8000-9000"`, or `"1-65535"` for all), and source filters. Sources can be CIDR addresses, Droplet IDs, Droplet tags, Kubernetes cluster IDs, or load balancer UIDs. Restrict SSH (port 22) to a management CIDR rather than `0.0.0.0/0`.

**Outbound rules** -- Same structure as inbound but controls egress traffic. Most applications need all-TCP outbound for API calls and package downloads, plus UDP 53 for DNS resolution. Restrict outbound only when compliance or data-loss prevention requires it.

**Multi-tier architecture** -- Use `sourceTags` in inbound rules to allow traffic only from specific tiers. For example, a database firewall can allow port 5432 only from Droplets tagged `web`, enforcing network segmentation without IP-based rules.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewall_id` | Unique firewall identifier (UUID) | API operations, audit logs |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web tier firewall** -- allows HTTP/HTTPS from all addresses, restricts SSH to a management CIDR, and permits all outbound traffic. Applied via tag to any Droplet tagged `web`. Start from the **Web Tier** preset.

**Database tier firewall** -- restricts inbound to a database port (default 5432) from Droplets tagged `web` only, with SSH limited to a management CIDR. Prevents direct internet access to database servers. Start from the **Database Tier** preset.

## Works With

This component operates independently and does not reference other components.