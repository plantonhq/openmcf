# Overview

The **AzureEventgridNamespaceTopic** component deploys one named CloudEvents stream inside an Azure Event Grid namespace (AzureEventgridNamespace). A namespace holds many topics with independent lifecycles: publishers and teams create and delete their own streams against the shared hub without touching it or each other -- exactly like consumer groups on an Event Hub.

## Purpose

- **Self-service streams**: each service owns its topic; the namespace stays platform-owned.
- **Retention as the one dial**: events persist 1-7 days for delivery; everything else is fixed by Azure (CloudEvents v1.0 schema, Custom publisher type).
- **Chart-ready wiring**: `namespace_id` defaults its reference to AzureEventgridNamespace's ID output.

## Key Features

- Full azurerm v5 surface: the topic's three arguments (namespace, name, retention) modeled exactly.
- Azure pins two properties the provider does not expose -- the event schema is always CloudEvents v1.0 and the publisher type is always "Custom"; both engines send exactly those values.
- Retention is the topic's ONLY updatable property; name and namespace changes replace it.

## Use Cases

- **Per-service event streams**: one topic per publishing service inside a shared environment hub.
- **Short-lived pipelines**: retention tuned down to a day for high-volume, quickly-consumed streams.
- **Tenant onboarding**: create a stream per tenant without touching the namespace.

## Future Enhancements

- Azure models namespace-topic event subscriptions as their own resource; the pinned Terraform provider does not ship it yet -- delivery wiring for namespace topics arrives when it does.
