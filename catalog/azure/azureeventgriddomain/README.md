# Overview

The **AzureEventgridDomain** component deploys an Azure Event Grid domain -- one publishing endpoint and one pair of access keys serving MANY event streams (domain topics). It is the multi-tenant pattern: a SaaS publishes every customer's events to the same endpoint, naming the topic per event, and each customer subscribes only to their own topic.

## Purpose

- **One endpoint, thousands of streams**: publishers integrate once; tenants come and go as domain topics without new endpoints, keys, or firewall changes.
- **Tenant isolation on the consumer side**: subscriptions attach per domain topic, so one tenant's handlers never see another tenant's events.
- **Choose the topic lifecycle**: auto-managed (topics materialize with their first subscription -- Azure's default) or pinned (topics exist only as declared AzureEventgridDomainTopic resources -- the governance posture).

## Key Features

- Full azurerm v5 surface: input schema (Event Grid, CloudEvents 1.0, or custom JSON with envelope mappings), the auto-create/auto-delete topic lifecycle flags, inbound IP rules (up to 128), public-network and local-auth switches, system- or user-assigned managed identity, tags.
- Secure-by-default guidance: the pinned-topics preset turns both lifecycle flags off so topics exist only by decision.
- Chart-ready: `resource_group` defaults its reference to AzureResourceGroup, identity IDs to AzureUserAssignedIdentity; outputs carry the shared endpoint and both access keys (as secrets), and `domain_id` is the reference target domain topics wire to.

## Use Cases

- **Multi-tenant SaaS eventing**: one domain, one topic per customer, per-customer subscriptions delivering to per-customer handlers.
- **Team-partitioned event bus**: one domain per platform, one topic per producing team, subscriptions crossing team boundaries by choice.
- **High-fanout ingestion**: a single publisher spreading events across many streams without managing many endpoints.

## Future Enhancements

- Event subscriptions (the consumer side) arrive with the AzureEventgridEventSubscription kind, attaching per domain topic via AzureEventgridDomainTopic's ID output.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
