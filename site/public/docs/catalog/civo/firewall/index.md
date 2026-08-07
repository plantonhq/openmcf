---
title: "Firewall"
description: "Firewall deployment documentation"
icon: "package"
order: 100
componentName: "civofirewall"
---

# Firewall on Civo

Deploys a stateful firewall on Civo Cloud with configurable inbound and outbound rules, VPC network binding, and tag-based instance targeting. Firewalls control traffic to and from compute instances within a VPC -- inbound traffic not matching a rule is denied, while outbound traffic is allowed by default when no egress rules are specified. Integrates with Planton's Provider Connections for Civo credential management and ValueFromRef for VPC network dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Firewall** -- a stateful firewall attached to the specified VPC network, with the configured inbound (ingress) and outbound (egress) rules
- **Ingress Rules** -- created for each entry in `inboundRules`; each rule defines the protocol, port range, source CIDRs, and action (allow or deny)
- **Egress Rules** -- created only when `outboundRules` entries are specified; controls outbound traffic by protocol, port range, and destination CIDRs

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo VPC network** in the target region. Provide the network ID directly or reference a CivoVpc Cloud Resource via ValueFromRef.
- **Network CIDR planning** -- know the IP ranges of your application and data tiers to write precise source/destination CIDR rules. Avoid `0.0.0.0/0` for non-public-facing ports.

## Deploy

### Console

Open the deployment store, find **Firewall on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Tier** preset in the [Presets](#presets) tab for a standard internet-facing web server firewall.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoFirewall
metadata:
  name: web-firewall
  org: acme-corp
  env: prod
spec:
  name: web-firewall
  networkId:
    value: "abc12345-6789-def0-1234-567890abcdef"
  inboundRules:
    - protocol: tcp
      portRange: "443"
      cidrs:
        - "0.0.0.0/0"
      action: allow
      label: https
  tags:
    - web
```

```shell
planton apply -f civo-firewall.yaml
```

This creates a firewall allowing inbound HTTPS traffic from any source, attached to the specified VPC. No egress rules are defined, so all outbound traffic is allowed by default. Instances tagged `web` in the same network inherit this firewall. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the firewall to a VPC deployed in the same InfraPipeline:

```yaml
spec:
  networkId:
    valueFrom:
      kind: CivoVpc
      name: app-network
      fieldPath: status.outputs.network_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC first, then provisions the firewall on it.

## Key Configuration

These are the most important decisions when configuring a Civo firewall. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Inbound rules** -- Each entry in `inboundRules` defines a protocol (`tcp`, `udp`, or `icmp`), port range (e.g., `"443"`, `"8000-9000"`), source CIDRs, and action (`allow` or `deny`). Traffic not matching any rule is denied. For web servers, allow ports 80 and 443 from `0.0.0.0/0` and restrict SSH (port 22) to your admin IP.

**Outbound rules** -- By default, all outbound traffic is allowed when no `outboundRules` are specified. Add egress rules only when you need to restrict outbound access -- for example, limiting database servers to communicate only with specific API endpoints.

**Tag-based targeting** -- The `tags` list specifies instance tag names. Any compute instance in the same VPC network with a matching tag automatically inherits this firewall. Use tags like `web`, `database`, or `bastion` to apply firewalls across groups of instances without manual assignment.

**Network binding** -- The `networkId` field ties the firewall to a specific VPC network. Firewalls only apply to instances within that network. Use ValueFromRef to reference a CivoVpc Cloud Resource when deploying the VPC and firewall together.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CivoVpc** | `networkId` | `status.outputs.network_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `firewall_id` | Unique identifier of the firewall in Civo | Compute instance firewall assignment, API references |
| `created_at_rfc3339` | Firewall creation timestamp in RFC 3339 format | Audit logs, lifecycle tracking |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web tier firewall** -- allows inbound HTTP (80), HTTPS (443), and restricted SSH (22) access from specified admin CIDRs. No egress restrictions. Targets instances tagged `web`. Covers internet-facing web servers, API gateways, and reverse proxies. Start from the **Web Tier** preset.

**Database tier firewall** -- restricts inbound access to PostgreSQL (5432) and MySQL (3306) ports from the application tier CIDR only. No SSH or public internet access. Targets instances tagged `database`. Enforces multi-tier security isolation. Start from the **Database Tier** preset.

## Works With

- [**Civo VPC**](/cloud-catalog/civo-vpc) -- provides the VPC network to which the firewall is attached