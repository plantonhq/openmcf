---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenTelemetry Operator"
type: "preset-list"
componentSlug: "opentelemetry-operator"
componentTitle: "OpenTelemetry Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-default"
    rank: "01"
    title: "Default preset"
    excerpt: "The standard operator install: the pinned chart (0.120.0 = operator v0.156.0) into its own `otel-operator` namespace, watching every namespace on the cluster. Nothing else is set — the chart creates..."
  - slug: "02-private-mirror"
    rank: "02"
    title: "Private-mirror preset"
    excerpt: "The air-gapped posture: both image seams mirrored, because they are DIFFERENT seams. `imageRegistry` replaces only the registry part of the operator's own manager image (the path stays the upstream..."
---

# OpenTelemetry Operator Presets

Ready-to-deploy configuration presets for OpenTelemetry Operator. Each preset is a complete manifest you can copy, customize, and deploy.
