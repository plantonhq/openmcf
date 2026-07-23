---
title: "Namespace App Reader"
description: "This preset grants a workload's ServiceAccount read-only access to pods, services, and ConfigMaps in one namespace. It creates a Role with a single `get/list/watch` rule and a RoleBinding pointing..."
type: "preset"
rank: "01"
presetSlug: "01-namespace-app-reader"
componentSlug: "rbac"
componentTitle: "RBAC"
provider: "kubernetes"
icon: "package"
order: 1
---

# Namespace App Reader

This preset grants a workload's ServiceAccount read-only access to pods, services, and ConfigMaps in one namespace. It creates a Role with a single `get/list/watch` rule and a RoleBinding pointing the ServiceAccount at it — the standard least-privilege grant for an application that inspects its own namespace.

## When to Use

- An application, sidecar, or agent that lists pods/services or reads ConfigMaps in its own namespace (service discovery, config watching, leader-election helpers)
- The default starting grant for any in-cluster workload identity — expand verb-by-verb, resource-by-resource as real needs appear

## Key Configuration Choices

- **`namespaceScope`** — the grant produces Role + RoleBinding, confined to one namespace; permissions never leak beyond it
- **Read-only verbs (`get,list,watch`)** — the canonical reader triple: `get` for single objects, `list` for collections, `watch` for change streams. Add write verbs (`create`, `update`, `patch`, `delete`) individually and deliberately
- **Core-group resources** — `apiGroups: [""]` is the core API group, home of `pods`, `services`, `configmaps`. Note Secrets are deliberately absent: `get` on `secrets` is effectively holding every credential in the namespace
- **ServiceAccount subject without a namespace** — in namespace scope the subject's namespace defaults to the grant's own; the `name` accepts a literal (as here) or a reference to a KubernetesServiceAccount resource, so identity and grant can deploy in one chart
- **Additive semantics** — this grant only adds permissions; it stacks with any other grants the ServiceAccount holds

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Namespace the grant applies to (Role, RoleBinding, and default subject namespace) | Your namespace management |
| `<your-service-account>` | Name of the ServiceAccount receiving the permissions | Your KubernetesServiceAccount resource or the workload's `spec.serviceAccountName` |

Also rename `app-reader` (`metadata.name`, which becomes the Role name) to reflect your workload, e.g. `<app>-reader`.

## Related Presets

- **02-grant-builtin-view** — broader read access via the built-in `view` role, no rules to maintain
- **03-cluster-operator** — cluster-scoped permissions including nodes and `/metrics`
- **04-aggregated-clusterrole** — a label-composed ClusterRole published without a binding
