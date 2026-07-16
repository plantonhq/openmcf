---
title: "Cluster Viewer"
description: "This preset grants a role read-only access across the whole cluster through the AWS-managed view policy -- no in-cluster RBAC objects, no ConfigMap edits."
type: "preset"
rank: "01"
presetSlug: "01-cluster-viewer"
componentSlug: "eks-access-entry"
componentTitle: "EKS Access Entry"
provider: "aws"
icon: "package"
order: 1
---

# Cluster Viewer

This preset grants a role read-only access across the whole cluster
through the AWS-managed view policy -- no in-cluster RBAC objects, no
ConfigMap edits.

## When to Use

- The default grant for engineers who inspect but do not operate
- Read-only CI (config drift checks, inventory) and dashboards
- The first access entry on a new cluster, before finer scopes exist

## Key Configuration Choices

- **`AmazonEKSViewPolicy`, cluster scope** -- read everything, change
  nothing; escalate per-namespace with `AmazonEKSAdminPolicy` instead
  of widening this grant
- **No `kubernetesGroups`** -- the AWS-managed policy is the whole
  authorization story; add groups only when your own RBAC bindings
  define something finer
- **No `userName`** -- AWS's session-templated default preserves the
  actual session name in audit logs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<entry-resource-name>` | Name for this grant (e.g. `platform-viewers`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) | Your deployment region |
| `<cluster-resource-name>` | Name of the AwsEksCluster resource | Your cluster manifest's `metadata.name` |
| `<principal-role-resource-name>` | Name of the AwsIamRole being granted access | Your role manifest's `metadata.name` |

## Common Additions

- A second association scoping `AmazonEKSEditPolicy` to specific
  namespaces for teams that deploy there

## Related Presets

- **02-namespace-admin** -- full admin inside a team's namespaces only
- **03-rbac-groups** -- bring-your-own RBAC group mapping
