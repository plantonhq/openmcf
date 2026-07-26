---
title: "Presets"
description: "Ready-to-deploy configuration presets for Grafana Loki"
type: "preset-list"
componentSlug: "grafana-loki"
componentTitle: "Grafana Loki"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-dev-single"
    rank: "01"
    title: "Dev single-node Loki"
    excerpt: "The smallest honest Loki: one monolithic replica on a filesystem volume with the nginx gateway. No object store, no tenancy — a dev loop or a single-node cluster where logs are convenient but not yet..."
  - slug: "02-production-scalable"
    rank: "02"
    title: "Production scalable Loki"
    excerpt: "The write/read/backend topology on object storage — how Loki runs in production. Ingest, query and compaction scale independently, and chunks live in an S3-compatible bucket (an in-cluster..."
  - slug: "03-multitenant"
    rank: "03"
    title: "Multi-tenant shared Loki"
    excerpt: "One Loki serving several teams with isolation: every push and query carries an `X-Scope-OrgID` tenant header, and the gateway enforces HTTP basic auth per tenant. Passwords are supplied as bcrypt..."
---

# Grafana Loki Presets

Ready-to-deploy configuration presets for Grafana Loki. Each preset is a complete manifest you can copy, customize, and deploy.
