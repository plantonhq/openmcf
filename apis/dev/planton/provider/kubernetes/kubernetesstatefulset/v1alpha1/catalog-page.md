# Kubernetes StatefulSet

Deploys a stateful application to a target Kubernetes cluster as an apps/v1 StatefulSet through a single declarative manifest: stable pod identity (`<name>-0`, `<name>-1`, ...), stable per-replica DNS through a derived headless governing Service, per-replica PersistentVolumeClaims from volume claim templates, quorum-aware disruption budgets, partition-based canary updates, and PVC retention policies. The IaC module derives the governing Service, selector labels, and label merging automatically.

## What Gets Created

When you deploy a KubernetesStatefulSet resource, Planton provisions:

- **StatefulSet** — the apps/v1 workload with the app container, sidecars, init containers, scheduling, and security configuration from the spec
- **Headless governing Service** — shares the workload's name; provides each replica's stable DNS name and load-balanced client access on the declared ports
- **PersistentVolumeClaims** — one per replica per entry in `volumeClaimTemplates`, named `<template>-<name>-<ordinal>`
- **PodDisruptionBudget** — when `availability.podDisruptionBudget.enabled` is true; shares the workload's name
- **Env Secret** — `<name>-env-secrets`, only when the spec carries literal secret env values
- **Namespace** — only when `createNamespace` is true
- **Labels** — standard Planton tracking labels merged with selector labels and any user-provided pod labels

External exposure is never created here: attach a first-class exposure kind (e.g. a Gateway API route like KubernetesHttpRoute) that references this workload's exported `service` output.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, a `KubernetesNamespace` resource referenced from `spec.namespace`, or `createNamespace: true`
- **A StorageClass** able to serve the volume claim templates (the cluster default is used when `storageClass` is unset)
- **A container image** for the stateful application, reachable from the cluster

## Quick Start

Create a file `statefulset.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStatefulSet
metadata:
  name: my-db
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStatefulSet.my-db
spec:
  namespace:
    value: databases
  container:
    app:
      image:
        repo: postgres
        tag: "16.3"
      ports:
        - name: db
          containerPort: 5432
          servicePort: 5432
      volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
          pvc:
            claimName: data
  volumeClaimTemplates:
    - name: data
      size: 10Gi
      accessModes:
        - ReadWriteOnce
```

Deploy:

```shell
planton apply -f statefulset.yaml
```

This creates a single-replica StatefulSet: pod `my-db-0` with its own 10Gi claim `data-my-db-0`, addressable in-cluster at `my-db-0.my-db.databases.svc.cluster.local`.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec.namespace` | `StringValueOrRef` | Target namespace — a literal name (`{ value: databases }`) or a reference to a `KubernetesNamespace` resource |
| `spec.container.app` | `WorkloadContainer` | The main application container. `image.repo` and `image.tag` are required; this path is the deploy-pipeline image injection point |

### Container Fields (`spec.container`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `app.ports` | `list` | `[]` | Container ports; `servicePort` publishes the port on the governing Service, named ports are referenced by probes |
| `app.env` | `ContainerEnv` | — | Plain `variables`, `secrets` (literal values materialize a managed Secret), and bulk `envFrom` imports |
| `app.livenessProbe` / `readinessProbe` / `startupProbe` | `Probe` | — | Health checks; readiness gates client traffic and update progression |
| `app.volumeMounts` | `list` | `[]` | Mounts carrying their own volume source (`pvc`, `configMap`, `secret`, `emptyDir`, `hostPath`); a `pvc.claimName` matching a volume claim template name mounts the per-replica storage |
| `app.lifecycle` | `WorkloadContainerLifecycle` | — | `postStart` / `preStop` hooks; `preStop` covers member handoff |
| `app.securityContext` | `WorkloadContainerSecurityContext` | — | Container-level hardening (non-root, read-only root filesystem, capabilities, seccomp) |
| `sidecars` | `list<WorkloadContainer>` | `[]` | Named sidecar containers with the full container surface |

### Pod Fields (`spec.pod`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `serviceAccount` | `StringValueOrRef` | namespace default | Identity pods run as — literal name or reference to a `KubernetesServiceAccount` resource |
| `initContainers` | `list<WorkloadContainer>` | `[]` | Run to completion, in order, before app containers start |
| `scheduling` | `WorkloadScheduling` | — | Node selection, tolerations, affinity/anti-affinity, topology spread |
| `securityContext` | `WorkloadPodSecurityContext` | — | Pod-level baseline including `fsGroup` for volume ownership |
| `terminationGracePeriodSeconds` | `int64` | `30` | Time between SIGTERM and SIGKILL; size to a clean member handoff |
| `dnsPolicy` / `dnsConfig` / `hostAliases` | — | — | DNS resolution controls |

### StatefulSet Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.createNamespace` | `bool` | `false` | Create the namespace if it does not exist |
| `spec.availability.replicas` | `int32` | `1` | Desired member count; no autoscaling by design |
| `spec.availability.podDisruptionBudget` | `object` | disabled | `minAvailable` or `maxUnavailable` (exactly one); set `minAvailable` to the quorum size for consensus systems |
| `spec.availability.minReadySeconds` | `int32` | `0` | Seconds a new pod must stay ready before the rollout proceeds |
| `spec.availability.revisionHistoryLimit` | `int32` | `10` | Retained ControllerRevisions for rollback |
| `spec.volumeClaimTemplates` | `list` | `[]` | Per-replica PVC templates: `name` (the claim name mounts reference), `size` (required), `storageClass`, `accessModes` (default `ReadWriteOnce`), `volumeMode` |
| `spec.updateStrategy.type` | `string` | `RollingUpdate` | `RollingUpdate` (reverse ordinal, one at a time) or `OnDelete` (operator-driven) |
| `spec.updateStrategy.partition` | `int32` | `0` | Canary-by-ordinal: only pods with ordinal >= partition update |
| `spec.updateStrategy.maxUnavailable` | `string` | `1` | Parallel ordinal updates; requires the `MaxUnavailableStatefulSet` feature gate |
| `spec.podManagementPolicy` | `string` | `OrderedReady` | `OrderedReady` (sequential scale operations) or `Parallel` (all at once) |
| `spec.pvcRetentionPolicy.whenDeleted` | `string` | `Retain` | `Retain` keeps PVCs when the StatefulSet is deleted; `Delete` removes them |
| `spec.pvcRetentionPolicy.whenScaled` | `string` | `Retain` | PVC fate for members removed by scale-down |
| `spec.ordinals.start` | `int32` | `0` | First replica's ordinal |

