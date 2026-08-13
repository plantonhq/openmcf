# Azure Compute Gallery Image -- Operational Guide

Judgment calls that matter when you run gallery images in production.

## The identity is forever; get Gen2 + trusted launch right on day one

Publisher/offer/SKU, OS type, Hyper-V generation, and the security flags are all create-only -- changing any of them replaces the definition and orphans every published version. The posture that ages best for new Linux/Windows fleets: `hyperVGeneration: V2` with `trustedLaunchSupported: true` (consumers CHOOSE trusted launch) rather than `trustedLaunchEnabled` (consumers are FORCED into it -- right only when compliance demands it, because it excludes every non-trusted-launch consumer forever).

## Treat versions as immutable releases

A version's source, replication mode, and name are create-only; only its target regions, exclude-from-latest, end-of-life, and tags move in place. The workflow that stays honest: the image pipeline publishes a NEW version per build (semver names -- the `1.2.0` form), `excludeFromLatest` quarantines a bad release instantly without unpublishing it, and removing the entry unpublishes it once nothing pins it. Never rebuild a version in place under the same name.

## Replicas are deployment throughput, not durability

`regionalReplicaCount` sizes CONCURRENT VM creation from that region's copy (Azure's guidance: one replica per ~20 concurrent creations, scale-set bursts more). One replica per region is fine for trickle deployments; raise it where scale sets fan out. Every replica bills storage in its region -- prune target regions no fleet deploys from, and remember `storage_account_type` per region is create-only in practice (the API cannot update it; remove and re-add the region to change it).

## Shallow replication is a dev-loop tool with a leash

`replicationMode: Shallow` publishes instantly by REFERENCING the source instead of copying it -- excellent for the image-bake inner loop and for very large images. The leash: the source snapshot/blob must outlive the version (deleting it breaks deploys), Shallow versions cannot use per-region disk encryption sets, and the replica count is effectively 1. Production releases use Full, always.

## Clearing end-of-life dates replaces the resource

`endOfLifeDate` is advisory metadata on both the definition and each version, and it updates in place -- but CLEARING a previously set date forces replacement (the provider cannot express "remove this property" to ARM). Set it when a real retirement date exists; if the date moves, set a NEW date rather than clearing it.
