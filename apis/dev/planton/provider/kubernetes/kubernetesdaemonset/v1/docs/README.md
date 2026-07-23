# Kubernetes DaemonSet: Research Documentation

## Introduction

A DaemonSet inverts the usual Kubernetes question. Every other workload kind asks "how many copies should run, and let the scheduler decide where"; a DaemonSet asks "which nodes need this, and run exactly one copy on each". The controller reconciles against node membership itself: a node joining the cluster gets the pod automatically, a node leaving takes it away, and there is no replica count anywhere because the node set IS the replica count.

That inversion exists because a whole category of software is node infrastructure rather than application: log shippers that read the node's files, metrics agents that observe the node's kernel, storage daemons that manage the node's disks, network plugins that program the node's routing. These programs need per-node presence, frequently need pieces of the host that ordinary pods are isolated from, and are meaningless to "scale" independently of the node count.

Planton's **KubernetesDaemonSet** component models the apps/v1 DaemonSet surface on the shared workload core — the same container and pod model every workload kind uses — so the DaemonSet-specific spec is strikingly small: an update strategy, two rollout knobs, and the shared pod's scheduling and host-access fields doing the node-targeting work.

## What Is Modeled — and What Deliberately Is Not

### The shared workload core

`spec.container` (app + sidecars) and `spec.pod` use the shared `WorkloadContainer` and `WorkloadPod` messages: probes, env (with secret materialization and cross-resource references), volume mounts carrying their own sources, lifecycle hooks, security contexts at both levels, scheduling, DNS, and termination behavior — identical across all five workload kinds.

### The DaemonSet-specific surface

- `update_strategy` — RollingUpdate/OnDelete with `max_unavailable` and `max_surge`
- `min_ready_seconds` and `revision_history_limit` — at spec level, since there is no availability block

### Deliberately not modeled

- **No replica count and no autoscaling.** Structural: the node set determines the pod count. "Scaling" a DaemonSet means changing which nodes match under `pod.scheduling`.
- **No Service and no ingress.** A DaemonSet's pods are node-local by intent; a load-balanced Service across them is almost always a design smell (which node's agent answered?). Clients reach a specific node's agent via `host_port` or `host_network`. When something genuinely needs to consume "whichever instance is on my node", the platform-native answer is the exported `selector_labels` output for NetworkPolicies and observability tooling.
- **No PodDisruptionBudget.** DaemonSet pods are not evicted by drains the way replicated workloads are (the eviction API treats them specially, and `kubectl drain` requires `--ignore-daemonsets`); a PDB would be dead configuration.
- **No ServiceAccount or RBAC creation.** `pod.service_account` references a KubernetesServiceAccount; cluster-read permissions that many agents need come from KubernetesRbac grants targeting that identity — visible, auditable resources instead of workload side effects.
- **No inline ConfigMaps.** Configuration composes from first-class KubernetesConfigMap resources mounted or imported by name.

## Per-Field Guidance for the Tricky Surfaces

### Node targeting via `pod.scheduling`

Which nodes run the agent is entirely a scheduling question, answered by three composable mechanisms:

- **`node_selector`** — the simple hard filter: `kubernetes.io/os: linux`, a node-pool label, a hardware marker. Right tool when the target set is "nodes with label X".
- **`node_affinity`** — the expressive form: operators (`In`, `NotIn`, `Exists`, `Gt`/`Lt`) and OR-ed requirement terms. Use it for "GPU nodes of either family" or "everything except the batch pool".
- **`tolerations`** — the unlock for tainted nodes. A toleration never attracts pods; it only permits scheduling where a taint would otherwise repel it. Two patterns matter for DaemonSets:
  - **Control-plane coverage.** Control-plane nodes carry `node-role.kubernetes.io/control-plane: NoSchedule`. Observability and security agents almost always want `operator: Exists` on that key — a monitoring pipeline that skips the control plane has a blind spot exactly where incidents are most interesting.
  - **Dedicated pools.** Nodes taints like `dedicated=gpu:NoSchedule` keep ordinary workloads off; the DaemonSets that serve those nodes (device plugins, specialized monitors) tolerate the taint and select the pool.

The DaemonSet controller itself adds tolerations for node-pressure taints (memory-pressure, disk-pressure, not-ready), so node agents keep running through conditions that evict ordinary pods — behavior you get for free and should not replicate in the spec.

### Host access patterns

