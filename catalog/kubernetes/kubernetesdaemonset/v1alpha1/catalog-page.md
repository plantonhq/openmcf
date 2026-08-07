# Kubernetes DaemonSet

Deploys a node agent to a target Kubernetes cluster as an apps/v1 DaemonSet through a single declarative manifest: exactly one pod on every node matching the pod's scheduling rules, with node targeting via selectors/tolerations/affinity, host access via HostPath mounts, host namespaces, and host ports, and rolling updates bounded by `maxUnavailable`/`maxSurge`. The IaC module handles selector labels and label merging automatically.

## What Gets Created

When you deploy a KubernetesDaemonSet resource, Planton provisions:

- **DaemonSet** — the apps/v1 workload with the agent container, sidecars, init containers, scheduling, host access, and security configuration from the spec
- **Env Secret** — `<name>-env-secrets`, only when the spec carries literal secret env values
- **Namespace** — only when `createNamespace` is true
- **Labels** — standard Planton tracking labels merged with selector labels and any user-provided pod labels

Nothing else: DaemonSets have no Service, ingress, HPA, or PDB — clients reach an agent on its node via `hostPort` or `hostNetwork`, and identity/permissions compose from `KubernetesServiceAccount` and `KubernetesRbac` resources.

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, a `KubernetesNamespace` resource referenced from `spec.namespace`, or `createNamespace: true`
- **A container image** for the agent, reachable from the cluster
- **Tolerations planned** for any tainted nodes (e.g. control-plane) the agent must cover

## Quick Start

Create a file `daemonset.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDaemonSet
metadata:
  name: log-collector
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesDaemonSet.log-collector
spec:
  namespace:
    value: observability
  container:
    app:
      image:
        repo: fluent/fluent-bit
        tag: "3.0.7"
      volumeMounts:
        - name: varlog
          mountPath: /var/log
          readOnly: true
          hostPath:
            path: /var/log
            type: Directory
```

Deploy:

```shell
planton apply -f daemonset.yaml
```

This runs one collector pod on every schedulable node, reading node logs through a read-only HostPath mount.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec.namespace` | `StringValueOrRef` | Target namespace — a literal name (`{ value: observability }`) or a reference to a `KubernetesNamespace` resource |
| `spec.container.app` | `WorkloadContainer` | The main agent container; `image.repo` and `image.tag` are required |

### Container Fields (`spec.container`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `app.ports` | `list` | `[]` | Container ports; `hostPort` exposes an endpoint at `<node-ip>:<port>` on every node (DaemonSets have no Service, so `servicePort` does not apply) |
| `app.env` | `ContainerEnv` | — | Plain `variables`, `secrets` (literal values materialize a managed Secret), and bulk `envFrom` imports |
| `app.livenessProbe` / `readinessProbe` / `startupProbe` | `Probe` | — | Health checks; keep them node-local so a sink outage does not restart the fleet |
| `app.volumeMounts` | `list` | `[]` | Mounts carrying their own volume source; `hostPath` sources (node logs, `/proc`, runtime sockets) are the classic DaemonSet pattern — prefer `readOnly: true` |
| `app.securityContext` | `WorkloadContainerSecurityContext` | — | Container hardening; add targeted capabilities after `drop: ["ALL"]` instead of `privileged` where possible |
| `sidecars` | `list<WorkloadContainer>` | `[]` | Named sidecar containers with the full container surface |

### Pod Fields (`spec.pod`)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `serviceAccount` | `StringValueOrRef` | namespace default | Identity pods run as — literal name or reference to a `KubernetesServiceAccount`; API permissions come from `KubernetesRbac` grants |
| `scheduling.nodeSelector` / `nodeAffinity` | — | all nodes | WHICH nodes run the agent |
| `scheduling.tolerations` | `list` | `[]` | Unlock tainted nodes — e.g. `node-role.kubernetes.io/control-plane` with `operator: Exists` for control-plane coverage |
| `hostNetwork` | `bool` | `false` | Run in the node's network namespace; pair with `dnsPolicy: ClusterFirstWithHostNet` if the agent resolves cluster services |
| `hostPid` | `bool` | `false` | Share the node's PID namespace (process-visibility agents) |
| `dnsPolicy` / `dnsConfig` / `hostAliases` | — | — | DNS resolution controls |
| `securityContext` | `WorkloadPodSecurityContext` | — | Pod-level security baseline |
| `terminationGracePeriodSeconds` | `int64` | `30` | Time between SIGTERM and SIGKILL |

### DaemonSet Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.createNamespace` | `bool` | `false` | Create the namespace if it does not exist |
| `spec.updateStrategy.type` | `string` | `RollingUpdate` | `RollingUpdate` (node by node) or `OnDelete` (per-node manual control) |
| `spec.updateStrategy.maxUnavailable` | `string` | `1` | Nodes whose agent may be down during an update — absolute or percentage |
| `spec.updateStrategy.maxSurge` | `string` | `0` | Nodes that may run old AND new pods during an update — gapless rollouts, but incompatible with exclusive host ports |
| `spec.minReadySeconds` | `int32` | `0` | Seconds a new pod must stay ready before the rollout proceeds — a fleet-wide flap detector |
| `spec.revisionHistoryLimit` | `int32` | `10` | Retained ControllerRevisions for rollback |

