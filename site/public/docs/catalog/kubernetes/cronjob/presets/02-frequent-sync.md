---
title: "Frequent Sync CronJob"
description: "This preset runs a synchronization task every 15 minutes with semantics tuned for high-frequency, latest-state-wins work: a stale run is replaced rather than protected, and a run that cannot start..."
type: "preset"
rank: "02"
presetSlug: "02-frequent-sync"
componentSlug: "cronjob"
componentTitle: "CronJob"
provider: "kubernetes"
icon: "package"
order: 2
---

# Frequent Sync CronJob

This preset runs a synchronization task every 15 minutes with semantics tuned for high-frequency, latest-state-wins work: a stale run is replaced rather than protected, and a run that cannot start promptly is skipped because the next one is minutes away.

## When to Use

- Cache warming, search-index refresh, or data synchronization where only the newest state matters
- Polling-style integrations that reconcile an external system on a short interval
- Any schedule frequent enough that skipping one slot is cheaper than running late

## Key Configuration Choices

- **`schedule: "*/15 * * * *"`** — every 15 minutes, on the quarter hour
- **`concurrencyPolicy: Replace`** — when a run overruns into the next slot, the stale run is cancelled and the fresh one starts. For a sync, the newest run always supersedes the old one; `Forbid` would skip fresh data to protect stale work, and `Allow` would race two syncs against the same target. Only choose `Replace` for idempotent work that tolerates mid-flight cancellation
- **`startingDeadlineSeconds: 300`** — a run that misses its slot by more than 5 minutes is skipped; with slots every 15 minutes, a late run has almost no value. The explicit deadline also keeps the controller's consecutive-missed-runs counter bounded — without any deadline, 100 consecutive misses (e.g. during prolonged suspension mishandling or controller downtime) stop scheduling entirely
- **`backoffLimit: 1`** — one retry per slot; persistent failures surface in the failed-jobs history rather than being retried into the next slot

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-sync-image>` | Image with your sync tooling | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-sync-command>` | The command that performs one sync pass | Your integration's documentation |

## Related Presets

- **01-nightly-backup** — Low-frequency schedule where overlap must be forbidden instead of replaced
- **03-monthly-report** — Indexed parallel fan-out on a monthly schedule
