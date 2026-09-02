# CloudflareLogpushJob guide

The judgment this guide protects you from: the entitlement refusal wears a misleading message, two of this resource's fields cannot be changed after the job exists, and the destination handshake deliberately routes through storage only you can read.

## "Exceeded max jobs allowed" means NO ENTITLEMENT, not too many jobs

On a plan without the Logpush entitlement, creating a job fails with 403 code 1004: "creating a new job (for http_requests dataset) is not allowed: exceeded max jobs allowed" -- measured on zones carrying ZERO jobs, on Free and Pro plans alike. The quota it refers to is the plan's job allowance, which is zero below Enterprise for most datasets (http_requests included). If you see this message on a fresh zone, the fix is the plan entitlement, not deleting jobs. The destination side is NOT what is gated: destination validation (including R2 with embedded credentials) succeeds on any plan.

## The job variant is locked at creation -- and the plan will not warn you

`kind` (normal batching vs `edge` Instant Logs) carries no replacement guard in the provider: changing it plans a clean in-place update, and Cloudflare then rejects the apply with HTTP 400 (the provider's own immutable-fields test records exactly this). Pick the variant when you create the job; changing it later means delete and recreate. `dataset` and the scope are honest about being immutable: changing either replaces the job in the plan.

`destination_conf`, despite being the field that LOOKS most locked-down, updates in place fine -- the provider's own tests change it three times over. Two cautions still apply: a NEW destination must prove ownership before logs flow (see below), and credentials rotated *inside* the URI count as an update, so rotate them deliberately rather than as a side effect of an unrelated apply.

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
