---
title: "Presets"
description: "Ready-to-deploy configuration presets for OpenFGA"
type: "preset-list"
componentSlug: "openfga"
componentTitle: "OpenFGA"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-with-postgres-datastore"
    rank: "01"
    title: "PostgreSQL datastore preset"
    excerpt: "The production OpenFGA shape: three stateless server replicas sharing a PostgreSQL datastore, schema migrations running as an init container in every pod (idempotent — `openfga migrate` gates each..."
---

# OpenFGA Presets

Ready-to-deploy configuration presets for OpenFGA. Each preset is a complete manifest you can copy, customize, and deploy.