## Examples

### Node Monitor with Host Namespaces

Host network and PID, read-only kernel interfaces, a host-port scrape endpoint, and targeted capabilities:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDaemonSet
metadata:
  name: node-monitor
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDaemonSet.node-monitor
spec:
  namespace:
    value: observability
  container:
    app:
      image:
        repo: prom/node-exporter
        tag: "v1.8.1"
      ports:
        - name: metrics
          containerPort: 9100
          hostPort: 9100
      volumeMounts:
        - name: proc
          mountPath: /host/proc
          readOnly: true
          hostPath:
            path: /proc
            type: Directory
        - name: sys
          mountPath: /host/sys
          readOnly: true
          hostPath:
            path: /sys
            type: Directory
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
          add: ["SYS_PTRACE"]
  pod:
    hostNetwork: true
    hostPid: true
    dnsPolicy: ClusterFirstWithHostNet
    scheduling:
      tolerations:
        - key: node-role.kubernetes.io/control-plane
          operator: Exists
          effect: NoSchedule
```

### Gapless Log Shipper Update

Surge-based rollout so no node ever lacks a running collector (note: no host ports, which surge cannot coexist with):

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDaemonSet
metadata:
  name: log-shipper
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDaemonSet.log-shipper
spec:
  namespace:
    value: observability
  container:
    app:
      image:
        repo: fluent/fluent-bit
        tag: "3.0.7"
      volumeMounts:
        - name: pod-logs
          mountPath: /var/log/pods
          readOnly: true
          hostPath:
            path: /var/log/pods
            type: DirectoryOrCreate
  updateStrategy:
    type: RollingUpdate
    maxUnavailable: "0"
    maxSurge: "1"
  minReadySeconds: 15
```

### Dedicated-Pool Agent

An agent that runs only on a tainted GPU pool:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesDaemonSet
metadata:
  name: gpu-device-agent
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDaemonSet.gpu-device-agent
spec:
  namespace:
    value: kube-system-ext
  container:
    app:
      image:
        repo: my-registry/gpu-agent
        tag: "2.1.0"
  pod:
    scheduling:
      nodeSelector:
        node-pool: gpu
      tolerations:
        - key: dedicated
          operator: Equal
          value: gpu
          effect: NoSchedule
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | Namespace the workload was deployed into |
| `daemonSetName` | `string` | Name of the DaemonSet object in the cluster |
| `selectorLabels` | `string` | Pod selector labels as `k=v,k=v`, for NetworkPolicy podSelectors, `kubectl -l`, and pod-affinity terms |

## Related Components

- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — N interchangeable replicas behind a Service
- [KubernetesStatefulSet](/docs/catalog/kubernetes/kubernetesstatefulset) — stateful members with identity and per-replica storage
- [KubernetesServiceAccount](/docs/catalog/kubernetes/kubernetesserviceaccount) — the identity pods run as; reference it from `spec.pod.serviceAccount`
- [KubernetesRbac](/docs/catalog/kubernetes/kubernetesrbac) — API permissions for agents, granted to the referenced identity
- [KubernetesConfigMap](/docs/catalog/kubernetes/kubernetesconfigmap) — agent configuration consumed via env or volume mounts
- [KubernetesNamespace](/docs/catalog/kubernetes/kubernetesnamespace) — provides the target namespace; reference it from `spec.namespace`
