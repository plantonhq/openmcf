---
title: "Presets"
description: "Ready-to-deploy configuration presets for ContainerRegistry"
type: "preset-list"
componentSlug: "containerregistry"
componentTitle: "ContainerRegistry"
provider: "alicloud"
icon: "package"
order: 200
presets:
  - slug: "01-basic-dev"
    rank: "01"
    title: "Basic Dev Registry"
    excerpt: "This preset creates a Basic-edition Container Registry instance billed pay-as-you-go, with a single auto-created private namespace. It is the right starting point for development and small-team..."
  - slug: "02-standard-production"
    rank: "02"
    title: "Standard Production Registry"
    excerpt: "This preset creates a Standard-edition Container Registry instance on a 12-month subscription, with separate private namespaces for platform, backend, and frontend images. It fits a production team..."
  - slug: "03-advanced-enterprise"
    rank: "03"
    title: "Advanced Enterprise Registry"
    excerpt: "This preset creates an Advanced-edition Container Registry instance on a 12-month subscription, with internal, shared, and deliberately public namespaces. It suits an organization that distributes..."
---

# ContainerRegistry Presets

Ready-to-deploy configuration presets for ContainerRegistry. Each preset is a complete manifest you can copy, customize, and deploy.
