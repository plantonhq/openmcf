# AWS Kafka Streaming

Managed Kafka on Amazon MSK, serverless-first: an IAM-authenticated cluster
behind a two-tier security-group contract, with a provisioned-broker arm one
toggle away for the workloads where steady throughput makes broker-hour
pricing the better deal. Both arms speak SASL/IAM over TLS on port 9098, so
producers and consumers move between them with nothing but a new bootstrap
string.

This is the event backbone for stream processing: change-data capture,
clickstream and telemetry ingestion, microservice event exchange, and the
feed for Firehose/Flink/Spark consumers downstream.

## Architecture

```
   producers / consumers
   (attach the client SG, assume an IAM role with kafka-cluster grants)
            |
            |  TLS + SASL/IAM, port 9098
            v
   [AwsSecurityGroup] cluster-sg  <- ingress ONLY from client-sg members
            |
            v
   [AwsMskServerlessCluster]      (default: per-GB billing, zero broker ops)
        -- XOR --                  one toggle: provisioned_enabled
   [AwsMskCluster]                (broker-hours you size: count x type x EBS)
            |
   private broker subnets (bring your own, 2+ AZs)
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|---|---|---|
| Client security group | `AwsSecurityGroup` | Attach to producers/consumers; membership IS the network grant. |
| Cluster security group | `AwsSecurityGroup` | Admits port 9098 only from client-group members — never a CIDR. |
| Serverless cluster | `AwsMskServerlessCluster` | The default arm. Rendered while `provisioned_enabled` is off. |
| Provisioned cluster | `AwsMskCluster` | The steady-throughput arm (SASL/IAM only). Rendered while `provisioned_enabled` is on. |

## Parameters

| Name | Description | Default | Required |
|---|---|---|---|
| `aws_region` | Region for the cluster and its security groups. | `us-east-1` | yes |
| `cluster_name` | Cluster name and companion-resource prefix. | `my-kafka` | yes |
| `vpc_id` | VPC for the security groups (the broker subnets' VPC). | example id | yes |
| `broker_subnet_ids` | Private subnets in different AZs for brokers/ENIs. | example ids | yes |
| `kafka_version` | Kafka version — provisioned arm only. | `3.6.0` | provisioned |
| `broker_count` | Total brokers; must be a multiple of the subnet count. | `2` | provisioned |
| `broker_instance_type` | Broker compute class. | `kafka.m7g.large` | provisioned |
| `broker_storage_gib` | EBS per broker (grows, never shrinks). | `100` | provisioned |
| `provisioned_enabled` | Provisioned cluster instead of serverless. | `false` | no |

## Composing with network-foundation

The chart deliberately takes its network as inputs — a Kafka cluster joins
the environment's network, it does not own one. With a network-foundation
deployment named `my-network` in the same environment, swap the literals for
references:

```yaml
vpcId:
  valueFrom:
    kind: AwsVpc
    name: my-network-vpc
    fieldPath: status.outputs.vpc_id
subnetIds:
  - valueFrom:
      kind: AwsSubnet
      name: my-network-private-us-east-1a
      fieldPath: status.outputs.subnet_id
  - valueFrom:
      kind: AwsSubnet
      name: my-network-private-us-east-1b
      fieldPath: status.outputs.subnet_id
```

## Client access (post-deploy)

Network reachability and identity are separate grants, by design:

1. **Network**: attach the chart's client security group
   (`<cluster_name>-client-sg`) to the producer/consumer workload (ECS
   service, EC2 instance, Lambda VPC config).
2. **Identity**: grant the workload's IAM role the Kafka actions, scoped to
   this cluster's ARN (`status.outputs.cluster_arn`) and its topics/groups:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "ConnectAndDescribe",
      "Effect": "Allow",
      "Action": ["kafka-cluster:Connect", "kafka-cluster:DescribeCluster"],
      "Resource": "<cluster_arn>"
    },
    {
      "Sid": "ProduceConsume",
      "Effect": "Allow",
      "Action": [
        "kafka-cluster:CreateTopic",
        "kafka-cluster:DescribeTopic",
        "kafka-cluster:WriteData",
        "kafka-cluster:ReadData"
      ],
      "Resource": "arn:aws:kafka:<region>:<account>:topic/<cluster_name>/*"
    },
    {
      "Sid": "ConsumerGroups",
      "Effect": "Allow",
      "Action": ["kafka-cluster:AlterGroup", "kafka-cluster:DescribeGroup"],
      "Resource": "arn:aws:kafka:<region>:<account>:group/<cluster_name>/*"
    }
  ]
}
```

Tighten `topic/<cluster_name>/*` to specific topic names as the topology
settles — per-topic IAM is MSK's whole authorization story, use it.

3. **Client config**: point the client at the cluster's
   `status.outputs.bootstrap_brokers_sasl_iam` with the standard IAM
   properties (`security.protocol=SASL_SSL`, `sasl.mechanism=AWS_MSK_IAM`,
   and the `aws-msk-iam-auth` callback handler on JVM clients).

## Serverless or provisioned, honestly

- **Serverless** bills per cluster-hour (~$0.75) plus per-GB in/out and
  partition-hours — and needs zero broker capacity math. It also caps
  throughput per cluster (single-digit-GB/s writes; check current quotas).
- **Provisioned** bills broker-hours + EBS regardless of traffic. Two
  `kafka.m7g.large` brokers run ~$300/month before storage.
- The crossover: sustained, around-the-clock throughput in the multi-MB/s
  range usually favors provisioned; anything spiky, growing, or idle
  overnight favors serverless. Start serverless, measure, and flip the
  toggle with real numbers — the client contract (IAM, 9098, the SG pair)
  is identical on both sides. Note the flip REPLACES the cluster: topics
  and offsets do not migrate; use MirrorMaker 2 or replay from source when
  continuity matters.

## Day-2 guidance

- **Provisioned depth when you need it**: the `AwsMskCluster` spec models
  the full surface — TLS/SCRAM auth arms, PrivateLink (`vpcConnectivity`),
  broker log delivery, Prometheus exporters, tiered storage, custom
  `serverProperties`. Manage the deployed resource directly as
  requirements land; the chart is the starting point, not a ceiling.
- **Monitoring**: provisioned clusters expose `enhancedMonitoring` levels
  and the JMX/node exporters for Prometheus; serverless publishes
  per-topic CloudWatch metrics out of the box.
- **Egress hygiene**: the client group's egress is already scoped to 9098;
  if your posture also pins destinations, replace its CIDR egress with a
  destination-SG rule pointing at the cluster group.
- **Schema management**: pair the cluster with AWS Glue Schema Registry
  (client-side config, no infrastructure) before the topic count grows
  past what tribal knowledge can carry.
