# Feature Flags

Boolean flags for gating pipeline behavior — delivered as a Kubernetes
ConfigMap so a flag flip is a config apply, not a DAG redeploy.

## When to use

- Rolling out a new code path (incremental load, new export format)
  behind a flag
- Toggling expensive steps (data quality checks) per environment
- Keeping a `dry_run` switch for safe pipeline testing

## What to customize

- `environment` — reference your `GcpCloudComposerEnvironment` resource
  (its `environment_name` output).
- `data` entries — your actual flags. Keep values quoted: unquoted
  `true`/`false` (and YAML 1.1's `on`/`off`/`yes`/`no`) parse as
  booleans, and ConfigMap values must be strings.
- `configMapName` — the name DAGs reference; immutable after creation.

## Important notes

- Flag changes update in place — DAG runs started after the apply see
  the new values.
- Data is plain; anything secret belongs in
  `GcpCloudComposerUserWorkloadsSecret`.

## Composes with

`GcpCloudComposerEnvironment` upstream (reference its
`environment_name` output). Pair with `01-dag-configuration` for
non-flag tuning values.
