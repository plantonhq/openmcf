# Bot Fight Mode

Turns on Bot Fight Mode -- the one Bot Management field every plan carries -- and leaves every other knob unmanaged. Destroy will not turn it off; apply `fight_mode: false` first if you want it off.

## When to Use

- First Bot Management resource on a free zone
- A zone that should challenge known-bot patterns without Super Bot Fight Mode
- The safest starting point: one field, reversible by applying false

## Key Configuration Choices

- **fight_mode: true** -- free-plan Bot Fight Mode. Mutually exclusive with Super Bot Fight Mode at Cloudflare; zones on SBFM plans manage the `sbfm_*` fields instead.
- **enable_js: true** -- required alongside `fight_mode` when the zone's JavaScript detections are off (a fresh zone's default): Cloudflare rejects Fight Mode alone with "cannot enable Fight_Mode while EnableJS is disabled".
- **Unset fields stay unmanaged** -- the module only sends what you set. Do not add SBFM or Enterprise fields on a free zone (the API omits them from responses and refresh reads as drift).
- **No-op destroy** -- deleting this resource leaves Fight Mode on. Revert before you retire.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|-------------|-------------|---------------|
| `zone_id.value` | The Cloudflare zone whose Bot Management configuration is managed | Cloudflare Dashboard -> zone Overview -> API section (right sidebar), or reference a CloudflareDnsZone via `value_from` instead |

## Related Presets

None yet -- Super Bot Fight Mode and Enterprise knobs are plan-gated and belong on a zone that already has the entitlement.
