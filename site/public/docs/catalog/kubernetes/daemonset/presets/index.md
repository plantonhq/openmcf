---
title: "Presets"
description: "Ready-to-deploy configuration presets for DaemonSet"
type: "preset-list"
componentSlug: "daemonset"
componentTitle: "DaemonSet"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-log-collector"
    rank: "01"
    title: "Log Collector"
    excerpt: "This preset deploys a log-shipping agent (the Fluent Bit / Vector / Filebeat shape) on every node in the cluster: it reads node and pod logs from HostPath mounts, buffers in a size-limited EmptyDir,..."
  - slug: "02-node-monitor"
    rank: "02"
    title: "Node Monitor"
    excerpt: "This preset deploys a node-metrics agent (the node-exporter shape) that observes the node itself: it joins the node's network and PID namespaces, reads `/proc` and `/sys` through read-only HostPath..."
  - slug: "03-hardened-agent"
    rank: "03"
    title: "Hardened Agent"
    excerpt: "This preset deploys a per-node agent that passes the Kubernetes restricted Pod Security Standard: non-root with a pinned UID, read-only root filesystem (with a writable EmptyDir for /tmp), all Linux..."
---

# DaemonSet Presets

Ready-to-deploy configuration presets for DaemonSet. Each preset is a complete manifest you can copy, customize, and deploy.
