# Overview

The **AzureEventgridEventSubscription** component deploys an Azure Event Grid event subscription -- the delivery instruction that completes the eventing story: "when events arrive at THIS source, filtered like THIS, deliver them THERE". It is the consumer side of every Event Grid kind in the catalog: custom topics, domains, domain topics, and system topics all become useful the moment a subscription routes their events to a handler.

## Purpose

- **Route events to handlers**: seven destination types -- Azure Functions, Event Hubs, Service Bus queues and topics, storage queues, relay hybrid connections, and HTTPS webhooks.
- **Filter at the platform**: event types, subject prefix/suffix, and a 19-operator field-level filter grammar keep noise away from handlers.
- **Deliver reliably**: tunable retry, storage-blob dead-lettering for events that exhaust retries, and managed-identity delivery to locked-down destinations.

## Key Features

- Full azurerm v5 surface across BOTH provider resources: one kind addresses any scope (custom topic, domain, domain topic, resource group, subscription) or a system topic, and the engines create the matching resource type.
- Chart-ready: `scope` defaults its reference to AzureEventgridTopic, `system_topic_id` to AzureEventgridSystemTopic, destination arms to AzureEventHub / AzureServiceBusQueue / AzureServiceBusTopic / AzureStorageAccount.
- Secrets handled properly: static delivery-property values (API keys handlers expect) are sensitive inputs referencing managed secrets, never plaintext.

## Use Cases

- **Queue fan-out**: a topic's events land in a storage queue a worker pool drains -- the cheapest at-least-once consumer.
- **Blob pipelines**: a storage system topic's BlobCreated events trigger Functions per upload.
- **Filtered webhooks**: only `orders/` events with an amount above a threshold reach the partner's HTTPS endpoint.

## Future Enhancements

- The Event Grid namespace family (MQTT, pull delivery) arrives as its own kinds; this subscription covers the classic push-delivery model end to end.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
