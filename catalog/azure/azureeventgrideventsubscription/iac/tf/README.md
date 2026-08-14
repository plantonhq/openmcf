# AzureEventgridEventSubscription Terraform Module

## Overview

Creates an Azure Event Grid event subscription -- the delivery instruction that routes events from a source to a handler, with filtering, retry, and dead-letter behavior. The spec's addressing choice selects which provider resource materializes: `scope` creates `azurerm_eventgrid_event_subscription` (attaches to any ARM resource by id -- a custom topic, domain, domain topic, resource group, or subscription), while `system_topic_id` creates `azurerm_eventgrid_system_topic_event_subscription` (a child of the system topic, addressed by resource group and topic name parsed from the referenced id). The two provider resources share one configuration grammar -- the module's two bodies are identical by design.

## Resources Created

Exactly one of:

- `azurerm_eventgrid_event_subscription` -- the scope-addressed subscription
- `azurerm_eventgrid_system_topic_event_subscription` -- the system-topic subscription

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridEventSubscriptionSpec fields; scope, system topic, destination, identity, and dead-letter references arrive as resolved literals

## Outputs

- `event_subscription_id` -- the subscription's ARM resource ID (shape follows the addressing choice)
- `event_subscription_name` -- the subscription's name
