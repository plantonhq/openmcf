# Kubernetes RBAC: Research Documentation

## Introduction

Role-Based Access Control is how Kubernetes decides who may do what. Every API request — from a human with kubectl, a controller, or a pod's ServiceAccount — is authorized against a set of roles (bundles of permission rules) reached through bindings (attachments of roles to subjects). Four object kinds implement this: Role and RoleBinding (namespaced), ClusterRole and ClusterRoleBinding (cluster-wide).

Three properties make the model predictable and are worth internalizing before writing any grant:

- **Deny by default**: a request is allowed only if some rule, reachable through a binding, allows it
- **Purely additive**: there are no deny rules; grants only ever add permissions, and multiple grants union together. Removing access means removing a grant, never adding a counter-rule
- **Evaluation order is irrelevant to outcome**: the authorizer checks ClusterRoleBindings, then RoleBindings in the request's namespace; any single matching rule allows

The friction in practice is not the model but the object topology: one logical decision — "give this team read access to this namespace" — becomes two to four coordinated objects whose names, namespaces, and `roleRef`s must line up exactly, with several cross-object rules (aggregation is ClusterRole-only, non-resource URLs are cluster-only, `roleRef` is immutable) enforced only at apply time.

Planton's **KubernetesRbac** component models the logical unit directly: one resource is one grant — *scope × role-source × subjects* — and the module materializes whichever Kubernetes objects that grant implies.

## Evolution and Historical Context

### Before RBAC: ABAC (pre-1.6)

Early Kubernetes shipped Attribute-Based Access Control: a JSON policy file on the API server, requiring an API server restart per policy change. Unmanageable at any scale.

### RBAC general availability (1.8, 2017)

RBAC reached GA in Kubernetes 1.8 and became the default authorizer everywhere. The four-object model (Role/ClusterRole × Binding/ClusterRoleBinding) has been essentially frozen since — a testament to the model, and the reason tooling above it is where improvement happens.

### Aggregated ClusterRoles (1.9+)

Kubernetes 1.9 added the `aggregationRule`: a ClusterRole whose rules are continuously composed by the RBAC controller from every ClusterRole matching its label selectors. This solved the CRD-extension problem — when an operator installs new resource types, it ships small ClusterRoles labeled `rbac.authorization.k8s.io/aggregate-to-view: "true"` (and `-to-edit`, `-to-admin`), and the built-in roles absorb the new permissions with no cluster-admin action. The same mechanism is available for building custom plugin-style permission sets.

### The built-in user-facing roles

Every conformant cluster ships four aggregated ClusterRoles intended for humans:

| Role | Grants | Notes |
|------|--------|-------|
| `view` | Read-only on most namespaced objects | Excludes Secrets (reading a Secret is effectively holding the credential) and RBAC objects |
| `edit` | Read/write on most namespaced objects | Excludes RBAC — an editor cannot expand their own access |
| `admin` | `edit` plus namespace-local RBAC management | Cannot modify the namespace itself or resource quotas |
| `cluster-admin` | `*` on `*`, everything everywhere | Bind sparingly, and never to `system:authenticated` |

Granting a built-in per-namespace is done with a RoleBinding whose `roleRef` points at the ClusterRole — the namespaced binding confines the cluster-wide role definition to one namespace. This indirection (binding kind ≠ role kind) is among the most common points of confusion in raw RBAC, and one the grant model absorbs.

### Escalation prevention

Since RBAC's GA, the API server enforces that you can only grant permissions you already hold (or hold the `escalate` verb for). This matters operationally: the credential running IaC must itself hold every permission its grants confer.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
kubectl create role app-reader -n production \
  --verb=get,list,watch --resource=pods,services,configmaps

kubectl create rolebinding app-reader-binding -n production \
  --role=app-reader --serviceaccount=production:app-identity

# Built-in role for a group
kubectl create rolebinding team-view -n production \
  --clusterrole=view --group=dev-team
```

**Pros:**
- The subcommands encode the role-vs-clusterrole distinction correctly

**Cons:**
- Two imperative commands per grant, no record, no reproducibility
- `kubectl edit` on live RBAC objects is how permission drift is born

**Verdict:** Debugging and break-glass only.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: app-reader
  namespace: production
rules:
  - apiGroups: [""]
    resources: ["pods", "services", "configmaps"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: app-reader-binding
  namespace: production
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: app-reader
subjects:
  - kind: ServiceAccount
    name: app-identity
    namespace: production
```

**Pros:**
- Declarative, version-controllable, complete surface

