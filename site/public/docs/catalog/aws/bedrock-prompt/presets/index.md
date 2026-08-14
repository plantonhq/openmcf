---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Prompt"
type: "preset-list"
componentSlug: "bedrock-prompt"
componentTitle: "Bedrock Prompt"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-single-text-prompt"
    rank: "01"
    title: "Single Text Prompt"
    excerpt: "This preset creates the simplest managed prompt: one text variant on Amazon Nova Micro, deterministic (temperature 0), with a single `{{input}}` variable — a prompt your applications invoke by ID..."
  - slug: "02-chat-with-tools"
    rank: "02"
    title: "Chat with Tools"
    excerpt: "This preset creates a two-variant prompt for A/B comparison: a plain text variant and a chat variant with system context and an order-lookup tool the model may call (auto tool choice) — with the chat..."
---

# Bedrock Prompt Presets

Ready-to-deploy configuration presets for Bedrock Prompt. Each preset is a complete manifest you can copy, customize, and deploy.
