---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Tools"
type: "preset-list"
componentSlug: "bedrock-agentcore-tools"
componentTitle: "Bedrock AgentCore Tools"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-sandboxed-code-interpreter"
    rank: "01"
    title: "Sandboxed Code Interpreter"
    excerpt: "This preset gives an agent the safest possible compute: a code interpreter with NO network access and no AWS credentials — model-written code can calculate, parse, and plot, and nothing else."
  - slug: "02-recorded-research-browser"
    rank: "02"
    title: "Recorded Research Browser"
    excerpt: "This preset gives an agent a cloud browser whose every session is recorded to S3 and whose traffic is cryptographically signed — plus a saved profile so sessions start already logged in to your docs..."
---

# Bedrock AgentCore Tools Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Tools. Each preset is a complete manifest you can copy, customize, and deploy.
