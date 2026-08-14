---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Memory"
type: "preset-list"
componentSlug: "bedrock-agentcore-memory"
componentTitle: "Bedrock AgentCore Memory"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-assistant-memory"
    rank: "01"
    title: "Assistant Memory"
    excerpt: "This preset gives an assistant the two memory shapes most conversations need: extracted facts per user and per-session summaries, over a 30-day raw-event window."
  - slug: "02-episodic-memory-with-streaming"
    rank: "02"
    title: "Episodic Memory with Streaming"
    excerpt: "This preset captures experience episodes with EPISODIC reflection — \"what happened and what worked\" — indexed by customer and streamed to Kinesis as records are written, for agents that should learn..."
---

# Bedrock AgentCore Memory Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Memory. Each preset is a complete manifest you can copy, customize, and deploy.
