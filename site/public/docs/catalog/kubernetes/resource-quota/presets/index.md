---
title: "Presets"
description: "Ready-to-deploy configuration presets for Resource Quota"
type: "preset-list"
componentSlug: "resource-quota"
componentTitle: "Resource Quota"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-namespace-governed"
    rank: "01"
    title: "Team Namespace Governed"
    excerpt: "This preset is the full governance pair and the safe pattern for compute caps: aggregate CPU/memory caps on the namespace (the ResourceQuota) paired with per-container defaults (the companion..."
  - slug: "02-object-count-caps"
    rank: "02"
    title: "Object Count Caps"
    excerpt: "This preset caps how MANY objects a namespace may hold — pods, Services, and PersistentVolumeClaims — without touching compute. It is the safest quota to introduce on a live namespace: object counts..."
  - slug: "03-besteffort-guard"
    rank: "03"
    title: "BestEffort Guard"
    excerpt: "This preset caps only the namespace's BestEffort pods — pods with no requests or limits at all — at 10. It is a scoped quota: the `best_effort` scope makes the quota track exclusively the naive pods,..."
---

# Resource Quota Presets

Ready-to-deploy configuration presets for Resource Quota. Each preset is a complete manifest you can copy, customize, and deploy.
