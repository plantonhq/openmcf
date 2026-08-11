# Scaleway Instance Security Group

Deploys a stateful firewall on Scaleway that controls inbound and outbound traffic to Instances. Security groups use ordered rules with accept/drop actions and configurable default policies, supporting both allowlist (deny-all + explicit accepts) and denylist (allow-all + explicit drops) models. Assigned to individual Instances via the Instance's `security_group_id` field.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Instance Security Group** -- a `scaleway_instance_security_group` in the specified zone with the configured default policies, statefulness, and SMTP security settings
- **Inbound Rules** -- created only when `inboundRules` entries are provided; ordered firewall rules controlling traffic TO instances (accept or drop by protocol, port range, and source IP)
- **Outbound Rules** -- created only when `outboundRules` entries are provided; ordered firewall rules controlling traffic FROM instances (accept or drop by protocol, port range, and destination IP)
- **Scaleway Tags** -- resource metadata tags (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Scaleway Account

- **A Scaleway account** with an active project and API access key pair (Access Key + Secret Key). The IaC module authenticates through the Scaleway provider configuration.
- **Choose an Availability Zone** -- security groups are zonal resources (`fr-par-1`, `nl-ams-1`, `pl-waw-1`). The zone must match the zone of the Instances that will use this security group. Cannot be changed after creation.
- **Plan your security model** -- decide between an allowlist model (`inboundDefaultPolicy: drop` with explicit accept rules) or a denylist model (`inboundDefaultPolicy: accept` with explicit drop rules) before defining rules. The allowlist model is recommended for production.

## Deploy

### Console

Open the deployment store, find **Scaleway Instance Security Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web Server** preset in the [Presets](#presets) tab for a public-facing web server with SSH, HTTP, and HTTPS rules.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: scaleway.planton.dev/v1
kind: ScalewayInstanceSecurityGroup
metadata:
  name: web-firewall
  org: acme-corp
  env: prod
spec:
  zone: fr-par-1
  stateful: true
  inboundDefaultPolicy: drop
  outboundDefaultPolicy: accept
  inboundRules:
    - action: accept
      protocol: TCP
      portRange: "443"
      ipRange: "0.0.0.0/0"
    - action: accept
      protocol: TCP
      portRange: "80"
      ipRange: "0.0.0.0/0"
```

```shell
planton apply -f scaleway-instance-security-group.yaml
```

This creates an allowlist security group that drops all inbound traffic except HTTP and HTTPS from any source. All outbound traffic is allowed. SSH is not configured -- add an inbound rule with a restricted `ipRange` for admin access. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a security group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Default policies** -- The `inboundDefaultPolicy` and `outboundDefaultPolicy` fields control what happens to traffic matching no rule. Set `inboundDefaultPolicy` to `drop` for production (allowlist model -- only explicitly accepted traffic is permitted). Keep `outboundDefaultPolicy` as `accept` unless you need strict egress control for compliance.

**Statefulness** -- The `stateful` field (default true) automatically permits return traffic for accepted connections. For example, accepting inbound TCP port 80 also allows the response packets outbound. Disable only for advanced stateless routing or network appliances.

**Rule ordering** -- Inbound and outbound rules are evaluated in order. The first matching rule wins. Place more specific rules (single IPs, narrow port ranges) before broader rules to ensure correct traffic classification.

**SMTP security** -- The `enableDefaultSecurity` field (default true) blocks outbound SMTP on ports 25, 465, and 587 to prevent spam abuse. Disable only if the Instance needs to send email directly and the Scaleway account is authorized for SMTP.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `security_group_id` | Zoned ID of the security group (`{zone}/{uuid}`) | Instance `security_group_id` field for firewall assignment |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web server firewall** -- An allowlist security group accepting SSH from a restricted admin CIDR plus HTTP/HTTPS from anywhere. All outbound traffic is allowed. The standard pattern for public-facing web servers, API servers, and reverse proxies. Start from the **Web Server** preset.

**Deny-all allowlist** -- A strict security group accepting only TCP from the private network range (10.0.0.0/8) and SSH from a single bastion host. The standard pattern for databases, caches, message queues, and backend workers that should never be directly internet-accessible. Start from the **Deny-All Allowlist** preset.

## Works With

This component operates independently and does not reference other components.