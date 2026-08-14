---
title: "Bedrock AgentCore Memory"
description: "Bedrock AgentCore Memory deployment documentation"
icon: "package"
order: 100
componentName: "awsbedrockagentcorememory"
---

# AWS Bedrock AgentCore Memory

Agent memory as managed infrastructure — a store that keeps raw session
events for a retention window you set and distills them into long-term
records (facts, summaries, preferences, episodes) through declarative
extraction strategies, encrypted and queryable by namespace.

## What Gets Created

- A memory with a 7–365 day short-term window for raw session events.
- Extraction strategies: built-in SEMANTIC / SUMMARIZATION /
  USER_PREFERENCE / EPISODIC pipelines, or CUSTOM prompt/model
  overrides — each partitioning its records by namespace templates.
- Optional KMS encryption, indexed metadata keys for filtered
  retrieval, and Kinesis streaming of records as they are written.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with AgentCore control-plane permissions
  (`bedrock-agentcore:CreateMemory` and its siblings).

### AWS Account

- Bedrock AgentCore available in the target region.
- For CUSTOM strategies or Kinesis delivery: an IAM role trusting
  `bedrock-agentcore.amazonaws.com` with model-invoke / stream-write
  access.

## Deploy

### Console

Create the resource from the AWS catalog, pick the region and retention
window, add strategies, and deploy.

### CLI

```bash
planton apply -f memory.yaml
```

## After Deploy

- `memory_id` / `memory_arn` identify the memory; `strategy_ids`
  carries each strategy.
- Agent code writes events and queries records through the AgentCore
  data plane; evaluation harnesses reference the memory by ARN.
- Strategy changes serialize through the parent memory and can take
  tens of minutes — batch them.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
