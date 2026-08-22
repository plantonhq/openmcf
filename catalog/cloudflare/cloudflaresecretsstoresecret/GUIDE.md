# CloudflareSecretsStoreSecret guide

The judgment this guide protects you from: the value only ever travels one way, the name is the contract with consumers, and the scopes list has exactly one legal shape.

## The value is write-only -- your secret store is the system of record

Cloudflare never returns a secret's value: not on read, not on import, not on refresh. Consequences worth internalizing: drift on the value is UNDETECTABLE (a dashboard edit behind IaC's back wins silently until the next apply overwrites it); an import lands with an empty value that must be re-asserted from configuration; and losing the source value means the secret can only be replaced, never recovered. Keep the value in a managed secret (`value` takes a reference) and treat that as the truth.

## Scopes are alphabetical, or they drift forever

The API returns scopes sorted alphabetically no matter how you send them, and the provider models the list as ordered -- an out-of-order config re-plans a change on every run, forever (the provider's own tests document it). This spec walls the order at validation time: `[access, ai_gateway, dex, workers]` is the full canonical sequence; list the subset you need in that order. If validation rejects your list, sort it -- do not fight the wall, it exists to keep your plans clean.

## Rotation is the whole point

Consumers reference the secret by store and NAME. Rotate by changing `value` (an in-place update) -- every Worker binding and gateway pick it up on next read, no consumer redeploys. Renaming, by contrast, REPLACES the secret under a new ID and breaks consumers that referenced the old name: treat `name` like an API contract.

## Scope minimally

`scopes` is an access-control list, not a formality: a secret scoped `[workers]` cannot be read by AI Gateway, and vice versa. Grant the surfaces that actually read it and nothing more -- widening later is a plain in-place update.

## Pairs well with

- [CloudflareSecretsStore](../cloudflaresecretsstore/README.md) -- the one-per-account vault this secret lives in.
- [CloudflareWorker](../cloudflareworker/README.md) -- secrets-store bindings that read the secret at runtime.
- [CloudflareAiGateway](../cloudflareaigateway/README.md) -- BYO-keys authentication backed by secrets like this one.
