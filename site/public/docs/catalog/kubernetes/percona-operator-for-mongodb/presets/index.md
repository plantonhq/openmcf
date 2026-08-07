---
title: "Presets"
description: "Ready-to-deploy configuration presets for Percona Operator for MongoDB"
type: "preset-list"
componentSlug: "percona-operator-for-mongodb"
componentTitle: "Percona Operator for MongoDB"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard"
    excerpt: "This preset installs the Percona Operator for MongoDB in its standard posture: the pinned `psmdb-operator` chart, own-namespace watch scope, telemetry off, structured logs, and explicit control-plane..."
---

# Percona Operator for MongoDB Presets

Ready-to-deploy configuration presets for Percona Operator for MongoDB. Each preset is a complete manifest you can copy, customize, and deploy.
