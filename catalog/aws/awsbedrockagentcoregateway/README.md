<p align="center">
  <img src="logo.svg" alt="AWS Bedrock AgentCore Gateway" width="80"/>
</p>

# AWS Bedrock AgentCore Gateway

Create and manage [Amazon Bedrock AgentCore gateways](https://docs.aws.amazon.com/bedrock-agentcore/latest/devguide/gateway.html)
— managed MCP (Model Context Protocol) front doors that turn your
existing APIs, Lambda functions, and MCP servers into tools any
MCP-speaking agent can discover and call through ONE authenticated URL.

## What Gets Created

- **A gateway** with IAM (SigV4) or custom JWT (OIDC) inbound
  authorization and an MCP endpoint URL.
- **Targets** (folded satellites) — each exposes one backend as MCP
  tools:
  - an **AgentCore agent runtime** (plain HTTP),
  - an **API Gateway REST stage** (one tool per route, filterable),
  - a **Lambda function** with explicitly-defined tools and JSON
    schemas,
  - an existing **remote MCP server**,
  - an **OpenAPI 3 schema** or **Smithy model** (inline or in S3).
- Per-target **outbound credentials**: API keys and OAuth tokens from
  AgentCore Identity providers, the caller's own IAM credentials, the
  gateway's role (SigV4), or JWT passthrough.
- Optional: Lambda **interceptors** on the request/response path, a
  Cedar **policy engine** evaluating every tool call, semantic tool
  search, and session/streaming tuning.

Creating a gateway is free — AWS bills per tool call at runtime.

## Destroy Ordering

AWS deletes a gateway's targets automatically before the gateway itself
— the provider manages the drain; zero orphans by construction.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
