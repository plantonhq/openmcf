---
title: "GPU ML Notebook"
description: "A GPU-accelerated JupyterLab environment for deep learning training, pinned to a pre-purchased capacity reservation."
type: "preset"
rank: "02"
presetSlug: "02-gpu-ml-notebook"
componentSlug: "vertex-ai-notebook"
componentTitle: "Vertex AI Notebook"
provider: "gcp"
icon: "package"
order: 2
---

# GPU ML Notebook

A GPU-accelerated JupyterLab environment for deep learning training,
pinned to a pre-purchased capacity reservation.

## What this preset creates

A Workbench instance named `gpu-training` in `us-central1-a` on an
`n1-standard-8` VM with one NVIDIA Tesla T4 (16 GB VRAM), a 200 GB SSD
boot disk, a 500 GB SSD data disk, and the TensorFlow GPU deep learning
image. The instance consumes GPU capacity from the named Compute Engine
reservation, so training never stalls on zonal GPU stockouts.

## When to use

- Training neural networks (CNNs, transformers)
- Fine-tuning pre-trained models
- GPU-accelerated data processing (RAPIDS, cuDF)
- Hyperparameter tuning experiments

## Remix ideas

- Drop the `reservationAffinity` block to run on-demand, or set
  `consumeReservationType: RESERVATION_ANY` to opportunistically use any
  matching open reservation.
- Switch `type` to `NVIDIA_L4` with a `g2-standard-*` machine type for
  better inference-oriented price/performance.
- Set `desiredState: STOPPED` between training runs — GPU compute
  billing stops, disks and environment persist.
- The `workbench-instances` image ships JupyterLab with the common ML stack; the service installs GPU drivers matching the accelerator config. Install PyTorch or TensorFlow in the notebook environment (or use a custom `containerImage`) for framework-specific work.
