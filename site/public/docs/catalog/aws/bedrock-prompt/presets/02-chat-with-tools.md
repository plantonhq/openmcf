---
title: "Chat with Tools"
description: "This preset creates a two-variant prompt for A/B comparison: a plain text variant and a chat variant with system context and an order-lookup tool the model may call (auto tool choice) — with the chat..."
type: "preset"
rank: "02"
presetSlug: "02-chat-with-tools"
componentSlug: "bedrock-prompt"
componentTitle: "Bedrock Prompt"
provider: "aws"
icon: "package"
order: 2
---

# Chat with Tools

This preset creates a two-variant prompt for A/B comparison: a plain
text variant and a chat variant with system context and an order-lookup
tool the model may call (auto tool choice) — with the chat variant as
the default.

## When to Use

- Comparing a simple formulation against a tool-assisted one on live
  traffic by flipping `defaultVariant`
- Chat surfaces where the model should decide when to call your tools

## What You Get

- A chat template with system instructions and a JSON-Schema-described
  tool — the description is what the model reads to decide tool use
- Variant metadata (`team: support`) for ownership annotations

## Customize

- Force structured output with `toolChoice.any: true` (the model MUST
  call some tool) or `toolChoice.toolName` (a specific one)
- Target an agent instead of a model: replace `modelId` with
  `agentAliasArn` referencing an `AwsBedrockAgent`'s `alias_arns` output
