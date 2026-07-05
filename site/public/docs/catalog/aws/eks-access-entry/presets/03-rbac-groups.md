---
title: "RBAC Group Mapping"
description: "This preset maps a principal onto Kubernetes RBAC groups you define -- the entry handles authentication; your own (Cluster)RoleBindings decide what those groups may do. The..."
type: "preset"
rank: "03"
presetSlug: "03-rbac-groups"
componentSlug: "eks-access-entry"
componentTitle: "EKS Access Entry"
provider: "aws"
icon: "package"
order: 3
---

# RBAC Group Mapping

This preset maps a principal onto Kubernetes RBAC groups you define --
the entry handles authentication; your own (Cluster)RoleBindings decide
what those groups may do. The bring-your-own-authorization pattern.

## When to Use

- Clusters with existing, carefully-built RBAC that AWS-managed
  policies cannot express (custom verbs, resource names, aggregated
  roles)
- Migrating aws-auth ConfigMap `mapRoles` entries -- the group mapping
  carries over one-to-one, minus the ConfigMap fragility
- Uniform group names across clusters, with per-cluster bindings

## Key Configuration Choices

- **`kubernetesGroups`** -- names only; nothing is created in-cluster.
  A group with no binding grants nothing (fail-safe). `system:` groups
  are rejected at validation, as AWS forbids them.
- **`userName` with `{{SessionName}}`** -- predictable prefix for
  bindings, real session name preserved for audit; omit the field
  entirely to keep AWS's default template
- **No `policyAssociations`** -- pure RBAC-side authorization; both
  paths can combine on one entry when needed

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<entry-resource-name>` | Name for this grant (e.g. `platform-operators`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |
| `<principal-role-resource-name>` | Name of the AwsIamRole being granted access | Your role manifest's `metadata.name` |
| `<rbac-group-name>` | The group your RBAC bindings reference | Your cluster's RBAC manifests |
| `<prefix>` | A username prefix for audit logs (e.g. `ops`) | Your naming convention |

## Common Additions

- A cluster-scoped `AmazonEKSViewPolicy` association as a read-only
  floor beneath the custom RBAC

## Related Presets

- **01-cluster-viewer** -- AWS-managed read-only, no RBAC needed
- **02-namespace-admin** -- AWS-managed namespace delegation
