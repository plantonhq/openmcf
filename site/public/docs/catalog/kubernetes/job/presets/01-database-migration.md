---
title: "Database Migration Job"
description: "This preset runs a schema migration to completion: one pod, one success required, a small retry budget, and automatic cleanup a day after the Job finishes. It is the canonical one-shot Job pattern —..."
type: "preset"
rank: "01"
presetSlug: "01-database-migration"
componentSlug: "job"
componentTitle: "Job"
provider: "kubernetes"
icon: "package"
order: 1
---

# Database Migration Job

This preset runs a schema migration to completion: one pod, one success required, a small retry budget, and automatic cleanup a day after the Job finishes. It is the canonical one-shot Job pattern — the same shape works for backfills, data fixes, and any other run-once task.

Identity is composed, not bundled: if the migration needs cloud or cluster permissions, reference a `KubernetesServiceAccount` from `spec.pod.serviceAccount` and grant permissions to that identity with `KubernetesRbac`.

## When to Use

- Database schema migrations run as part of a release
- One-time backfills, data corrections, or setup tasks
- Any task that must run exactly once to success and then stop

## Key Configuration Choices

- **Single completion, single pod** (`completions: 1`, `parallelism: 1`) — one successful run finishes the Job
- **`restartPolicy: Never`** — each failed attempt leaves its pod behind for `kubectl logs`, and the Job controller creates a fresh pod for the retry
- **`backoffLimit: 2`** — migrations that fail twice almost never succeed on the third try; fail fast and investigate instead of retrying six times (the Kubernetes default)
- **`activeDeadlineSeconds: 3600`** — a hung migration is killed after an hour rather than blocking the release forever
- **`ttlSecondsAfterFinished: 86400`** — the finished Job and its pods are deleted after 24 hours, leaving a full day to inspect logs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-migration-image>` | Image containing the migration tooling | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-migration-command>` | The command that runs the migration | Your migration tool's documentation |
| `<your-database-connection-string>` | Database connection URL | Your database configuration; move to a secret reference for production |

## Related Presets

- **02-parallel-batch** — Fan out partitioned work across Indexed pods
- **03-resilient-batch** — Distinguish unrecoverable failures from infrastructure disruption
