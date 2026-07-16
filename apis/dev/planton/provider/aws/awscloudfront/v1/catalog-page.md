# AWS CloudFront

Deploys an Amazon CloudFront distribution — the global CDN front door that terminates TLS at the edge, caches responses close to viewers, and routes requests to one or more origins: S3 buckets (kept fully private via Origin Access Control), load balancers, API endpoints, or resources inside your VPC.

## What Gets Created

- **CloudFront Distribution** — origins, primary/failover origin groups, the default cache behavior plus path-matched ordered behaviors, custom domains with your certificate, custom error pages, geo restrictions, access logs, and WAF attachment.
- **Origin Access Controls** (for S3 origins that ask for one) — CloudFront signs origin requests with SigV4 so buckets stay private; one OAC per requesting origin.
- **Monitoring Subscription** (when `enableAdditionalMetrics` is set) — CloudWatch additional metrics: cache hit rate, origin latency, per-status error rates.

## Prerequisites

- **AWS credentials** configured via a Planton provider config
- For custom domains: an **ACM certificate in us-east-1** covering every alias (reference an `AwsCertManagerCert`)
- For private S3 origins: a **bucket policy** allowing the distribution's ARN (the `distribution_arn` output)

## Quick Start

Create a file `cdn.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsCloudFront
metadata:
  name: my-website-cdn
spec:
  region: us-east-1
  origins:
    - originId: s3-site
      domainName: my-site.s3.us-east-1.amazonaws.com
      s3Origin:
        createOriginAccessControl: true
  defaultCacheBehavior:
    targetOriginId: s3-site
    viewerProtocolPolicy: redirect-to-https
    compress: true
    cachePolicyId: 658327ea-f89d-4fab-a63d-7e88639e58f6  # Managed-CachingOptimized
  defaultRootObject: index.html
```

Deploy:

```shell
planton apply -f cdn.yaml
```

This serves the private bucket on the distribution's `*.cloudfront.net` domain with HTTPS redirect and the AWS-managed static-content cache policy.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `region` | `string` | Provider-connection region. CloudFront is global; its certificates and CLOUDFRONT-scope WAF ACLs live in `us-east-1`. |
| `origins` | `list` | At least one content source; each with a unique `originId` and at most one type arm (`s3Origin` / `customOrigin` / `vpcOrigin`). |
| `defaultCacheBehavior` | `object` | How unmatched requests are cached and forwarded; targets an origin or origin group by ID. |

### Cache Behaviors

| Field | Type | Description |
|-------|------|-------------|
| `defaultCacheBehavior.targetOriginId` | `string` | The `originId` or `originGroupId` this behavior routes to (validated to resolve). |
| `defaultCacheBehavior.viewerProtocolPolicy` | `string` | `redirect-to-https` (websites), `https-only` (APIs), or `allow-all`. |
| `defaultCacheBehavior.cachePolicyId` | `string` | Modern caching: a managed or custom cache policy ID. Mutually exclusive with `forwardedValues`. |
| `defaultCacheBehavior.forwardedValues` | `object` | Legacy caching: inline query/header/cookie forwarding with per-behavior TTLs. |
| `defaultCacheBehavior.functionAssociations` | `list` | CloudFront Functions on viewer events (max 2). |
| `defaultCacheBehavior.lambdaFunctionAssociations` | `list` | Lambda@Edge version ARNs on any of the four events (max 4). |
| `orderedCacheBehaviors` | `list` | Path-matched behaviors (`pathPattern` + `behavior`), evaluated in order before the default. |

### Origins

| Field | Type | Description |
|-------|------|-------------|
| `origins[].originId` | `string` | Your stable handle, unique across origins and groups. |
| `origins[].domainName` | `string` | The DNS name CloudFront connects to (regional S3 endpoint, LB DNS name, ...). |
| `origins[].s3Origin.createOriginAccessControl` | `bool` | Provision + attach an OAC (the recommended private-bucket shape). |
| `origins[].customOrigin.protocolPolicy` | `string` | `https-only` (recommended), `http-only` (S3 website endpoints), `match-viewer`. |
| `origins[].vpcOrigin.vpcOriginId` | `string` | Route to a provisioned CloudFront VPC origin. |
| `origins[].originShield` | `object` | Extra regional caching layer that collapses edge requests to the origin. |
| `originGroups[]` | `list` | Two-member primary/failover pairs with the status codes that trigger failover. |

### Custom Domains

| Field | Type | Description |
|-------|------|-------------|
| `aliases` | `list(string)` | The CNAMEs the distribution answers for; require a custom certificate covering them. |
| `viewerCertificate.acmCertificateArn` | `ref` | ACM certificate in us-east-1 (references `AwsCertManagerCert.status.outputs.cert_arn`). |
| `viewerCertificate.minimumProtocolVersion` | `string` | TLS floor; defaults to `TLSv1.2_2021`. |

### Distribution Knobs

| Field | Type | Description |
|-------|------|-------------|
| `enabled` | `bool` | Default `true`; `false` keeps the distribution deployed-but-dark. |
| `priceClass` | `string` | `PriceClass_All` (default), `PriceClass_200`, `PriceClass_100`. |
| `httpVersion` | `string` | `http2` (default), `http2and3`, `http3`, `http1.1`. |
| `isIpv6Enabled` | `bool` | Free dual-stack serving. |
| `webAclArn` | `ref` | CLOUDFRONT-scope WAF Web ACL ARN (references `AwsWafWebAcl`). |
| `customErrorResponses` | `list` | Map origin errors to custom pages with error-caching control. |
| `geoRestriction` | `object` | Country whitelist/blacklist. |
| `logging` | `object` | Standard access logs to an S3 bucket domain name. |
| `enableAdditionalMetrics` | `bool` | CloudWatch additional metrics (billed per CloudWatch rates). |
| `waitForDeployment` | `bool` | Default `true`; block until propagated to every edge (typically 5-15 min). |

## Stack Outputs

| Output | Description |
|--------|-------------|
| `distribution_id` | The distribution ID — invalidations and monitoring key on it |
| `distribution_arn` | The ARN — WAF associations and S3 bucket policies reference it |
| `domain_name` | The `*.cloudfront.net` domain — the DNS alias target |
| `hosted_zone_id` | CloudFront's global Route53 zone ID for alias records |
| `status` | `Deployed` once propagated to every edge location |

## Composition

- Point DNS with an `AwsRoute53DnsRecord` alias record built from `domain_name` + `hosted_zone_id`.
- Front a private `AwsS3Bucket` (regional endpoint + OAC + bucket policy on `distribution_arn`).
- Attach an `AwsCertManagerCert` (us-east-1) for custom domains and an `AwsWafWebAcl` (CLOUDFRONT scope) for edge protection.
- Route `/api/*` to an `AwsAlb` through an ordered behavior with `Managed-CachingDisabled` while static paths serve from S3.
