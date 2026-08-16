# AwsBedrockAgentCoreIdentity — Terraform/OpenTofu module

Deploys an Amazon Bedrock AgentCore identity-and-access bundle:
workload identities (`aws_bedrockagentcore_workload_identity`), vaulted
credentials (`aws_bedrockagentcore_api_key_credential_provider`,
`aws_bedrockagentcore_oauth2_credential_provider`), and a Cedar policy
engine with policies (`aws_bedrockagentcore_policy_engine`,
`aws_bedrockagentcore_policy`).

Module facts worth knowing before editing:

- **The spec's vendor enum selects the provider block** — six
  structurally-identical vendor blocks render from one `vendor` field
  via `locals.oauth2_vendor_values` (CUSTOM/GITHUB/... →
  CustomOauth2/GithubOauth2/...); only CUSTOM carries `oauth_discovery`.
- **The write-only argument variants (`api_key_wo`, `client_id_wo`,
  `client_secret_wo`) are never rendered** — the spec's sensitive fields
  arrive JIT-resolved, and the plain arguments let the provider detect
  rotation.
- **A policy is a structural child of its engine** — created after,
  destroyed before, keyed by the engine's computed ID.
- **The token-vault CMK is deliberately not rendered** — an
  account-level settings singleton outside this kind.

Outputs mirror the Pulumi module key-for-key: `workload_identity_arns`,
`api_key_provider_arns`, `api_key_secret_arns`, `oauth2_provider_arns`,
`oauth2_client_secret_arns`, `policy_engine_id`, `policy_engine_arn`,
`policy_ids`.
