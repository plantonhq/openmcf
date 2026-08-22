# Tag-Targeted Fleet CPU Alert

This preset watches CPU across every droplet carrying the `web` tag -- membership tracks the tag automatically, so replacements and scale-ups are covered the moment they exist, with no manifest change.

## When to Use

- Droplet fleets behind load balancers or autoscale pools, where members come and go
- The standard "something is running hot" signal routed to a team inbox

## Key Configuration Choices

- **Tags over id lists** -- id-targeted policies watch exactly the droplets listed; tag targeting follows the fleet.
- **`window: 10m`** -- long enough to ignore boot spikes; tighten to `5m` for latency-sensitive services.
- **Email delivery** -- add `slack` rows (channel + webhook URL, kept secret) for channel paging.

## What You Get

One policy covering the whole tagged fleet, its `alert_id` exported for reference.
