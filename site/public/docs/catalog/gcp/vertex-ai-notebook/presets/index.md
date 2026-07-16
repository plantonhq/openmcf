---
title: "Presets"
description: "Ready-to-deploy configuration presets for Vertex AI Notebook"
type: "preset-list"
componentSlug: "vertex-ai-notebook"
componentTitle: "Vertex AI Notebook"
provider: "gcp"
icon: "package"
order: 200
presets:
  - slug: "01-basic-notebook"
    rank: "01"
    title: "Basic Notebook"
    excerpt: "The minimal managed JupyterLab environment: a CPU-only Workbench instance in the ambient project on the service's default Workbench image."
  - slug: "02-gpu-ml-notebook"
    rank: "02"
    title: "GPU ML Notebook"
    excerpt: "A GPU-accelerated JupyterLab environment for deep learning training, pinned to a pre-purchased capacity reservation."
  - slug: "03-secure-private-notebook"
    rank: "03"
    title: "Secure Private Notebook"
    excerpt: "The hardened posture for regulated data: private networking, CMEK, Shielded VM, Confidential Computing, and per-user credentials — every security lever the platform models, composed from first-class..."
---

# Vertex AI Notebook Presets

Ready-to-deploy configuration presets for Vertex AI Notebook. Each preset is a complete manifest you can copy, customize, and deploy.
