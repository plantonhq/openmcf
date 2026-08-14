---
title: "Single Text Prompt"
description: "This preset creates the simplest managed prompt: one text variant on Amazon Nova Micro, deterministic (temperature 0), with a single `{{input}}` variable — a prompt your applications invoke by ID..."
type: "preset"
rank: "01"
presetSlug: "01-single-text-prompt"
componentSlug: "bedrock-prompt"
componentTitle: "Bedrock Prompt"
provider: "aws"
icon: "package"
order: 1
---

# Single Text Prompt

This preset creates the simplest managed prompt: one text variant on
Amazon Nova Micro, deterministic (temperature 0), with a single
`{{input}}` variable — a prompt your applications invoke by ID instead
of embedding the string in code.

## When to Use

- Moving the first prompt out of application code and into managed,
  reviewable infrastructure
- Single-step tasks: summarize, classify, extract, rewrite

## What You Get

- A prompt whose DRAFT you edit declaratively; consumers invoking by
  prompt ID pick up changes on the draft, or pin published versions for
  stability

## Customize

- Add variables to the template and `inputVariables` together — AWS
  matches them at invocation
- Raise `maxTokens`/`temperature` for longer or more creative outputs
