# AzureTrafficManagerProfile Terraform Module

## Overview

Creates a Traffic Manager profile -- Azure's DNS-based traffic director. The profile owns a public DNS name (`{relative_name}.trafficmanager.net`) and answers each lookup with the address of one of its endpoints (`AzureTrafficManagerEndpoint` resources referencing this profile), chosen by routing method and endpoint health. Because the steering happens in DNS, Traffic Manager fronts anything with a resolvable address -- across regions, clouds, and on-premises -- and is never in the data path.

## Resources Created

- `azurerm_traffic_manager_profile` -- the profile with its DNS identity and health-probe configuration

## Variables

The generated `variables.tf` mirrors the proto contract:

- `metadata` -- Planton resource metadata (name, org, env, labels, tags)
- `spec` -- the AzureTrafficManagerProfileSpec fields; the resource group reference arrives as a resolved literal

## Outputs

- `traffic_manager_profile_id` -- the profile's ARM resource ID (what endpoints reference)
- `traffic_manager_profile_name` -- the profile's resource name
- `fqdn` -- the profile's public DNS name (what users resolve / your domain CNAMEs to)

## Behavior Notes

- **Traffic Manager is GLOBAL** -- the provider pins the ARM location to `"global"` itself; the spec carries no region.
- **`dns_config.relative_name` is globally unique across ALL of Azure** (the shared trafficmanager.net namespace) and FIXED at creation -- Azure rejects a taken name at apply time.
- **MultiValue routing requires `max_return`** (spec-validated, mirroring the provider's own create check); the module sends it only when set.
- **The fast probe interval (10s) narrows the timeout to an explicit 5-9** (spec-validated, mirroring the provider). Probe cadence defaults (30/10/3) are always sent explicitly so plans stay deterministic.
- **Everything except name, resource group, and the DNS relative name updates in place** -- including the routing method, which re-steers traffic without touching endpoints.
- **Billing**: per million DNS queries plus per-endpoint health-probe charges (fast-interval probes and Traffic View bill extra). The profile object itself is cheap at rest.

## Required Permissions

See [`../permissions.yaml`](../permissions.yaml) for the least-privilege actions the deploying principal needs.
