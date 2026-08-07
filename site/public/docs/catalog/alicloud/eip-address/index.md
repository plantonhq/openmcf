---
title: "EIP Address"
description: "EIP Address deployment documentation"
icon: "package"
order: 100
componentName: "alicloudeipaddress"
---

# AliCloud EIP Address

Deploys an Alibaba Cloud Elastic IP Address (EIP) -- a static, public IPv4 address that persists independently of any cloud resource. An EIP can be associated with NAT gateways, load balancers, VPN gateways, and ECS instances, and can be released and re-associated without changing the address.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **EIP** -- an `alicloud_eip_address` resource in the specified region with configurable bandwidth, ISP line type, and metering settings

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **Region selection** -- the EIP can only be associated with resources in the same region.
- **ISP and metering decisions** -- both `isp` and `internetChargeType` are immutable after creation. Choose carefully based on workload characteristics.

## Deploy

### Console

Open the deployment store, find **AliCloud EIP Address**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including region, bandwidth, and ISP.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudEipAddress
metadata:
  name: nat-eip
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  addressName: prod-nat-eip
  description: EIP for production NAT gateway
  bandwidth: 10
  tags:
    purpose: nat
    team: platform
```

```shell
planton apply -f alicloud-eip.yaml
```

This allocates a 10 Mbps EIP using BGP multi-line and PayByTraffic metering. A Stack Job tracks the provisioning in real time.

### InfraChart

EIPs are standalone resources with no upstream dependencies. Downstream components reference the EIP via ValueFromRef:

```yaml
spec:
  eipId:
    valueFrom:
      kind: AliCloudEipAddress
      name: nat-eip
      fieldPath: status.outputs.eip_id
```

The InfraPipeline resolves the dependency graph and allocates the EIP before any dependent resources.

## Key Configuration

These are the most important decisions when configuring an EIP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Bandwidth** -- The `bandwidth` field sets the maximum outbound bandwidth in Mbps (1-1000, default 5). With PayByTraffic, this is a ceiling. With PayByBandwidth, this is the reserved allocation.

**Metering method** -- The `internetChargeType` field is immutable after creation. "PayByTraffic" (default) bills per GB transferred -- best for bursty workloads. "PayByBandwidth" bills for reserved bandwidth -- best for steady, high-throughput workloads.

**ISP line type** -- The `isp` field selects the network routing. "BGP" (default) works in all regions. "BGP_PRO" provides optimized routing for China mainland. Single-carrier options (ChinaTelecom, ChinaUnicom, ChinaMobile) are available for specific needs. Immutable after creation.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `eip_id` | EIP allocation ID assigned by Alibaba Cloud | AliCloudNatGateway, AliCloudApplicationLoadBalancer, AliCloudVpnGateway |
| `ip_address` | Allocated public IPv4 address | DNS record configuration, application endpoints |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard EIP** -- A default 5 Mbps BGP EIP with PayByTraffic metering, suitable for NAT gateways and non-production load balancers. Start from the **Standard** preset.

**High-bandwidth production** -- A 100 Mbps BGP_PRO EIP with PayByBandwidth metering for production ALBs with predictable traffic. Start from the **High Bandwidth** preset.

## Works With

- [**AliCloud NAT Gateway**](/cloud-catalog/ali-cloud-nat-gateway) -- associate this EIP for SNAT outbound internet access
- [**AliCloud Application Load Balancer**](/cloud-catalog/ali-cloud-application-load-balancer) -- use this EIP for internet-facing ALB
- [**AliCloud Network Load Balancer**](/cloud-catalog/ali-cloud-network-load-balancer) -- use this EIP for internet-facing NLB
- [**AliCloud VPN Gateway**](/cloud-catalog/ali-cloud-vpn-gateway) -- use this EIP for VPN gateway public endpoint
