---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redshift Serverless Workgroup"
type: "preset-list"
componentSlug: "redshift-serverless-workgroup"
componentTitle: "Redshift Serverless Workgroup"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-capped-dev"
    rank: "01"
    title: "Capped Development Workgroup"
    excerpt: "This preset creates a cost-bounded serverless workgroup: the smallest practical RPU baseline (8) with a hard scaling ceiling (32 RPU), attached to a namespace from the resource graph. Billing follows..."
  - slug: "02-price-performance-production"
    rank: "02"
    title: "Price-Performance Production Workgroup"
    excerpt: "This preset creates a production workgroup where AWS owns the capacity baseline: the price-performance target sits at the balanced level (50), a 512-RPU ceiling bounds worst-case spend, all..."
  - slug: "03-private-endpoint-domain"
    rank: "03"
    title: "Governed Workgroup with Private Endpoints and a Custom Domain"
    excerpt: "This preset creates a production serverless workgroup whose spend, reachability, and identity are governed declaratively: usage limits deactivate queries when daily RPU-hours exceed the cap and log..."
---

# Redshift Serverless Workgroup Presets

Ready-to-deploy configuration presets for Redshift Serverless Workgroup. Each preset is a complete manifest you can copy, customize, and deploy.
