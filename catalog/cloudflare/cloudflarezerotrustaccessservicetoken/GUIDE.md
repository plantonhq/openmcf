# CloudflareZeroTrustAccessServiceToken guide

The judgment this guide protects you from: the client secret is shown twice in the token's life -- at create, and at each rotation -- and never again. If you did not capture it, the token still exists and policies still match its ID, but no client can authenticate until you rotate.

## Capture the secret in the same change that creates the token

`client_secret` is a computed, sensitive output. Cloudflare returns it on create and when `client_secret_version` changes; Read always overwrites the API response with the value already in state. Import lands with an empty secret. The operational contract is: write `status.outputs.client_id` and `status.outputs.client_secret` to your secret store before the apply is considered done. Logs must never print the secret.

## Rotation is a pair, not a version bump

`client_secret_version` and `previous_client_secret_expires_at` are both-or-neither. Setting one without the other is rejected here (the provider's AlsoRequires pair). Leaving both unset is the normal non-rotating state -- Cloudflare treats the initial secret as version 1 (the field defaults to 1 when omitted).

To rotate:

1. Increment `client_secret_version` (1 → 2 on the first rotation).
2. Set `previous_client_secret_expires_at` to when the OLD secret should die. Extend it into the future so services can migrate; set it in the past to kill a compromised secret immediately.
3. Apply. Capture the **new** `client_secret` output. The `service_token_id` and `client_id` stay the same, so Access policy `service_token` rules keep matching.

Do not destroy-and-recreate to rotate. That mints a new token ID and every policy that listed the old ID stops matching.

## Duration vs rotation

`duration` is how long the token itself is valid (`8760h` default, or `forever`). Rotation is how you replace the secret inside that lifetime. A `forever` token still needs rotation when the secret leaks.

## API-token auth on the rotation path

The provider's own rotation acceptance tests unset `CLOUDFLARE_API_TOKEN` -- "the Access service does not yet support API tokens" for that path. Planton's Cloudflare harness authenticates with an API token. Create of a non-rotating token is the well-trodden path; a rotation apply may 403. If it does, record the defect and keep rotation as an offline-proven shape until the endpoint accepts tokens. Do not switch the harness to a legacy key to paper over it.

## Pairs well with

- [CloudflareZeroTrustAccessPolicy](../cloudflarezerotrustaccesspolicy/README.md) -- a `service_token` include rule that lists this token's ID.
- [CloudflareZeroTrustAccessApplication](../cloudflarezerotrustaccessapplication/README.md) -- the app the machine client calls.
- [CloudflareZeroTrustAccessIdentityProvider](../cloudflarezerotrustaccessidentityprovider/README.md) -- human sign-in, the other door.
