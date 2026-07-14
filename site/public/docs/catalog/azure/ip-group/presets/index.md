---
title: "Presets"
description: "Ready-to-deploy configuration presets for IP Group"
type: "preset-list"
componentSlug: "ip-group"
componentTitle: "IP Group"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-branch-offices"
    rank: "01"
    title: "Branch-Office Address Set"
    excerpt: "This preset creates an IP Group carrying the CIDR ranges of an organization's branch offices. It is the building block of maintainable firewall policy: instead of enumerating branch CIDRs inline in..."
  - slug: "02-on-prem-datacenter"
    rank: "02"
    title: "On-Premises Datacenter Ranges"
    excerpt: "This preset creates an IP Group carrying the private ranges of an on-premises datacenter -- the destination half of the classic hybrid egress policy (\"branches and Azure workloads may reach the..."
---

# IP Group Presets

Ready-to-deploy configuration presets for IP Group. Each preset is a complete manifest you can copy, customize, and deploy.
