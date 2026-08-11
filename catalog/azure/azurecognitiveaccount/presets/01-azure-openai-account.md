# Azure OpenAI Account

This preset creates an `OpenAI`-kind S0 account -- the container Azure OpenAI model deployments (gpt-4o, embeddings) are created onto. The account object itself carries no idle cost; billing follows deployment usage.

## When to Use

- The starting point for any Azure OpenAI workload
- Teams that deploy and call models (AzureCognitiveDeployment rides on this)
- Accounts that may later grow into AI Foundry (`OpenAI` upgrades to `AIServices` in place)

## Key Configuration Choices

- **Kind OpenAI, SKU S0** -- the standard pairing for model deployments
- **Custom subdomain set at creation** -- required for Entra ID auth and network ACLs; changing it later replaces the account
- **Region matters** -- Azure OpenAI model availability differs per region; deploy where your models exist

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the account | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `acme-openai-prod` | Example account name (also the subdomain) -- replace with your own globally unique name | Your naming convention; alphanumeric start, then alphanumerics, periods, dashes, underscores |
