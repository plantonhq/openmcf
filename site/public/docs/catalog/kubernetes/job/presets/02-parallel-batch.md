---
title: "Parallel Batch Job (Indexed)"
description: "This preset partitions work across 10 numbered completions, running 3 pods at a time. Each pod receives its completion index (0–9) through the `batch.kubernetes.io/job-completion-index` annotation..."
type: "preset"
rank: "02"
presetSlug: "02-parallel-batch"
componentSlug: "job"
componentTitle: "Job"
provider: "kubernetes"
icon: "package"
order: 2
---

# Parallel Batch Job (Indexed)

This preset partitions work across 10 numbered completions, running 3 pods at a time. Each pod receives its completion index (0–9) through the `batch.kubernetes.io/job-completion-index` annotation and the `JOB_COMPLETION_INDEX` environment variable, so each worker can claim its own slice of the data — shard N of a table, file N of a manifest, page N of an export.

Retries are budgeted per index rather than per Job: one flaky partition burns only its own retries, and the Job tolerates up to two permanently failed partitions before giving up.

## When to Use

- Processing partitioned data where each partition maps to one index
- Fan-out batch work that benefits from bounded parallelism
- Workloads where partitions fail independently and one bad shard should not sink the run

## Key Configuration Choices

- **`completionMode: Indexed` + `completions: 10`** — ten numbered partitions; the Job succeeds when every index has one successful pod
- **`parallelism: 3`** — at most three pods run at once, bounding cluster load; the controller works through the remaining indexes as pods finish
- **`backoffLimitPerIndex: 2`** — each index gets its own retry budget of 2; without this, one flaky partition could exhaust the global budget and fail indexes that never ran
- **`maxFailedIndexes: 2`** — the run tolerates two permanently failed indexes; a third terminates the whole Job. Requires `backoffLimitPerIndex`
- **`restartPolicy: Never`** — required for per-index backoff; each attempt gets a fresh pod, and failed pods remain for debugging

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-worker-image>` | Worker image that reads its partition from `JOB_COMPLETION_INDEX` | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-worker-command>` | Command that processes one partition | Your batch tooling |

## Related Presets

- **01-database-migration** — Single-completion one-shot Job
- **03-resilient-batch** — Failure-policy rules that classify failures before they burn retries
