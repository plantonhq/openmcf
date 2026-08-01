---
title: "Presets"
description: "Ready-to-deploy configuration presets for Kyverno"
type: "preset-list"
componentSlug: "kyverno"
componentTitle: "Kyverno"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-cluster"
    rank: "01"
    title: "Dev cluster preset"
    excerpt: "The smallest declarable Kyverno: the chart defaults end to end. All four controllers run single-replica, the engine generates and rotates its own webhook certificates (no prerequisites), the policy..."
  - slug: "02-production-ha"
    rank: "02"
    title: "Production HA preset"
    excerpt: "Kyverno sized for a cluster where admission availability matters. The admission controller runs three replicas — it sits on the cluster's WRITE PATH, and with the default fail-closed policy posture..."
  - slug: "03-airgapped-mirror"
    rank: "03"
    title: "Air-gapped mirror preset"
    excerpt: "Kyverno for clusters that cannot reach public registries. One field does the heavy lifting: `imageRegistry` sets the chart's global registry, rerouting EVERY Kyverno image — the four controllers, the..."
---

# Kyverno Presets

Ready-to-deploy configuration presets for Kyverno. Each preset is a complete manifest you can copy, customize, and deploy.
