# AliCloudRocketmqInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudRocketmqInstanceSpec defines the configuration for an Alibaba Cloud
RocketMQ 5.x instance with bundled topics and consumer groups.

RocketMQ is Alibaba Cloud's distributed messaging and streaming platform,
supporting normal, FIFO, delayed, and transactional messages. This component
targets the RocketMQ 5.x API (2022-08-01), which provides VPC-integrated
instances with configurable throughput tiers, billing modes, and internet
access -- a significant upgrade over the legacy ONS API.

The instance edition is controlled by series_code (standard, professional,
ultimate) and sub_series_code (cluster_ha, single_node, serverless). These
are ForceNew in the provider, meaning edition changes require instance
replacement.

This component bundles topics and consumer groups (per DD07 composite
bundling) because they are meaningless without a parent instance. ACL
accounts and permission rules are intentionally excluded -- security
configuration has an independent lifecycle.

The provider's deeply nested network_info block is partially flattened for
better YAML authoring experience: vpc_id and vswitch_id are promoted to the
spec root (consistent with all other AliCloud components), while internet
access settings remain in an optional nested message since their fields are
conditionally relevant.

Provider resources:
  Terraform: alicloud_rocketmq_instance + alicloud_rocketmq_topic + alicloud_rocketmq_consumer_group
  Pulumi:    rocketmq.RocketMQInstance + rocketmq.RocketMQTopic + rocketmq.ConsumerGroup

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudRocketmqInstance
metadata:
  name: alicloudrocketmqinstance-demo
spec:
  region: cn-hangzhou
  seriesCode: standard
  subSeriesCode: single_node
  vpcId:
    value: vpc-demo123
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.seriesCode` | `string` | yes |  |  |
| `spec.subSeriesCode` | `string` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` | yes |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.instanceName` | `string` |  |  |  |
| `spec.remark` | `string` |  |  |  |
| `spec.paymentType` | `string` |  | `PayAsYouGo` |  |
| `spec.period` | `int32` |  |  |  |
| `spec.periodUnit` | `string` |  |  |  |
| `spec.autoRenew` | `bool` |  |  |  |
| `spec.autoRenewPeriod` | `int32` |  |  |  |
| `spec.vswitchId` | `string \| valueFrom` |  |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.securityGroupId` | `string` |  |  |  |
| `spec.internetInfo` | `AliCloudRocketmqInternetInfo` |  |  |  |
| `spec.internetInfo.enabled` | `bool` |  |  |  |
| `spec.internetInfo.flowOutType` | `string` |  | `payByTraffic` |  |
| `spec.internetInfo.flowOutBandwidth` | `int32` |  |  |  |
| `spec.msgProcessSpec` | `string` |  |  |  |
| `spec.productInfo` | `AliCloudRocketmqProductInfo` |  |  |  |
| `spec.productInfo.messageRetentionTime` | `int32` |  |  |  |
| `spec.productInfo.sendReceiveRatio` | `double` |  |  |  |
| `spec.productInfo.autoScaling` | `bool` |  |  |  |
| `spec.productInfo.traceOn` | `bool` |  |  |  |
| `spec.productInfo.storageEncryption` | `bool` |  |  |  |
| `spec.productInfo.storageSecretKey` | `string` |  |  |  |
| `spec.ipWhitelists` | `[]string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.topics` | `[]AliCloudRocketmqTopic` |  |  |  |
| `spec.topics[].topicName` | `string` | yes |  |  |
| `spec.topics[].messageType` | `string` |  | `NORMAL` |  |
| `spec.topics[].remark` | `string` |  |  |  |
| `spec.topics[].maxSendTps` | `int32` |  |  |  |
| `spec.consumerGroups` | `[]AliCloudRocketmqConsumerGroup` |  |  |  |
| `spec.consumerGroups[].consumerGroupId` | `string` | yes |  |  |
| `spec.consumerGroups[].deliveryOrderType` | `string` |  |  |  |
| `spec.consumerGroups[].remark` | `string` |  |  |  |
| `spec.consumerGroups[].maxReceiveTps` | `int32` |  |  |  |
| `spec.consumerGroups[].consumeRetryPolicy` | `AliCloudRocketmqConsumeRetryPolicy` |  |  |  |
| `spec.consumerGroups[].consumeRetryPolicy.retryPolicy` | `string` |  | `DefaultRetryPolicy` |  |
| `spec.consumerGroups[].consumeRetryPolicy.maxRetryTimes` | `int32` |  |  |  |
| `spec.consumerGroups[].consumeRetryPolicy.deadLetterTargetTopic` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the RocketMQ instance will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.seriesCode

