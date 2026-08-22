# CloudflareSecretsStore guide

The judgment this guide protects you from: the store is a one-per-account singleton whose destruction takes every secret with it -- it is permanent infrastructure wearing a two-field spec.

## One per account: create once, adopt otherwise

Cloudflare allows a single Secrets Store per account. A create against an account that already has one (dashboard-created, or another team's) fails at the API. The correct move is adopt-by-import (`{account_id}/{store_id}`), never a second create. In an organization, decide early WHICH manifest owns the store -- everything else references it.

## Renaming is destruction

Both fields force replacement, and replacing the store destroys every secret inside it -- every Worker binding and AI Gateway key breaks at once. Pick a boring, permanent name ("account-secrets") and never touch it. If you truly must rename: create nothing until every secret has been re-created under the new store and every consumer re-pointed.

## Destroy last, if ever

The store's blast radius is every secret it holds. In teardown flows, destroy secrets and their consumers first; leave the store for the very end -- or leave it standing, it is free.

## Pairs well with

- [CloudflareSecretsStoreSecret](../cloudflaresecretsstoresecret/README.md) -- the secrets inside the store, each rotating on its own cadence.
- [CloudflareWorker](../cloudflareworker/README.md) -- secrets-store bindings that read secrets at runtime.
- [CloudflareAiGateway](../cloudflareaigateway/README.md) -- BYO-keys authentication backed by this store.
