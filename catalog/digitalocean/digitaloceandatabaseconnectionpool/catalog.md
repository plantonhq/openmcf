# Connection Pool on DigitalOcean

Creates a PgBouncer connection pool on a DigitalOcean managed PostgreSQL cluster -- the endpoint applications connect to when direct connections would exhaust the cluster's limit. Supports transaction, session, and statement pooling modes, dedicated or inbound-user authentication, and any logical database on the cluster. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Connection Pool** -- a named PgBouncer pool on the referenced PostgreSQL cluster, with its own endpoint port
- **Pool Identity** -- configured only when `user` is set; the pool authenticates as that cluster user. Omitted creates DigitalOcean's inbound-user pool where clients bring their own credentials

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A DigitalOceanDatabaseCluster** -- the owning PostgreSQL cluster, referenced by name (or an existing cluster's UUID as a literal). Pools exist only on PostgreSQL clusters.

### DigitalOcean Account

- Nothing beyond the cluster: pools are a free feature of managed PostgreSQL.

## After You Deploy

Wire applications to the pool's `uri`/`private_uri` secret outputs (note the pool's own port -- it differs from the cluster's). Remember every field is create-only: a later size or mode change replaces the pool and drops its live connections, so schedule such edits like brief connection outages.
