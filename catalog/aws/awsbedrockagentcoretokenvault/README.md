<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Token Vault" width="80"/>
</p>

# AWS Bedrock AgentCore Token Vault

Manage the encryption key on the region's [AgentCore token vault](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/security-token-vault.html)
— the store [AgentCore Identity](../awsbedrockagentcoreidentity)
keeps OAuth tokens and API keys in.

This is a **settings singleton**: AWS provisions ONE default vault
per account+region and this component sets which KMS key encrypts it.
Deploy at most one instance per vault. `metadata.name` never reaches
AWS.

## What Gets Managed

- **The vault's key ownership**: your own KMS key
  (`CustomerManagedKey` — you control rotation, policy, and
  revocation) or AWS's (`ServiceManagedKey`, the default posture).

Destroying this component is a **no-op at AWS**: the last-applied key
setting remains in effect. To return to AWS-managed encryption, apply
an instance with `ServiceManagedKey` — destroy alone does not revert
it. The setting is free; a customer-managed KMS key carries KMS's
flat monthly per-key charge.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
