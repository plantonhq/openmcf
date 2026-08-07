# Cost-Optimized Batch

An ephemeral Dataproc cluster tuned for batch Spark jobs: a small
on-demand base, a Spot secondary group with machine-type flexibility and
a standard/spot capacity mix, an attached autoscaling policy, and
aggressive idle auto-delete.

## When to use

- Batch ETL jobs that tolerate preemption
- Nightly or scheduled data processing pipelines
- Large-scale transformations where cost matters more than latency
- Workloads with Spark dynamic allocation enabled

## What to customize

- `projectId` / `region` — your project and region.
- `autoscalingPolicyUri` — the preset references a
  `GcpDataprocAutoscalingPolicy` named `batch-autoscaling` by
  `valueFrom`; deploy the policy first (see its `02-aggressive-batch`
  preset), or swap in a literal `value:` with the policy's full resource
  name (`projects/{project}/locations/{region}/autoscalingPolicies/{id}`).
- `secondaryWorkerConfig.numInstances` — size to your job's
  parallelism; updates in place.
- `instanceFlexibilityPolicy.instanceSelectionList` — machine types the
  group may draw from, in preference order (lower rank wins).

## Key configuration

- **1 master + 2 primary workers** — the minimal stable, HDFS-carrying base
- **10 Spot secondaries** — up to ~80% cheaper than on-demand for
  fault-tolerant burst capacity
- **Instance flexibility policy** — ranked machine-type fallbacks keep
  the group schedulable when one type's spot capacity dries up;
  `standardCapacityBase: 2` keeps a preemption-proof floor on demand
- **Attached autoscaling policy** — the shared policy resizes the
  worker groups on YARN memory pressure; the attachment updates in place
- **10-minute graceful decommission** — running tasks finish before
  scale-down removes their node
- **15-minute idle auto-delete** — the cluster self-destructs shortly
  after the last job completes

## Important notes

- Spot VMs can be reclaimed at any time; Spark's task retry plus the
  external shuffle service absorb preemptions for batch workloads.
- Secondary workers carry no HDFS data, so preemption never loses data.

## Related presets

- **01-dev-jupyter** — interactive cluster for development
- **02-ha-production** — high-availability cluster for production
- **04-spark-on-gke** — run Spark as pods on an existing GKE cluster
