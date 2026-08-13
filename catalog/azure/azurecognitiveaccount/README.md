# Overview

The **AzureCognitiveAccount** component deploys an Azure AI services account -- the container every Azure AI capability is provisioned and billed through: Azure OpenAI model deployments, the multi-service AI Services account that backs AI Foundry projects and agents, and the single-service accounts (Speech, Vision, Language, Content Safety, and friends). The account owns the endpoint, the access keys, the network perimeter, and the responsible-AI (content-filter) policy.

## Purpose

- **One container for AI capacity**: the account is what Azure bills and rate-limits; model deployments (AzureCognitiveDeployment) and AI Foundry projects (AzureCognitiveAccountProject) are created onto it.
- **The kind decides the product**: `OpenAI` hosts Azure OpenAI model deployments; `AIServices` is the multi-service account (and the only kind that supports AI Foundry projects and agent network injection); the remaining values are single-service accounts.
- **Security posture in one place**: custom subdomain + network ACLs for the perimeter, Entra-ID-only auth (`localAuthEnabled: false`), customer-managed keys, restricted outbound with an FQDN allowlist.
- **Responsible AI as configuration**: content-filter policies and custom blocklists deploy as composed children; model deployments select a policy by name.

## Key Features

- Full azurerm v5 surface: 37 account kinds, network ACLs with VNet rules and trusted-services bypass, agent network injection, customer-managed keys, user-owned storage, Metrics Advisor / QnA Maker / custom question answering parameters.
- Kind-gated contracts enforced upfront by validation (project management only on `AIServices`, bypass only on the AI kinds, and friends) -- errors land in seconds, not at deploy time.
- Composed responsible-AI children with name-keyed `rai_blocklist_ids` / `rai_policy_ids` outputs.
- Sensitive outputs (`primary_access_key`, `secondary_access_key`) marked secret in both engines.

## Use Cases

- **Azure OpenAI**: an `OpenAI`-kind account plus AzureCognitiveDeployment resources for gpt-4o / embeddings.
- **AI Foundry**: an `AIServices`-kind account with `projectManagementEnabled: true`, plus AzureCognitiveAccountProject workspaces per team.
- **Single-service AI**: Speech, Vision, Language, Content Safety accounts with the same perimeter and key discipline.

## Future Enhancements

- Account connections (the AI Foundry connection surface for storage, search, and partner services) -- deferred until the Pulumi provider can express them; the upgrade is additive.
- A typed reference for `customQuestionAnsweringSearchServiceId` once the Azure AI Search kind registers.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
