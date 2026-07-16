---
title: "Presets"
description: "Ready-to-deploy configuration presets for Region Network Endpoint Group"
type: "preset-list"
componentSlug: "region-network-endpoint-group"
componentTitle: "Region Network Endpoint Group"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-cloud-run-neg"
    rank: "01"
    title: "Cloud Run behind a Load Balancer"
    excerpt: "The most common serverless NEG: put a Cloud Run service behind a global external Application Load Balancer so it can serve a custom domain with Cloud CDN, Cloud Armor, and IAP in front of it."
  - slug: "02-private-service-connect-neg"
    rank: "02"
    title: "Private Service Connect Backend"
    excerpt: "A PSC network endpoint group fronting a published producer service or a Google API — the way to reach a private, PSC-published backend from a load balancer without exposing it to the internet."
  - slug: "03-internet-fqdn-neg"
    rank: "03"
    title: "External Internet Origin"
    excerpt: "An internet NEG that lets a Google Cloud load balancer front an external origin — an on-prem service or a third-party API reached by FQDN — so it can sit behind Cloud CDN, Cloud Armor, and a single..."
---

# Region Network Endpoint Group Presets

Ready-to-deploy configuration presets for Region Network Endpoint Group. Each preset is a complete manifest you can copy, customize, and deploy.
