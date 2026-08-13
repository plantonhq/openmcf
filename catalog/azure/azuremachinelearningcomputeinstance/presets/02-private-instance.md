# Private Instance

This preset creates a VNet-placed workstation with no public IP -- the hardened posture for estates where personal compute must stay off the internet.

## When to Use

- Regulated estates where workstations may not carry public IPs
- Workspaces whose data sits behind private endpoints in the same VNet
- Pairing with ExpressRoute / VPN access so scientists reach the instance privately

## Key Configuration Choices

- **`nodePublicIpEnabled: false`** -- no public IP; reachability comes from the VNet (peering, VPN, or ExpressRoute)
- **`subnetId`** -- REQUIRED with public IP off on a plain workspace. On a MANAGED-NETWORK workspace, remove both the subnet and keep public IP off -- Azure networks the instance itself, and the provider enforces this against the live workspace at apply time
- **`personal` + `assignToUser`** -- ownership works the same as the personal preset

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-machine-learning-workspace-id>` | ARM ID of the parent workspace | `AzureMachineLearningWorkspace` status outputs (`machine_learning_workspace_id`), or reference it with valueFrom |
| `assignToUser.tenantId` (example UUID) | The owner's Entra ID tenant | `az account show --query tenantId` |
| `assignToUser.objectId` (example UUID) | The owner's Entra ID object ID | `az ad user show --id <owner-email> --query id` |
| `<your-subnet-id>` | ARM ID of the workstation subnet | `AzureSubnet` status outputs (`subnet_id`), or reference it with valueFrom |

Everything on this resource is fixed at creation -- decide the network posture before provisioning, not after.
