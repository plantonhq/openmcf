# AzureEventgridNamespaceTopic Terraform Module

## Overview

Creates one named CloudEvents stream inside an Azure Event Grid namespace. Many topics share one namespace with independent lifecycles -- like consumer groups on an Event Hub.

## Resources Created

- `azurerm_eventgrid_namespace_topic` -- the stream (namespace binding, name, delivery retention)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridNamespaceTopicSpec fields; the namespace reference arrives as a resolved literal

## Outputs

- `namespace_topic_id` -- the topic's ARM resource ID (`{namespace_id}/topics/{name}`)
- `namespace_topic_name` -- the topic's name

## Behavior Notes

- **Azure pins the schema and publisher type**: CloudEvents v1.0 and "Custom" -- the provider sends both; neither is configurable.
- **Retention is the only updatable property** (1-7 days, platform default 7 via `coalesce`); name and namespace changes replace the topic -- and drop its buffered events.
- **No tags**: the provider exposes none on this resource.
