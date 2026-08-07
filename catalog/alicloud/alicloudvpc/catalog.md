# AliCloud VPC

Deploys an Alibaba Cloud Virtual Private Cloud with a configurable IPv4 CIDR block, optional IPv6 dual-stack support, resource group assignment, and automatic tag management. The VPC is the networking foundation for VSwitches, security groups, NAT gateways, load balancers, and Kubernetes clusters on Alibaba Cloud.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPC** -- an `alicloud_vpc` resource with the specified CIDR block, name, region, and optional IPv6 configuration
- **VRouter** -- automatically created by Alibaba Cloud as part of VPC creation, responsible for routing traffic between VSwitches
- **System Route Table** -- the default route table associated with the VRouter, containing system routes for intra-VPC communication

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **CIDR block planning** -- the primary IPv4 CIDR cannot be changed after creation. Choose a range that accommodates future growth and avoids overlap with other VPCs if using VPC peering or CEN.
- **Region selection** -- choose the region closest to your workloads (e.g., cn-hangzhou, cn-shanghai, us-west-1, ap-southeast-1).

## Deploy

### Console

Open the deployment store, find **AliCloud VPC**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including region, CIDR block, and IPv6 settings.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudVpc
metadata:
  name: platform-vpc
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcName: platform-vpc
  cidrBlock: "10.0.0.0/16"
  description: Production VPC for platform workloads
  tags:
    team: platform
```

```shell
planton apply -f alicloud-vpc.yaml
```

This creates a VPC with a /16 CIDR block in cn-hangzhou. A Stack Job tracks the provisioning in real time.

### InfraChart

No upstream dependencies are required for a VPC. Downstream components reference the VPC via ValueFromRef:

```yaml
spec:
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
```

The InfraPipeline resolves the dependency graph and provisions the VPC before any dependent resources.

## Key Configuration

These are the most important decisions when configuring a VPC. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CIDR block** -- The `cidrBlock` field defines the private address space. Must be in 10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16 with a mask length of 8-28. This is immutable after creation -- plan for growth.

**IPv6 dual-stack** -- The `enableIpv6` field allocates a /56 IPv6 CIDR block to the VPC. VSwitches can then be assigned IPv6 subnets. Enable at creation time if you anticipate IPv6 workloads.

**Resource group** -- The `resourceGroupId` field places the VPC in an Alibaba Cloud resource group for access control and cost attribution. If omitted, the account's default resource group is used.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `vpc_id` | VPC ID assigned by Alibaba Cloud | AliCloudVswitch, AliCloudSecurityGroup, AliCloudNatGateway, AliCloudKubernetesCluster |
| `vpc_name` | VPC name as created | Display and tagging |
| `cidr_block` | Primary IPv4 CIDR block | VSwitch CIDR planning, security group rules |
| `router_id` | VRouter ID auto-created with the VPC | Advanced routing configuration |
| `route_table_id` | System route table ID | Custom route entries |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard production** -- A /16 VPC with organizational tags. Start from the **Standard Production** preset.

**Development** -- A minimal /16 VPC using the 192.168.0.0 range for isolated dev environments. Start from the **Development** preset.

**Dual-stack IPv6** -- A VPC with IPv6 enabled for modern applications that need native IPv6 connectivity. Start from the **Dual Stack IPv6** preset.

## Works With

- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- creates subnets within this VPC
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- creates security groups bound to this VPC
- [**AliCloud NAT Gateway**](/cloud-catalog/ali-cloud-nat-gateway) -- provides outbound internet access for private VSwitches
- [**AliCloud EIP Address**](/cloud-catalog/ali-cloud-eip-address) -- standalone public IPs for NAT gateways and load balancers
- [**AliCloud Application Load Balancer**](/cloud-catalog/ali-cloud-application-load-balancer) -- deploys ALBs in this VPC
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- deploys managed Kubernetes clusters in this VPC
