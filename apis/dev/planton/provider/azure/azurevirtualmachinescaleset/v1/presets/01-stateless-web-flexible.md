# Stateless Web Fleet (Flexible)

This preset creates a FLEXIBLE-orchestration Linux fleet built for stateless web workloads: three zone-spread instances on ephemeral OS disks, accelerated networking, load-balancer pool membership from the member side, health-checked rolling upgrades, and automatic instance repair.

## When to Use

- Web and API tiers behind an Azure Load Balancer that scale horizontally and hold no local state
- Fleets that should upgrade safely (health-checked batches) and heal themselves (automatic repair)
- New workloads generally -- FLEXIBLE is Azure's recommended orchestration for anything that does not need a UNIFORM-only capability

## Key Configuration Choices

- **Ephemeral OS disk** (`diffDiskSettings` + `READ_ONLY` caching) -- free, fast local OS disks wiped on stop/reimage; exactly right when instances are cattle stamped from the image
- **Member-side pool membership** (`loadBalancerBackendAddressPoolIds`) -- references the load balancer's `status.outputs.backend_pool_ids.<pool>` map output, so the fleet joins the pool without the LB knowing about it
- **Health extension + ROLLING upgrades + repair** -- the `ApplicationHealthLinux` extension is the load-bearing piece: rolling batches pause on unhealth, and repair replaces instances that stop reporting healthy
- **Zone spreading with `platformFaultDomainCount: 1`** -- zones are the resilience unit; one fault domain per zone is the FLEXIBLE contract

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the subnet and load balancer) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | `AzureResourceGroup` status outputs |
| `<your-fleet-name>` | Scale-set name, unique within the resource group | Your naming convention |
| `<your-ssh-public-key>` | OpenSSH public key for the admin account | `ssh-keygen -t ed25519` |
| `<subnet-resource-id>` | Full ARM ID of the instances' subnet | `AzureSubnet` status outputs |
| `<load-balancer-backend-pool-id>` | The pool's ARM ID | `AzureLoadBalancer` `status.outputs.backend_pool_ids.<pool>` |

## Related Presets

- **02-spot-batch** -- interruption-tolerant batch capacity at spot prices
- **03-windows-uniform-rolling** -- a Windows UNIFORM fleet with automatic OS image upgrades
