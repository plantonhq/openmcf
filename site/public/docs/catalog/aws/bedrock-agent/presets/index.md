---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Agent"
type: "preset-list"
componentSlug: "bedrock-agent"
componentTitle: "Bedrock Agent"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-support-agent-with-tools"
    rank: "01"
    title: "Support Agent with Tools"
    excerpt: "This preset creates a customer-support agent on Amazon Nova Micro with a return-control order-lookup tool, the reserved `AMAZON.UserInput` capability (so the agent can ask clarifying questions), and..."
  - slug: "02-supervisor-with-collaborators"
    rank: "02"
    title: "Supervisor with Collaborators"
    excerpt: "This preset creates a multi-agent supervisor: it answers general questions itself (with product-docs retrieval), delegates billing questions to a specialist agent through its `live` alias, and..."
---

# Bedrock Agent Presets

Ready-to-deploy configuration presets for Bedrock Agent. Each preset is a complete manifest you can copy, customize, and deploy.
