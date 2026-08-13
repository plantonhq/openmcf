# MSK Consumer with Shared Provisioned Pollers

This preset consumes an Amazon MSK topic with a pinned poller fleet that is SHARED across mappings: every mapping naming the same `pollerGroupName` draws from (and is jointly capped by) one provisioned fleet instead of provisioning its own. That is the cost lever for many low-traffic topics — one warm fleet serves them all, with schema-validated records and full Kafka metrics.

## When to Use

- Consuming many Kafka topics where per-mapping poller fleets would sit idle
- Predictable-throughput Kafka ingestion (provisioned pollers instead of reactive scaling)
- Consumers that must validate records against a schema registry before invocation
- Resuming an existing consumer group's committed offsets (`kafkaConsumerGroupId`)

## Key Configuration Choices

- **`pollerGroupName`** — the fleet-sharing key. Mappings with the same group name in the same account/region share one fleet bounded by the group's poller limits.
- **Provisioned pollers 2–50** — a warm floor of 2 keeps latency flat; the ceiling caps spend. Both bounds apply to the GROUP when a name is set.
- **Schema registry** — records are validated (VALUE attribute) against a Glue schema registry and delivered as JSON; bad records never bill function time.
- **`TRIM_HORIZON`** — processes the topic's retained backlog on first attach; committed group offsets take precedence for partitions the group already consumed.
- **Full metrics** — `EventCount`, `ErrorCount`, and the Kafka-only `KafkaMetrics` (consumer-lag visibility). Metrics are billed.

## Placeholders to Replace

| Placeholder | Description | Example |
|-------------|-------------|---------|
| `<aws-region>` | Region of the cluster, function, and mapping | `us-west-2` |
| `<glue-schema-registry-arn>` | Glue schema registry ARN (or a Confluent HTTPS URL) | `arn:aws:glue:us-west-2:123456789012:registry/events` |

The `valueFrom` references assume an `orders-processor` function and an `events-cluster` MSK cluster managed as Planton resources; replace with your own names or literal ARNs.

## Common Additions

- Add `filters` to discard irrelevant records before invocation (no function billing for filtered records)
- Add `onFailureDestinationArn` (SQS, SNS, or S3 for Kafka sources) to capture failed batches
- Add `sourceAccessConfigurations` for SASL/mTLS when the cluster requires client auth
- Set `startingPosition: LATEST` for consumers that only care about new records

## Related Presets

- **01-sqs-worker** — use for queue-based work distribution
- **02-kinesis-consumer** — use for Kinesis stream processing with shard-level parallelism
