# AzureEventgridDomain Terraform Module

## Overview

Creates an Azure Event Grid domain -- one publishing endpoint and one pair of access keys serving many event streams (domain topics), the multi-tenant pattern. The domain's name becomes a public DNS hostname (`{name}.{region}.eventgrid.azure.net`), unique across all Azure customers in the region.

## Resources Created

- `azurerm_eventgrid_domain` -- the domain (endpoint, access keys, topic-lifecycle flags, optional managed identity, optional IP firewall)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureEventgridDomainSpec fields; the resource group and identity references arrive as resolved literals

## Outputs

- `domain_id` -- the domain's ARM resource ID
- `domain_name` -- the domain's name
- `endpoint` -- the HTTPS publish endpoint (shared by every topic in the domain)
- `primary_access_key` / `secondary_access_key` -- the SAS keys (sensitive)
- `identity_principal_id` -- the system-assigned identity's principal (empty when no identity)

## Behavior Notes

- **Create-only surfaces**: name, region, resource group, input schema, and both input-mapping blocks replace the domain when changed -- and a replaced domain gets a new endpoint hostname, so publishers must be repointed.
- **Domain-topic lifecycle**: `auto_create_topic_with_first_subscription` / `auto_delete_topic_with_last_subscription` are always sent (platform defaults mirror Azure's auto-managed posture). Set both false for the pinned-topics governance posture, declaring each topic as an `AzureEventgridDomainTopic` resource.
- **Input mappings only mean something on `CustomEventSchema`** -- the module sends a mapping block only when it carries at least one field.
- **`public_network_access_enabled` / `local_auth_enabled` are always sent** (platform defaults mirror Azure's `true`). Local auth inverts to ARM's `disableLocalAuth` inside the provider.
- **Inbound IP rules are attribute-mode**: the provider clears rules on update by sending an empty list, so rule deletions propagate exactly.
- **The IP rule action vocabulary is `Allow` only** at azurerm v5 -- deny rules do not exist on this resource.
- **Identity is SystemAssigned XOR UserAssigned** -- no combined mode on this resource.
- **Billing**: free at rest; per-operation pricing on publishes and deliveries.

## Required Permissions

The deploying principal needs `Microsoft.EventGrid/domains/*` (EventGrid Contributor covers it). Reading the access keys requires `Microsoft.EventGrid/domains/listKeys/action`.
