# Kubernetes Job

Runs work to completion on a Kubernetes cluster as a batch/v1 Job through a single declarative manifest: one-shot tasks like migrations and backfills, or parallel batch processing with Indexed completions, per-index retry budgets, pod failure policies, and success policies. The IaC module handles namespace resolution, secret materialization, and label merging automatically.

Jobs front no Service and create no ingress — pods run, finish, and are found afterward through the exported label selector. For scheduled work, use KubernetesCronJob; for always-on services, KubernetesDeployment.

## What Gets Created

When you deploy a KubernetesJob resource, Planton provisions:

- **Job** — a batch/v1 Job with the full set of batch controls (parallelism, completions, completion mode, retry budgets, deadlines, TTL, failure and success policies)
- **Namespace** (optional) — created when `createNamespace` is true
- **Env Secret** (conditional) — literal secret env values across all containers are materialized into one workload-scoped Secret named `<metadata.name>-env-secrets`
- **Image Pull Secret** (conditional) — docker-registry credentials materialized as `<metadata.name>-image-pull`
- **Labels** — standard Planton tracking labels merged with any user-provided pod labels

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that exists, is created via `createNamespace`, or is referenced as a `KubernetesNamespace` resource
- **A container image** containing the batch workload; the pod succeeds when every container exits 0

## Quick Start

Create a file `job.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesJob
metadata:
  name: db-migration
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesJob.db-migration
spec:
  namespace:
    value: backend
  container:
    app:
      image:
        repo: ghcr.io/acme/migrate
        tag: v3.2.0
      command: ["./migrate.sh"]
  restartPolicy: Never
  backoffLimit: 2
  ttlSecondsAfterFinished: 86400
```

Deploy:

```shell
planton apply -f job.yaml
```

This runs the migration once to success, retries at most twice, and cleans up the finished Job after 24 hours.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec.namespace` | `StringValueOrRef` | Namespace to run the Job in. Literal name or reference to a `KubernetesNamespace` resource. |
| `spec.container.app` | `WorkloadContainer` | The main application container: image (repo + tag required), command/args, env, resources, probes, volume mounts, security context. Its exit code decides pod success. |

### Batch Controls (optional)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.parallelism` | `int32` | `1` | Maximum pods running at once. |
| `spec.completions` | `int32` | `1` | Successful completions required; in Indexed mode, also the number of indexes. |
| `spec.completionMode` | `string` | `NonIndexed` | `NonIndexed` (interchangeable pods) or `Indexed` (each pod gets an index 0..completions−1 for partitioned work). |
| `spec.backoffLimit` | `int32` | `6` | Pod-failure retries before the Job fails, with exponential back-off. |
| `spec.backoffLimitPerIndex` | `uint32` | — | Per-index retry budget (Indexed mode, `restartPolicy: Never` only). |
| `spec.maxFailedIndexes` | `uint32` | — | Failed indexes tolerated before the Job terminates; requires `backoffLimitPerIndex`. |
| `spec.activeDeadlineSeconds` | `int64` | — | Hard wall-clock limit for the whole Job; exceeded ⇒ all pods killed, Job failed. |
| `spec.ttlSecondsAfterFinished` | `int32` | — | Seconds after finish before the Job and pods are auto-deleted. |
| `spec.suspend` | `bool` | `false` | Create without running; suspending a running Job deletes active pods and resets the deadline timer. |
| `spec.restartPolicy` | `string` | `Never` | `Never` (fresh pod per attempt) or `OnFailure` (in-place container restart). `Always` is invalid for Jobs. |
| `spec.podFailurePolicy` | `object` | — | Ordered rules classifying pod failures: `FailJob`, `Ignore`, `Count`, or `FailIndex` triggered by exit codes or pod conditions. Requires `restartPolicy: Never`. |
| `spec.successPolicy` | `object` | — | Early-success rules for Indexed Jobs (`succeededIndexes`, `succeededCount`). |

### Pod and Workload (optional)

| Field | Type | Description |
|-------|------|-------------|
| `spec.createNamespace` | `bool` | Create the namespace when it does not exist. |
| `spec.container.sidecars[]` | `WorkloadContainer` | Named sidecar containers. **Every sidecar must exit for the Job to complete.** |
| `spec.pod.serviceAccount` | `StringValueOrRef` | Identity pods run as — literal name or `KubernetesServiceAccount` reference. |
| `spec.pod.initContainers[]` | `WorkloadContainer` | Run to completion, in order, before the app starts. |
| `spec.pod.scheduling` | `object` | Node selector, tolerations, affinity, topology spread. |
| `spec.pod.securityContext` | `object` | Pod-level user/group identity and filesystem ownership. |

## Examples

### Indexed Parallel Batch

Ten numbered partitions, three at a time, with per-index retries:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesJob
metadata:
  name: shard-processor
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesJob.shard-processor
spec:
  namespace:
    value: batch
  container:
    app:
      image:
        repo: ghcr.io/acme/shard-worker
        tag: v1.4.0
      command: ["/bin/sh", "-c", "process --shard $JOB_COMPLETION_INDEX"]
  completionMode: Indexed
  completions: 10
  parallelism: 3
  restartPolicy: Never
  backoffLimitPerIndex: 2
  maxFailedIndexes: 2
  ttlSecondsAfterFinished: 86400
```

### Failure-Aware Batch on Spot Nodes

Fail fast on the application's "unrecoverable" exit code; never charge node disruption against the retry budget:

```yaml
apiVersion: kubernetes.planton.dev/v1
kind: KubernetesJob
metadata:
  name: resilient-batch
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesJob.resilient-batch
spec:
  namespace:
    value: batch
  container:
    app:
      image:
        repo: ghcr.io/acme/batch
        tag: v2.0.1
  restartPolicy: Never
  backoffLimit: 4
  podFailurePolicy:
    rules:
      - action: FailJob
        onExitCodes:
          operator: In
          values: [42]
      - action: Ignore
        onPodConditions:
          - type: DisruptionTarget
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | Namespace the Job was created in |
| `jobName` | `string` | Name of the Job object in the cluster |
| `selectorLabels` | `string` | Job-name label selector as `k=v` — ready for `kubectl get pods -l` / `kubectl logs -l` post-run inspection |

## Related Components

- [KubernetesCronJob](/docs/catalog/kubernetes/kubernetescronjob) — the scheduled counterpart; each run stamps out a Job like this one
- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — for always-on services
- [KubernetesServiceAccount](/docs/catalog/kubernetes/kubernetesserviceaccount) — the identity pods run as; reference it from `spec.pod.serviceAccount`
- [KubernetesRbac](/docs/catalog/kubernetes/kubernetesrbac) — grants permissions to that identity
- [KubernetesSecret](/docs/catalog/kubernetes/kubernetessecret) — holds credentials the Job consumes via `secretRef` env entries
