# Kubernetes CronJob

Runs work on a recurring schedule as a batch/v1 CronJob through a single declarative manifest. Scheduling controls (cron expression, time zone, concurrency, history retention) live at the top level; the work itself — containers, pod configuration, and the full set of batch controls — lives in `jobTemplate`, which can express everything a standalone KubernetesJob can. The IaC module handles namespace resolution, secret materialization, and label merging automatically.

CronJobs front no Service and create no ingress — each scheduled run creates a Job whose pods run to completion. Concurrency defaults to `Forbid` (skip a run while the previous one is still going), deliberately safer than the upstream `Allow` default.

## What Gets Created

When you deploy a KubernetesCronJob resource, Planton provisions:

- **CronJob** — a batch/v1 CronJob with the schedule, time zone, concurrency policy, starting deadline, suspension flag, and history limits
- **Namespace** (optional) — created when `createNamespace` is true
- **Env Secret** (conditional) — literal secret env values across all containers are materialized into one workload-scoped Secret named `<metadata.name>-env-secrets`
- **Image Pull Secret** (conditional) — docker-registry credentials materialized as `<metadata.name>-image-pull`
- **Labels** — standard Planton tracking labels merged with any user-provided pod labels

## Prerequisites

- **Kubernetes credentials** configured via environment variables or Planton provider config
- **A Kubernetes namespace** that exists, is created via `createNamespace`, or is referenced as a `KubernetesNamespace` resource
- **A container image** containing the scheduled workload; each run's pod succeeds when every container exits 0

## Quick Start

Create a file `cronjob.yaml`:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesCronJob
metadata:
  name: nightly-cleanup
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesCronJob.nightly-cleanup
spec:
  namespace:
    value: backend
  schedule: "0 3 * * *"
  timeZone: America/New_York
  jobTemplate:
    container:
      app:
        image:
          repo: ghcr.io/acme/cleanup
          tag: v1.8.0
        command: ["./cleanup.sh"]
    restartPolicy: Never
    backoffLimit: 2
```

Deploy:

```shell
planton apply -f cronjob.yaml
```

This runs the cleanup every night at 03:00 New York time, never overlapping with a still-running previous run.

## Configuration Reference

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `spec.namespace` | `StringValueOrRef` | Namespace to run the CronJob in. Literal name or reference to a `KubernetesNamespace` resource. |
| `spec.schedule` | `string` | Standard 5-field cron expression, e.g. `"0 3 * * *"` (daily at 03:00) or `"*/15 * * * *"` (every 15 minutes). |
| `spec.jobTemplate.container.app` | `WorkloadContainer` | The main application container: image (repo + tag required), command/args, env, resources, volume mounts, security context. **Deploy-target contract:** pipelines inject built images at `spec.jobTemplate.container.app.image`. |

### Scheduling Controls (optional)

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.timeZone` | `string` | controller-local | IANA zone name (e.g. `Asia/Kolkata`) the schedule is evaluated in. Set it whenever wall-clock time matters. |
| `spec.startingDeadlineSeconds` | `int64` | — | How late a missed run may start before it is skipped. Also keeps the controller's 100-consecutive-missed-runs cutoff bounded on frequent schedules. |
| `spec.concurrencyPolicy` | `string` | `Forbid` | `Forbid` (skip the new run), `Allow` (run concurrently), or `Replace` (cancel the running Job, start fresh). |
| `spec.suspend` | `bool` | `false` | Stop scheduling future runs; Jobs already running are unaffected. |
| `spec.successfulJobsHistoryLimit` | `int32` | `3` | Completed successful Jobs retained for log inspection. |
| `spec.failedJobsHistoryLimit` | `int32` | `1` | Failed Jobs retained; keep at least one so failure logs survive. |

### Job Template (optional)

