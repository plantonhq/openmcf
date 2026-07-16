# GCP Project IAM Grants: Why Additive Is the Only Composable Mode

## Three Write Modes, One Safe Choice

GCP's project IAM policy can be written three ways, and the difference is not cosmetic — it decides whether independent tools can coexist on one project:

1. **Additive member** — "ensure this one (role, member) pair exists." Creating it merges the pair into the policy; destroying it subtracts exactly that pair. Every other binding is untouched.
2. **Authoritative binding** — "this role's member list is exactly this." Anyone else's grant on the same role is silently removed on the next apply.
3. **Authoritative policy** — "the project's entire IAM policy is exactly this." The nuclear option: one apply can lock every human and service out of a project.

In a platform where multiple charts, teams, and tools each own a slice of a project's access, modes 2 and 3 are fight-by-design: two owners of the same role clobber each other on every reconcile, and the loser finds out via an outage. Planton therefore models **only the additive member**. This is a deliberate stopping line, not a coverage gap — the authoritative modes are omitted because their semantics are hostile to composition, which is the platform's core promise.

## Anatomy of a Grant

A grant is a triple — project, role, member — optionally narrowed by a condition:

- **The role half** can be a predefined role (`roles/logging.logWriter`) or a custom role's fully-qualified name (`projects/<project>/roles/<roleId>`). Custom roles are themselves first-class nodes (GcpIamCustomRole), and their `name` output is exactly the string a grant expects.
- **The member half** is an identity in IAM member syntax. In infrastructure code the overwhelmingly common member is a service account, and GcpServiceAccount's `member` output emits the ready-made `serviceAccount:<email>` string so no consumer ever assembles the prefix by hand.
- **The condition**, when present, is a CEL expression that must evaluate true for the grant to apply (expiry dates, resource-name prefixes, request attributes). Critically, the condition is part of the grant's *identity*: the same role granted with and without a condition are two independent policy entries. This is why conditions are immutable here — "editing" a condition is really replacing one grant with a different one, and the API models it the same way.

## Why Every Field Is Immutable

The IAM API has no "update grant" operation. What looks like an edit is always a remove-and-add on the policy, and any tooling that pretends otherwise invents an update semantic the platform beneath it does not have. This component mirrors reality: all fields are create-time, and any change replaces the grant atomically. The practical consequence is pleasant — a grant is cheap, fast, and stateless to replace, and there is never a partially-updated grant.

## Grants as Graph Edges

Modeled as its own node, a grant makes the access relationship *visible* in the resource graph:

```
GcpServiceAccount ──member──▶ GcpProjectIamMember ◀──role── GcpIamCustomRole
                                     │
                                  project
                                     ▼
                                 GcpProject
```

Compare this with role lists embedded inside an identity (which GcpServiceAccount also supports for convenience): embedded lists deploy in one shot but the graph cannot show *what* was granted where, grants cannot be added or removed independently of the identity's lifecycle, and a role granted to a non-service-account member has no home at all. The standalone grant node covers all of it: user and group grants, federated principals, conditional grants, custom roles, and cross-chart composition where the identity, the role, and the grant are owned by different manifests.

## The 90/10 Coverage Decision

| Provider surface | Modeled | Notes |
|---|---|---|
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; falls back to the provider's default project (made concrete via the provider client config, since the underlying resource requires an explicit project) |
| `role` | ✅ `role` | `StringValueOrRef` → GcpIamCustomRole `name` output; literal for predefined roles |
| `member` | ✅ `member` | `StringValueOrRef` → GcpServiceAccount `member` output; literal for all other principal forms; format validated at deploy time (the value usually arrives via a reference resolved only then) |
| `condition` | ✅ `condition` | title + expression required, description optional |
| `google_project_iam_binding` / `_policy` | ❌ | Deliberately excluded — authoritative clobber semantics (see above) |
| `google_project_iam_audit_config` | ❌ | Audit-log configuration is a different concern from access grants; revisit on concrete pull |

## Operational Guidance

- **Public grants**: `allUsers` and `allAuthenticatedUsers` pass validation because they are legitimate (public websites, public datasets) — but they make the resource world-readable. Most organizations should pair them with organization policy constraints.
- **Case sensitivity**: IAM treats most member emails case-insensitively but federated `principal://` identifiers case-sensitively; the modules pass values through untouched, so keep casing consistent with the identity provider.
- **Etag output**: the exported policy etag is a fingerprint of the policy version after the grant — useful when correlating a grant with audit-log entries about the policy change.
