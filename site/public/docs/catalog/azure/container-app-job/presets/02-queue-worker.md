---
title: "Event-Triggered Queue Worker"
description: "This preset creates a queue-draining job: KEDA polls an Azure Storage Queue every 30 seconds, and when the queue holds more than 5 messages it starts executions (up to 10 at once). Each execution..."
type: "preset"
rank: "02"
presetSlug: "02-queue-worker"
componentSlug: "container-app-job"
componentTitle: "Container App Job"
provider: "azure"
icon: "package"
order: 2
---

# Event-Triggered Queue Worker

This preset creates a queue-draining job: KEDA polls an Azure Storage Queue every 30 seconds, and when the queue holds more than 5 messages it starts executions (up to 10 at once). Each execution processes its batch and exits -- the queue-worker model where executions ARE the scaling unit.

## When to Use

- Message and task processing where each batch should run to completion
- Bursty workloads where zero queue depth should mean zero cost
- Work that needs per-execution isolation (a poison message fails one execution, not a long-lived worker)

## Key Configuration Choices

- **Event trigger with azure-queue scaler** (`eventTrigger.scale.rules`) -- Queue depth drives execution fan-out; any KEDA scaler (Kafka, Service Bus, RabbitMQ) works the same way
- **Execution cap** (`maxExecutions: 10`) -- The cost and concurrency ceiling under backlog pressure
- **10-minute deadline** (`replicaTimeoutInSeconds: 600`) -- Bounded batch size; size it to the slowest legitimate batch
- **Two retries** (`replicaRetryLimit: 2`) -- Transient failures re-run before the execution is marked failed
- **System-assigned identity** -- Pulls the image keylessly; grant it roles instead of managing credentials

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the environment's) | Your environment's configuration |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<storage-account-connection-string>` | Connection string KEDA uses to read queue depth | `AzureStorageAccount` status outputs (primary_connection_string) |

## Related Presets

- **01-scheduled-batch** -- Use instead when work fires on a fixed schedule rather than queue pressure
- **03-on-demand** -- Use instead when executions are started manually or by a CI/CD system
