---
title: "CloudWatch Logging Role"
description: "This preset sets the region's API Gateway logging role — the prerequisite for execution/access logging on every REST API stage in the region."
type: "preset"
rank: "01"
presetSlug: "01-cloudwatch-logging-role"
componentSlug: "api-gateway-account-settings"
componentTitle: "API Gateway Account Settings"
provider: "aws"
icon: "package"
order: 1
---

# CloudWatch Logging Role

This preset sets the region's API Gateway logging role — the
prerequisite for execution/access logging on every REST API stage in
the region.

## When to Use

- Before enabling stage logging on any REST API in the region
- Once per region (this is a settings singleton)

## What You Get

- The account-level CloudWatch role set region-wide
- Stage logging on AwsRestApiGateway starts working

## Customize

- Point `cloudwatchRoleArn` at a role trusting
  `apigateway.amazonaws.com` with the managed
  `AmazonAPIGatewayPushToCloudWatchLogs` policy

## Composing

```yaml
# The role this preset expects:
apiVersion: aws.planton.dev/v1alpha1
kind: AwsIamRole
metadata:
  name: apigw-logging-role
spec:
  region: <aws-region>
  managedPolicyArns:
    - value: arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs
  trustPolicy:
    Version: "2012-10-17"
    Statement:
      - Effect: Allow
        Principal:
          Service: apigateway.amazonaws.com
        Action: sts:AssumeRole
```
