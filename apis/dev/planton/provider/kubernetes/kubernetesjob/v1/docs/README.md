# Kubernetes Job: Research Documentation

## Introduction

The batch/v1 Job is Kubernetes' run-to-completion primitive: create pods, run them until enough succeed, never restart the finished result. Everything interesting about Jobs lives in the gap between that one-line description and production reality — how failures are counted, how partitioned work is distributed, how infrastructure disruption is distinguished from application error, and how finished work is cleaned up.

Planton's **KubernetesJob** component models the full modern JobSpec on top of the shared workload container/pod core, so a user who knows how to configure a Deployment's containers already knows how to configure a Job's — only the batch controls are new.

## The Completion Model

### NonIndexed: interchangeable pods

In the default `NonIndexed` mode, pods are fungible workers. The Job succeeds after `completions` pods (default 1) exit successfully, with at most `parallelism` running at once. Which pod performed which unit of work is invisible to Kubernetes — coordination, if needed, is the application's problem (a work queue, a claim table).

This is the right mode for two shapes of work:

- **Single completion** (`completions: 1`) — migrations, backfills, one-shot scripts
- **Worker-pool draining** — N identical workers pulling from an external queue until it is empty; here `completions` equals the worker count and each worker exits when the queue is drained

### Indexed: numbered partitions

`Indexed` mode gives every completion a stable identity: each pod receives an index from 0 to completions−1, exposed through the `batch.kubernetes.io/job-completion-index` annotation, the `JOB_COMPLETION_INDEX` environment variable, and the pod hostname. Each index must succeed exactly once for the Job to complete.

This turns the Job controller into a static work distributor: shard N of a table, file N of a manifest, partition N of a dataset. No queue infrastructure needed — the partition assignment IS the index. Indexed mode is also the foundation for the per-index failure controls and success policies below, none of which make sense when pods are interchangeable.

## Failure Handling, Layer by Layer

### Layer 1: restart_policy — inside the pod

