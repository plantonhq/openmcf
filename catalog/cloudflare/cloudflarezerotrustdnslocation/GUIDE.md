# CloudflareZeroTrustDnsLocation guide

The judgment this guide protects you from: updates are FULL REPLACES at the API -- what the manifest omits does not merely stay unmanaged, it can actively reset.

## Omitting max_ttl on update resets it

Unlike the settings singletons elsewhere in Zero Trust, a location update sends the whole object: omitting `max_ttl` on an update RESETS the TTL behavior to inherit. If you manage a TTL override, keep it declared forever. The spec pairs `ttl_secs` with mode override by validation, so a manifest cannot get the pairing wrong.

## client_default is an account-level lever

`client_default: true` makes this location the attribution target for traffic from unregistered sources -- account-wide, one location at a time. Flipping it on a new location silently changes how every unattributed query is filtered. Treat it like changing a default route: deliberately, in a change window, never on a scratch location.

## Token-gated DoH is the roaming shape

The DoH subdomain is world-reachable by design; `require_token` on the doh endpoint rejects resolvers that merely discovered the URL. For office networks the source-network allowlists carry the gating; for roaming devices the token does.

## Never pin the shared destination pool

`dns_destination_ips_id` unset lets Cloudflare auto-assign the shared IPv4 destination pair. The docs' example UUID IS that shared pool -- copying it into a manifest pins what auto-assign would have done anyway, and blocks a future move to dedicated IPs. Set the field only for a dedicated (BYOIP/dedicated-resolver) mapping you actually hold.

## Pairs well with

- [CloudflareZeroTrustGatewayPolicy](../cloudflarezerotrustgatewaypolicy/README.md) -- the rules filtering this location's queries.
- [CloudflareZeroTrustGatewaySettings](../cloudflarezerotrustgatewaysettings/README.md) -- the account-wide Gateway posture.
