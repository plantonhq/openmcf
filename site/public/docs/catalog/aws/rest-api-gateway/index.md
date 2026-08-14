---
title: "REST API Gateway"
description: "REST API Gateway deployment documentation"
icon: "package"
order: 100
componentName: "awsrestapigateway"
---

# AWS REST API Gateway

Amazon API Gateway REST APIs (v1) — mapping templates, JSON Schema
models, request validation, per-method caching and throttling, and
EDGE / REGIONAL / PRIVATE endpoints — declared as one resource: the
API, its method tree, a single stage with an explicit deployment, and
the API-scoped satellites.

HTTP APIs (the leaner, cheaper alternative) are
[AWS HTTP API Gateway](/cloud-catalog/aws-http-api-gateway).

## What Gets Created

- A REST API from typed `routes` (the modules derive the resource tree
  from the paths) or an `openapi` document — exactly one form.
- One explicit deployment (hashed from the definition so every spec
  change redeploys) and one stage.
- Named authorizers, models, request validators, gateway responses, a
  resource policy, documentation, and an optional generated client
  certificate.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with API Gateway control-plane
  permissions (`apigateway:POST` and its siblings).

### AWS Account

- Backend targets already exist: a Lambda for `AWS` / `AWS_PROXY`
  integrations, an HTTP URL for `HTTP` / `HTTP_PROXY`, or an
  [AWS REST API VPC Link](/cloud-catalog/aws-rest-api-vpc-link) for
  private NLB backends.
- TOKEN / REQUEST authorizers need a Lambda invoke URI
  (`status.outputs.invoke_arn` on [AWS Lambda](/cloud-catalog/aws-lambda)).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, define
routes (or paste OpenAPI), and deploy.

### CLI

```bash
planton apply -f rest-api.yaml
```

## After Deploy

- `rest_api_id` / `execution_arn` / `invoke_url` identify the API;
  custom domains map onto `rest_api_id` + `stage_name`.
- Usage plans attach to this API's stage via
  [AWS REST API Usage Plan](/cloud-catalog/aws-rest-api-usage-plan).
- REST APIs deploy by snapshot, not auto-deploy — a spec change
  creates a new deployment and the stage moves to it.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
