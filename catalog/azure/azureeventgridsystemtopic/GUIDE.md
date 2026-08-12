# Azure Event Grid System Topic -- Operational Guide

Judgment calls that matter when you run system topics in production.

## The topic is a singleton claim on its source

Azure allows exactly one system topic per source resource per topic type. In a multi-team subscription that makes the topic SHARED INFRASTRUCTURE: the first team to create it owns its resource group placement and lifecycle, and every other team attaches subscriptions to it. Decide deliberately where it lives (a platform-owned group beats an app team's group that might get deleted), and never let two deployment pipelines race to create topics on the same source -- the loser fails with a conflict, not a merge.

## Deleting the topic deletes every subscription on it

The topic's create-only surface (source, type, name, region) means innocent-looking changes are replacements -- and a replaced topic silently drops all attached subscriptions, including other teams'. Treat topic replacement like a coordinated migration: enumerate subscriptions first (`az eventgrid system-topic event-subscription list`), recreate them after.

## Region is inherited, not chosen

The topic must sit in its source's region -- there is no latency or placement decision to make, just a contract to satisfy. The one wrinkle is global sources: subscriptions and resource groups emit from Azure's control plane, so their topics take `Global`. A region mismatch fails at deploy time with a clear error; the spec comment carries the rule so it never gets that far.

## Enable identity before the subscriptions need it

Subscriptions deliver AS the topic's identity when targeting locked-down destinations (a storage queue with shared keys disabled, an identity-secured Event Hub). The identity must exist on the topic BEFORE such a subscription is created, and its data-plane grants on the destination must land first too. The combined mode (system + user assigned) exists precisely for migrations -- keep the old identity granted while the new one rolls out.

## Not every resource type emits every event

The topic type names a catalog entry, not a promise that your scenario is covered -- e.g. storage accounts emit blob and queue events but not table events. Check the source type's event schema (`az eventgrid topic-type list-event-types --name <type>`) before designing a pipeline around an event that does not exist.
