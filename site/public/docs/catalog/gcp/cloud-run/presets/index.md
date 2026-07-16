---
title: "Presets"
description: "Ready-to-deploy configuration presets for Cloud Run"
type: "preset-list"
componentSlug: "cloud-run"
componentTitle: "Cloud Run"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-public-api-service"
    rank: "01"
    title: "Public API Service"
    excerpt: "This preset creates a public, scale-to-zero HTTP API: unauthenticated callers reach the service directly on its run.app URL, instances appear with traffic and disappear when idle."
  - slug: "02-private-vpc-service"
    rank: "02"
    title: "Private VPC-Connected Backend"
    excerpt: "This preset creates an internal backend service wired into a VPC: IAM-authenticated callers only, direct VPC egress to private resources, Cloud SQL over managed Unix sockets, credentials from Secret..."
  - slug: "03-gpu-inference"
    rank: "03"
    title: "GPU Inference Service"
    excerpt: "This preset creates a GPU-backed model-serving endpoint: one NVIDIA L4 per instance, scale-to-zero so idle GPUs cost nothing, instance-based billing so the model stays resident, and IAM-authenticated..."
---

# Cloud Run Presets

Ready-to-deploy configuration presets for Cloud Run. Each preset is a complete manifest you can copy, customize, and deploy.
