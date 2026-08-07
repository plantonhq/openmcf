# AliCloud Security Group

Deploys an Alibaba Cloud Security Group with bundled ingress and egress rules in a VPC. The component provisions the security group and its traffic rules as a single atomic unit, ensuring the firewall is always created with its intended access policy. Security groups are stateful virtual firewalls that control traffic for ECS instances, ACK worker nodes, RDS instances, and other VPC-aware resources.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Security Group** -- an `alicloud_security_group` resource bound to the specified VPC with configurable inner access policy and tags
- **Security Group Rules** -- one `alicloud_security_group_rule` per entry in `rules`, defining ingress and egress traffic policies

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC** -- the security group must belong to a VPC. Create one with AliCloudVpc or reference an existing VPC via ValueFromRef.
- **Rule planning** -- rules are evaluated by priority (1-100, lower = higher priority). The first matching rule wins. All rule fields except description are immutable after creation.

## Deploy

### Console

Open the deployment store, find **AliCloud Security Group**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, rules, and inner access policy.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudSecurityGroup
metadata:
  name: web-sg
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  securityGroupName: web-tier
  description: Web tier allowing HTTP/HTTPS inbound
  rules:
    - type: ingress
      ipProtocol: tcp
      portRange: "443/443"
      cidrIp: "0.0.0.0/0"
      description: Allow HTTPS from anywhere
    - type: egress
      ipProtocol: all
      portRange: "-1/-1"
      cidrIp: "0.0.0.0/0"
      description: Allow all outbound
```

```shell
planton apply -f alicloud-security-group.yaml
```

This creates a security group with HTTPS inbound and unrestricted outbound. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a network stack, use ValueFromRef to wire the VPC dependency:

```yaml
spec:
  region: cn-hangzhou
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  securityGroupName: web-tier
  rules:
    - type: ingress
      ipProtocol: tcp
      portRange: "443/443"
      cidrIp: "0.0.0.0/0"
      description: Allow HTTPS
```

The InfraPipeline resolves the dependency graph and provisions the VPC before this security group.

## Key Configuration

These are the most important decisions when configuring a security group. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Inner access policy** -- The `innerAccessPolicy` field controls whether instances in the same security group can communicate freely. "Accept" (default) allows all intra-group traffic. "Drop" requires explicit rules for intra-group communication -- use this for database tiers.

**Rule priority** -- Each rule has a `priority` field (1-100). Lower numbers are evaluated first. Use distinct priorities to create predictable rule ordering.

**Protocol and ports** -- The `ipProtocol` and `portRange` fields define what traffic the rule matches. For TCP/UDP, specify explicit ports (e.g., "443/443"). For ICMP, GRE, or ALL, use "-1/-1".

**CIDR vs security group source** -- Rules can match by `cidrIp` (IP range) or `sourceSecurityGroupId` (SG-to-SG). Use SG references for intra-VPC service-to-service rules.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `security_group_id` | Security group ID assigned by Alibaba Cloud | AliCloudEcsInstance, AliCloudKubernetesCluster |
| `security_group_name` | Security group name as created | Display and tagging |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web tier** -- Allows HTTP (80) and HTTPS (443) inbound from the internet with unrestricted outbound. Start from the **Web Tier** preset.

**Database tier** -- Restricts inbound to MySQL (3306), PostgreSQL (5432), and Redis (6379) from VPC-internal CIDR only, with inner access set to Drop. Start from the **Database Tier** preset.

**Bastion host** -- Allows SSH (22) inbound from trusted CIDRs and SSH/database outbound to VPC instances. Start from the **Bastion Host** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC this security group belongs to
- [**AliCloud ECS Instance**](/cloud-catalog/ali-cloud-ecs-instance) -- associate ECS instances with this security group
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- use this security group for ACK worker nodes
