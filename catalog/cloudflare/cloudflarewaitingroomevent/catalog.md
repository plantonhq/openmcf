# Cloudflare Waiting Room Event

A scheduled window on a waiting room during which optional overrides temporarily replace the room's settings. Events live on their own cadence -- created and deleted per launch while the room persists. Unset overrides inherit the room. Times are RFC3339; start must be at least one minute before end.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Waiting room event** -- one `cloudflare_waiting_room_event` on the named room and zone

## Prerequisites

- **A CloudflareWaitingRoom** -- this event's `waitingRoomId` foreign key; create the room first
- **The same zone** the room belongs to -- `zoneId` is required alongside the room id
- **A Cloudflare API token** with Zone → Waiting Rooms → Edit
- **A window that satisfies the time rules** -- start ≥ 1 minute before end; prequeue ≥ 5 minutes before start if set

## Quick Start

A launch window with no overrides -- the room's settings stay in charge:

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWaitingRoomEvent
metadata:
  name: launch
  org: acme-corp
  env: prod
spec:
  waitingRoomId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  zoneId:
    value: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  name: launch
  eventStartTime: "2026-09-01T10:00:00Z"
  eventEndTime: "2026-09-01T14:00:00Z"
```

```shell
planton apply -f waiting-room-event.yaml
```

The room queues with its own thresholds during this window. Add overrides only when the launch needs different math.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `waitingRoomId` | StringValueOrRef | The room the event runs on. Can reference a CloudflareWaitingRoom via `valueFrom` (defaults to `status.outputs.waiting_room_id`). | Required. |
| `zoneId` | StringValueOrRef | The zone the room belongs to. Can reference a CloudflareDnsZone via `valueFrom` (defaults to `status.outputs.zone_id`). | Required. |
| `name` | string | Display name. | Required, min length 1. |
| `eventStartTime` | string | When the window opens. | Required. RFC3339. At least one minute before `eventEndTime`. |
| `eventEndTime` | string | When the window closes. | Required. RFC3339. |

### Optional Fields

Every override is **null-means-inherit** -- unset leaves the room's value in charge.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `prequeueStartTime` | string | unset (no prequeue) | RFC3339. At least five minutes before start. |
| `shuffleAtEventStart` | bool | unset | Randomize prequeued users at the door. Requires a prequeue. |
| `description` | string | unset | Dashboard note. |
| `suspended` | bool | unset | The window comes and goes with no effect. |
| `customPageHtml` | string | inherit | Override the room's queue page (Advanced entitlement applies). |
| `disableSessionRenewal` | bool | inherit | Override the room's session-renewal setting. |
| `newUsersPerMinute` | int32 | inherit | Override the room's rate. Floor 200. Set together with `totalActiveUsers` (both or neither). |
| `totalActiveUsers` | int32 | inherit | Override the room's cap. Floor 200. Set together with `newUsersPerMinute`. |
| `queueingMethod` | string | inherit | `fifo`, `random`, `passthrough`, `reject`. |
| `sessionDuration` | int32 | inherit | 1–30 minutes. |
| `turnstileAction` | string | inherit | `log` or `infinite_queue`. |
| `turnstileMode` | string | inherit | `off`, `invisible`, `visible_non_interactive`, `visible_managed`. |

## Examples

### Launch window, no overrides

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWaitingRoomEvent
metadata:
  name: launch
  org: acme-corp
  env: prod
spec:
  waitingRoomId:
    valueFrom:
      kind: CloudflareWaitingRoom
      name: checkout
      fieldPath: status.outputs.waiting_room_id
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  name: launch
  eventStartTime: "2026-09-01T10:00:00Z"
  eventEndTime: "2026-09-01T14:00:00Z"
```

### Prequeue with a shuffled door and a users-pair override

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWaitingRoomEvent
metadata:
  name: on-sale
  org: acme-corp
  env: prod
spec:
  waitingRoomId:
    valueFrom:
      kind: CloudflareWaitingRoom
      name: checkout
      fieldPath: status.outputs.waiting_room_id
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: example-zone
      fieldPath: status.outputs.zone_id
  name: on-sale
  eventStartTime: "2026-09-01T10:00:00Z"
  eventEndTime: "2026-09-01T14:00:00Z"
  prequeueStartTime: "2026-09-01T09:30:00Z"
  shuffleAtEventStart: true
  newUsersPerMinute: 400
  totalActiveUsers: 2000
```

## Destroy Semantics

Destroy is a real delete. The event is removed; the room is untouched and returns to its own settings. Set `suspended: true` to cancel a window without losing the ID.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `event_id` | string | The created event's ID |
| `waiting_room_id` | string | The room the event runs on |
| `zone_id` | string | The zone the waiting room belongs to |

## Related Components

- [Cloudflare Waiting Room](/docs/catalog/cloudflare/cloudflarewaitingroom) -- the room this event overrides; create it first
- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- `zoneId` foreign key
