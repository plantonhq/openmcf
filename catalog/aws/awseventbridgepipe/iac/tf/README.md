# AwsEventBridgePipe — Terraform/OpenTofu module

Manages one EventBridge Pipe (`aws_pipes_pipe`): source → (filter) → (enrichment) → target.

Module facts worth knowing before editing:

- **The source is fixed for life** — `source` and the per-family positions (`starting_position`, `topic_name`, `queue_name`, self-managed Kafka's `consumer_group_id`, `additional_bootstrap_servers`) replace the pipe on change; the **target swaps in place**.
- **`desired_state` is the pause lever** — STOPPED halts consumption without deleting; stream positions are kept. Creates and state flips can take minutes (the provider waits up to 30).
- **Credentials are Secrets Manager ARNs** — Kafka/MQ auth fields carry references, never credential values (the spec's patterns enforce it).
- **Removing `input_template` genuinely clears it** — the provider pre-seeds the empty value on update so a dropped template does not linger at AWS.
- **`assign_public_ip` is a string enum at the provider** (`ENABLED`/`DISABLED`) for ECS targets — the module maps the spec's bool.
- **`include_execution_data` is a list at the provider** (`["ALL"]`) — the module maps the spec's bool. TRACE + execution data logs payloads; mind the sensitivity.
- **One family block per side** — the spec's CELs guarantee at most one source family and one target family, mirroring the provider's ConflictsWith lists.

Outputs mirror the Pulumi module key-for-key: `pipe_arn`, `pipe_name` (import ID).
