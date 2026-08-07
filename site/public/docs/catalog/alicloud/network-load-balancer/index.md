---
title: "Network Load Balancer"
description: "Network Load Balancer deployment documentation"
icon: "package"
order: 100
componentName: "alicloudnetworkloadbalancer"
---

# AliCloud Network Load Balancer

Deploys an Alibaba Cloud Network Load Balancer (NLB) with bundled server groups and listeners. NLB is a modern Layer 4 load balancer for TCP, UDP, and TCPSSL traffic, designed for ultra-high performance and low latency. The NLB, server groups, and listeners are deployed as a single atomic unit because an NLB without at least one server group and listener is non-functional.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **NLB** -- an `alicloud_nlb_load_balancer` resource placed across at least two Availability Zones for high availability, with configurable address type and cross-zone settings
- **Server Groups** -- one `alicloud_nlb_server_group` per entry in `serverGroups`, each with its own health check, scheduling algorithm, connection draining, and client IP preservation. Groups are created empty -- backend membership is managed externally
- **Listeners** -- one `alicloud_nlb_listener` per entry in `listeners`, each binding a port/protocol to a server group. TCPSSL listeners support TLS termination and mutual TLS

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **An existing VPC** with at least two VSwitches in different Availability Zones -- NLB requires multi-AZ deployment.
- **A server certificate** (for TCPSSL listeners) -- obtain from Alibaba Cloud Certificate Management Service (CAS).
- **EIP addresses** (optional) -- for stable public IPs on internet-facing NLBs. If omitted, Alibaba Cloud auto-assigns IPs.

## Deploy

### Console

Open the deployment store, find **AliCloud Network Load Balancer**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including VPC, zone mappings, server groups, and listeners.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudNetworkLoadBalancer
metadata:
  name: tcp-nlb
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  vpcId:
    value: vpc-bp1234567890
  zoneMappings:
    - zoneId: cn-hangzhou-a
      vswitchId:
        value: vsw-zone-a
    - zoneId: cn-hangzhou-b
      vswitchId:
        value: vsw-zone-b
  serverGroups:
    - name: tcp-backend
      healthCheck:
        healthCheckEnabled: true
  listeners:
    - listenerPort: 80
      listenerProtocol: TCP
      serverGroupName: tcp-backend
```

```shell
planton apply -f alicloud-nlb.yaml
```

This creates an internet-facing NLB with a TCP listener forwarding to the tcp-backend server group. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a network stack, use ValueFromRef to wire VPC, VSwitch, and optional EIP dependencies:

```yaml
spec:
  region: cn-hangzhou
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: platform-vpc
      fieldPath: status.outputs.vpc_id
  zoneMappings:
    - zoneId: cn-hangzhou-a
      vswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: nlb-vswitch-a
          fieldPath: status.outputs.vswitch_id
      allocationId:
        valueFrom:
          kind: AliCloudEipAddress
          name: nlb-eip-a
          fieldPath: status.outputs.eip_id
    - zoneId: cn-hangzhou-b
      vswitchId:
        valueFrom:
          kind: AliCloudVswitch
          name: nlb-vswitch-b
          fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph and provisions VPC, VSwitches, and EIPs before the NLB.

## Key Configuration

These are the most important decisions when configuring an NLB. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Address type** -- The `addressType` field determines whether the NLB is internet-facing ("Internet", default) or VPC-internal ("Intranet").

**Cross-zone balancing** -- The `crossZoneEnabled` field (default true) distributes traffic across all healthy backends in all zones. Disable it to keep traffic within the zone where it was received.

**Connection draining** -- Each server group supports `connectionDrainEnabled` and `connectionDrainTimeout` to allow in-flight connections to complete when backends are removed.

**Client IP preservation** -- The `preserveClientIpEnabled` field (default true) lets backends see the real client IP. Disable for scenarios where this causes asymmetric routing.

**TCPSSL listeners** -- TCPSSL listeners support TLS termination via `certificateIds` and optional mutual TLS via `caEnabled` and `caCertificateIds`.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** | `zoneMappings[].vswitchId` | `status.outputs.vswitch_id` |
| **AliCloudEipAddress** | `zoneMappings[].allocationId` | `status.outputs.eip_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `load_balancer_id` | NLB instance ID (e.g., nlb-xxxxx) | Monitoring, advanced configuration |
| `dns_name` | Auto-assigned DNS name for the NLB | CNAME target for custom domain DNS records |
| `server_group_ids` | Map of server group names to IDs | Backend attachment |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Internet TCP** -- A public-facing NLB with a TCP listener and basic health checks. Start from the **Internet TCP** preset.

**Internal TCP with draining** -- A VPC-internal NLB with connection draining, source IP hashing, and Proxy Protocol for backend IP visibility. Start from the **Internal TCP Drain** preset.

**TCPSSL production** -- An internet-facing NLB with stable EIPs, TLS termination, least-connections scheduling, and HTTP health checks. Start from the **TCPSSL Production** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- the VPC this NLB belongs to
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides zone-specific IP allocation for NLB nodes
- [**AliCloud EIP Address**](/cloud-catalog/ali-cloud-eip-address) -- stable public IPs for internet-facing NLBs
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- network security rules for backend instances
- [**AliCloud DNS Record**](/cloud-catalog/ali-cloud-dns-record) -- DNS records pointing to the NLB's dns_name
