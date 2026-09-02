# AwsAppSyncApi — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Pick the arm by the client conversation

GraphQL is request/response with subscriptions bolted on; Events is pub/sub first. If clients ask questions about data, use the graphql arm; if clients broadcast and listen on channels (chat, presence, live dashboards), the events arm is simpler, cheaper, and has no schema to maintain. Migrating between them is a replacement — they are different AWS objects.

## The schema is applied async and never drift-checked

AppSync applies the schema through an asynchronous creation call, and the Terraform provider performs no drift detection on it: an out-of-band schema edit stays invisible until your next in-band change. Treat the manifest as the only writer. Schema errors surface at APPLY time, not plan — a failed apply names the offending SDL line.

## Managed type definitions are create-time text

AWS rewrites a type's SDL definition server-side into its own whitespace form (indentation stripped, braces re-placed, blank lines injected — live-caught), and the provider reads that rewritten form back with no suppression. Both engines therefore ignore textual drift on `definition` to keep plans converging — which means an in-place edit to a type's definition body does NOT propagate. To change a managed type, remove its entry, apply, and re-add it with the new text — the same replace-the-entry choreography its `format` field already requires. Types that evolve often belong in the schema SDL instead, where the whole document re-applies on any change.

## Resolver changes serialize per API

The provider takes a per-API lock around every resolver mutation and retries conflicts for two minutes. A manifest that changes thirty resolvers applies them one at a time — a big rollout is slow, not stuck. The same lock covers types and the schema.

## GraphQL names have no hyphens anywhere

The API's own name (`api_name`), data source names, function names, and type/field names all follow GraphQL's `[A-Za-z_][0-9A-Za-z_]*` charset. This is also what keeps the provider's hyphen-joined import IDs (`{api_id}-{type}-{field}`) unambiguous. The Events arm differs: its API name comes from metadata.name and channel namespace names allow interior hyphens.

## API keys: the secret shows once

AWS returns a key's secret only at creation; every read after that returns the key ID. The `api_key_ids` output carries IDs (never secrets) — fetch the secret from the console/CLI at creation time, store it in your secret manager, and rotate by overlap: add a key, roll clients, delete the old. Keys live at most 365 days; AWS rounds expiry down to the hour.

## The cache is a billed instance, not a flag

`graphql.cache` provisions a Redis-backed instance billed per hour while it exists — SMALL is the sensible default. Both encryption flags are decided at creation; changing either replaces the cache (a cold start, not an outage). PER_RESOLVER_CACHING only caches resolvers whose `caching` block opts in.

## Custom domains ride us-east-1 certificates

Like CloudFront, AppSync accepts only us-east-1 ACM certificates for custom domains, regardless of the API's region. After apply, point DNS at the `appsync_domain_name` output. The domain associates with exactly one API at a time; the association re-points in place, and both directions of that operation can take tens of minutes at AWS's side.

## EventBridge data sources: replace, don't edit

At this provider pin, updating an EventBridge data source silently drops its bus configuration (an upstream defect). The modules treat these as replace-to-change: rename the entry to force recreation instead of editing it in place.

## MERGED APIs own nothing

A merged API serves its source APIs' schemas — defining schema, types, functions, or resolvers on the merged surface is rejected (validation stops you before AWS does). AUTO_MERGE propagates source changes automatically; MANUAL_MERGE waits for you to trigger merges. The execution role needs appsync:SourceGraphQL plus start-merge permissions on the sources.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
