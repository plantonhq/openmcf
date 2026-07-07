# AwsCloudFront

The **AwsCloudFront** resource provisions an Amazon CloudFront distribution — the global CDN front door that terminates TLS at the edge, caches responses close to viewers, and routes requests to one or more origins (S3 buckets, load balancers, API endpoints, or anything HTTP-addressable).

## The Model

The spec mirrors CloudFront's own composition:

- **`origins`** declare WHERE content comes from. Each origin carries a user-chosen `originId` (the stable handle behaviors target) and one origin-type arm: `s3Origin` (REST endpoint, pair with origin access control), `customOrigin` (anything HTTP — load balancers, APIs, S3 *website* endpoints), or `vpcOrigin` (a provisioned CloudFront VPC origin reaching private resources with zero public exposure).
- **`originGroups`** compose two origins into a primary/failover pair; behaviors can target a group's ID exactly like an origin's.
- **`defaultCacheBehavior`** (plus path-matched **`orderedCacheBehaviors`**, first match wins) declares HOW requests are matched, cached, and forwarded — protocol policy, methods, compression, edge functions, and the caching configuration.
- **`viewerCertificate`** + **`aliases`** put the distribution on your own domain.

Validation proves every behavior target and origin-group member resolves to a declared origin, so a dangling reference is caught at manifest time instead of at deploy time.

## Private S3 Origins: Origin Access Control

The modern way to front S3 is an **Origin Access Control** — CloudFront signs origin requests with SigV4 so the bucket stays fully private. Set `s3Origin.createOriginAccessControl: true` and the module provisions and attaches the OAC; then allow the distribution's ARN in the bucket policy (`cloudfront.amazonaws.com` principal with an `AWS:SourceArn` condition on the `distribution_arn` output). An existing OAC can be attached by ID, and an existing legacy Origin Access Identity path is accepted (never created — OAC supersedes it).

## Caching: Two Generations

Each behavior chooses exactly one caching generation:

- **Modern (recommended)** — `cachePolicyId` (+ `originRequestPolicyId`, `responseHeadersPolicyId`). AWS ships managed policies covering most cases: `Managed-CachingOptimized` (`658327ea-f89d-4fab-a63d-7e88639e58f6`) for static content, `Managed-CachingDisabled` (`4135ea2d-6df8-44a3-9df3-4b5a84be39ad`) for APIs, and origin-request policy `Managed-AllViewer` (`216adef6-5c7f-47e4-b989-5492eafa07d3`) to forward everything.
- **Legacy** — the inline `forwardedValues` block with per-behavior TTLs. Kept for existing configurations; mutually exclusive with a cache policy.

## Edge Compute

- **CloudFront Functions** (`functionAssociations`) — sub-millisecond JavaScript at every edge, viewer-request/viewer-response only. The lightweight choice for URL rewrites, redirects, and header manipulation.
- **Lambda@Edge** (`lambdaFunctionAssociations`) — full Lambda at regional edges, all four event types. Requires a numbered function VERSION ARN in us-east-1.

## Custom Domains

Set `aliases` plus `viewerCertificate` with the ACM arm (a `StringValueOrRef` to an `AwsCertManagerCert` — the certificate **must live in us-east-1** and cover every alias) or the legacy IAM arm. SNI-only serving and the `TLSv1.2_2021` protocol floor are the defaults. Point DNS at the distribution with Route53 alias records built from the `domain_name` and `hosted_zone_id` outputs.

## Everything Else

- **`customErrorResponses`** — replace origin errors with custom pages (e.g. map S3's 403-for-missing-object to a 404, or to `200 /index.html` for SPAs) and control error caching.
- **`geoRestriction`** — allow or deny viewers by country.
- **`logging`** — standard access logs to an S3 bucket (the bucket needs ACLs enabled).
- **`webAclArn`** — a CLOUDFRONT-scope WAF Web ACL by ARN (referenceable from an `AwsWafWebAcl`).
- **`priceClass`** — the cost/latency dial (`PriceClass_All` default, `PriceClass_200`, `PriceClass_100`).
- **`httpVersion`** / **`isIpv6Enabled`** — protocol surface (`http2and3` is the safe way to adopt HTTP/3; IPv6 costs nothing).
- **`enableAdditionalMetrics`** — CloudWatch additional metrics (cache hit rate, origin latency, per-status error rates).
- **`enabled`** / **`waitForDeployment`** / **`retainOnDelete`** — operational knobs; deploys propagate to every edge location (typically 5-15 minutes).

## Stack Outputs

| Output | Description |
|--------|-------------|
| `distribution_id` | What invalidation requests and monitoring key on |
| `distribution_arn` | The WAF-association and bucket-policy join key |
| `domain_name` | The `d123abc.cloudfront.net` target for DNS |
| `hosted_zone_id` | CloudFront's global Route53 zone (`Z2FDTNDATAQYW2`), exported so alias records compose without hardcoding |
| `status` | `Deployed` or `InProgress` |

## Deliberately Not Modeled (candidate kinds on demand)

Cache/origin-request/response-headers policies, CloudFront Functions, key groups and public keys (signed URLs), VPC origins, real-time log configurations, field-level encryption profiles, continuous-deployment policies, the multi-tenant distribution family, anycast IP lists, and mTLS trust stores each have independent lifecycles and are candidate first-class kinds. The spec carries their IDs/ARNs today, so future kinds compose with zero rework.
