# Overview

The **AzureEventgridSystemTopic** component deploys an Azure Event Grid system topic -- the subscription surface for events Azure itself publishes about one of your resources. A storage account announcing blob creations, a resource group announcing resource writes, a Key Vault announcing secret expiries: the system topic is where those built-in event streams become subscribable.

## Purpose

- **Turn Azure services into publishers**: no code changes, no polling -- the platform already emits the events; the system topic exposes them.
- **One topic per source, shared by every consumer**: Azure allows one system topic per source resource per topic type, so teams attach their own subscriptions to the shared topic.
- **Identity for secured delivery**: system-assigned, user-assigned, or both at once -- subscriptions deliver as the topic's identity to locked-down destinations.

## Key Features

- Full azurerm v5 surface: source binding (any supported Azure resource by ARM ID), topic type, managed identity (including the combined mode), tags.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, identity IDs to AzureUserAssignedIdentity; the `system_topic_id` output is the wiring edge an AzureEventgridEventSubscription's `system_topic_id` references.
- Region contract taught in-line: the topic follows the source's region; global sources (subscriptions, resource groups) use `Global`.

## Use Cases

- **Blob-triggered pipelines**: subscribe to a storage account's BlobCreated events and fan them into Functions or queues.
- **Resource governance**: watch a resource group's lifecycle events for audit trails or drift detection.
- **Secret rotation automation**: react to a Key Vault's near-expiry events before credentials lapse.

## Future Enhancements

- Delivery wiring lives in AzureEventgridEventSubscription -- point its `system_topic_id` at this component's ID output.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
