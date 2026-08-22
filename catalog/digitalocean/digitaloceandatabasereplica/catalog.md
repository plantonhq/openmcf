# Database Read Replica on DigitalOcean

Creates a single-node read-only replica of a DigitalOcean managed database cluster -- in the primary's region for read scaling or in a different region for geographically local reads -- with optional VPC placement and custom storage. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster and VPC dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Read-Only Replica** -- a full single-node managed database of the configured size, continuously following the primary
- **VPC Network Attachment** -- configured only when `vpc` is provided; places the replica's private endpoint in the named VPC (the REPLICA region's VPC for cross-region replicas)
- **Custom Storage** -- configured only when `storageSizeMib` is provided; must stay at or above the primary's storage
- **DigitalOcean Tags** -- your `tags` plus resource metadata tags applied automatically -- note replica tags are create-only upstream (a retag replaces the replica)

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A DigitalOceanDatabaseCluster** -- the primary (PostgreSQL or MySQL), referenced by name (or an existing cluster's UUID as a literal).

### DigitalOcean Account

- **Budget for a second node** -- the replica bills hourly like a single-node cluster of its slug, from creation, regardless of read traffic.

## After You Deploy

Point read-heavy consumers at the replica's own endpoint outputs (`host`/`port`, or the `uri`/`private_uri` secrets). Size and storage can grow in place later; everything else -- including tags -- replaces the replica, which reseeds from the primary without touching primary data.
