---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Evaluation"
type: "preset-list"
componentSlug: "bedrock-agentcore-evaluation"
componentTitle: "Bedrock AgentCore Evaluation"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-code-evaluator"
    rank: "01"
    title: "Code Evaluator"
    excerpt: "This preset creates a single code-based evaluator backed by a Lambda — the cheapest first evaluation object. Create does not invoke the function."
  - slug: "02-llm-judge-harness"
    rank: "02"
    title: "LLM Judge and Harness"
    excerpt: "This preset pairs an LLM-as-a-judge evaluator with a Bedrock harness — the shape for repeatable agent scoring runs. Creating the objects is free; AWS bills when you start a run."
---

# Bedrock AgentCore Evaluation Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Evaluation. Each preset is a complete manifest you can copy, customize, and deploy.
