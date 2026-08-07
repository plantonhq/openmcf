---
title: "Presets"
description: "Ready-to-deploy configuration presets for Altinity ClickHouse Operator"
type: "preset-list"
componentSlug: "altinity-clickhouse-operator"
componentTitle: "Altinity ClickHouse Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard preset"
    excerpt: "The one-per-cluster engine install: the Altinity operator watching every namespace, with real operator credentials and modest container sizing. Declare it once; every KubernetesClickHouse resource in..."
  - slug: "02-namespace-scoped"
    rank: "02"
    title: "Namespace-scoped preset"
    excerpt: "The tenancy posture: an operator that watches exactly one namespace and holds only namespace-scoped Roles instead of cluster-wide RBAC. On a shared cluster this is what lets a team run their own..."
  - slug: "03-private-mirror"
    rank: "03"
    title: "Private mirror preset"
    excerpt: "The air-gap posture: every image the install pulls re-pointed at a private registry, with a pull secret for all of them. Beyond the usual mirror hygiene there is one image here that genuinely needs..."
---

# Altinity ClickHouse Operator Presets

Ready-to-deploy configuration presets for Altinity ClickHouse Operator. Each preset is a complete manifest you can copy, customize, and deploy.
