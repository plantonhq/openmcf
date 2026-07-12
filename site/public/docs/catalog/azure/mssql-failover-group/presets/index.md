---
title: "Presets"
description: "Ready-to-deploy configuration presets for MSSQL Failover Group"
type: "preset-list"
componentSlug: "mssql-failover-group"
componentTitle: "MSSQL Failover Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-automatic-failover"
    rank: "01"
    title: "Automatic Failover Group"
    excerpt: "This preset creates a failover group with automatic failover and a 60-minute grace period -- the standard production DR shape. Azure fails the group over to the partner on its own when it detects a..."
  - slug: "02-manual-failover"
    rank: "02"
    title: "Manual Failover Group"
    excerpt: "This preset creates a failover group with manual failover -- an operator (or an automation pipeline) initiates every failover. This is the choice when you want a human decision in the loop before..."
---

# MSSQL Failover Group Presets

Ready-to-deploy configuration presets for MSSQL Failover Group. Each preset is a complete manifest you can copy, customize, and deploy.
