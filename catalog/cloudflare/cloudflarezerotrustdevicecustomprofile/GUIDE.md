# CloudflareZeroTrustDeviceCustomProfile guide

The judgment this guide protects you from: profile assignment is an ordering competition, and the settings body carries the same replacement semantics as the default profile.

## Precedence is a fleet-wide ordering, manage it like one

Every custom profile on the account competes in one precedence space; the LOWEST value wins for a matching device. Leave gaps (100, 200, 300...) so a future profile can slot between two existing ones without renumbering, and record the account's precedence plan somewhere the next operator will read -- two profiles with adjacent values and overlapping matches is how a device lands in the wrong split tunnel.

## The match expression fails open to the default profile

A device matching NO custom profile gets the default profile -- so a typo in `match` does not strand devices, it silently baselines them. After changing a match expression, verify a known device actually landed on this profile (the Zero Trust dashboard shows the applied profile per device) rather than trusting the apply.

## The fallback list here is per-profile, and still full-replacement

Unlike the account-default list, this profile's fallback rows apply only to its matched devices -- but the same replacement rule holds: the declared list IS the list. Rows ride the profile; deleting the profile retires them with it, so there is no separate cleanup.

## Deleting a profile is a routing change, not just a removal

Matched devices immediately fall back to the default profile (or the next-lowest-precedence match). If this profile carved lab networks out of the tunnel, its deletion sends that traffic THROUGH the tunnel on every matched device at once. Treat profile deletion with the same care as profile creation.

## Pairs well with

- [CloudflareZeroTrustDeviceDefaultProfile](../cloudflarezerotrustdevicedefaultprofile/README.md) -- the baseline this profile overrides.
- [CloudflareZeroTrustDevicePostureRule](../cloudflarezerotrustdeviceposturerule/README.md) -- health checks that can gate what these devices reach.
- [CloudflareZeroTrustTunnelVirtualNetwork](../cloudflarezerotrusttunnelvirtualnetwork/README.md) -- the networks `virtual_networks` scopes device access to.
