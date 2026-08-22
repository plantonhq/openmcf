# CloudflareAccountApiToken guide

The judgment this guide protects you from: the token IS the credential, you get exactly one chance to capture it, and deleting it breaks whatever was using it that same second.

## One chance at the value -- after that, rotate

Cloudflare returns the token's secret in the create response and never again. No read, no refresh, and no import recovers it: an imported token arrives with configuration only and an empty value. So there is no such thing as "recovering" a lost token -- the recovery procedure IS rotation (delete and recreate, or roll it in the dashboard), followed by updating every consumer. Capture the value into a managed secret on the first apply and let consumers read it from there.

## Deleting the token is an immediate outage for its consumers

There is no grace period and no provider-side deletion guard. The moment the token is destroyed, every pipeline, script, and service using it starts getting 401s. Retire consumers first, or use the safer intermediate step: set `status: disabled`. That suspends the credential while keeping it on file, so you can watch what breaks and re-enable in seconds if you were wrong. Confirming a revocation with `disabled` before deleting is the habit worth building.

## Two grant shapes, and picking the wrong one silently over-grants

Each entry under `policies[].resources` is either a whole-resource grant or a nested scoping:

- `permission: "*"` on `com.cloudflare.api.account.<id>` grants across the ACCOUNT -- every zone it contains, present and future.
- `subresources: {"com.cloudflare.api.account.zone.*": "*"}` scopes into the account's zones, and naming a specific zone id instead of `*` narrows it to that zone.

Both are valid and both come straight from Cloudflare's own examples, which is exactly why the spec types them separately instead of accepting an opaque JSON string: the difference between "this account" and "the zones in this account" is the difference between a broad token and a narrow one, and it should be visible in review. Deny policies override allow policies covering the same resources -- useful for carving one zone out of an otherwise account-wide grant.

## Permission groups are UUIDs, and the list is the source of truth

`permission_group_ids` takes Cloudflare's own group identifiers, not names. Fetch them with `GET /accounts/{account_id}/tokens/permission_groups` (filterable by `name` and `scope`) and pin the UUIDs you need. This catalog deliberately does not model that read-only registry -- it is Cloudflare's data, it changes as products ship, and a copy here would rot.

## Order is not yours to control

Cloudflare canonically re-orders both policies and the permission groups inside them, and the provider carries hand-written machinery to match its response back to your configuration. Practical consequence: treat both lists as sets. Reordering them in the manifest is a no-op, and an imported token's policies come back in canonical order rather than the order you applied -- the provider's own import tests ignore the field for that reason.

## Statuses you can set, and statuses you can only observe

You set `active` or `disabled`. Cloudflare additionally reports `expired` and `revoked (exposed)` -- the latter when it detects a leaked token in the wild. Those are outcomes, not settings, and the provider drops such a token from state and recreates it on the next apply. If a token turns up `revoked (exposed)`, treat it as a real incident: the credential was published somewhere.

## The minting token needs its own permission

Creating tokens requires the credential doing the creating to carry Account API Tokens → Write. That is easy to forget when the rest of a pipeline's token is scoped to DNS or Workers, and the failure is a flat authorization error at create time.

## Pairs well with

- [CloudflareZeroTrustAccessServiceToken](../cloudflarezerotrustaccessservicetoken/README.md) -- machine credentials for Access-protected applications (a different trust domain from API tokens).
- [CloudflareSecretsStoreSecret](../cloudflaresecretsstoresecret/README.md) -- somewhere for the minted value to live so Workers can use it.
