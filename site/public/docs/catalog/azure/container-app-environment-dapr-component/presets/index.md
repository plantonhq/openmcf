---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Environment Dapr Component"
type: "preset-list"
componentSlug: "container-app-environment-dapr-component"
componentTitle: "Container App Environment Dapr Component"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-blob-state-store"
    rank: "01"
    title: "Blob Storage State Store"
    excerpt: "This preset registers a Dapr state store backed by Azure Blob Storage. Dapr-enabled apps whose `dapr.app_id` appears in `scopes` call the Dapr state API with the component name (`statestore`) and..."
  - slug: "02-servicebus-pubsub"
    rank: "02"
    title: "Service Bus Pub/Sub (Keyless)"
    excerpt: "This preset registers a Dapr pub/sub component backed by Azure Service Bus, authenticating with a managed identity instead of a connection string. Publisher and subscriber apps use Dapr's pub/sub API..."
---

# Container App Environment Dapr Component Presets

Ready-to-deploy configuration presets for Container App Environment Dapr Component. Each preset is a complete manifest you can copy, customize, and deploy.
