# Weekly Random Drills

This preset creates the restore-proof workhorse: every Monday 05:00
UTC, restore a RANDOM recovery point from the last 30 days for every
backed-up EBS volume, validate for 4 hours, delete the copy.

## When to Use

- Turning "we have backups" into "we have restores" with zero manual
  drills
- Feeding restore-time metrics into Backup Audit Manager's
  restore-time controls

## What You Get

- Random-within-window selection — older recovery points get exercised
  too, the stronger proof
- A 4-hour validation window per restored copy (then AWS deletes it)
- Coverage of every EBS recovery point in every vault

## Customize

- Add a selection per additional resource type (EC2, RDS, S3, ...) —
  each carries its own role and coverage
- Narrow `includeVaults` to vault ARNs when only some tiers need
  drills
- Tag-scope coverage with `protectedResourceConditions` instead of
  `protectedResourceArns: ["*"]` (exactly one of the two)
- Tests bill as real restores — cadence and coverage are the cost
  levers
