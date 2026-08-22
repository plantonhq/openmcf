# Zone with automatic installation

The zero-code shape: measure a Cloudflare zone and let the edge inject the beacon into every proxied page -- no snippet to paste, no deploy to coordinate. The zone must actually be proxied (orange-clouded); on a DNS-only zone Cloudflare has no response to inject into and `auto_install` silently does nothing. Point `zone_tag` at a `CloudflareDnsZone` reference in real manifests so the site follows the zone it belongs to.
