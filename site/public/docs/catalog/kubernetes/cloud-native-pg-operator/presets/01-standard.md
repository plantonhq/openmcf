---
title: "Standard"
description: "This preset installs the CloudNativePG operator in its standard posture: the operator release alone, cluster-wide watch scope, pinned chart version, sized and prioritized for a production control..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "cloud-native-pg-operator"
componentTitle: "Cloud Native PG Operator"
provider: "kubernetes"
icon: "package"
order: 1
---

# Standard

This preset installs the CloudNativePG operator in its standard posture:
the operator release alone, cluster-wide watch scope, pinned chart
version, sized and prioritized for a production control plane. No backup
plugin — databases declared with KubernetesPostgres run, replicate, and
fail over, but their backup blocks need the plugin arm (see
02-with-backup-plugin). One installation per cluster: the CRDs and
webhooks are cluster-scoped singletons, and the release name is fixed to
`cnpg`.

## When to Use

- Any cluster that will run KubernetesPostgres databases and does not
  need object-store backups yet
- The 30-second choice: this is the standard first CloudNativePG
  installation

## Key Configuration Choices

- **Operator only, no plugin** — the Barman Cloud plugin (and its
  cert-manager dependency) stays out until a database needs backups;
  enabling it later is a spec change, not a reinstall
- **`replicas: 1`** (chart default) — extra replicas are leader-elected
  warm standbys that shorten failover of the operator itself; they add
  no reconciliation throughput (`max_concurrent_reconciles` is that
  knob)
- **Explicit resources** — the chart ships no requests/limits by
  default; a control-plane component should have both
- **`priority_class_name: system-cluster-critical`** — databases stop
  failing over without their operator; it should outlive stateless
  workloads under node pressure
- **`namespace: cnpg-system` + `create_namespace: true`** — the upstream
  convention, in a namespace this resource creates and owns
- **CRDs keep-on-uninstall** — the chart stamps
  `helm.sh/resource-policy: keep` on every CRD unconditionally, so
  removing the component never cascade-deletes the databases

## Placeholders to Replace

None — this preset deploys as-is.

## Related Presets

- **02-with-backup-plugin** — the same installation plus the Barman
  Cloud plugin, for clusters whose databases declare backup blocks
