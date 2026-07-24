---
title: "Presets"
description: "Ready-to-deploy configuration presets for Solr Operator"
type: "preset-list"
componentSlug: "solr-operator"
componentTitle: "Solr Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard preset"
    excerpt: "The default install: one cluster-wide Solr operator plus the bundled zookeeper-operator, on the stable 0.9.1 chart. This is the shape for a cluster you administer — install once, and every..."
  - slug: "02-existing-zookeeper-operator"
    rank: "02"
    title: "Existing zookeeper operator preset"
    excerpt: "The install for clusters that already run a zookeeper-operator: `install: false` skips the bundled dependency (whose fixed-name, cluster-scoped RBAC would conflict with the existing one), and..."
  - slug: "03-namespace-fenced"
    rank: "03"
    title: "Namespace fenced preset"
    excerpt: "An operator that reconciles KubernetesSolr resources only in an explicit namespace list, sized for a governed platform: two leader-elected replicas (one active, one warm standby) and declared..."
---

# Solr Operator Presets

Ready-to-deploy configuration presets for Solr Operator. Each preset is a complete manifest you can copy, customize, and deploy.
