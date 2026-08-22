# CloudflareWaitingRoom guide

The judgment this guide protects you from: the 200 floors are Cloudflare's, Advanced fields fail at apply on a plan without the add-on, and `bypass_rules` is the room's entire rule list -- destroy empties it, including rules you did not create.

## Floors of 200, session 1–30

`new_users_per_minute` and `total_active_users` have a Cloudflare floor of 200. A room that admits 50 people per minute is not a waiting room this API will create -- validation rejects it here so you do not discover that at apply. `session_duration` is 1–30 minutes (Cloudflare's default is 5).

`queue_all: true` is the "close the shop" switch: everyone queues regardless of thresholds. Use it for a maintenance window instead of deleting the room.

## Advanced fields fail at the API, not here

These need the Waiting Rooms Advanced add-on: `additional_routes`, `custom_page_html`, `disable_session_renewal`, `json_response_enabled`, a non-default `queueing_method`, `turnstile_action: infinite_queue`, and any `turnstile_mode` other than `off`. The entitlement wall is Cloudflare's. Setting one of these on a plan without the add-on fails the apply.

The free-plan-safe shape is host + path + 200/200 + `fifo` + default Turnstile. That is the preset and the e2e live arm. Advanced fields live in `e2e/advanced-plan.yaml` for offline plan only.

## bypass_rules is the WHOLE table

Cloudflare stores waiting-room rules as a separate per-room list with full-replacement updates. This kind folds that list into the room. Every apply replaces the room's entire rule set with exactly `bypass_rules`. Rules added in the dashboard disappear on the next apply.

The action is fixed to `bypass_waiting_room` -- the module supplies it so manifests never repeat a constant. `enabled` defaults true (unlike snippet rules).

Destroy PUTs `[]` and then deletes the room. Dashboard rules, API rules, rules from a previous manifest -- gone. If you need to stop bypassing without deleting the room, set each rule's `enabled: false` and apply.

## Events are a different kind

A scheduled launch window that overrides rates, HTML, or Turnstile is `CloudflareWaitingRoomEvent`, not a field of the room. Events live on their own cadence -- created and deleted per launch while the room persists. Wire `waiting_room_id` via `value_from`.

## Pairs well with

- [CloudflareWaitingRoomEvent](../cloudflarewaitingroomevent/README.md) -- the scheduled override window for a launch or on-sale.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`.
