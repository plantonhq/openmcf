# CloudflareZeroTrustMcpPortal guide

The judgment this guide protects you from: the portal is the AUDIENCE layer -- curation decisions belong here or at the server, and putting them in the wrong layer makes future audiences inherit the wrong defaults.

## Server-level vs portal-level curation

Disable a tool at the SERVER (`CloudflareZeroTrustMcpServer.updated_tools`) when nobody should ever call it through Cloudflare. Disable it at the PORTAL when this audience shouldn't. A tool disabled at the server cannot be re-enabled by any portal -- the server is the floor, portals only subtract.

## The slug is the URL is the identity

`portal_id` is immutable and (like the hostname) part of what clients configure. Renaming the slug replaces the portal -- every AI client configured against it breaks. Pick audience-stable slugs ("eng-tools", not "q3-pilot").

## The hostname must be yours

The portal is HOSTED by Cloudflare on your zone (unlike MCP servers, whose hostnames point at external upstreams), so the portal hostname must belong to a domain the account serves. A hostname on a domain the account does not own is rejected at create with API error 7012 "invalid_domain" (live-measured).

## Overrides only work after the server has synced

Cloudflare validates every prompt/tool override against the server's actually-synced inventory: an override naming a tool the server has not synced is rejected at write with error 7001 "Tool 'x' does not exist on server" (live-measured). Register the server, let its background sync succeed, then add overrides -- an override written in the same change as a brand-new server usually races the sync and fails.

## Rows are a set; some overrides are write-only

Server rows carry no order -- the backend returns its own canonical order, and reordering manifest rows never plans a change. Within a row's prompt/tool overrides, `alias` and `description` are WRITE-ONLY per portal (the API returns them under different names, so the provider never reads them back): they cannot drift, and they land empty on import -- re-assert from configuration.

## Adopting an existing portal

Import restores the portal's scalars but NOT its configured server rows: the GET returns server rows in a different shape than the config writes (no `server_id`, server-echoed defaults), so the first post-adoption apply re-asserts the declared rows -- an idempotent replace of the portal's server list, not a real change (live-measured; the import recipes tolerate exactly this).

## on_behalf is a trust decision

With `on_behalf: true` (Cloudflare's default), Cloudflare authenticates to the upstream FOR the user -- one upstream credential, held by Cloudflare, shared by every portal user. Set it false when the upstream must know WHICH user is calling (per-user audit or per-user entitlements upstream).

## Pairs well with

- [CloudflareZeroTrustMcpServer](../cloudflarezerotrustmcpserver/README.md) -- the registrations this portal publishes.
- [CloudflareZeroTrustOrganization](../cloudflarezerotrustorganization/README.md) -- the login experience in front of the portal.
