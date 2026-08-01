---
title: "Presets"
description: "Ready-to-deploy configuration presets for Keycloak"
type: "preset-list"
componentSlug: "keycloak"
componentTitle: "Keycloak"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-standard"
    rank: "01"
    title: "Standard preset"
    excerpt: "The production Keycloak shape: two instances clustering through the operator's discovery Service, a PostgreSQL database referenced from a `KubernetesPostgres` resource (host = its read-write Service,..."
  - slug: "02-dev-sandbox"
    rank: "02"
    title: "Dev-sandbox preset"
    excerpt: "The smallest Keycloak that starts: the `dev-mem` embedded H2 database (in memory), the plain-HTTP listener, and strict hostname resolution off so the server answers on whatever host reaches it..."
---

# Keycloak Presets

Ready-to-deploy configuration presets for Keycloak. Each preset is a complete manifest you can copy, customize, and deploy.
