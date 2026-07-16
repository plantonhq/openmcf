---
title: "EC2 GPU Job"
description: "An EC2 GPU job definition for ML training and CUDA workloads — one pinned GPU, raised file limits, large shared memory, and a day-long timeout."
type: "preset"
rank: "02"
presetSlug: "02-ec2-gpu-job"
componentSlug: "batch-job-definition"
componentTitle: "Batch Job Definition"
provider: "aws"
icon: "package"
order: 2
---

# EC2 GPU Job

An EC2 GPU job definition for ML training and CUDA workloads — one pinned GPU, raised file limits, large shared memory, and a day-long timeout.

## When to Use

- ML training, inference batches, and CUDA compute on EC2 GPU compute environments
- Any workload needing the EC2-only container controls Fargate cannot offer

## What It Configures

- **`gpus: 1`** — one GPU pinned exclusively to the container (the compute environment needs GPU instance types and an `ECS_AL2_NVIDIA` image type)
- **8 vCPU / 60 GiB sizing** — matched to a p3.2xlarge-class instance
- **`nofile` ulimit + 8 GiB /dev/shm** — the two limits data loaders hit first
- **Init process** — reaps zombie processes from frameworks that fork workers
- **24-hour timeout, 2 attempts** — long-running training with one retry

## What to Customize

- Replace the region/image/role placeholders
- Size `vcpus`/`memoryMib`/`gpus` to the target instance family; multi-GPU jobs raise `gpus`
- Mount training data via `volumes` (EFS file system + access point references) instead of baking it into the image
