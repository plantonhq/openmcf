---
title: "Presets"
description: "Ready-to-deploy configuration presets for CloudWatch Composite Alarm"
type: "preset-list"
componentSlug: "cloudwatch-composite-alarm"
componentTitle: "CloudWatch Composite Alarm"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-shared-cause-outage"
    rank: "01"
    title: "Preset: Shared-Cause Outage Page"
    excerpt: "**Use case:** Page on-call once for an outage with many symptoms."
  - slug: "02-maintenance-suppressed"
    rank: "02"
    title: "Preset: Maintenance-Suppressed Paging"
    excerpt: "**Use case:** Stop paging during planned maintenance without deleting or disabling alarms."
---

# CloudWatch Composite Alarm Presets

Ready-to-deploy configuration presets for CloudWatch Composite Alarm. Each preset is a complete manifest you can copy, customize, and deploy.
