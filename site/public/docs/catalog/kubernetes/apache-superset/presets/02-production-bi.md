---
title: "Production BI"
description: "The full analytics posture: two web replicas, a Celery worker pair for async SQL Lab queries and thumbnails, beat firing scheduled alerts & reports, dashboard-level RBAC — over a composed PostgreSQL..."
type: "preset"
rank: "02"
presetSlug: "02-production-bi"
componentSlug: "apache-superset"
componentTitle: "Apache Superset"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production BI

The full analytics posture: two web replicas, a Celery worker pair for
async SQL Lab queries and thumbnails, beat firing scheduled alerts &
reports, dashboard-level RBAC — over a composed PostgreSQL metadata
database and a composed Valkey cache/broker.

Both credentials ride environment references to the composed
resources' own Secrets (the chart's bring-your-own mechanism); the
bundled database/redis subcharts — which ride frozen legacy image
lines upstream — never ship from this kind. Connect data sources from
the UI or API; they are encrypted at rest with the module-generated
(and deliberately STABLE) session-signing key.
