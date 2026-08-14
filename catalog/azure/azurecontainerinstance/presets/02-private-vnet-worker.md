# Private VNet Worker

This preset runs a worker container inside your virtual network -- a private IP in a delegated subnet, pulling its image from Azure Container Registry as a managed identity, with no public surface and no secrets in the manifest.

## When to Use

- Queue processors, internal API consumers, and anything that must reach VNet-internal services
- Workloads whose image lives in a private registry and whose Azure access should be identity-based
- Replacing a "utility VM" that exists only to run one containerized worker

## Key Configuration Choices

- **`ipAddressType: Private` + `subnetId`** -- the subnet must carry the `Microsoft.ContainerInstance/containerGroups` delegation (the AzureSubnet kind's `delegations` field); Azure serializes group operations per subnet, so parallel deploys into one subnet queue up
- **No `dnsNameLabel`** -- private groups have no public DNS; the spec rejects the combination outright
- **The identity registry credential** -- grant the identity AcrPull on the registry instead of shipping a password; the same identity serves the workload's own Azure access
- **Identity is the one live-updatable setting** -- you can rotate or add identities without replacing the group

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-resource-group>` | The Planton name of your `AzureResourceGroup` resource | Planton console |
| `<your-delegated-subnet>` | The Planton name of the delegated `AzureSubnet` | Planton console (the subnet needs the containerGroups delegation) |
| `<your-user-assigned-identity>` | The Planton name of your `AzureUserAssignedIdentity` | Planton console |
| `<your-container-registry>` | The Planton name of your `AzureContainerRegistry` | Planton console (or replace `valueFrom` with `value:` and the login server) |

## Related Presets

- **Public Web Container** -- the internet-facing posture
- **Run-Once Job** -- the Never-restart batch shape
