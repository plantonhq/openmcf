---
title: "Presets"
description: "Ready-to-deploy configuration presets for API Gateway Account Settings"
type: "preset-list"
componentSlug: "api-gateway-account-settings"
componentTitle: "API Gateway Account Settings"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-cloudwatch-logging-role"
    rank: "01"
    title: "CloudWatch Logging Role"
    excerpt: "This preset sets the region's API Gateway logging role — the prerequisite for execution/access logging on every REST API stage in the region."
  - slug: "02-no-logging-posture"
    rank: "02"
    title: "No-Logging Posture"
    excerpt: "This preset explicitly manages the region with NO API Gateway logging role — applying it clears any role previously set by anyone."
---

# API Gateway Account Settings Presets

Ready-to-deploy configuration presets for API Gateway Account Settings. Each preset is a complete manifest you can copy, customize, and deploy.
