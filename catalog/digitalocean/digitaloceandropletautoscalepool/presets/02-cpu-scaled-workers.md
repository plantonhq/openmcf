# CPU-Scaled Workers

This preset runs a worker fleet that grows from two to eight droplets as CPU load rises and shrinks back when it falls -- elasticity for batch processors, queue consumers, and bursty backends.

## When to Use

- Queue-driven or batch workloads whose load arrives in waves
- Web tiers with pronounced daily peaks where paying for peak capacity all day is waste

## Key Configuration Choices

- **min 2 / max 8** -- the floor keeps the service alive at zero load; the ceiling is the cost cap under sustained (or runaway) load. Set the max from budget.
- **CPU target 0.7** -- scale-ups trigger as the fleet averages past 70% CPU; the 10-minute cooldown stops flapping on short spikes.
- **Agent enabled** -- non-negotiable for dynamic pools: the utilization metrics ARE the agent's telemetry.
- **Bootstrap via userData** -- members are cattle; every boot configures itself identically, and local state is one scale-in from gone.

## What You Get

A fleet whose size follows real load between your bounds -- billing between 2 and 8 droplets instead of a fixed peak-sized fleet.
