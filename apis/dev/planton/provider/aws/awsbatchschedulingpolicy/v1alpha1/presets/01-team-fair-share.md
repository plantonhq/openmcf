# Team Fair Share

A three-share policy demonstrating every fairness dial: weighted team shares (remember: LOWER weight = MORE capacity), a wildcard ad-hoc class, hour-long usage decay, and headroom for teams not currently running.

## When to Use

- Multiple teams submitting to one queue on shared compute
- Preventing bulk backfill jobs from starving interactive submissions

## What It Configures

- **`teamData` at 0.5** — entitled to twice the capacity of `teamMl` at 1.0
- **`adhoc*` at 2.0** — every identifier starting with "adhoc" aggregates into one half-capacity share
- **`shareDecaySeconds: 3600`** — a team's usage over the past hour counts against its allocation
- **`computeReservation: 10`** — capacity held back so a quiet team's first job starts promptly

## What to Customize

- Replace `<aws-region>` and the share identifiers with your team/workload taxonomy
- Attach to queues via `schedulingPolicy` — noting a queue can replace but never remove its policy
- Jobs participate by submitting with a `shareIdentifier`; unlisted identifiers run at weight 1.0
