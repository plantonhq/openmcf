# AwsBedrockPrompt — Terraform/OpenTofu module

Deploys an Amazon Bedrock prompt (`aws_bedrockagent_prompt`) from Prompt
Management, with its variants (text or chat templates, model or
agent-alias execution targets, tools, and inference settings).

Module facts worth knowing before editing:

- **`template_type` is derived**, never spec surface: exactly one of a
  variant's `text`/`chat` arms is set (CEL guarded) and the module sends
  TEXT or CHAT accordingly. The same derivation picks `model_id` versus
  the `gen_ai_resource` agent-alias block.
- **Only the DRAFT version is managed** — AWS assigns a new version string
  on every update; there is no prompt-version resource at the pin.
- **One-value vocabularies are module constants**: cache point type
  `default`.
- **The template tree mirrors the AwsBedrockFlow module's inline-prompt
  rendering** arm-for-arm (upstream shares the same Go models between the
  two resources) — change them together.
- **No timeouts**: the upstream resource is plain CRUD with no waiters.

Outputs mirror the Pulumi module key-for-key: `prompt_id`, `prompt_arn`,
`draft_version`.
