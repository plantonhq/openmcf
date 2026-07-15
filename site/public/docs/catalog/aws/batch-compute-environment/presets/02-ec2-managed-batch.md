---
title: "EC2 Managed Batch"
description: "An EC2 On-Demand compute environment with `optimal` instance selection and the BEST_FIT_PROGRESSIVE strategy — the configuration that keeps day-2 infrastructure changes cheap (in-place updates..."
type: "preset"
rank: "02"
presetSlug: "02-ec2-managed-batch"
componentSlug: "batch-compute-environment"
componentTitle: "Batch Compute Environment"
provider: "aws"
icon: "package"
order: 2
---

# EC2 Managed Batch

An EC2 On-Demand compute environment with `optimal` instance selection and the BEST_FIT_PROGRESSIVE strategy — the configuration that keeps day-2 infrastructure changes cheap (in-place updates instead of environment replacement).

## When to Use

- Jobs that exceed Fargate's sizing (more than 16 vCPU / 120 GiB per job)
- Workloads needing GPUs, custom AMIs, privileged containers, or kernel tuning
- Steady batch pipelines where EC2 pricing beats Fargate at sustained utilization

## What It Configures

- **`optimal` instance types** — Batch picks from the C, M, and R families to match each job's resource requirements
- **BEST_FIT_PROGRESSIVE** — cheapest fitting types, falling forward on capacity; one of the three strategies that (with the service-linked role) keeps the environment eligible for in-place infrastructure updates
- **Scale-to-zero floor** — `minVcpus: 0` keeps no instances warm when idle
- **Update policy** — running jobs get 60 minutes to finish before instances are replaced during updates

## What to Customize

- Replace the region/subnet/security-group placeholders and the ECS instance profile ARN (or reference an `AwsIamInstanceProfile`)
- Pin explicit `instanceTypes` (e.g. `["c5.xlarge", "c5.2xlarge"]`) for predictable performance; add a `launchTemplate` for custom AMIs or user data
- Set a non-zero `minVcpus` to keep capacity warm for latency-sensitive queues
