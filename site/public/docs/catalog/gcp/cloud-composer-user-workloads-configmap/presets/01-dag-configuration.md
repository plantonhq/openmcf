---
title: "DAG Configuration"
description: "Runtime configuration for a pipeline — endpoints, batch sizing, retry tuning — delivered as a Kubernetes ConfigMap into a Composer environment. The DAG reads values by key instead of hard-coding them."
type: "preset"
rank: "01"
presetSlug: "01-dag-configuration"
componentSlug: "cloud-composer-user-workloads-configmap"
componentTitle: "Cloud Composer User Workloads ConfigMap"
provider: "gcp"
icon: "package"
order: 1
---

# DAG Configuration

Runtime configuration for a pipeline — endpoints, batch sizing, retry
tuning — delivered as a Kubernetes ConfigMap into a Composer
environment. The DAG reads values by key instead of hard-coding them.

## When to use

Any DAG whose behavior should be tunable without editing and
redeploying DAG code: target endpoints, batch sizes, timeouts,
timezones.

## What to customize

- `environment` — reference your `GcpCloudComposerEnvironment` resource
  (its `environment_name` output).
- `data` entries — your pipeline's actual knobs; all values are
  strings, so quote numbers and booleans in YAML.
- `configMapName` — the name DAGs reference; immutable after creation.

## Important notes

- Data is plain (readable by anyone with cluster access) — use
  `GcpCloudComposerUserWorkloadsSecret` for credentials.
- Data updates in place; a config change is an apply, not a recreate.

## Composes with

`GcpCloudComposerEnvironment` upstream (reference its
`environment_name` output). Pair with a user workloads Secret for the
credential half of the same pipeline's configuration.
