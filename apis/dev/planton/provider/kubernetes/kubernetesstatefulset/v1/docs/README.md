# Kubernetes StatefulSet: Research Documentation

## Introduction

A StatefulSet is Kubernetes' answer to a category of software that a Deployment structurally cannot serve: systems whose replicas are not interchangeable. A database's replica 0 is not a spare copy of replica 1 — it has its own data, its own place in a replication topology, and clients that must be able to find *that specific member* again after a restart. The apps/v1 StatefulSet contract provides exactly three guarantees that make this possible:

1. **Stable identity** — pods are named `<name>-0` through `<name>-(N-1)` and keep those names across restarts and rescheduling
2. **Stable network identity** — a headless governing Service gives each pod a DNS name of the form `<name>-<ordinal>.<service>.<namespace>.svc.cluster.local` that always resolves to that member
3. **Stable storage** — each pod gets its own PersistentVolumeClaims stamped from templates, and reattaches to the same claims wherever it reschedules

Everything else in the StatefulSet API — ordered startup, reverse-ordinal updates, retention policies — exists to protect those three guarantees during change. Planton's **KubernetesStatefulSet** component models that surface completely, on top of the shared workload core (the same container and pod model every workload kind uses), so the stateful-specific surface stays small and everything learned on one workload kind transfers.

## What Is Modeled — and What Deliberately Is Not

### The shared workload core

`spec.container` (app + sidecars) and `spec.pod` use the shared `WorkloadContainer` and `WorkloadPod` messages. Probes, env (with secret materialization and cross-resource references), volume mounts carrying their own volume sources, lifecycle hooks, container and pod security contexts, scheduling, DNS, and termination behavior are documented once and behave identically across Deployment, StatefulSet, DaemonSet, Job, and CronJob.

### The StatefulSet-specific surface

- `volume_claim_templates` — per-replica storage templates (name, storage class, size, access modes, volume mode)
- `update_strategy` — RollingUpdate/OnDelete, `partition`, `max_unavailable`
- `pod_management_policy` — OrderedReady/Parallel
- `pvc_retention_policy` — Retain/Delete on deletion and on scale-down
- `ordinals.start` — alternate ordinal numbering
- `availability` — replicas, PodDisruptionBudget, `min_ready_seconds`, `revision_history_limit`

### Deliberately not modeled

- **No autoscaling.** A Deployment's availability block carries HPA settings; the StatefulSet's does not. Membership changes in stateful systems require application-aware orchestration — data re-sync, rebalancing, quorum reconfiguration — that a CPU-threshold controller knows nothing about. Scaling a stateful system is a deliberate operation, so it is a spec change, not a controller loop.
- **No ingress.** Exposure is composed through first-class ingress kinds referencing the exported Service handle. The workload spec stays exposure-free by design, and the resource graph shows every exposure decision explicitly.
- **No ServiceAccount or RBAC creation.** `pod.service_account` *references* a KubernetesServiceAccount; permissions come from KubernetesRbac grants targeting that identity.
- **No inline ConfigMaps.** Configuration is composed: KubernetesConfigMap resources are mounted or imported by name.
- **No `serviceName` knob.** Upstream requires naming the governing Service; the module derives it from the resource name. One less field, and the headless Service can never point at the wrong workload.

## Per-Field Guidance for the Tricky Surfaces

### Update strategy and partition canaries

StatefulSet rolling updates always proceed from the highest ordinal down, one pod at a time, waiting for each replacement to become ready. `partition` turns this into a canary mechanism with no extra machinery: during a rolling update, only pods with ordinal >= partition receive the new template.

With 5 replicas and `partition: 4`, pushing a new image updates only pod `-4`. Validate it against production traffic, then lower the partition step by step — each decrement rolls exactly the ordinals it newly covers — until 0 rolls the remainder. Raising the partition back stops the rollout where it stands; pods below the partition never saw the new template.

`max_unavailable` (values above 1 update multiple ordinals in parallel) requires the `MaxUnavailableStatefulSet` feature gate, which is not enabled on every cluster, and has little effect under `OrderedReady` pod management. Leave it empty unless you run large fleets on clusters where the gate is confirmed on.

`OnDelete` hands full control to the operator: pods are recreated with the new template only when someone deletes them. This is the right mode when each member's update needs coordinated steps (drain the member from the cluster, delete, verify rejoin) that no controller ordering can express.

### PVC retention trade-offs

The Kubernetes default — Retain on both events — makes storage strictly more durable than the workload: delete the StatefulSet and the PVCs (and data) remain; re-create it with the same name and the members re-adopt their volumes. The cost is manual cleanup and, until then, ongoing storage spend.

`pvc_retention_policy` lets each event trade that safety for hygiene independently:

