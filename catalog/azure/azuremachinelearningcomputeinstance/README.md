# Overview

The **AzureMachineLearningComputeInstance** component creates a compute instance on an Azure Machine Learning workspace -- a single always-on VM serving as one data scientist's cloud workstation: notebooks, interactive debugging, and small jobs, pre-wired into the workspace's data and tooling.

## Purpose

- **Personal ML workstations as declarative infrastructure**: VM size, ownership, identity, and networking -- provisioned per team member, reviewed and versioned like everything else.
- **The admin-provisions-for-the-team pattern**: `assignToUser` creates an instance locked to a specific colleague -- the shape platform teams standardize on.
- **Typed references end-to-end**: the workspace, subnet, and user-assigned identities all wire by reference -- chart-ready.
- **Honest cost surface**: an instance is a running VM billing around the clock unless stopped -- documented where it bites, not discovered on the invoice.

## Key Features

- Full azurerm v5 surface: VM size (case-insensitive on the provider), personal authorization with user assignment, system/user-assigned identity, SSH access with the service-assigned username and port surfaced as outputs, VNet placement, node public-IP and local-auth toggles.
- The provider's replacement contract recorded loudly: this resource has NO update path -- every change (tags included) replaces the instance, and its OS disk and local files go with it.
- The service's placement rules modeled honestly: instances always run in the workspace's region, and their names are reserved region-wide per subscription.

## Use Cases

- **Per-scientist workstations**: one instance per team member, assigned by the platform team.
- **Notebook-first development**: the compute behind ML studio notebooks and VS Code remote sessions.
- **Interactive debugging**: a personal VM with the workspace's datastores and environments mounted.

## Future Enhancements

- Schedule-based auto-stop as the provider grows a schedule surface.
