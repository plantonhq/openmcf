---
title: "LLM Judge and Harness"
description: "This preset pairs an LLM-as-a-judge evaluator with a Bedrock harness — the shape for repeatable agent scoring runs. Creating the objects is free; AWS bills when you start a run."
type: "preset"
rank: "02"
presetSlug: "02-llm-judge-harness"
componentSlug: "bedrock-agentcore-evaluation"
componentTitle: "Bedrock AgentCore Evaluation"
provider: "aws"
icon: "package"
order: 2
---

# LLM Judge and Harness

This preset pairs an LLM-as-a-judge evaluator with a Bedrock harness
— the shape for repeatable agent scoring runs. Creating the objects is
free; AWS bills when you start a run.

## When to Use

- You have Bedrock model access in the region
- You want a named test bench (model + system prompt) evaluation
  runs execute against

## What You Get

- A SESSION-level categorical judge (`pass` / `fail`)
- A Bedrock harness on the named AgentCore IAM role
- IDs in `evaluator_ids` / `harness_ids` for later runs

## Customize

- Swap the model ID for the foundation model your account has enabled
- Point `executionRoleArn` at a role trusting
  `bedrock-agentcore.amazonaws.com` with `bedrock:InvokeModel`
- Add tools or an online evaluation config when you are ready to
  sample production sessions
