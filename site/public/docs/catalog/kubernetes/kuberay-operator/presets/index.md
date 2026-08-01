---
title: "Presets"
description: "Ready-to-deploy configuration presets for KubeRay Operator"
type: "preset-list"
componentSlug: "kuberay-operator"
componentTitle: "KubeRay Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default"
    rank: "01"
    title: "Default preset"
    excerpt: "The standard operator install: the pinned chart (1.6.2 = operator v1.6.2) into its own `ray-system` namespace, watching every namespace on the cluster, leader election on (the chart default — safe..."
  - slug: "02-private-mirror"
    rank: "02"
    title: "Private-mirror preset"
    excerpt: "The air-gapped posture for the operator: `imageRegistry` replaces only the registry part of the operator's own image (the path stays the upstream `kuberay/operator`; the default registry is quay.io)..."
---

# KubeRay Operator Presets

Ready-to-deploy configuration presets for KubeRay Operator. Each preset is a complete manifest you can copy, customize, and deploy.
