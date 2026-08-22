# CloudflareD1Database guide

Operational judgment for D1 databases. The README covers what each field is; this covers how the pieces interact.

## Placement is fixed at creation

Both the region hint and the jurisdiction are creation-time decisions: changing either replaces the database, which destroys its data. Pick placement before the database holds anything real, and treat a placement change in a plan as the destructive event it is.

## Region and jurisdiction are two answers to one question

Both fix where the primary lives — region for latency, jurisdiction for data residency — so the spec forbids setting both. If a residency regime applies ("eu", "fedramp"), jurisdiction wins and the region choice belongs to Cloudflare within that boundary.

## Read replication changes the Worker's contract

Enabling `read_replication: auto` is not free performance: a Worker reading a replicated database must use the D1 Sessions API to get sequential consistency. Enable it together with the application change, not ahead of it.

## Schema lives with the application, not here

This kind provisions the container. Tables, indexes, and migrations are applied by Wrangler (or the application's own migration step) against the database id this kind outputs — never bake schema into infrastructure.
