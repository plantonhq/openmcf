---
title: "Jaeger-compatible Tempo"
description: "Tempo with the four legacy Jaeger receiver protocols opened alongside OTLP and the Jaeger-UI-compatible query sidecar enabled — a migration posture for fleets still emitting Jaeger while they move to..."
type: "preset"
rank: "03"
presetSlug: "03-jaeger-compat"
componentSlug: "grafana-tempo"
componentTitle: "Grafana Tempo"
provider: "kubernetes"
icon: "package"
order: 3
---

# Jaeger-compatible Tempo

Tempo with the four legacy Jaeger receiver protocols opened alongside OTLP
and the Jaeger-UI-compatible query sidecar enabled — a migration posture
for fleets still emitting Jaeger while they move to OTLP, and for tooling
that speaks the Jaeger query API.

**When to use:** you have existing Jaeger instrumentation or Jaeger-native
tooling and want Tempo as the backend without rewriting emitters first.

**Direction of travel:** OTLP is the 2026 standard; treat the Jaeger
receivers as a bridge and disable them once emitters are on OTLP.
