---
title: "Presets"
description: "Ready-to-deploy configuration presets for Batch Job Definition"
type: "preset-list"
componentSlug: "batch-job-definition"
componentTitle: "Batch Job Definition"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-fargate-container-job"
    rank: "01"
    title: "Fargate Container Job"
    excerpt: "A Fargate job definition with the full production posture: split identities, parameterized command, Spot-safe retry discrimination, and a hard timeout."
  - slug: "02-ec2-gpu-job"
    rank: "02"
    title: "EC2 GPU Job"
    excerpt: "An EC2 GPU job definition for ML training and CUDA workloads — one pinned GPU, raised file limits, large shared memory, and a day-long timeout."
  - slug: "03-eks-pod-job"
    rank: "03"
    title: "EKS Pod Job"
    excerpt: "A Batch-on-EKS job definition — the workload half of an EKS-attached compute environment: a hardened pipeline pod with an init container, in-memory scratch, and secret-projected configuration."
---

# Batch Job Definition Presets

Ready-to-deploy configuration presets for Batch Job Definition. Each preset is a complete manifest you can copy, customize, and deploy.
