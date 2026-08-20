# CloudflareZeroTrustMcpServer guide

The judgment this guide protects you from: the credentials are WRITE-ONLY and the identity is IMMUTABLE -- both are safety features that surprise you exactly once.

## Write-only credentials: IaC is the rotation path

`auth_credentials` and `client_secret` are stored AES-GCM-encrypted at Cloudflare and never returned by any read. Consequences: out-of-band rotation (dashboard, API) is INVISIBLE to IaC -- state keeps the old value and no drift ever shows; and import lands with empty credentials that must be re-asserted from configuration. Rotate by changing the managed-secret value and re-applying -- the modules send the full body on update, so the new value lands.

## The immutable trio

`server_id`, `hostname`, and `auth_type` all force replacement. A replacement gets a NEW registration -- portals referencing the old `server_id` must be re-pointed, and users lose in-flight OAuth grants. Upgrading a server's auth from bearer to oauth is therefore a migration (new registration, portal re-point, old registration retired), not an edit.

## Sync status is information, not health

After registration, Cloudflare syncs the upstream server's prompt/tool inventory in the background. A `status` of "error" means the upstream did not answer -- common for servers behind ACLs that only admit Cloudflare's egress after DNS/network setup. The registration object is real either way; fix reachability and the next sync recovers. Never gate deployment success on sync state.

## Overrides are allow-list thinking

`updated_tools`/`updated_prompts` rows with `enabled: false` remove sharp tools (deletes, writes) from what users can invoke -- and the portal layer can disable more per audience. Disable at the SERVER for "nobody should ever call this through Cloudflare"; disable at the PORTAL for "this audience shouldn't".

## Pairs well with

- [CloudflareZeroTrustMcpPortal](../cloudflarezerotrustmcpportal/README.md) -- the user-facing endpoint publishing this server.
- [CloudflareZeroTrustGatewaySettings](../cloudflarezerotrustgatewaysettings/README.md) -- Gateway filtering for upstream traffic (secure_web_gateway).
