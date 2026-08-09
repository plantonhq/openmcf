# OpenStack Security Group

Deploys a Neutron security group on OpenStack with optional inline firewall rules for controlling ingress and egress traffic on instances and ports. The security group supports stateful and stateless modes, zero-trust baselines via default rule deletion, and inline rules keyed for stable IaC state management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Security Group** -- a virtual firewall with configurable description, stateful/stateless mode, and default rule deletion
- **Inline Security Group Rules** -- created only when `rules` entries are specified; one rule per entry, each provisioned as a separate resource keyed by the rule's `key` field for stable IaC state
- **OpenStack Tags** -- user-defined tags applied to the security group for filtering and organization in the OpenStack API and Horizon dashboard

## Before You Deploy

### OpenStack Account

- **Rule planning** -- decide whether to manage rules inline (via the `rules` field in this component) or as standalone OpenStackSecurityGroupRule resources. Inline rules are simpler for self-contained groups. Standalone rules provide DAG-level visibility and cross-group references in InfraCharts.
- **Stateless support** -- if you plan to use `stateful: false`, confirm your OpenStack deployment supports stateless security groups. Not all Neutron backends implement this feature.
- **Default rule behavior** -- OpenStack creates two default egress rules on every new security group (allow all IPv4 and IPv6 outbound). Decide whether to keep them or set `deleteDefaultRules: true` for a zero-trust baseline.

## Deploy

### Console

Open the deployment store, find **OpenStack Security Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Server** preset in the [Presets](#presets) tab to pre-populate a common configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackSecurityGroup
metadata:
  name: web-sg
  org: acme-corp
  env: prod
spec:
  description: "Web server security group"
  rules:
    - key: allow-https
      direction: ingress
      ethertype: IPv4
      protocol: tcp
      portRangeMin: 443
      portRangeMax: 443
      remoteIpPrefix: "0.0.0.0/0"
```

```shell
planton apply -f security-group.yaml
```

This creates a security group with one inline rule allowing inbound HTTPS from all sources. OpenStack's default egress rules (allow all IPv4/IPv6 outbound) remain. No SSH or HTTP rules are configured.

## Key Configuration

These are the most important decisions when configuring a security group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Default rule deletion** -- Set `deleteDefaultRules: true` for a zero-trust baseline. OpenStack normally creates two default egress rules (allow all IPv4 and IPv6). Deleting them forces you to explicitly define every allowed traffic flow. This is a create-time setting.

**Stateful vs stateless** -- `stateful` defaults to the OpenStack deployment setting (typically `true`). Stateful groups automatically allow return traffic. Stateless groups require explicit rules for both directions but offer better performance for high-throughput workloads.

**Inline rules** -- Each `rules` entry requires a unique `key` (e.g., `allow-ssh`, `egress-all-ipv4`) used as the IaC resource identifier. Rules support protocol, port ranges, CIDR-based filtering (`remoteIpPrefix`), and security group-based filtering (`remoteGroupId`). The two remote sources are mutually exclusive per rule.

**Rule management strategy** -- For standalone security groups, inline rules keep everything in one manifest. For InfraChart deployments where rules need DAG-level dependency tracking (e.g., cross-referencing security group IDs via ValueFromRef), use the separate OpenStackSecurityGroupRule component instead.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `security_group_id` | UUID of the security group in OpenStack | Security group rule references, network port assignments |
| `name` | Name of the security group | Instance security group assignments (Compute API uses names, not UUIDs) |
| `region` | OpenStack region where the security group was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web server** -- Allows SSH from a trusted CIDR, HTTP and HTTPS from anywhere, with OpenStack's default egress rules kept. The standard pattern for web-facing instances and load balancer targets. Start from the **Web Server** preset.

**Restrictive (zero-trust)** -- Deletes OpenStack's default egress rules and explicitly defines all allowed flows. Suitable for compliance-sensitive environments (PCI-DSS, HIPAA) where every traffic path must be documented and auditable. Start from the **Restrictive** preset.

## Works With

This component operates independently and does not reference other components.