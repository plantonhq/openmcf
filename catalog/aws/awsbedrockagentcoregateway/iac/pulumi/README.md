# AwsBedrockAgentCoreGateway — Pulumi module (Go)

Deploys an Amazon Bedrock AgentCore gateway (`bedrock.AgentcoreGateway`)
with its target satellites (`bedrock.AgentcoreGatewayTarget`).

Module facts worth knowing before editing:

- **The bridge types the tool-schema recursion per tree position** — the
  input and output schemas carry parallel but DISTINCT Go types at every
  level (`...InputSchemaPropertyItemsArgs` vs
  `...OutputSchemaPropertyItemsArgs`), so `schema.go` holds two
  position-specific builders of one shape; change them together.
- **AWS deletes targets before the gateway at destroy** — the provider
  manages the drain; the module orders creates via `DependsOn` only.
- **One-value vocabularies are module constants**: protocol type (MCP,
  never sent), search type (SEMANTIC), exception level (DEBUG).
- **`JwtPassthrough` is an empty args struct** — presence IS the
  configuration.
- **Raw-JSON schema leaves marshal Structs to the provider's
  normalized-JSON strings** (`structToJson`).
- **Input structs follow the SDK doc examples verbatim** (`&XArgs{...}`
  forms) — the generated `XPtr(...)` forms compile but panic the
  marshaler at preview.

Outputs mirror the Terraform module key-for-key: `gateway_id`,
`gateway_arn`, `gateway_url`, `workload_identity_arn`, `target_ids`.
