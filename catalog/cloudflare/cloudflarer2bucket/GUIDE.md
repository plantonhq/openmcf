# CloudflareR2Bucket guide

Operational judgment for R2 buckets. The README covers what each field is; this covers how the pieces interact.

## Destroy requires an empty bucket

R2 refuses to delete a bucket holding objects, and the provider does no draining — a teardown of a non-empty bucket fails at the API with a raw error. Empty the bucket first (a lifecycle rule with a short expiry, or an explicit purge) and treat "who empties the bucket" as part of every decommission plan.

## Jurisdiction is identity, not configuration

The jurisdiction ("default", "eu", "fedramp") is fixed at creation and joins the import ID as a mandatory third segment — a default-jurisdiction bucket imports with the literal `default`. Every bucket-scoped companion (CORS, lifecycle, lock, notifications, domains) is created inside the same jurisdiction automatically.

## The location hint only counts once

`location` is honored the first time a bucket with that name is created and is best-effort even then. Changing it later changes nothing — the provider deliberately preserves the original placement rather than recreating the bucket.

## r2.dev is for development; custom domains are for production

The managed `r2.dev` public domain is rate-limited, and the provider cannot destroy or even re-read it once enabled — it stays until manually disabled in the dashboard. Serve production traffic from `custom_domains` on a zone in the same account, which also unlocks Cloudflare's cache and TLS controls.

## Object lock is a commitment

A lock rule retains matching objects for its whole window — nothing, including a bucket teardown, deletes them early. Locked objects also keep the bucket non-empty, which blocks destroy until the retention expires. Scope lock rules by prefix and keep retention as short as compliance allows.

## Event notifications ride the Queues entitlement

Notifications deliver into a Cloudflare Queue, so they inherit the queue's plan gating (see the CloudflareQueue guide) and the queue must exist before the notification that names it.
