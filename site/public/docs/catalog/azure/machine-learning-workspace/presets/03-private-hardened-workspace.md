---
title: "Private Hardened Workspace"
description: "This preset deploys a locked-down workspace: public network access off, managed-network isolation in approved-outbound mode, and an explicit outbound allowlist (package indexes plus Key Vault). The..."
type: "preset"
rank: "03"
presetSlug: "03-private-hardened-workspace"
componentSlug: "machine-learning-workspace"
componentTitle: "Machine Learning Workspace"
provider: "azure"
icon: "package"
order: 3
---

# Private Hardened Workspace

This preset deploys a locked-down workspace: public network access off, managed-network isolation in approved-outbound mode, and an explicit outbound allowlist (package indexes plus Key Vault). The compliance posture for regulated estates.

## When to Use

- Regulated workloads where workspace egress must be enumerated
- Estates standardizing on private endpoints for data-plane access
- Any workspace handling data that must not exfiltrate through arbitrary egress

## Key Configuration Choices

- **`ALLOW_ONLY_APPROVED_OUTBOUND`** -- only the listed rules (plus Azure's own built-in required rules) are reachable; tightening into this mode works in place, loosening back out is rejected once the network provisions -- decide deliberately
- **`provisionOnCreationEnabled: true`** -- the managed VNet materializes at deploy time so quota/provisioning failures land early, not on the team's first job
- **Package-index FQDNs allowed** -- pip installs keep working inside the isolation boundary; extend the list per your toolchain (conda, CRAN, internal mirrors)
- **Access is via private endpoints** -- with public access off, plan the workspace's own private endpoints separately; a no-public-IP serverless compute here would also need a serverless subnet (validated at manifest time)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group-name>` | The resource group to create the workspace in | `AzureResourceGroup` status outputs (`resource_group_name`), or reference it with valueFrom |
| `<your-application-insights-id>` | ARM ID of the Application Insights component | `AzureApplicationInsights` status outputs (`application_insights_id`), or reference it with valueFrom |
| `<your-key-vault-id>` | ARM ID of the Key Vault | `AzureKeyVault` status outputs (`key_vault_id`), or reference it with valueFrom |
| `<your-storage-account-id>` | ARM ID of the storage account (general-purpose, HNS off) | `AzureStorageAccount` status outputs (`storage_account_id`), or reference it with valueFrom |

The `name` field carries a realistic example (`acme-ml-private`) because the workspace name is pattern-validated -- replace it with your own.
