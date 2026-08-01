---
title: "Presets"
description: "Ready-to-deploy configuration presets for Keycloak Operator"
type: "preset-list"
componentSlug: "keycloak-operator"
componentTitle: "Keycloak Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-operator"
    rank: "01"
    title: "Operator preset"
    excerpt: "The official Keycloak Operator in its default posture: installed into the `keycloak` namespace with the NAMESPACED watch — the operator reconciles only Keycloak declarations living beside it, so..."
---

# Keycloak Operator Presets

Ready-to-deploy configuration presets for Keycloak Operator. Each preset is a complete manifest you can copy, customize, and deploy.
