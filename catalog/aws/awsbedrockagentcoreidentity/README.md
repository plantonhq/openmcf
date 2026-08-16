<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Identity" width="80"/>
</p>

# AWS Bedrock AgentCore Identity

Create and manage [Amazon Bedrock AgentCore identity and access](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/identity.html)
— the credentials and authorization control plane your agents,
gateways, and tools authenticate through, as ONE bundle:

## What Gets Created

- **Workload identities** — named identities AgentCore workloads
  present when calling other services, with OAuth2 return-URL
  allow-lists.
- **API key credential providers** — vaulted API keys (AWS stores the
  key in Secrets Manager under the service's token vault; consumers
  reference the provider ARN, never the secret).
- **OAuth2 credential providers** — vaulted client credentials for
  well-known vendors (GitHub, Google, Microsoft, Salesforce, Slack) or
  any CUSTOM OIDC provider (discovery URL or spelled-out endpoints).
- **A Cedar policy engine** with its **policies** — the authorization
  engine gateways evaluate tool calls against (LOG_ONLY or ENFORCE).

Every arm is optional — author the ones this bundle owns.

## Secrets Handling

The API key / client credential values are sensitive fields: supply
managed-secret references resolved just-in-time at deploy, never
literals. Rotating a value updates the vault in place; consumers keep
referencing the same provider ARN.

The account/region **token-vault CMK** is deliberately NOT part of this
kind — it is an account-level settings singleton.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
