# Kubernetes Deployment: Research Documentation

## Introduction

The apps/v1 Deployment is the most-used workload controller in Kubernetes: it manages a set of interchangeable pod replicas, replaces them gradually on update, and pairs naturally with a Service for load-balanced traffic and a HorizontalPodAutoscaler for elastic capacity. Virtually every stateless service in production Kubernetes runs as a Deployment.

The upstream surface, however, is enormous. A raw Deployment manifest embeds the entire core/v1 PodSpec — dozens of fields spanning containers, scheduling, security, DNS, and node-level plumbing — plus the Deployment's own rollout machinery. Planton's **KubernetesDeployment** component curates that surface into a typed, validated spec: the fields production services actually use are fully modeled with early validation; the fields that are imperative tooling, controller machinery, or immature upstream features are deliberately absent. This document records the curation rationale and the operational mechanics — rollouts, autoscaling, disruption budgets, hardening, scheduling — that the spec is designed around.

## Design Rationale: Curating apps/v1

### The full-surface baseline

Everything a stateless service needs from the upstream API is expressible: the complete container shape (image, command/args, ports, env with multiple value sources, resources, all three probes, volume mounts, lifecycle hooks, security context) applies identically to the app container, every sidecar, and every init container — which is exactly how the Kubernetes API treats them. Pod-level configuration covers identity, scheduling (node selectors, tolerations, node/pod affinity, topology spread), pod security context, termination, DNS, host aliases, host namespaces, priority, and runtime class. The Deployment's own controls — replicas, strategy, `minReadySeconds`, `revisionHistoryLimit`, `progressDeadlineSeconds`, `paused` — are all present.

### What is deliberately not modeled, and why

Three categories of upstream fields are intentionally absent:

**Imperative tooling, not declarative config.**

- `ephemeralContainers` — the mechanism behind `kubectl debug`. Ephemeral containers are injected into *running* pods and cannot even be specified at creation time; they have no place in a declarative manifest.
- `stdin` / `tty` on containers — support for interactive `kubectl attach` sessions.

**Controller- and operator-facing machinery.**

- `nodeName`, `schedulingGates`, `readinessGates` — levers for custom schedulers and controllers, not for application authors. A manifest that pins `nodeName` has opted out of scheduling entirely, which is almost never what a service owner means.
- `hostname` / `subdomain` — manual pod-identity fields. On a Deployment, replicas are interchangeable by definition; giving them individual DNS identities contradicts the kind's contract. (StatefulSet, where stable identity is the point, derives these fields from its own machinery — they are controller-owned on every workload kind.)

**Immature or niche upstream features.**

- `resizePolicy` — in-place pod vertical resize. Worth modeling once the feature is GA and supported by both IaC engine providers; until then it would be a field that silently does nothing on most clusters.
- `hostUsers`, `os`, `overhead`, `resourceClaims`, `preemptionPolicy`, `enableServiceLinks`, `setHostnameAsFQDN`, `shareProcessNamespace` — toggles without established production demand. Each unused field in a spec is cognitive cost paid by every reader; they can be added when a real use case arrives.

The principle throughout: **absence is a feature**. A curated spec is only trustworthy if users can assume that what is present works everywhere and what is absent was excluded for a reason.

### Composition boundaries

Two upstream-adjacent concerns are excluded not because they are immature but because they belong to other kinds:

- **Identity**: `serviceAccountName` exists only as a *reference* — a literal name or a link to a KubernetesServiceAccount resource. ServiceAccount creation, workload-identity annotations (GKE/EKS/AKS federation), and RBAC grants are owned by KubernetesServiceAccount and KubernetesRbac. A workload that creates its own identity makes permissions invisible side effects; a workload that references identity makes them auditable resources.
- **Exposure**: there is no ingress, gateway, or certificate configuration in the spec. The workload exports `service`, `kube_endpoint`, and `selector_labels`; exposure kinds reference those outputs. This keeps the workload portable across exposure technologies (Ingress controllers, Gateway API, meshes) and keeps every exposed hostname visible in the resource graph.

## Zero-Downtime Rollout Mechanics

A rolling update with defaults (25% surge, 25% unavailable) *will* drop requests. Zero-downtime requires four pieces working together, and the failure mode of omitting any one is instructive:

