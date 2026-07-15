# AWS Batch Scheduling Policy: Concepts

A scheduling policy converts a job queue from FIFO into a fair-share scheduler. This reference covers the fairness model, the dials, and the operational one-way door.

## The Fairness Model

Jobs carry a `shareIdentifier` at submission. Batch tracks each share's usage and schedules the next job from whichever share is furthest BELOW its fair allocation. The allocation is weight-derived and inverted from intuition:

- **`weight_factor` is a divisor, not a share count.** A 0.5-weight share is entitled to twice the capacity of a 1.0 share; a 2.0-weight share to half. Range 0.0001-999.9999; unlisted shares run at 1.0.
- **Within a share**, jobs order by scheduling priority (the job definition's `scheduling_priority` or SubmitJob's override), then submission time.

## The Dials

- **`share_decay_seconds` (0-604800)** — the sliding window over which PAST usage counts. At 0, only currently-running jobs matter; a team that just finished a huge burst immediately competes as an equal. At 3600, that burst counts against them for an hour. Longer windows smooth fairness across bursty workloads.
- **`compute_reservation` (0-99)** — headroom for shares not currently running, computed as `(reservation/100)^N` of capacity where N is the number of ACTIVE shares. Reservation 10 with one active share holds back 10%; with two active shares 1%. It guarantees a quiet team's first job does not queue behind a busy team's backlog.
- **Wildcard identifiers** — `adhoc*` aggregates every identifier with that prefix into ONE share. Useful for treating a whole class of ad-hoc submissions as a single budget.

## Shareable by Design

The policy is account-level and referenced by ARN, so one policy can govern many queues. Changing a weight updates every attached queue's behavior immediately — define organizational fairness once, not per queue.

## The One-Way Door

A queue that has ever had a scheduling policy can swap it for another but can never return to FIFO — AWS rejects removing `scheduling_policy_arn`. The modules surface this in the queue's spec comment; plan queue fairness before first deploy where possible.

## Composition

| Consumed by | Via | Output referenced |
|-------------|-----|-------------------|
| AwsBatchJobQueue | `scheduling_policy` | `scheduling_policy_arn` |

The submission side (SubmitJob's `shareIdentifier`) is data plane — jobs choose their share at runtime; the policy only defines how shares divide capacity.
