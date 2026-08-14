# AwsBedrockAgentCoreTokenVault — Terraform/OpenTofu module

Sets the KMS key on the region's AgentCore token vault
(`aws_bedrockagentcore_token_vault_cmk`) — the store AgentCore
Identity keeps OAuth tokens and API keys in.

Module facts worth knowing before editing:

- **One default vault per account+region**, named "default"; an unset
  spec id targets it. Changing the id replaces the configuration onto
  the other vault.
- **Destroy is a no-op at AWS.** The provider has no delete API for
  the setting; whatever key type was last applied remains. Reverting
  to AWS-managed encryption is an apply with `ServiceManagedKey`,
  never a destroy — the E2E profile's verify-absent asserts the vault
  still exists.
- **The key pairing is CEL-enforced** (CustomerManagedKey requires
  the ARN, ServiceManagedKey forbids it), so the module never sends a
  stray ARN.
- **No tags.** The upstream resource carries no tags argument.

Outputs mirror the Pulumi module key-for-key: `token_vault_id`,
`key_type`, `kms_key_arn`.
