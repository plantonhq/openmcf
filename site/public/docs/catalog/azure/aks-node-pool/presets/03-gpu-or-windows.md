---
title: "GPU Node Pool Reserved by Taint"
description: "This preset creates a single-node GPU pool (NVIDIA T4) with the AKS-managed GPU driver installed, labeled and tainted so only GPU workloads schedule onto it. Scale the fixed count -- or switch to..."
type: "preset"
rank: "03"
presetSlug: "03-gpu-or-windows"
componentSlug: "aks-node-pool"
componentTitle: "AKS Node Pool"
provider: "azure"
icon: "package"
order: 3
---

# GPU Node Pool Reserved by Taint

This preset creates a single-node GPU pool (NVIDIA T4) with the AKS-managed GPU driver installed, labeled and tainted so only GPU workloads schedule onto it. Scale the fixed count -- or switch to autoscaling -- as GPU demand grows.

## When to Use

- ML inference and training, video transcoding, or CUDA workloads
- Teams adding accelerated compute to an existing cluster without touching general pools
- Any special-hardware pool that must be reserved for the workloads that need it (the same label + taint pattern applies to Windows pools -- swap `vmSize` for a Windows-supported size, set `osType: WINDOWS`, and note the cluster must carry a `windowsProfile`)

## Key Configuration Choices

- **GPU VM size** (`vmSize: Standard_NC4as_T4_v3`) -- 1× NVIDIA T4, 4 vCPUs; the smallest current-generation GPU size, good for inference
- **Managed driver** (`gpuDriver: INSTALL`) -- AKS installs the NVIDIA driver; set `NONE` if you run the GPU operator yourself
- **Reserved by taint** (`nodeTaints: ["sku=gpu:NoSchedule"]`) -- only pods with the matching toleration schedule here, so expensive GPU nodes never fill with ordinary pods
- **Labeled for selection** (`nodeLabels: {workload: gpu}`) -- pods target the pool with a nodeSelector
- **Fixed single node** (`nodeCount: 1`) -- GPU capacity is deliberate; raise the count or enable autoscaling as demand grows

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aks-cluster-arm-id>` | ARM ID of the parent AKS cluster | `AzureAksCluster` `status.outputs.cluster_id` (or reference it with `valueFrom`) |

## Related Presets

- **01-on-demand-general** -- The general-purpose pool most workloads run on
- **02-spot-cost-optimized** -- Combine the taint pattern with SPOT priority for interruptible GPU batch jobs
