---
title: "JWT Gateway with an OAuth Backend"
description: "This preset admits OIDC bearer tokens from your identity provider and fronts a partner's remote MCP server, obtaining OAuth tokens from a vaulted AgentCore Identity credential provider on every call..."
type: "preset"
rank: "02"
presetSlug: "02-jwt-gateway-with-oauth-backend"
componentSlug: "bedrock-agentcore-gateway"
componentTitle: "Bedrock AgentCore Gateway"
provider: "aws"
icon: "package"
order: 2
---

# JWT Gateway with an OAuth Backend

This preset admits OIDC bearer tokens from your identity provider and
fronts a partner's remote MCP server, obtaining OAuth tokens from a
vaulted AgentCore Identity credential provider on every call — no
credential material in the manifest or the gateway.

## When to Use

- Agents running outside AWS (no SigV4) calling through the gateway
- Backends that demand OAuth client credentials you refuse to scatter

## What You Get

- Inbound: only tokens matching your issuer and audience pass
- Outbound: AWS fetches and refreshes the partner token from the vault;
  rotating the client secret touches only the Identity component
- DYNAMIC listing keeps the tool set live as the partner adds tools

## Customize

- Add `customClaims` rules (EQUALS/CONTAINS/CONTAINS_ANY) for
  finer-grained token gating
- Use `jwtPassthrough: true` credentials instead when the backend wants
  the caller's own token forwarded
