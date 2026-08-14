# Azure Event Grid Domain Topic

Deploys one named event stream (domain topic) inside an Azure Event Grid domain -- the per-tenant mailbox of the multi-tenant pattern. It integrates with Planton's Provider Connections for Azure credential management and ValueFromRef for dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Event Grid domain topic** -- a named stream under the domain (`{domain_id}/topics/{name}`), the scope event subscriptions attach to

## Before You Deploy

### Planton Setup

- **Azure Provider Connection** -- an active connection in the Connect module with credentials for the target Azure subscription.
- **An Event Grid domain** -- reference an AzureEventgridDomain's `domain_id` output (or provide a literal domain ARM ID).

### Azure Subscription

- **A domain topic has no endpoint or keys of its own** -- publishers POST to the DOMAIN's endpoint and name this topic in the event; the domain's auth settings govern publishing.
- **Everything is create-only** -- the topic is pure addressing; a name change replaces it (briefly interrupting its subscriptions, nothing else).
- **Pinning versus auto-managing**: an explicitly declared topic persists regardless of the domain's `auto_create`/`auto_delete` flags; the full governance posture sets both flags false on the domain.
- **Free** -- operations are billed on the domain.

## Deploy

### Console

Open the deployment store, find **Azure Event Grid Domain Topic**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Tenant Stream** preset in the [Presets](#presets) tab.

### CLI

```bash
planton apply -f eventgrid-domain-topic.yaml
```

## After Deploy

The `domain_topic_id` output is the scope event subscriptions attach to -- create the tenant's subscriptions against it next (events published to a topic with no subscription are dropped, not queued). The topic appears on the domain's **Topics** blade, and a test event naming the topic shows up on the domain's **Metrics** blade.
