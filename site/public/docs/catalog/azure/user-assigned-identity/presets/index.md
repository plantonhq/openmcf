---
title: "Presets"
description: "Ready-to-deploy configuration presets for User Assigned Identity"
type: "preset-list"
componentSlug: "user-assigned-identity"
componentTitle: "User Assigned Identity"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard Managed Identity"
    excerpt: "This preset creates a plain user-assigned managed identity -- the anchor of every keyless-auth story on Azure. The identity is deliberately just the identity: what it may do and who may act as it are..."
  - slug: "02-ci-deployer"
    rank: "02"
    title: "CI Deployer Identity"
    excerpt: "This preset creates the identity at the center of keyless CI/CD -- the identity a pipeline authenticates as when it deploys to Azure. It is the first of three composable pieces; the other two attach..."
  - slug: "03-governance-tagged"
    rank: "03"
    title: "Governance-Tagged, Regionally Isolated Identity"
    excerpt: "This preset creates an identity shaped for regulated environments: regional isolation restricts token issuance to the identity's own region (a data-residency / blast-radius control), and governance..."
---

# User Assigned Identity Presets

Ready-to-deploy configuration presets for User Assigned Identity. Each preset is a complete manifest you can copy, customize, and deploy.
