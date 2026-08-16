# AWS REST API Domain

An API Gateway custom domain for REST APIs — callers hit
`https://api.example.com/orders` instead of the execute-api endpoint,
TLS terminates on your certificate, and base-path mappings fan the
hostname's paths out across APIs and stages.

## What Gets Created

- A custom domain bound to an ACM certificate (REGIONAL, EDGE, or
  PRIVATE).
- Base-path mappings onto [AWS REST API Gateway](/cloud-catalog/aws-rest-api-gateway)
  stages.
- VPC-endpoint access associations when the domain is PRIVATE.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with API Gateway domain permissions.

### AWS Account

- An ACM certificate covering the hostname. REGIONAL / PRIVATE
  certificates live in the same region; EDGE certificates live in
  `us-east-1`.
- The REST APIs and stages the mappings will point at.

## Deploy

### Console

Create the resource from the AWS catalog, pick the hostname and
certificate, add base-path mappings, and deploy.

### CLI

```bash
planton apply -f rest-api-domain.yaml
```

## After Deploy

- Point a Route 53 alias at `regional_domain_name` (REGIONAL / PRIVATE)
  or `cloudfront_domain_name` (EDGE). The matching zone IDs are
  outputs.
- Changing `domain_name` replaces the domain (AWS has no rename).

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
