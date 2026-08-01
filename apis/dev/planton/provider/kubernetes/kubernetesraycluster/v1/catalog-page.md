# Ray Cluster

The shared distributed-compute service that ML and data teams (and
their agents) connect to: notebooks and applications dial the client
port (10001), jobs submit through the dashboard port (8265, the Job
Submission API), and Ray's own processes coordinate through the GCS
port (6379) on the head. This component declares the RayCluster CR
the KubeRay operator reconciles into a head pod, worker groups, and
the head Service — a `KubernetesKubeRayOperator` whose watch scope
covers the namespace is the prerequisite.

## Highlights

- **The state truth, told up front** — the head's GCS holds the
  cluster's control state in memory; losing the head loses every job
  and actor unless GCS fault tolerance puts that state in a composed
  `KubernetesValkey`, after which a replaced head recovers instead of
  rebuilding empty.
- **Version/image lockstep by construction** — the default image
  derives from `rayVersion`, so the operator's command shaping and
  the running Ray can only diverge when you override the image
  deliberately (a mismatch fails at runtime, not at apply).
- **Secure by default** — token auth is the catalog default even
  though the operator's own default is an open cluster: anyone who
  can reach port 8265 on an unauthenticated Ray cluster runs
  arbitrary code; the bearer-token Secret is the exported credential
  handle.
- **Production heads stay unloaded** — the modules render
  `num-cpus: 0` so application work stays off the control plane (a
  task-loaded head starves the GCS); single-node labs flip one field.
- **Accelerators as limits, autoscaled from zero** — GPU groups
  declare `extraResourceLimits` (limits-only, the Kubernetes
  extended-resource contract), and Ray's own application-aware
  autoscaler materializes them only when a task asks.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
