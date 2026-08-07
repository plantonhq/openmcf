# OpenStack Security Group Rule

Deploys a standalone Neutron security group rule on OpenStack, adding a single firewall rule to an existing security group. This is the DAG-visible counterpart to inline rules defined in OpenStackSecurityGroup -- use this component when individual rules need to be independently managed or when cross-security-group references need InfraChart wiring via ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Neutron Security Group Rule** -- a single firewall rule with configurable direction, ethertype, protocol, port range, and remote source (CIDR or security group reference)

## Before You Deploy

### OpenStack Account

- **Security group** -- an existing Neutron security group to add the rule to. Provide the security group ID directly or reference an OpenStackSecurityGroup Cloud Resource via ValueFromRef.
- **Remote security group** (optional) -- if restricting traffic by membership in another security group (instead of CIDR), have the remote security group ID ready or reference it via ValueFromRef.
- **Rule uniqueness** -- OpenStack rejects duplicate rules (same direction, protocol, port range, ethertype, and remote source) within the same security group.

## Deploy

### Console

Open the deployment store, find **OpenStack Security Group Rule**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Allow SSH** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: openstack.planton.dev/v1
kind: OpenStackSecurityGroupRule
metadata:
  name: allow-https
  org: acme-corp
  env: prod
spec:
  securityGroupId:
    value: "<security-group-id>"
  direction: ingress
  ethertype: IPv4
  protocol: tcp
  portRangeMin: 443
  portRangeMax: 443
  remoteIpPrefix: "0.0.0.0/0"
```

```shell
planton apply -f security-group-rule.yaml
```

This creates an ingress rule allowing HTTPS traffic from any IPv4 source. No description, remote security group reference, or region override is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the rule to a security group deployed in the same InfraPipeline:

```yaml
spec:
  securityGroupId:
    valueFrom:
      kind: OpenStackSecurityGroup
      name: web-sg
      fieldPath: status.outputs.security_group_id
  remoteGroupId:
    valueFrom:
      kind: OpenStackSecurityGroup
      name: app-sg
      fieldPath: status.outputs.security_group_id
```

The InfraPipeline resolves the dependency graph, deploys both security groups first, then provisions the rule with the resolved IDs.

## Key Configuration

These are the most important decisions when configuring a security group rule. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Direction and ethertype** -- `direction` must be `ingress` or `egress`. `ethertype` must be `IPv4` or `IPv6`. These two fields define the fundamental scope of the rule.

**Protocol and port range** -- `protocol` accepts `tcp`, `udp`, `icmp`, `icmpv6`, or any IANA protocol name/number. If omitted, the rule applies to all protocols. Port ranges require a protocol -- for ICMP, `portRangeMin` is the ICMP type and `portRangeMax` is the ICMP code.

**Remote source** -- `remoteIpPrefix` restricts traffic to a CIDR (e.g., `0.0.0.0/0` for all, `10.0.0.0/8` for private ranges). `remoteGroupId` restricts traffic to instances in another security group. The two are mutually exclusive. Security group references are the key advantage of standalone rules over inline rules, enabling InfraChart DAG wiring.

**Immutability** -- All fields are ForceNew. Changing any field (direction, protocol, port range, remote source) recreates the rule.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OpenStackSecurityGroup** | `securityGroupId` | `status.outputs.security_group_id` |
| **OpenStackSecurityGroup** (optional) | `remoteGroupId` | `status.outputs.security_group_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `rule_id` | UUID of the security group rule | Import, audit, and debugging |
| `security_group_id` | UUID of the security group this rule belongs to | Audit, topology reference |
| `direction` | Direction of the rule (ingress or egress) | Monitoring labels, automation |
| `protocol` | IP protocol of the rule | Monitoring labels, automation |
| `port_range_min` | Lower bound of the port range (or ICMP type) | Audit, documentation |
| `port_range_max` | Upper bound of the port range (or ICMP code) | Audit, documentation |
| `region` | OpenStack region where the rule was created | Region-aware downstream resources |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Allow SSH** -- Permits inbound SSH (TCP port 22) from a trusted CIDR. The standard pattern for administrative access to instances. Restrict `remoteIpPrefix` to your VPN or bastion CIDR instead of `0.0.0.0/0`. Start from the **Allow SSH** preset.

**Allow HTTPS** -- Permits inbound HTTPS (TCP port 443) from any IPv4 source. The standard pattern for web-facing services, load balancers, and reverse proxies. Start from the **Allow HTTPS** preset.

## Works With

- [**OpenStack Security Group**](/cloud-catalog/openstack-security-group) -- provides the security group ID that this rule is added to, and optionally the remote security group ID for group-based traffic filtering