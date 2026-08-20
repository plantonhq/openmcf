# AWS API Gateway Account Settings

The region's API Gateway account object — a settings singleton whose
one lever is the CloudWatch Logs role REST API stages log through.
Stage-level logging on
[AWS REST API Gateway](/cloud-catalog/aws-rest-api-gateway) silently
does nothing until this role is set.

## What Gets Managed

- The region-wide CloudWatch logging role for API Gateway (one
  account object per region; deploy at most one instance per region).

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with API Gateway account permissions.

### AWS Account

- An IAM role trusting `apigateway.amazonaws.com` with CloudWatch
  Logs write permissions ([AWS IAM Role](/cloud-catalog/aws-iam-role)
  with the managed `AmazonAPIGatewayPushToCloudWatchLogs` policy).

## Deploy

### Console

Create the resource from the AWS catalog, pick the logging role, and
deploy.

### CLI

```bash
planton apply -f apigw-account-settings.yaml
```

## After Deploy

- Enable execution/access logging on REST API stages — it works now.
- Destroying this component resets the role (region-wide API Gateway
  logging stops). The account object itself always exists.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