Node agents need slices of the host, and the spec makes each slice a separate, explicit grant:

- **Node filesystem** — HostPath volume mounts on the container: `/var/log` and `/var/log/pods` for log collectors, `/proc` and `/sys` (mounted under `/host/...`, read-only) for node monitors, the container runtime socket for container-aware agents. Prefer `read_only: true` everywhere the agent only observes; a writable runtime-socket mount is effectively node-root.
- **Host network** — `pod.host_network: true` puts the agent in the node's network namespace: it sees real interfaces and binds node ports directly. Pair with `dns_policy: ClusterFirstWithHostNet` whenever the agent also talks to cluster services, or its DNS silently degrades to the node's resolver.
- **Host PID** — `pod.host_pid: true` exposes the node's process table, for process-visibility and security agents.
- **Node-IP exposure** — per-port `host_port` publishes an endpoint at `<node-ip>:<port>` on every node without host networking. This is the scrape-endpoint pattern: each node's agent is individually addressable, which is exactly what per-node metrics collection wants.
- **Capabilities over privilege** — the container `security_context` can add targeted capabilities (`SYS_PTRACE` for process inspection, `NET_ADMIN` for network programming) after `drop: ["ALL"]`. Reserve `privileged: true` for the few agents that genuinely manage devices or kernel state; it is root on the node and disables most isolation.

The grants compose independently — a log collector needs only HostPath; a node monitor adds host namespaces and a host port; a CNI component adds capabilities. Reviewing a DaemonSet manifest means reading exactly which slices of the node it was given.

### The surge-versus-hostPort interaction

`update_strategy` defaults to RollingUpdate with `max_unavailable: 1`: one node's agent down at a time, the safe default everywhere. Two levers modify it:

- **`max_unavailable` above 1 (or a percentage)** speeds rollouts across large fleets when a per-node agent gap is tolerable.
- **`max_surge`** flips the trade: the new pod starts on a node BEFORE the old one stops, so there is never a moment without a running agent — the right choice for log shippers and security agents where a gap means lost data. The cost is doubled per-node resource usage during the update, and one hard constraint: **surge requires the old and new pods to coexist on the same node, so it cannot be combined with exclusive host ports.** Two pods cannot bind `<node-ip>:9100` at once — a surging update of a hostPort agent deadlocks with the new pod unschedulable. Host-port agents therefore live with the brief per-node gap (softened by `min_ready_seconds` as a flap detector); gapless agents must expose themselves some other way (host network with SO_REUSEPORT-aware binding, or no node-local endpoint at all).

`OnDelete` remains the escape hatch for agents whose restart must be coordinated with node maintenance: the controller updates nothing until someone deletes a pod, giving operators per-node control of the rollout.

## Production Best Practices

1. **Budget resources for the fleet, not the pod.** Requests and limits multiply by the node count; an unbounded collector on 200 nodes is 200 unbounded collectors competing with workloads.
2. **Tolerate the control plane deliberately.** Decide per agent whether control-plane coverage is wanted, and write the toleration (or its absence) as a reviewed choice.
3. **Mount host paths read-only unless the agent provably writes.** And treat writable mounts of runtime sockets or kernel interfaces as the privilege grants they are.
4. **Prefer capabilities to `privileged`.** Most agents document the exact capabilities they need; granting them individually keeps the blast radius legible.
5. **Use `min_ready_seconds` as a rollout flap detector.** 10–30 seconds catches a crash-looping agent build before the rollout replaces it on every node in the fleet.
6. **Keep probes node-local.** An agent's liveness should reflect the agent's own health; probing through to its remote sink turns a sink outage into a fleet-wide restart storm.

## Conclusion

KubernetesDaemonSet is the smallest of the workload kinds because the shared core already carries almost everything a node agent needs — the kind itself contributes the one-pod-per-node reconciliation model and its update semantics. Its omissions are structural, not gaps: no replicas (nodes are the replicas), no Service (agents are node-local), no PDB (drains treat daemons specially), no bundled identity or configuration (both compose from first-class kinds). What remains in a manifest is precisely the interesting part: which nodes, which slices of the host, and how updates roll across the fleet.

## References

- [Kubernetes DaemonSet Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
- [Perform a Rolling Update on a DaemonSet](https://kubernetes.io/docs/tasks/manage-daemon/update-daemon-set/)
- [Taints and Tolerations](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [DaemonSet API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/daemon-set-v1/)
