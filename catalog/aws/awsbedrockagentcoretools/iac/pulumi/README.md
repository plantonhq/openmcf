# AwsBedrockAgentCoreTools — Pulumi module (Go)

Deploys an Amazon Bedrock AgentCore tools bundle: browsers
(`bedrock.AgentcoreBrowser`), browser profiles
(`bedrock.AgentcoreBrowserProfile`), and code interpreters
(`bedrock.AgentcoreCodeInterpreter`).

Module facts worth knowing before editing:

- **AWS exposes NO update for any of the three resources** — every
  argument is ForceNew at the provider; a spec change recreates the
  tool.
- **The browser and code interpreter share the certificate shape** (a
  Secrets Manager location); the enterprise-policy and recording shapes
  are browser-only.
- **VPC placement pairs with mode VPC** (spec-validated both ways);
  SANDBOX exists only on code interpreters.
- **Iteration is name-sorted per arm** for deterministic previews.
- **Input structs follow the SDK doc examples verbatim** (`&XArgs{...}`
  forms) — the generated `XPtr(...)` forms compile but panic the
  marshaler at preview.

Outputs mirror the Terraform module key-for-key: `browser_ids`,
`browser_arns`, `browser_profile_ids`, `browser_profile_arns`,
`code_interpreter_ids`, `code_interpreter_arns`.
