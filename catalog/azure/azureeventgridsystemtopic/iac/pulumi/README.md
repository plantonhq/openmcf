# AzureEventgridSystemTopic Pulumi Module

## Overview

Creates an Azure Event Grid system topic -- the subscription surface for events Azure itself publishes about one of your resources (a storage account's blob events, a resource group's lifecycle events, a Key Vault's secret expiries). Azure allows one system topic per source resource per topic type; event subscriptions attach to it to route events to handlers.

## Resources Created

- `eventgrid.SystemTopic` -- the system topic (source binding, optional managed identity)

## Outputs

- `system_topic_id` -- the system topic's ARM resource ID (the target an event subscription's `system_topic_id` references)
- `system_topic_name` -- the system topic's name
- `metric_resource_id` -- the GUID-style Azure Monitor identifier for the topic's metrics
- `identity_principal_id` -- the system-assigned identity's principal (empty when no system-assigned identity)

## Behavior Notes

- **Everything but identity and tags is create-only**: name, region, resource group, source, and topic type replace the topic when changed -- and a replaced topic drops every subscription attached to it.
- **One system topic per source resource per topic type** -- a second create against the same source fails at deploy time; teams sharing a source share its system topic.
- **The region must match the source's region**; global sources (subscriptions via `Microsoft.Resources.Subscriptions`, resource groups via `Microsoft.Resources.ResourceGroups`) require `Global`.
- **The identity supports the combined mode** (`SystemAssigned, UserAssigned`) -- unlike the Event Grid publisher kinds.
- **The bridged SDK still carries the deprecated `sourceArmResourceId` alias**; the module uses the v5 name (`sourceResourceId`), which both engines render identically.
- **Billing**: free at rest; per-operation pricing on deliveries.

## Required Permissions

The deployer permissions this module needs are cataloged in [`../permissions.yaml`](../permissions.yaml).
