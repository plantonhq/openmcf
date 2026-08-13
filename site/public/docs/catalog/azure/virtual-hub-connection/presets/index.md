---
title: "Presets"
description: "Ready-to-deploy configuration presets for Virtual Hub Connection"
type: "preset-list"
componentSlug: "virtual-hub-connection"
componentTitle: "Virtual Hub Connection"
provider: "azure"
icon: "package"
order: 200
presets:
  - slug: "01-any-to-any-attachment"
    rank: "01"
    title: "Any-to-Any Attachment"
    excerpt: "This preset creates the simplest attachment: the spoke joins the hub with ARM's default routing -- associated with and propagating to the hub's built-in default route table, so every connected..."
  - slug: "02-isolated-spoke"
    rank: "02"
    title: "Isolated Spoke"
    excerpt: "This preset creates the classic isolation attachment: the spoke is routed by an isolated route table (so it never learns other isolated spokes' routes) and propagates its prefixes only to tables..."
---

# Virtual Hub Connection Presets

Ready-to-deploy configuration presets for Virtual Hub Connection. Each preset is a complete manifest you can copy, customize, and deploy.
