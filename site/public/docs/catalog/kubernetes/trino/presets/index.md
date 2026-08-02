---
title: "Presets"
description: "Ready-to-deploy configuration presets for Trino"
type: "preset-list"
componentSlug: "trino"
componentTitle: "Trino"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-sql-playground"
    rank: "01"
    title: "SQL playground"
    excerpt: "The two-line Trino: a coordinator and two workers with PASSWORD authentication on (module-generated admin credential in the `trino-auth` Secret — upstream's open, anyone-can-query default never..."
  - slug: "02-federated-warehouse"
    rank: "02"
    title: "Federated warehouse"
    excerpt: "The production posture: a composed PostgreSQL warehouse queryable through Trino (and JOIN-able against any other catalog), autoscaled workers (up to 6 on CPU) that drain running queries before..."
---

# Trino Presets

Ready-to-deploy configuration presets for Trino. Each preset is a complete manifest you can copy, customize, and deploy.
