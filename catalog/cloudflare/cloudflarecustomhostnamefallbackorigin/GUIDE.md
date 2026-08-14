# CloudflareCustomHostnameFallbackOrigin guide

Operational judgment for the zone fallback origin. The README covers what each field is; this covers how the pieces interact.

## One per zone, identity is the zone

There is no resource id. The API identity is the zone (`GET /zones/{zone_id}/custom_hostnames/fallback_origin`). Two manifests targeting the same zone are the same object — the second apply overwrites the first. Put exactly one of these in a SaaS zone's chart.

## Create equals update, and it is asynchronous

The write path is PUT. Create and update are the same call, so re-applying an identical spec is meant to be idempotent. Status then moves through `pending_deployment` before `active` — a just-applied origin is not yet serving. Read `status` (and `errors`) before assuming traffic flows.

## The origin hostname must already resolve

Cloudflare expects the origin hostname to exist as a DNS record on this zone (typically a proxied A/CNAME to the real backend). This kind does not create that record. If the DNS record kind is not yet deployed, the fallback origin can still be written but will sit pending or error until the name exists.

## Token credentials can 403 the write

Some Cloudflare API tokens are allowed to GET this singleton and forbidden to PUT it. A token-only harness that re-applies for idempotency may 403 even though the first create succeeded. If a write fails with 403 after a successful read, it is a credential scope problem, not a spec problem.
