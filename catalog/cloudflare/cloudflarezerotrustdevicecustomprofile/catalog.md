# Cloudflare Zero Trust Device Custom Profile

A targeted WARP device profile: the default profile's settings body, applied only to devices matched by a wirefilter expression. Lowest precedence wins; deleting the profile returns its devices to the default.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Custom device profile** -- one `cloudflare_zero_trust_device_custom_profile`
- **Per-profile fallback list** -- one `cloudflare_zero_trust_device_custom_profile_local_domain_fallback` when `fallbackDomains` is declared (full replacement of this profile's list)

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustDeviceCustomProfile
metadata:
  name: contractor-devices
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: contractor-devices
  match: identity.groups.name == "contractors"
  precedence: 100
```

```shell
planton apply -f custom-profile.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | The profile's name. | Required. |
| `match` | string | The wirefilter device selector. | Required; selectors: identity.email, identity.groups.*, identity.service_token_uuid, identity.saml_attributes, network, os.name, os.version. |
| `precedence` | int | Evaluation order; LOWER wins. | Required, >= 1 (our tightening -- the API rejects a create without it). |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `enabled` | bool | Whether the profile applies. | Unset keeps Cloudflare's default (on). |
| `exclude` / `include` | list | Split-tunnel routes. | Mutually exclusive lists; each entry address XOR host. |
| `serviceModeV2` | object | The client mode. | `mode` required when set; `port` 1-65535. |
| `virtualNetworks` | object | Reachable Zero Trust virtual networks. | `allowed` min 1 (vnet references); `defaultVirtualNetworkId` required. |
| `dnsSearchSuffixes` | list | Short-hostname resolution suffixes. | `suffix` required per row. |
| `fallbackDomains` | list | The folded per-profile fallback list. | FULL REPLACEMENT -- the declared list is exactly what exists. |
| `lanAllowMinutes`, `lanAllowSubnetSize` | int | The LAN access window. | >= 0; subnet size <= 128. |
| `switchLocked`, `allowedToLeave`, ... | bool | The toggle body. | Unset never sends (keeps Cloudflare's default). |

## Destroy Semantics

Real delete: matched devices fall back to the default profile; the per-profile fallback rows retire with the profile.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `policy_id` | string | The Cloudflare-assigned profile identifier |
| `gateway_unique_id` | string | The Gateway-side profile identifier |

## Related Components

- [Cloudflare Zero Trust Device Default Profile](/docs/catalog/cloudflare/cloudflarezerotrustdevicedefaultprofile) -- the account baseline
- [Cloudflare Zero Trust Device Posture Rule](/docs/catalog/cloudflare/cloudflarezerotrustdeviceposturerule) -- device health checks
