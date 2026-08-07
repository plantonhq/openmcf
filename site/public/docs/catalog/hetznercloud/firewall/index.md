---
title: "Firewall"
description: "Firewall deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudfirewall"
---

# Hetzner Cloud Firewall

Deploys a stateful firewall with inline rules that control inbound and outbound network traffic for Hetzner Cloud servers. Firewalls are deny-by-default for inbound traffic -- when applied to a server, all inbound packets not matching a rule are dropped, while outbound traffic is allowed unless explicitly restricted. Supports TCP, UDP, ICMP, ESP, and GRE protocols with CIDR-based source and destination filtering. Hetzner Cloud allows up to 50 rules per firewall.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Firewall** -- a single `hcloud_firewall` resource with inline rules defining allowed traffic directions, protocols, ports, and CIDR blocks

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.
- **Network planning** -- decide which ports and protocols to allow for inbound traffic, and which CIDR blocks should be permitted as sources.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Firewall**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including rule definitions.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudFirewall
metadata:
  name: web-firewall
  org: acme-corp
  env: prod
spec:
  rules:
    - direction: in
      protocol: tcp
      port: "22"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "Allow SSH from anywhere"
    - direction: in
      protocol: tcp
      port: "80"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "Allow HTTP"
    - direction: in
      protocol: tcp
      port: "443"
      sourceIps:
        - "0.0.0.0/0"
        - "::/0"
      description: "Allow HTTPS"
```

```shell
planton apply -f hetznercloud-firewall.yaml
```

This creates a firewall allowing SSH, HTTP, and HTTPS inbound traffic from all sources. A Stack Job tracks the provisioning in real time. Reference the firewall in HetznerCloudServer manifests via `firewall_ids`.

### InfraChart

When deploying as part of a server environment, use ValueFromRef to wire the firewall to servers:

```yaml
# In the HetznerCloudServer manifest:
spec:
  firewallIds:
    - valueFrom:
        kind: HetznerCloudFirewall
        name: web-firewall
        fieldPath: status.outputs.firewall_id
```

The InfraPipeline resolves the dependency graph, creates the firewall first, then provisions servers with the firewall applied.

## Key Configuration

These are the most important decisions when configuring a firewall. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Rules** -- The `rules` field is a repeated list of firewall rules. Each rule specifies a `direction` (in or out), `protocol` (tcp, udp, icmp, esp, gre), optional `port` (required for TCP/UDP -- accepts single port, range, or "any"), and CIDR blocks for `sourceIps` (inbound) or `destinationIps` (outbound). An empty rules list creates a firewall that blocks all inbound and allows all outbound traffic.

**Direction** -- Inbound rules (`in`) require `sourceIps`. Outbound rules (`out`) require `destinationIps`. Use `["0.0.0.0/0", "::/0"]` to match all IPv4 and IPv6 traffic.

**Port** -- Required for TCP and UDP protocols. Accepts a single port (`"80"`), a range (`"80-443"`), or `"any"` for all ports. Must not be set for ICMP, ESP, or GRE.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewall_id` | Hetzner Cloud numeric ID of the firewall | HetznerCloudServer `firewallIds` field |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web server firewall** -- Allow SSH (22), HTTP (80), and HTTPS (443) inbound from all sources. The standard starting point for web-facing servers.

**Database firewall** -- Allow inbound connections only from private network CIDR blocks on database ports (e.g., 5432 for PostgreSQL, 3306 for MySQL). Restrict SSH to a bastion IP range.

## Works With

- [**Hetzner Cloud Server**](/cloud-catalog/hetznercloud-server) -- servers reference this firewall via `firewallIds` to apply traffic rules
