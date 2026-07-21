---
title: "Immutable Versioned Configuration"
description: "This preset creates an immutable ConfigMap with a versioned name (`app-config-v1`). Configuration changes ship as a NEW ConfigMap (`app-config-v2`) plus a workload update pointing at the new name —..."
type: "preset"
rank: "02"
presetSlug: "02-immutable-versioned"
componentSlug: "configmap"
componentTitle: "ConfigMap"
provider: "kubernetes"
icon: "package"
order: 2
---

# Immutable Versioned Configuration

This preset creates an immutable ConfigMap with a versioned name (`app-config-v1`). Configuration changes ship as a NEW ConfigMap (`app-config-v2`) plus a workload update pointing at the new name — never as an in-place edit.

## When to Use

- Production configuration where accidental edits must be impossible
- Roll-forward/rollback discipline: each config version is a distinct, auditable object, and rollback is just pointing the workload back at the previous name
- Large clusters where API server watch load matters — the kubelet stops watching immutable ConfigMaps

## Key Configuration Choices

- **`immutable: true`** — data cannot be updated after creation (only metadata can change); updating requires delete-and-recreate, which is exactly what the versioned-name pattern turns into a feature
- **Versioned name (`app-config-v1`)** — bump the suffix for every config change. Because the workload's reference changes too, the rollout is atomic: pods either see the old version or the new one, never a half-propagated mix
- **No propagation surprises** — mutable ConfigMaps refresh volume mounts eventually and env vars never; the versioned pattern sidesteps both, since consumers restart on the reference change

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the ConfigMap (must match the consuming workload's namespace) | Your namespace management |
| `<your-log-level>` | Application log level, e.g. `warn` for production | Your application's configuration reference |
| `<your-max-connections>` | Connection pool ceiling, e.g. `"100"` (quoted — all values are strings) | Your application's configuration reference |

Also rename `app-config-v1` (in `metadata.name` and `spec.name`) to your own `<app>-config-v<N>` scheme.

## Related Presets

- **01-app-config** — mutable configuration with a properties file
- **03-binary-payload** — base64-encoded binary entries alongside text
