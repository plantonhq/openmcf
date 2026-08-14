---
title: "GPT-4o Chat"
description: "This preset deploys `gpt-4o-mini` on the GlobalStandard SKU -- the standard starting point for chat and completion workloads: per-token billing, no idle cost, and the widest regional model..."
type: "preset"
rank: "01"
presetSlug: "01-gpt-4o-chat"
componentSlug: "cognitive-deployment"
componentTitle: "Cognitive Deployment"
provider: "azure"
icon: "package"
order: 1
---

# GPT-4o Chat

This preset deploys `gpt-4o-mini` on the GlobalStandard SKU -- the standard starting point for chat and completion workloads: per-token billing, no idle cost, and the widest regional model availability.

## When to Use

- Chat, completion, and tool-calling application backends
- Development and production alike (scale `capacity` in place)
- Anywhere you want Azure-wide routing rather than single-region capacity

## Key Configuration Choices

- **GlobalStandard** -- per-token billing with Azure-wide routing; capacity is a rate limit in thousands of tokens-per-minute
- **Version unset** -- tracks the model's current default version; pin `model.version` + `versionUpgradeOption: NO_AUTO_UPGRADE` for compliance workloads
- **The deployment name is the API contract** -- applications pass "chat", not the model name; the model behind it can change freely

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-cognitive-account-id>` | ARM ID of the parent account (kind OpenAI or AIServices) | `AzureCognitiveAccount` status outputs (`cognitive_account_id`), or reference it with valueFrom |