1. **`maxUnavailable: "0"`** — the Deployment never scales below the desired count. Without it, old pods are killed before their replacements serve.
2. **`maxSurge: "1"`** (or more) — the Deployment may exceed the desired count during the roll, giving it room to make progress. The spec enforces that surge and unavailable are not both zero — that configuration cannot make progress at all.
3. **A readiness probe** — this is what "available" means. Without a probe, a pod counts as ready the instant its containers start, so the controller happily replaces serving pods with pods that have not finished booting. With `maxUnavailable: 0` and a probe, the controller waits for each new pod to actually pass readiness before removing an old one.
4. **A `preStop` sleep** — the subtle half of the problem. Pod termination and endpoint removal are *concurrent*, not sequential: when a pod is deleted, the kubelet starts termination while the endpoint controller separately removes the pod from Service endpoints, and load balancers converge on that change asynchronously. A pod that exits immediately on SIGTERM closes connections that in-flight routing still points at. A short `preStop` sleep (the kubelet-native `sleep` action requires no sleep binary in the image — it works in distroless images) keeps the process serving during the propagation window. Size `terminationGracePeriodSeconds` to cover the hook *plus* the app's own drain time — the grace clock starts before the hook runs.

Two supporting controls round this out. `minReadySeconds` (10–30s) is a cheap flap detector: a new pod must stay ready that long before counting as available, which stops a crash-on-first-request regression from replacing the whole fleet before the first crash registers. `progressDeadlineSeconds` bounds how long a rollout may make no progress before the Deployment's conditions mark it failed — the signal deployment pipelines should watch instead of sleeping.

## HPA Metric Semantics

The two supported autoscaling targets have deliberately different shapes, mirroring how the autoscaling/v2 API is actually used:

