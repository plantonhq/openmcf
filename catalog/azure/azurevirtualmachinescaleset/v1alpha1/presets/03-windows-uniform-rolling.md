# Windows Fleet with Automatic OS Upgrades (Uniform)

This preset creates a UNIFORM-orchestration Windows Server fleet that patches itself at the image level: automatic OS image upgrades roll new `latest` releases across zone-balanced instances in health-probe-gated batches, with automatic instance repair replacing anything that stops answering.

## When to Use

- Windows web/application tiers behind an Azure Load Balancer that should track new OS image releases without operator involvement
- Fleets that need UNIFORM-only capabilities: automatic OS image upgrades, the LB health probe as the health signal, strict zone balancing
- Organizations with existing Windows Server licenses (Azure Hybrid Benefit is on by default here)

## Key Configuration Choices

- **`orchestrationMode: UNIFORM`** -- automatic OS image upgrades and `healthProbeId` are UNIFORM capabilities; FLEXIBLE fleets use the health extension and manual image rolls instead
- **`automaticOsUpgrade` + `version: latest`** -- Azure watches the image publisher and rolls new releases batch by batch, gated by the probe; rollback stays enabled (the safer default)
- **`healthProbeId`** -- references the load balancer's `status.outputs.probe_ids.<probe>` map output; the probe is the health signal for upgrades AND repair
- **`zoneBalance: true`** -- strict instance balance across all three zones, trading scale-out speed for exact spread
- **Short `computerNamePrefix`** -- Windows caps computer names at 15 characters including the per-instance suffix

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | `AzureResourceGroup` status outputs |
| `<your-fleet-name>` | Scale-set name, unique within the resource group | Your naming convention |
| `<admin-password>` | Admin password (8-123 chars, 3 of 4 complexity classes) -- prefer a secret reference | Your secret manager / Config Manager |
| `<subnet-resource-id>` | Full ARM ID of the instances' subnet | `AzureSubnet` status outputs |
| `<load-balancer-backend-pool-id>` | The pool's ARM ID | `AzureLoadBalancer` `status.outputs.backend_pool_ids.<pool>` |
| `<load-balancer-probe-id>` | The health probe's ARM ID | `AzureLoadBalancer` `status.outputs.probe_ids.<probe>` |

## Related Presets

- **01-stateless-web-flexible** -- the Linux FLEXIBLE counterpart with the health extension
- **02-spot-batch** -- interruption-tolerant batch capacity at spot prices
