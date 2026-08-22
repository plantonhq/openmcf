# AzureEventgridTopic Pulumi Module

## Overview

Creates an Azure Event Grid custom topic -- the HTTPS endpoint an application publishes its own events to, from which event subscriptions fan events out to handlers (Functions, webhooks, queues, Event Hubs). The topic's name becomes a public DNS hostname (`{name}.{region}.eventgrid.azure.net`), unique across all Azure customers in the region.

## Resources Created

- `eventgrid.Topic` -- the topic (endpoint, access keys, optional managed identity, optional IP firewall)

## Outputs

- `topic_id` -- the topic's ARM resource ID
- `topic_name` -- the topic's name
- `endpoint` -- the HTTPS publish endpoint
- `primary_access_key` / `secondary_access_key` -- the SAS keys (sensitive)
- `identity_principal_id` -- the system-assigned identity's principal (empty when no identity)

## Behavior Notes

- **Create-only surfaces**: name, region, resource group, input schema, and both input-mapping blocks replace the topic when changed -- and a replaced topic gets a new endpoint hostname, so publishers must be repointed.
- **Input mappings only mean something on `CustomEventSchema`** -- the module sends a mapping block only when it carries at least one field, mirroring the Terraform module.
- **`public_network_access_enabled` / `local_auth_enabled` are always sent** (platform defaults mirror Azure's `true`). Local auth inverts to ARM's `disableLocalAuth` inside the provider.
- **Inbound IP rules are built unconditionally** so an emptied list clears the rules on update, exactly like the Terraform module.
- **The IP rule action vocabulary is `Allow` only** at the pinned provider line -- deny rules do not exist on this resource.
- **Identity is SystemAssigned XOR UserAssigned** -- no combined mode on this resource.
- **Billing**: free at rest; per-operation pricing on publishes and deliveries.

## Required Permissions

The deployer permissions this module needs are cataloged in [`../permissions.yaml`](../permissions.yaml).
