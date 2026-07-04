---
title: "Private Service Connect to Google APIs"
description: "The non-load-balancer face of the global forwarding rule: with `loadBalancingScheme: NONE` and the literal target `all-apis`, the rule becomes a Private Service Connect endpoint — VPC workloads reach..."
type: "preset"
rank: "03"
presetSlug: "03-psc-google-apis"
componentSlug: "global-forwarding-rule"
componentTitle: "Global Forwarding Rule"
provider: "gcp"
icon: "package"
order: 3
---

# Private Service Connect to Google APIs

The non-load-balancer face of the global forwarding rule: with `loadBalancingScheme: NONE` and the literal target `all-apis`, the rule becomes a Private Service Connect endpoint — VPC workloads reach Google APIs over an internal IP without touching the public internet.

## When to Use

- VPC Service Controls perimeters (use target `vpc-sc` to restrict to VPC-SC-supported APIs)
- Locked-down networks with no default internet route that still need Google APIs
- Keeping API traffic on private paths for compliance

## Remix Notes

- The name doubles as the DNS/service-directory handle and is limited to 20 characters (letters and digits) for PSC-for-Google-APIs rules.
- `ipAddress` must be an INTERNAL `GcpGlobalAddress` with `purpose: PRIVATE_SERVICE_CONNECT` in the same VPC — reference both via `valueFrom`.
- Set `noAutomateDnsZone: true` only when you manage the `googleapis.com` private DNS zone yourself; by default GCP creates it for you.
- For producer services (not Google APIs), point `target` at the producer's service attachment URI instead — the `pscConnectionStatus` output shows whether the producer ACCEPTED the connection.
