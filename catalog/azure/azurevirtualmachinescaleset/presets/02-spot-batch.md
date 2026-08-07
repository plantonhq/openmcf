# Spot Batch Fleet (Flexible, Mixed SKUs)

This preset creates a FLEXIBLE-orchestration spot fleet for interruption-tolerant batch work: ten instances drawn from three interchangeable VM sizes (capacity-optimized), a two-instance guaranteed on-demand base, DELETE eviction so autoscaling replaces reclaimed capacity, and a 15-minute pre-termination drain signal.

## When to Use

- Batch processing, CI runners, render farms, queue workers -- anything that tolerates interruption and restarts
- Workloads where compute cost dominates and spot's discount (often 60-90%) is worth the eviction risk
- Fleets that must keep SOME guaranteed capacity (the on-demand base) while the rest rides spot prices

## Key Configuration Choices

- **`skuName: Mix` + `skuProfile`** -- spot capacity comes and goes per size; drawing from three interchangeable sizes with `CAPACITY_OPTIMIZED` allocation is the difference between a resilient spot fleet and one that starves
- **`priorityMix`** -- `baseRegularCount: 2` guarantees two on-demand instances before any spot math applies; `regularPercentageAboveBase: 0` makes everything above the base spot
- **`evictionPolicy: DELETE`** -- evicted instances are removed (not deallocated), so replacement capacity comes from a fresh allocation attempt across the size mix
- **`terminationNotification: PT15M`** -- workers get the scheduled event and drain their current item before the instance disappears

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | `AzureResourceGroup` status outputs |
| `<your-fleet-name>` | Scale-set name, unique within the resource group | Your naming convention |
| `<your-ssh-public-key>` | OpenSSH public key for the admin account | `ssh-keygen -t ed25519` |
| `<subnet-resource-id>` | Full ARM ID of the instances' subnet | `AzureSubnet` status outputs |
| `<base64-encoded-cloud-init>` | Worker bootstrap, base64-encoded | `base64 -w0 cloud-init.yaml` |

## Related Presets

- **01-stateless-web-flexible** -- health-checked web serving on FLEXIBLE orchestration
- **03-windows-uniform-rolling** -- a Windows UNIFORM fleet with automatic OS image upgrades
