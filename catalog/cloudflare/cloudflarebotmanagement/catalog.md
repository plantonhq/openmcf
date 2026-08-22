# Cloudflare Bot Management

Manages a zone's Bot Management configuration -- the singleton switchboard for automated traffic, from free Bot Fight Mode through Super Bot Fight Mode and Enterprise knobs to AI-crawler controls. Unset fields are not managed. Destroy abandons the last-applied values; Cloudflare has no delete for this surface.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Bot Management configuration** -- one `cloudflare_bot_management` on the zone (a zone singleton; Cloudflare's identity for this surface IS the zone ID)

## Prerequisites

- **A Cloudflare zone** -- typically a CloudflareDnsZone whose `zone_id` output this resource references
- **A Cloudflare API token** with Zone → Bot Management → Edit
- **The right zone plan for gated fields** -- `fight_mode` is free; `sbfm_*` needs Pro/Business; Enterprise knobs need Enterprise. Setting a field the plan does not include fails at the API, and non-entitled zones omit those fields from responses (refresh drift)

## Quick Start

Turn on Bot Fight Mode and touch nothing else:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareBotManagement
metadata:
  name: fight-mode
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  fightMode: true
```

```shell
planton apply -f bot-management.yaml
```

Every other Bot Management field stays exactly as it was -- unset means unmanaged. Destroy will not turn Fight Mode off; apply `fightMode: false` first if you want it off.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone whose Bot Management configuration is managed. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. At least one setting must also be configured -- a resource that manages nothing is rejected. |

### Optional Fields

All settings are optional and default to **not managed** -- the module never sends a field you did not set.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `fightMode` | bool | not managed | Bot Fight Mode (free). Mutually exclusive with Super Bot Fight Mode at Cloudflare. |
| `sbfmDefinitelyAutomated` | string | not managed | SBFM action on definitely-automated traffic: `allow`, `block`, `managed_challenge`. Pro+. |
| `sbfmLikelyAutomated` | string | not managed | SBFM action on likely-automated traffic. Business+. |
| `sbfmVerifiedBots` | string | not managed | SBFM action on verified bots: `allow` or `block`. Blocking verified bots blocks search indexing. |
| `sbfmStaticResourceProtection` | bool | not managed | Extend SBFM to static assets. Can catch legitimate hotlinked traffic. |
| `optimizeWordpress` | bool | not managed | Tune SBFM for WordPress. |
| `autoUpdateModel` | bool | not managed | Enterprise: adopt new bot-detection ML models automatically. |
| `suppressSessionScore` | bool | not managed | Enterprise: stop tracking the session's highest bot score. |
| `enableJs` | bool | not managed | Enterprise: run invisible JavaScript detections. |
| `bmCookieEnabled` | bool | not managed | Enterprise: allow the Bot Management cookie (Cloudflare defaults this to true). |
| `aiBotsProtection` | string | not managed | `block`, `disabled`, or `only_on_ad_pages`. |
| `crawlerProtection` | string | not managed | `enabled` or `disabled`. |
| `contentBotsProtection` | string | not managed | `block` or `disabled`. |
| `cfRobotsVariant` | string | not managed | `off` or `policy_only`. |
| `isRobotsTxtManaged` | bool | not managed | Serve a Cloudflare-managed robots.txt prepended to the origin's. |

## Examples

### Bot Fight Mode (free)

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareBotManagement
metadata:
  name: fight-mode
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  fightMode: true
```

### Super Bot Fight Mode (Pro/Business)

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareBotManagement
metadata:
  name: sbfm
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  sbfmDefinitelyAutomated: managed_challenge
  sbfmVerifiedBots: allow
```

## Destroy Semantics

Destroy is a NO-OP at Cloudflare. The state entry disappears; the live configuration stays exactly as last applied. To retire a setting, apply its off value first, then destroy. Non-entitled zones omit gated fields from API responses -- a manifest that set those fields will refresh as drift.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `zone_id` | string | The zone whose Bot Management configuration is managed. The configuration is a zone singleton -- the zone ID is its identity. |

## Related Components

- [Cloudflare IP Access Rule](/docs/catalog/cloudflare/cloudflareipaccessrule) -- a static IP/ASN/country decision
- [Cloudflare Ruleset](/docs/catalog/cloudflare/cloudflareruleset) -- skip rules for verified-bot exceptions
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key; the zone plan gates SBFM and Enterprise fields
