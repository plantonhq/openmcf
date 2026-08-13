# Data Flow Compute

This preset creates the managed data-flow compute -- the serverless Spark engine mapping data flows transform on, joined to the factory's managed virtual network with a small warm pool.

## When to Use

- Any factory that runs mapping data flows (the flows need this engine to execute)
- Flows that must reach private endpoints (the managed-virtual-network posture)
- Back-to-back flow schedules where cluster startup latency matters (the warm pool)

## Key Configuration Choices

- **`virtualNetworkEnabled: true`** -- the compute joins the factory's managed virtual network so data flows can reach managed private endpoints; the FACTORY must be deployed with `managedVirtualNetworkEnabled: true`, or Azure rejects this runtime
- **`timeToLiveMin: 10`** -- clusters stay warm 10 minutes after a run, so consecutive flows skip the several-minute cold start; warm minutes bill like run minutes, so match this to your inter-run gap
- **Size defaults apply** -- omitting `computeType`/`coreCount` gives General compute with 8 cores (Azure's smallest data flow cluster); scale up only when flows actually spill

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-data-factory>` | The Planton name of your `AzureDataFactory` resource | Planton console (or replace `valueFrom` with `value:` and the factory's ARM ID) |
| `<your-region>` | The Azure region for the compute, e.g. `eastus` -- or `AutoResolve` to run each flow nearest its data | Your factory's region (`az account list-locations -o table` for the menu) |

## Related Presets

- **Self-Hosted Bridge** -- the agent registration for data Azure cannot reach directly.
- **SSIS Runtime** -- the managed engine for lift-and-shift SSIS packages.
