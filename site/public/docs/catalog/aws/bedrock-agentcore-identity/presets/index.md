---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Identity"
type: "preset-list"
componentSlug: "bedrock-agentcore-identity"
componentTitle: "Bedrock AgentCore Identity"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-vaulted-credentials"
    rank: "01"
    title: "Vaulted Credentials"
    excerpt: "This preset vaults the two credential shapes agent tools most commonly need — an API key and a GitHub OAuth client — as managed-secret references resolved just-in-time at deploy. Gateway targets and..."
  - slug: "02-cedar-authorization"
    rank: "02"
    title: "Cedar Authorization"
    excerpt: "This preset stands up a Cedar policy engine with a permit/forbid pair — read tools for everyone, mutations denied off-hours — ready to attach to a gateway in LOG_ONLY mode."
---

# Bedrock AgentCore Identity Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Identity. Each preset is a complete manifest you can copy, customize, and deploy.
