# GcpProject: Design Notes

## What This Component Models

One `google_project` — GCP's Layer-0 container. The component is the root
of the GCP composition graph: every other GCP kind's `projectId` is a
reference that can resolve this kind's `project_id` output, so environment
charts create the project first and everything else composes into it.

The spec owns exactly what the project resource owns: identity, hierarchy
placement, billing linkage, labels/tags, the default-network decision, API
enablement, and the deletion policy.

## Why IAM Is Not Bundled

A project-owner grant inside the project resource is a bundled black box:
it duplicates the first-class `GcpProjectIamMember` kind, hides an IAM
grant where audits do not look for it, and models only one hardcoded role.
Access control composes as one additive `GcpProjectIamMember` per grant —
each grant independently owned, condition-capable, and removable without
touching the project.

## Identity Is Deterministic

`projectId` is the project's permanent, globally-unique identity, chosen
by the user and immutable — the same manifest produces the same project on
either engine, byte for byte. Uniqueness belongs in the ID itself (a team
or sequence suffix chosen by the user), never in generated randomness: a
nondeterministically-suffixed ID would make the same manifest converge to
different projects depending on which engine deployed it, on the one kind
whose output feeds every other kind's project reference.

## Deletion Policy Semantics

The provider's real switch is three-way, and the spec models it directly:

- `DELETE` (default) — destroy shuts the project down; GCP holds it in a
  30-day pending-deletion window during which it can be restored.
- `PREVENT` — destroy fails. The setting for foundation projects whose
  loss would cascade.
- `ABANDON` — the resource leaves state and the project lives on
  unmanaged. The safe hand-off when ownership moves to another team or
  tool. A boolean "protection" flag would make this path inexpressible.

The default is passed explicitly on both engines so destroy semantics
never depend on a provider's client-side default.

## The Default Network Decision

GCP auto-creates a "default" VPC in new projects unless told otherwise.
The spec's `autoCreateNetwork` defaults to false — the auto-created
network is an implicit, permissive resource that explicit `GcpVpcNetwork`
components replace — with the provider's documented caveat taught in the
field comment: one network slot of quota is still consumed momentarily
during creation even when false.

## Tags vs Labels

- **Labels** are mutable key/value metadata and the primary
  cost-allocation dimension in billing exports. User labels merge beneath
  the platform's attribution labels, identically on both engines.
- **Resource-manager tags** (`tagKeys/{id}` → `tagValues/{id}`) drive org
  policies and IAM conditions. On this resource they bind at CREATE TIME
  only — changing them later recreates the project — so post-creation tag
  changes belong to out-of-band tag bindings, not this spec.

## Scope Boundaries

- Per-grant IAM: `GcpProjectIamMember`.
- Org/folder management: the organization and folder nodes themselves are
  not modeled by this kind; `parentId` takes their numeric IDs.
