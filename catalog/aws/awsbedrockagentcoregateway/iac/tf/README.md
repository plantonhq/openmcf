# AwsBedrockAgentCoreGateway — Terraform/OpenTofu module

Deploys an Amazon Bedrock AgentCore gateway
(`aws_bedrockagentcore_gateway`) with its target satellites
(`aws_bedrockagentcore_gateway_target`).

Module facts worth knowing before editing:

- **`spec.targets` is `any`-typed in variables.tf** — the tool-schema
  tree's raw-JSON members defeat HCL's object-type unification, so every
  target attribute access in main.tf is `try()`-based sparse access
  against the tfvars converter's key omission (zero values are absent,
  never empty).
- **AWS deletes targets before the gateway at destroy** — the provider
  manages the drain; the module adds no ordering of its own beyond
  create-after via the gateway-ID reference.
- **One-value vocabularies are module constants**: `protocol_type` (MCP,
  provider-computed and never sent), `search_type` (SEMANTIC on
  `enable_semantic_search`), `exception_level` (DEBUG on
  `expose_debug_exceptions`).
- **`jwt_passthrough` renders an EMPTY provider block** — presence IS
  the configuration.
- **The three-level tool-schema render bottoms out in `jsonencode()`d
  Struct leaves** (`items_json`/`properties_json`) — exactly where the
  provider's own schema stops typing.

Outputs mirror the Pulumi module key-for-key: `gateway_id`,
`gateway_arn`, `gateway_url`, `workload_identity_arn`, `target_ids`.
