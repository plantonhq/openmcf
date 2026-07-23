---
title: "Deployment"
description: "Deployment deployment documentation"
icon: "package"
order: 100
componentName: "kubernetesdeployment"
---

# Kubernetes Deployment

Deploys a long-running, stateless application to a target cluster as an apps/v1 Deployment fronted by a ClusterIP Service, with optional horizontal pod autoscaling, rollout strategy control, and a pod disruption budget. The full container surface is available — sidecars, init containers, probes, lifecycle hooks, security contexts, volume mounts — along with pod scheduling and hardening. Identity and external exposure are composed from other components rather than embedded.

## What Gets Created

When you deploy a KubernetesDeployment resource, Planton provisions:

- **Deployment** — an apps/v1 Deployment with the specified containers, probes, volume mounts, scheduling, security contexts, and rollout strategy
- **Service** — a ClusterIP Service mapping each `servicePort` to its `containerPort`, created only when `container.app.ports` is non-empty; port-less workers get no Service
- **HorizontalPodAutoscaler** — created only when `availability.horizontalPodAutoscaling.enabled` is `true`, scaling between `availability.replicas` and `maxReplicas`
- **PodDisruptionBudget** — created only when `availability.podDisruptionBudget.enabled` is `true`, bound to the workload's selector labels
- **Env Secret** — one Opaque Secret materializing literal values from `env.secrets` across the app container, sidecars, and init containers; entries using `secretRef` are wired directly to existing Secrets instead
- **Image Pull Secret** — a `kubernetes.io/dockerconfigjson` Secret, created only when registry credentials are supplied
- **Namespace** — created only when `createNamespace` is `true`

Not created by design: ServiceAccounts and RBAC (reference a `KubernetesServiceAccount` from `spec.pod.serviceAccount`; grant permissions with `KubernetesRbac`), and ingress of any form (attach exposure components to the exported `service` / `kubeEndpoint` outputs).

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that already exists, a `KubernetesNamespace` resource referenced from `spec.namespace`, or `createNamespace: true`
- **A container image** accessible from the cluster (public registry, or private with pull credentials)
- **Metrics server** installed in the cluster if enabling horizontal pod autoscaling

## Quick Start

Create a file `deployment.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: my-app
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesDeployment.my-app
spec:
  namespace:
    value: my-namespace
  container:
    app:
      image:
        repo: nginx
        tag: "1.27"
      ports:
        - name: http
          containerPort: 80
          servicePort: 80
      readinessProbe:
        httpGet:
          path: /
          portNumber: 80
```

Deploy:

```shell
planton apply -f deployment.yaml
```

This creates a single-replica Deployment and a ClusterIP Service named `my-app` in `my-namespace`.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec.namespace` | `StringValueOrRef` | Target namespace. Accepts a literal name (`{ value: my-namespace }`) or a reference to a `KubernetesNamespace` resource. |
| `spec.container.app` | `WorkloadContainer` | The main application container. Its `image.repo` and `image.tag` are required; its ports drive the Service. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.createNamespace` | `bool` | `false` | Create the namespace if it does not exist. Leave `false` when the namespace pre-exists or is owned by a `KubernetesNamespace` resource. |
| `spec.version` | `string` | `main` | Deployment track, stamped as a pod label so multiple tracks of one app coexist in a namespace. Deployment pipelines set this from the git branch. Lowercase letters, numbers, hyphens; max 30 chars. |
| `spec.container.sidecars` | `WorkloadContainer[]` | `[]` | Sidecar containers running alongside the app in every replica. Full container surface applies; each sidecar must be named. |
| `spec.pod.serviceAccount` | `StringValueOrRef` | namespace `default` | ServiceAccount pods run as — a literal name or a reference to a `KubernetesServiceAccount` resource's `status.outputs.service_account_name`. |
| `spec.pod.automountServiceAccountToken` | `bool` | cluster default | Set `false` to withhold the Kubernetes API token from pods that never call the API. |
| `spec.pod.initContainers` | `WorkloadContainer[]` | `[]` | Run to completion, in order, before app containers start — migrations, config templating, dependency waits. |
| `spec.pod.scheduling` | `WorkloadScheduling` | — | Node selectors, tolerations, node/pod affinity, and topology spread constraints. |
| `spec.pod.securityContext` | `WorkloadPodSecurityContext` | — | Pod-level baseline every container inherits: `runAsNonRoot`, `fsGroup`, supplemental groups, sysctls, seccomp. |
| `spec.pod.terminationGracePeriodSeconds` | `int64` | `30` | Seconds between SIGTERM and SIGKILL; size it to cover `preStop` hooks plus drain time. |
| `spec.availability.replicas` | `int32` | `1` | Desired replica count; the floor when autoscaling is enabled. |
| `spec.availability.horizontalPodAutoscaling` | `KubernetesDeploymentHpa` | disabled | `enabled`, `maxReplicas`, and at least one of `targetCpuUtilizationPercent` (1–100) or `targetMemoryUtilization` (quantity, e.g. `1Gi`). |
| `spec.availability.strategy` | `KubernetesDeploymentStrategy` | RollingUpdate 25%/25% | `type` (`RollingUpdate` or `Recreate`), `maxUnavailable`, `maxSurge`. `"0"` unavailable with `"1"` surge is the zero-downtime pattern; surge and unavailable cannot both be zero. |
| `spec.availability.podDisruptionBudget` | `KubernetesDeploymentPodDisruptionBudget` | disabled | `enabled` plus exactly one of `minAvailable` or `maxUnavailable` (absolute or percentage). |
| `spec.availability.minReadySeconds` | `int32` | `0` | Seconds a new pod must stay ready before counting as available — a cheap flap detector during rollouts. |
| `spec.availability.revisionHistoryLimit` | `int32` | `10` | Old ReplicaSets retained for rollback; `0` disables rollback. |
| `spec.availability.progressDeadlineSeconds` | `int32` | `600` | Seconds a rollout may stall before being marked failed in the Deployment's conditions. |
| `spec.availability.paused` | `bool` | `false` | Pause rollouts to batch several spec changes into one. |

