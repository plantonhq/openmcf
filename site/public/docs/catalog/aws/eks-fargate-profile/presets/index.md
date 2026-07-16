---
title: "Presets"
description: "Ready-to-deploy configuration presets for EKS Fargate Profile"
type: "preset-list"
componentSlug: "eks-fargate-profile"
componentTitle: "EKS Fargate Profile"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-namespace"
    rank: "01"
    title: "Namespace on Fargate"
    excerpt: "This preset sends every pod in one namespace to Fargate -- the simplest serverless slice of a cluster, with the pod execution role and private subnets attached by reference."
  - slug: "02-labeled-workloads"
    rank: "02"
    title: "Label-Scoped Fargate"
    excerpt: "This preset runs only opted-in pods on Fargate -- pods that carry the selector's labels -- while the rest of the namespace keeps running on node groups. The mixed-compute pattern for shared..."
---

# EKS Fargate Profile Presets

Ready-to-deploy configuration presets for EKS Fargate Profile. Each preset is a complete manifest you can copy, customize, and deploy.
