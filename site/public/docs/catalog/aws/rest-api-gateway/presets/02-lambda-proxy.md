---
title: "Lambda Proxy Orders API"
description: "This preset fronts a Lambda with a REST API `POST /orders` method using AWS_PROXY integration. API Gateway forwards the full request and the function returns the HTTP response."
type: "preset"
rank: "02"
presetSlug: "02-lambda-proxy"
componentSlug: "rest-api-gateway"
componentTitle: "REST API Gateway"
provider: "aws"
icon: "package"
order: 2
---

# Lambda Proxy Orders API

This preset fronts a Lambda with a REST API `POST /orders` method using
AWS_PROXY integration. API Gateway forwards the full request and the
function returns the HTTP response.

## When to Use

- Serverless APIs where Lambda owns request handling
- Routes that should require an API key (pair with
  AwsRestApiUsagePlan)

## What You Get

- One POST method integrated with the named AwsLambda
- `apiKeyRequired: true` so a usage plan can meter callers
- An explicit deployment and a `prod` stage

## Customize

- Point `uri.valueFrom.name` at your AwsLambda resource
- Add GET / PUT / DELETE routes on the same path
- Drop `apiKeyRequired` for a public method; add a TOKEN authorizer
  when callers present a bearer token instead of a key
