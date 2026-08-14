---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cognitive Account"
type: "preset-list"
componentSlug: "cognitive-account"
componentTitle: "Cognitive Account"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-azure-openai-account"
    rank: "01"
    title: "Azure OpenAI Account"
    excerpt: "This preset creates an `OpenAI`-kind S0 account -- the container Azure OpenAI model deployments (gpt-4o, embeddings) are created onto. The account object itself carries no idle cost; billing follows..."
  - slug: "02-ai-foundry-account"
    rank: "02"
    title: "AI Foundry Account"
    excerpt: "This preset creates an `AIServices`-kind account with project management enabled and a system-assigned identity -- the foundation AI Foundry team workspaces (AzureCognitiveAccountProject) and agents..."
  - slug: "03-private-hardened-account"
    rank: "03"
    title: "Private Hardened Account"
    excerpt: "This preset locks the account down on every axis the service offers: deny-by-default network ACLs with explicit IP and subnet allowances, Entra-ID-only authentication (keys disabled), and restricted..."
---

# Cognitive Account Presets

Ready-to-deploy configuration presets for Cognitive Account. Each preset is a complete manifest you can copy, customize, and deploy.
