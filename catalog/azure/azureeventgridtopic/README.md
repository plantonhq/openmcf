# Overview

The **AzureEventgridTopic** component deploys an Azure Event Grid custom topic -- the HTTPS endpoint an application publishes its own events to. Event subscriptions fan those events out to handlers (Functions, webhooks, queues, Event Hubs) with filtering and retry handled by Event Grid; the publisher fires one POST and moves on.

## Purpose

- **Decouple producers from consumers**: the application publishes facts ("order placed"); who reacts, how many handlers, and where they live changes without touching the publisher.
- **Push, not poll**: Event Grid delivers with per-subscription filtering, retry, and dead-lettering -- no consumer polling loops.
- **Secure the publish edge**: Entra ID or SAS keys, an IP firewall, and a public-access switch for private-endpoint-only topics.

## Key Features

- Full azurerm v5 surface: input schema (Event Grid, CloudEvents 1.0, or your own custom JSON with envelope mappings), inbound IP rules (up to 128), public-network and local-auth switches, system- or user-assigned managed identity, tags.
- Secure-by-default guidance: the locked-down preset ships local auth off (Entra-only publishing) with an IP allowlist.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, identity IDs to AzureUserAssignedIdentity; outputs carry the endpoint and both access keys (as secrets) for wiring publishers.

## Use Cases

- **Application integration events**: services publish domain events one POST at a time; subscribers filter by event type or subject prefix.
- **Cross-cloud eventing**: a CloudEvents-schema topic feeds handlers that also consume from other clouds' buses.
- **Custom-schema ingestion**: keep an existing JSON event shape and map its fields onto the envelope instead of rewriting producers.

## Future Enhancements

- Event subscriptions (the consumer side) arrive with the AzureEventgridEventSubscription kind, wiring handlers to this topic's ID output.
