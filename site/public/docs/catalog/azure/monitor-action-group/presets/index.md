---
title: "Presets"
description: "Ready-to-deploy configuration presets for Monitor Action Group"
type: "preset-list"
componentSlug: "monitor-action-group"
componentTitle: "Monitor Action Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-oncall-team"
    rank: "01"
    title: "On-Call Team Notifications"
    excerpt: "This preset creates the human-notification action group: email, SMS, and Azure mobile app push for the on-call rotation. Metric alerts, query alerts, and Service Health alerts all reference this one..."
  - slug: "02-automation-hooks"
    rank: "02"
    title: "Automation Hooks (Webhook + Function)"
    excerpt: "This preset creates a machine-only action group: an Entra-authenticated webhook into an incident platform and an HTTP-triggered Azure Function for programmatic remediation. Alert rules typically..."
  - slug: "03-role-fanout"
    rank: "03"
    title: "ARM Role Fan-Out"
    excerpt: "This preset notifies every user assigned an ARM role on the subscription -- here the built-in Owner role. RBAC membership becomes the notification list, so joiners and leavers are handled by access..."
---

# Monitor Action Group Presets

Ready-to-deploy configuration presets for Monitor Action Group. Each preset is a complete manifest you can copy, customize, and deploy.
