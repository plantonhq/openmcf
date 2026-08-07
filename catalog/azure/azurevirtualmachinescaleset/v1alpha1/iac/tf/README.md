# AzureVirtualMachineScaleSet Terraform Module

## Overview

This Terraform module provisions an Azure Virtual Machine Scale Set
using the `azurerm` provider (`~> 4.0`). One spec, three dispatch
branches: the module reads `orchestration_mode` and the OS profile and
creates exactly one of `azurerm_orchestrated_virtual_machine_scale_set`
(FLEXIBLE — the default), `azurerm_linux_virtual_machine_scale_set`, or
`azurerm_windows_virtual_machine_scale_set` (UNIFORM), mirroring how ARM
has one scale-set resource type whose orchestration mode gates its
capabilities.

The template-versus-reference line the spec draws is realized here:
per-instance NICs, disks, and optional public IPs are stamped from
inline templates (they live and die with each instance), while
everything shared — subnets, load-balancer backend pools and NAT rules
(resolved through the load balancer's name-keyed map outputs), NSGs,
user-assigned identities, public IP prefixes, Key Vaults, and the
UNIFORM health probe — arrives as resolved reference values in the
tfvars.

Spec-level validation enforces ARM's conditional matrix (mode gates,
rolling-upgrade pairings, spot pairings, confidential-disk chains), so
the module maps fields without re-validating them. Enum values arrive as
FULL proto value names (`READ_ONLY`, `LINUX_AUTOMATIC_BY_PLATFORM`,
`DATA_PREMIUM_LRS`) and are mapped to ARM's casing through the
`locals.tf` maps. This module realizes ranked mixed-SKU profiles
(`sku_profile` with PRIORITIZED ranks), which the Pulumi module's pinned
SDK cannot express — see the `PARITY-EXCEPTION` note beside the
`sku_profile` block in `main.tf`.

Lifecycle notes worth knowing before operating this resource:

- The orchestration mode, OS flavor, and name are the fleet's identity —
  changing any of them replaces the scale set.
- Upgrade mode MANUAL (the default) applies template changes to NEW
  instances only; ROLLING requires health monitoring (the application
  health extension, or the health probe on UNIFORM sets).
- Removing a zone from `zones` replaces the set; adding one does not.
- `single_placement_group` can go true→false but never back.

## Resources Created

Exactly one of:

- `azurerm_orchestrated_virtual_machine_scale_set.flexible` — FLEXIBLE orchestration (either OS, nested under `os_profile.linux_configuration` / `windows_configuration`)
- `azurerm_linux_virtual_machine_scale_set.linux_uniform` — UNIFORM + Linux
- `azurerm_windows_virtual_machine_scale_set.windows_uniform` — UNIFORM + Windows

## Variables

| Variable | Type | Description |
|----------|------|-------------|
| `metadata` | object | Planton metadata (name, org, env) |
| `spec` | object | Scale-set specification |

Key spec fields:

| Field | Required | Description |
|-------|----------|-------------|
| `region` | yes | Azure region; must match every referenced subnet and load balancer |
| `resource_group` | yes | Resource group name |
| `name` | yes | Scale-set name (instance computer names derive from `os_profile.computer_name_prefix` or this) |
| `orchestration_mode` | no | `FLEXIBLE` (default) or `UNIFORM`; selects the dispatch branch |
| `sku_name` / `instances` | yes / no | The size and count; FLEXIBLE `Mix` activates `sku_profile` |
| `os_profile` | yes | Exactly one of `linux` / `windows` with the OS's auth and management surface |
| `os_disk` / `data_disks` | yes / no | Disk templates (ephemeral OS disks, dialed Ultra/PremiumV2 performance, CMK/confidential encryption) |
| `network_interfaces` | yes | NIC templates with subnet, LB pool/NAT-rule, App-Gateway, ASG, and per-instance public-IP surfaces |
| `upgrade_policy` | no | Mode, rolling batches, UNIFORM automatic OS upgrades + health probe |
| `spot` | no | Presence makes the fleet spot: eviction, bid, UNIFORM restore, FLEXIBLE priority mix |
| `platform_fault_domain_count` | FLEXIBLE: yes | 1 with zones, or the region's max for regional spreading |
| `identity` / `security` / `extensions` / `placement` / `zones` / `tags` | no | The management, security, and governance surfaces |

## Outputs

| Output | Description |
|--------|-------------|
| `scale_set_id` | Full ARM ID of the scale set — what an AzureVirtualMachine's `availability.virtual_machine_scale_set_id` references (FLEXIBLE attach) and what autoscale/monitoring scope to |
| `scale_set_name` | The scale set's name as deployed |
| `unique_id` | The scale set's globally unique ARM-assigned identifier |
| `system_assigned_identity_principal_id` | The system-assigned identity's principal ID (UNIFORM sets with a SYSTEM_ASSIGNED identity); empty otherwise |

## Usage

The module is invoked by the Planton CLI with a generated tfvars file:

```bash
planton tofu plan --manifest fleet.yaml --module-dir catalog/azure/azurevirtualmachinescaleset/v1alpha1/iac/tf
planton tofu apply --manifest fleet.yaml --module-dir catalog/azure/azurevirtualmachinescaleset/v1alpha1/iac/tf
```
