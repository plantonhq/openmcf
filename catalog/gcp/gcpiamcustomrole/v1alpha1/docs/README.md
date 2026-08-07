# GCP IAM Custom Roles: Least Privilege as a First-Class Node

## The Problem with Predefined Roles

Google ships more than a thousand predefined IAM roles, and almost none of them fit any real workload exactly. The pattern every security review finds is the same: a service account that needs to write objects into one class of bucket holds `roles/storage.admin` — which also lets it delete buckets, rewrite IAM on them, and read everything in the project. The workload needed two permissions; the role granted fifty.

The alternatives people reach for first are all worse:

- **Stacking narrower predefined roles** gets closer but still over-grants, and the grant list becomes an unreadable pile of role names whose union nobody can reason about.
- **Basic roles** (`roles/viewer`, `roles/editor`, `roles/owner`) are the legacy sledgehammers — Google's own documentation tells you not to use them in production.
- **Doing nothing** until the audit, which is how most projects actually behave.

A custom role is the correct tool: a named permission set with exactly the verbs the workload uses, defined once, granted everywhere it applies.

## The Role Lifecycle (and Its Sharp Edges)

Custom roles have three behaviors worth understanding before automating them:

**1. Identity is immutable; content is not.** The `roleId` (and the owning project) can never change — the full name `projects/<project>/roles/<roleId>` is the handle every grant references. But `permissions`, `title`, `description`, and `stage` all update in place, and a permission edit propagates instantly to every existing grant of the role. This is the operational win: tightening or extending access is a one-line change in one place.

**2. Deletion is soft.** Deleting a custom role does not free its ID. The role enters a deleted state for up to 14 days (visible with a `deleted: true` flag), during which grants of it are rejected but the ID remains reserved. Creating a role with the ID of a soft-deleted role *undeletes* it and patches it to the new definition. Automation that treats create-after-delete as a fresh create will converge correctly — but only if the tooling handles the undelete path, which both of this component's modules (and the underlying provider) do natively.

**3. Stages are labels, not gates.** `ALPHA`/`BETA`/`GA`/`DEPRECATED` communicate maturity to humans browsing the console; they do not alter behavior. The one exception is `DISABLED`: a disabled role stays defined but all of its bindings behave as if the role were empty — an IAM kill switch that is far safer than deleting the role (no 14-day ID reservation dance to undo).

## Project Scope vs Organization Scope

Custom roles exist at two parents: a project or an organization. This component deliberately models the **project-scoped** role:

- Project-scoped roles can only be granted on resources within that project — which is exactly the blast-radius property least-privilege design wants.
- Org-scoped roles are for platform teams standardizing a role across many projects; they demand org-level `iam.roles*` permissions that most automation should not hold.

The two are different resources with different parents and different authority requirements; conflating them into one kind with a mode switch would make the common case (project) carry the risk profile of the rare case (org).

## The 90/10 Coverage Decision

The underlying `google_project_iam_custom_role` resource has a compact surface, and this component models all of it:

| Provider field | Modeled | Notes |
|---|---|---|
| `role_id` | ✅ `roleId` | Same validation as the API (3-64 chars, no hyphens) |
| `project` | ✅ `projectId` | `StringValueOrRef` → GcpProject; falls back to the provider default project |
| `title` | ✅ `title` | Required, max 100 chars |
| `description` | ✅ `description` | Max 256 chars |
| `permissions` | ✅ `permissions` | Min 1 item; blank strings filtered defensively in the modules |
| `stage` | ✅ `stage` | Default GA |
| `deleted` (computed) | ✅ output | Surfaced in stack outputs |
| `deletion_policy` | ❌ | A Terraform-provider-level abandon-vs-delete lever, not a property of the role itself; Planton's lifecycle management owns this concern |

## Composition: the Role Is the Leaf

A custom role grants nothing by itself — it is the *definition* half of access control. The *grant* half is a separate node (GcpProjectIamMember) that binds this role to an identity. Keeping them separate is deliberate:

- One role, many grants: the same `logBucketWriter` role granted to five service accounts is five independent grant nodes referencing one role node.
- The role's `name` output (`projects/<project>/roles/<roleId>`) is exactly the string IAM grants expect — the grant's `role` field references it with no assembly.
- Permission changes live on the role; membership changes live on the grants. Two different change velocities, two different owners, two different nodes.

## Operational Guidance

- **Naming**: use camelCase for `roleId` (`logBucketWriter`, not `log-bucket-writer` — hyphens are rejected by the API). Keep the `title` human ("Log Bucket Writer") since that is what the console shows in policy pickers.
- **Discovering permissions**: `gcloud iam list-testable-permissions <full-resource-name>` lists what can be granted on a resource type; permissions marked `NOT_SUPPORTED` cannot appear in custom roles.
- **Auditing**: the role's `description` is your future auditor's first (often only) context. Write who should hold the role and why it exists.
- **Deprecation flow**: mark `stage: DEPRECATED` first (signals intent, changes nothing), move grants off it, then `DISABLED` (kill switch), then delete.
