# Overview

The **AzureEventgridDomainTopic** component deploys one named event stream inside an Azure Event Grid domain -- the per-tenant mailbox of the multi-tenant pattern. Publishers address it by naming the topic in events sent to the DOMAIN's shared endpoint; subscribers attach event subscriptions to the topic's own ARM ID.

## Purpose

- **Make tenant onboarding an explicit act**: with the domain's auto-create/auto-delete flags off, each stream exists only as a declared resource -- auditable, chart-managed, reviewable.
- **Give subscriptions their target**: event subscriptions scope to a domain topic's ID; this kind's output is that wiring edge.
- **Keep lifecycles independent**: tenants join and leave as topic creates and destroys, without touching the domain or each other.

## Key Features

- The full azurerm v5 surface (the resource is pure addressing: the parent domain and a name).
- Chart-ready: `domain_id` defaults its reference to AzureEventgridDomain's `domain_id` output, so a domain and its topics compose in one manifest set.
- Free at rest and creates in seconds -- operations bill on the domain.

## Use Cases

- **SaaS tenant onboarding**: a chart creates the tenant's domain topic (and its subscriptions) when the tenant signs up, and destroys them when the tenant leaves.
- **Pinned team streams**: a platform domain where each producing team's stream is declared, not materialized as a side effect.

## Future Enhancements

- Event subscriptions (the consumer side) arrive with the AzureEventgridEventSubscription kind, attaching to this kind's `domain_topic_id` output.
