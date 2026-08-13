---
title: "Metered API Keys"
description: "This preset covers one REST API stage with a 1000-request daily quota and one enabled API key — the starting shape for a metered consumer."
type: "preset"
rank: "01"
presetSlug: "01-metered-api-keys"
componentSlug: "rest-api-usage-plan"
componentTitle: "REST API Usage Plan"
provider: "aws"
icon: "package"
order: 1
---

# Metered API Keys

This preset covers one REST API stage with a 1000-request daily quota
and one enabled API key — the starting shape for a metered consumer.

## When to Use

- The first usage plan in an environment
- Partner or mobile-app callers that should be counted, not just
  authenticated

## What You Get

- A plan covering the named AwsRestApiGateway's `prod` stage
- A daily quota of 1000 requests
- One enabled API key (`default-consumer`)

## Customize

- Point `restApiId` at your REST API
- Add more `apiKeys` entries (one per consumer)
- Methods that should require the key need `apiKeyRequired: true` on
  the gateway
