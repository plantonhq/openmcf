# AwsBedrockAgentCoreRuntime — Terraform/OpenTofu module

Deploys an Amazon Bedrock AgentCore agent runtime
(`aws_bedrockagentcore_agent_runtime`) with its named endpoints
(`aws_bedrockagentcore_agent_runtime_endpoint`) and optional resource
policy (`aws_bedrockagentcore_resource_policy`).

Module facts worth knowing before editing:

- **Every spec change creates a new runtime VERSION in place**;
  switching the artifact arm (code ↔ container) replaces the runtime
  (provider-enforced `RequiresReplaceIf`).
- **Endpoints key on their names** — the for_each key, the
  `endpoint_arns` output key, and the AWS identity are all the same
  string (there is no separate endpoint ID).
- **The resource policy attaches to the runtime's own ARN** — the
  provider resource accepts any AgentCore ARN; this module scopes it to
  the runtime it deploys and renders the Struct via `jsonencode()`.
- **The JWT authorizer's single-member wrapper is flattened** in the
  spec (`custom_jwt_authorizer` directly), as are the protocol and
  request-header one-attribute blocks.
- **HTTP is AWS's default protocol** — the block renders only on an
  explicit `server_protocol` so the module never fights the default.

Outputs mirror the Pulumi module key-for-key: `agent_runtime_id`,
`agent_runtime_arn`, `agent_runtime_version`, `workload_identity_arn`,
`endpoint_arns`.
