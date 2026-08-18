# AzureTrafficManagerEndpoint Pulumi Module

## Overview

Creates one Traffic Manager endpoint -- a destination the referenced profile (`AzureTrafficManagerProfile`) steers traffic to. The endpoint type is whichever variant block the spec carries (validation guarantees exactly one): a public Azure resource by ARM ID (`azure`), a DNS name or IP address (`external`), or another Traffic Manager profile composing routing trees (`nested`).

## Resources Created

Exactly one of `network.TrafficManagerAzureEndpoint`, `network.TrafficManagerExternalEndpoint`, or `network.TrafficManagerNestedEndpoint`.

## Outputs

- `endpoint_id` -- the endpoint's ARM resource ID (`{profile_id}/{TYPE}/{name}`)
- `endpoint_name` -- the endpoint's name within its profile

## Behavior Notes

- **Which shared fields matter depends on the PROFILE's routing method** -- weight (Weighted), priority (Priority), geo_mappings (Geographic; every code claimed by exactly one endpoint), subnets (Subnet; no overlaps) -- all evaluated by Azure at apply time.
- **Priority is sent only when set** -- unset lets Azure assign the next free value in creation order. Weight defaults to 1 and is always sent, so both engines send identical wire shapes.
- **`endpoint_location` (external/nested only) is REQUIRED by the service under Performance routing**; azure endpoints derive their region from the target resource and carry no location argument.
- **Nested endpoints carry no always-serve switch** -- the provider exposes none (child health is the point of nesting).
- **The subnet/header builders are typed per variant** -- the SDK generates a distinct Go type per resource for structurally identical blocks, hence three small builders instead of one.
- **Endpoints carry NO ARM tags on any engine** and are free at rest -- billable meters live on the profile.
