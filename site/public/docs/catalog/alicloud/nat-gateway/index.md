---
title: "NAT Gateway"
description: "NAT Gateway deployment documentation"
icon: "package"
order: 100
componentName: "alicloudnatgateway"
---

# AliCloud NAT Gateway

Deploys an Alibaba Cloud Enhanced NAT Gateway with bundled EIP association and SNAT entries. The component provisions all three resources as a single atomic unit, enabling private VSwitch traffic to reach the internet through a managed NAT service. The NAT Gateway, EIP binding, and SNAT rules are deployed together because a NAT Gateway without an EIP and at least one SNAT entry is non-functional.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **NAT Gateway** -- an `alicloud_nat_gateway` resource placed in the specified VPC and VSwitch with configurable type, billing, and deletion protection
- **EIP Association** -- an `alicloud_eip_association` binding the provided Elastic IP to the NAT Gateway
- **SNAT Entries** -- one `alicloud_snat_entry` per entry in `snatEntries`, mapping private VSwitch or CIDR traffic to the associated EIP for outbound internet access

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC** -- the NAT Gateway must belong to a VPC. Create one with AliCloudVpc.
- **A VSwitch for NAT placement** -- the Enhanced NAT Gateway requires placement in a VSwitch within the VPC. Create one with AliCloudVswitch.
- **An Elastic IP** -- the NAT Gateway needs an EIP for outbound traffic. Create one with AliCloudEipAddress. The EIP must be in the same region.
- **VSwitch(es) for SNAT** -- identify which VSwitches need outbound internet access. Each becomes a SNAT entry.

## Deploy

### Console

Open the deployment store, find **AliCloud NAT Gateway**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, VSwitch, EIP, and SNAT entries.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudNatGateway
metadata:
  name: platform-nat
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  vswitchId:
    value: vsw-nat-zone-a
  natGatewayName: platform-nat
  eipId:
    value: eip-abc123
  snatEntries:
    - sourceVswitchId:
        value: vsw-app-zone-a
      snatEntryName: app-zone-a
```

```shell
planton apply -f alicloud-nat-gateway.yaml
```

This creates an Enhanced NAT Gateway with one SNAT entry enabling the app VSwitch to reach the internet. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a network stack, use ValueFromRef to wire VPC, VSwitch, and EIP dependencies:

```yaml
spec:
  region: cn-hangzhou
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: nat-vswitch
      fieldPath: status.outputs.vswitch_id
  natGatewayName: platform-nat
  eipId:
    valueFrom:
      kind: AliCloudEipAddress
      name: nat-eip
      fieldPath: status.outputs.eip_id
  snatEntries:
    - sourceVswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: app-vswitch-a
          fieldPath: status.outputs.vswitch_id
      snatEntryName: app-zone-a
    - sourceVswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: app-vswitch-b
          fieldPath: status.outputs.vswitch_id
      snatEntryName: app-zone-b
```

The InfraPipeline resolves the dependency graph and provisions VPC, VSwitches, and EIP before the NAT Gateway.

## Key Configuration

These are the most important decisions when configuring a NAT Gateway. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**NAT type** -- The `natType` field defaults to "Enhanced", which supports VSwitch placement and higher performance. "Normal" is legacy and not recommended.

**SNAT source** -- Each SNAT entry accepts either `sourceVswitchId` (NAT an entire VSwitch, the common case) or `sourceCidr` (NAT a specific CIDR for fine-grained control). These are mutually exclusive per entry.

**Billing** -- The `internetChargeType` field defaults to "PayByLcu" (capacity units, scales automatically). "PayBySpec" uses a fixed tier (Small/Middle/Large/XLarge.1) with predictable pricing.

**Deletion protection** -- The `deletionProtection` field prevents accidental deletion. Enable this for production NAT Gateways.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** | `vswitchId` | `status.outputs.vswitch_id` |
| **AliCloudEipAddress** | `eipId` | `status.outputs.eip_id` |
| **AliCloudVswitch** | `snatEntries[].sourceVswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `nat_gateway_id` | NAT Gateway ID (e.g., ngw-xxxxx) | Monitoring, advanced routing |
| `nat_gateway_name` | NAT Gateway name as created | Display and tagging |
| `snat_table_id` | SNAT table ID | Adding SNAT entries outside Planton |
| `forward_table_id` | Forward (DNAT) table ID | Adding DNAT entries outside Planton |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Single VSwitch NAT** -- One NAT Gateway with one SNAT entry for a single VSwitch. The simplest configuration for development or single-zone deployments. Start from the **Single VSwitch** preset.

**Multi-AZ production** -- A NAT Gateway with deletion protection and SNAT entries for multiple VSwitches across Availability Zones. Start from the **Multi AZ Production** preset.

**CIDR-based SNAT** -- Fine-grained SNAT using CIDR blocks instead of whole VSwitches, useful when only a subset of addresses in a VSwitch needs internet access. Start from the **CIDR Based SNAT** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC this NAT Gateway belongs to
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- placement VSwitch for the NAT Gateway and SNAT sources
- [**AliCloud EIP Address**](/cloud-catalog/ali-cloud-eip-address) -- provides the public IP for outbound traffic
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- network security rules for instances using NAT
- [**AliCloud Kubernetes Cluster**](/cloud-catalog/ali-cloud-kubernetes-cluster) -- Kubernetes worker nodes often use NAT for outbound access
