# Kubernetes CronJob: Research Documentation

## Introduction

A CronJob is a controller that manufactures Jobs on a clock. That framing carries the whole design: nothing about the *work* is CronJob-specific — the template inside is a full batch/v1 JobSpec — and everything CronJob-specific is about *time*: when runs fire, in which time zone, what happens when they overlap, how late a run may start, and how much finished history is kept.

Planton's **KubernetesCronJob** component models both halves at full depth: the scheduling controls at the top level and the complete Job surface (Indexed completions, per-index retries, pod failure policies, success policies) inside `job_template`, on the same shared workload container/pod core as every other workload kind.

## The Scheduling Machinery

### Cron expressions, without macros

The schedule is a standard 5-field cron expression — minute, hour, day-of-month, month, day-of-week: `"0 3 * * *"` is daily at 03:00, `"*/15 * * * *"` every 15 minutes, `"0 9 * * 1-5"` weekdays at 09:00. Two overlapping day fields (day-of-month and day-of-week) OR together when both are restricted — a classic cron trap worth remembering. Prefer standard expressions over `@daily`-style macros, which are not portable across controller versions.

### Time zones

Historically, CronJob schedules ran in the time zone of the kube-controller-manager process — an implementation detail that changed meaning when clusters moved or control planes were rebuilt. The `timeZone` field (stable since Kubernetes 1.27) fixes the schedule to an IANA zone name like `Europe/Berlin`, DST transitions included. The operational rule is simple: if a human cares what wall-clock time the run fires, set the zone explicitly; the controller's local clock is usually UTC, but that is convention, not contract.

DST corner cases are worth knowing for schedules between 01:00 and 03:00: when clocks skip forward over the scheduled time, the run is skipped for that day; when clocks fall back, schedules that resolve twice fire twice. Scheduling outside the transition window avoids the question entirely.

### The missed-run model and the 100-run cutoff

The controller computes missed schedule times whenever it reconciles. Runs can be missed for real reasons: controller downtime, a suspended CronJob being resumed, an unschedulable cluster. `startingDeadlineSeconds` decides what happens next — a missed run may still start within the deadline, and counts as failed-and-skipped beyond it. Unset means no deadline: late runs start whenever possible.

There is a hard backstop behind this: if the controller finds **more than 100 consecutive missed runs**, it stops scheduling the CronJob entirely and emits an event, refusing to guess which of the missed runs matter. With no deadline set, the 100-miss window stretches back to the last run ever started — so a frequent schedule (every minute) suspended for two hours trips the cutoff. An explicit `startingDeadlineSeconds` bounds the window the controller examines (only misses within the deadline count), which is why frequent schedules should always set one; 300 seconds on a 15-minute schedule keeps the counter permanently bounded.

### Concurrency: the three policies

`concurrencyPolicy` governs the moment the next run comes due while the previous run's Job is still active:

- **`Allow`** (the upstream default) — both run. Correct only when runs are independent — and the reason many teams discover their "hourly" job is running six copies
- **`Forbid`** — skip the new run. The safe default for anything touching shared state: backups, migrations, billing. Its failure mode is silence — a hung run blocks the schedule with no error — which is why `activeDeadlineSeconds` in the template is the essential companion
- **`Replace`** — cancel the running Job, start the fresh one. Right for latest-state-wins synchronization, where a stale run has no value; requires idempotent work that tolerates mid-flight cancellation

This component defaults to `Forbid`, deliberately diverging from upstream's `Allow`: overlapping cron runs are the classic scheduled-workload incident, so overlap is opt-in.

### Suspension

`suspend: true` stops the controller from creating new Jobs; work already in flight is unaffected — it pauses the schedule, not the workload. Note the interaction with the missed-run model above: a long suspension accumulates missed runs, and on resume, the deadline (or the 100-run cutoff) decides what happens.

## Retention: History Limits vs TTL

Two mechanisms bound the residue of finished runs, and they compose:

- **History limits** (CronJob level) — `successfulJobsHistoryLimit` (Kubernetes default 3) and `failedJobsHistoryLimit` (default 1) keep the most recent N finished Jobs; older ones are deleted, newest-first retention. This is *count-based* and is the CronJob-native mechanism
- **`ttlSecondsAfterFinished`** (template level) — each Job deletes itself N seconds after finishing. This is *time-based* and runs independently of the history limits; a TTL shorter than the schedule interval can empty the history before the limits ever apply

The usual production choice is history limits alone: they already bound retention, and a week of nightly successes (`successfulJobsHistoryLimit: 7`) with a few failures kept for post-mortems covers the debugging need. Reach for the TTL when Jobs are large or frequent enough that count-based retention still leaves too much behind. Keep `failedJobsHistoryLimit` at 1 or higher always — a failed run whose Job was already deleted is a failure with no logs.

