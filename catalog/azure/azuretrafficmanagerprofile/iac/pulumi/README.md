# AzureTrafficManagerProfile Pulumi Module

## Overview

Creates a Traffic Manager profile -- Azure's DNS-based traffic director. The profile owns a public DNS name (`{relative_name}.trafficmanager.net`) and answers each lookup with the address of one of its endpoints (`AzureTrafficManagerEndpoint` resources referencing this profile), chosen by routing method and endpoint health.

## Resources Created

- `network.TrafficManagerProfile` -- the profile with its DNS identity and health-probe configuration. (The SDK's `trafficmanager.Profile` is the deprecated legacy token for the same ARM object -- this module deliberately uses the current `network` token.)

## Outputs

- `traffic_manager_profile_id` -- the profile's ARM resource ID (what endpoints reference)
- `traffic_manager_profile_name` -- the profile's resource name
- `fqdn` -- the profile's public DNS name (what users resolve / your domain CNAMEs to)

## Behavior Notes

- **Traffic Manager is GLOBAL** -- the provider pins the ARM location to `"global"` itself; the spec carries no region.
- **`dns_config.relative_name` is globally unique across ALL of Azure** and FIXED at creation -- Azure rejects a taken name at apply time.
- **MultiValue routing requires `max_return`** (spec-validated, mirroring the provider's own create check); sent only when set.
- **The fast probe interval (10s) narrows the timeout to an explicit 5-9** (spec-validated, mirroring the provider). Probe cadence defaults (30/10/3) and the DNS TTL default (60, the portal's own) are always sent explicitly so both engines send identical wire shapes.
- **Everything except name, resource group, and the DNS relative name updates in place** -- including the routing method.
- **Billing**: per million DNS queries plus per-endpoint health-probe charges (fast-interval probes and Traffic View bill extra).

## Required Permissions

The deploying principal needs `Microsoft.Network/trafficManagerProfiles/*` (Network Contributor on the resource group covers it).
