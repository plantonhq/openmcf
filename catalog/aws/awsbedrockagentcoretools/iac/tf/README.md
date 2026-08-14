# AwsBedrockAgentCoreTools — Terraform/OpenTofu module

Deploys an Amazon Bedrock AgentCore tools bundle: browsers
(`aws_bedrockagentcore_browser`), browser profiles
(`aws_bedrockagentcore_browser_profile`), and code interpreters
(`aws_bedrockagentcore_code_interpreter`).

Module facts worth knowing before editing:

- **AWS exposes NO update for any of the three resources** — every
  argument is ForceNew at the provider; a spec change recreates the
  tool. The module renders plain arguments and lets the provider drive
  the replace.
- **The browser and code interpreter share the certificate shape** (a
  Secrets Manager location); the enterprise-policy and recording shapes
  are browser-only.
- **VPC placement pairs with mode VPC** (spec-validated both ways);
  SANDBOX exists only on code interpreters.
- **`browser_signing` and `recording.enabled` render only on explicit
  choices** so the module never fights AWS's defaults.

Outputs mirror the Pulumi module key-for-key: `browser_ids`,
`browser_arns`, `browser_profile_ids`, `browser_profile_arns`,
`code_interpreter_ids`, `code_interpreter_arns`.
