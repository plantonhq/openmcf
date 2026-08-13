---
title: "Presets"
description: "Ready-to-deploy configuration presets for SageMaker Notebook Instance"
type: "preset-list"
componentSlug: "sagemaker-notebook-instance"
componentTitle: "SageMaker Notebook Instance"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-starter-notebook"
    rank: "01"
    title: "Starter Notebook"
    excerpt: "This preset gives a data scientist a ready Jupyter workstation on the cheapest current-generation instance (~$0.05/hour), with the everyday Python stack installed once at creation."
  - slug: "02-locked-down-gpu-notebook"
    rank: "02"
    title: "Locked-Down GPU Notebook"
    excerpt: "This preset puts a GPU under Jupyter with the security posture tightened — no root for users, IMDSv2 only, the current Amazon Linux 2023 platform, and the team's repository cloned as the working..."
---

# SageMaker Notebook Instance Presets

Ready-to-deploy configuration presets for SageMaker Notebook Instance. Each preset is a complete manifest you can copy, customize, and deploy.