**Cons:**
- Every grant is 2+ documents whose names and kinds must line up by hand; `roleRef` typos fail at apply
- Cross-object rules (aggregation scope, non-resource URL scope, `roleRef` immutability) discovered at apply time
- Subject strings (`system:serviceaccount:...`) assembled by hand

**Verdict:** The baseline; correctness rests on conventions and review.

### Level 2: Terraform

```hcl
resource "kubernetes_role_v1" "app_reader" {
  metadata {
    name      = "app-reader"
    namespace = "production"
  }
  rule {
    api_groups = [""]
    resources  = ["pods", "services", "configmaps"]
    verbs      = ["get", "list", "watch"]
  }
}

resource "kubernetes_role_binding_v1" "app_reader" {
  metadata {
    name      = "app-reader-binding"
    namespace = "production"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "Role"
    name      = kubernetes_role_v1.app_reader.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = "app-identity"
    namespace = "production"
  }
}
```

**Pros:**
- Full IaC lifecycle; role/binding name coupling via references

**Cons:**
- Still two resources per grant, four resource types to choose among
- `role_ref.kind` is a plain string; scope rules unvalidated until apply
- `roleRef` immutability surfaces as a confusing apply-time replace

**Verdict:** Production-grade lifecycle; the four-object topology remains the author's problem.

### Level 3: Pulumi

```go
role, _ := rbacv1.NewRole(ctx, "app-reader", &rbacv1.RoleArgs{ /* ... */ })
_, _ = rbacv1.NewRoleBinding(ctx, "app-reader-binding", &rbacv1.RoleBindingArgs{
    RoleRef: rbacv1.RoleRefArgs{
        ApiGroup: pulumi.String("rbac.authorization.k8s.io"),
        Kind:     pulumi.String("Role"),
        Name:     role.Metadata.Name().Elem(),
    },
    Subjects: rbacv1.SubjectArray{ /* ... */ },
})
```

**Pros:**
- Full programming language; names flow between resources

**Cons:**
- Same topology burden and stringly-typed kinds as Terraform

**Verdict:** Excellent IaC; no abstraction over the grant.

### Other Methods

**Helm:** charts template Role+RoleBinding pairs behind `rbac.create: true` values — fine inside a chart, not a general grant-management approach.

**Specialized tools** (audit2rbac, rbac-manager, rbac-lookup): valuable for auditing and generating from observed traffic; complementary to, not a replacement for, declarative grant management.

## Comparative Analysis

| Aspect | kubectl | YAML | Terraform | Pulumi | Planton |
|--------|---------|------|-----------|--------|---------|
| Objects per grant | 2 commands | 2–4 documents | 2 resources | 2 resources | 1 resource |
| Scope/kind selection | Subcommand | By hand (4 kinds) | By hand + strings | By hand + strings | One `scope` choice |
| Cross-object rules | At apply | At apply | At apply | At apply | Schema + CEL |
| Built-in role grants | `--clusterrole` flag | roleRef by hand | role_ref by hand | RoleRef by hand | `existingRole` + `isClusterRole` |
| SA subjects as references | No | No | Within-state only | Within-state only | Cross-resource refs |
| Dual IaC | N/A | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### One grant, three orthogonal choices

The spec is shaped after the decision, not the objects:

1. **Scope** (`oneof`): `namespaceScope` (with a namespace value-or-ref) or `clusterScope`. This one choice determines whether the module creates Role/RoleBinding or ClusterRole/ClusterRoleBinding — the four-kind selection problem disappears
2. **Role source** (`oneof`): `createRole` (rules and/or an aggregation rule, optional explicit name) or `existingRole` (name + `isClusterRole` flag, covering both "bind built-in `view` in a namespace" and "bind a custom Role")
3. **Subjects** (repeated): ServiceAccounts (name and namespace as value-or-refs), users, or groups — each subject exactly one of the three, matching rbac/v1 Subject semantics

Every combination is materialized correctly: role-only (publish a reusable role, no binding), binding-only (existing role + subjects), or both.

### Kubernetes' cross-object rules as schema rules

The rules Kubernetes enforces at apply time — or worse, silently — are CEL validations on the spec:

- Exactly one scope; exactly one role source
- `existingRole` requires subjects (binding a role to nobody deploys nothing)
- `aggregationRule` requires cluster scope (namespaced Roles cannot aggregate)
- `nonResourceUrls` require cluster scope (those paths exist outside any namespace)
- ServiceAccount subjects in cluster scope must set a namespace explicitly (there is nothing to default from)
- Per-rule: resources XOR non-resource URLs; every rule must grant something; a role must have rules or aggregation

