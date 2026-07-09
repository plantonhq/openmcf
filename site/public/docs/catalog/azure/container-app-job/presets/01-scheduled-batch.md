---
title: "Scheduled Batch Job"
description: "This preset creates a nightly batch job that runs on a cron schedule inside a Container App Environment. Each execution runs a single replica to completion (up to a 1-hour deadline), pulls its image..."
type: "preset"
rank: "01"
presetSlug: "01-scheduled-batch"
componentSlug: "container-app-job"
componentTitle: "Container App Job"
provider: "azure"
icon: "package"
order: 1
---

# Scheduled Batch Job

This preset creates a nightly batch job that runs on a cron schedule inside a Container App Environment. Each execution runs a single replica to completion (up to a 1-hour deadline), pulls its image from a private registry via managed identity, and reads its database connection string from Key Vault -- no plaintext credentials anywhere.

## When to Use

- Nightly reports, data exports, or aggregation runs
- Periodic cleanup and retention enforcement
- Any run-to-completion work that fires on a fixed schedule (UTC cron)

## Key Configuration Choices

- **Schedule trigger** (`scheduleTrigger.cronExpression: "0 2 * * *"`) -- Executions start at 02:00 UTC daily; standard five-field cron format
- **1-hour deadline** (`replicaTimeoutInSeconds: 3600`) -- A replica running longer is terminated and counted as failed
- **One retry** (`replicaRetryLimit: 1`) -- A failed replica gets one more attempt before the execution is marked failed
- **Key Vault secret** (`keyVaultSecretId` + `identity`) -- The connection string never lives in the job definition; the identity reads it just-in-time
- **Managed-identity registry pull** (`registries[].identity`) -- No registry password to rotate

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the environment's) | Your environment's configuration |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<your-user-assigned-identity-id>` | ARM ID of the user-assigned identity | `AzureUserAssignedIdentity` status outputs |

## Related Presets

- **02-queue-worker** -- Use instead when work arrives as queue messages rather than on a schedule
- **03-on-demand** -- Use instead when executions are started manually or by a CI/CD system
