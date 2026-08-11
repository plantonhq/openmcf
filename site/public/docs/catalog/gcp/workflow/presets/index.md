---
title: "Presets"
description: "Ready-to-deploy configuration presets for Workflow"
type: "preset-list"
componentSlug: "workflow"
componentTitle: "Workflow"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-http-orchestrator"
    rank: "01"
    title: "HTTP Orchestrator"
    excerpt: "Call internal services in order with per-step OIDC auth and retries — the canonical Workflows shape: the workflow is the coordination skeleton, the services do the work."
  - slug: "02-event-handler"
    rank: "02"
    title: "Event Handler"
    excerpt: "An event-driven step function: a GcpEventarcTrigger (destination: workflow) starts one execution per event; the workflow decodes the CloudEvent and fans out the work."
---

# Workflow Presets

Ready-to-deploy configuration presets for Workflow. Each preset is a complete manifest you can copy, customize, and deploy.