## The Template Is a Full Job

Everything documented for KubernetesJob applies inside `job_template`, unchanged:

- **Indexed completions** — a monthly report fanning out six partitions per run, each pod reading `JOB_COMPLETION_INDEX`, is just an Indexed Job stamped out on a schedule
- **Per-index retries** — `backoff_limit_per_index` and `max_failed_indexes` budget failures per partition
- **Pod failure policies** — the DisruptionTarget pattern (ignore infrastructure-caused failures, fail fast on unrecoverable exit codes) matters *more* on CronJobs, because scheduled work tends to run on spot capacity at off-peak hours, exactly when autoscalers churn nodes
- **Success policies** — leader/worker topologies per run

The one deliberate difference: the template carries no `suspend` of its own — suspension lives at the CronJob level, where it pauses the schedule itself.

Two template-level settings deserve CronJob-specific emphasis. `activeDeadlineSeconds` is what makes `Forbid` safe (a hung run is killed instead of blocking every future run), and the sidecar rule is sharper here: a sidecar that never exits keeps every run from completing — and under `Forbid`, one bad run stalls the entire schedule.

## Deployment Methods Landscape

### Level 0: kubectl

```bash
kubectl create cronjob nightly-report --image=ghcr.io/acme/report:v2 --schedule="0 3 * * *" -- ./report.sh
```

Immediate and imperative; expresses nothing about concurrency, time zone, deadlines, or retention — every one of which is where CronJobs actually go wrong.

### Level 1: Raw YAML

The full surface, validated only at the API server. CronJobs stack the Job kind's cross-field rules (failure policy requires `restartPolicy: Never`, per-index backoff requires Indexed mode) on top of their own (cron expression syntax, valid concurrency policies, valid time zone names) — each discovered one apply at a time.

### Level 2/3: Terraform / Pulumi

Full IaC lifecycle with state and drift detection; schemas remain stringly-typed where CronJobs are most error-prone (schedule expressions, policy enums), so mistakes still surface at apply.

### The Planton approach

The spec validates the cron expression shape, concurrency policy, and every Job-level cross-field rule at schema level with CEL, before anything reaches a cluster. The scheduling/work split is explicit in the message structure, and the template reuses the shared workload core — configure a container once, and the knowledge transfers across all five workload kinds. Both Pulumi and Terraform modules consume the same validated spec with feature parity.

## Production Best Practices

1. **Set `timeZone` on every schedule a human reasons about**: "03:00" should mean 03:00 somewhere specific, through DST and cluster rebuilds
2. **Match the concurrency policy to the work's semantics**: shared mutable target → `Forbid`; latest-state-wins sync → `Replace`; genuinely independent runs → `Allow`. Never leave the choice implicit
3. **Pair `Forbid` with a template deadline**: `activeDeadlineSeconds` converts "hung run blocks the schedule forever" into a bounded, visible failure
4. **Set `startingDeadlineSeconds` on frequent schedules**: it skips valueless stale runs and keeps the 100-consecutive-missed-runs cutoff from ever triggering
5. **Keep failed history**: at least one failed Job retained, or failures leave no logs behind
6. **Design runs idempotent**: every policy except `Forbid` can cancel or duplicate work, and even `Forbid` skips runs; scheduled work that cannot tolerate a skip or a rerun is misdesigned
7. **Verify sidecars exit**: the one spec mistake that silently stalls an entire schedule
8. **Compose identity**: reference a `KubernetesServiceAccount` from `job_template.pod.service_account`, grant permissions via `KubernetesRbac`

## Conclusion

CronJobs concentrate their failure modes in time, not in containers: overlap, drift, missed-run pileups, and silent stalls. A spec that makes the time-side controls explicit — zone-pinned schedules, opt-in overlap, bounded lateness, bounded retention — and validates them before deployment removes the incident classes that follow scheduled workloads around, while the full Job surface inside the template keeps the work itself as expressive as any one-shot batch run.

## References

- [Kubernetes CronJobs Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)
- [Running Automated Tasks with a CronJob](https://kubernetes.io/docs/tasks/job/automated-tasks-with-cron-jobs/)
- [Kubernetes Jobs Documentation](https://kubernetes.io/docs/concepts/workloads/controllers/job/)
- [Pod Failure Policy](https://kubernetes.io/docs/concepts/workloads/controllers/job/#pod-failure-policy)
- [CronJob API Reference](https://kubernetes.io/docs/reference/kubernetes-api/workload-resources/cron-job-v1/)
