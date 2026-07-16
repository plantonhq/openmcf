---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Run Job"
type: "preset-list"
componentSlug: "cloud-run-job"
componentTitle: "Cloud Run Job"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-batch-etl-basic"
    rank: "01"
    title: "Batch ETL — basic"
    excerpt: "A single-task Cloud Run job for nightly extract-transform-load work. One execution runs one task that must succeed for the run to succeed (`taskCount: 1`). Retries are disabled (`maxRetries: 0`) so a..."
  - slug: "02-parallel-vpc-cleanup"
    rank: "02"
    title: "Parallel VPC cleanup"
    excerpt: "Ten cleanup tasks with five running concurrently — the standard map-reduce shape on Cloud Run jobs. Tasks reach private APIs through direct VPC egress (`networkInterfaces` + `PRIVATE_RANGES_ONLY`)..."
  - slug: "03-gpu-batch-inference"
    rank: "03"
    title: "GPU batch inference"
    excerpt: "Eight inference tasks with four GPUs active at a time. Each task gets one `nvidia-l4` via `nodeSelector`; container limits meet Cloud Run's GPU minimums (4 CPU / 16Gi). `gpuZonalRedundancyDisabled:..."
---

# Cloud Run Job Presets

Ready-to-deploy configuration presets for Cloud Run Job. Each preset is a complete manifest you can copy, customize, and deploy.
