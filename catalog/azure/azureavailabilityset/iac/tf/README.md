# AzureAvailabilitySet Terraform Module

## Overview

Creates an availability set -- the classic pre-zones placement grouping that spreads VMs across separate fault domains (power/network/rack) and update domains (planned-maintenance batches) so one hardware failure or maintenance window cannot take them all down.

## Resources Created

- `azurerm_availability_set` -- the availability set (domain counts, managed alignment, optional proximity placement group, tags)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureAvailabilitySetSpec fields; the resource group and proximity-placement-group references arrive as resolved literals

## Outputs

- `availability_set_id` -- the set's ARM resource ID (what VMs reference to join it)
- `availability_set_name` -- the set's name

## Behavior Notes

- **The whole configuration is fixed at creation** -- only tags update in place; everything else forces replacement.
- **Unset optionals ride the provider defaults**: 5 update domains, 3 fault domains, `managed = true`. The module passes null so the provider's own defaults apply.
- **`managed = true` aligns fault domains with managed-disk storage** -- the right value for every managed-disk VM. Some regions support fewer than 3 fault domains; Azure rejects a count the region cannot provide.
- **An availability set is free** -- it bills nothing itself.
- **VMs join at creation** (AzureVirtualMachine's `availability.availability_set_id`); a VM cannot move into or out of a set in place.
