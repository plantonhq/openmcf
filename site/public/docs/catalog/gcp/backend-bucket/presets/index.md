---
title: "Presets"
description: "Ready-to-deploy configuration presets for Backend Bucket"
type: "preset-list"
componentSlug: "backend-bucket"
componentTitle: "Backend Bucket"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-cdn-static-assets"
    rank: "01"
    title: "CDN-Cached Static Assets"
    excerpt: "The workhorse backend bucket: fingerprinted static assets (JS/CSS/images with content hashes in their paths) served through Cloud CDN with sensible TTLs, negative caching for missing-asset storms,..."
  - slug: "02-plain-origin"
    rank: "02"
    title: "Plain Origin (No CDN)"
    excerpt: "A backend bucket with edge caching off: the load balancer proxies every request to the bucket. The simplest way to put GCS content behind the same global VIP, hostname, and URL map as the rest of an..."
  - slug: "03-signed-url-private-cdn"
    rank: "03"
    title: "Private Content via CDN Signed URLs"
    excerpt: "A backend bucket serving private media through Cloud CDN, gated by signed URLs: expiring, tamper-proof links minted with the named signing key. The application signs each URL; the edge verifies the..."
---

# Backend Bucket Presets

Ready-to-deploy configuration presets for Backend Bucket. Each preset is a complete manifest you can copy, customize, and deploy.
