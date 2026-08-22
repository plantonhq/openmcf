# Firewall on DigitalOcean

Deploys a DigitalOcean Cloud Firewall with configurable inbound and outbound rules, applied to Droplets by reference or by tag. Firewalls filter traffic at the network edge before it reaches the Droplet, providing stateful packet inspection without consuming Droplet resources. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for Droplet, load balancer, and Kubernetes cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Cloud Firewall** -- a stateful, default-deny network firewall with the configured inbound and outbound rules, applied to the referenced Droplets and/or Droplet tags
- **Inbound Rules** -- which sources (CIDR addresses, Droplet tags, Droplets, Kubernetes clusters, load balancers) can reach the protected Droplets on which protocols and ports
- **Outbound Rules** -- which destinations the protected Droplets can reach on which protocols and ports

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **Droplets or Droplet tags to protect** -- target Droplets by reference (`dropletIds`, up to 10) or by tag (`tags`, up to 5; any Droplet carrying the tag is protected automatically, and DigitalOcean creates tags implicitly). Tag targeting is the production standard for anything long-lived.
- **Network planning** -- decide which protocols, ports, and sources your application needs. Firewalls deny everything not explicitly allowed, and at least one rule (in either direction) is required.

## Deploy

### Console

Open the deployment store, find **Firewall on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Tier** preset in the [Presets](#presets) tab for a web-facing firewall with HTTPS/HTTP inbound and restricted SSH.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanFirewall
metadata:
  name: web-firewall
  org: acme-corp
  env: prod
spec:
  firewallName: web-firewall
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
      portRange: all
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

**Targeting strategy** -- Use `tags` for dynamic targeting (any Droplet with the tag inherits the firewall) or `dropletIds` for explicit targeting by reference to `DigitalOceanDroplet` resources or literal numeric IDs. Tag targeting scales better where Droplets are created and destroyed by automation. The API caps targeting at 5 tags and 10 Droplets per firewall.

**Inbound rules** -- Each rule specifies a `protocol` (tcp, udp, or icmp), a `portRange` (`"443"`, `"8000-9000"`, or `all`; omitted for icmp), and source filters: CIDR addresses (IPv4 and IPv6), Droplet tags, and references to Droplets, Kubernetes clusters, or load balancers. Restrict SSH (port 22) to a management CIDR rather than `0.0.0.0/0`.

**Outbound rules** -- Same structure as inbound but controls egress. Most applications need all-TCP outbound plus UDP 53 for DNS; restrict egress on high-security tiers (databases) to contain exfiltration paths.

**Multi-tier architecture** -- Use `sourceTags` in inbound rules to allow traffic only from specific tiers: a database firewall allowing port 5432 only from Droplets tagged `web` enforces segmentation with zero IP management.

## Outputs and Dependencies

### What This Component Consumes

Optional references, resolved automatically at deploy time:

| Field | References | Purpose |
|-------|-----------|---------|
| `dropletIds` | DigitalOceanDroplet `droplet_id` | Droplets the firewall protects |
| `inboundRules[].sourceDropletIds` / `outboundRules[].destinationDropletIds` | DigitalOceanDroplet `droplet_id` | Droplets allowed as traffic sources/destinations |
| `inboundRules[].sourceKubernetesIds` / `outboundRules[].destinationKubernetesIds` | DigitalOceanKubernetesCluster `cluster_id` | Clusters allowed as traffic sources/destinations |
| `inboundRules[].sourceLoadBalancerUids` / `outboundRules[].destinationLoadBalancerUids` | DigitalOceanLoadBalancer `load_balancer_id` | Load balancers allowed as traffic sources/destinations |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewall_id` | Unique firewall identifier (UUID) | API operations, imports, audit logs |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web tier firewall** -- HTTPS/HTTP from everywhere, SSH from a management CIDR, open egress, tag-targeted. Start from the **Web Tier** preset.

**Database tier firewall** -- inbound only from the web tier's tag on the database port, SSH from a management CIDR, egress restricted to DNS and HTTPS. Start from the **Database Tier** preset.

## Works With

- **DigitalOceanDroplet** -- the Droplets this firewall protects and the rule-level source/destination Droplets
- **DigitalOceanLoadBalancer** -- rule-level sources so backends accept traffic only from their balancer
- **DigitalOceanKubernetesCluster** -- rule-level sources/destinations for cluster-to-Droplet traffic
