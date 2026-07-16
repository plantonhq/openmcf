# Serverless Auto-Pausing Database

This preset creates a serverless database that bills compute per second,
scales between a warm floor (0.5 vCores) and the SKU ceiling (2 vCores),
and pauses entirely after an hour of inactivity -- storage is the only
cost while paused.

## When to Use

- Development, test, and preview environments that sit idle most of the
  day
- Spiky or unpredictable production workloads where average utilization
  is well below peak

## Key Configuration Choices

- **`GP_S_Gen5_2`** -- the `_S_` infix is the serverless variant; the
  trailing number caps the vCores it can scale to
- **`autoPauseDelayInMinutes: 60`** -- the first connection after a
  pause pays a resume delay (seconds to a minute); set `-1` to disable
  auto-pause and keep per-second billing without pauses
- **`minCapacity: 0.5`** -- the always-warm floor; raise it when cold
  query latency matters more than idle cost

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<server-resource-name>` | The AzureMssqlServer's Planton resource name | Your server composition |
| `dev-scratch-db` | The database's name on the server | Your application |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
