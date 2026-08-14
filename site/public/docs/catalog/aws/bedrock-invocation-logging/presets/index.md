---
title: "Presets"
description: "Ready-to-deploy configuration presets for Bedrock Invocation Logging"
type: "preset-list"
componentSlug: "bedrock-invocation-logging"
componentTitle: "Bedrock Invocation Logging"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-full-fidelity-audit"
    rank: "01"
    title: "Full-Fidelity Audit"
    excerpt: "This preset delivers invocation logs to BOTH destinations: CloudWatch for querying (with S3 spillover for oversized payloads) and S3 for retention — the canonical audit posture for production Bedrock..."
  - slug: "02-s3-archive-only"
    rank: "02"
    title: "S3 Archive Only"
    excerpt: "This preset archives text and image invocation logs to S3 without a CloudWatch arm — the low-cost retention posture when nobody queries logs interactively."
---

# Bedrock Invocation Logging Presets

Ready-to-deploy configuration presets for Bedrock Invocation Logging. Each preset is a complete manifest you can copy, customize, and deploy.