`string` · required

Instance edition series. Determines feature set and throughput ceiling.
"standard" -- cost-effective for light workloads (dev/test).
"professional" -- production workloads with higher throughput.
"ultimate" -- mission-critical workloads with the highest throughput.
ForceNew: changing this replaces the instance.

- rule: series_code must be one of: standard, professional, ultimate
- rule: {"required":true}

### spec.subSeriesCode

`string` · required

Instance deployment architecture within the chosen series.
"cluster_ha" -- multi-node cluster with high availability.
"single_node" -- single node for development and testing.
"serverless" -- auto-scaling serverless deployment.
ForceNew: changing this replaces the instance.

- rule: sub_series_code must be one of: cluster_ha, single_node, serverless
- rule: {"required":true}

### spec.vpcId

`string | valueFrom` · required

VPC ID where the RocketMQ instance is deployed.
ForceNew: changing this replaces the instance.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.instanceName

`string`

Human-readable instance name. 2-64 characters.
If omitted, defaults to metadata.name.

- rule: instance_name must be between 2 and 64 characters when set

### spec.remark

`string`

Instance description or remark.

### spec.paymentType

`string` · optional (explicit presence)

Billing method.
"PayAsYouGo" -- pay-as-you-go (post-paid).
"Subscription" -- prepaid with commitment period.
Default: "PayAsYouGo"

- default: `PayAsYouGo`
- rule: payment_type must be one of: PayAsYouGo, Subscription

### spec.period

`int32` · optional (explicit presence)

Subscription period length. Only applicable when payment_type is
"Subscription". Used together with period_unit.

- rule: period must be positive when set

### spec.periodUnit

`string` · optional (explicit presence)

Subscription period unit.
"Month" or "Year". Only applicable when payment_type is "Subscription".

- rule: period_unit must be one of: Month, Year

### spec.autoRenew

`bool` · optional (explicit presence)

Enable auto-renewal for Subscription instances.

### spec.autoRenewPeriod

`int32` · optional (explicit presence)

Auto-renewal period. Valid values: 1, 2, 3, 6, 12.
Only applicable when auto_renew is true.

- rule: auto_renew_period must be one of: 1, 2, 3, 6, 12

### spec.vswitchId

`string | valueFrom`

VSwitch ID for the instance's VPC endpoint. When provided, the instance
is placed in this VSwitch's availability zone within the VPC.
For serverless instances, at least two VSwitches are recommended.
ForceNew: changing this replaces the instance.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.securityGroupId

`string`

Security group ID to associate with the instance's VPC endpoint.
ForceNew: changing this replaces the instance.

### spec.internetInfo

`AliCloudRocketmqInternetInfo`

Internet access configuration. When omitted, internet access is disabled
and the instance is only accessible within the VPC.

### spec.internetInfo.enabled

`bool` · optional (explicit presence)

Enable public internet access. When true, the instance gets a public
endpoint in addition to the VPC endpoint.
Default: false (internet disabled)

### spec.internetInfo.flowOutType

`string` · optional (explicit presence)

Billing type for public network outbound traffic. Only relevant when
enabled is true.
"payByBandwidth" -- fixed bandwidth billing (set flow_out_bandwidth).
"payByTraffic" -- usage-based traffic billing.
Default: "payByTraffic"

- default: `payByTraffic`
- rule: flow_out_type must be one of: payByBandwidth, payByTraffic

### spec.internetInfo.flowOutBandwidth

`int32` · optional (explicit presence)

Public network bandwidth in Mb/s. Range: 1-1000.
Only applicable when flow_out_type is "payByBandwidth".

- rule: flow_out_bandwidth must be between 1 and 1000

### spec.msgProcessSpec

`string`

Message processing specification that determines the instance's
throughput capacity. Examples: "rmq.s1.micro", "rmq.p2.4xlarge",
"rmq.u2.4xlarge". Required when product_info is set.
The valid values depend on the selected series_code.

### spec.productInfo

`AliCloudRocketmqProductInfo`

Advanced product configuration for message retention, auto-scaling,
tracing, and encryption at rest. When omitted, provider defaults apply.

### spec.productInfo.messageRetentionTime

`int32` · optional (explicit presence)

Duration of message retention in hours. Controls how long messages are
stored before expiration. Longer retention increases storage costs.

- rule: message_retention_time must be positive when set

### spec.productInfo.sendReceiveRatio

`double` · optional (explicit presence)

Ratio of message send to receive capacity. Range: 0.2-0.5.

