# Spot Cost-Optimized Batch

An EC2 Spot compute environment using the SPOT_PRICE_CAPACITY_OPTIMIZED strategy — up to ~90% below On-Demand pricing for retry-tolerant batch workloads, with instance selection that balances price against interruption risk.

## When to Use

- Fault-tolerant, retryable jobs: ETL, media transcoding, simulation sweeps, ML data prep
- Large-scale processing where compute cost dominates
- As the FIRST environment in a job queue's order, with an On-Demand environment as overflow

## What It Configures

- **Five instance types across three families** — more Spot pools means fewer interruptions and better prices
- **SPOT_PRICE_CAPACITY_OPTIMIZED** — the AWS-recommended Spot strategy; also keeps the environment eligible for in-place updates. No Spot Fleet role is needed (that is a BEST_FIT-only requirement)
- **Three Availability Zones** — capacity diversity is the main interruption defense
- **No bid percentage** — defaults to the On-Demand price cap; the realized price is usually far lower

## What to Customize

- Replace the region/subnet/security-group placeholders and the ECS instance profile ARN
- Pair jobs with a retry strategy on their `AwsBatchJobDefinition` that RETRYs on `"Host EC2*"` status reasons (Spot reclaims) while EXITing on real failures
- Compose the overflow pattern: an `AwsBatchJobQueue` with this environment at order 1 and an On-Demand environment at order 2
