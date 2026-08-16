# DigitalOcean Database Cluster -- Operational Guide

Judgment calls that matter when you run managed databases on DigitalOcean.

## Pick the engine slug, not the marketing name

The `engine` values are DigitalOcean's own API slugs: PostgreSQL is `pg`, never `postgres`. Redis and Valkey are separate slugs for the same caching product line — DigitalOcean treats them as interchangeable and migrates Redis clusters toward Valkey; new caches should start on `valkey`.

## Node count is an engine decision

- **PostgreSQL / MySQL / MongoDB**: 1 node for dev, 3 for production failover. 2 buys a standby without quorum; most teams go straight to 3.
- **Redis / Valkey**: 1 node is normal — caches tolerate a failover gap. Add standbys only when cache warm-up is expensive.
- **Kafka**: 3 is the floor; DigitalOcean rejects less.
- **OpenSearch**: 1–15; go multi-node when the index must survive node loss rather than for query speed.

DigitalOcean enforces these server-side; the spec only enforces the universal minimum of 1 so new engine rules never require a contract change.

## Storage: plan for growth, not shrinkage

`storageGib` only ever grows. Two practical rules:

- Prefer `storageAutoscale` over hand-managed increments — DigitalOcean grows the disk at your threshold with a one-hour cooldown.
- If you grow `sizeSlug` while `storageGib` is unset, the cluster adopts the new slug's (larger) default storage automatically. A stale explicit `storageGib` smaller than the new slug's default is invalid — unset it when upsizing.

Note `storageAutoscale` currently deploys through the Terraform provisioner only; the Pulumi bridge rejects it loudly (no silent drop). Choose the provisioner accordingly or manage growth manually on Pulumi stacks.

## Upgrades are one-way and live

Raising `engineVersion` performs an in-place major upgrade on the running cluster. There is no downgrade, and no blue-green: take a backup-restore copy first if the application's compatibility is unproven. Region changes are similar — a live migration, not a recreate — expect elevated latency while it runs.

## VPC placement is create-only

`vpc` cannot be changed after the cluster exists. Decide network placement first; retrofitting means a new cluster plus a data migration (`backupRestore` gives you the copy).

## Restoring from backup

`backupRestore.databaseName` names the SOURCE cluster; omit `backupCreatedAt` to take the newest backup. The block acts only at creation and is never reported back by DigitalOcean — it is provisioning input, not ongoing configuration.

## Connection strings: public vs private

`connection_uri`/`host` traverse the public internet (TLS-required); `private_uri`/`private_host` resolve only inside the cluster's VPC. Applications in the same VPC should always use the private pair — lower latency and no public exposure. The default user's password is in both URIs; treat them as secrets.

## Eviction policy semantics

`evictionPolicy` values mirror Redis maxmemory policies with underscores (`allkeys_lru`, `volatile_ttl`, ...). Removing the field from a cluster that had one resets the policy to `noeviction` — it does not "keep the last value".

## What is deliberately NOT here

Users, logical databases, connection pools, read replicas, firewall (trusted-sources) rules, per-engine config parameters, Kafka topics, and log sinks are separate DigitalOcean resources with independent lifecycles. Manage them as their own resources rather than expecting cluster fields.
