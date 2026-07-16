---
title: "Presets"
description: "Ready-to-deploy configuration presets for EKS Access Entry"
type: "preset-list"
componentSlug: "eks-access-entry"
componentTitle: "EKS Access Entry"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-cluster-viewer"
    rank: "01"
    title: "Cluster Viewer"
    excerpt: "This preset grants a role read-only access across the whole cluster through the AWS-managed view policy -- no in-cluster RBAC objects, no ConfigMap edits."
  - slug: "02-namespace-admin"
    rank: "02"
    title: "Namespace Admin"
    excerpt: "This preset delegates a team full admin inside its own namespaces plus read-only visibility across the cluster -- the multi-tenant delegation pattern, expressed entirely through AWS-managed policies."
  - slug: "03-rbac-groups"
    rank: "03"
    title: "RBAC Group Mapping"
    excerpt: "This preset maps a principal onto Kubernetes RBAC groups you define -- the entry handles authentication; your own (Cluster)RoleBindings decide what those groups may do. The..."
---

# EKS Access Entry Presets

Ready-to-deploy configuration presets for EKS Access Entry. Each preset is a complete manifest you can copy, customize, and deploy.
