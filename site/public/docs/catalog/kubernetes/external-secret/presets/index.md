---
title: "Presets"
description: "Ready-to-deploy configuration presets for External Secret"
type: "preset-list"
componentSlug: "external-secret"
componentTitle: "External Secret"
provider: "kubernetes"
icon: "package"
order: 200
presets:
  - slug: "01-explicit-keys"
    rank: "01"
    title: "Explicit Keys (Application Credentials)"
    excerpt: "This preset syncs two named fields (`username`, `password`) from one structured backend entry into a Kubernetes Secret, refreshed hourly. Each mapping is explicit and reviewable — the standard form..."
  - slug: "02-extract-json-document"
    rank: "02"
    title: "Extract a JSON Document (Bulk Pull with Rewrite)"
    excerpt: "This preset pulls ALL properties of one structured backend entry — a JSON document of related credentials — into a Kubernetes Secret in a single `dataFrom.extract`, with a regex rewrite that strips..."
  - slug: "03-docker-registry-template"
    rank: "03"
    title: "Docker Registry Pull Secret (Template)"
    excerpt: "This preset syncs registry credentials from the backend and TEMPLATES them into a `kubernetes.io/dockerconfigjson` Secret — the shape `imagePullSecrets` expects. The synced `username`/`password` keys..."
---

# External Secret Presets

Ready-to-deploy configuration presets for External Secret. Each preset is a complete manifest you can copy, customize, and deploy.
