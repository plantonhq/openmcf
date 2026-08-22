# Pulumi Module: Cloudflare DNS Zone

Provisions a `cloudflare.Zone` together with its zone-singleton satellites —
`cloudflare.ZoneDnsSettings`, `cloudflare.ZoneDnssec`, `cloudflare.ZoneHold`,
`cloudflare.ZoneSubscription` — and an inline set of `cloudflare.DnsRecord`
resources whose lifecycle tracks the zone. Records with independent lifecycles
are better modeled as standalone `cloudflarednsrecord` resources; the inline
surface is identical in depth.

## Layout

```
iac/pulumi/
├── main.go            # entrypoint (loads stack-input, calls module.Resources)
├── Pulumi.yaml
└── module/
    ├── main.go            # Resources(): provider setup + zone()
    ├── locals.go          # stack-input references
    ├── zone.go            # cloudflare.Zone + satellite orchestration + outputs
    ├── records.go         # inline records + the typed-data builder
    ├── dns_settings.go    # cloudflare.ZoneDnsSettings
    ├── dnssec.go          # cloudflare.ZoneDnssec (status "active" when enabled)
    ├── hold.go            # cloudflare.ZoneHold (created only when enabled)
    ├── subscription.go    # cloudflare.ZoneSubscription (created only when set)
    └── outputs.go         # output constant names
```

## Inputs

A `CloudflareDnsZoneStackInput` (target + provider config). Required spec
fields: `zoneName`, `accountId`.

`spec.records[]` entries are either simple records (`content`) or structured
records whose typed oneof case (srv/caa/cert/…) is translated by
`buildRecordData` in [module/records.go](./module/records.go) into the
provider's single union `data` object. `ttl` 0 maps to 1 (automatic — the
provider requires ttl ≥ 1), and top-level `priority` is sent for MX only
(SRV/URI/HTTPS/SVCB carry theirs inside their structured data).

## Outputs

- `zone_id` — the zone ID (also the import identity of every zone-singleton
  satellite).
- `nameservers` — the Cloudflare-assigned nameservers.
- `status` — zone status (`pending` until the registrar delegates).
- `dnssec_*` — DS/DNSKEY material to enter at the registrar (empty strings when
  DNSSEC is off, so the output contract is always satisfied).

## Requirements

- The Cloudflare provider is configured from the stack-input provider config /
  `CLOUDFLARE_API_TOKEN`. The token needs **Zone → Zone → Edit** and
  **Zone → DNS → Edit**; add **Zone → Zone Settings → Edit** for the hold, and
  **Billing → Edit** (Billing Write) when `spec.subscription` is set.
- `vanity_name_servers` needs a Business/Enterprise plan; `foundation_dns` is a
  paid add-on; `multi_provider` and `secondary_overrides` are plan-gated.
