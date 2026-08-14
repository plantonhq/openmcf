# Azure Event Grid Namespace Topic -- Operational Guide

Judgment calls that matter when you run namespace topics in production.

## The topic is the tenant boundary, the namespace is the landlord

Namespace topics exist so teams and tenants onboard streams WITHOUT touching the shared hub. Keep it that way operationally: the namespace lives in a platform-owned resource group and pipeline; each topic lives with the service that publishes to it. Folding topic creation into the namespace's pipeline recreates the bottleneck the resource model exists to remove.

## Retention is a delivery buffer, not an archive

The 1-7 day window is how long Event Grid holds events for delivery -- it is not storage. If consumers can be down longer than the window, or you need replay beyond it, land the events somewhere durable (an Event Hub, a storage queue) instead of stretching retention.

## Deleting a topic is deleting its unread events

A topic delete (and any replace -- name and namespace are create-only) drops everything still buffered in it. Drain or stop publishers first; the delete itself succeeds regardless.

## Subscriptions on namespace topics live outside the catalog today

Azure models namespace-topic event subscriptions as their own ARM resource, and the pinned Terraform provider does not ship it yet. AzureEventgridEventSubscription addresses the CLASSIC resources only. Until the provider catches up, manage namespace-topic subscriptions with `az eventgrid namespace topic event-subscription` (or route through a classic topic via the namespace's MQTT bridge when that fits the shape).
