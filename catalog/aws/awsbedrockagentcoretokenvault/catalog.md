# AWS Bedrock AgentCore Token Vault

The encryption posture of the region's AgentCore token vault — the
store AgentCore Identity keeps OAuth tokens and API keys in. A
settings singleton: one default vault per region, at most one
instance deployed per vault.

## What Gets Managed

- The vault's key ownership: your KMS key
  ([AWS KMS Key](/cloud-catalog/aws-kms-key), `CustomerManagedKey`)
  or AWS's (`ServiceManagedKey`).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore and KMS permissions.

### AWS Account

- For the customer-managed posture: a symmetric KMS key in the same
  region.

## Deploy

### Console

Create the resource from the AWS catalog, pick the key type, and
deploy.

### CLI

```bash
planton apply -f agentcore-token-vault.yaml
```

## After Deploy

- Every credential AgentCore Identity stores in the region is
  encrypted under the configured key.
- Revoking a customer-managed key's grants locks AgentCore out of
  every stored credential — treat key policy changes as outage-grade.
- **Destroy does not revert the setting.** Apply `ServiceManagedKey`
  to return to AWS-managed encryption before retiring the key.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
