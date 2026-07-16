# Static Website Bucket

A publicly readable multi-region bucket serving a static site: index and
404 pages configured, CORS opened for the application origin, and public
access granted the additive way (one explicit `allUsers` reader).

## What this preset creates

A US multi-region bucket with website configuration and a single additive
`roles/storage.objectViewer` grant to `allUsers`. Uniform bucket-level
access stays on — public read comes from the IAM grant, never from legacy
object ACLs.

## Prerequisites

None — but the project's organization policy must permit public access
(`publicAccessPrevention: inherited` defers to it; the grant fails if the
org forbids public buckets).

## Composing HTTPS serving

The bucket's website endpoint is HTTP-only. For production HTTPS, compose
the L7 load-balancer family instead: point a `GcpBackendBucket` at this
bucket's `bucket_id` output, route it from a `GcpUrlMap`, and terminate
TLS on a `GcpTargetHttpsProxy` — CDN caching comes from the backend
bucket's `cdnPolicy`.

## Remix ideas

- Add a lifecycle rule deleting objects under a `previews/` prefix after
  30 days for PR-preview deployments.
- Set `forceDestroy: true` for ephemeral preview-site buckets that CI
  tears down.
