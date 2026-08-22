# Cloudflare Zero Trust Device Default Profile

The account's default WARP device profile: the settings every enrolled device receives unless a custom profile matches it first. Applying edits the always-existing singleton; destroy is a no-op.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Default device profile** -- one `cloudflare_zero_trust_device_default_profile` (a PATCH upsert of the account singleton)
- **Local-DNS fallback list** -- one `cloudflare_zero_trust_device_default_profile_local_domain_fallback` when `fallback_domains` is declared (full replacement of the account's list)
- **Zone certificate toggle** -- one `cloudflare_zero_trust_device_default_profile_certificates` when `zone_certificates` is declared (zone-scoped)

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit (the zone-certificates fold additionally needs Zone → SSL and Certificates → Edit)

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustDeviceDefaultProfile
metadata:
  name: default-device-profile
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  switchLocked: true
  allowedToLeave: false
```

```shell
planton apply -f default-profile.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `exclude` / `include` | list | Split-tunnel routes. | Mutually exclusive lists; each entry address XOR host. |
| `serviceModeV2` | object | The client mode. | `mode` required when set; `port` 1-65535. |
| `virtualNetworks` | object | Reachable Zero Trust virtual networks. | `allowed` min 1 (vnet references); `defaultVirtualNetworkId` required. |
| `dnsSearchSuffixes` | list | Short-hostname resolution suffixes. | `suffix` required per row. |
| `fallbackDomains` | list | The folded local-DNS fallback list. | FULL REPLACEMENT -- the declared list is exactly what exists. |
| `zoneCertificates` | object | Per-zone WARP client certificates. | `zoneId` reference + explicit `enabled` (no delete exists). |
| `autoConnect`, `captivePortal` | int | Reconnect/captive-portal timers (seconds). | >= 0. |
| `lanAllowMinutes`, `lanAllowSubnetSize` | int | The LAN access window. | >= 0; subnet size <= 128. |
| `switchLocked`, `allowedToLeave`, `allowModeSwitch`, `allowUpdates`, ... | bool | The toggle body. | Unset never sends (keeps Cloudflare's default). |
| `tunnelProtocol` | string | wireguard or masque. | API-owned value set; empty keeps the account default. |

## Destroy Semantics

No-op on all three surfaces: the last-applied values stand. Removing the resource abandons management, never reverts.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `account_id` | string | The account (the singleton's identity) |
| `gateway_unique_id` | string | The Gateway-side profile identifier |
| `policy_id` | string | The profile's policy identifier |

## Related Components

- [Cloudflare Zero Trust Device Custom Profile](/docs/catalog/cloudflare/cloudflarezerotrustdevicecustomprofile) -- targeted overrides
- [Cloudflare Zero Trust Device Posture Rule](/docs/catalog/cloudflare/cloudflarezerotrustdeviceposturerule) -- device health checks
