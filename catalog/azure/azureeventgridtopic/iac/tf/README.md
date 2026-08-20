# AzureEventgridTopic Terraform Module

## Overview

Creates an Azure Event Grid custom topic -- the HTTPS endpoint an application publishes its own events to, from which event subscriptions fan events out to handlers (Functions, webhooks, queues, Event Hubs). The topic's name becomes a public DNS hostname (`{name}.{region}.eventgrid.azure.net`), unique across all Azure customers in the region.

## Resources Created

- `azurerm_eventgrid_topic` -- the topic (endpoint, access keys, optional managed identity, optional IP firewall)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridTopicSpec fields; the resource group and identity references arrive as resolved literals

## Outputs

- `topic_id` -- the topic's ARM resource ID
- `topic_name` -- the topic's name
- `endpoint` -- the HTTPS publish endpoint
- `primary_access_key` / `secondary_access_key` -- the SAS keys (sensitive)
- `identity_principal_id` -- the system-assigned identity's principal (empty when no identity)

## Behavior Notes

- **Create-only surfaces**: name, region, resource group, input schema, and both input-mapping blocks replace the topic when changed -- and a replaced topic gets a new endpoint hostname, so publishers must be repointed.
- **Input mappings only mean something on `CustomEventSchema`** -- the module sends a mapping block only when it carries at least one field.
- **`public_network_access_enabled` / `local_auth_enabled` are always sent** (platform defaults mirror Azure's `true`). Local auth inverts to ARM's `disableLocalAuth` inside the provider.
- **Inbound IP rules are attribute-mode**: the provider clears rules on update by sending an empty list, so rule deletions propagate exactly.
- **The IP rule action vocabulary is `Allow` only** at azurerm v5 -- deny rules do not exist on this resource.
- **Identity is SystemAssigned XOR UserAssigned** -- the provider has no combined mode on this resource.
- **Billing**: free at rest; per-operation pricing on publishes and deliveries.

## Required Permissions

The deployer permissions this module needs are cataloged in [`../permissions.yaml`](../permissions.yaml).
