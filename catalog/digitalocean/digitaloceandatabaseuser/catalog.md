# Database User on DigitalOcean

Creates an additional user on a DigitalOcean managed database cluster, with the MySQL authentication plugin choice and per-topic Kafka / per-index OpenSearch access-control lists. DigitalOcean generates the password (and Kafka mTLS certificate pair) server-side; they surface as secret stack outputs for application wiring. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Database User** -- a named user on the referenced cluster, with a server-generated password
- **MySQL Auth Plugin** -- configured only when `mysqlAuthPlugin` is set; chooses between DigitalOcean's modern default (`caching_sha2_password`) and the legacy plugin for old clients
- **Kafka / OpenSearch ACLs** -- configured only when `settings` is set; grants per-topic or per-index permissions from closed permission lists

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A DigitalOceanDatabaseCluster** -- the owning cluster, referenced by name (or an existing cluster's UUID as a literal).

### DigitalOcean Account

- Nothing beyond the cluster: users are a free feature of managed databases.

## After You Deploy

Read the `password` output (a secret) into your application configuration -- it never appears in the manifest. On Kafka clusters, wire the `access_cert` / `access_key` pair for mutual TLS. Remember the ACLs you declared here are the source of truth: DigitalOcean never reports them back.
