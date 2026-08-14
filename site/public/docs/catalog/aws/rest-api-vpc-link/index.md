---
title: "REST API VPC Link"
description: "REST API VPC Link deployment documentation"
icon: "package"
order: 100
componentName: "awsrestapivpclink"
---

# AWS REST API VPC Link

An API Gateway v1 VPC link — the NLB attachment that lets REST API
integrations reach private services inside a VPC without exposing them
to the internet. One link is shared by many APIs.

HTTP APIs use a different link ([AWS HTTP API VPC Link](/cloud-catalog/aws-http-api-vpc-link))
that attaches to subnets directly. The two are not interchangeable.

## What Gets Created

- An API Gateway v1 VPC link fronting exactly one Network Load
  Balancer ([AWS NLB](/cloud-catalog/aws-nlb)).
- AWS-managed network attachment to that balancer (create-time
  immutable).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with API Gateway VPC-link permissions.

### AWS Account

- An internal Network Load Balancer the link will front. REST API VPC
  links accept NLB only (not ALB, not Cloud Map).

## Deploy

### Console

Create the resource from the AWS catalog, pick the NLB, and deploy.
Provisioning takes several minutes.

### CLI

```bash
planton apply -f rest-api-vpc-link.yaml
```

## After Deploy

- REST API integrations set `connectionType: VPC_LINK` and
  `connectionId` to `vpc_link_id`.
- Changing the target NLB replaces the link. Creating the link is
  free; the NLB bills as usual.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
