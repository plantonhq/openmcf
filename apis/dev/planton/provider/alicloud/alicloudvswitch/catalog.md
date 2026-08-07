# AliCloud VSwitch

Deploys an Alibaba Cloud VSwitch (subnet) within an existing VPC, bound to a single Availability Zone with a dedicated IPv4 CIDR block, optional IPv6 dual-stack support, and automatic tag management. The VSwitch is the mandatory network placement target for ECS instances, databases, Kubernetes clusters, NAT gateways, and load balancers on Alibaba Cloud.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VSwitch** -- an `alicloud_vswitch` resource with the specified VPC, Availability Zone, CIDR block, name, and optional IPv6 configuration

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC** -- the VSwitch's `vpcId` must reference a valid AliCloudVpc. Create one first or reference an existing VPC via ValueFromRef.
- **CIDR block planning** -- the VSwitch CIDR must be a subset of the parent VPC's CIDR block with a mask length of 16-29. Cannot overlap with other VSwitches in the same VPC.
- **Availability Zone selection** -- each VSwitch is permanently bound to one zone (e.g., cn-hangzhou-a). Resources deployed into this VSwitch run in that zone.

## Deploy

### Console

Open the deployment store, find **AliCloud VSwitch**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, zone, and CIDR block.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudVswitch
metadata:
  name: app-vswitch-a
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  zoneId: cn-hangzhou-a
  cidrBlock: "10.0.0.0/24"
  vswitchName: app-tier-zone-a
  description: Application tier VSwitch in zone A
```

```shell
planton apply -f alicloud-vswitch.yaml
```

This creates a /24 VSwitch in cn-hangzhou-a within the specified VPC. A Stack Job tracks the provisioning in real time.

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
  zoneId: cn-hangzhou-a
  cidrBlock: "10.0.0.0/24"
  vswitchName: app-tier-zone-a
```

The InfraPipeline resolves the dependency graph and provisions the VPC before this VSwitch.

## Key Configuration

These are the most important decisions when configuring a VSwitch. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**VPC binding** -- The `vpcId` field links this VSwitch to its parent VPC. Use a ValueFromRef to establish a declarative dependency. Changing this after creation forces replacement.

**Availability Zone** -- The `zoneId` field determines where resources in this VSwitch physically run. Choose based on latency requirements and HA strategy. Changing this forces replacement.

**CIDR block** -- The `cidrBlock` field carves out an address range from the parent VPC. A /24 gives 256 addresses, a /20 gives 4,096. This is immutable after creation.

**IPv6** -- The `enableIpv6` field allocates a /64 IPv6 segment to this VSwitch. The parent VPC must have IPv6 enabled first.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vswitch_id` | VSwitch ID assigned by Alibaba Cloud | AliCloudNatGateway, AliCloudEcsInstance, AliCloudKubernetesCluster, AliCloudRdsInstance |
| `vswitch_name` | VSwitch name as created | Display and tagging |
| `cidr_block` | IPv4 CIDR block of the VSwitch | Security group rules, CIDR planning |
| `zone_id` | Availability Zone of the VSwitch | Zone-aware resource placement |
| `ipv6_cidr_block` | IPv6 CIDR block (when IPv6 is enabled) | IPv6 resource configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development single-zone** -- A /24 VSwitch in one zone for non-production workloads. Start from the **Dev Single Zone** preset.

**Production application tier** -- A /20 VSwitch with a large address space for Kubernetes node pools or ECS auto-scaling groups. Start from the **Prod App Tier** preset.

**IPv6 dual-stack** -- A VSwitch with IPv6 enabled for workloads requiring native IPv6. The parent VPC must have IPv6 enabled. Start from the **IPv6 Enabled** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the parent VPC this VSwitch belongs to
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- controls traffic for resources in this VSwitch
- [**AliCloud NAT Gateway**](/cloud-catalog/ali-cloud-nat-gateway) -- provides outbound internet access for private VSwitches
- [**AliCloud ECS Instance**](/cloud-catalog/ali-cloud-ecs-instance) -- compute instances launched in this VSwitch
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- managed Kubernetes clusters using this VSwitch for node placement
- [**AliCloud RDS Instance**](/cloud-catalog/ali-cloud-rds-instance) -- database instances placed in this VSwitch
