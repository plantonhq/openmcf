# TestCloudResourceGeneric — internal documentation

This kind is the catalog's standing canary and the certification suite's
primary vehicle. It is registered like any real kind (enum entry, id prefix
`tcrg`, grammar-validated version) and carries the full per-kind file shape —
protos, dual-engine IaC modules, presets — so every pipeline that walks the
catalog walks this kind too. When a repo-wide transformation or a machinery
change breaks a pipeline, it breaks here first, in a kind no user depends on.

## What it deliberately exercises

- **Every generic field class** in `spec.proto`: scalar defaults, absent
  defaults, nested messages, maps, repeated fields, sensitive fields, and
  value-or-reference strings including a fully annotated foreign-key fixture
  (`annotated_ref`, self-referential to this kind's own outputs).
- **The envelope contract**: `api.proto` carries the same const-validated
  `apiVersion`/`kind` pair as every real kind.
- **The IaC lifecycle without a cloud**: `stack_input.proto` has no provider
  configuration; both engines produce deterministic outputs from inputs.

## What stays out of user surfaces

Release content zips, the catalog site, and module auto-tagging exclude the
`_test` provider explicitly (guarded in CI). This kind must never appear in
anything a user or their agent browses.
