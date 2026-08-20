# Cloudflare Zero Trust Device Posture Rule

A device posture rule: a health check WARP evaluates on enrolled devices, which Access and Gateway policies can require before admitting a device. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Posture rule** -- one `cloudflare_zero_trust_device_posture_rule`

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit
- For `*_s2s` types: a configured posture integration (its UUID goes in `input.connectionId`)

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustDevicePostureRule
metadata:
  name: macos-current
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: macos-current
  type: os_version
  match:
    - platform: mac
  input:
    version: 14.4.1
    operator: ">="
```

```shell
planton apply -f posture-rule.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex; replaces on change. |
| `name` | string | The rule's name. | Required (our tightening -- the provider marks it optional). |
| `type` | string | The check to run. | One of the provider's 23 types (enum-walled). |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `match` | list | Platform targeting. | `platform` in windows/mac/linux/android/ios/chromeos. |
| `input` | object | The check's parameters. | Enum walls per field (operators, compliance_status, state, auth_state, risk_level, network_status, operational_state, extended_key_usage, trust_stores); which fields a type reads is API-owned -- see field comments. |
| `schedule` | string | Polling frequency. | Default 5m; minimum 1m. |
| `expiration` | string | Result validity window. | Empty keeps results valid until overwritten. |

## Destroy Semantics

Real delete. Policies referencing the rule stop evaluating it -- review consumers first.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `rule_id` | string | The Cloudflare-assigned UUID policies reference |

## Related Components

- [Cloudflare Zero Trust Access Policy](/docs/catalog/cloudflare/cloudflarezerotrustaccesspolicy) -- Access-side consumers
- [Cloudflare Zero Trust Gateway Policy](/docs/catalog/cloudflare/cloudflarezerotrustgatewaypolicy) -- Gateway-side consumers
- [Cloudflare Zero Trust Device Default Profile](/docs/catalog/cloudflare/cloudflarezerotrustdevicedefaultprofile) -- the client the checks run under
