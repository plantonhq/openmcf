# AI Foundry Account

This preset creates an `AIServices`-kind account with project management enabled and a system-assigned identity -- the foundation AI Foundry team workspaces (AzureCognitiveAccountProject) and agents are built on. It carries the full multi-service AI surface, model deployments included.

## When to Use

- Teams organizing AI work into AI Foundry projects (agents, evaluations, files)
- Platform teams provisioning one account many teams share through isolated projects
- Users of the retired `azurerm_ai_services` resource -- this account is its migration target

## Key Configuration Choices

- **Kind AIServices + projectManagementEnabled** -- the pairing projects require; validation also requires the identity block
- **System-assigned identity** -- created with the account; its principal ID surfaces in the outputs for Key Vault / storage grants
- **Custom subdomain at creation** -- keeps token auth and network hardening open

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | Resource group for the account | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `acme-foundry-prod` | Example account name (also the subdomain) -- replace with your own globally unique name | Your naming convention |
