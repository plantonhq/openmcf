# KubernetesRbac

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

**KubernetesRbacSpec** defines a complete RBAC grant: "give these subjects these permissions
in this scope". One resource bundles the role definition and its binding — the unit in which
RBAC is actually reasoned about — instead of splitting Role, RoleBinding, ClusterRole, and
ClusterRoleBinding into four coordinating resources.

Three orthogonal choices shape a grant:

  1. **Scope** — `namespace_scope` produces Role/RoleBinding objects (permissions inside one
     namespace); `cluster_scope` produces ClusterRole/ClusterRoleBinding objects
     (cluster-wide permissions, or rules over cluster-scoped resources and non-resource URLs).
  2. **Role** — `create_role` defines new PolicyRules (optionally with ClusterRole
     aggregation); `existing_role` binds to a role that already exists, most commonly the
     built-in "view" / "edit" / "admin" / "cluster-admin" ClusterRoles.
  3. **Subjects** — who receives the permissions: ServiceAccounts (referenced, so charts can
     create the identity and its grant in one run), users, or groups. Omitting subjects
     creates the role definition with no binding — useful for publishing a reusable role
     that other grants bind to later.

Kubernetes authorizes by evaluating ClusterRoleBindings first, then RoleBindings in the
request's namespace, denying by default — every grant is purely additive; there are no
"deny" rules in Kubernetes RBAC.

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesRbac
metadata:
  name: test-rbac
spec:
  namespaceScope:
    namespace:
      value: default
  createRole:
    rules:
      - verbs: ["get", "list"]
        apiGroups: [""]
        resources: ["configmaps"]
  subjects:
    - serviceAccount:
        name:
          value: test-sa
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespaceScope` | `KubernetesRbacNamespaceScope` |  |  |  |
| `spec.namespaceScope.namespace` | `string \| valueFrom` |  |  | KubernetesNamespace (`spec.name`) |
| `spec.clusterScope` | `KubernetesRbacClusterScope` |  |  |  |
| `spec.createRole` | `KubernetesRbacRoleDefinition` |  |  |  |
| `spec.createRole.name` | `string` |  |  |  |
| `spec.createRole.rules` | `[]KubernetesRbacPolicyRule` |  |  |  |
| `spec.createRole.rules[].verbs` | `[]string` | yes |  |  |
| `spec.createRole.rules[].apiGroups` | `[]string` |  |  |  |
| `spec.createRole.rules[].resources` | `[]string` |  |  |  |
| `spec.createRole.rules[].resourceNames` | `[]string` |  |  |  |
| `spec.createRole.rules[].nonResourceUrls` | `[]string` |  |  |  |
| `spec.createRole.aggregationRule` | `KubernetesRbacAggregationRule` |  |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors` | `[]KubernetesRbacLabelSelector` | yes |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors[].matchLabels` | `map<string, string>` |  |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions` | `[]KubernetesRbacLabelSelectorRequirement` |  |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].key` | `string` | yes |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].operator` | `string` |  |  |  |
| `spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].values` | `[]string` |  |  |  |
| `spec.existingRole` | `KubernetesRbacExistingRole` |  |  |  |
| `spec.existingRole.name` | `string` | yes |  |  |
| `spec.existingRole.isClusterRole` | `bool` |  |  |  |
| `spec.subjects` | `[]KubernetesRbacSubject` |  |  |  |
| `spec.subjects[].serviceAccount` | `KubernetesRbacServiceAccountSubject` |  |  |  |
| `spec.subjects[].serviceAccount.name` | `string \| valueFrom` | yes |  | KubernetesServiceAccount (`spec.name`) |
| `spec.subjects[].serviceAccount.namespace` | `string \| valueFrom` |  |  | KubernetesNamespace (`spec.name`) |
| `spec.subjects[].user` | `string` |  |  |  |
| `spec.subjects[].group` | `string` |  |  |  |
| `spec.labels` | `map<string, string>` |  |  |  |
| `spec.annotations` | `map<string, string>` |  |  |  |

## Field Details

### spec.namespaceScope

`KubernetesRbacNamespaceScope`

Grant applies within a single namespace (Role + RoleBinding).

### spec.namespaceScope.namespace

`string | valueFrom`

The namespace the grant applies to. Accepts a literal namespace name or a reference to a
KubernetesNamespace resource, so a chart can create the namespace and its access grants
in one run. When omitted, the grant lands in the cluster's `default` namespace.

