# AWS CloudFront

Amazon CloudFront is AWS's global content delivery network: 400+ edge locations that terminate TLS close to viewers, cache responses, and route requests back to origins. This document describes how Planton models CloudFront distributions: the origin/behavior composition, how private S3 serving works, and the operational characteristics that shape deployments.

## Distribution identity

CloudFront distributions have no name in AWS — identity is the generated ID (`E2ABC...`). `metadata.name` drives the `Name` identity tag, and consumers compose through stack outputs: `domain_name` + `hosted_zone_id` for Route53 alias records, `distribution_arn` for WAF associations and S3 bucket policies, `distribution_id` for cache invalidations.

## The origin/behavior composition

A distribution is a routing table over content sources:

- **Origins** declare where content comes from. Each has a user-chosen `origin_id` — the stable handle everything else targets — and one type arm:

| Arm | For | Notes |
|-----|-----|-------|
| `s3_origin` | S3 buckets via the REST endpoint | Pair with Origin Access Control to keep the bucket private |
| `custom_origin` | Anything HTTP | Load balancers, API endpoints, S3 *website* endpoints (which speak plain HTTP only) |
| `vpc_origin` | Private resources in your VPC | Requires a provisioned CloudFront VPC origin |

- **Origin groups** pair two origins into primary/failover with configurable trigger status codes.
- **Cache behaviors** declare how requests are matched and handled. The default behavior catches everything; ordered behaviors match path patterns first (`/api/*`, `*.jpg`), in list order. Each behavior targets one origin or group by ID — validation proves every target resolves, so a typo cannot reach AWS.

## Private S3 serving: OAC over OAI

Origin Access Control is the modern mechanism (SigV4 signing; works with SSE-KMS in all regions) and what the spec creates when asked (`create_origin_access_control: true`). The distribution's ARN must then be allowed in the bucket policy — `cloudfront.amazonaws.com` principal with an `AWS:SourceArn` condition. The legacy Origin Access Identity is accepted by path for buckets already wired to one, but never created.

## The two caching generations

- **Modern (recommended)**: policy IDs. A cache policy owns the cache key and TTLs; an origin-request policy owns what is forwarded without joining the cache key; a response-headers policy manipulates response headers. AWS-managed policies cover most needs — `Managed-CachingOptimized` for static content, `Managed-CachingDisabled` for APIs, `Managed-AllViewer` for full forwarding.
- **Legacy**: the inline `forwarded_values` block plus per-behavior TTLs, where everything forwarded joins the cache key. Kept for existing configurations; each behavior must choose exactly one generation (AWS rejects both and neither).

## Edge compute

- **CloudFront Functions** — sub-millisecond JavaScript at all 400+ edges, viewer events only, no network access. URL rewrites, redirects, header manipulation.
- **Lambda@Edge** — full Lambda at regional edge caches, all four events (viewer-request/response run on every request with a 5s limit; origin-request/response run on cache misses with 30s). Requires a numbered function version ARN in us-east-1.

## The us-east-1 gravity

CloudFront is global, but its control plane lives in us-east-1: viewer certificates (ACM) must be issued there, CLOUDFRONT-scope WAF Web ACLs must be created there, and Lambda@Edge functions must be published there — regardless of where origins live.

## Operational characteristics

- **Slow propagation**: every create/update pushes configuration to every edge location — typically 5-15 minutes. `wait_for_deployment` (default true) blocks until `Deployed`; disabling it returns early while edges converge.
- **Deletion requires disabling first**: AWS refuses to delete an enabled distribution. The IaC engines handle the disable-then-delete dance automatically, which makes destroys as slow as deploys.
- **Cache invalidation is an operation, not configuration**: prefer content-hashed filenames; invalidate by distribution ID when needed.

## Stack outputs

| Output | Description |
|--------|-------------|
| `distribution_id` | Invalidations and monitoring key on it |
| `distribution_arn` | WAF associations and bucket policies reference it |
| `domain_name` | The DNS alias target (`d123abc.cloudfront.net`) |
| `hosted_zone_id` | CloudFront's global Route53 zone (`Z2FDTNDATAQYW2`) |
| `status` | `Deployed` once propagated everywhere |
