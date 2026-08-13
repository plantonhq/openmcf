---
title: "Bedrock AgentCore Tools"
description: "Bedrock AgentCore Tools deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockagentcoretools"
---

# AWS Bedrock AgentCore Tools

Agent hands as managed infrastructure — cloud browsers with recording
and enterprise policy control, saved browser profiles, and sandboxed
code interpreters, so agents browse and compute inside AWS-managed
isolation instead of on your hosts.

## What Gets Created

- Cloud browsers: PUBLIC or VPC egress, optional S3 session recording,
  traffic signing, Chrome enterprise policies, mTLS certificates.
- Browser profiles: saved cookies/logins sessions start from.
- Code interpreters: SANDBOX (no network), PUBLIC, or VPC sandboxes for
  model-written code.

Tools are free to create; AWS bills per session at runtime.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore tool permissions
  (`bedrock-agentcore:CreateBrowser`, `CreateBrowserProfile`,
  `CreateCodeInterpreter` and their siblings).

### AWS Account

- Bedrock AgentCore available in the target region.
- For recordings/policies/certificates: an IAM role trusting
  `bedrock-agentcore.amazonaws.com` with the matching S3 / Secrets
  Manager access.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region, author the
tools this bundle owns, and deploy.

### CLI

```bash
planton apply -f tools.yaml
```

## After Deploy

- The output maps carry each tool's IDs and ARNs — agent code starts
  sessions against them, and evaluation harnesses reference them as
  browser/code-interpreter tools.
- Remember every field change recreates the tool (AWS exposes no
  update) — in-flight sessions on the old tool finish; new sessions land
  on the replacement.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
