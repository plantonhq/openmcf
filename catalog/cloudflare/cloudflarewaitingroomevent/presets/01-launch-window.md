# Launch window

A waiting-room event with a name and a time window, and no overrides -- the room's settings stay in charge for the launch. Times are RFC3339; start must be at least one minute before end.

## When to Use

- A product launch or on-sale that should use the room's existing thresholds
- First event on a room
- A window you will delete after the date (events live on their own cadence)

## Key Configuration Choices

- **No overrides** -- unset means inherit. Add `new_users_per_minute` and `total_active_users` together (both or neither) only when the launch needs different math.
- **No prequeue** -- `shuffle_at_event_start` requires a prequeue. Add `prequeue_start_time` at least five minutes before start if you want an early line.
- **waiting_room_id and zone_id** -- both required, and they must agree. Prefer `value_from` once the room is a Planton resource.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `waiting_room_id.value` | The waiting room this event runs on | A CloudflareWaitingRoom's `status.outputs.waiting_room_id`, or the dashboard room id |
| `zone_id.value` | The zone the room belongs to | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `event_start_time` | When the window opens (RFC3339) | Your launch schedule |
| `event_end_time` | When the window closes (RFC3339) | Your launch schedule |

## Related Presets

None on this kind -- pair with the `CloudflareWaitingRoom` **01-basic-room** preset for the room this event overrides.