- references: KubernetesNamespace (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.clusterScope

`KubernetesRbacClusterScope`

Grant applies cluster-wide (ClusterRole + ClusterRoleBinding).

### spec.createRole

`KubernetesRbacRoleDefinition`

Create a new Role (namespace scope) or ClusterRole (cluster scope) with the given rules.

- rule: at least one rule is required unless aggregation_rule is set

### spec.createRole.name

`string` · optional (explicit presence)

The name of the created Role/ClusterRole (its `metadata.name`). When omitted, the
component's own metadata.name is used. Must be a valid Kubernetes object name.

- rule: Role name must be a valid DNS subdomain (lowercase alphanumeric, hyphens, and dots, no leading/trailing dots or hyphens)
- rule: {"string":{"maxLen":"253"}}

### spec.createRole.rules

`[]KubernetesRbacPolicyRule`

The permission rules. Each rule independently grants a set of verbs over a set of
resources (or non-resource URLs); there is no rule ordering and no deny semantics.
May be empty only when `aggregation_rule` is set (the aggregation controller then owns
the rules).

- rule: a rule grants either resources or non_resource_urls, not both
- rule: a rule must list resources (with api_groups) or non_resource_urls

### spec.createRole.rules[].verbs

`[]string` · required

The actions granted, e.g. "get", "list", "watch", "create", "update", "patch", "delete",
"deletecollection" — or special verbs like "impersonate", "bind", "escalate", "use".
"*" grants all verbs. At least one is required.

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.createRole.rules[].apiGroups

`[]string`

The API groups containing the resources, e.g. "" (core group, for pods/services/
configmaps), "apps", "batch", "networking.k8s.io", or a CRD's group. "*" matches all
groups. Required for resource rules; leave empty for non-resource URL rules.

### spec.createRole.rules[].resources

`[]string`

The resources the verbs apply to, as lowercase plurals: "pods", "deployments",
"secrets". Subresources use a slash: "pods/log", "deployments/scale". "*" matches all
resources in the listed groups. Required for resource rules; leave empty for
non-resource URL rules.

### spec.createRole.rules[].resourceNames

`[]string`

Optional whitelist of object names the rule is limited to, e.g. a single ConfigMap's
name. Empty means all objects of the listed resources. Name restrictions cannot apply
to "create" or "deletecollection" (Kubernetes evaluates those before a name exists).

### spec.createRole.rules[].nonResourceUrls

`[]string`

Non-resource URL paths this rule grants access to, e.g. "/healthz", "/metrics",
"/api/*". A trailing "*" wildcard is allowed only as the full final path segment.
Cluster scope only — these paths exist outside any namespace. A rule grants either
resources or non-resource URLs, never both.

### spec.createRole.aggregationRule

`KubernetesRbacAggregationRule`

ClusterRole aggregation (cluster scope only). Instead of listing rules directly, the
ClusterRole's rules are continuously composed by the controller from every ClusterRole
matching these label selectors — the mechanism behind the built-in "view"/"edit"/"admin"
roles absorbing CRD permissions. When set, directly listed `rules` are controller-managed
and will be overwritten.

### spec.createRole.aggregationRule.clusterRoleSelectors

`[]KubernetesRbacLabelSelector` · required

Label selectors locating the ClusterRoles to aggregate. A ClusterRole matching ANY
selector contributes its rules. At least one selector is required.

- rule: {"repeated":{"minItems":"1"}}

### spec.createRole.aggregationRule.clusterRoleSelectors[].matchLabels

`map<string, string>`

Exact-match label requirements: every key/value pair listed must be present on the
ClusterRole, e.g. {"rbac.example.com/aggregate-to-monitoring": "true"}.

### spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions

`[]KubernetesRbacLabelSelectorRequirement`

Set-based label requirements, for selections exact-match cannot express
(key existence, value-in-set).

- rule: values must be non-empty for In/NotIn and empty for Exists/DoesNotExist

### spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].key

`string` · required

The label key the requirement applies to.

- rule: {"string":{"minLen":"1"}}

### spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].operator

`string`

The relationship between the key and the values: "In" (value must be one of `values`),
"NotIn" (must not be), "Exists" (key present, `values` must be empty), or
"DoesNotExist" (key absent, `values` must be empty).

- rule: {"string":{"in":["In","NotIn","Exists","DoesNotExist"]}}

### spec.createRole.aggregationRule.clusterRoleSelectors[].matchExpressions[].values

`[]string`

The values compared against the label's value. Required (non-empty) for In/NotIn;
must be empty for Exists/DoesNotExist.

### spec.existingRole

`KubernetesRbacExistingRole`

Bind to a role that already exists in the cluster instead of creating one.

### spec.existingRole.name

`string` · required

The name of the existing role, e.g. "view", "edit", "admin", "cluster-admin", or any
custom role already in the cluster.

- rule: {"string":{"minLen":"1","maxLen":"253"}}

### spec.existingRole.isClusterRole

`bool`

Whether the referenced role is a ClusterRole. Only meaningful in namespace scope, where
a RoleBinding may point at either a namespaced Role (false) or a ClusterRole (true —
how built-in roles like "view" are granted per-namespace). In cluster scope the
reference is always a ClusterRole and this flag is ignored.

### spec.subjects

`[]KubernetesRbacSubject`

Who receives the permissions. Each subject is a ServiceAccount, user, or group. When at
least one subject is present, a RoleBinding (namespace scope) or ClusterRoleBinding
(cluster scope) is created pointing every subject at the role. When empty, only the role
definition is created (requires `create_role` — binding to an existing role with no
subjects would deploy nothing).

- rule: Exactly one subject must be set (service_account, user, or group)

### spec.subjects[].serviceAccount

`KubernetesRbacServiceAccountSubject`

A Kubernetes ServiceAccount — the subject kind for workloads.

### spec.subjects[].serviceAccount.name

`string | valueFrom` · required

The ServiceAccount's name. Accepts a literal name or a reference to a
KubernetesServiceAccount resource, so a chart can create the identity and its
permissions in one run.

- references: KubernetesServiceAccount (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesServiceAccount, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.subjects[].serviceAccount.namespace

`string | valueFrom`

The ServiceAccount's namespace. When omitted, defaults to the grant's own namespace
(namespace scope). Cluster-scoped grants must set it explicitly — a ServiceAccount
always lives in some namespace, and cluster scope provides no default.

containment_exempt: locates the SUBJECT being granted access, not the
grant itself — the grant's own home is the scope's namespace field.

- references: KubernetesNamespace (`spec.name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.subjects[].user

`string`

A user name as asserted by the cluster's authenticator (OIDC claim, certificate CN,
cloud IAM principal mapping). Kubernetes has no User objects — this is a plain string
matched at authentication time.

### spec.subjects[].group

`string`

A group name as asserted by the cluster's authenticator, e.g. an OIDC groups claim or
a built-in group like "system:authenticated".

### spec.labels

`map<string, string>`

Additional labels to apply to every created RBAC object.
These are merged with standard Planton labels for resource tracking and governance.

### spec.annotations

`map<string, string>`

Additional annotations to apply to every created RBAC object.

## Validation Rules

- `scope_required`: Exactly one scope must be set (namespace_scope or cluster_scope)
- `role_required`: Exactly one role source must be set (create_role or existing_role)
- `existing_role_requires_subjects`: subjects must not be empty when binding to an existing role
- `aggregation_requires_cluster_scope`: aggregation_rule is only valid with cluster_scope (namespaced Roles cannot aggregate)
- `non_resource_urls_require_cluster_scope`: non_resource_urls are only valid with cluster_scope (they cannot be namespaced)
- `cluster_scope_sa_subjects_need_namespace`: service_account subjects must set namespace explicitly in cluster_scope grants

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesRbac, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.role_name` | `string` | The name of the role in the grant: the created Role/ClusterRole, or the existing role that was bound to. |
| `status.outputs.role_kind` | `string` | The Kubernetes kind of the role in the grant: "Role" or "ClusterRole". |
| `status.outputs.binding_name` | `string` | The name of the created binding. Empty when the grant defined a role with no subjects (no binding is created). |
| `status.outputs.binding_kind` | `string` | The Kubernetes kind of the created binding: "RoleBinding" or "ClusterRoleBinding". Empty when no binding is created. |
| `status.outputs.namespace` | `string` | The namespace the grant applies to. Empty for cluster-scoped grants. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespaceScope.namespace` | KubernetesNamespace | `spec.name` |
| `spec.subjects[].serviceAccount.name` | KubernetesServiceAccount | `spec.name` |
| `spec.subjects[].serviceAccount.namespace` | KubernetesNamespace | `spec.name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
