---
title: "Presets"
description: "Ready-to-deploy configuration presets for Grafana"
type: "preset-list"
componentSlug: "grafana"
componentTitle: "Grafana"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-dashboards"
    rank: "01"
    title: "Dev dashboards preset"
    excerpt: "The smallest useful Grafana: one ephemeral replica, chart-generated admin credentials, and a single provisioned Prometheus datasource — a real dashboard endpoint for a dev loop without any production..."
  - slug: "02-persistent-team"
    rank: "02"
    title: "Persistent team preset"
    excerpt: "The single stateful instance most teams actually want: a 10Gi volume under Grafana's embedded database so hand-built dashboards, users and preferences survive pod restarts, a Prometheus datasource..."
  - slug: "03-ha-external-db"
    rank: "03"
    title: "HA external database preset"
    excerpt: "The production posture: two Grafana replicas behind one Service with ALL state — dashboards, users, sessions, preferences — in an external Postgres. Any pod can die, be drained or be upgraded and..."
---

# Grafana Presets

Ready-to-deploy configuration presets for Grafana. Each preset is a complete manifest you can copy, customize, and deploy.
