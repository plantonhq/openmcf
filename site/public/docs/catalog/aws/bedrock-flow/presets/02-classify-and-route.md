---
title: "Classify and Route"
description: "This preset creates a branching flow: an inline prompt classifies the request, a Condition node routes docs questions through a knowledge-base query (retrieval + generation on Nova Lite) and..."
type: "preset"
rank: "02"
presetSlug: "02-classify-and-route"
componentSlug: "bedrock-flow"
componentTitle: "Bedrock Flow"
provider: "aws"
icon: "package"
order: 2
---

# Classify and Route

This preset creates a branching flow: an inline prompt classifies the
request, a Condition node routes docs questions through a knowledge-base
query (retrieval + generation on Nova Lite) and everything else to a raw
reply — the canonical shape for intent-routed assistants.

## When to Use

- Requests that split into a retrieval-worthy class and a pass-through
  class
- Learning conditional connections: the preset pairs each conditional
  edge with the data edge that feeds the branch target

## What You Get

- A six-node graph with both connection types and a `default` else-arm
- A knowledge-base node composing an `AwsBedrockKnowledgeBase` by
  reference — the chart orders the deployments

## Customize

- Replace the KnowledgeBase branch with an Agent node
  (`agent.agentAliasArn` from an `AwsBedrockAgent`'s `alias_arns`
  output) to delegate instead of retrieve
- Add more conditions (up to 5 per node) for finer routing
