---
title: "Mock Health API"
description: "This preset creates the simplest REST API: a single `GET /health` method with a MOCK integration that returns `{\"ok\":true}`. No backend to provision — the right starting point before wiring Lambda or..."
type: "preset"
rank: "01"
presetSlug: "01-mock-health"
componentSlug: "rest-api-gateway"
componentTitle: "REST API Gateway"
provider: "aws"
icon: "package"
order: 1
---

# Mock Health API

This preset creates the simplest REST API: a single `GET /health` method
with a MOCK integration that returns `{"ok":true}`. No backend to
provision — the right starting point before wiring Lambda or HTTP.

## When to Use

- First REST API in an environment
- Contract stubs while backends are still being built
- Proving the resource tree, stage, and invoke URL without a Lambda bill

## What You Get

- A REGIONAL REST API (the default endpoint type)
- One derived resource (`/health`) and one GET method
- An explicit deployment and a `prod` stage

## Customize

- Add more `routes` with MOCK, HTTP, or AWS_PROXY integrations
- Set `endpointConfiguration.type` to EDGE or PRIVATE when the callers
  need it
- Attach an authorizer from the second preset once a Lambda exists
