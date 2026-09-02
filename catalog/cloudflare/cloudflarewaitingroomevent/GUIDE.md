# CloudflareWaitingRoomEvent guide

The judgment this guide protects you from: scheduled events are an Enterprise advanced add-on (the room itself is Business+ -- a Free/Pro zone rejects rooms with code 1034, measured live), overrides are null-means-inherit, and the time rules are real walls. An unset field does not mean "zero" -- it means the room stays in charge. A window that violates the minute/five-minute floors fails at manifest time.

## Own cadence: the event is not a field of the room

Events are created and deleted per launch while the room persists. That is why they are their own kind. Hang many events on one room over a year; destroy each when the window is over. Destroying the room first leaves events pointing at a deleted object -- delete the events first, or accept the lookup failure.

`waiting_room_id` and `zone_id` are both required. The room id defaults to a `CloudflareWaitingRoom` reference; the zone id defaults to a `CloudflareDnsZone` reference. They must agree -- an event on room A in zone B fails at the API.

## Null means inherit

Every override is optional. Unset fields are never sent, so the room's value stays in charge during the window. Setting `new_users_per_minute: 400` without `total_active_users` is rejected -- Cloudflare requires the pair together or neither, because the two numbers are one piece of math.

The same inherit rule applies to `custom_page_html`, `disable_session_renewal`, `queueing_method`, `session_duration`, and the Turnstile fields. The room's Advanced entitlement still applies to Advanced overrides -- a `custom_page_html` on an event fails at the API if the zone lacks the add-on.

## RFC3339, and the clocks have floors

Times are RFC3339 (`2026-09-01T10:00:00Z`). The spec validates three rules up front so a bad window fails at manifest time, not as an opaque API error:

- `event_start_time` must be at least one minute before `event_end_time`
- `prequeue_start_time`, when set, must be at least five minutes before start
- `shuffle_at_event_start` requires a prequeue -- there is nothing to shuffle otherwise

Shuffle is the fairness tool for on-sales where arrival time is a bot advantage. It only makes sense when the event's `queueing_method` respects order (`fifo`). Shuffling a `random` queue is wasted work.

## Destroy is a real delete

Destroy removes the event. The room is untouched. `suspended: true` lets the window come and go with no effect -- use that for a cancelled launch you may reschedule.

## Pairs well with

- [CloudflareWaitingRoom](../cloudflarewaitingroom/README.md) -- the room this event overrides; create it first and wire `waiting_room_id` via `value_from`.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
