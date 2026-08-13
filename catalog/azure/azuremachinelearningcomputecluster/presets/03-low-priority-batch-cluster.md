# Low-Priority Batch Cluster

This preset creates a wide, cheap, evictable cluster for fault-tolerant batch work -- eight spot-class nodes at a fraction of dedicated cost, releasing five minutes after their last job.

## When to Use

- Hyperparameter sweeps, retraining fleets, and other embarrassingly-parallel batch work
- Any job that checkpoints and resumes (or is cheap to simply rerun)
- Cost-sensitive estates where GPU/CPU hours dominate the ML bill

## Key Configuration Choices

- **`LOW_PRIORITY`** -- the contract, not just a discount: Azure evicts nodes at any time, taking running work with them; your training code must checkpoint/resume or tolerate reruns. Priority is ForceNew -- converting to dedicated later means replacing the cluster
- **`maxNodeCount: 8`** -- wide fan-out suits sweeps; low-priority quota is separate from dedicated quota
- **`PT5M` idle duration** -- batch fleets rarely benefit from warm nodes; release fast

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-machine-learning-workspace-id>` | ARM ID of the parent workspace | `AzureMachineLearningWorkspace` status outputs (`machine_learning_workspace_id`), or reference it with valueFrom |

Never put an unresumable multi-hour job on this cluster -- one eviction costs more than the discount saved.