- **`when_deleted: Delete`** is right for reproducible state — caches, test fixtures, anything re-derivable — and wrong for the only copy of a database. The volumes are removed with the workload; recovery is whatever your backups make possible.
- **`when_scaled: Delete`** changes scale-down semantics: the removed member's PVC goes with it, so a later scale-up starts that ordinal empty and it must re-sync from peers. Retain (default) means a scaled-down-then-up member rejoins with its data intact — usually the safer and faster path for replicated databases.

A practical composition: `when_deleted: Delete` + `when_scaled: Retain` keeps day-to-day scaling safe while ensuring full teardowns don't strand volumes.

### Pod management policy

`OrderedReady` (default) creates pods one at a time, each waiting for the previous to be ready, and deletes in reverse. Systems that bootstrap sequentially — a primary that must exist before replicas join — need this ordering.

`Parallel` launches and deletes all pods at once. Systems that run their own membership coordination (Cassandra, Kafka, Elasticsearch) gain nothing from sequential creation and bootstrap N× faster in parallel. Two things to keep straight:

- The policy affects only **scale operations** (creation and deletion), never updates — rolling updates remain one-at-a-time reverse-ordinal regardless
- Changing the policy on an existing StatefulSet requires replacing the object; choose it up front

### Quorum-aware PodDisruptionBudgets

For quorum systems, the PDB is not an availability nicety — it is the mechanism that makes node drains safe. A 3-member etcd or Kafka/KRaft cluster tolerates exactly one member down; two down means quorum loss, which for some systems is an outage and for others is data loss.

Set `min_available` to the quorum size (e.g. `"2"` for 3 members): the eviction API will then refuse to take a second member while one is already down, and cluster upgrades serialize around your quorum automatically. Prefer the absolute form over percentages for small member counts — `"66%"` of 3 rounds in ways that are easy to get wrong, while `"2"` says exactly what the system needs. Pair the PDB with hard anti-affinity on `kubernetes.io/hostname` so involuntary single-node failures also cost at most one member.

### Sizing termination for member handoff

`pod.termination_grace_period_seconds` deserves explicit attention on stateful systems: the default 30 seconds is calibrated for stateless HTTP servers, not for a database flushing its WAL or a broker handing off partition leadership. The grace clock starts before `pre_stop` hooks run, so size it to cover the hook plus the application's own worst-case clean shutdown. An OrderedReady rolling update multiplies any shortfall by the replica count.

## Storage Mechanics Worth Knowing

- Each template stamps one PVC per replica, named `<template>-<name>-<ordinal>`. The container mounts it with a `pvc` volume mount whose `claim_name` equals the template name — the StatefulSet controller handles the per-pod binding.
- Templates are effectively immutable on the object; growing volumes later requires a StorageClass with `allowVolumeExpansion` and editing the PVCs, not the template. Pick sizes with headroom.
- An empty `storage_class` uses the cluster default. Reference a class created by a KubernetesStorageClass resource to pin performance characteristics explicitly.
- `fs_group` in the pod security context is the standard fix for non-root databases on provisioned volumes (which arrive root-owned); `fs_group_change_policy: OnRootMismatch` avoids paying a recursive chown on every pod start of a large volume.

## Production Best Practices

1. **Set the PDB to quorum from day one.** Retro-fitting a disruption budget after the first bad drain is the classic stateful incident.
2. **Spread members across failure domains.** Hard anti-affinity by hostname for small clusters; topology spread constraints across zones for larger ones.
3. **Canary with partitions, not parallel fleets.** The partition mechanism gives ordinal-level control with zero extra infrastructure.
4. **Keep readiness strictly about serving.** For quorum members, readiness gates both client traffic and update progression — a readiness probe that checks cluster-wide health (not just this member) can deadlock a rolling update.
5. **State your retention policy explicitly.** Even when choosing the Retain default, writing it in the spec documents the data-lifecycle decision where the next operator will look.
6. **Give termination real time.** Size `termination_grace_period_seconds` to the measured clean-shutdown time of your slowest member.

## Conclusion

KubernetesStatefulSet models the full apps/v1 StatefulSet surface on the shared workload core, and its omissions are as deliberate as its fields: no autoscaling (membership is an application decision), no ingress (exposure composes through first-class kinds), no bundled identity (ServiceAccounts are referenced, not created). The result is a spec where the stateful-specific decisions — storage templates, update ordering, retention, quorum protection — are the visible, reviewable heart of the manifest.

## References

- [Kubernetes StatefulSets Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/)
- [StatefulSet Basics Tutorial](https://kubernetes.io/docs/tutorials/stateful-application/basic-stateful-set/)
- [PersistentVolumeClaim Retention](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#persistentvolumeclaim-retention)
- [Pod Disruption Budgets](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)
- [StatefulSet API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/stateful-set-v1/)
