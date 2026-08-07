# KubernetesGrafana Guide

The judgment this guide carries: standalone Grafana earns its node when
it reads MORE than one source — and its default state is a trap: every
hand-made dashboard lives in an ephemeral local database that vanishes on
pod restart unless persistence is declared.

## Standalone hub vs the stack's bundled Grafana

One kube-prometheus-stack and nothing else to look at? Its bundled
Grafana (on by default there) is the simpler path — skip this kind. The
moment dashboards must read Loki, Tempo, ClickHouse, Postgres, or a
second Prometheus, THIS kind is the composition hub: declare one
standalone Grafana with a `datasources` entry per source, each `url`
wired by `valueFrom` to the source's exported endpoint — the full wired
example lives in the
[observability-stack pattern](../../_patterns/observability-stack.md).
Never run both for the same audience.

## Declare state, or lose it

UI-authored dashboards, users, and preferences live in an embedded
SQLite on local disk, and the chart's default is EPHEMERAL — a pod
restart erases everything hand-made (the reference page states it).
Either declare persistence/database (the spec's own arms), or treat
dashboards as code: `dashboards` provisioning renders them from the
manifest, present from first boot, immune to restarts. For anything
beyond a scratch environment, one of the two is part of the proposal.

## Credentials

The chart generates the admin password once, into the `<name>` Secret —
consume by reference; or point `adminSecret` at a Secret you manage
(e.g. a KubernetesExternalSecret projection). Never inline.

## On the diagram

Grafana renders as the hub with a labeled datasource edge into every
source — the observability topology is readable from the graph alone.
Provisioned dashboards live inside the node; hand-made ones are invisible
AND ephemeral, a double reason to provision.

## Pairs well with

- KubernetesKubePrometheusStack / KubernetesLoki / KubernetesTempo — the
  three standard datasources (pattern above).
- KubernetesExternalSecret — managed admin credentials.
- KubernetesIngress / route kinds — exposing the UI, composed as always.
