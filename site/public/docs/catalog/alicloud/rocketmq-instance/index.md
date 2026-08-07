---
title: "RocketMQ Instance"
description: "RocketMQ Instance deployment documentation"
icon: "package"
order: 100
componentName: "alicloudrocketmqinstance"
---

# AliCloud RocketMQ Instance

Deploys an Alibaba Cloud RocketMQ 5.x instance with bundled topics and consumer groups. The component provisions the instance, its topics, and its consumer groups as a single atomic unit. RocketMQ 5.x provides VPC-integrated instances with configurable edition tiers, billing modes, internet access, message tracing, encryption at rest, and four message types (Normal, FIFO, Delay, Transaction). The component integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to VPCs and VSwitches.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **RocketMQ Instance** -- an `alicloud_rocketmq_instance` with the selected edition series, deployment architecture, VPC networking, billing mode, and optional internet access, throughput specification, and encryption
- **Topics** -- one `alicloud_rocketmq_topic` per entry in `topics`, each with a configured message type (NORMAL, FIFO, DELAY, or TRANSACTION)
- **Consumer Groups** -- one `alicloud_rocketmq_consumer_group` per entry in `consumerGroups`, each with delivery ordering, retry policy, and optional dead-letter topic configuration
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- A VPC for the instance's network endpoint. RocketMQ 5.x instances require VPC placement. Provide the VPC ID directly or reference an AliCloudVpc Cloud Resource via ValueFromRef.
- A VSwitch (optional but recommended) to control the instance's availability zone placement within the VPC. For serverless instances, at least two VSwitches across availability zones are recommended.
- A KMS key ID (optional) if enabling encryption at rest via `productInfo.storageEncryption`.

## Deploy

### Console

Open the deployment store, find **AliCloud RocketMQ Instance**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Production HA** preset in the [Presets](#presets) tab to pre-populate a professional-tier cluster with HA, message tracing, and sample topics.

### CLI

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRocketmqInstance
metadata:
  name: app-messaging
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  seriesCode: professional
  subSeriesCode: cluster_ha
  vpcId:
    value: "vpc-abc123"
  topics:
    - topicName: order-events
      messageType: NORMAL
    - topicName: payment-events
      messageType: FIFO
  consumerGroups:
    - consumerGroupId: order-processor
      deliveryOrderType: Concurrently
    - consumerGroupId: payment-processor
      deliveryOrderType: Orderly
```

```shell
planton apply -f rocketmq.yaml
```

This creates a professional-tier HA RocketMQ instance with two topics and two consumer groups. Billing defaults to PayAsYouGo. Internet access is disabled, so the instance is only accessible within the VPC.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the instance to a VPC and VSwitch deployed in the same InfraPipeline:

```yaml
spec:
  vpcId:
    valueFrom:
      kind: AliCloudVpc
      name: messaging-vpc
      fieldPath: status.outputs.vpc_id
  vswitchId:
    valueFrom:
      kind: AliCloudVswitch
      name: messaging-vswitch
      fieldPath: status.outputs.vswitch_id
```

The InfraPipeline resolves the dependency graph, deploys the VPC and VSwitch first, then provisions the RocketMQ instance with the resolved values.

## Key Configuration

These are the most important decisions when configuring a RocketMQ instance. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Edition and architecture** -- Set `seriesCode` to `standard` (dev/test), `professional` (production), or `ultimate` (mission-critical) and `subSeriesCode` to `cluster_ha` (multi-node HA), `single_node` (development), or `serverless` (auto-scaling). Both fields are immutable after creation -- edition changes require instance replacement.

**Message types** -- Each topic's `messageType` is immutable after creation: `NORMAL` (no ordering guarantee), `FIFO` (strict ordering per message group), `DELAY` (scheduled delivery), `TRANSACTION` (two-phase commit). Match the topic type to the consumer group's `deliveryOrderType` -- `Orderly` for FIFO topics, `Concurrently` for NORMAL topics.

**Billing and subscription** -- Set `paymentType` to `PayAsYouGo` (default, post-paid) or `Subscription` (prepaid with `period` and `periodUnit`). Subscription offers cost savings for stable workloads. Enable `autoRenew` to prevent accidental service interruption on expiry.

**Internet access** -- Configure `internetInfo.enabled: true` to expose a public endpoint alongside the VPC endpoint. Set `flowOutType` to `payByBandwidth` (fixed, set `flowOutBandwidth` in Mb/s) or `payByTraffic` (usage-based, default). Internet configuration is immutable after creation.

**Retry and dead-letter** -- Each consumer group's `consumeRetryPolicy` controls failure handling. `DefaultRetryPolicy` uses exponential backoff (recommended). `FixedRetryPolicy` uses a constant interval. Set `maxRetryTimes` (0-1000) and `deadLetterTargetTopic` to route exhausted messages for manual investigation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudVpc** | `vpcId` | `status.outputs.vpc_id` |
| **AliCloudVswitch** (optional) | `vswitchId` | `status.outputs.vswitch_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `instance_id` | The RocketMQ instance ID assigned by Alibaba Cloud | Monitoring dashboards, ACL configuration |
| `tcp_endpoint` | TCP endpoint for VPC-internal access | Producer and consumer client configuration within the VPC |
| `internet_endpoint` | TCP endpoint for public internet access (empty when disabled) | External producer/consumer access, cross-region replication |
| `topic_ids` | Map of topic names to their resource IDs | Topic management, monitoring per-topic metrics |
| `consumer_group_ids` | Map of consumer group IDs to their resource IDs | Consumer group management, per-group monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Development single node** -- A standard-tier single-node instance for development and testing. Minimal configuration with only a VPC reference. Start from the **Development Single Node** preset.

**Production HA** -- A professional-tier cluster with HA, message tracing, 7-day message retention, IP whitelists, and sample NORMAL and FIFO topics with matching consumer groups. Start from the **Production HA** preset.

**Enterprise encrypted** -- An ultimate-tier cluster with Subscription billing, 14-day message retention, auto-scaling, encryption at rest via KMS, internet access, and TRANSACTION and DELAY topics with dead-letter routing. Start from the **Enterprise Encrypted** preset.

## Works With

- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- provides the VPC for the instance's network endpoint
- [**AliCloud VSwitch**](/cloud-catalog/ali-cloud-vswitch) -- provides availability zone placement within the VPC