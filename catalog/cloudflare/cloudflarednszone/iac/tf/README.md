# Terraform Module: Cloudflare DNS Zone

Provisions a `cloudflare_zone` together with its zone-singleton satellites —
`cloudflare_zone_dns_settings`, `cloudflare_zone_dnssec`, `cloudflare_zone_hold`,
`cloudflare_zone_subscription` — and an inline set of `cloudflare_dns_record`
resources whose lifecycle tracks the zone. Records with independent lifecycles
are better modeled as standalone `cloudflarednsrecord` resources; the inline
surface is identical in depth.

## Resources

- `cloudflare_zone` — the zone itself (name, owning account, type, paused,
  vanity nameservers).
- `cloudflare_dns_record` (one per `spec.records[]` entry) — inline records,
  including typed structured data for the 13 non-content record types.
- `cloudflare_zone_dns_settings` (`count`-gated on `spec.dns_settings`) —
  CNAME flattening, zone mode, SOA, nameserver set, internal-DNS fallback.
- `cloudflare_zone_dnssec` (`count`-gated on `spec.dnssec.enabled`) — zone
  signing; the DS material surfaces through outputs for the registrar.
- `cloudflare_zone_hold` (`count`-gated on `spec.hold.enabled`) — blocks the
  zone's hostname (optionally subdomains) from being created as a zone in any
  other account.
- `cloudflare_zone_subscription` (`count`-gated on `spec.subscription`) — the
  zone's rate plan. A paid plan bills real money and needs Billing Write token
  scope.

## Inputs

- `metadata` — name/labels.
- `spec` — see [variables.tf](./variables.tf). Required: `zone_name`,
  `account_id`.
- `spec.records[]` — a record is either a simple record (`content`) or a
  structured record whose typed case (srv/caa/cert/…) arrives as a top-level
  attribute on the record entry (the tfvars converter emits the active oneof
  case by its own name; a `data` wrapper key never appears). The flatten in
  [records.tf](./records.tf) rebuilds the provider's single union `data`
  object from whichever case is set. `ttl` 0 maps to 1 (automatic — the
  provider requires ttl ≥ 1), and top-level `priority` is sent for MX only
  (SRV/URI/HTTPS/SVCB carry theirs inside their structured data).

## Outputs

| Output | Description |
|---|---|
| `zone_id` | The zone ID (also the import identity of every zone-singleton satellite) |
| `nameservers` | The Cloudflare-assigned nameservers |
| `status` | Zone status (`pending` until the registrar delegates to Cloudflare) |
| `dnssec_status`, `dnssec_ds`, `dnssec_digest`, `dnssec_digest_type`, `dnssec_digest_algorithm`, `dnssec_algorithm`, `dnssec_key_tag`, `dnssec_public_key`, `dnssec_flags` | DS/DNSKEY material to enter at the registrar (empty when DNSSEC is off) |

## Requirements

- The provider reads `CLOUDFLARE_API_TOKEN` from the environment. The token
  needs **Zone → Zone → Edit** and **Zone → DNS → Edit**; add
  **Zone → Zone Settings → Edit** for the hold, and **Billing → Edit** (Billing
  Write) when `spec.subscription` is set.
- `vanity_name_servers` needs a Business/Enterprise plan; `foundation_dns` is a
  paid add-on; `multi_provider` and `secondary_overrides` are plan-gated.
