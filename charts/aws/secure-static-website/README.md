# AWS Secure Static Website

A production static site the way AWS security guidance says to build one: a
fully private S3 bucket served exclusively through CloudFront with Origin
Access Control, TLS on your own domain via ACM, Route 53 alias records, and
an optional AWS-managed-rules WAF at the edge. Nothing in this architecture
is publicly readable except the CDN itself.

This is the chart for marketing sites, documentation, single-page apps, and
any "build once, serve from the edge" frontend — the highest-traffic thing
most teams run, and the one they most often get subtly wrong (public buckets,
legacy origin identities, HTTP-only origins). Here the secure shape is the
only shape.

## Architecture

```
                          ┌───────────────────┐  alias A + AAAA   ┌──────────────────┐
        visitors ────────▶│  AwsRoute53Zone   │──────────────────▶│                  │
   https://www.example.com│  (or existing)    │                   │                  │
                          └───────────────────┘                   │   AwsCloudFront  │
                          ┌───────────────────┐   TLS (us-east-1) │   OAC origin     │
                          │ AwsCertManagerCert│──────────────────▶│   price class    │
                          │  DNS-validated    │                   │   SPA fallback   │
                          └───────────────────┘                   └────────┬─────────┘
                          ┌───────────────────┐   inspect first           │ SigV4-signed
                          │   AwsWafWebAcl    │◀──────────────────────────┤ GetObject
                          │  managed packs    │                           ▼
                          └───────────────────┘                  ┌──────────────────┐
                                                                 │   AwsS3Bucket    │
                                                                 │  private, OAC-   │
                                                                 │  only policy     │
                                                                 └──────────────────┘
```

Deployment order derives from the references: zone → certificate (validates
through the zone) and WAF in parallel → distribution (consumes both) → alias
records (read the distribution's address).

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Content bucket | `AwsS3Bucket` | Private origin for built site files — CloudFront-only bucket policy, all public access blocked |
| CDN | `AwsCloudFront` | Edge delivery, TLS termination, OAC-signed origin fetches, SPA fallback (conditional) |
| Certificate | `AwsCertManagerCert` | DNS-validated TLS for the custom domain, pinned to us-east-1 (conditional) |
| Hosted zone | `AwsRoute53Zone` | Public DNS zone for the domain (conditional — bring your own instead) |
| Alias records | `AwsRoute53DnsRecord` ×2 | A + AAAA aliases pointing the domain at the distribution (conditional) |
| Web ACL | `AwsWafWebAcl` | CLOUDFRONT-scope managed-rules firewall (conditional) |

## Parameters

| Parameter | Description | Default | Type |
|-----------|-------------|---------|------|
| `site_name` | Name prefix for the CDN-side resources | `my-site` | string |
| `bucket_name` | Globally unique content bucket name | `my-org-site` | string |
| `aws_account_id` | Account id scoping the bucket's CloudFront trust condition | `123456789012` | string |
| `aws_region` | Content bucket region (cert/WAF are us-east-1 by AWS contract) | `us-east-1` | string |
| `domain_name` | The FQDN the site serves on | `www.example.com` | string |
| `dns_zone_name` | Hosted-zone domain (when the chart creates the zone) | `example.com` | string |
| `existing_zone_id` | Existing hosted-zone id (when `dns_zone_enabled` is off) | `Z0123456789EXAMPLE` | string |
| `price_class` | CloudFront edge coverage: `PriceClass_100` / `200` / `All` | `PriceClass_100` | string |
| `custom_domain_enabled` | Own domain + certificate + DNS records | `true` | bool |
| `dns_zone_enabled` | Chart creates the hosted zone vs bring-your-own | `true` | bool |
| `waf_enabled` | Managed-rules WAF in front of the CDN (~$6+/mo) | `true` | bool |
| `spa_mode` | Route unknown paths to /index.html for client-side routers | `true` | bool |

## After deploying

1. **Fresh zone only** (`dns_zone_enabled: true`): delegate the domain at your
   registrar to the zone's `status.outputs.nameservers`. Until delegation is
   live, certificate DNS-validation cannot complete — on a brand-new domain
   do this immediately after the zone appears (the certificate deployment
   waits for validation).
2. **Upload the site**: sync your build output to the bucket, e.g.
   `aws s3 sync ./dist s3://<bucket_name>`. The `AwsS3ObjectSet` kind can
   also manage the files as infrastructure if you prefer declarative content.
3. **Invalidate on deploy**: static assets are cached at the edge by design.
   Have CI create an invalidation for `/*` (or use hashed filenames and only
   invalidate `/index.html`) after each upload:
   `aws cloudfront create-invalidation --distribution-id <id> --paths "/*"`.
4. Browse to `https://<domain_name>` — or, with the custom domain off, to the
   distribution's `status.outputs.domain_name`.

## Day-2 guidance

- **Tighten the bucket policy to the exact distribution** once it exists:
  replace the account-scoped `ArnLike` condition value
  (`arn:aws:cloudfront::<account>:distribution/*`) with the distribution's
  `status.outputs.distribution_arn`. The account scoping already blocks
  cross-account abuse; pinning the ARN also isolates the bucket from other
  distributions in your own account.
- **Adding www/apex twins**: request the certificate with the extra name in
  `alternateDomainNames`, add it to the distribution's `aliases`, and create
  a second pair of alias records. For a proper 301 (rather than serving
  duplicate content), attach a small CloudFront Function to the twin.
- **Watching the WAF**: both packs emit CloudWatch metrics per rule. To trial
  a new pack without risk, add it with `overrideAction: count` first and
  promote to `none` when the sampled requests look right.
- **Custom 404 page** (with `spa_mode: false`): add a
  `customErrorResponses` entry mapping 403 to your `/404.html` with
  `responseCode: 404` — remember S3-via-OAC signals "not found" as 403.
- **Tearing down**: the bucket refuses deletion while it holds site files
  unless `forceDestroy` is set — intentional friction; empty it (or set the
  flag deliberately) before destroying the stack.

---

© Planton. Licensed under [Apache-2.0](../../../LICENSE).
