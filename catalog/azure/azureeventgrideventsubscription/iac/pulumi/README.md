# AzureEventgridEventSubscription Pulumi Module

## Overview

Creates an Azure Event Grid event subscription -- the delivery instruction that routes events from a source to a handler, with filtering, retry, and dead-letter behavior. The spec's addressing choice selects which provider resource materializes: `scope` creates `eventgrid.EventSubscription` (attaches to any ARM resource by id -- a custom topic, domain, domain topic, resource group, or subscription), while `system_topic_id` creates `eventgrid.SystemTopicEventSubscription` (a child of the system topic, addressed by resource group and topic name parsed from the referenced id). The two provider resources share one configuration grammar -- the module's two builders are identical by design.

## Resources Created

Exactly one of:

- `eventgrid.EventSubscription` -- the scope-addressed subscription
- `eventgrid.SystemTopicEventSubscription` -- the system-topic subscription

## Outputs

- `event_subscription_id` -- the subscription's ARM resource ID (shape follows the addressing choice)
- `event_subscription_name` -- the subscription's name

## Behavior Notes

- **Addressing, name, and delivery schema are create-only** -- changing any of them replaces the subscription.
- **Exactly one destination arm** (spec-enforced): Azure Function, Event Hub, hybrid connection, Service Bus queue, Service Bus topic, storage queue, or HTTPS webhook.
- **Webhook creates perform a validation handshake** -- the endpoint must be live and answer it, or the create fails.
- **Delivery properties are ignored by Azure on storage-queue destinations** (queue messages carry no custom properties); static values marked `secret` return "Hidden" from Azure's read APIs.
- **`retry_policy` is sent only when set** -- Azure's defaults (30 attempts / 1440 minutes) echo back on read otherwise.
- **The advanced filter caps at 25 total values per subscription** (Azure's service limit) -- the Terraform engine rejects an overflow at plan time on scope-addressed subscriptions; this engine surfaces Azure's own deploy-time rejection.
- **Engine-shape note**: the bridged SDK still names the id-arm destinations `eventhubEndpointId` / `hybridConnectionEndpointId` / `serviceBusQueueEndpointId` / `serviceBusTopicEndpointId`, where the pinned Terraform provider renamed them to `eventhub_id` / `hybrid_connection_id` / `service_bus_queue_id` / `service_bus_topic_id` at v5. Both engines write the identical ARM destination object.
- **Billing**: free at rest; per-operation pricing on deliveries.

## Required Permissions

The deployer permissions this module needs are cataloged in [`../permissions.yaml`](../permissions.yaml).
