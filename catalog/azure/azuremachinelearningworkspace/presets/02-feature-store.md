# Feature Store

This preset deploys a FEATURE_STORE-flavor workspace -- the managed feature store backing online/offline feature serving for ML pipelines. Same companion services as a regular workspace; the flavor changes what the workspace is for.

## When to Use

- Centralizing feature definitions shared across training and serving
- Estates adopting the Azure ML managed feature store
- Alongside (not instead of) regular workspaces -- feature stores serve features; regular workspaces train models

## Key Configuration Choices

- **`kind: FEATURE_STORE` + the featureStore block** -- the pairing is required in both directions (validated at manifest time, mirroring the provider's own contract)
- **Spark runtime pinned** -- the materialization compute's Spark version; drop the field to track the service default
- **Connection names unset** -- the offline/online store connections can be named later as the store's backing services land

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | The resource group to create the workspace in | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-application-insights-id>` | ARM ID of the Application Insights component | `AzureApplicationInsights` status outputs (`application_insights_id`), or reference it with valueFrom |
| `<your-key-vault-id>` | ARM ID of the Key Vault | `AzureKeyVault` status outputs (`key_vault_id`), or reference it with valueFrom |
| `<your-storage-account-id>` | ARM ID of the storage account (general-purpose, HNS off) | `AzureStorageAccount` status outputs (`storage_account_id`), or reference it with valueFrom |

The `name` field carries a realistic example (`acme-feature-store`) because the workspace name is pattern-validated -- replace it with your own.