Each container (app, sidecar, init) additionally supports `command`/`args`, `workingDir`, `imagePullPolicy`, `env` (variables, secrets, `envFrom`), `resources`, `livenessProbe`/`readinessProbe`/`startupProbe`, `volumeMounts` (each mount carries its volume source: ConfigMap, Secret, EmptyDir, HostPath, or PVC), `lifecycle` (`postStart`/`preStop`, including a kubelet-native `sleep`), and a full `securityContext`.

## Examples

### Production Web Service with Autoscaling

CPU-based autoscaling between 2 and 10 replicas, zero-downtime rollouts, and a disruption budget:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: api-server
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDeployment.api-server
spec:
  namespace:
    value: production
  container:
    app:
      image:
        repo: ghcr.io/acme/api-server
        tag: v2.3.1
      resources:
        requests:
          cpu: 100m
          memory: 256Mi
        limits:
          cpu: 1000m
          memory: 1Gi
      ports:
        - name: http
          containerPort: 8080
          servicePort: 80
          appProtocol: http
      readinessProbe:
        httpGet:
          path: /healthz
          portNumber: 8080
        periodSeconds: 5
      lifecycle:
        preStop:
          sleep:
            seconds: 10
  availability:
    replicas: 2
    horizontalPodAutoscaling:
      enabled: true
      maxReplicas: 10
      targetCpuUtilizationPercent: 70
    strategy:
      type: RollingUpdate
      maxUnavailable: "0"
      maxSurge: "1"
    podDisruptionBudget:
      enabled: true
      minAvailable: "1"
```

### Background Worker

No ports, so no Service is created — the workload is unreachable by design and receives work over connections it initiates:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: queue-worker
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDeployment.queue-worker
spec:
  namespace:
    value: workers
  container:
    app:
      image:
        repo: ghcr.io/acme/worker
        tag: v1.8.0
      env:
        variables:
          - name: WORKER_CONCURRENCY
            value: "4"
        secrets:
          - name: QUEUE_CONNECTION_STRING
            secretRef:
              name: queue-secrets
              key: connection-string
```

### Hardened Service with Composed Identity

Passes the restricted Pod Security Standard and runs as a `KubernetesServiceAccount` referenced by output:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesDeployment
metadata:
  name: payments
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesDeployment.payments
spec:
  namespace:
    value: payments
  container:
    app:
      image:
        repo: ghcr.io/acme/payments
        tag: v4.1.0
      ports:
        - name: http
          containerPort: 8080
          servicePort: 80
      readinessProbe:
        httpGet:
          path: /healthz
          portNumber: 8080
      volumeMounts:
        - name: tmp
          mountPath: /tmp
          emptyDir:
            sizeLimit: 128Mi
      securityContext:
        readOnlyRootFilesystem: true
        allowPrivilegeEscalation: false
        capabilities:
          drop: ["ALL"]
  pod:
    serviceAccount:
      valueFrom:
        kind: KubernetesServiceAccount
        name: payments-identity
        fieldPath: status.outputs.service_account_name
    automountServiceAccountToken: false
    securityContext:
      runAsNonRoot: true
      runAsUser: 10001
      runAsGroup: 10001
      fsGroup: 10001
      seccompProfile:
        type: RuntimeDefault
    scheduling:
      topologySpreadConstraints:
        - maxSkew: 1
          topologyKey: topology.kubernetes.io/zone
          whenUnsatisfiable: ScheduleAnyway
  availability:
    replicas: 3
    strategy:
      type: RollingUpdate
      maxUnavailable: "0"
      maxSurge: "1"
    podDisruptionBudget:
      enabled: true
      minAvailable: "2"
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | The namespace the workload was deployed into |
| `deploymentName` | `string` | Name of the Deployment object as created in the cluster |
| `service` | `string` | The Kubernetes Service fronting the replicas; empty when the app container exposes no ports |
| `selectorLabels` | `string` | Pod selector labels as `k=v,k=v` — the exact labels the Service selects on, ready for NetworkPolicy podSelectors, `kubectl get pods -l`, and pod-affinity terms |
| `portForwardCommand` | `string` | Ready-to-run `kubectl port-forward` command for reaching the workload without external exposure |
| `kubeEndpoint` | `string` | In-cluster DNS endpoint of the Service (e.g. `my-app.my-ns.svc.cluster.local`) — the handle exposure components and sibling workloads connect to |

## Related Components

- [KubernetesServiceAccount](/docs/catalog/kubernetes/serviceaccount) — the identity pods run as; reference it from `spec.pod.serviceAccount` so workload-identity bindings and pull secrets live on the identity, not the workload
- [KubernetesRbac](/docs/catalog/kubernetes/rbac) — grants Kubernetes API permissions to the workload's identity
- [KubernetesHttpRoute](/docs/catalog/kubernetes/http-route) — publishes the workload at a hostname by referencing its exported `service` / `kubeEndpoint` outputs
- [KubernetesNamespace](/docs/catalog/kubernetes/namespace) — provides the target namespace; reference it from `spec.namespace` to deploy both in one run
- [KubernetesConfigMap](/docs/catalog/kubernetes/configmap) / [KubernetesSecret](/docs/catalog/kubernetes/secret) — configuration and credentials consumed via env sources or volume mounts
- [KubernetesStatefulSet](/docs/catalog/kubernetes/statefulset) — for workloads needing stable identity or per-replica storage
