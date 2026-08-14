# AwsBedrockAgentCoreRuntime — Pulumi module (Go)

Deploys an Amazon Bedrock AgentCore agent runtime
(`bedrock.AgentcoreAgentRuntime`) with its named endpoints
(`bedrock.AgentcoreAgentRuntimeEndpoint`) and optional resource policy
(`bedrock.AgentcoreResourcePolicy`).

Module facts worth knowing before editing:

- **Every spec change creates a new runtime VERSION in place**;
  switching the artifact arm (code ↔ container) replaces the runtime
  (provider-enforced).
- **Endpoints key on their names** — iteration is name-sorted for
  deterministic previews, and each name is the AWS identity (there is no
  separate endpoint ID).
- **The resource policy attaches to the runtime's own ARN** — the
  module marshals the Struct to the provider's normalized-JSON string.
- **The override entries' private endpoint is a DISTINCT bridge type**
  from the top-level one (`...PrivateEndpointOverridePrivateEndpointArgs`)
  — the two builder helpers mirror one shape; change them together.
- **Input structs follow the SDK doc examples verbatim** (`&XArgs{...}`
  forms) — the generated `XPtr(...)` forms compile but panic the
  marshaler at preview.

Outputs mirror the Terraform module key-for-key: `agent_runtime_id`,
`agent_runtime_arn`, `agent_runtime_version`, `workload_identity_arn`,
`endpoint_arns`.
