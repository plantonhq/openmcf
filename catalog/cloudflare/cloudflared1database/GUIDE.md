# CloudflareD1Database guide

Operational judgment for D1 databases. The README covers what each field is; this covers how the pieces interact.

## Placement is fixed at creation

Both the region hint and the jurisdiction are creation-time decisions. They fail differently when edited later, and both protect your data (live-verified 2026-08-26). Editing `region` does nothing: Cloudflare never returns the hint, so both engines deliberately ignore post-create changes — the alternative was a plan that destroys the database (data included) to move a hint. Editing `jurisdiction` still plans a REPLACE, which destroys the data — treat it as the destructive event it is. Either way: pick placement before the database holds anything real; changing it for real means creating a new database and migrating.

## Adopting an existing database

Import works by `{account_id}/{database_id}` and lands cleanly — including `read_replication`, which Cloudflare reports even when never configured. The one field an import cannot restore is the region hint (never returned by the API); it stays unmanaged after adoption, which is harmless because it is inert post-create anyway.

## Region and jurisdiction are two answers to one question

Both fix where the primary lives — region for latency, jurisdiction for data residency — so the spec forbids setting both. If a residency regime applies ("eu", "fedramp"), jurisdiction wins and the region choice belongs to Cloudflare within that boundary.

## Read replication changes the Worker's contract

Enabling `read_replication: auto` is not free performance: a Worker reading a replicated database must use the D1 Sessions API to get sequential consistency. Enable it together with the application change, not ahead of it.

## Schema lives with the application, not here

This kind provisions the container. Tables, indexes, and migrations are applied by Wrangler (or the application's own migration step) against the database id this kind outputs — never bake schema into infrastructure.
