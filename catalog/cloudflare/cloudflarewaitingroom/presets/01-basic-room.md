# Basic checkout queue

A free-plan-safe waiting room on `shop.example.com/checkout` at Cloudflare's 200/200 floors, fifo admission, default Turnstile. No Advanced fields -- those fail at apply on a zone without the add-on.

## When to Use

- First waiting room on a zone
- A checkout or ticket path that must not melt the origin
- A starting point before you add bypass rules or a `CloudflareWaitingRoomEvent`

## Key Configuration Choices

- **200/200 floors** -- Cloudflare will not create a room below these. Raise them; do not lower them.
- **No Advanced fields** -- `additional_routes`, `custom_page_html`, non-fifo `queueing_method`, and non-off Turnstile need the Waiting Rooms Advanced add-on.
- **No bypass_rules** -- add them in this same manifest when an office CIDR should skip the queue. The list is the room's entire rule set.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone the room belongs to | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |
| `host` | The host the room protects | Your public hostname (e.g. shop.example.com) |
| `path` | The path the room protects | The subtree that should queue (e.g. /checkout) |

## Related Presets

- Pair with the `CloudflareWaitingRoomEvent` **01-launch-window** preset for a scheduled override of these rates