Everything a standalone KubernetesJob expresses, minus `suspend` (which lives at the CronJob level):

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `spec.jobTemplate.parallelism` / `completions` | `int32` | `1` / `1` | Pods in flight per run and successes required per run. |
| `spec.jobTemplate.completionMode` | `string` | `NonIndexed` | `Indexed` gives each pod a completion index for partitioned work. |
| `spec.jobTemplate.backoffLimit` | `int32` | `6` | Pod-failure retries before a run's Job fails. |
| `spec.jobTemplate.backoffLimitPerIndex` / `maxFailedIndexes` | `uint32` | — | Per-index retry budgets (Indexed mode, `restartPolicy: Never`). |
| `spec.jobTemplate.activeDeadlineSeconds` | `int64` | — | Hard wall-clock cap per run — essential with `Forbid`, where a hung run blocks all future runs. |
| `spec.jobTemplate.ttlSecondsAfterFinished` | `int32` | — | Per-Job auto-delete; usually left unset in favor of the history limits. |
| `spec.jobTemplate.restartPolicy` | `string` | `Never` | `Never` or `OnFailure`; `Always` is invalid for Jobs. |
| `spec.jobTemplate.podFailurePolicy` | `object` | — | Ordered failure-classification rules (`FailJob`, `Ignore`, `Count`, `FailIndex`). |
| `spec.jobTemplate.successPolicy` | `object` | — | Early-success rules for Indexed runs. |
| `spec.jobTemplate.container.sidecars[]` | `WorkloadContainer` | — | Named sidecars. **Every sidecar must exit** or no run ever completes. |
| `spec.jobTemplate.pod` | `WorkloadPod` | — | ServiceAccount reference, init containers, scheduling, security hardening, DNS, termination. |

## Examples

### Frequent Sync with Replace

Every 15 minutes; a stale overrunning sync is cancelled in favor of the fresh one:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesCronJob
metadata:
  name: index-sync
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.KubernetesCronJob.index-sync
spec:
  namespace:
    value: search
  schedule: "*/15 * * * *"
  startingDeadlineSeconds: 300
  concurrencyPolicy: Replace
  jobTemplate:
    container:
      app:
        image:
          repo: ghcr.io/acme/index-sync
          tag: v0.9.2
    restartPolicy: Never
    backoffLimit: 1
```

### Indexed Parallel Runs

Each monthly run fans out six numbered partitions, three at a time:

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesCronJob
metadata:
  name: monthly-report
  annotations:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: prod.KubernetesCronJob.monthly-report
spec:
  namespace:
    value: reporting
  schedule: "0 2 1 * *"
  timeZone: Europe/Berlin
  jobTemplate:
    container:
      app:
        image:
          repo: ghcr.io/acme/report-worker
          tag: v4.1.0
        command: ["/bin/sh", "-c", "render --section $JOB_COMPLETION_INDEX"]
    completionMode: Indexed
    completions: 6
    parallelism: 3
    restartPolicy: Never
    backoffLimitPerIndex: 2
    activeDeadlineSeconds: 14400
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `namespace` | `string` | Namespace the CronJob was created in |
| `cronJobName` | `string` | Name of the CronJob object in the cluster |
| `schedule` | `string` | The effective cron expression — the deployed truth, for dependents and audits |

## Related Components

- [KubernetesJob](/docs/catalog/kubernetes/kubernetesjob) — the one-shot counterpart; each scheduled run stamps out a Job
- [KubernetesDeployment](/docs/catalog/kubernetes/kubernetesdeployment) — for always-on services
- [KubernetesServiceAccount](/docs/catalog/kubernetes/kubernetesserviceaccount) — the identity run pods use; reference it from `spec.jobTemplate.pod.serviceAccount`
- [KubernetesRbac](/docs/catalog/kubernetes/kubernetesrbac) — grants permissions to that identity
- [KubernetesSecret](/docs/catalog/kubernetes/kubernetessecret) — holds credentials runs consume via `secretRef` env entries
