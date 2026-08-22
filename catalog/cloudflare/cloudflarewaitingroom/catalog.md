# Cloudflare Waiting Room

A virtual queue in front of a host+path that admits visitors at a controlled rate. Bypass rules ride along in the same manifest and replace the room's entire rule list on every apply. Advanced add-on fields fail at the API on plans without the add-on.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Waiting room** -- one `cloudflare_waiting_room` on the zone
- **Bypass rules** -- one `cloudflare_waiting_room_rules` when `bypass_rules` is non-empty; the list is the room's entire rule set (action fixed to `bypass_waiting_room`)

## Prerequisites

- **A Cloudflare zone** -- typically a CloudflareDnsZone whose `zone_id` output this resource references
- **A Cloudflare API token** with Zone → Waiting Room → Edit
- **Thresholds at or above 200** -- Cloudflare's floor for `newUsersPerMinute` and `totalActiveUsers`
- **The Waiting Rooms Advanced add-on** -- only if you set Advanced fields (`additionalRoutes`, `customPageHtml`, non-default `queueingMethod`, non-off Turnstile, and kin)

## Quick Start

A basic checkout queue at the 200/200 floor:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWaitingRoom
metadata:
  name: checkout-queue
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: checkout-queue
  host: shop.example.com
  path: /checkout
  newUsersPerMinute: 200
  totalActiveUsers: 200
```

```shell
planton apply -f waiting-room.yaml
```

A scheduled launch window that overrides these rates is a separate `CloudflareWaitingRoomEvent`, not a field of the room.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `zoneId` | StringValueOrRef | The zone the waiting room belongs to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `name` | string | Display name. | Required, min length 1. |
| `host` | string | Host the room protects. | Required. |
| `newUsersPerMinute` | int32 | New admissions per minute before queueing. | Required. Floor 200. |
| `totalActiveUsers` | int32 | Concurrent users on the route before queueing. | Required. Floor 200. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `path` | string | `/` | Path the room protects. |
| `sessionDuration` | int32 | 5 | Minutes a slot stays valid (1–30). |
| `suspended` | bool | unset | Pause queueing without deleting the room. |
| `queueAll` | bool | unset | Send everyone to the queue regardless of thresholds. |
| `queueingMethod` | string | `fifo` | `fifo`, `random`, `passthrough`, `reject`. Non-default needs Advanced. |
| `queueingStatusCode` | int32 | 200 | `200`, `202`, or `429`. |
| `cookieAttributes` | object | unset | `{samesite, secure}` of the `__cfwaitingroom` cookie. |
| `cookieSuffix` | string | unset | Suffix so several rooms on one host do not collide cookies. |
| `customPageHtml` | string | unset | Custom queue page (Advanced). |
| `defaultTemplateLanguage` | string | `en-US` | Language of Cloudflare's default page (38-language wall). |
| `description` | string | unset | Shown in the dashboard. |
| `disableSessionRenewal` | bool | unset | Expire the session from entry, not last request (Advanced). |
| `jsonResponseEnabled` | bool | unset | JSON queue response for API clients (Advanced). |
| `additionalRoutes` | object[] | empty | Extra host+path covered by the same quota (Advanced). |
| `enabledOriginCommands` | string[] | empty | Currently only `revoke`. |
| `turnstileAction` | string | `log` | `log` or `infinite_queue` (Advanced). |
| `turnstileMode` | string | `invisible` | `off`, `invisible`, `visible_non_interactive`, `visible_managed`. Non-off needs Advanced. |
| `bypassRules` | object[] | empty | The room's entire bypass-rule table. Each rule needs `expression`. Action is fixed to `bypass_waiting_room`. |

## Examples

### Basic checkout queue

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWaitingRoom
metadata:
  name: checkout-queue
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: checkout-queue
  host: shop.example.com
  path: /checkout
  newUsersPerMinute: 200
  totalActiveUsers: 200
  bypassRules:
    - expression: 'ip.src in {203.0.113.0/24}'
      description: Office network skips the queue
```

## Destroy Semantics

Destroy is a real delete of the room. The folded bypass-rules resource PUTs an empty list first, which wipes rules added outside this manifest too. Events that referenced `waiting_room_id` start failing their lookup. Suspend or `queueAll` is the reversible alternative when you want the ID to stay.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `waiting_room_id` | string | The created room's ID -- what events and the import recipe reference |
| `zone_id` | string | The zone the waiting room belongs to |

## Related Components

- [Cloudflare Waiting Room Event](/docs/catalog/cloudflare/cloudflarewaitingroomevent) -- a scheduled window that temporarily overrides this room's rates
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
