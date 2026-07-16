# AWS MSK Serverless: Kafka Without Capacity Planning

## Introduction

Running Apache Kafka has always carried a capacity-planning tax. Provisioned clusters — even fully managed ones like Amazon MSK — make you choose broker counts, instance types, and storage volumes up front, then watch utilization to decide when to scale. For teams with spiky, unpredictable, or simply *new* streaming workloads, that tax is paid before the first message flows.

MSK Serverless removes the tax: AWS manages brokers, storage, and the Kafka version behind a single cluster endpoint, scales capacity with actual throughput, and bills per data in/out and storage consumed. The trade is control for operations — no custom `server.properties`, no instance selection, and exactly one authentication scheme (SASL/IAM).

## The Shape of the Resource

A serverless cluster is one of the smallest resources in the AWS catalog, and honestly so:

- **Networking** — the subnets where the cluster places its network interfaces and the security groups attached to them. This is the entire configurable surface.
- **Authentication** — SASL/IAM, always on. AWS offers no alternative, so the component does not model it as a choice: both IaC engines enable it unconditionally.
- **Immutability** — everything except tags is create-time. Changing subnets or security groups replaces the cluster. Plan network placement before the first deploy.

## Serverless vs Provisioned MSK

| Dimension | `AwsMskServerlessCluster` | `AwsMskCluster` (provisioned) |
|---|---|---|
| Capacity | AWS-managed, auto-scaling | You choose broker count, instance type, storage |
| Billing | Per throughput + storage consumed | Per broker-hour + storage |
| Authentication | SASL/IAM only | SASL/IAM, SASL/SCRAM, mTLS, unauthenticated |
| Encryption in transit | TLS, always | TLS / TLS_PLAINTEXT / PLAINTEXT |
| Public access | No | Optional (SERVICE_PROVIDED_EIPS) |
| Multi-VPC PrivateLink | No | Yes (`vpc_connectivity`) |
| Custom server.properties | No | Yes (folded configuration) |
| Tiered storage | No (AWS manages storage) | Optional |
| Kafka version | AWS-managed | You choose and upgrade |
| Mutability | Immutable (except tags) | Rich in-place update surface |

Rule of thumb: start serverless for new event-driven architectures and low-duty-cycle pipelines; move to provisioned when sustained throughput makes broker-hours cheaper, or when you need an auth scheme or connectivity feature serverless lacks.

## Composition

The cluster composes onto the same networking graph as every other data-plane component:

- **Subnets** come from first-class `AwsSubnet` nodes (2+ AZs recommended).
- **Security groups** come from first-class `AwsSecurityGroup` nodes. The ingress rule for the SASL/IAM listener (port 9098) lives on the referenced group — typically allowing the application tier's security group — keeping network policy sharable and auditable outside the cluster.
- **Client access** is IAM: workloads produce and consume with `kafka-cluster:*` permissions scoped to the exported `cluster_arn`. There are no passwords, no certificates, and nothing to rotate.

## Connecting Clients

The single exported endpoint list is `bootstrap_brokers_sasl_iam` (port 9098). Kafka clients need:

1. The AWS MSK IAM auth library for their language (e.g. `aws-msk-iam-sasl-signer`).
2. `security.protocol=SASL_SSL` and `sasl.mechanism=AWS_MSK_IAM`.
3. IAM permissions: `kafka-cluster:Connect` plus topic/group actions (`CreateTopic`, `WriteData`, `ReadData`, `AlterGroup`, ...) scoped to the cluster ARN.

## Quotas Worth Knowing

Serverless MSK enforces service quotas (per cluster) on maximum ingress/egress throughput, partition counts, and retention. Check the current AWS quotas page before committing a high-throughput workload — sustained loads beyond the quotas are the signal to move to a provisioned cluster.
