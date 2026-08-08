# AzureVirtualMachine Terraform Module

## Overview

This Terraform module provisions an Azure Virtual Machine using the
`azurerm` provider (`~> 5.0`). ARM models Linux and Windows VMs as separate
management surfaces (different auth contracts, patch vocabularies, and OS
settings), so the module deploys exactly one of
`azurerm_linux_virtual_machine` or `azurerm_windows_virtual_machine` from
the spec's explicit OS discriminator (`os_profile.linux` XOR
`os_profile.windows`), plus one
`azurerm_virtual_machine_data_disk_attachment` per referenced data disk.

The VM is deliberately just the machine (matching Azure's own model):

- Network presence comes from referenced network interfaces
  (`network_interface_ids`, at least one; the first is primary); public
  IPs, NSG filtering, and subnet placement live NIC-side.
- Data volumes are referenced managed disks realized as attachment
  resources -- the disk (and its data) outlives the VM, and detaching is
  just removing the spec entry.
- Only the OS disk is inline: it is born and dies with the VM, unless the
  VM boots from an existing referenced OS disk (`os_managed_disk_id`), in
  which case the OS profile carries no authentication fields (spec-level
  validation enforces the pairing).

Lifecycle notes worth knowing before operating this module: name, region,
zone, image source, admin credentials, `custom_data`, and the
security/confidential posture are the VM's identity -- changing any of them
replaces the VM (the OS disk with it; data disks and NICs survive, which is
exactly why they are referenced). Resizing (`size`) reboots in place. Spot
settings are fixed at creation.

The provider block is empty: credentials come from the standard `ARM_*`
environment variables (or ambient Azure CLI/MSI auth), supplied by the
platform at run time.

## Resources Created

- `azurerm_linux_virtual_machine.main` OR
  `azurerm_windows_virtual_machine.main` -- exactly one, from the OS
  discriminator
- `azurerm_virtual_machine_data_disk_attachment.main` -- one per
  `data_disk_attachments` entry, keyed by LUN

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) -- feeds the derived resource tags |
| `spec` | object | Virtual machine specification |

Key spec fields (enum fields take the spec enums' name strings; references
are resolved to literal ARM IDs by the platform before the module runs):

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; must match every referenced NIC and disk |
| `resource_group` | yes | Resource group name |
| `name` | yes | VM name, unique within the resource group (1-64 chars Linux, 1-15 Windows) |
| `size` | yes | VM size (SKU), e.g. `Standard_D2s_v3` |
| `network_interface_ids` | yes | Resolved ARM IDs of the attached NICs; the first is primary |
| `os_profile` | yes | Exactly one of `linux` (SSH-first auth, patch mode, license type) or `windows` (password auth, automatic updates, hotpatching, timezone, WinRM, unattend content, license type), plus optional `computer_name` |
| `os_disk` | yes | Caching + storage SKU (required), optional size/name, ephemeral `diff_disk_settings`, CMK encryption sets, confidential `security_encryption_type`, Write Accelerator |
| `source_image_reference` / `source_image_id` / `os_managed_disk_id` | one of | Exactly one image source: marketplace coordinates, custom/gallery image ID, or an existing OS disk (the golden-disk path) |
| `data_disk_attachments` | no | Referenced first-class disks mounted at LUNs (0-63) with caching; realized as attachment resources |
| `identity` | no | `SYSTEM_ASSIGNED`, `USER_ASSIGNED` (+ resolved identity ARM IDs), or `SYSTEM_AND_USER_ASSIGNED` |
| `spot` | no | Presence makes the VM a spot instance: `eviction_policy` (`DEALLOCATE`/`DELETE`), `max_bid_price` (default -1) |
| `availability` | no | Zone XOR availability set, proximity placement, capacity reservation, dedicated host/group, Flexible scale-set attach + fault domain |
| `security` | no | Secure boot + vTPM (trusted launch), encryption at host |
| `patching` | no | Assessment mode, reboot setting, safety-check bypass (the latter two require the OS profile's AUTOMATIC_BY_PLATFORM patch mode) |
| `boot_diagnostics` | no | Presence enables; empty URI uses Azure-managed storage |
| `gallery_applications` | no | VM Applications installed at deployment (up to 100, ordered) |
| `termination_notification` / `os_image_notification` | no | Scheduled events; presence enables |
| `plan` | no | Marketplace purchase plan (third-party images) |
| `custom_data` / `user_data` | no | Base64: cloud-init (secret, first boot, replaces on change) vs IMDS-readable data (updatable, never secret) |
| `extensions_time_budget`, `provision_vm_agent`, `allow_extension_operations` | no | Extension provisioning budget and agent gates (Azure defaults: PT1H30M, true, true) |
| `disk_controller_type` | no | `SCSI` / `NVME`; unset applies Azure's default for the size/image |
| `additional_capabilities` | no | Ultra SSD attachability, hibernation |
| `secrets` | no | Key Vault certificates installed at provisioning (Windows entries name a certificate store) |
| `edge_zone` | no | Edge Zone pinning |
| `tags` | no | User tags merged over the metadata-derived tags (user wins on collision) |

Only explicit spec choices are ever sent to ARM (null otherwise), so an
unspecified field and Azure's default deploy identically.

## Outputs

| Output | Description |
|--------|-------------|
| `vm_id` | Full ARM ID of the VM |
| `vm_name` | The VM's name as deployed |
| `virtual_machine_guid` | The 128-bit unique GUID Azure assigns the VM |
| `private_ip_address` | Primary private IP across the attached NICs |
| `public_ip_address` | Primary public IP across the attached NICs (empty for private-only VMs) |
| `computer_name` | The OS hostname the VM booted with |
| `system_assigned_identity_principal_id` | Principal ID of the system-assigned identity (empty when the identity type doesn't include it) |

## Usage

```hcl
module "azure_virtual_machine" {
  source = "./iac/tf"

  metadata = {
    name = "app-vm"
    org  = "mycompany"
    env  = "production"
  }

  spec = {
    region         = "eastus"
    resource_group = "prod-rg"
    name           = "app-vm"
    size           = "Standard_D2s_v3"

    network_interface_ids = [
      "/subscriptions/.../providers/Microsoft.Network/networkInterfaces/app-vm-nic"
    ]

    os_profile = {
      linux = {
        admin_username = "azureuser"
        ssh_public_keys = [
          { public_key = "ssh-ed25519 AAAA..." }
        ]
      }
    }

    os_disk = {
      caching              = "READ_WRITE"
      storage_account_type = "PREMIUM_LRS"
    }

    source_image_reference = {
      publisher = "Canonical"
      offer     = "ubuntu-24_04-lts"
      sku       = "server"
      version   = "latest"
    }

    availability = {
      zone = "1"
    }

    boot_diagnostics = {}
  }
}
```

## Required Permissions

The deploying credential needs
`Microsoft.Compute/virtualMachines/write` on the resource group -- held via
Virtual Machine Contributor, Contributor, or Owner. Attaching NICs and data
disks additionally exercises `Microsoft.Network/networkInterfaces/join/action`
and `Microsoft.Compute/disks/write`. Encryption at host requires the
subscription to have the `EncryptionAtHost` feature registered.
