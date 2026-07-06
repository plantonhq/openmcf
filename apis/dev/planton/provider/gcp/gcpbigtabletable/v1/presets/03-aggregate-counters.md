# Aggregate Counter Table

A metering/analytics table using Bigtable's server-side aggregate cells:
the `intsum` family increments atomically at write time, eliminating
read-modify-write races and application-side counter logic.

## When to use

Usage metering, rate counters, leaderboard tallies — anywhere many
writers increment the same logical counter concurrently.

## What to customize

- The aggregate `type` — `intsum` (running total), `intmin`/`intmax`
  (extremes), or `inthll` (approximate distinct counts).
- Retention on both families.

## Composes with

`GcpBigtableInstance` upstream (reference its `instance_name` output).
