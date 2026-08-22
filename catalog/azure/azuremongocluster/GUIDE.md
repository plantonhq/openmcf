# Azure Mongo Cluster -- Operational Guide

Judgment calls that matter when you run Mongo vCore clusters in production.

## vCore or RU? Answer it before anything else

Azure sells two MongoDB-compatible products. This one (vCore) is a REAL MongoDB engine on dedicated compute: predictable hourly cost, full query surface, community drivers and tools work unchanged. The other (AzureCosmosdbAccount with the Mongo API) is Cosmos's request-unit model: serverless economics, per-operation pricing, a compatibility layer rather than the engine. Lift-and-shift MongoDB workloads and anything using change streams, transactions, or aggregation-heavy queries belong HERE; spiky tiny workloads that would idle a vCore belong there.

## The mode you create with is the mode you die with

Azure never returns `create_mode` on reads, so the provider replaces the cluster on ANY mode change -- there is no "promote this replica" or "re-parent this clone" through configuration. Promotion of a GeoReplica is an Azure-side operation (portal/CLI); after promoting, import the promoted cluster rather than editing the replica's manifest into Default mode (that is a replacement that would destroy it).

## Free and M25 are sandboxes with walls

The two burstable tiers reject zone-redundant HA and sharding past one shard (enforced at manifest time here), and upgrades AWAY from them stage a tier-first update the provider performs itself -- expect two update waves in one apply. One Free cluster per subscription is Azure's own cap. The Free tier also refuses Microsoft Entra ID authentication outright -- ARM rejects the create with "Microsoft Entra ID authentication is not supported for 'Free' cluster tier" (enforced at manifest time here), so a sandbox that needs passwordless app sign-in via AzureMongoClusterUser grants starts at M10. Never let a Free-tier proof-of-concept quietly become production: the tier ceiling arrives as throttling, not as an error.

## Password rotation is an update; username is forever

`administrator_password` rotates in place (reference a secret store so rotation is a reference change, not a manifest edit). `administrator_username` is create-only. The rotation-friendly posture for applications is to not use the administrator at all: grant each app's identity an AzureMongoClusterUser and keep the admin credential for break-glass.

## Storage grows, never shrinks

`storage_size_in_gb` moves up in place; there is no path down short of dump-and-restore into a new cluster. Size for the working set, not the someday set -- growth is a one-line change later.

## The connection strings substitute real credentials

Azure returns connection strings with a `<user>:<password>` placeholder; the engines substitute the actual administrator credentials into the outputs. Treat `connection_string`/`connection_strings` as secrets end to end (they are marked sensitive on every surface) and wire consumers by reference so a password rotation propagates.

## Adopting an existing cluster: expect one convergence apply

Azure never returns `create_mode` or the admin password on a read, so a freshly imported cluster's state cannot carry them. The first apply after an import therefore shows one in-place update re-asserting exactly those two values from your manifest -- that is expected convergence, not drift (setting the same password again is a cloud-side no-op), and it never proposes a replace. From the second apply on, the plan is clean.
