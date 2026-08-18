# AwsEventBridgePipe — Pulumi module (Go)

Manages one EventBridge Pipe (`pipes.Pipe`): source → (filter) → (enrichment) → target.

Module facts worth knowing before editing:

- **The source is fixed for life** — `Source` and the per-family positions (`StartingPosition`, `TopicName`, `QueueName`, self-managed Kafka's `ConsumerGroupId`, `AdditionalBootstrapServers`) replace the pipe on change; the **target swaps in place**.
- **`DesiredState` is the pause lever** — STOPPED halts consumption without deleting; stream positions are kept. Creates and state flips can take minutes (the provider waits up to 30).
- **Credentials are Secrets Manager ARNs** — Kafka/MQ auth fields carry references, never credential values (the spec's patterns enforce it).
- **`AssignPublicIp` is a string enum at the provider** (`ENABLED`/`DISABLED`) for ECS targets — the module maps the spec's bool.
- **`IncludeExecutionDatas` is a list at the provider** (`["ALL"]`) — the module maps the spec's bool. TRACE + execution data logs payloads; mind the sensitivity.
- **`PathParameterValues` is a single string at this SDK** — the bridge collapsed the one-entry list AWS models; the spec's single `path_parameter_value` maps directly.
- **One family block per side** — the spec's CELs guarantee at most one source family and one target family, mirroring the provider's ConflictsWith lists.

Outputs mirror the Terraform module key-for-key: `pipe_arn`, `pipe_name` (import ID).
