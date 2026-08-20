# CloudflareLogpushJob guide

The judgment this guide protects you from: two of this resource's fields cannot be changed after the job exists, and the destination handshake deliberately routes through storage only you can read.

## The destination cannot be repointed -- recreate instead

`destination_conf` looks editable: the provider does not mark it for replacement, so a changed URI plans a clean in-place update. Cloudflare then rejects the request with HTTP 400 (the provider's own immutable-fields test records exactly this). To move a job to a new destination, delete it and create a new one -- and remember credentials rotated *inside* the URI count as a change. `dataset` is honest about being immutable: changing it replaces the job.

## Ownership proof passes through your bucket, on purpose

Most destinations must prove you own them before Cloudflare will write there. The flow has a manual middle step by design:

1. Set `generate_ownership_challenge: true` and deploy. Cloudflare drops a challenge file into the destination; the outputs report its name.
2. Read that file from your own storage and copy the token out of it.
3. Put the token in `ownership_challenge` and deploy again.

No API can do step 2 for you -- reading the file is what proves control. Same-account R2 destinations skip the handshake entirely, which is why the live proof lane uses one. The challenge record itself is one-shot: Cloudflare never deletes it, its plan output says so outright, and it cannot be imported or refreshed.

## A declared job ships logs

Cloudflare creates jobs DISABLED by default. That default makes a declared-but-idle job the easy accident, so this component sends `enabled = true` when you do not say otherwise. Set `enabled: false` explicitly to keep a paused job on file.

## Bound the volume before you turn it on

`http_requests` on a busy zone is a firehose, and the bill lands on your destination, not on Cloudflare. Two knobs bound it before the first byte ships: `filter` (a Cloudflare filter expression -- restrict to a hostname, a status class, a path prefix) and `output_options.sample_rate` (0.1 ships a tenth). Both are cheaper to set now than to discover on next month's storage invoice.

## Scope must match the dataset

Zone datasets (`http_requests`, `firewall_events`, `dns_logs`, `nel_reports`, `spectrum_events`, `page_shield_events`, `zaraz_events`) need `zone_id`; everything else needs `account_id`. The spec enforces exactly-one, but it cannot know which one a given dataset wants -- a mismatch fails at the API. Note the provider silently prefers the account when both are set, which is exactly the ambiguity this spec's exactly-one rule removes.

## Deleting the job is silent

Nothing warns you that log delivery stopped -- the objects already in the destination stay, so a dashboard glance looks fine. If a job matters for audit or compliance, pair it with a `CloudflareNotificationPolicy` on `failing_logpush_job_disabled_alert` so Cloudflare tells you when a job stops shipping.

## Pairs well with

- [CloudflareR2Bucket](../cloudflarer2bucket/README.md) -- the same-account destination that skips the ownership handshake.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- the scope for every zone dataset.
- [CloudflareNotificationPolicy](../cloudflarenotificationpolicy/README.md) -- alerting when a job fails or is disabled.
