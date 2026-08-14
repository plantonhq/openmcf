---
title: "Bedrock AgentCore Identity"
description: "Bedrock AgentCore Identity deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockagentcoreidentity"
---

# AWS Bedrock AgentCore Identity

Agent credentials and authorization as managed infrastructure — vaulted
API keys and OAuth clients your gateways and tools reference by ARN
(never by value), workload identities your agents present, and a Cedar
policy engine authorizing every tool call.

## What Gets Created

- Workload identities with OAuth2 return-URL allow-lists.
- Vaulted credentials: API keys and OAuth2 clients (well-known vendors
  or CUSTOM OIDC) stored by AWS in Secrets Manager under the token
  vault — rotation never touches consumers.
- A Cedar policy engine with permit/forbid policies, attachable to
  gateways in LOG_ONLY or ENFORCE mode.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore identity permissions
  (`bedrock-agentcore:CreateWorkloadIdentity`,
  `CreateApiKeyCredentialProvider`, `CreateOauth2CredentialProvider`,
  `CreatePolicyEngine`, `CreatePolicy` and their siblings).

### AWS Account

- Bedrock AgentCore available in the target region.
- Credential VALUES as managed secrets (resolved just-in-time at
  deploy) — never literals in manifests.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, author the
arms this bundle owns, and deploy.

### CLI

```bash
planton apply -f identity.yaml
```

## After Deploy

- The output maps carry each arm's ARNs: gateway targets consume
  `api_key_provider_arns` / `oauth2_provider_arns`; a gateway's policy
  engine attachment consumes `policy_engine_arn`.
- The vaulted secret ARNs (`api_key_secret_arns`,
  `oauth2_client_secret_arns`) locate the Secrets Manager entries AWS
  manages for you.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
