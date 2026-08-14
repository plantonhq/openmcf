# CloudflareDnsZone guide

Operational judgment for configuring the zone well. The README covers what each field is; this covers how the pieces interact.

## Pending vs. active zones

A newly created zone is `pending` until your registrar delegates the domain to the nameservers in `status.outputs.nameservers`. Nearly everything works on a pending zone -- records, rulesets, DNS settings, DNSSEC configuration, most zone-scoped resources -- but nothing resolves publicly and proxied records serve no traffic. Check `status.outputs.status` before treating the zone as live, and expect 1-24h of registrar propagation after updating nameservers.

## Inline records vs. standalone CloudflareDnsRecord

Both surfaces carry identical depth (all 21 types, typed structured data, tags, settings), so the choice is purely lifecycle:

- **Inline** when the records live and die with the zone: bootstrap A/AAAA records, MX + SPF/DKIM TXTs, CAA issuance policy, verification TXTs. One manifest, one apply.
- **Standalone** when a record changes on its own cadence (a CI-managed deploy record, per-service SRVs owned by different teams) or is created by another chart wiring `zone_id` via ValueFromRef.

Renaming an inline record's `name` or `type` replaces the record (the resource key includes both) -- brief NXDOMAIN between destroy and create. For records where that matters, standalone resources give finer control.

## Simple content vs. typed data blocks

A record's value comes from exactly one place, and it must match `type`:

- `content` for A, AAAA, CNAME, MX, NS, PTR, TXT, OPENPGPKEY -- the presentation string.
- A typed block named after the record type for SRV, CAA, CERT, DNSKEY, DS, HTTPS, LOC, NAPTR, SMIMEA, SSHFP, SVCB, TLSA, URI (e.g. `srv: {priority, weight, port, target}`).

Top-level `priority` is only for MX. SRV, URI, HTTPS, and SVCB carry their own priority inside their typed block -- setting the top-level field for those types does nothing.

## TTL and the proxy

`ttl: 1` (or leaving it unset) means "automatic", which is what proxied records should use -- Cloudflare answers proxied queries with its own edge IPs and manages the TTL itself. Explicit TTLs (30-86400s) matter only for gray-cloud records. `settings.flatten_cname` is likewise a gray-cloud lever: proxied CNAMEs are always flattened.

## DNSSEC is a two-step handshake

`dnssec.enabled: true` makes Cloudflare sign the zone, but the chain of trust completes only when you enter the DS material (`dnssec_ds` / `dnssec_digest` / `dnssec_key_tag` outputs) at your registrar. Do it in that order; entering DS records for an unsigned zone breaks resolution. `multi_signer` and `presigned` are for multi-provider and Cloudflare-as-secondary setups -- leave both off for the normal single-provider case.

## The hold is a takeover guard, not a lock

`hold.enabled: true` stops the zone's hostname from being ADDED as a zone in any other Cloudflare account -- protection while a domain migrates between accounts or when subdomains delegate to other teams (`include_subdomains: true` extends the guard to every subdomain, including SSL-for-SaaS custom hostnames). It does not freeze the zone's own configuration. A future-dated `hold_after` temporarily releases the hold until that instant -- a planned migration window that re-arms itself.

## Subscription is a billing action

`subscription.rate_plan` is real money on anything above `free`, effective at apply, and the deploying token needs Billing Write scope (a plain Zone-Edit token gets a 403 on just this satellite). Plan choice also gates other levers in this spec: `vanity_name_servers` needs Business+, `foundation_dns` is a paid add-on, and `multi_provider`/`secondary_overrides` are plan-gated. If your organization manages plans centrally, leave `subscription` unset and let the account default stand.

## Zone types pair with specific settings

- `full` (default): Cloudflare hosts DNS entirely. The normal case.
- `partial`: CNAME setup for partner-hosted zones; most zone-wide DNS settings don't apply.
- `secondary`: Cloudflare mirrors an external primary -- pair with `dns_settings.secondary_overrides` to allow proxied overrides, and note the zone-transfer configuration itself is Enterprise territory.
- `internal`: internal-only resolution -- pair with `dns_settings.internal_dns.reference_zone_id` (a ValueFromRef to another CloudflareDnsZone) for fallback resolution.
