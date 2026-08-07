---
title: "Presets"
description: "Ready-to-deploy configuration presets for Apache Superset"
type: "preset-list"
componentSlug: "apache-superset"
componentTitle: "Apache Superset"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-team-bi"
    rank: "01"
    title: "Team BI"
    excerpt: "The smallest honest Superset: the web application against a composed PostgreSQL metadata database (the ONE required input — dashboards, charts, users and the encrypted datasource credentials live..."
  - slug: "02-production-bi"
    rank: "02"
    title: "Production BI"
    excerpt: "The full analytics posture: two web replicas, a Celery worker pair for async SQL Lab queries and thumbnails, beat firing scheduled alerts & reports, dashboard-level RBAC — over a composed PostgreSQL..."
---

# Apache Superset Presets

Ready-to-deploy configuration presets for Apache Superset. Each preset is a complete manifest you can copy, customize, and deploy.
