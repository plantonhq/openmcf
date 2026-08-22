# DigitalOcean Container Registry -- Operational Guide

Judgment calls that matter when you run a container registry on DigitalOcean.

## There is exactly one registry, so treat it as shared infrastructure

A DigitalOcean account holds ONE registry — not one per team, not one per environment. That makes this component account-level infrastructure with many consumers, like a VPC: one owner, one manifest, and everyone else references its outputs. Splitting environments happens inside the registry through repository naming (`registry.digitalocean.com/acme/staging-api`), never through a second registry. If a create fails with "already exists", the account already has its registry — import it rather than fighting the API.

## Default credentials live ~50 years; decide that on purpose

An unset `expirySeconds` mints a credential at the API maximum — roughly fifty years, effectively forever, and revocable only by deleting the credentials resource. That is sometimes right (a cluster's pull secret that must never rot) and often wrong (anything handled by humans). Set an explicit lifetime for anything outside automated pull paths, and remember rotation is a re-apply: changing `expirySeconds` re-mints the token in place.

## Read-only is the default for a reason

`write: false` credentials can pull but never push — the right shape for everything that RUNS images (Kubernetes clusters, Droplets, App Platform). Push access belongs to the build pipeline alone. If one consumer needs push and five need pull, resist widening the shared credential: mint the write credential here and distribute pull access separately (DOKS clusters can integrate with DOCR natively, without this credential at all).

## The credential knobs are unrecoverable — the manifest is the record

DigitalOcean never reports back how a credential was minted: `write` and `expirySeconds` exist only in your manifest and the provisioner state, and the credentials resource cannot be imported at the current provider pin (a fresh token would have to be minted just to observe one). Practical consequence: keep the manifest authoritative and never hand-mint credentials in the control panel alongside it — there is no way to reconcile the two later.

## Tier changes are live; everything else is a replacement

`subscriptionTier` moves up or down without touching stored images — downgrades just fail if you exceed the smaller tier's storage. `name` and `region` are create-only: changing either replaces the registry, and a replaced registry means every stored image is gone and every consumer's image reference breaks. Garbage collection of untagged images is an on-demand action in DigitalOcean's control panel/API, not a registry attribute — nothing in this manifest schedules it.
