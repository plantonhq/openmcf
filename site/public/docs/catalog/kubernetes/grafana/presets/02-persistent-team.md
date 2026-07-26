---
title: "Persistent team preset"
description: "The single stateful instance most teams actually want: a 10Gi volume under Grafana's embedded database so hand-built dashboards, users and preferences survive pod restarts, a Prometheus datasource..."
type: "preset"
rank: "02"
presetSlug: "02-persistent-team"
componentSlug: "grafana"
componentTitle: "Grafana"
provider: "kubernetes"
icon: "package"
order: 2
---

# Persistent team preset

The single stateful instance most teams actually want: a 10Gi volume
under Grafana's embedded database so hand-built dashboards, users and
preferences survive pod restarts, a Prometheus datasource wired by
reference to the cluster's metrics stack, sized resources, and the
public root URL set for composed exposure.

The volume is ReadWriteOnce, which makes this a one-replica shape by
design — the spec enforces it. That is not a limitation to work
around with a bigger number: two Grafanas sharing one SQLite file
corrupt it, and two with separate files silently split the team's
dashboards. When this instance becomes load-bearing enough to need
HA, the move is the 03 preset — state into an external database,
`storage` removed, replicas raised — and the migration is a Grafana
database export/import, so doing it before the dashboard count grows
is cheaper than after.

The datasource uses the reference form: it resolves to the named
KubernetesKubePrometheusStack's exported Prometheus endpoint and
gives the deployment a real dependency edge — the stack deploys
first, and renaming it updates this Grafana instead of leaving a
dead literal URL behind.

Change first: replace the `root_url` placeholder with the real
hostname and compose the ingress or gateway route over the exported
`service` handle; then let other teams ship dashboards through
ConfigMaps labeled `grafana_dashboard: "1"` — the sidecar (on by
default) discovers them cluster-wide with no edits here.

See [02-persistent-team.yaml](./02-persistent-team.yaml) for the
manifest.
