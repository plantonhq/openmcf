# CloudflareZeroTrustMcpPortal guide

The judgment this guide protects you from: the portal is the AUDIENCE layer -- curation decisions belong here or at the server, and putting them in the wrong layer makes future audiences inherit the wrong defaults.

## Server-level vs portal-level curation

Disable a tool at the SERVER (`CloudflareZeroTrustMcpServer.updated_tools`) when nobody should ever call it through Cloudflare. Disable it at the PORTAL when this audience shouldn't. A tool disabled at the server cannot be re-enabled by any portal -- the server is the floor, portals only subtract.

## The slug is the URL is the identity

`portal_id` is immutable and (like the hostname) part of what clients configure. Renaming the slug replaces the portal -- every AI client configured against it breaks. Pick audience-stable slugs ("eng-tools", not "q3-pilot").

## Rows are a set; some overrides are write-only

Server rows carry no order -- the backend returns its own canonical order, and reordering manifest rows never plans a change. Within a row's prompt/tool overrides, `alias` and `description` are WRITE-ONLY per portal (the API returns them under different names, so the provider never reads them back): they cannot drift, and they land empty on import -- re-assert from configuration.

## on_behalf is a trust decision

With `on_behalf: true` (Cloudflare's default), Cloudflare authenticates to the upstream FOR the user -- one upstream credential, held by Cloudflare, shared by every portal user. Set it false when the upstream must know WHICH user is calling (per-user audit or per-user entitlements upstream).

## Pairs well with

- [CloudflareZeroTrustMcpServer](../cloudflarezerotrustmcpserver/README.md) -- the registrations this portal publishes.
- [CloudflareZeroTrustOrganization](../cloudflarezerotrustorganization/README.md) -- the login experience in front of the portal.
