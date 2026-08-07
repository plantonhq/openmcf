# GPU Dedicated

This preset declares a tainted, on-demand GPU pool constrained to the g5
instance family. The `nvidia.com/gpu` taint is the standard
dedicated-pool pattern: only pods that tolerate it schedule onto — and
therefore trigger provisioning of — GPU nodes, so expensive accelerators
never launch for ordinary workloads. Consolidation is restricted to empty
nodes and node lifetime is shortened to a week.

## When to Use

- ML training/inference, transcoding, or any workload requesting
  `nvidia.com/gpu` resources
- Clusters that must guarantee GPU capacity is never consumed (or
  provisioned) by non-GPU pods

## Key Configuration Choices

- **Taint `nvidia.com/gpu=true:NoSchedule`** — pods must tolerate it to
  land on the pool; the device-plugin DaemonSet that exposes the GPUs
  must tolerate it too
- **`instance-family In [g5]`** — NVIDIA A10G-class capacity; widen the
  list (g6, p4d, ...) as your accelerator needs change
- **On-demand only** — GPU jobs rarely tolerate spot reclaims mid-run;
  build a separate spot GPU pool if your jobs checkpoint
- **`consolidationPolicy: WhenEmpty` + `consolidateAfter: 5m`** — only
  nodes with no workload pods are consolidated, and a finished node
  lingers five minutes for follow-on jobs instead of churning
- **`expireAfter: 168h`** — 7-day maximum lifetime instead of the CRD's
  720h default; expensive nodes are recycled faster
- **`limits."nvidia.com/gpu": "16"`** — pool-wide accelerator ceiling

## Placeholders to Replace

| Placeholder             | Description                            | Where to Find                                             |
| ----------------------- | -------------------------------------- | --------------------------------------------------------- |
| `<ec2-node-class-name>` | Name of the EC2NodeClass to build from | `metadata.name` of your `KubernetesKarpenterEc2NodeClass` |

## Related Presets

- **01-general-purpose-on-demand** — the baseline pool for everything else
- **02-spot-diversified** — cost-optimized pool for interruption-tolerant work
