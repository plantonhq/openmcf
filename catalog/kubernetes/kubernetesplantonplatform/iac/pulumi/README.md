# KubernetesPlantonPlatform Pulumi Module

Declares one `PlantonPlatform` custom resource that the Planton operator
reconciles into a running self-hosted platform, plus the optional owning
namespace. The CR is rendered UNTYPED (`apiextensions.NewCustomResource`
with a plain nested map built in `platform_cr.go`) — the exact semantic
twin of the Terraform module's `kubectl_manifest` + `yamlencode` locals.

## Module Behavior

- **Keys render only when declared** — the operator's own defaulting
  stays authoritative for everything unset, the same posture as the
  `planton` umbrella chart's verbatim pass-through. Three-state
  optionals (the default-true toggles, the defaulted scalars) render
  exactly when the proto field is PRESENT: an explicit `enabled: true`
  is faithfully forwarded even though it matches the CRD default,
  because presence is the user's deliberate statement.
- **The CR is namespaced and named from `metadata.name`** — the prefix
  of every object the operator creates for the platform, which is why
  the outputs (gateway Service, setup-code Secret, the two exact
  commands) derive deterministically at declaration time.
- **Destroy is garbage collection, not operator work** — every
  operator-created object is owner-referenced to the CR (the operator
  has no finalizers), so deletion completes even when the operator is
  already gone. The 15-minute delete timeout is headroom, never an
  expected wait.
- **`version` is always explicit** — required by the spec (mirroring the
  CRD), never defaulted by the module: a module default would turn a
  catalog update into a silent whole-platform upgrade on the next apply.

## Why Untyped

The PlantonPlatform schema is consumed as data rather than a generated
SDK type: the operator's CRD is the validation authority (server-side
apply surfaces its rejections verbatim), and the render-only-declared
posture would be fought by any typed struct's zero values.
