---
title: "Production sized preset"
description: "The full production posture: an HA Prometheus pair (each replica scrapes and stores the complete target set — duplication, not sharding), a quorum-safe three-replica Alertmanager gossip cluster,..."
type: "preset"
rank: "03"
presetSlug: "03-production-sized"
componentSlug: "kube-prometheus-stack"
componentTitle: "Kube Prometheus Stack"
provider: "kubernetes"
icon: "package"
order: 3
---

# Production sized preset

The full production posture: an HA Prometheus pair (each replica
scrapes and stores the complete target set — duplication, not
sharding), a quorum-safe three-replica Alertmanager gossip cluster,
30-day/80GiB dual-bounded retention on 100Gi volumes, persistent
bundled Grafana, and explicit resources on every component so the
stack cannot be evicted by the workloads it is supposed to watch.

Retention is deliberately two-dimensional: `retention` trims by age,
`retention_size` trims by bytes BELOW the volume size — when the
volume itself fills, Prometheus crash-loops instead of trimming.
Memory is the number to watch: it scales with active series, and an
undersized limit is the most common cause of a crash-looping
Prometheus on busy clusters. Treat the 2Gi/4Gi shape here as a
starting point and raise it as target count grows.

The control-plane scraper posture matches the managed-cloud preset
because most production clusters are managed — on a self-hosted
control plane (kubeadm, datacenter) DROP the `control_plane_scrapers`
and `default_rules` blocks so the controller-manager, scheduler,
etcd and kube-proxy are scraped and their curated alerts stay armed.

Change first: `external_labels.cluster` to the cluster's real name
(multi-cluster backends and federation key on it), then
`alertmanager.config_yaml` with real notification routes — until
then alerts are visible in the UIs but notify nobody. Add
`prometheus.remote_write` when a long-term or managed backend enters
the picture.

See [03-production-sized.yaml](./03-production-sized.yaml) for the
manifest.
