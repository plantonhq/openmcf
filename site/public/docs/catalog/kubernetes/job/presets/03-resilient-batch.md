---
title: "Resilient Batch Job (Failure Policy)"
description: "This preset classifies pod failures before they consume the retry budget. Two rules, evaluated in order: an application-signaled unrecoverable error (exit code 42) fails the Job immediately, and..."
type: "preset"
rank: "03"
presetSlug: "03-resilient-batch"
componentSlug: "job"
componentTitle: "Job"
provider: "kubernetes"
icon: "package"
order: 3
---

# Resilient Batch Job (Failure Policy)

This preset classifies pod failures before they consume the retry budget. Two rules, evaluated in order: an application-signaled unrecoverable error (exit code 42) fails the Job immediately, and infrastructure-caused failures (node drains, preemption) are ignored entirely. Everything else falls through to the default — count against `backoffLimit`.

The result is a Job that fails fast when retrying is pointless, survives cluster maintenance without burning retries, and still retries genuine transient errors.

## When to Use

- Batch work on clusters with node autoscaling, spot/preemptible nodes, or frequent drains
- Workloads that can detect unrecoverable conditions (bad input, missing config) and signal them with a dedicated exit code
- Any Job where "why did it fail" should determine "should it retry"

## Key Configuration Choices

- **Rule 1 — `FailJob` on exit code 42 (`In` operator)** — the application exits 42 when it knows retrying cannot succeed (malformed input, invalid configuration); the Job fails immediately with all pods terminated, instead of retrying four times to the same end. Pick any non-zero code and make it your workload's contract; exit code 0 can never appear in an `In` list because it means success
- **Rule 2 — `Ignore` on `DisruptionTarget`** — pods evicted by drains, preemption, or taint changes carry the `DisruptionTarget` condition; those failures say nothing about the workload, so a replacement pod is created and the retry budget is untouched
- **Rule order matters** — the first matching rule wins; put the most specific classification first
- **`restartPolicy: Never`** — required whenever `podFailurePolicy` is set; in-place container restarts would never surface as pod failures for the rules to inspect
- **`backoffLimit: 4`** — the budget now only counts genuine application failures, so a modest number suffices

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-batch-image>` | Batch image that exits 42 on unrecoverable errors | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-batch-command>` | The command that performs the batch work | Your batch tooling |

## Related Presets

- **01-database-migration** — Simple one-shot Job without failure classification
- **02-parallel-batch** — Indexed fan-out with per-index retry budgets
