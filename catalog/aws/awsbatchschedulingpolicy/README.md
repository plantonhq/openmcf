# AwsBatchSchedulingPolicy

An AWS Batch fair-share scheduling policy — the reusable rules that divide a job queue's compute capacity across share identifiers instead of processing jobs first-in-first-out.

## What It Is

Without a scheduling policy, one team's burst of ten thousand jobs starves every other submitter on the queue. With one, jobs carry a **share identifier** at submission (per team, per workload class) and Batch balances dispatch so each share receives capacity proportional to its weight.

The policy is a standalone, shareable object: many [job queues](../awsbatchjobqueue/README.md) can reference the same policy, so an organization's fairness rules are defined once. Every dial updates in place on a live policy.

## When to Use It

| Use Case | Description |
|----------|-------------|
| **Multi-team shared capacity** | Divide one compute environment's capacity across teams by weight. |
| **Workload-class fairness** | Keep bulk backfill from starving interactive/adhoc jobs on one queue. |
| **New-submitter headroom** | `compute_reservation` holds capacity back for shares not yet running. |

## When NOT to Use It

| Need | Use Instead |
|------|-------------|
| **Strict FIFO processing** | No policy — FIFO is the queue's default. |
| **Whole tiers of priority** | Separate queues with different `priority` values on one environment. |
| **Isolating compute entirely** | Separate compute environments per team. |

## Key Facts

- **Lower weight = MORE capacity.** A share with `weight_factor` 0.5 receives twice the capacity of a 1.0 share. Unlisted shares default to 1.0.
- **Wildcard prefixes.** `analytics*` aggregates `analyticsDaily` and `analyticsAdhoc` into one share. Up to 500 distributions per policy.
- **`share_decay_seconds`** makes fairness account for history (up to 7 days) — "you had the cluster all morning" counts against a share.
- **`compute_reservation`** reserves `(reservation/100)^N` of capacity for inactive shares, shrinking as more shares activate.
- **Attachment is a one-way door on the queue**: a queue's policy can be replaced but never removed (recreate the queue to return to FIFO).
- **Jobs opt in at submission** — SubmitJob's `shareIdentifier` (and `schedulingPriorityOverride`, or the job definition's `scheduling_priority`) decide how a job participates.

## Spec Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `region` | string | **Yes** | AWS region; attachable only to same-region queues. |
| `compute_reservation` | int | No | 0-99: % of capacity held back for inactive shares. |
| `share_decay_seconds` | int | No | 0-604800: sliding window of past usage counted against shares. |
| `share_distributions` | list (≤500) | No | `{share_identifier, weight_factor}` entries; identifiers may end with `*`. |

## Outputs

| Field | Description |
|-------|-------------|
| `scheduling_policy_arn` | What job queues reference through `scheduling_policy`. |
| `scheduling_policy_name` | The policy's name (from `metadata.name`). |

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsBatchSchedulingPolicy
metadata:
  name: team-fair-share
  org: my-org
spec:
  region: us-west-2
  computeReservation: 10
  shareDecaySeconds: 3600
  shareDistributions:
    - shareIdentifier: teamData
      weightFactor: 0.5   # twice the capacity of a 1.0 share
    - shareIdentifier: teamMl
      weightFactor: 1.0
    - shareIdentifier: adhoc*
      weightFactor: 2.0   # half the capacity -- backfill class
```

Attach it to a queue:

```yaml
spec:
  schedulingPolicy:
    valueFrom:
      kind: AwsBatchSchedulingPolicy
      name: team-fair-share
      fieldPath: status.outputs.scheduling_policy_arn
```

See docs/README.md for the fairness math and operational guidance.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
