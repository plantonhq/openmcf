# AzureTrafficManagerEndpoint Terraform Module

## Overview

Creates one Traffic Manager endpoint -- a destination the referenced profile (`AzureTrafficManagerProfile`) steers traffic to. The endpoint type is whichever variant block the spec carries (validation guarantees exactly one): a public Azure resource by ARM ID (`azure`), a DNS name or IP address (`external`), or another Traffic Manager profile composing routing trees (`nested`).

## Resources Created

Exactly one of:

- `azurerm_traffic_manager_azure_endpoint` -- when `spec.azure` is set
- `azurerm_traffic_manager_external_endpoint` -- when `spec.external` is set
- `azurerm_traffic_manager_nested_endpoint` -- when `spec.nested` is set

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureTrafficManagerEndpointSpec fields; the profile reference and target references arrive as resolved literals

## Outputs

- `endpoint_id` -- the endpoint's ARM resource ID (`{profile_id}/{TYPE}/{name}`, coalesced across the three variants)
- `endpoint_name` -- the endpoint's name within its profile

## Behavior Notes

- **Which shared fields matter depends on the PROFILE's routing method** -- weight steers Weighted profiles, priority steers Priority profiles, geo_mappings steer Geographic profiles (every code claimed by exactly one endpoint), subnets steer Subnet profiles (no overlaps) -- all evaluated by Azure at apply time.
- **Priority is sent only when set** -- unset lets Azure assign the next free value in creation order (the service owns that default). Weight defaults to 1 and is always sent.
- **`endpoint_location` (external/nested only) is REQUIRED by the service under Performance routing** -- external targets carry no discoverable region; azure endpoints derive theirs from the target resource and have no location argument at all.
- **Nested endpoints carry no always-serve switch** -- the provider exposes none (child health is the point of nesting); azure/external endpoints default it off.
- **Endpoints carry NO ARM tags on any engine** -- the platform's derived tags land on the owning profile instead.
- **Name, profile, type, and subnet claims are fixed at creation**; everything else (target, weight, priority, enabled, geo claims, headers) updates in place.
- **Endpoints are free at rest** -- each billable meter (health probes, queries) lives on the profile.

## Required Permissions

The deploying principal needs `Microsoft.Network/trafficManagerProfiles/*` (endpoints are profile children; Network Contributor on the resource group covers it).
