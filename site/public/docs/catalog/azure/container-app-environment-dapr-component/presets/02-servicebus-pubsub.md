---
title: "Service Bus Pub/Sub (Keyless)"
description: "This preset registers a Dapr pub/sub component backed by Azure Service Bus, authenticating with a managed identity instead of a connection string. Publisher and subscriber apps use Dapr's pub/sub API..."
type: "preset"
rank: "02"
presetSlug: "02-servicebus-pubsub"
componentSlug: "container-app-environment-dapr-component"
componentTitle: "Container App Environment Dapr Component"
provider: "azure"
icon: "package"
order: 2
---

# Service Bus Pub/Sub (Keyless)

This preset registers a Dapr pub/sub component backed by Azure Service Bus, authenticating with a managed identity instead of a connection string. Publisher and subscriber apps use Dapr's pub/sub API with the component name (`pubsub`); application code stays broker-agnostic, and no broker credential exists anywhere in the deployment.

## When to Use

- Event-driven communication between Container Apps (order placed -> billing + shipping)
- Decoupling producers from consumers with at-least-once delivery
- Teams standardizing on Dapr's pub/sub API so the broker can change without code changes

## Key Configuration Choices

- **Keyless authentication** (`namespaceName` + `azureClientId`) -- The apps' user-assigned identity authenticates to Service Bus via Entra; grant it the **Azure Service Bus Data Sender** and **Data Receiver** roles on the namespace. `azureClientId` references the identity's `client_id` output, so the wiring survives identity recreation. The connection-string alternative (a `connectionString` metadata entry backed by a component secret) still works but ships a SAS credential -- prefer keyless.
- **IaC owns the topology** (`disableEntityManagement: "true"`) -- Dapr consumes queues provisioned as first-class `AzureServiceBusQueue` resources and never creates entities itself. This is also what lets the identity's grant stay at Sender/Receiver instead of Data Owner.
- **Component type** (`pubsub.azure.servicebus.queues`) -- Queue-per-consumer pub/sub over first-class queues. Use `pubsub.azure.servicebus.topics` when consumers should fan out through Service Bus topic subscriptions instead; the metadata keys are identical.
- **Both sides scoped** (`scopes`) -- Publisher and subscriber app ids are listed; unlisted apps cannot use the component

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<service-bus-namespace>` | The Service Bus namespace name | `AzureServiceBusNamespace` spec `namespace_name` |
| `<workload-identity>` | The apps' user-assigned identity resource name | The `AzureUserAssignedIdentity` attached to the apps |
| `<publisher-app-id>` / `<subscriber-app-id>` | The consuming apps' dapr.app_id values | Each `AzureContainerApp`'s `dapr.app_id` field |

## Related Presets

- **01-blob-state-store** -- The state building block on Azure Blob Storage