`restart_policy` decides what the kubelet does when a container exits non-zero. `Never` (this component's default) leaves the failed pod in place and lets the Job controller create a replacement — one pod per attempt, every failure inspectable with `kubectl logs`. `OnFailure` restarts the container in place, reusing the pod; cheaper, but the failed container's state is gone.

`Never` is the production default for a second reason: `pod_failure_policy` requires it, because in-place container restarts never surface as pod failures for policy rules to match.

### Layer 2: backoff_limit — the global retry budget

Failed pods are retried with exponential back-off (10s, 20s, 40s… capped at 6 minutes) until `backoff_limit` failures accumulate (Kubernetes defaults to 6). Then the whole Job is marked failed and its pods terminated.

One subtlety worth designing around: when `backoff_limit_per_index` is set, upstream stops counting failures globally and the plain `backoff_limit` becomes effectively unlimited unless set explicitly.

### Layer 3: backoff_limit_per_index — retries per partition

For Indexed Jobs, a global budget has a failure mode: one pathologically flaky partition can exhaust the budget and fail indexes that never got to run. `backoff_limit_per_index` gives every index its own budget, so each partition succeeds or fails on its own merits. `max_failed_indexes` then bounds the damage — the Job terminates once that many indexes have failed for good; unset, every index runs to its own outcome and the Job finishes with whatever mix results.

### Layer 4: pod_failure_policy — classify before counting

The retry budget treats all failures identically, but failures are not identical. A pod killed by a node drain says nothing about the workload; a pod that exits with "invalid configuration" will fail every retry. The pod failure policy inserts a classification step: ordered rules, first match wins, unmatched failures fall through to default counting.

Each rule pairs one **action** with exactly one **trigger** (a container exit-code check or a pod-condition check, never both):

- **`FailJob`** — terminate the whole Job immediately, ignoring the budget. For application-signaled unrecoverable errors: the workload exits with a contract code (say, 42) meaning "retrying cannot help"
- **`Ignore`** — the failure does not count; a replacement pod is created. For infrastructure-caused failures
- **`Count`** — explicit default handling; useful to terminate a rule chain
- **`FailIndex`** — mark only this pod's index failed, without retrying it. Indexed mode only, requires `backoff_limit_per_index`

### The DisruptionTarget pattern

The canonical rule pair, and the reason the feature exists:

```yaml
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

Kubernetes stamps the `DisruptionTarget` condition on pods terminated by node drains, preemption, or taint eviction. Ignoring those failures makes batch work first-class on spot/preemptible nodes and autoscaled clusters: maintenance churns pods freely without ever burning the retry budget, while genuine application failures still count. Rule order matters — the most specific classification (the exit-code contract) goes first, because the first matching rule wins.

Exit-code rules have two guardrails: `values` cannot contain 0 for the `In` operator (0 is success and is never considered), and `container_name` optionally scopes the check to one container.

## Success Policies: Leader/Worker Topologies

By default, an Indexed Job succeeds only when every index succeeds. `success_policy` relaxes that: the Job is declared succeeded as soon as **any** rule is satisfied, and lingering pods are terminated.

Each rule sets `succeeded_indexes` (comma-separated integers or ranges, e.g. `"0"` or `"0-2,4"`), `succeeded_count` (minimum number of successes), or both — when combined, only successes within the listed indexes are counted.

The motivating topology is leader/worker batch frameworks (MPI-style): index 0 is the driver whose exit code carries the verdict; the workers are its helpers. `succeeded_indexes: "0"` declares the Job done the moment the leader succeeds, instead of waiting for workers that may have nothing left to do. Quorum shapes (`succeeded_count: 5` of 10 replicated attempts) fall out of the same mechanism.

## Deadlines, Suspension, and Cleanup

- **`active_deadline_seconds`** is a whole-Job wall clock, counted from the Job's start. When it expires, all pods are killed and the Job fails regardless of remaining retries — it outranks every retry control. It is the standard guard against runaway batch work
- **`suspend`** stops pod creation. Suspending a running Job deletes its active pods (the workload must tolerate that) and — easy to miss — resets the start time, and with it the `active_deadline_seconds` timer
- **`ttl_seconds_after_finished`** deletes the Job and its pods after it finishes, success or failure. 0 deletes immediately; unset keeps the Job forever. The production balance is retention long enough for post-mortem log access (a day is common), because once the Job is deleted, `kubectl logs` has nothing to read

## The Sidecar Caveat

A Job pod completes only when **all** of its containers exit. A log shipper or proxy sidecar that runs an infinite loop keeps the pod Running after the app finishes — the Job never completes, the deadline eventually kills it, and the run is marked failed despite the work having succeeded.

The component models sidecars as full containers (probes, mounts, security context, lifecycle hooks — anything the app container can express), and its spec documentation carries the warning explicitly: use termination-aware sidecars, or have the app signal them to shut down when its work is done. Init containers have no such caveat — they run to completion by definition, before the app starts, which makes them the safe place for setup work.

## Deployment Methods Landscape

### Level 0: kubectl

```bash
kubectl create job my-migration --image=ghcr.io/acme/migrate:v3 -- ./migrate.sh
```

Immediate, imperative, unrepeatable. Fine for experiments; nothing about the retry posture, deadline, or cleanup is expressed.

### Level 1: Raw YAML

The full surface is available, but validation happens only at the API server, and Jobs have unusually many cross-field rules: `pod_failure_policy` requires `restartPolicy: Never`, `backoff_limit_per_index` requires Indexed mode, exit-code lists cannot contain 0 with `In`, each policy rule takes exactly one trigger. Raw YAML discovers each of these at apply time, one round-trip per mistake.

### Level 2/3: Terraform / Pulumi

Full IaC lifecycle with state, drift detection, and composition. The provider schemas are untyped where it matters, though: policy actions, operators, and completion modes are plain strings, and the cross-field rules above still surface only at apply.

### The Planton approach

The spec moves the API server's cross-field rules to validation time (CEL expressions on the schema), models the batch surface one-to-one with upstream naming, and inherits the shared workload core — so container, pod, scheduling, and hardening configuration is identical across every workload kind. Both Pulumi and Terraform modules consume the same validated spec with feature parity.

## Production Best Practices

1. **Choose the completion mode deliberately**: external queue → NonIndexed worker pool; static partitions → Indexed. Retrofitting index awareness into a queue-draining design (or vice versa) is a rewrite
2. **`restart_policy: Never` unless pod churn is a problem**: attempt-per-pod debuggability is worth more than pod reuse, and failure policies require it
3. **Budget retries per index for partitioned work**: pair `backoff_limit_per_index` with `max_failed_indexes` so one bad shard neither sinks nor stalls the run
4. **Define an exit-code contract**: reserve one code for "do not retry" and wire it to `FailJob`; always pair it with an `Ignore` on `DisruptionTarget` so infrastructure noise stays out of the budget
5. **Set both a deadline and a TTL**: `active_deadline_seconds` bounds the run, `ttl_seconds_after_finished` bounds the residue
6. **Keep sidecars terminating**: verify every sidecar exits when the app does — it is the one spec mistake that turns a succeeded workload into a failed Job
7. **Compose identity**: reference a `KubernetesServiceAccount` from `pod.service_account` and grant permissions via `KubernetesRbac`; the Job creates neither

## Conclusion

The Job spec rewards understanding its layers: restart policy inside the pod, retry budgets above it, failure classification above that, and success policies at the top. Modeled fully and validated early, those layers make batch work on Kubernetes predictable — including on the churning, spot-priced clusters where batch work actually runs.

## References

- [Kubernetes Jobs Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
- [Indexed Jobs for Parallel Processing](https://kubernetes.io/docs/tasks/job/indexed-parallel-processing-static/)
- [Pod Failure Policy](https://kubernetes.io/docs/concepts/workloads/controllers/job/#pod-failure-policy)
- [Backoff Limit per Index](https://kubernetes.io/docs/concepts/workloads/controllers/job/#backoff-limit-per-index)
- [Job Success Policy](https://kubernetes.io/docs/concepts/workloads/controllers/job/#success-policy)
- [Automatic Cleanup for Finished Jobs](https://kubernetes.io/docs/concepts/workloads/controllers/ttlafterfinished/)
- [Job API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/job-v1/)
