# AliCloudSecurityGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudSecurityGroupSpec defines the configuration for an Alibaba Cloud
Security Group with bundled security rules.

A security group is a stateful virtual firewall that controls inbound and
outbound traffic for ECS instances, ACK worker nodes, RDS instances, and
other VPC-aware resources. Each security group belongs to exactly one VPC,
and resources can be associated with up to five security groups.

This component bundles the security group with its rules (per DD07 composite
bundling) because a security group without rules is an open door -- the rules
are the substance of the resource.

Security group rules are evaluated by priority (lowest number = highest
priority, range 1-100). The first matching rule determines whether traffic
is allowed or dropped. If no rule matches, the default behavior depends on
the inner_access_policy for intra-group traffic, and deny for everything else.

Provider resources:
  Terraform: alicloud_security_group + alicloud_security_group_rule
  Pulumi:    ecs.SecurityGroup + ecs.SecurityGroupRule

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudSecurityGroup
metadata:
  name: alicloudsecuritygroup-demo
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-demo123
  securityGroupName: demo-sg
  description: Demo security group for local testing
  rules:
    - type: ingress
      ipProtocol: tcp
      portRange: "443/443"
      cidrIp: "0.0.0.0/0"
      description: Allow HTTPS
    - type: egress
      ipProtocol: all
      portRange: "-1/-1"
      cidrIp: "0.0.0.0/0"
      description: Allow all outbound
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.securityGroupName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.innerAccessPolicy` | `string` |  | `Accept` |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.rules` | `[]AliCloudSecurityGroupRule` |  |  |  |
| `spec.rules[].type` | `string` | yes |  |  |
| `spec.rules[].ipProtocol` | `string` | yes |  |  |
| `spec.rules[].portRange` | `string` |  | `-1/-1` |  |
| `spec.rules[].cidrIp` | `string` |  |  |  |
| `spec.rules[].sourceSecurityGroupId` | `string` |  |  |  |
| `spec.rules[].priority` | `int32` |  | `1` |  |
| `spec.rules[].policy` | `string` |  | `accept` |  |
| `spec.rules[].description` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the security group will be created.
Must match the region of the VPC referenced by vpc_id.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.vpcId

`string | valueFrom` · required

VPC ID that this security group belongs to.
All security group rules automatically apply to VPC-internal traffic
(nic_type "intranet" is implied and hardcoded in the IaC modules).

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.securityGroupName

`string` · required

Security group name. 2-128 characters; must start with a letter; can
contain Unicode characters, digits, colons, underscores, periods, hyphens.
Maps to the provider field `security_group_name`.

- rule: {"required":true,"string":{"minLen":"2","maxLen":"128"}}

### spec.description

`string`

Human-readable description of the security group's purpose.
2-256 characters; cannot start with http:// or https://.

### spec.innerAccessPolicy

`string` · optional (explicit presence)

Controls whether instances within the same security group can communicate
with each other freely. When "Accept", all intra-group traffic is allowed.
When "Drop", intra-group traffic follows normal rule evaluation.

Uses the provider's exact casing: "Accept" or "Drop".
Default: "Accept"

- default: `Accept`
- rule: inner_access_policy must be one of: Accept, Drop

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the security group is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the security group.

### spec.rules

`[]AliCloudSecurityGroupRule`

Security rules that define allowed or dropped traffic flows.
Each rule specifies direction (ingress/egress), protocol, port range,
source/destination, and an accept/drop decision.

Rules are evaluated by priority within their direction. Lower priority
numbers are evaluated first (1 = highest priority). The first matching
rule wins.

### spec.rules[].type

`string` · required

Traffic direction: "ingress" for inbound, "egress" for outbound.

- rule: type must be one of: ingress, egress
- rule: {"required":true}

### spec.rules[].ipProtocol

`string` · required

IP protocol for this rule.
For "tcp" and "udp", port_range must specify explicit ports.
For "icmp", "gre", and "all", port_range must be "-1/-1".

- rule: ip_protocol must be one of: tcp, udp, icmp, gre, all
- rule: {"required":true}

### spec.rules[].portRange

`string` · optional (explicit presence)

Port range in "start/end" format.
For tcp/udp: valid range 1-65535, e.g. "80/80" (single port), "8080/8090" (range).
For icmp/gre/all: must be "-1/-1".
Default: "-1/-1" (all ports)

- default: `-1/-1`

### spec.rules[].cidrIp

`string`

IPv4 CIDR block for the traffic source (ingress) or destination (egress).
Example: "0.0.0.0/0" (any), "10.0.0.0/8" (private range).
At least one of cidr_ip or source_security_group_id must be specified.

### spec.rules[].sourceSecurityGroupId

`string`

Source (for ingress) or destination (for egress) security group ID.
Used for SG-to-SG rules within the same VPC.
Mutually exclusive with cidr_ip in the provider -- when both are set,
cidr_ip takes precedence.

### spec.rules[].priority

`int32` · optional (explicit presence)

Rule evaluation priority. Lower number = higher priority.
Range: 1 to 100.
Default: 1 (highest priority)

- default: `1`
- rule: {"int32":{"lte":100,"gte":1}}

### spec.rules[].policy

`string` · optional (explicit presence)

Action to take when this rule matches: "accept" to allow, "drop" to block.
Default: "accept"

- default: `accept`
- rule: policy must be one of: accept, drop

### spec.rules[].description

`string`

Human-readable description of this rule's purpose.
This is the only field that can be updated after creation without
recreating the rule.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudSecurityGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.security_group_id` | `string` | The security group ID assigned by Alibaba Cloud. Referenced by downstream components (EcsInstance, AckManagedCluster, etc.) via StringValueOrRef. |
| `status.outputs.security_group_name` | `string` | The security group name as created. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudEcsInstance | `spec.securityGroupIds` | `status.outputs.security_group_id` |
| AliCloudFunction | `spec.vpcConfig.securityGroupId` | `status.outputs.security_group_id` |
| AliCloudKubernetesCluster | `spec.securityGroupId` | `status.outputs.security_group_id` |
| AliCloudKubernetesNodePool | `spec.securityGroupIds` | `status.outputs.security_group_id` |
| AliCloudSaeApplication | `spec.securityGroupId` | `status.outputs.security_group_id` |

## See Also

- [Overview](../README.md)
