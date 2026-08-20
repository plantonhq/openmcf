# AWS REST API Usage Plan

API Gateway usage plans meter and throttle REST API consumers: each
plan covers one or more stages, sets quota and throttle ceilings, and
admits callers through API keys attached to the plan.

## What Gets Created

- A usage plan covering the named
  [AWS REST API Gateway](/cloud-catalog/aws-rest-api-gateway) stages.
- Optional quota (requests per day / week / month) and throttle
  ceilings, including per-method throttles.
- API keys created and attached to the plan. Key values are secrets
  and are not exported.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with API Gateway usage-plan permissions.

### AWS Account

- The REST APIs and stages the plan will cover. Methods that should
  require a key need `apiKeyRequired: true` on
  [AWS REST API Gateway](/cloud-catalog/aws-rest-api-gateway).

## Deploy

### Console

Create the resource from the AWS catalog, pick the stages, set quota
and throttle, add keys, and deploy.

### CLI

```bash
planton apply -f usage-plan.yaml
```

## After Deploy

- Distribute key *values* from the AWS console (they are not stack
  outputs). `api_key_ids` identify the keys for later rotation.
- A key admitted by this plan may call exactly the stages listed in
  `apiStages`.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
