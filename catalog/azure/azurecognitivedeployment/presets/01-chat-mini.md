# Chat (Mini Model)

This preset deploys a current mini-class chat model (`gpt-5.4-mini`) on the GlobalStandard SKU -- the standard starting point for chat and completion workloads: per-token billing, no idle cost, and the widest regional model availability.

## When to Use

- Chat, completion, and tool-calling application backends
- Development and production alike (scale `capacity` in place)
- Anywhere you want Azure-wide routing rather than single-region capacity

## Key Configuration Choices

- **GlobalStandard** -- per-token billing with Azure-wide routing; capacity is a rate limit in thousands of tokens-per-minute
- **Version unset** -- tracks the model's current default version; pin `model.version` + `versionUpgradeOption: NO_AUTO_UPGRADE` for compliance workloads
- **The deployment name is the API contract** -- applications pass "chat", not the model name; the model behind it can change freely
- **Model names age** -- Azure rejects models whose catalog lifecycle is "Deprecating" for NEW deployments (`ServiceModelDeprecating`), well before their final retirement date. Check the current catalog with `az cognitiveservices model list -l <region> --query "[?model.lifecycleStatus=='GenerallyAvailable'].model.name"` and swap the model here freely; existing deployments keep running.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-cognitive-account-id>` | ARM ID of the parent account (kind OpenAI or AIServices) | `AzureCognitiveAccount` status outputs (`cognitive_account_id`), or reference it with valueFrom |
