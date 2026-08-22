# Logical Database on DigitalOcean

Creates an additional logical database inside a DigitalOcean managed database cluster -- the namespace an application's tables live in, separate from other workloads sharing the same cluster. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Logical Database** -- a named database inside the referenced cluster, addressable by clients, users, and connection pools

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A DigitalOceanDatabaseCluster** -- the owning cluster, referenced by name (or an existing cluster's UUID as a literal).

### DigitalOcean Account

- Nothing beyond the cluster: logical databases are a free feature of managed databases.

## After You Deploy

Point application connection strings and DigitalOceanDatabaseConnectionPool resources at the `database_name`. Know the one hazard: both fields are create-only, so a rename later replaces the database and drops its data -- treat renames as data migrations, never edits.
