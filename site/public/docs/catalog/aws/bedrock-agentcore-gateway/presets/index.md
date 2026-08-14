---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Gateway"
type: "preset-list"
componentSlug: "bedrock-agentcore-gateway"
componentTitle: "Bedrock AgentCore Gateway"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-openapi-tools-gateway"
    rank: "01"
    title: "OpenAPI Tools Gateway"
    excerpt: "This preset turns an existing REST API into agent tools by pointing the gateway at its OpenAPI schema in S3 — no code changes to the API, IAM (SigV4) callers only, semantic tool search on."
  - slug: "02-jwt-gateway-with-oauth-backend"
    rank: "02"
    title: "JWT Gateway with an OAuth Backend"
    excerpt: "This preset admits OIDC bearer tokens from your identity provider and fronts a partner's remote MCP server, obtaining OAuth tokens from a vaulted AgentCore Identity credential provider on every call..."
---

# Bedrock AgentCore Gateway Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Gateway. Each preset is a complete manifest you can copy, customize, and deploy.
