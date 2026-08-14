---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Runtime"
type: "preset-list"
componentSlug: "bedrock-agentcore-runtime"
componentTitle: "Bedrock AgentCore Runtime"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-code-bundle-agent"
    rank: "01"
    title: "Code-Bundle Agent"
    excerpt: "This preset hosts a Python agent from an S3 code bundle on the managed runtime — no container image to build or maintain — with public egress and one floating `live` endpoint."
  - slug: "02-container-agent-in-vpc"
    rank: "02"
    title: "Container Agent in a VPC"
    excerpt: "This preset hosts a container-image agent whose sessions attach to your subnets — the pattern for agents that query private databases or internal APIs — with session lifecycle caps and per-session..."
---

# Bedrock AgentCore Runtime Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Runtime. Each preset is a complete manifest you can copy, customize, and deploy.
