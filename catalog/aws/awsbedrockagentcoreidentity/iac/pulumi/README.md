# AwsBedrockAgentCoreIdentity — Pulumi module (Go)

Deploys an Amazon Bedrock AgentCore identity-and-access bundle:
workload identities (`bedrock.AgentcoreWorkloadIdentity`), vaulted
credentials (`bedrock.AgentcoreApiKeyCredentialProvider`,
`bedrock.AgentcoreOauth2CredentialProvider`), and a Cedar policy engine
with policies (`bedrock.AgentcorePolicyEngine`,
`bedrock.AgentcorePolicy`).

Module facts worth knowing before editing:

- **The spec's vendor enum selects the provider config member** — the
  switch in `identity.go` maps CUSTOM/GITHUB/... to
  CustomOauth2/GithubOauth2/... AND the matching one of six
  structurally-identical members; only CUSTOM carries `OauthDiscovery`.
- **The write-only argument variants (`ApiKeyWo`, `ClientIdWo`,
  `ClientSecretWo`) are never set** — the spec's sensitive fields arrive
  JIT-resolved, and the plain arguments let the provider detect
  rotation.
- **A policy is a structural child of its engine** — `DependsOn` the
  engine, keyed by its computed ID.
- **Iteration is name-sorted per arm** for deterministic previews.
- **Input structs follow the SDK doc examples verbatim** (`&XArgs{...}`
  forms) — the generated `XPtr(...)` forms compile but panic the
  marshaler at preview.

Outputs mirror the Terraform module key-for-key: `workload_identity_arns`,
`api_key_provider_arns`, `api_key_secret_arns`, `oauth2_provider_arns`,
`oauth2_client_secret_arns`, `policy_engine_id`, `policy_engine_arn`,
`policy_ids`.
