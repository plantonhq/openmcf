# Fixed-Size Dedicated Pool

This preset adds a fixed two-node pool reserved for system workloads: a `NoSchedule` taint keeps ordinary pods off the nodes, and a matching label lets system components target them explicitly.

## When to Use

- Isolating system components (ingress controllers, observability agents, operators) from application churn
- Workloads that need a stable node count rather than autoscaling
- Any pool whose pods should be placed by explicit toleration + selector, never by default scheduling

## Key Configuration Choices

- **Fixed node count** (`nodeCount: 2`, no `autoScale`) -- the pool holds exactly two nodes; scale it by editing the count.
- **Taint + label pairing** (`dedicated=system:NoSchedule` with `workload: system`) -- pods need a matching toleration to land here AND a nodeSelector to prefer it; the taint alone keeps everyone else out. The taint's `effect` must be spelled exactly as Kubernetes does (`NoSchedule`, `PreferNoSchedule`, `NoExecute`).
- **Cluster reference** (`cluster.valueFrom`) -- resolves the owning cluster's UUID from a `DigitalOceanKubernetesCluster` resource's outputs at deploy time.

## Placeholders to Replace

- `metadata.name` / `nodePoolName` -- your pool's name (unique within the cluster).
- `cluster.valueFrom.name` -- the name of your `DigitalOceanKubernetesCluster` resource (or replace the block with `value: <the cluster UUID>` for a cluster created outside Planton).
