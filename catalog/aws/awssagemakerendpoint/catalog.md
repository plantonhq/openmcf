# AWS SageMaker Endpoint

Real-time inference as one declarative resource — the endpoint and its
capacity definition folded together, serving your SageMaker models
through serverless or instance-backed variants, with weighted traffic
splits, shadow testing, data capture, and guarded blue/green or rolling
deployments.

## What Gets Created

- An endpoint plus its endpoint configuration — the modules roll
  name-suffixed configurations on every capacity change and repoint
  the endpoint, so updates never break the running fleet.
- 1–10 production variants, each serverless (no idle charge,
  per-inference billing) or instance-backed with managed instance scaling, routing
  strategy, and startup/download timeouts; optional shadow variants
  for mirror-traffic testing.
- Optional: data capture to S3, async inference with S3 delivery and
  SNS notifications, KMS-encrypted volumes, and blue/green or rolling
  deployment policies with CloudWatch-alarm auto-rollback.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with SageMaker control-plane permissions
  (`sagemaker:CreateEndpoint`, `sagemaker:CreateEndpointConfig`, and
  their siblings).

### AWS Account

- A SageMaker model for each variant to serve (or an
  `execution_role_arn` on the endpoint for inference-component
  endpoints, where components attach models later).
- For data capture, async inference, or core dumps: the S3 buckets
  (and SNS topics) the configuration points at.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, add a
variant pointing at your model, choose serverless or instance capacity,
and deploy.

### CLI

```bash
planton apply -f endpoint.yaml
```

## After Deploy

- `endpoint_name` / `endpoint_arn` identify the endpoint — clients
  invoke it by name; `endpoint_config_name` /
  `endpoint_config_arn` echo the configuration currently in service.
- Capacity changes mint a new configuration and repoint the endpoint
  in place — shape the crossing with a `deployment` policy and
  rollback alarms for production fleets.
- Serverless variants bill per inference and cost nothing idle;
  instance variants bill per instance-hour from the moment the
  endpoint is InService.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
