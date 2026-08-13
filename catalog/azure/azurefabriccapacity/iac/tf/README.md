# AzureFabricCapacity Terraform Module

## Overview

Creates a Microsoft Fabric capacity -- the billing and compute anchor of Microsoft Fabric. Workspaces assign themselves to the capacity; its F-SKU sets how much compute every workload on it shares.

## Resources Created

- `azurerm_fabric_capacity` -- the capacity (F-SKU, administrators, tags)

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureFabricCapacitySpec fields; the resource group reference arrives as a resolved literal

## Outputs

- `fabric_capacity_id` -- the capacity's ARM resource ID
- `fabric_capacity_name` -- the capacity's name (what Fabric workspaces assign themselves to)

## Behavior Notes

- **A running capacity bills PER HOUR from the moment it exists** (F2 is the smallest; F2048 is a thousand times that). The SKU scales up and down IN PLACE -- start small.
- **The SKU tier is always "Fabric"** (Azure's only value at v5) -- deliberately not part of the spec; the module sends it explicitly.
- **At least one administrator is required** -- Azure rejects a capacity created with none (the spec enforces the list non-empty at all times; clearing every administrator on a running paid capacity is a lockout).
- **This is azurerm's entire Fabric surface**: workspaces and items live in Microsoft's dedicated `fabric` Terraform provider, the Fabric portal, or its APIs.
- **Create and delete are polled long-running operations** (30-minute provider timeouts).
