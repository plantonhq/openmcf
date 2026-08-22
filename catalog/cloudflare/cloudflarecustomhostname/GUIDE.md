# CloudflareCustomHostname guide

Operational judgment for SaaS custom hostnames. The README covers what each field is; this covers how the pieces interact.

## Ownership proof is the customer's job

Onboarding a hostname does not make it live. Cloudflare returns TXT (and sometimes HTTP) ownership-verification records; the customer must create them on *their* DNS (or serve the HTTP body) before the hostname activates. The stack outputs those records so a chart or a ticket can hand them over. Until then the hostname sits in `pending` / `pending_validation`.

## SSL is a nested lifecycle

The hostname resource owns its SSL settings. Changing validation method, certificate type, or swapping a custom cert is an SSL-settings edit, not a hostname replace — but a custom (BYO) certificate is Enterprise-gated and will 403 on a non-Enterprise account even though the proto accepts it. Keep BYO certs off the default path.

## The fallback origin is a different kind

Traffic for a custom hostname that has no more-specific origin goes to the zone's fallback origin (CloudflareCustomHostnameFallbackOrigin). That singleton is per-zone, not per-hostname. Set it once for the SaaS zone; do not try to express it on each hostname.

## Soft-delete still answers GET

A hostname in `pending_deletion` or `deleted` is still readable. An automation that treats "the id exists" as "the hostname is live" will lie after a destroy. Wait for a real 404 (or an absent-status) before considering it gone.
