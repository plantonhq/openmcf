---
title: "Presets"
description: "Ready-to-deploy configuration presets for RBAC"
type: "preset-list"
componentSlug: "rbac"
componentTitle: "RBAC"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-namespace-app-reader"
    rank: "01"
    title: "Namespace App Reader"
    excerpt: "This preset grants a workload's ServiceAccount read-only access to pods, services, and ConfigMaps in one namespace. It creates a Role with a single `get/list/watch` rule and a RoleBinding pointing..."
  - slug: "02-grant-builtin-view"
    rank: "02"
    title: "Grant Built-in `view` to a Group"
    excerpt: "This preset binds the Kubernetes built-in `view` ClusterRole to a group, confined to one namespace. No role is created — only a RoleBinding — giving everyone in the group read-only access to most..."
  - slug: "03-cluster-operator"
    rank: "03"
    title: "Cluster Operator"
    excerpt: "This preset grants an operator-style ServiceAccount cluster-wide read access to nodes and namespaces, plus the `/metrics` endpoint. It creates a ClusterRole with two rules and a ClusterRoleBinding —..."
  - slug: "04-aggregated-clusterrole"
    rank: "04"
    title: "Aggregated ClusterRole"
    excerpt: "This preset publishes a ClusterRole whose rules are continuously composed by the RBAC controller from every ClusterRole carrying a matching label — with NO subjects, so no binding is created. It is a..."
---

# RBAC Presets

Ready-to-deploy configuration presets for RBAC. Each preset is a complete manifest you can copy, customize, and deploy.