A grant that validates will apply; the class of "the YAML was accepted but the objects don't do what was meant" errors shrinks accordingly.

### PolicyRule fidelity

`KubernetesRbacPolicyRule` mirrors rbac/v1 PolicyRule exactly — `verbs`, `apiGroups`, `resources` (subresources as `pods/log`), `resourceNames`, `nonResourceUrls` — so there is no translation layer to learn and no upstream capability lost. The documentation semantics (e.g. `resourceNames` cannot constrain `create` because authorization happens before a name exists) are stated on the fields themselves.

### Composition with KubernetesServiceAccount

ServiceAccount subject names are `StringValueOrRef`s targeting KubernetesServiceAccount resources. An infra chart creates the namespace, the identity, and the grant in one run, with names flowing through the resource graph instead of being retyped.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates resource creation — Kubernetes provider, then role and/or binding per the grant shape
- **`locals.go`**: Resolves the scope namespace and subject references; computes role/binding names and kinds
- **`outputs.go`**: Exports `role_name`, `role_kind`, `binding_name`, `binding_kind`, `namespace`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes which of the four object kinds the grant implies
- **`main.tf`**: Conditionally creates `kubernetes_role_v1` / `kubernetes_cluster_role_v1` and `kubernetes_role_binding_v1` / `kubernetes_cluster_role_binding_v1`
- **`outputs.tf`**: Exports role and binding names and kinds

### Resource Count

The component creates **zero, one, or two Kubernetes objects** depending on the grant shape: a role (when `createRole`), a binding (when subjects are present) — never more. The complexity is in the shape selection and validation, not in orchestration.

## Production Best Practices

### Rule discipline

1. **Least verbs**: `get,list,watch` for readers; add write verbs individually. `*` belongs only in true admin roles
2. **Narrow resources**: list resources explicitly; wildcard resources silently absorb future CRDs into old grants
3. **Use `resourceNames` for point access**: e.g. `get` on one named Secret — and remember it cannot constrain `create` or `deletecollection`
4. **Treat Secrets as privileged**: even built-in `view` excludes them; `get` on `secrets` is effectively holding every credential in the namespace

### Grant topology

1. **Prefer namespace scope**: reach for cluster scope only for cluster-scoped resources, non-resource URLs, or genuinely cluster-wide permissions
2. **Prefer built-ins for humans**: `view`/`edit`/`admin` per-namespace to groups covers most team access with zero maintained rules
3. **Bind groups, not users**: membership changes then live in the identity provider
4. **One grant per intent**: many small grants audit and revoke cleanly; mega-roles do neither

### Operational

1. **Mind escalation prevention**: the IaC credential must hold every permission it grants
2. **`roleRef` is immutable**: changing a grant's role source means the binding is replaced — momentary revocation during the swap; plan changes to hot paths accordingly
3. **Audit periodically**: `kubectl auth can-i --list --as=system:serviceaccount:<ns>:<name>` answers "what can this identity actually do" — the union of all grants, which no single resource shows
4. **No deny rules exist**: "must not" requirements are admission policy (ValidatingAdmissionPolicy, OPA/Kyverno), not RBAC

## Conclusion

Kubernetes RBAC's model — additive, deny-by-default, four object kinds — has been stable for years; the recurring failure mode is coordination between the objects, not the model. KubernetesRbac collapses the coordination: one resource states the grant (scope × role-source × subjects), the schema enforces the cross-object rules Kubernetes would only reveal at apply time, and the module materializes exactly the objects the grant implies. Built-in roles, aggregation, non-resource URLs, and role-only publication are all first-class shapes of the same resource.

## References

- [Using RBAC Authorization](https://kubernetes.io/docs/reference/access-authn-authz/rbac/)
- [User-Facing Roles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#user-facing-roles)
- [Aggregated ClusterRoles](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#aggregated-clusterroles)
- [Privilege Escalation Prevention](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#privilege-escalation-prevention-and-bootstrapping)
- [RBAC Good Practices](https://kubernetes.io/docs/concepts/security/rbac-good-practices/)
- [Checking API Access (kubectl auth can-i)](https://kubernetes.io/docs/reference/access-authn-authz/authorization/#checking-api-access)
- [Pulumi Kubernetes RBAC](https://www.pulumi.com/registry/packages/kubernetes/api-docs/rbac/v1/)
- [Terraform kubernetes_role_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/role_v1)
