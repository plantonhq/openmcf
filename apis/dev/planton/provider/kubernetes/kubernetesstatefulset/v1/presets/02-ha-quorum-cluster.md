# Highly Available Quorum Cluster

This preset deploys a three-member clustered stateful system — the shape of Kafka, etcd, ZooKeeper, or any consensus-based store. Each member gets a stable name (`my-quorum-cluster-0/-1/-2`), a stable per-replica DNS name through the governing Service, and its own 20Gi volume. Hard anti-affinity puts every member on a different node, and the disruption budget guarantees a node drain can never break quorum.

Members discover each other through the per-replica DNS names exported as the `podDnsTemplate` output: replica 0 is `my-quorum-cluster-0.my-quorum-cluster.<namespace>.svc.cluster.local`, and clustered clients build their member lists by substituting ordinals.

## When to Use

- Quorum-based systems: Kafka/KRaft, etcd, ZooKeeper, Consul
- Replicated databases where members are not interchangeable and must survive node maintenance

## Key Configuration Choices

- **3 replicas + PDB `minAvailable: "2"`** — with a 3-member quorum of 2, voluntary disruptions (drains, upgrades) can never take the cluster below quorum; Kubernetes simply refuses to evict the second pod
- **Required anti-affinity on `kubernetes.io/hostname`** — one member per node, so a single node failure costs at most one member; the `app: my-quorum-cluster` labels match the workload's own selector labels (keep them in sync with `metadata.name`)
- **`podManagementPolicy: Parallel`** — all three pods launch at once instead of one at a time; correct for systems that run their own membership coordination (Kafka, Cassandra) and materially faster to bootstrap. Keep the `OrderedReady` default for systems whose members must join one by one. This policy affects only scale-up/scale-down — rolling updates still proceed one ordinal at a time
- **`pvcRetentionPolicy: Retain/Retain`** — stated explicitly even though it is the Kubernetes default: member data survives both deletion and scale-down, so a scaled-down member that later returns rejoins with its data instead of re-syncing from peers
- **`terminationGracePeriodSeconds: 120`** — clustered members need time to hand off leadership and flush state; size this to your system's clean-shutdown worst case
- **Two ports (client + peer)** — both are published on the governing Service; member-to-member traffic uses the per-replica DNS names

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-cluster-image>` | Clustered system image | Your container registry |
| `<your-image-tag>` | Image tag or version | Your registry or CI/CD pipeline output |

## Related Presets

- **01-database** — Single-replica starting point
- **03-hardened-database** — Restricted-profile security hardening with a composed ServiceAccount identity