- **CPU targets `Utilization`** — a percentage of the container's CPU *request*, averaged across replicas. `targetCpuUtilizationPercent: 70` means "add replicas when average CPU exceeds 70% of requests, remove them when it falls below." This works because CPU is elastic in both directions: load spreads across more replicas, per-replica CPU falls, scale-in has a real signal. It also means resource requests are the autoscaler's denominator — a service with no CPU request cannot CPU-autoscale, and an inflated request silently mutes the signal.
- **Memory targets `AverageValue`** — an absolute per-pod quantity (e.g. `1Gi`), not a percentage. Memory-as-percentage invites a trap the absolute form at least makes visible: **memory is a poor scale-in signal**. Most runtimes (JVM heaps, Go's runtime, caches, arenas) grow memory under load and never return it to the OS. After a scale-out, per-pod memory stays high even when traffic subsides, so a memory-based autoscaler ratchets up and never scales in — capacity costs without elasticity. Memory targets are appropriate for the narrow class of workloads whose working set genuinely tracks concurrent load; for everything else, scale on CPU and treat memory limits as the OOM guardrail.

The `replicas` field is the autoscaler's floor and `maxReplicas` its ceiling. Setting the floor at the PDB's `minAvailable` or above keeps the disruption budget satisfiable at minimum scale (see below).

## PDB Selector Discipline

A PodDisruptionBudget is only as correct as its selector. The module binds the PDB to the workload's own selector labels — the same immutable label set the Deployment's selector and the Service use — which closes two classic failure modes:

- **Over-broad selectors**: a hand-written PDB matching `app: payments` in a namespace where three tracks of the payments service run will count all three workloads' pods against one budget, blocking drains that should be safe or permitting disruptions that are not.
- **Selector drift**: pod labels that change per release (a version or build label in the selector) orphan the PDB — it matches zero pods, and depending on the direction of the mismatch either blocks every drain in the namespace or protects nothing.

The selector labels deliberately exclude the deployment-track version label for exactly this reason: selectors must be stable across releases, and the version changes on every pipeline run.

Budget arithmetic deserves equal care. The spec enforces exactly one of `minAvailable` / `maxUnavailable`. For an N-replica service, `minAvailable: N-1` (as a number, not a percentage) permits exactly one voluntary disruption at a time — the usual intent. Two configurations to avoid: `minAvailable` equal to the replica count makes the workload *undrainable* — cluster upgrades stall until a human intervenes; and a PDB on a single-replica Deployment with `minAvailable: 1` is the same trap in its most common costume. A PDB also protects nothing during *involuntary* disruptions (node crashes, OOM kills) — it is a contract with the eviction API, not an availability guarantee.

## Pod-Security Hardening

The container and pod security contexts cover the complete **restricted** profile of the Kubernetes Pod Security Standards, so a fully hardened workload is expressible without escape hatches. The checklist, and what each item actually buys:

| Setting | Why |
|---|---|
| `runAsNonRoot: true` | Refuses to start images that silently default to UID 0 — the failure is at pod start, loud and attributable, instead of a root process in production |
| Pinned `runAsUser` / `runAsGroup` (e.g. 10001) | Deterministic file ownership; no dependence on the image's `USER` directive |
| `readOnlyRootFilesystem: true` | The container cannot modify its own image at runtime; pair with a size-limited EmptyDir mounted at `/tmp` for legitimate scratch writes |
| `allowPrivilegeEscalation: false` | Blocks setuid binaries and file capabilities from granting more privilege than the parent process holds |
| `capabilities.drop: ["ALL"]` | The restricted baseline; add back only what the app demonstrably needs (`NET_BIND_SERVICE` for ports below 1024 — or better, bind 8080 and map via `servicePort`) |
| `seccompProfile.type: RuntimeDefault` | The runtime's default syscall filter; near-zero compatibility cost for ordinary services |
| `fsGroup` (pod level) | Volume ownership for non-root writers — the standard fix for permission-denied on mounted volumes; `fsGroupChangePolicy: OnRootMismatch` skips the recursive chown when ownership already matches, dramatically faster pod starts on large volumes |
| `automountServiceAccountToken: false` | An app that never calls the Kubernetes API should not carry credentials for it; this is the pod-level half of least privilege, complementing RBAC |

The layering rule: pod-level security context is the baseline every container inherits; container-level settings override per field. Sidecars and init containers accept the full security context, so hardening does not stop at the app container.

## Scheduling Patterns

For spreading replicas across failure domains, two mechanisms overlap and the choice matters:

- **Self-anti-affinity** (`podAntiAffinity` on the workload's own selector labels across `kubernetes.io/hostname`) is *binary*: a node either has a matching pod or it does not. In `required` form it caps replicas at the number of nodes — the eleventh replica of a ten-node cluster is unschedulable forever, and rollouts with surge can deadlock because the surge pod has nowhere to go. The `preferred` form degrades gracefully but offers no bound on how uneven placement may get.
- **Topology spread constraints** bound *skew*: `maxSkew: 1` across `topology.kubernetes.io/zone` means no zone may run more than one replica above the least-loaded zone, at any replica count. `whenUnsatisfiable` chooses hard (`DoNotSchedule`) or soft (`ScheduleAnyway`) enforcement per constraint. When `matchLabels` is omitted, the module defaults to the workload's own selector labels — self-spreading, the overwhelmingly common intent.

The practical guidance: topology spread for replica distribution (it scales past the node count and bounds rather than forbids), reserved anti-affinity for genuine exclusion ("never co-locate with the batch tier"), and `required` pod anti-affinity only with full awareness that it can deadlock rollouts. The remaining mechanisms compose rather than compete: `nodeSelector`/node affinity choose candidate nodes, tolerations unlock tainted ones, and spread constraints distribute replicas across whatever set results.

## Implementation Landscape

Both IaC implementations — Pulumi (Go) and Terraform (HCL) — have feature parity and identical logic:

- Namespace resolution (literal or reference), with optional creation
- Satellite Secrets created **before** the Deployment: the env Secret (literal secret env values collected across app, sidecars, and init containers into one Opaque Secret) and the docker-registry pull Secret. Ordering matters — a pod that starts before its env Secret exists crashes
- The Deployment, with the pod volume list derived from the union of all containers' volume mounts (two containers sharing an EmptyDir simply declare the same mount name and source)
- The ClusterIP Service, created only when the app container declares ports — port-less workers get no Service at all
- HPA and PDB when enabled, both bound to the workload's selector labels
- Outputs: namespace, deployment name, Service name, selector labels, port-forward command, and in-cluster endpoint

## Conclusion

KubernetesDeployment is a curated, opinionated rendering of the most-used workload API in Kubernetes: full coverage of the production surface, deliberate absence of imperative and immature fields, and hard composition boundaries around identity and exposure. The operational patterns it is designed for — zero-downtime rollouts, CPU-led autoscaling, selector-disciplined disruption budgets, restricted-profile hardening, skew-bounded spreading — are each expressible in a few lines because the spec was shaped around them.

## References

- [Kubernetes Deployments Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
- [Horizontal Pod Autoscaling](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
- [Pod Disruption Budgets](https://kubernetes.io/docs/concepts/workloads/pods/disruptions/)
- [Pod Security Standards](https://kubernetes.io/docs/concepts/security/pod-security-standards/)
- [Pod Topology Spread Constraints](https://kubernetes.io/docs/concepts/scheduling-eviction/topology-spread-constraints/)
- [Container Lifecycle Hooks](https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/)
- [Deployment API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/deployment-v1/)