## Examples

### Three-Member Quorum Cluster

Anti-affinity across nodes, a quorum-guarding PDB, and parallel bootstrap:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStatefulSet
metadata:
  name: kraft-broker
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesStatefulSet.kraft-broker
spec:
  namespace:
    value: streaming
  container:
    app:
      image:
        repo: my-registry/kafka
        tag: "3.7.0"
      ports:
        - name: client
          containerPort: 9092
          servicePort: 9092
        - name: peer
          containerPort: 9093
          servicePort: 9093
      volumeMounts:
        - name: data
          mountPath: /var/lib/kafka
          pvc:
            claimName: data
  pod:
    scheduling:
      podAntiAffinity:
        required:
          - matchLabels:
              app: kraft-broker
            topologyKey: kubernetes.io/hostname
    terminationGracePeriodSeconds: 120
  availability:
    replicas: 3
    podDisruptionBudget:
      enabled: true
      minAvailable: "2"
  podManagementPolicy: Parallel
  volumeClaimTemplates:
    - name: data
      size: 50Gi
      accessModes:
        - ReadWriteOnce
```

### Canary Update by Partition

Only the highest ordinal receives new templates until the partition is lowered:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStatefulSet
metadata:
  name: search-cluster
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesStatefulSet.search-cluster
spec:
  namespace:
    value: search
  container:
    app:
      image:
        repo: my-registry/search
        tag: "2.4.1"
      ports:
        - name: http
          containerPort: 9200
          servicePort: 9200
  availability:
    replicas: 5
  updateStrategy:
    type: RollingUpdate
    partition: 4
```

### Ephemeral State with Automatic Cleanup

A reproducible cache whose volumes are removed with the workload:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesStatefulSet
metadata:
  name: render-cache
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesStatefulSet.render-cache
spec:
  namespace:
    value: rendering
  container:
    app:
      image:
        repo: my-registry/render-cache
        tag: "1.2.0"
      ports:
        - name: cache
          containerPort: 6379
          servicePort: 6379
      volumeMounts:
        - name: cache-data
          mountPath: /data
          pvc:
            claimName: cache-data
  pvcRetentionPolicy:
    whenDeleted: Delete
    whenScaled: Delete
  volumeClaimTemplates:
    - name: cache-data
      size: 5Gi
      accessModes:
        - ReadWriteOnce
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | Namespace the workload was deployed into |
| `statefulSetName` | `string` | Name of the StatefulSet object in the cluster |
| `service` | `string` | The headless governing Service — stable per-pod DNS and client access |
| `selectorLabels` | `string` | Pod selector labels as `k=v,k=v`, for NetworkPolicies and `kubectl -l` |
| `portForwardCommand` | `string` | Ready-to-run `kubectl port-forward` command for local access |
| `kubeEndpoint` | `string` | In-cluster DNS endpoint of the Service |
| `podDnsTemplate` | `string` | Template for each replica's stable DNS name — substitute the ordinal to address a specific member |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — the stateless counterpart; use it when replicas are interchangeable
- [KubernetesDaemonSet](/docs/catalog/kubernetes/kubernetesdaemonset) — one pod per node, for node agents
- [KubernetesServiceAccount](/docs/catalog/kubernetes/kubernetesserviceaccount) — the identity pods run as; reference it from `spec.pod.serviceAccount`
- [KubernetesConfigMap](/docs/catalog/kubernetes/kubernetesconfigmap) — configuration consumed via env or volume mounts
- [KubernetesSecret](/docs/catalog/kubernetes/kubernetessecret) — confidential values referenced from container env
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace`
- [KubernetesStorageClass](/docs/catalog/kubernetes/storage-class) — pin volume performance characteristics for claim templates
