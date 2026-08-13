---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitoring Notification Channel"
type: "preset-list"
componentSlug: "monitoring-notification-channel"
componentTitle: "Monitoring Notification Channel"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-oncall-email"
    rank: "01"
    title: "On-call Email"
    excerpt: "The simplest paging path: incidents from every referencing alert policy land in one inbox. No external credentials, no webhook plumbing."
  - slug: "02-slack-channel"
    rank: "02"
    title: "Slack Channel"
    excerpt: "Posts incident notifications into a Slack channel — the team-visibility tier of alerting, where everyone sees opens and closes without being paged individually."
  - slug: "03-pagerduty-service"
    rank: "03"
    title: "PagerDuty Service"
    excerpt: "Routes incidents into a PagerDuty service — the escalation tier, where PagerDuty's own rotations, overrides, and escalation policies take over."
---

# Monitoring Notification Channel Presets

Ready-to-deploy configuration presets for Monitoring Notification Channel. Each preset is a complete manifest you can copy, customize, and deploy.
