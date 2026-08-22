# CloudflareZoneSettings guide

The judgment this guide protects you from: zone settings have no delete at Cloudflare, so this resource's destroy restores nothing -- and every field you set becomes a value you own until you explicitly set it to something else. Read the destroy section before you deploy, not after.

## Destroy is a no-op: revert before you retire

Cloudflare's zone settings API has no delete operation. The provider's destroy for a zone setting drops it from state and walks away; the live value keeps serving. The same is true for `waitingRoomCrawlerBypass`. Of the whole spec, only two companions actually reset on destroy: managed transforms (destroy disables them) and URL normalization (destroy restores Cloudflare defaults).

The consequence: if you set `securityLevel: under_attack` during an incident and later delete the resource, the zone stays in under-attack mode indefinitely, with nothing in your infrastructure declaring it. The discipline is revert-before-retire: before removing a field or destroying the resource, set every managed field to the value the zone should keep, apply, and only then remove it. Treat the manifest's last applied state as what the zone will carry forever unless something else changes it.

## Unset means unmanaged, by design

The module never sends defaults. Every settings field is presence-tracked (`optional`, message, or repeated), and only fields you set produce an API write. This is not laziness -- it is the only safe design given the no-delete contract above. If the module sent "sensible defaults" for the other 62 settings when you managed one, it would take ownership of values it can never give back, and the first apply would silently overwrite whatever the dashboard, another team, or a Cloudflare plan default had put there.

Corollary: this resource composes with dashboard-managed settings. Managing `alwaysUseHttps` here does not stop anyone from flipping `rocketLoader` in the dashboard. But for a field you DO manage, every apply reasserts your value -- the manifest wins any tug-of-war over managed fields, one apply at a time.

Validation enforces the same philosophy at the floor: a spec with only `zoneId` is rejected, because a resource that manages nothing would deploy nothing.

## Plan-gated settings fail at apply, not silently

Cloudflare reports a setting the zone's plan cannot edit as editable=false, and the module surfaces that as an apply error. Nothing is skipped and nothing is billed. `advancedDdos`, `orangeToOrange`, `prefetchPreload`, `responseBuffering`, `sortQueryStringForCache`, `trueClientIpHeader`, and `proxyReadTimeout` are Enterprise territory; `polish`, `mirage`, and `imageResizing` need Pro or above. If an apply fails on one of these, the fix is to remove the field or upgrade the plan (a `CloudflareDnsZone` `subscription` decision, and real money) -- not to retry.

## The boundary with CloudflareDnsZone

Two kinds touch "zone settings" and the split is deliberate: `CloudflareDnsZone`'s `dnsSettings` owns how the zone resolves (SOA tuning, NS TTL, DNS-level CNAME flattening behavior, DNSSEC); this kind owns how Cloudflare serves HTTP for the zone (terminate, cache, transform, protect). If you are reaching for a resolver-facing knob here, it lives on the zone resource. The `cnameFlattening` field in this spec is the serving-side setting id of the same name -- Cloudflare exposes it through the settings API, so it is managed here.

On the diagram this reads correctly too: the zone node creates and resolves; a `CloudflareZoneSettings` node hanging off its `zone_id` declares the serving posture. One settings resource per zone -- it is a singleton, and the only output is the `zone_id` passthrough.

## Quirks the module absorbs so you don't have to

- **`sslRecommender` uses the value form.** The provider documents an enabled-attribute form for this setting, but that form fails the provider's own validation at v5.23.0. The module sends the plain on/off value form on writes, which the API accepts. You just set a boolean.
- **`zeroRtt` is setting id `0rtt`.** Proto identifiers cannot start with a digit, so this is the one field whose name differs from its Cloudflare setting id. Everything else matches the API's vocabulary exactly.
- **`aegis` and `automaticPlatformOptimization` value shapes come from Cloudflare's Go SDK.** The Terraform provider's own schema files do not describe these settings' object values; the module builds them to the SDK's shapes. APO in particular requires every member of its object on every write -- there are no server-side defaults for omitted members, which is why the spec makes the whole object explicit.
- **`longLivedGrpc` is real but undocumented.** It does not appear in the API's documented setting table at provider v5.23.0, but the settings endpoint accepts it and the provider's regression tests exercise it. Modeled as a normal toggle; if Cloudflare ever drops it, expect the failure at apply, not at validation.
- **The provider is pinned to v5.23.** The value-form quirk and the undocumented-setting evidence above are verified against that version. Bumping the provider is a deliberate act: re-check `sslRecommender`'s write form first, since it is the known point where documented and actual behavior diverge.

## Pairs well with

- [CloudflareDnsZone](../cloudflarednszone/README.md) -- the zone itself; wire `zoneId` via `valueFrom` and the reference defaults to `status.outputs.zone_id`.
- [CloudflareCacheSettings](../cloudflarecachesettings/README.md) -- cache rules and cache reserve; this kind's `cacheLevel`/`browserCacheTtl`/`edgeCacheTtl` set the zone-wide baseline those rules override.
- [CloudflareZoneTlsSettings](../cloudflarezonetlssettings/README.md) -- advanced TLS surface beyond this spec's `ssl`/`minTlsVersion`/`ciphers` toggles.
- [CloudflareRuleset](../cloudflareruleset/README.md) -- when a setting should differ per path or per request (`http_config_settings` phase), instead of zone-wide here.
