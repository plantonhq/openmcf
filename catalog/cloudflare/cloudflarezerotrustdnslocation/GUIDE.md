# CloudflareZeroTrustDnsLocation guide

The judgment this guide protects you from: updates are FULL REPLACES at the API -- what the manifest omits does not merely stay unmanaged, it can actively reset.

## Omitting max_ttl on update resets it

Unlike the settings singletons elsewhere in Zero Trust, a location update sends the whole object: omitting `max_ttl` on an update RESETS the TTL behavior to inherit. If you manage a TTL override, keep it declared forever. The spec pairs `ttl_secs` with mode override by validation, so a manifest cannot get the pairing wrong.

## client_default is an account-level lever

`client_default: true` makes this location the attribution target for traffic from unregistered sources -- account-wide, one location at a time. Flipping it on a new location silently changes how every unattributed query is filtered. Treat it like changing a default route: deliberately, in a change window, never on a scratch location.

## Your first location is permanent (live-measured)

Cloudflare auto-promotes the FIRST DNS location ever created on an account to the account default -- even when `client_default` is unset. And the current default can neither be deleted (API error 1217 "Cannot delete default location") nor demoted in place (error 1216): the only way to move the default is promoting ANOTHER location with `client_default: true`. Two consequences: destroying a resource that holds the default FAILS until a replacement default exists, and an account that has ever had a location can never return to zero locations. Plan the first location on an account as a permanent fixture, not a scratch object.

## Declared endpoints want real network lists (provider defect)

At Terraform provider v5.23.0/v5.24.0, a declared endpoint (doh/dot/ipv6) whose `networks` list is EMPTY re-plans a cosmetic in-place update forever: Cloudflare drops empty lists on read, the config re-asserts them, and the plan never settles (a null list is worse -- it crashes the apply outright, which is why the modules always send known values). Declare real source networks on every endpoint you enable, or accept the permanent no-op diff until the provider models these lists as computed-aware types.

## Token-gated DoH is the roaming shape

The DoH subdomain is world-reachable by design; `require_token` on the doh endpoint rejects resolvers that merely discovered the URL. For office networks the source-network allowlists carry the gating; for roaming devices the token does.

## Never pin the shared destination pool

`dns_destination_ips_id` unset lets Cloudflare auto-assign the shared IPv4 destination pair. The docs' example UUID IS that shared pool -- copying it into a manifest pins what auto-assign would have done anyway, and blocks a future move to dedicated IPs. Set the field only for a dedicated (BYOIP/dedicated-resolver) mapping you actually hold.

## Pairs well with

- [CloudflareZeroTrustGatewayPolicy](../cloudflarezerotrustgatewaypolicy/README.md) -- the rules filtering this location's queries.
- [CloudflareZeroTrustGatewaySettings](../cloudflarezerotrustgatewaysettings/README.md) -- the account-wide Gateway posture.
