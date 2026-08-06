---
title: "RBAC"
description: "RBAC deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesrbac"
---

# Kubernetes RBAC

Deploys a complete RBAC grant — "give these subjects these permissions in this scope" — through a single declarative manifest. One resource covers what raw Kubernetes splits across Role, RoleBinding, ClusterRole, and ClusterRoleBinding; the IaC module materializes exactly the objects the grant implies.

## What Gets Created

Depending on the grant shape, Planton provisions:

- **Role** (namespace scope) or **ClusterRole** (cluster scope) — when `createRole` defines new rules or an aggregation rule
- **RoleBinding** (namespace scope) or **ClusterRoleBinding** (cluster scope) — when at least one subject is present, pointing every subject at the role (created or existing)
- **Labels/Annotations** — standard Planton tracking labels merged with user-provided labels, applied to every created object

Grant shapes: `createRole` + subjects creates role and binding; `createRole` with no subjects creates only the role (a reusable definition); `existingRole` + subjects creates only the binding.

## Kubernetes RBAC Semantics

- **Deny by default, purely additive**: a request is allowed only if some rule reachable through a binding allows it. There are no deny rules — grants union together, and removing access means removing a grant
- **Built-in roles**: every cluster ships `view`, `edit`, `admin` (bound per-namespace) and `cluster-admin` (bound cluster-wide). Grant them with `existingRole: { name: view, isClusterRole: true }`
- **Aggregation**: a ClusterRole can compose its rules from every ClusterRole matching label selectors — the mechanism behind the built-ins absorbing CRD permissions. Cluster scope only

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config — the credential must itself hold every permission the grant confers (Kubernetes escalation prevention)
- **A Kubernetes namespace** for namespace-scoped grants (existing, or a `KubernetesNamespace` reference)
- **Subjects**: ServiceAccounts existing or deployed in the same chart; users and groups are authenticator-asserted strings needing no object

## Quick Start

Create a file `rbac.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesRbac
metadata:
  name: app-reader
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesRbac.app-reader
spec:
  namespaceScope:
    namespace:
      value: backend
  createRole:
    rules:
      - verbs: ["get", "list", "watch"]
        apiGroups: [""]
        resources: ["pods", "services"]
  subjects:
    - serviceAccount:
        name:
          value: app-identity
```

Deploy:

```shell
planton apply -f rbac.yaml
```

