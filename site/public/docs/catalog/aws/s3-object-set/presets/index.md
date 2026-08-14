---
title: "Presets"
description: "Ready-to-deploy configuration presets for S3 Object Set"
type: "preset-list"
componentSlug: "s3-object-set"
componentTitle: "S3 Object Set"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-config-files"
    rank: "01"
    title: "Configuration Files"
    excerpt: "This preset uploads a set of application configuration files to an S3 bucket. It demonstrates the inline content pattern for managing JSON and YAML configuration files as infrastructure, with proper..."
  - slug: "02-static-website-assets"
    rank: "02"
    title: "Static Website Assets"
    excerpt: "This preset uploads a small static-site asset set — an HTML entry page, a fingerprint-friendly stylesheet, and a redirect marker for a moved page — with the cache posture each asset class deserves...."
  - slug: "03-encrypted-compliance-drop"
    rank: "03"
    title: "Encrypted Compliance Drop"
    excerpt: "This preset writes an audit artifact under WORM (write-once-read-many) retention: customer-managed KMS encryption, a SHA256 upload checksum, and COMPLIANCE-mode Object Lock that nobody — including..."
  - slug: "04-promote-golden-artifacts"
    rank: "04"
    title: "Promote Golden Artifacts"
    excerpt: "This preset copies released build artifacts from a golden artifacts bucket into an environment's bucket, server-side — the bytes move inside S3 and never through the deploy host, so artifact size is..."
---

# S3 Object Set Presets

Ready-to-deploy configuration presets for S3 Object Set. Each preset is a complete manifest you can copy, customize, and deploy.
