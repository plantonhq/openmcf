# Service Bus Pub/Sub

This preset registers a Dapr pub/sub component backed by Azure Service Bus topics. Publisher and subscriber apps use Dapr's pub/sub API with the component name (`pubsub`); Dapr manages the topics, subscriptions, and delivery -- application code stays broker-agnostic.

## When to Use

- Event-driven communication between Container Apps (order placed -> billing + shipping)
- Decoupling producers from consumers with at-least-once delivery
- Teams standardizing on Dapr's pub/sub API so the broker can change without code changes

## Key Configuration Choices

- **Component type** (`pubsub.azure.servicebus.topics`) -- Any Dapr pub/sub component works the same way (pubsub.kafka, pubsub.redis, ...); the metadata keys change per broker
- **Secret-backed connection string** (`metadata[].secretName` -> `secrets[]`) -- The broker credential travels as a component secret, never a plain metadata value
- **Both sides scoped** (`scopes`) -- Publisher and subscriber app ids are listed; unlisted apps cannot use the component

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<container-app-environment-id>` | ARM ID of the Container App Environment | `AzureContainerAppEnvironment` status outputs |
| `<service-bus-connection-string>` | Service Bus namespace connection string | `AzureServiceBusNamespace` status outputs |
| `<publisher-app-id>` / `<subscriber-app-id>` | The consuming apps' dapr.app_id values | Each `AzureContainerApp`'s `dapr.app_id` field |

## Related Presets

- **01-blob-state-store** -- The state building block on Azure Blob Storage
