# Personal Dev Instance

This preset creates a personal workstation assigned to a specific team member -- the admin-provisions-for-the-team shape: one general-purpose VM, locked to its owner, with a system identity for credential-free data access.

## When to Use

- Provisioning per-scientist workstations as part of onboarding
- Standardizing the team on a reviewed, versioned instance shape
- Any time the deploying principal is NOT the instance's intended user

## Key Configuration Choices

- **`authorizationType: personal` + `assignToUser`** -- the instance belongs to the assigned user from first boot; without `assignToUser`, the deploying credentials would own it
- **`STANDARD_DS3_V2`** -- a workstation size: this VM bills around the clock, so heavy training belongs on a compute cluster instead
- **`SYSTEM_ASSIGNED` identity** -- grant it Storage Blob Data Reader on the team's datastores so notebooks read data credential-free
- **No `ssh` block** -- the SSH port stays closed; the studio and VS Code remote cover most needs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-machine-learning-workspace-id>` | ARM ID of the parent workspace | `AzureMachineLearningWorkspace` status outputs (`machine_learning_workspace_id`), or reference it with valueFrom |
| `assignToUser.tenantId` (example UUID) | The owner's Entra ID tenant | `az account show --query tenantId` |
| `assignToUser.objectId` (example UUID) | The owner's Entra ID object ID | `az ad user show --id <owner-email> --query id` |

The `name` carries a realistic example (`alice-dev`) because instance names are reserved region-wide -- name it after its owner. The two `assignToUser` IDs carry example UUIDs (the tenant ID is format-validated at manifest time) -- both must be replaced with the owner's real values.
