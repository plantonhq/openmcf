# Spot-First With On-Demand Overflow

The canonical Batch cost pattern: jobs try the Spot environment first (order 1) and spill to On-Demand (order 2) only when Spot capacity runs out. Most work rides ~90%-discounted compute; nothing waits when Spot dries up.

## When to Use

- Retry-tolerant pipelines where cost matters but throughput must not stall
- Any workload already running Spot that occasionally hits capacity gaps

## What It Configures

- **Two environments in preference order** — the scheduler always prefers the lower order and falls forward on capacity, never load-balances
- **Same family requirement respected** — both referenced environments must be EC2-based (or both Fargate-based); AWS rejects mixed queues

## What to Customize

- Replace `<aws-region>` and both environment references (a `SPOT` or `FARGATE_SPOT` environment first, its On-Demand sibling second)
- Give jobs a retry strategy on their `AwsBatchJobDefinition` that RETRYs on `"Host EC2*"` status reasons so Spot reclaims re-run automatically
