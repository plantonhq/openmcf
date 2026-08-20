# AWS Bedrock AgentCore Runtime

Your agent code as managed infrastructure — a serverless,
session-isolated host for any agent framework, with versioned rollouts
through named endpoints, AWS-managed scaling, and per-second billing
only while sessions execute.

## What Gets Created

- An AgentCore runtime executing your artifact: a container image (ECR)
  or an AWS-managed code bundle (Python/Node source in S3 — no image to
  build or maintain).
- Named serving endpoints that float on the latest version or pin an
  explicit one — promote a new version by re-pointing an endpoint, never
  by redeploying consumers.
- Optional VPC networking, inbound JWT (OIDC) authorization with custom
  claim checks, filesystem mounts (EFS / S3 access points / per-session
  scratch), and a resource policy on the runtime's ARN.

Creating a runtime is free; AWS bills CPU/memory per second only while
sessions execute.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore control-plane permissions
  (`bedrock-agentcore:CreateAgentRuntime` and its siblings).

### AWS Account

- Bedrock AgentCore available in the target region.
- An IAM role trusting `bedrock-agentcore.amazonaws.com` with pull/read
  access to your artifact (ECR image or S3 code bundle).

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, point at your
artifact, and deploy.

### CLI

```bash
planton apply -f runtime.yaml
```

## After Deploy

- `agent_runtime_id` / `agent_runtime_arn` identify the runtime;
  `endpoint_arns` carries each named endpoint.
- Invoke through the AgentCore data plane (`InvokeAgentRuntime`) using
  an endpoint qualifier.
- Every spec change advances `agent_runtime_version`; floating endpoints
  pick it up automatically.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
