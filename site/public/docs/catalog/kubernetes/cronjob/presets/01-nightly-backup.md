---
title: "Nightly Backup CronJob"
description: "This preset runs a backup every night at 03:00 in an explicit time zone, never overlaps with a still-running previous backup, keeps a week of successful run history, and pulls the database credential..."
type: "preset"
rank: "01"
presetSlug: "01-nightly-backup"
componentSlug: "cronjob"
componentTitle: "CronJob"
provider: "kubernetes"
icon: "package"
order: 1
---

# Nightly Backup CronJob

This preset runs a backup every night at 03:00 in an explicit time zone, never overlaps with a still-running previous backup, keeps a week of successful run history, and pulls the database credential from an existing Kubernetes Secret rather than embedding it.

The spec splits in two: scheduling controls (`schedule`, `timeZone`, `concurrencyPolicy`, history limits) live at the top level, and the work itself — the container, retries, and deadline — lives in `jobTemplate`.

## When to Use

- Nightly database or filesystem backups
- Any scheduled task where the wall-clock time matters and overlap is dangerous

## Key Configuration Choices

- **`schedule: "0 3 * * *"` + `timeZone`** — daily at 03:00 in the named IANA zone. Without `timeZone`, the schedule follows the controller's local clock (usually UTC, but not guaranteed) — always set it when the run time matters
- **`concurrencyPolicy: Forbid`** — if last night's backup is somehow still running, tonight's is skipped rather than racing it against the same target. This is also the kind's default; it is spelled out here because it is load-bearing for backups
- **`activeDeadlineSeconds: 7200` in the template** — the guard that makes Forbid safe: a hung backup is killed after two hours instead of silently blocking every future run
- **`successfulJobsHistoryLimit: 7` / `failedJobsHistoryLimit: 3`** — a week of successful runs and the last three failures stay inspectable via `kubectl logs`
- **Secret via `secretRef`** — the credential comes from an existing Kubernetes Secret (managed by `KubernetesSecret` or an external secrets operator); nothing sensitive lives in this manifest
- **`backoffLimit: 2`** — two retries per night; a backup that fails three times needs a human, not more retries

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-iana-time-zone>` | IANA time zone name, e.g. `America/New_York` | The zone your backup window is defined in |
| `<your-container-registry>/<your-backup-image>` | Image with your backup tooling | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-backup-command>` | The command that performs the backup | Your backup tool's documentation |
| `<your-backup-destination>` | Where backups are written (bucket, path) | Your storage configuration |
| `<your-existing-secret-name>` / `<your-secret-key>` | Kubernetes Secret and key holding the database credential | `kubectl get secrets -n <namespace>` or your `KubernetesSecret` resource |

## Related Presets

- **02-frequent-sync** — High-frequency schedule where only the latest run matters
- **03-monthly-report** — Indexed parallel fan-out on a monthly schedule