- rule: send_receive_ratio must be between 0.2 and 0.5

### spec.productInfo.autoScaling

`bool` · optional (explicit presence)

Enable auto-scaling of throughput capacity.

### spec.productInfo.traceOn

`bool` · optional (explicit presence)

Enable message trace functionality for debugging and monitoring
message flows.

### spec.productInfo.storageEncryption

`bool` · optional (explicit presence)

Enable encryption at rest for stored messages.
ForceNew: changing this replaces the instance.

### spec.productInfo.storageSecretKey

`string`

KMS key for encryption at rest. Only applicable when
storage_encryption is true.
ForceNew: changing this replaces the instance.

### spec.ipWhitelists

`[]string`

IP addresses or CIDR blocks allowed to access the instance.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID (per DD05).

### spec.tags

`map<string, string>`

Tags to apply to the RocketMQ instance.

### spec.topics

`[]AliCloudRocketmqTopic`

Topics to create within the instance.

### spec.topics[].topicName

`string` · required

Topic name. Must match ^[%a-zA-Z0-9_-]+$ and be unique within the
instance. ForceNew: changing this replaces the topic.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.topics[].messageType

`string` · optional (explicit presence)

Message type determines how messages in this topic are processed.
"NORMAL" -- standard messages with no ordering guarantee.
"FIFO" -- messages delivered in strict FIFO order.
"DELAY" -- messages delivered after a configurable delay.
"TRANSACTION" -- messages with two-phase commit semantics.
Default: "NORMAL"
ForceNew: changing this replaces the topic.

- default: `NORMAL`
- rule: message_type must be one of: NORMAL, FIFO, DELAY, TRANSACTION

### spec.topics[].remark

`string`

Human-readable description of the topic's purpose.

### spec.topics[].maxSendTps

`int32` · optional (explicit presence)

Maximum send TPS (transactions per second) for the topic.

### spec.consumerGroups

`[]AliCloudRocketmqConsumerGroup`

Consumer groups to create within the instance.

### spec.consumerGroups[].consumerGroupId

`string` · required

Consumer group ID. Must be unique within the instance.
ForceNew: changing this replaces the consumer group.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.consumerGroups[].deliveryOrderType

`string` · optional (explicit presence)

Message delivery order type.
"Concurrently" -- messages delivered in parallel (higher throughput).
"Orderly" -- messages delivered in order per message group.

- rule: delivery_order_type must be one of: Concurrently, Orderly

### spec.consumerGroups[].remark

`string`

Human-readable description of the consumer group's purpose.

### spec.consumerGroups[].maxReceiveTps

`int32` · optional (explicit presence)

Maximum receive TPS for the consumer group.

### spec.consumerGroups[].consumeRetryPolicy

`AliCloudRocketmqConsumeRetryPolicy`

Retry policy for failed message consumption. When omitted, defaults to
DefaultRetryPolicy with 16 retries.

### spec.consumerGroups[].consumeRetryPolicy.retryPolicy

`string` · optional (explicit presence)

Retry strategy.
"DefaultRetryPolicy" -- exponential backoff (recommended).
"FixedRetryPolicy" -- fixed interval between retries.
Default: "DefaultRetryPolicy"

- default: `DefaultRetryPolicy`
- rule: retry_policy must be one of: DefaultRetryPolicy, FixedRetryPolicy

### spec.consumerGroups[].consumeRetryPolicy.maxRetryTimes

`int32` · optional (explicit presence)

Maximum number of retry attempts. Range: 0-1000.
Default: 16 (provider default).

- rule: {"int32":{"lte":1000,"gte":0}}

### spec.consumerGroups[].consumeRetryPolicy.deadLetterTargetTopic

`string`

Dead-letter topic name. When a message exhausts all retries, it is
delivered to this topic for manual investigation.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudRocketmqInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.instance_id` | `string` | The RocketMQ instance ID assigned by Alibaba Cloud. |
| `status.outputs.tcp_endpoint` | `string` | TCP endpoint for VPC-internal access. Applications within the same VPC use this endpoint to produce and consume messages. |
| `status.outputs.internet_endpoint` | `string` | TCP endpoint for public internet access. Empty when internet access is disabled on the instance. |
| `status.outputs.topic_ids` | `map<string, string>` | Map of topic names to their resource IDs. Keys are the topic_name values specified in spec.topics[]. |
| `status.outputs.consumer_group_ids` | `map<string, string>` | Map of consumer group IDs to their resource IDs. Keys are the consumer_group_id values specified in spec.consumer_groups[]. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
