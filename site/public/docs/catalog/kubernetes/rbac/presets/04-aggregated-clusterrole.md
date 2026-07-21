---
title: "Aggregated ClusterRole"
description: "This preset publishes a ClusterRole whose rules are continuously composed by the RBAC controller from every ClusterRole carrying a matching label — with NO subjects, so no binding is created. It is a..."
type: "preset"
rank: "04"
presetSlug: "04-aggregated-clusterrole"
componentSlug: "rbac"
componentTitle: "RBAC"
provider: "kubernetes"
icon: "package"
order: 4
---

# Aggregated ClusterRole

This preset publishes a ClusterRole whose rules are continuously composed by the RBAC controller from every ClusterRole carrying a matching label — with NO subjects, so no binding is created. It is a role definition other grants bind to later, and other teams extend by shipping labeled ClusterRoles.

## When to Use

- Building a plugin-style permission set: define `monitoring-aggregate` once, and let each component contribute its own permissions by shipping a small ClusterRole labeled `rbac.example.com/aggregate-to-monitoring: "true"` — the aggregate grows automatically, with no central edits
- Mirroring how the built-ins work: `view`/`edit`/`admin` are aggregated ClusterRoles that absorb CRD permissions exactly this way (via `rbac.authorization.k8s.io/aggregate-to-view` and friends)
- Publishing a reusable role for later binding — bind it from another grant with `existingRole: { name: monitoring-aggregate, isClusterRole: true }`

## Key Configuration Choices

- **`clusterScope` is mandatory** — aggregation exists only on ClusterRoles; the Kubernetes API has no `aggregationRule` on namespaced Roles, and the schema enforces this
- **`aggregationRule` instead of `rules`** — the controller owns the rules: any ClusterRole matching ANY selector contributes its rules to this one, continuously. Directly listed `rules` would be overwritten by the controller, so a role definition provides rules OR aggregation, not both meaningfully
- **`matchLabels` selector** — exact-match label requirements; a ClusterRole must carry every listed key/value pair to match. `matchExpressions` (set-based: `In`, `NotIn`, `Exists`, `DoesNotExist`) are available for selections exact-match cannot express, and AND together with `matchLabels` within one selector
- **NO subjects** — deliberately: this grant creates only the role definition. Omitting subjects is valid with `createRole` (with `existingRole` it would deploy nothing and is rejected)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `rbac.example.com/aggregate-to-monitoring` | Your aggregation label key — pick a domain you own plus an `aggregate-to-<role>` suffix, by convention | Your own naming; contributing ClusterRoles must carry the identical label |

Also rename `monitoring-aggregate` (`metadata.name`, which becomes the ClusterRole name) to reflect the permission set you are composing.

## Related Presets

- **01-namespace-app-reader** — namespace-confined read access for a workload
- **02-grant-builtin-view** — bind an existing (built-in) role to a group
- **03-cluster-operator** — cluster-scoped rules with subjects, bound directly
