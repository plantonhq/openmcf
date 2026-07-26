---
title: "Dev single cluster preset"
description: "The smallest honest kube-prometheus-stack for a kind cluster, a laptop lab, or a private single-node environment: small Prometheus and Alertmanager volumes, short retention, and the control-plane..."
type: "preset"
rank: "01"
presetSlug: "01-dev-single-cluster"
componentSlug: "kube-prometheus-stack"
componentTitle: "Kube Prometheus Stack"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev single cluster preset

The smallest honest kube-prometheus-stack for a kind cluster, a laptop
lab, or a private single-node environment: small Prometheus and
Alertmanager volumes, short retention, and the control-plane scrapers
that cannot succeed on those platforms turned off (with their matching
rule groups disabled so the target list and alert set stay truthful).

Everything else is the component default — cluster-wide ServiceMonitor
discovery (so every catalog component's `service_monitor_enabled`
toggle lights up without extra wiring), Alertmanager on, the bundled
Grafana on with a chart-generated admin Secret, and both exporters.
Resource names stay predictable because the modules pin the chart
fullname to `metadata.name`; keep that name at 26 characters or fewer
(the chart silently truncates beyond that).

Port-forward Prometheus or the bundled Grafana from the stack outputs
when you want a browser; nothing is exposed outside the cluster.
Grow first by raising `prometheus.disk_size` / `retention`, then by
switching to the managed-cloud or production-sized presets when the
cluster shape changes.

See [01-dev-single-cluster.yaml](./01-dev-single-cluster.yaml) for the
manifest.
