# Tag-Scoped Drills

This preset creates the focused monthly drill: on the 1st at 04:00
UTC, restore the LATEST recovery point of everything tagged
`tier: critical`, across snapshot and continuous points, with an
8-hour validation window.

## When to Use

- Proving the recovery story for the resources that actually matter,
  at a cadence that keeps costs negligible
- Regimes that require documented drills for critical systems

## What You Get

- Latest-within-7-days selection — the drill answers "can we restore
  what we'd actually restore in an incident"
- Tag-driven coverage (`aws:ResourceTag/tier = critical`) — resources
  opt in by tagging, no plan edits
- Both point types covered, so continuous (PITR) backups get proven
  too

## Customize

- Switch to `RANDOM_WITHIN_WINDOW` with a wider window to also prove
  retention depth
- Add `restoreMetadataOverrides` (lowercase keys) to restore into an
  isolated subnet or downsized instance type
- Add selections for the other critical resource types — one per type
