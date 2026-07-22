---
title: "Production Observability (Hubble + Prometheus + WireGuard)"
description: "This preset is the production hardening layer for a Cilium installation: the full Hubble stack (relay, UI, and the core metric families), Prometheus telemetry with ServiceMonitors, transparent..."
type: "preset"
rank: "04"
presetSlug: "04-production-observability"
componentSlug: "cilium"
componentTitle: "Cilium"
provider: "kubernetes"
icon: "package"
order: 4
---

# Production Observability (Hubble + Prometheus + WireGuard)

This preset is the production hardening layer for a Cilium installation:
the full Hubble stack (relay, UI, and the core metric families), Prometheus
telemetry with ServiceMonitors, transparent WireGuard encryption of
pod-to-pod traffic, the eBPF bandwidth manager, HA operator replicas, and
deliberate resource sizing for both the agent DaemonSet and the operator.
Combine its fields with the cluster-posture preset that matches your
environment (primary CNI or chaining) — this preset intentionally leaves
IPAM/routing at chart defaults.

## When to Use

- Production clusters that need flow-level visibility: who talked to whom,
  what was dropped and why, DNS/HTTP behavior per service
- Clusters running the Prometheus operator (kube-prometheus-stack) that
  should scrape Cilium and Hubble automatically
- Compliance postures requiring encryption in transit between pods without
  per-application TLS work

## Key Configuration Choices

- **Hubble `relay` + `ui`** — relay aggregates flows cluster-wide; the UI
  serves the live service map. The spec's CEL rule makes the dependency
  explicit: the UI reads flows exclusively through relay.
- **`hubble.metrics: [dns, drop, tcp, flow, http]`** — the core metric
  families for production triage (drop reasons, DNS failures, HTTP rates).
  `metricsServiceMonitor` stays false here — flip it to true to scrape
  Hubble metrics through the same Prometheus operator.
- **`prometheus.enabled` + `serviceMonitor: true`** — agent and operator
  `/metrics` with automatic scrape discovery. The Prometheus operator CRDs
  MUST exist on the cluster or the Helm release fails to install.
- **WireGuard encryption** — transparent pod-to-pod encryption with
  automatic key management; no application changes, no certificate estate.
- **`bandwidthManager.enabled: true`** — pods' `egress-bandwidth`
  annotations enforced in eBPF instead of noisy token-bucket qdiscs.
- **`policyEnforcementMode: default`** — enforcement only where policies
  select pods; move to `always` for default-deny postures.
- **`operator.replicas: 2` + resources** — the chart's HA default made
  explicit (requires 2+ nodes: the replicas carry pod anti-affinity), with
  requests/limits set deliberately instead of the chart's unset defaults.
  Memory-only limits avoid CPU throttling of the dataplane components.

## Placeholders to Replace

No placeholders — this preset is directly deployable. Review the resource
requests/limits against your node sizes and cluster scale.

## Prerequisites

- Prometheus operator CRDs on the cluster (`serviceMonitor: true` makes
  the release fail without them)
- 2+ schedulable nodes for the two operator replicas
- Kernels with WireGuard support (any modern distribution kernel)

## Related Presets

- **01-kind-dev-cluster** — the minimal local posture
- **02-eks-chaining** — the EKS chaining posture to combine with
- **03-self-managed-primary-kpr** — the primary-CNI posture to combine with
