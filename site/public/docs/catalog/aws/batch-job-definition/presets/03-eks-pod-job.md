---
title: "EKS Pod Job"
description: "A Batch-on-EKS job definition — the workload half of an EKS-attached compute environment: a hardened pipeline pod with an init container, in-memory scratch, and secret-projected configuration."
type: "preset"
rank: "03"
presetSlug: "03-eks-pod-job"
componentSlug: "batch-job-definition"
componentTitle: "Batch Job Definition"
provider: "aws"
icon: "package"
order: 3
---

# EKS Pod Job

A Batch-on-EKS job definition — the workload half of an EKS-attached compute environment: a hardened pipeline pod with an init container, in-memory scratch, and secret-projected configuration.

## When to Use

- Batch workloads scheduled onto an EKS cluster (a compute environment with `eksConfiguration`) instead of ECS/Fargate
- Teams standardizing on Kubernetes-native controls — service accounts (IRSA), pod security contexts, and Kubernetes secrets

## What It Configures

- **Main + init container** — the init container runs to completion (fetch reference data) before the pipeline starts
- **Kubernetes sizing** — cpu/memory requests and limits (Batch schedules the job by these; its EKS counterpart of `vcpus`/`memoryMib`)
- **Hardened security context** — non-root UID/GID asserted, privilege escalation off, read-only root filesystem
- **`hostNetwork: false`** — the pod gets its own network namespace (AWS's Batch-pod default is `true`); required for VPC-CNI pod networking
- **`serviceAccountName`** — IRSA / Pod Identity grants the job's code AWS permissions, the EKS-native counterpart of `jobRole`
- **Volumes** — tmpfs scratch (`emptyDir` with `medium: Memory`, counted against the container's memory) and a Kubernetes secret projected read-only
- **2 attempts, 2-hour timeout** — retry and wall-clock guardrails work identically across both workload arms

## What to Customize

- Replace the region/image placeholders; point `secretName` at a secret in the job's namespace (the namespace comes from the compute environment's `eksConfiguration`)
- Size `resources` to the cluster's node types; add `nvidia.com/gpu` limits for GPU nodes
- Drop `hostNetwork: false` to accept AWS's default (node network namespace) for egress-only workloads at higher pod density
