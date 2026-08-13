---
title: "OpenAPI Tools Gateway"
description: "This preset turns an existing REST API into agent tools by pointing the gateway at its OpenAPI schema in S3 — no code changes to the API, IAM (SigV4) callers only, semantic tool search on."
type: "preset"
rank: "01"
presetSlug: "01-openapi-tools-gateway"
componentSlug: "bedrock-agentcore-gateway"
componentTitle: "Bedrock AgentCore Gateway"
provider: "aws"
icon: "package"
order: 1
---

# OpenAPI Tools Gateway

This preset turns an existing REST API into agent tools by pointing the
gateway at its OpenAPI schema in S3 — no code changes to the API, IAM
(SigV4) callers only, semantic tool search on.

## When to Use

- Exposing an existing internal API to agents without touching it
- The first gateway: one backend, IAM auth, nothing exotic

## What You Get

- One MCP URL agents connect to; AWS derives one tool per documented
  operation with your descriptions as the model-facing prose
- Outbound calls signed with the gateway's own role (SigV4)

## Customize

- Switch `authorizerType` to CUSTOM_JWT with `customJwtAuthorizer` to
  admit OIDC bearer tokens (agents outside AWS)
- Add more targets — a Lambda with explicit tool schemas, a remote MCP
  server, an AgentCore runtime — behind the same URL
- Attach a Cedar `policyEngine` in LOG_ONLY to observe authorization
  decisions before enforcing them
