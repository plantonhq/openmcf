---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Feature Group"
type: "preset-list"
componentSlug: "sagemaker-feature-group"
componentTitle: "SageMaker Feature Group"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-realtime-serving-features"
    rank: "01"
    title: "Realtime Serving Features"
    excerpt: "This preset is an online-only feature group: a handful of customer features served at low latency, with stale records expiring 30 days after their event time."
  - slug: "02-training-and-serving-features"
    rank: "02"
    title: "Training and Serving Features"
    excerpt: "This preset is a dual-store feature group: every write serves online at low latency AND lands in S3 under an auto-created Glue table for training datasets and point-in-time queries."
---

# SageMaker Feature Group Presets

Ready-to-deploy configuration presets for SageMaker Feature Group. Each preset is a complete manifest you can copy, customize, and deploy.
