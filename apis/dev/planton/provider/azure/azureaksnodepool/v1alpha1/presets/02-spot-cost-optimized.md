# Spot Cost-Optimized Node Pool

This preset creates a cost-optimized AKS user node pool using Azure Spot VMs, which provide 30-90% savings over on-demand pricing. The pool scales to zero when idle and up to 20 nodes under load. Spot VMs can be evicted when Azure needs capacity, so this pool is only suitable for fault-tolerant, stateless workloads.

## When to Use

- Batch processing, CI/CD runners, and other interruptible workloads
- Dev/test environments where occasional eviction is acceptable
- Stateless services with proper retry logic and graceful shutdown handling
- Cost-sensitive workloads that can be rescheduled to on-demand pools during eviction

## Key Configuration Choices

- **Spot priority** (`priority: SPOT`) -- 30-90% cost savings over on-demand; nodes can be evicted at any time and automatically carry the `scalesetpriority=spot:NoSchedule` taint
- **Delete on eviction** (`evictionPolicy: EVICTION_DELETE`) -- evicted VMs are deleted and billing stops (vs. DEALLOCATE, which keeps disks billing)
- **No price ceiling** (`spotMaxPrice` unset = -1) -- pay up to the on-demand price, never price-evicted; capacity eviction still applies
- **Scale-to-zero** (`minCount: 0`) -- no cost when there are no pods to schedule
- **3 availability zones** (`zones`) -- spreads across zones for better Spot capacity availability
- **User mode** (`mode: USER`) -- Spot cannot be used for system pools

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aks-cluster-arm-id>` | ARM ID of the parent AKS cluster | `AzureAksCluster` `status.outputs.cluster_id` (or reference it with `valueFrom`) |

## Related Presets

- **01-on-demand-general** -- Use instead for workloads that cannot tolerate eviction (production APIs, stateful services)
- **03-gpu-or-windows** -- Use instead for GPU workloads reserved by taint
