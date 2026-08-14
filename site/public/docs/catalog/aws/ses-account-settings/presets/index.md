---
title: "Presets"
description: "Ready-to-deploy configuration presets for SES Account Settings"
type: "preset-list"
componentSlug: "ses-account-settings"
componentTitle: "SES Account Settings"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-reputation-defaults"
    rank: "01"
    title: "Reputation Defaults"
    excerpt: "This preset sets the recommended account-wide suppression posture: hard bounces and spam complaints automatically suppress recipient addresses across every send from the account."
  - slug: "02-vdm-analytics"
    rank: "02"
    title: "VDM Analytics"
    excerpt: "This preset layers the Virtual Deliverability Manager on top of the reputation defaults: engagement dashboards plus Guardian delivery optimization."
---

# SES Account Settings Presets

Ready-to-deploy configuration presets for SES Account Settings. Each preset is a complete manifest you can copy, customize, and deploy.
