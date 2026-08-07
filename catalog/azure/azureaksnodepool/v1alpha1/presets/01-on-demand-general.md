# On-Demand General Purpose Node Pool

This preset creates a general-purpose AKS user node pool with on-demand (regular) VMs, autoscaling from 2 to 10 nodes across 3 availability zones. This is the standard configuration for production application workloads that need reliable, non-preemptible compute.

## When to Use

- Production application workloads that cannot tolerate node eviction
- General-purpose services (web apps, APIs, microservices) with moderate CPU/memory needs
- Teams adding a dedicated user node pool to an existing AKS cluster
- Workloads requiring high availability across availability zones

## Key Configuration Choices

- **On-demand VMs** (`priority` unset = REGULAR) -- no risk of eviction; suitable for all workload types
- **Standard_D4s_v5** (`vmSize`) -- 4 vCPUs, 16 GiB RAM; balanced general-purpose compute
- **Autoscaling 2-10** (`autoScalingEnabled: true`, `minCount: 2`, `maxCount: 10`) -- always at least 2 nodes for HA; scales up under load
- **3 availability zones** (`zones`) -- distributes nodes across zones for the 99.95% SLA
- **User mode** (`mode: USER`) -- runs application workloads, separated from the cluster's system pool

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aks-cluster-arm-id>` | ARM ID of the parent AKS cluster | `AzureAksCluster` `status.outputs.cluster_id` (or reference it with `valueFrom`) |

## Related Presets

- **02-spot-cost-optimized** -- Use instead for fault-tolerant, stateless workloads that can tolerate eviction for 30-90% savings
- **03-gpu-or-windows** -- Use instead for GPU workloads reserved by taint
