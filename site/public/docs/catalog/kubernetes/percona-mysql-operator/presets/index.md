---
title: "Presets"
description: "Ready-to-deploy configuration presets for Percona MySQL Operator"
type: "preset-list"
componentSlug: "percona-mysql-operator"
componentTitle: "Percona MySQL Operator"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard"
    excerpt: "This preset installs the Percona Operator for MySQL (XtraDB Cluster) in its standard posture: the pinned `pxc-operator` chart, own-namespace watch scope, telemetry off, structured logs, and explicit..."
---

# Percona MySQL Operator Presets

Ready-to-deploy configuration presets for Percona MySQL Operator. Each preset is a complete manifest you can copy, customize, and deploy.
