# AwsBedrockAgentCoreMemory — Terraform/OpenTofu module

Deploys an Amazon Bedrock AgentCore memory
(`aws_bedrockagentcore_memory`) with its strategy satellites
(`aws_bedrockagentcore_memory_strategy`).

Module facts worth knowing before editing:

- **Strategy writes serialize through the parent memory** — AWS
  processes them one at a time per memory and the provider holds a
  per-memory lock with a 45m default timeout; the module orders
  strategies after the memory and adds nothing else.
- **`namespace_templates` is always sent** — the provider pairs it
  exactly-one with the deprecated `namespaces` twin (never rendered;
  excluded-deprecated), so an unset pair fails at plan.
- **MEMORY_RECORDS is the only stream content type** AWS defines — a
  module constant on the Kinesis delivery render.
- **`indexed_key` and `encryption_key_arn` replace the memory**
  (provider-enforced ForceNew).
- **The deprecated strategy-level `memory_execution_role_arn` is never
  sent** — the memory-level role is the living surface.

Outputs mirror the Pulumi module key-for-key: `memory_id`, `memory_arn`,
`strategy_ids`.
