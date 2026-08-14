---
title: "Presets"
description: "Ready-to-deploy configuration presets for REST API Gateway"
type: "preset-list"
componentSlug: "rest-api-gateway"
componentTitle: "REST API Gateway"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-mock-health"
    rank: "01"
    title: "Mock Health API"
    excerpt: "This preset creates the simplest REST API: a single `GET /health` method with a MOCK integration that returns `{\"ok\":true}`. No backend to provision — the right starting point before wiring Lambda or..."
  - slug: "02-lambda-proxy"
    rank: "02"
    title: "Lambda Proxy Orders API"
    excerpt: "This preset fronts a Lambda with a REST API `POST /orders` method using AWS_PROXY integration. API Gateway forwards the full request and the function returns the HTTP response."
---

# REST API Gateway Presets

Ready-to-deploy configuration presets for REST API Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
