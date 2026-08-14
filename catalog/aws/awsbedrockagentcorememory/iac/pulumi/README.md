# AwsBedrockAgentCoreMemory — Pulumi module (Go)

Deploys an Amazon Bedrock AgentCore memory (`bedrock.AgentcoreMemory`)
with its strategy satellites (`bedrock.AgentcoreMemoryStrategy`).

Module facts worth knowing before editing:

- **Strategy writes serialize through the parent memory** — AWS
  processes them one at a time per memory and the provider holds a
  per-memory lock with a 45m default timeout; the module orders
  strategies after the memory (name-sorted for deterministic previews)
  and adds nothing else.
- **`NamespaceTemplates` is always sent** — the provider pairs it
  exactly-one with the deprecated `Namespaces` twin (never rendered;
  excluded-deprecated), so an unset pair fails at preview.
- **MEMORY_RECORDS is the only stream content type** AWS defines — a
  module constant on the Kinesis delivery render.
- **Input structs follow the SDK doc examples verbatim** (`&XArgs{...}`
  forms) — the generated `XPtr(...)` forms compile but panic the
  marshaler at preview.

Outputs mirror the Terraform module key-for-key: `memory_id`,
`memory_arn`, `strategy_ids`.
