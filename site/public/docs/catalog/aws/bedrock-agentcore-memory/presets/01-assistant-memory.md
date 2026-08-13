---
title: "Assistant Memory"
description: "This preset gives an assistant the two memory shapes most conversations need: extracted facts per user and per-session summaries, over a 30-day raw-event window."
type: "preset"
rank: "01"
presetSlug: "01-assistant-memory"
componentSlug: "bedrock-agentcore-memory"
componentTitle: "Bedrock AgentCore Memory"
provider: "aws"
icon: "package"
order: 1
---

# Assistant Memory

This preset gives an assistant the two memory shapes most conversations
need: extracted facts per user and per-session summaries, over a 30-day
raw-event window.

## When to Use

- The first memory for a conversational assistant
- Personalization that should survive individual sessions

## What You Get

- SEMANTIC extraction distilling standalone facts into
  `/facts/{actorId}`
- Per-session summaries under `/summaries/{actorId}/{sessionId}` for
  fast context rebuilds

## Customize

- Add a USER_PREFERENCE strategy when the agent should remember stated
  preferences distinctly from facts
- Raise `eventExpiryDays` (up to 365) if you replay raw events; records
  extracted by strategies outlive the window either way
- Add `indexedKeys` up front for filtered retrieval — they are
  create-time structure and replace the memory when changed
