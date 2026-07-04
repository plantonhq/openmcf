---
title: "Windows Server with Trusted Launch"
description: "This preset creates a zonal Windows Server 2022 VM with trusted launch (secure boot + vTPM), password authentication sourced from a secret reference, and Azure Hybrid Benefit licensing. The Windows..."
type: "preset"
rank: "02"
presetSlug: "02-windows-server"
componentSlug: "virtual-machine"
componentTitle: "Virtual Machine"
provider: "azure"
icon: "package"
order: 2
---

# Windows Server with Trusted Launch

This preset creates a zonal Windows Server 2022 VM with trusted launch (secure boot + vTPM), password authentication sourced from a secret reference, and Azure Hybrid Benefit licensing. The Windows management surface -- automatic updates, hotpatching, timezone, WinRM -- lives in the `windows` profile, mirroring ARM's own per-OS model.

## When to Use

- Windows Server application workloads (IIS, .NET Framework services, Active Directory members)
- Organizations bringing existing Windows Server licenses (Azure Hybrid Benefit)
- Any Windows VM on a Gen2 image, which should take trusted launch by default

## Key Configuration Choices

- **`adminPassword` by reference** -- Windows requires password authentication; source it from a secret (e.g. a Config Manager entry), never a manifest literal
- **`licenseType: WINDOWS_SERVER`** -- Azure Hybrid Benefit cuts the compute bill substantially when you hold a license with Software Assurance; drop the field for pay-as-you-go
- **Trusted launch** (`secureBootEnabled` + `vtpmEnabled`) -- the modern security baseline for Gen2 images; fixed at creation, so decide it here
- **Computer-name limit** -- Windows hostnames cap at 15 characters and default to the VM name; use `osProfile.computerName` if the VM name must be longer
- **Patch orchestration** -- the image's default is Windows Update (`AUTOMATIC_BY_OS`); set `patchMode: AUTOMATIC_BY_PLATFORM` plus `patching.rebootSetting` to hand scheduling to Azure Update Manager
- **RDP exposure** -- keep the NIC private and reach the VM through Azure Bastion or a VPN; a public IP on 3389 is an attack magnet

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region (must match the NIC's region) | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-network-interface-resource-name>` | Planton metadata name of the `AzureNetworkInterface` | Your NIC resource |
| `<your-secret-resource-name>` | The secret carrying the admin password | Your secret management |
