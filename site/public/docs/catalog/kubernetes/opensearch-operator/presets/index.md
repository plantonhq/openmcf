---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenSearch Operator"
type: "preset-list"
componentSlug: "opensearch-operator"
componentTitle: "OpenSearch Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard preset"
    excerpt: "The default install: one cluster-wide operator watching every namespace, the stable 2.8.0 chart/operator pairing, chart-default manager sizing. This is the shape for a cluster you administer —..."
  - slug: "02-namespace-scoped"
    rank: "02"
    title: "Namespace scoped preset"
    excerpt: "An operator fenced to a single namespace with namespace-scoped RBAC: `watch_namespace` limits what it reconciles, `use_role_bindings` swaps ClusterRoleBindings for RoleBindings in the watched..."
  - slug: "03-private-mirror"
    rank: "03"
    title: "Private mirror preset"
    excerpt: "The standard cluster-wide operator, but every image byte comes from your own registry: an `image` override pointing at the mirrored manager image, `image_pull_secrets` naming the registry credential,..."
---

# OpenSearch Operator Presets

Ready-to-deploy configuration presets for OpenSearch Operator. Each preset is a complete manifest you can copy, customize, and deploy.
