---
title: "S3 Static Website with Origin Access Control"
description: "This preset serves a private S3 bucket through CloudFront using an Origin Access Control (OAC) -- the modern way to front S3. The bucket never becomes public: CloudFront signs its origin requests..."
type: "preset"
rank: "01"
presetSlug: "01-s3-static-website"
componentSlug: "cloudfront"
componentTitle: "CloudFront"
provider: "aws"
icon: "package"
order: 1
---

# S3 Static Website with Origin Access Control

This preset serves a private S3 bucket through CloudFront using an Origin Access Control (OAC) -- the modern way to front S3. The bucket never becomes public: CloudFront signs its origin requests with SigV4, and the bucket policy allows only this distribution's ARN.

## When to Use

- Static websites and single-page applications built to S3
- Asset/download distribution from a private bucket
- Any S3-backed content that should only be reachable through the CDN

## Key Configuration Choices

- **Origin Access Control** (`s3Origin.createOriginAccessControl: true`) -- the module provisions the OAC and attaches it to the origin; add a bucket policy allowing `cloudfront.amazonaws.com` with an `AWS:SourceArn` condition on the distribution's ARN (the `distribution_arn` output)
- **Managed cache policy** (`cachePolicyId: 658327ea-...`) -- `Managed-CachingOptimized`, the AWS-managed policy tuned for static content; no custom policy to maintain
- **SPA-friendly error mapping** -- S3 returns 403 for missing objects when accessed via OAC; mapping it to `200 /index.html` lets client-side routers own unknown paths (use `404` + an error page for a classic multi-page site)
- **Default CloudFront certificate** -- the distribution serves on `*.cloudfront.net`; add `aliases` + `viewerCertificate` for a custom domain (see preset 02)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<bucket-name>` | The S3 bucket holding the site content | `AwsS3Bucket` outputs or the S3 console |
| `<bucket-region>` | The bucket's region (the REGIONAL endpoint avoids redirect latency) | The bucket's configuration |

## Related Presets

- **02-custom-domain-cdn** -- Use when the site should serve on your own domain with an ACM certificate
