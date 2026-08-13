---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Inference Profile"
type: "preset-list"
componentSlug: "bedrock-inference-profile"
componentTitle: "Bedrock Inference Profile"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-per-application-tracking"
    rank: "01"
    title: "Per-Application Tracking"
    excerpt: "This preset creates one application inference profile routing straight to a foundation model — the per-consumer attribution unit: give each service its own copy of this preset (named after the..."
  - slug: "02-cross-region-routing"
    rank: "02"
    title: "Cross-Region Routing"
    excerpt: "This preset creates an application profile sourced from an AWS system-defined geo inference profile (`us.` prefix) — invocations then ride AWS's cross-region capacity pools for higher availability..."
---

# Bedrock Inference Profile Presets

Ready-to-deploy configuration presets for Bedrock Inference Profile. Each preset is a complete manifest you can copy, customize, and deploy.
