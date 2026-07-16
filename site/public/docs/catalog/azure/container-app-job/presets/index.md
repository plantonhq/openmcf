---
title: "Presets"
description: "Ready-to-deploy configuration presets for Container App Job"
type: "preset-list"
componentSlug: "container-app-job"
componentTitle: "Container App Job"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-scheduled-batch"
    rank: "01"
    title: "Scheduled Batch Job"
    excerpt: "This preset creates a nightly batch job that runs on a cron schedule inside a Container App Environment. Each execution runs a single replica to completion (up to a 1-hour deadline), pulls its image..."
  - slug: "02-queue-worker"
    rank: "02"
    title: "Event-Triggered Queue Worker"
    excerpt: "This preset creates a queue-draining job: KEDA polls an Azure Storage Queue every 30 seconds, and when the queue holds more than 5 messages it starts executions (up to 10 at once). Each execution..."
  - slug: "03-on-demand"
    rank: "03"
    title: "On-Demand Job (Manual Trigger)"
    excerpt: "This preset creates a manually triggered job -- the database-migration model. The job definition lives in the environment ready to run; executions start on demand from the CLI (`az containerapp job..."
---

# Container App Job Presets

Ready-to-deploy configuration presets for Container App Job. Each preset is a complete manifest you can copy, customize, and deploy.
