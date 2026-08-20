# CloudflareZeroTrustDeviceDefaultProfile guide

The judgment this guide protects you from: this is the ACCOUNT-WIDE default for every enrolled device, all three of its surfaces survive destroy, and one of them replaces a list Cloudflare seeded for you.

## Destroy reverts nothing, on any surface

The profile is a settings singleton (PATCH upsert, no delete), the fallback list's destroy is a no-op, and the certificate toggle has no delete at all. Removing this resource abandons the last-applied state -- devices keep behaving as configured. To actually revert, apply the values you want; never assume a destroy cleaned up.

## The fallback list is full-replacement -- and Cloudflare pre-seeded it

Every account ships with a default local-domain fallback list (localhost, home.arpa, and kin). The moment `fallback_domains` is declared, YOUR list replaces the whole thing -- including those seeded rows. Fetch the current list first (`GET accounts/{id}/devices/policy/fallback_domains`), fold the rows you want to keep into the manifest, and only then apply. Declaring one row does not add a row; it deletes everything else.

## switch_locked and allowed_to_leave are lock-in levers

`switch_locked: true` removes users' ability to turn WARP off; `allowed_to_leave: false` removes their ability to unenroll. Together they hard-lock the fleet. Roll them out after split-tunnel and fallback routing are proven, not before -- a bad route with a locked switch is a helpdesk incident on every device at once.

## The certificate toggle is zone-scoped and permanent

`zone_certificates` enables WARP client certificate provisioning on ONE zone (the one surface here that is not account-scoped). Cloudflare offers no delete and no import for it: to turn it off, apply `enabled: false` -- removing the block leaves it however it last was.

## exclude and include are modes, not filters to combine

A profile runs in exclude mode (everything tunnels except the list) or include mode (nothing tunnels except the list). The spec enforces the exclusivity; when switching modes, remember the semantics of every existing entry invert.

## Pairs well with

- [CloudflareZeroTrustDeviceCustomProfile](../cloudflarezerotrustdevicecustomprofile/README.md) -- group-specific overrides of this baseline.
- [CloudflareZeroTrustDevicePostureRule](../cloudflarezerotrustdeviceposturerule/README.md) -- the health checks Access and Gateway policies demand from these devices.
- [CloudflareZeroTrustTunnelVirtualNetwork](../cloudflarezerotrusttunnelvirtualnetwork/README.md) -- the networks `virtual_networks` scopes device access to.
