---
title: "GPU batch inference"
description: "Eight inference tasks with four GPUs active at a time. Each task gets one `nvidia-l4` via `nodeSelector`; container limits meet Cloud Run's GPU minimums (4 CPU / 16Gi). `gpuZonalRedundancyDisabled:..."
type: "preset"
rank: "03"
presetSlug: "03-gpu-batch-inference"
componentSlug: "cloud-run-job"
componentTitle: "Cloud Run Job"
provider: "gcp"
icon: "package"
order: 3
---

# GPU batch inference

Eight inference tasks with four GPUs active at a time. Each task gets one `nvidia-l4` via `nodeSelector`; container limits meet Cloud Run's GPU minimums (4 CPU / 16Gi). `gpuZonalRedundancyDisabled: true` trades zonal redundancy for cheaper GPU capacity — appropriate for fault-tolerant batch scoring where a failed shard can be re-run.

Requires regional GPU quota and the GEN2 execution environment.
