# Monthly Report CronJob (Indexed)

This preset generates a partitioned report on the first of every month: each scheduled run stamps out an Indexed Job with six numbered completions (one per report section, region, or data shard), running three at a time. Every pod reads its partition from the `JOB_COMPLETION_INDEX` environment variable.

It demonstrates the composability of the spec: everything a standalone `KubernetesJob` can express — Indexed mode, per-index retries, deadlines — is equally expressible inside a CronJob's `jobTemplate`.

## When to Use

- Periodic reports or exports naturally partitioned by region, tenant, or shard
- Scheduled heavy batch work that benefits from bounded parallel fan-out
- Monthly/weekly aggregations too large for one sequential pod

## Key Configuration Choices

- **`schedule: "0 2 1 * *"` + `timeZone`** — 02:00 on the 1st of each month in your zone; billing-style deadlines make the explicit time zone essential
- **`completionMode: Indexed` + `completions: 6`** — six numbered partitions per run; the run succeeds when every index has one successful pod
- **`parallelism: 3`** — at most three partitions in flight, bounding memory and database load
- **`backoffLimitPerIndex: 2`** — one flaky partition retries alone instead of exhausting a shared budget
- **`activeDeadlineSeconds: 14400`** — the whole run is capped at four hours; with `Forbid`, a hung run would otherwise block next month's report
- **History limits 3/1** — the last three successful monthly runs remain inspectable

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace | Your namespace management or `KubernetesNamespace` resource |
| `<your-iana-time-zone>` | IANA time zone name, e.g. `Europe/Berlin` | The zone your reporting deadline is defined in |
| `<your-container-registry>/<your-report-image>` | Report generator that reads `JOB_COMPLETION_INDEX` | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |
| `<your-report-command>` | Command that renders one partition | Your reporting tooling |

## Related Presets

- **01-nightly-backup** — Daily single-pod schedule with secret-backed credentials
- **02-frequent-sync** — High-frequency schedule with Replace concurrency
