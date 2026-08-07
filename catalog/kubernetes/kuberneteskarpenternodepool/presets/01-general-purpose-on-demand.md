# General Purpose On-Demand

This preset declares the default fleet most clusters start with: amd64
on-demand nodes from current-generation compute-, general- and
memory-optimized instance families, with consolidation active and a CPU
ceiling on the whole pool. Clusters routinely run several NodePools — this
is the baseline one; the machine template (AMIs, subnets, security groups,
IAM) lives in the `KubernetesKarpenterEc2NodeClass` it references.

## When to Use

- The first NodePool of a new Karpenter installation
- General workloads without spot tolerance, GPU needs, or dedicated-pool
  isolation

## Key Configuration Choices

- **On-demand only** (`karpenter.sh/capacity-type In [on-demand]`) —
  predictable capacity with no interruption exposure
- **Category + generation requirements instead of instance-type lists** —
  `instance-category In [c, m, r]` with `instance-generation Gt 4` gives
  Karpenter a wide, current-generation catalog to right-size from
- **`consolidationPolicy: WhenEmptyOrUnderutilized`** (CRD default, made
  explicit) — nodes are replaced or removed whenever cheaper capacity
  fits the pods; `consolidateAfter: 1m` adds a brief settle window over
  the immediate `0s` default
- **`expireAfter: 720h`** (CRD default) — 30-day maximum node lifetime as
  the eventually-consistent recycling mechanism
- **`limits.cpu: "1000"`** — provisioning stops at 1000 vCPUs; the pool's
  cost ceiling
- **`nodeClassRef` uses spec defaults** (`karpenter.k8s.aws` /
  `EC2NodeClass`) with a literal `value:` name

## Placeholders to Replace

| Placeholder             | Description                            | Where to Find                                     |
| ----------------------- | -------------------------------------- | ------------------------------------------------- |
| `<ec2-node-class-name>` | Name of the EC2NodeClass to build from | `metadata.name` of your `KubernetesKarpenterEc2NodeClass` |

## Related Presets

- **02-spot-diversified** — the cost-optimized spot companion pool
- **03-gpu-dedicated** — tainted GPU pool for accelerated workloads
