---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock AgentCore Token Vault"
type: "preset-list"
componentSlug: "bedrock-agentcore-token-vault"
componentTitle: "Bedrock AgentCore Token Vault"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-customer-managed-key"
    rank: "01"
    title: "Customer-Managed Key"
    excerpt: "This preset encrypts the region's AgentCore token vault with your own KMS key — you control rotation, policy, and revocation for every credential agents store."
  - slug: "02-service-managed-revert"
    rank: "02"
    title: "Service-Managed Key (Revert)"
    excerpt: "This preset returns the region's token vault to AWS-owned encryption — the default posture, and the ONLY way back from a customer-managed key (destroying the component does not revert the setting)."
---

# Bedrock AgentCore Token Vault Presets

Ready-to-deploy configuration presets for Bedrock AgentCore Token Vault. Each preset is a complete manifest you can copy, customize, and deploy.
