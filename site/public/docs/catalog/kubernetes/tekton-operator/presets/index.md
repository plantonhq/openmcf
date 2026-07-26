---
title: "Presets"
description: "Ready-to-deploy configuration presets for Tekton Operator"
type: "preset-list"
componentSlug: "tekton-operator"
componentTitle: "Tekton Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-operator"
    rank: "01"
    title: "Tekton Operator preset"
    excerpt: "The complete install, which is deliberately tiny: the operator is a lifecycle manager, not the product. It brings the Tekton CRDs and the controller that turns a `TektonConfig` declaration into..."
---

# Tekton Operator Presets

Ready-to-deploy configuration presets for Tekton Operator. Each preset is a complete manifest you can copy, customize, and deploy.
