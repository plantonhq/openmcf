---
title: "Presets"
description: "Ready-to-deploy configuration presets for Percona Operator for MySQL"
type: "preset-list"
componentSlug: "percona-operator-for-mysql"
componentTitle: "Percona Operator for MySQL"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard"
    excerpt: "This preset installs the Percona Operator for MySQL (XtraDB Cluster) in its standard posture: the pinned `pxc-operator` chart, own-namespace watch scope, telemetry off, structured logs, and explicit..."
---

# Percona Operator for MySQL Presets

Ready-to-deploy configuration presets for Percona Operator for MySQL. Each preset is a complete manifest you can copy, customize, and deploy.
