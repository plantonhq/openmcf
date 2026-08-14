---
title: "Edge Custom Domain"
description: "This preset fronts a REST API at a CloudFront-distributed hostname. The ACM certificate must live in `us-east-1` — that is CloudFront's region, regardless of where the API lives."
type: "preset"
rank: "02"
presetSlug: "02-edge-domain"
componentSlug: "rest-api-domain"
componentTitle: "REST API Domain"
provider: "aws"
icon: "package"
order: 2
---

# Edge Custom Domain

This preset fronts a REST API at a CloudFront-distributed hostname.
The ACM certificate must live in `us-east-1` — that is CloudFront's
region, regardless of where the API lives.

## When to Use

- Callers in many geographies who benefit from CloudFront's edge
- Hostnames that already sit behind CloudFront DNS

## What You Get

- An EDGE domain bound to a us-east-1 ACM certificate
- A root mapping onto the named REST API's `prod` stage
- `cloudfront_domain_name` / `cloudfront_zone_id` outputs for the
  Route 53 alias

## Customize

- The certificate AwsCertManagerCert must be in us-east-1
- Point DNS at the CloudFront target, not the regional one
- Prefer the regional preset unless you specifically need EDGE
