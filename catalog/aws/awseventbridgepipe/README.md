# AwsEventBridgePipe

One EventBridge Pipe: a managed point-to-point integration that reads from one source (SQS, Kinesis, DynamoDB streams, MSK/self-managed Kafka, ActiveMQ/RabbitMQ), optionally filters and enriches events in flight, and delivers to one target — sixteen-plus AWS services, no event bus in between.

## Highlights

- **Full provider depth**: all seven source families and eleven target families, each as a CEL-gated union mirroring the provider's own mutual-exclusion lists; enrichment (Lambda, Step Functions express, API destinations), event filtering, input templates, KMS, and execution logging to CloudWatch Logs / Firehose / S3.
- **The source is fixed for life** — changing it (or a family's stream position, topic, or queue) replaces the pipe; the target swaps in place. Taught on every affected field.
- **Credentials are references, never values**: Kafka/MQ auth fields take Secrets Manager secret ARNs, pattern-enforced at validation.
- **`desired_state` is the pause lever**: STOPPED halts consumption without deleting; stream positions survive. Pipes bill per event — an idle pipe costs nothing.
- **Bare polymorphic references**: source/target/DLQ ARNs carry no default kind (no single kind dominates) — a `valueFrom` on these fields states its `kind:` explicitly.

## Both Engines

Both modules render the single `aws_pipes_pipe` / `pipes.Pipe` identically — the ECS `assign_public_ip` string enum and the `include_execution_data` list both map from honest spec booleans — and export the same outputs: `pipe_arn`, `pipe_name` (import ID).

## Chart Wiring

`role_arn` → AwsIamRole `role_arn`; stream DLQs → AwsSqsQueue `queue_arn`; ECS targets → AwsEcsTaskDefinition `task_definition_arn`, AwsSubnet/AwsSecurityGroup ids; log destinations → AwsCloudwatchLogGroup `log_group_arn`, AwsKinesisFirehose `delivery_stream_arn`, AwsS3Bucket `bucket_id`; `kms_key_identifier` → AwsKmsKey `key_arn`. Source/target/enrichment ARNs wire to any producing kind with an explicit `kind:`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
