# AzureEventgridSystemTopic Terraform Module

## Overview

Creates an Azure Event Grid system topic -- the subscription surface for events Azure itself publishes about one of your resources (a storage account's blob events, a resource group's lifecycle events, a Key Vault's secret expiries). Azure allows one system topic per source resource per topic type; event subscriptions attach to it to route events to handlers.

## Resources Created

- `azurerm_eventgrid_system_topic` -- the system topic (source binding, optional managed identity)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridSystemTopicSpec fields; the resource group, source, and identity references arrive as resolved literals

## Outputs

- `system_topic_id` -- the system topic's ARM resource ID (the target an event subscription's `system_topic_id` references)
- `system_topic_name` -- the system topic's name
- `metric_resource_id` -- the GUID-style Azure Monitor identifier for the topic's metrics
- `identity_principal_id` -- the system-assigned identity's principal (empty when no system-assigned identity)
