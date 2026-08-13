---
title: "Bedrock AgentCore Gateway"
description: "Bedrock AgentCore Gateway deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockagentcoregateway"
---

# AWS Bedrock AgentCore Gateway

Your existing APIs as agent tools — a managed MCP front door that
converts REST APIs, Lambda functions, OpenAPI/Smithy schemas, and remote
MCP servers into tools agents discover and call through one
authenticated URL, with per-target credentials and optional Cedar
authorization on every call.

## What Gets Created

- An MCP gateway with IAM (SigV4) or custom JWT (OIDC) inbound auth.
- One target per backend: agent runtimes, API Gateway stages, Lambda
  tools with JSON schemas, remote MCP servers, OpenAPI/Smithy schemas.
- Per-target outbound credentials (vaulted API keys / OAuth via
  AgentCore Identity, caller or gateway IAM, JWT passthrough), Lambda
  interceptors, and a Cedar policy engine in LOG_ONLY or ENFORCE mode.

Creating a gateway is free; AWS bills per tool call at runtime.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore control-plane permissions
  (`bedrock-agentcore:CreateGateway` and its siblings).

### AWS Account

- Bedrock AgentCore available in the target region.
- An IAM role trusting `bedrock-agentcore.amazonaws.com` that can reach
  your targets (invoke Lambdas, call API Gateway, sign SigV4).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and inbound
auth, add targets, and deploy.

### CLI

```bash
planton apply -f gateway.yaml
```

## After Deploy

- `gateway_url` is the MCP URL agents connect to; `target_ids` carries
  each target.
- Point any MCP client (an AgentCore runtime agent, Claude, an IDE) at
  the URL with credentials matching your inbound auth.
- Roll out authorization safely: attach the policy engine in LOG_ONLY,
  read the decisions, then flip to ENFORCE.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
