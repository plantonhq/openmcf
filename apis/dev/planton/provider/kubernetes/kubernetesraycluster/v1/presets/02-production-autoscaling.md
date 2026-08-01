# Production autoscaling preset

The production shape, all four load-bearing choices made:

**An unloaded head.** `scheduleTasksOnHead: false` (the default)
renders `num-cpus: 0` so the scheduler keeps application work off the
head — a task-loaded head starves the GCS and destabilizes the whole
cluster. Size the head for coordination (upstream guidance starts at
4 CPU / 8Gi and scales with cluster size).

**Two worker groups under Ray's own autoscaler.** A CPU group and a
GPU group, each scaling between its bounds on the Ray scheduler's
demand — application-aware, unlike HPA. The GPU group starts at ZERO
and materializes only when a task asks for GPUs; accelerators are
declared through `extraResourceLimits` (limits-only, the Kubernetes
extended-resource contract) paired with the node selector and
toleration that land those pods on accelerator nodes.

**GCS fault tolerance on a composed store.** The head's control state
lives in a `KubernetesValkey` (referenced by name — the FK defaults
wire its endpoint and its `<name>-auth` Secret; the key is the ACL
username, `default` unless you declared users), so a replaced head
RECOVERS jobs, actors, and workers instead of rebuilding an empty
cluster. Deploy the Valkey first, in this same namespace — the
credential secretKeyRef cannot cross namespaces.

**Token auth by default.** Every API surface (dashboard, job
submission, client) requires the bearer token from the
operator-provisioned Secret named exactly after this resource (key
`auth_token`). An open Ray cluster executes arbitrary code for anyone
who reaches port 8265 — never disable auth outside a fenced lab.

Change first: worker group sizing and bounds; the Valkey reference
(`ray-state`) to your actual store's name.

See [02-production-autoscaling.yaml](./02-production-autoscaling.yaml)
for the manifest.
