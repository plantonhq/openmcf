---
title: "Presets"
description: "Ready-to-deploy configuration presets for EFS Access Point"
type: "preset-list"
componentSlug: "efs-access-point"
componentTitle: "EFS Access Point"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-app-data"
    rank: "01"
    title: "Application Data Access Point"
    excerpt: "Enforced POSIX identity (1000:1000) with a dedicated `/app-data` root directory, auto-created on first mount. The standard shape for giving one application least-privilege access to a shared file..."
---

# EFS Access Point Presets

Ready-to-deploy configuration presets for EFS Access Point. Each preset is a complete manifest you can copy, customize, and deploy.
