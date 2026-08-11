---
title: "PagerDuty Service"
description: "Routes incidents into a PagerDuty service — the escalation tier, where PagerDuty's own rotations, overrides, and escalation policies take over."
type: "preset"
rank: "03"
presetSlug: "03-pagerduty-service"
componentSlug: "monitoring-notification-channel"
componentTitle: "Monitoring Notification Channel"
provider: "gcp"
icon: "package"
order: 3
---

# PagerDuty Service

Routes incidents into a PagerDuty service — the escalation tier, where
PagerDuty's own rotations, overrides, and escalation policies take over.

## What it configures

- `type: pagerduty` with the service integration key in
  `sensitiveLabels.serviceKey` — a managed secret. The key comes from the
  PagerDuty service's "Events API" / Google Cloud Monitoring integration.
- `deletionPolicy: PREVENT` — destroying the production paging path
  should be a deliberate act, so destroy fails until the policy is
  relaxed.

## Adjust before deploying

- **serviceKey** — supply through the platform's secret handling (the
  literal in this preset is a placeholder shape, not a working key).

## When to choose something else

For informational alerts that need no human escalation, the **On-call
Email** or **Slack Channel** presets avoid PagerDuty noise; policies can
reference several channels at once, so severity tiers usually combine
them.
