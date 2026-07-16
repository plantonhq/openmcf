---
title: "Namespace Admin"
description: "This preset delegates a team full admin inside its own namespaces plus read-only visibility across the cluster -- the multi-tenant delegation pattern, expressed entirely through AWS-managed policies."
type: "preset"
rank: "02"
presetSlug: "02-namespace-admin"
componentSlug: "eks-access-entry"
componentTitle: "EKS Access Entry"
provider: "aws"
icon: "package"
order: 2
---

# Namespace Admin

This preset delegates a team full admin inside its own namespaces plus
read-only visibility across the cluster -- the multi-tenant delegation
pattern, expressed entirely through AWS-managed policies.

## When to Use

- Namespace-per-team clusters where teams operate their own workloads
- Delegating deploy rights to a service's CI role without cluster-wide
  power
- Replacing hand-built RBAC roles that reimplemented "admin, but only
  here"

## Key Configuration Choices

- **`AmazonEKSAdminPolicy` with namespace scope** -- admin stops at the
  namespace boundary; add namespaces to the list as the team grows
  (associations update in place)
- **`AmazonEKSViewPolicy` cluster-wide alongside** -- teams can see
  shared context (nodes, CRDs, other namespaces' service endpoints)
  without touching it; drop this association for stricter isolation
- **Two associations, one entry** -- each materializes as its own
  provider resource keyed by policy name, so re-scoping one never
  disturbs the other

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<entry-resource-name>` | Name for this grant (e.g. `team-a-admin`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |
| `<team-role-resource-name>` | Name of the team's AwsIamRole | Your role manifest's `metadata.name` |
| `<team-namespace>` | The team's namespace (wildcards like `team-a-*` allowed) | Your cluster's namespace layout |

## Common Additions

- More namespaces on the admin association's list
- `AmazonEKSEditPolicy` instead of AdminPolicy for deploy-but-not-
  configure delegation

## Related Presets

- **01-cluster-viewer** -- the read-only baseline grant
- **03-rbac-groups** -- bring-your-own RBAC group mapping
