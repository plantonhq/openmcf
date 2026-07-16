---
title: "Presets"
description: "Ready-to-deploy configuration presets for Project"
type: "preset-list"
componentSlug: "project"
componentTitle: "Project"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-standard-production"
    rank: "01"
    title: "Standard Production Project"
    excerpt: "This preset creates a production-grade project under a folder: billing linked, the default network suppressed, a hardening baseline of APIs pre-enabled, and destroy blocked by `deletionPolicy:..."
  - slug: "02-development"
    rank: "02"
    title: "Development Project"
    excerpt: "This preset creates a lightweight project for development environments: billing linked, a minimal API set, and the default DELETE deletion policy so teardown is one command."
---

# Project Presets

Ready-to-deploy configuration presets for Project. Each preset is a complete manifest you can copy, customize, and deploy.
