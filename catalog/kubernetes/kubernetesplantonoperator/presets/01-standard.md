# Standard

This preset installs the Planton operator in its standard posture: the
manager alone, chart defaults (one replica, leader election on, the
chart's own resource requests/limits), pinned chart version. Installing
the operator deploys NO platform — declare each platform with a
KubernetesPlantonPlatform resource, which this operator reconciles. One
installation per cluster: the operator enforces that itself at startup,
and the release name is fixed to `planton-operator`.

## When to Use

- Any cluster that will run one or more self-hosted Planton platforms
- The 30-second choice: this is the standard first Planton operator
  installation

## Key Configuration Choices

- **Manager only** — platforms are separate KubernetesPlantonPlatform
  declarations with their own lifecycles; platform teams own this
  install, application teams declare platforms
- **Chart defaults kept** — one replica (extras are leader-elected warm
  standbys, not capacity), leader election on, requests 10m/256Mi and
  limits 500m/512Mi (the real work happens in the platforms' own
  workloads, never in the operator)
- **`namespace: planton-operator` + `create_namespace: true`** — the
  convention, in a namespace this resource creates and owns
- **CRD keeps on uninstall** — the module owns the
  `plantonplatforms.planton.ai` CRD with keep-on-uninstall semantics, so
  removing the operator never cascade-deletes platform declarations

## Placeholders to Replace

None — this preset deploys as-is.

## Related Presets

- **02-ha** — two leader-elected replicas for faster operator failover
  on clusters where reconcile latency matters
