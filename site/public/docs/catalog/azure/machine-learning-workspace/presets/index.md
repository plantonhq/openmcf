---
title: "Presets"
description: "Ready-to-deploy configuration presets for Machine Learning Workspace"
type: "preset-list"
componentSlug: "machine-learning-workspace"
componentTitle: "Machine Learning Workspace"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-team-workspace"
    rank: "01"
    title: "Team Workspace"
    excerpt: "This preset deploys a standard team workspace with a system-assigned identity on the three required companion services -- the starting point for most ML estates: public access on (the default), no..."
  - slug: "02-feature-store"
    rank: "02"
    title: "Feature Store"
    excerpt: "This preset deploys a FEATURE_STORE-flavor workspace -- the managed feature store backing online/offline feature serving for ML pipelines. Same companion services as a regular workspace; the flavor..."
  - slug: "03-private-hardened-workspace"
    rank: "03"
    title: "Private Hardened Workspace"
    excerpt: "This preset deploys a locked-down workspace: public network access off, managed-network isolation in approved-outbound mode, and an explicit outbound allowlist (package indexes plus Key Vault). The..."
---

# Machine Learning Workspace Presets

Ready-to-deploy configuration presets for Machine Learning Workspace. Each preset is a complete manifest you can copy, customize, and deploy.
