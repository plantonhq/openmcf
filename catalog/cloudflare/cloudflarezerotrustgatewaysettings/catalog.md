# Cloudflare Zero Trust Gateway Settings

The account's Secure Web Gateway configuration: the settings panel behind every Gateway policy, the logging controls, and the PAC files -- three Cloudflare surfaces folded into one component.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Gateway configuration** -- one `cloudflare_zero_trust_gateway_settings` (only when `settings` is declared)
- **Gateway logging** -- one `cloudflare_zero_trust_gateway_logging` (only when `logging` is declared)
- **PAC files** -- one `cloudflare_zero_trust_gateway_pacfile` per `pacFiles` row

## Prerequisites

- **Zero Trust enabled on the account** (the team-name onboarding step)
- **A Cloudflare API token** with Account → Zero Trust → Edit
- **An activated Gateway certificate** BEFORE enabling `tlsDecrypt`, `fips.tls`, or deep `bodyScanning` (error 2211 otherwise)

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareZeroTrustGatewaySettings
metadata:
  name: acme-gateway
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  settings:
    activityLog:
      enabled: true
  logging:
    redactPii: true
```

```shell
planton apply -f gateway-settings.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex. |

### The settings tree (each sub-object optional; unset = unmanaged)

| Field | Description | Validation |
|-------|-------------|------------|
| `activityLog.enabled` | Master activity-log switch. | |
| `antivirus` | Upload/download scanning + fail_closed + user notification. | |
| `blockPage` | The branded block page (colors, text, mailto) or redirect. | `mode`: customized_block_page / redirect_uri; redirect needs `targetUri`. |
| `bodyScanning.inspectionMode` | deep or shallow. | Certificate-gated when deep. |
| `browserIsolation` | Non-identity + URL-triggered isolation. | |
| `certificate.id` | The TLS-inspection certificate (UUID; nil UUID = Cloudflare Root CA). | Must be ACTIVATED before tlsDecrypt. |
| `extendedEmailMatching.enabled` | Alias matching in policies. | |
| `fips.tls` | FIPS-approved TLS only. | Certificate-gated. |
| `hostSelector.enabled` | Host selectors in HTTP policies. | |
| `inspection.mode` | static or dynamic. | |
| `maxTtlSecs` | DNS TTL cap. | 60-36000. |
| `protocolDetection.enabled` | Non-standard-port protocol detection. | |
| `sandbox` | Sandboxed downloads + fallback allow/block. | |
| `tlsDecrypt.enabled` | THE TLS-inspection switch. | Error 2211 without an activated certificate. |

### logging

| Field | Description |
|-------|-------------|
| `redactPii` | Redact PII from activity logs. |
| `settingsByRuleType.{dns,http,l4}` | `logAll` / `logBlocks` per firewall type. The complete tree is always sent. |

### pacFiles rows

| Field | Description | Validation |
|-------|-------------|------------|
| `name` | The file's name. | Required; the module's per-row key. |
| `contents` | The PAC JavaScript body. | Required. |
| `slug` | The URL slug. | IMMUTABLE -- baked into the public URL. |
| `description` | Free-form. | |

## Destroy Semantics

Settings and logging: NO-OP (the live configuration stays as last applied). PAC files: real per-file delete.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `account_id` | string | The account applied to (the singleton's identity) |

## Related Components

- [Cloudflare Zero Trust Gateway Policy](/docs/catalog/cloudflare/cloudflarezerotrustgatewaypolicy) -- the filtering rules
- [Cloudflare Zero Trust DNS Location](/docs/catalog/cloudflare/cloudflarezerotrustdnslocation) -- where DNS traffic enters
