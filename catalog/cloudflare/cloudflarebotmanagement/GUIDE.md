# CloudflareBotManagement guide

The judgment this guide protects you from: Bot Management has no delete at Cloudflare, so this resource's destroy restores nothing -- and every field you set becomes a value you own until you explicitly set it to something else. Read the destroy section before you deploy, not after.

## Destroy is a no-op: revert before you retire

Cloudflare's Bot Management API has no delete operation. Destroy drops the resource from state and walks away; the live configuration keeps serving. Cloudflare's own tooling warns the same.

The consequence: if you set `fight_mode: true` and later delete the resource, Fight Mode stays on indefinitely, with nothing in your infrastructure declaring it. The discipline is revert-before-retire: before removing a field or destroying the resource, set every managed field to the value the zone should keep, apply, and only then remove it.

## Unset means unmanaged, by design

The module never sends defaults. Every settings field is presence-tracked, and only fields you set produce an API write. If the module sent "sensible defaults" for the other fourteen fields when you managed one, it would take ownership of values it can never give back, and the first apply would silently overwrite whatever the dashboard or a plan default had put there.

Validation enforces the same philosophy at the floor: a spec with only `zone_id` is rejected, because a resource that manages nothing would deploy nothing.

## Plan families are Cloudflare's wall

Fields group into four families. Setting a field the zone's plan does not include fails at the API, not here -- the plan walls move, and they are Cloudflare's.

- **Free:** `fight_mode`. Mutually exclusive with Super Bot Fight Mode -- zones on SBFM plans manage the `sbfm_*` fields instead.
- **Pro/Business (Super Bot Fight Mode):** `sbfm_definitely_automated`, `sbfm_likely_automated` (Business+), `sbfm_verified_bots`, `sbfm_static_resource_protection`, `optimize_wordpress`.
- **Enterprise Bot Management:** `auto_update_model`, `suppress_session_score`, `enable_js`, `bm_cookie_enabled`.
- **AI and crawler controls:** `ai_bots_protection`, `crawler_protection`, `content_bots_protection`, `cf_robots_variant`, `is_robots_txt_managed`.

Blocking `sbfm_verified_bots` blocks search indexing -- `allow` is the safe default posture. `sbfm_static_resource_protection` can catch legitimate hotlinked-asset traffic; Cloudflare's own docs caution the same.

## Non-entitled zones omit fields: refresh drift

On a zone whose plan does not include a field, the API omits that field from responses entirely. The provider then refreshes a manifest that set the field as a perpetual plan diff. Manage only what the plan includes. The provider's own issue tracker records this class; it is not a Planton bug and it will not clean up on retry.

## Pairs well with

- [CloudflareIpAccessRule](../cloudflareipaccessrule/README.md) -- a static IP/ASN/country decision when scored bot traffic is the wrong tool.
- [CloudflareRuleset](../cloudflareruleset/README.md) -- skip rules for verified-bot exceptions under `content_bots_protection`.
- [CloudflareDnsZone](../cloudflarednszone/README.md) -- wire `zone_id` via `value_from`; the zone plan is what gates SBFM and Enterprise fields.
