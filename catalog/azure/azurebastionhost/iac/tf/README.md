# AzureBastionHost Terraform Module

## Overview

Creates an Azure Bastion host -- the managed jump service that opens RDP/SSH sessions to virtual machines over their private addresses. Dedicated-infrastructure SKUs (Basic/Standard/Premium) deploy into the network's `AzureBastionSubnet` with a Standard static public IP; the Developer SKU attaches to a virtual network on Azure-shared infrastructure.

## Resources Created

- `azurerm_bastion_host` -- the Bastion host

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureBastionHostSpec fields; the resource group, subnet, public IP, and virtual network references arrive as resolved literals; the SKU arrives as the enum NAME and is mapped to the ARM wire value in `locals.tf`

## Outputs

- `bastion_host_id` -- the host's ARM resource ID
- `bastion_host_name` -- the host's name
- `dns_name` -- the DNS name sessions connect through
- `private_only_enabled` -- whether the host deployed private-only (a Premium host without a public IP)

## Usage

The module is executed by the Planton platform with a tfvars file converted from the manifest. To run it standalone, provide `metadata` and `spec` variables matching the generated `variables.tf`.

## Behavior Notes

- **The subnet must be named EXACTLY `AzureBastionSubnet`** and sized /26 or larger -- ARM validates the name at deploy time (the requirement lives on the referenced subnet, so no offline gate can check it).
- **The host binds its public IP EXCLUSIVELY** (Standard SKU, static allocation) -- never share the address with another consumer.
- **SKU upgrades are in-place; downgrades REPLACE the host** (the provider forces a new resource -- Azure has no downgrade path).
- **`kerberos_enabled` is applied at CREATE only** -- the provider has no update path for it and silently ignores later changes; plan Kerberos up front or replace the host to change it.
- **A Premium host without a public IP deploys PRIVATE-ONLY** (reachable only from connected networks) -- surfaced in the `private_only_enabled` output.
- **Billing starts at provisioning** (~10-minute creates): hourly per SKU, plus per-scale-unit on Standard/Premium. Developer is free.
- **The SKU/feature matrix is spec-validated** (features require Standard/Premium; session recording requires Premium; scale units are fixed at 2 below Standard) -- invalid combinations are rejected before any cloud call.
