---
title: "Log Collector"
description: "This preset deploys a log-shipping agent (the Fluent Bit / Vector / Filebeat shape) on every node in the cluster: it reads node and pod logs from HostPath mounts, buffers in a size-limited EmptyDir,..."
type: "preset"
rank: "01"
presetSlug: "01-log-collector"
componentSlug: "daemonset"
componentTitle: "DaemonSet"
provider: "kubernetes"
icon: "package"
order: 1
---

# Log Collector

This preset deploys a log-shipping agent (the Fluent Bit / Vector / Filebeat shape) on every node in the cluster: it reads node and pod logs from HostPath mounts, buffers in a size-limited EmptyDir, and tolerates the control-plane taint so control-plane logs are collected too.

A DaemonSet has no replica count and no Service — node membership IS the replica count, and log collectors initiate their own connections outward.

## When to Use

- Cluster-wide log shipping to a central sink (Elasticsearch, Loki, S3, a SaaS backend)
- Any agent that must read files from every node's filesystem

## Key Configuration Choices

- **Read-only HostPath mounts of `/var/log` and `/var/log/pods`** — the standard sources for node and container logs; `readOnly: true` because a collector should never be able to modify node logs
- **`type: Directory` vs `DirectoryOrCreate`** — `/var/log` must already exist on any node (fail loudly if not); `/var/log/pods` is created by the kubelet, so `DirectoryOrCreate` avoids a race on fresh nodes
- **Control-plane toleration with `operator: Exists`** — control-plane nodes carry a `NoSchedule` taint; without this toleration their logs are silently missing from the pipeline
- **Size-limited EmptyDir buffer** — spill space when the sink is slow or unreachable; the limit stops a broken sink from filling the node disk
- **Explicit resource limits** — this pod runs on every node, so its footprint multiplies by the node count; 500m/512Mi caps a misbehaving collector before it competes with workloads

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-log-collector-image>` | Log collector image (e.g., `fluent/fluent-bit`) | Your container registry |
| `<your-image-tag>` | Image tag or version | Your registry or CI/CD pipeline output |

## Related Presets

- **02-node-monitor** — Host-namespace metrics agent with a host port
- **03-hardened-agent** — Restricted-profile agent for security-sensitive clusters
