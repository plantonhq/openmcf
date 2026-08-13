# Team Workspace

This preset deploys a standard team workspace with a system-assigned identity on the three required companion services -- the starting point for most ML estates: public access on (the default), no network isolation, Basic tier.

## When to Use

- A team's first ML workspace, development or production
- Estates without private-networking requirements
- The foundation under datastores, compute, and (later) model endpoints

## Key Configuration Choices

- **System-assigned identity** -- the simple path; Azure creates and rotates it with the workspace. Grant it storage/vault access after creation (or switch to user-assigned identities to grant beforehand)
- **Companion services referenced, not embedded** -- swap the `value:` placeholders for `valueFrom` references when the storage account, key vault, and insights component are Planton-managed
- **No managed-network block** -- isolation stays disabled and is read back; move to the Private Hardened Workspace preset for locked-down estates (the tightening works in place; loosening back does not)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | The resource group to create the workspace in | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-application-insights-id>` | ARM ID of the Application Insights component | `AzureApplicationInsights` status outputs (`application_insights_id`), or reference it with valueFrom |
| `<your-key-vault-id>` | ARM ID of the Key Vault | `AzureKeyVault` status outputs (`key_vault_id`), or reference it with valueFrom |
| `<your-storage-account-id>` | ARM ID of the storage account (general-purpose, HNS off) | `AzureStorageAccount` status outputs (`storage_account_id`), or reference it with valueFrom |

The `name` field carries a realistic example (`acme-ml-team`) because the workspace name is pattern-validated -- replace it with your own 3-33 character name.
