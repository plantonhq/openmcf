---
title: "Presets"
description: "Ready-to-deploy configuration presets for Role Assignment"
type: "preset-list"
componentSlug: "role-assignment"
componentTitle: "Role Assignment"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-identity-resource-group-grant"
    rank: "01"
    title: "Managed Identity Grant on a Resource Group"
    excerpt: "This preset grants a user-assigned managed identity a built-in role on a resource group -- the most common authorization pattern in composed Azure environments: the identity a workload runs as gets..."
  - slug: "02-abac-conditioned-grant"
    rank: "02"
    title: "ABAC-Conditioned Data Grant"
    excerpt: "This preset grants a data-plane role narrowed by an Azure attribute-based access control (ABAC) condition -- the role's permissions apply only when the condition evaluates true. The template shows..."
  - slug: "03-custom-role-subscription-grant"
    rank: "03"
    title: "Custom Role at Subscription Scope"
    excerpt: "This preset assigns a custom role -- referenced by its role definition ID, the exact form custom roles require -- to a group at subscription scope. This is the governance-team pattern: define a..."
---

# Role Assignment Presets

Ready-to-deploy configuration presets for Role Assignment. Each preset is a complete manifest you can copy, customize, and deploy.