This creates a Role `app-reader` and a RoleBinding in the `backend` namespace, granting the `app-identity` ServiceAccount read access to pods and services.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `spec.namespaceScope` / `spec.clusterScope` | `oneof` | The grant's scope. `namespaceScope` produces Role/RoleBinding; `clusterScope` produces ClusterRole/ClusterRoleBinding. | Exactly one must be set |
| `spec.createRole` / `spec.existingRole` | `oneof` | The role source. `createRole` defines new rules; `existingRole` binds to a role already in the cluster. | Exactly one must be set |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.namespaceScope.namespace` | `StringValueOrRef` | `default` | Namespace the grant applies to. Literal (`{ value: my-namespace }`) or a `KubernetesNamespace` reference. |
| `spec.createRole.name` | `string` | component's `metadata.name` | Name of the created Role/ClusterRole. DNS subdomain. |
| `spec.createRole.rules` | `list` | `[]` | Permission rules (see below). At least one required unless `aggregationRule` is set. |
| `spec.createRole.aggregationRule.clusterRoleSelectors` | `list` | — | Label selectors (`matchLabels` and/or `matchExpressions`); the ClusterRole's rules are continuously composed from every ClusterRole matching ANY selector. Cluster scope only; directly listed rules become controller-managed. |
| `spec.existingRole.name` | `string` | — | Name of the existing role, e.g. `view`, `edit`, `admin`, `cluster-admin`, or any custom role. Requires at least one subject. |
| `spec.existingRole.isClusterRole` | `bool` | `false` | Whether the referenced role is a ClusterRole. Meaningful in namespace scope (how built-ins are granted per-namespace); ignored in cluster scope, where the reference is always a ClusterRole. |
| `spec.subjects[]` | `list` | `[]` | Recipients. Each entry is exactly one of `serviceAccount`, `user`, or `group`. Empty means role-only (requires `createRole`). |
| `spec.subjects[].serviceAccount.name` | `StringValueOrRef` | — | ServiceAccount name — literal or a `KubernetesServiceAccount` reference. |
| `spec.subjects[].serviceAccount.namespace` | `StringValueOrRef` | grant's namespace | The ServiceAccount's namespace. REQUIRED in cluster scope (nothing to default from). |
| `spec.subjects[].user` | `string` | — | User name as asserted by the cluster's authenticator (OIDC claim, certificate CN, cloud IAM mapping). Kubernetes has no User objects. |
| `spec.subjects[].group` | `string` | — | Group name as asserted by the authenticator, e.g. an OIDC groups claim or `system:authenticated`. |
| `spec.labels` / `spec.annotations` | `map<string, string>` | `{}` | Applied to every created RBAC object. |

### Policy Rules (`spec.createRole.rules[]`)

Each rule independently grants verbs over resources OR non-resource URLs (never both in one rule):

| Field | Description |
|-------|-------------|
| `verbs` | Actions granted: `get`, `list`, `watch`, `create`, `update`, `patch`, `delete`, `deletecollection`, special verbs (`impersonate`, `bind`, `escalate`, `use`), or `*`. At least one required. |
| `apiGroups` | API groups of the resources: `""` (core — pods, services, configmaps), `apps`, `batch`, `networking.k8s.io`, a CRD's group, or `*`. |
| `resources` | Lowercase plurals: `pods`, `deployments`, `secrets`; subresources as `pods/log`; `*` for all in the listed groups. |
| `resourceNames` | Optional whitelist of object names. Cannot constrain `create`/`deletecollection` (authorization happens before a name exists). |
| `nonResourceUrls` | URL paths: `/metrics`, `/healthz`, `/api/*` (trailing `*` only as the full final segment). Cluster scope only. |

Schema-enforced constraints: aggregation and `nonResourceUrls` require `clusterScope`; ServiceAccount subjects in `clusterScope` must set `namespace`; every rule must grant something; `existingRole` requires subjects.

## Examples

### Grant Built-in `view` to a Group

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesRbac
metadata:
  name: team-view
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesRbac.team-view
spec:
  namespaceScope:
    namespace:
      value: production
  existingRole:
    name: view
    isClusterRole: true
  subjects:
    - group: dev-team
```

Creates only a RoleBinding: the built-in `view` ClusterRole, confined to the `production` namespace, for everyone the authenticator places in `dev-team`.

### Cluster-Scoped Operator Permissions

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesRbac
metadata:
  name: node-monitor
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesRbac.node-monitor
spec:
  clusterScope: {}
  createRole:
    rules:
      - verbs: ["get", "list", "watch"]
        apiGroups: [""]
        resources: ["nodes", "namespaces"]
      - verbs: ["get"]
        nonResourceUrls: ["/metrics"]
  subjects:
    - serviceAccount:
        name:
          value: monitoring-agent
        namespace:
          value: monitoring
```

Note the ServiceAccount subject sets `namespace` explicitly — required in cluster scope.

### Aggregated ClusterRole (no binding)

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesRbac
metadata:
  name: monitoring-aggregate
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesRbac.monitoring-aggregate
spec:
  clusterScope: {}
  createRole:
    aggregationRule:
      clusterRoleSelectors:
        - matchLabels:
            rbac.example.com/aggregate-to-monitoring: "true"
```

Publishes a ClusterRole whose rules are composed from every ClusterRole carrying the label; other grants bind to it via `existingRole`.

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `roleName` | `string` | The role in the grant — the created Role/ClusterRole, or the existing role bound to |
| `roleKind` | `string` | `Role` or `ClusterRole` |
| `bindingName` | `string` | The created binding; empty when the grant defined a role with no subjects |
| `bindingKind` | `string` | `RoleBinding` or `ClusterRoleBinding`; empty when no binding is created |
| `namespace` | `string` | The namespace the grant applies to; empty for cluster-scoped grants |

## Related Components

- [KubernetesServiceAccount](/docs/catalog/kubernetes/serviceaccount) — the identity most grants target; reference it from `subjects[].serviceAccount.name` to create identity and grant in one run
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the namespace for namespace-scoped grants
