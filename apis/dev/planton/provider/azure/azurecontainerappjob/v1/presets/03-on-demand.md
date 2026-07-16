# On-Demand Job (Manual Trigger)

This preset creates a manually triggered job -- the database-migration model. The job definition lives in the environment ready to run; executions start on demand from the CLI (`az containerapp job start`), an SDK call, or a CI/CD pipeline step, and each runs a single replica to completion.

## When to Use

- Database migrations gated on a deploy pipeline step
- One-off operational tasks that should run inside the environment's network boundary
- Anything a human or pipeline decides to run, rather than a schedule or event

## Key Configuration Choices

- **Manual trigger** (`manualTrigger`) -- Executions start only when asked; `parallelism: 1` + `replicaCompletionCount: 1` is the run-one-copy-to-success shape
- **30-minute deadline** (`replicaTimeoutInSeconds: 1800`) -- A hung migration is terminated rather than blocking forever
- **Entrypoint override** (`command` + `args`) -- The image's default entrypoint is replaced with the migration command
- **Key Vault secret via system identity** -- The database credential is read just-in-time; nothing to rotate in the job definition

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the environment's) | Your environment's configuration |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |

## Related Presets

- **01-scheduled-batch** -- Use instead when the work should fire on a fixed schedule
- **02-queue-worker** -- Use instead when work arrives as queue messages
