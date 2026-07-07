# Spark on GKE

A Dataproc-on-GKE virtual cluster: Spark workloads run as Kubernetes
pods on an existing GKE cluster instead of dedicated Compute Engine VMs,
sharing the GKE cluster's capacity, autoscaling, and operational
tooling.

## When to use

- Platform teams standardizing all compute on Kubernetes
- Sharing one GKE cluster's capacity between Spark and other workloads
- Reusing GKE-native operations (node auto-provisioning, monitoring,
  policy) for Spark
- Avoiding a second VM fleet when a well-run GKE cluster already exists

## What to customize

- `projectId` / `region` — your project and region.
- `gkeClusterTarget` — the preset references a `GcpGkeCluster` named
  `platform-gke`; swap in your own reference or a literal `value:` with
  the cluster's fully qualified resource name
  (`projects/{project}/locations/{location}/clusters/{cluster}`).
- `nodePoolTarget` — the preset schedules every role onto one pool
  (`roles: [DEFAULT]`). Split `CONTROLLER`, `SPARK_DRIVER`, and
  `SPARK_EXECUTOR` onto dedicated pools when driver and executor
  capacity need independent sizing.
- `componentVersion.SPARK` — pin a version compatible with the target
  GKE version.

## Key configuration

- **Virtual arm** — `virtualClusterConfig` replaces `clusterConfig`;
  the two are mutually exclusive
- **Existing GKE composition** — the GKE cluster and node pool are
  first-class resources (`GcpGkeCluster`, `GcpGkeNodePool`) referenced
  by output, not created here
- **DEFAULT role mapping** — one pool hosts every Dataproc role; the
  pool may pre-exist or be created by Dataproc via `nodePoolConfig`
- **No user labels** — the Dataproc API rejects labels on virtual
  clusters; the spec validation catches this pre-deploy
- **Fully immutable** — any change to the virtual arm replaces the
  virtual cluster; the underlying GKE cluster and pools are untouched

## Related presets

- **01-dev-jupyter** — lightweight GCE-arm cluster for development
- **02-ha-production** — high-availability GCE-arm cluster
- **03-cost-optimized-batch** — Spot secondaries and autoscaling on the GCE arm
