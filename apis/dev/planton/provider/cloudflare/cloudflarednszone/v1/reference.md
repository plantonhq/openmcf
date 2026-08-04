# CloudflareDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1`

**CloudflareDnsZoneSpec** defines a Cloudflare DNS zone: the domain, its
account, common zone-level options, optional zone-wide DNS settings and DNSSEC,
and an optional set of DNS records managed alongside the zone. Records may also
be managed independently as first-class CloudflareDnsRecord resources; the
embedded list is a convenience for records whose lifecycle tracks the zone.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneName` | `string` | yes |  |  |
| `spec.accountId` | `string` | yes |  |  |
| `spec.paused` | `bool` |  |  |  |
| `spec.records` | `[]CloudflareDnsZoneRecord` |  |  |  |
| `spec.records[].name` | `string` | yes |  |  |
| `spec.records[].type` | `enum` | yes |  |  |
| `spec.records[].content` | `string` | yes |  |  |
| `spec.records[].proxied` | `bool` |  |  |  |
| `spec.records[].ttl` | `int32` |  |  |  |
| `spec.records[].priority` | `int32` |  |  |  |
| `spec.records[].comment` | `string` |  |  |  |
| `spec.type` | `enum` |  |  |  |
| `spec.vanityNameServers` | `[]string` |  |  |  |
| `spec.dnsSettings` | `CloudflareDnsZoneDnsSettings` |  |  |  |
| `spec.dnsSettings.flattenAllCnames` | `bool` |  |  |  |
| `spec.dnsSettings.foundationDns` | `bool` |  |  |  |
| `spec.dnsSettings.multiProvider` | `bool` |  |  |  |
| `spec.dnsSettings.secondaryOverrides` | `bool` |  |  |  |
| `spec.dnsSettings.nsTtl` | `uint32` |  |  |  |
| `spec.dnsSettings.zoneMode` | `enum` |  |  |  |
| `spec.dnsSettings.soa` | `CloudflareDnsZoneSoa` |  |  |  |
| `spec.dnsSettings.soa.expire` | `uint32` |  |  |  |
| `spec.dnsSettings.soa.minTtl` | `uint32` |  |  |  |
| `spec.dnsSettings.soa.mname` | `string` |  |  |  |
| `spec.dnsSettings.soa.refresh` | `uint32` |  |  |  |
| `spec.dnsSettings.soa.retry` | `uint32` |  |  |  |
| `spec.dnsSettings.soa.rname` | `string` |  |  |  |
| `spec.dnsSettings.soa.ttl` | `uint32` |  |  |  |
| `spec.dnsSettings.nameservers` | `CloudflareDnsZoneNameservers` |  |  |  |
| `spec.dnsSettings.nameservers.nsSet` | `uint32` |  |  |  |
| `spec.dnsSettings.nameservers.type` | `string` |  |  |  |
| `spec.dnsSettings.internalDns` | `CloudflareDnsZoneInternalDns` |  |  |  |
| `spec.dnsSettings.internalDns.referenceZoneId` | `string \| valueFrom` |  |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.dnssec` | `CloudflareDnsZoneDnssec` |  |  |  |
| `spec.dnssec.enabled` | `bool` |  |  |  |
| `spec.dnssec.multiSigner` | `bool` |  |  |  |
| `spec.dnssec.presigned` | `bool` |  |  |  |
| `spec.dnssec.useNsec3` | `bool` |  |  |  |

## Field Details

### spec.zoneName

`string` · required

The fully qualified domain name of the DNS zone (e.g., "example.com").

- rule: zone_name must be a valid fully qualified domain name
- rule: {"required":true}

### spec.accountId

`string` · required

The Cloudflare account identifier under which to create the zone.

- rule: {"required":true}

### spec.paused

`bool`

Whether the zone is created paused. A paused zone uses Cloudflare DNS only
and receives no security or performance (proxy/CDN/WAF) benefits.

### spec.records

`[]CloudflareDnsZoneRecord`

DNS records managed alongside this zone. For records with independent
lifecycles (or that need the full record feature set), prefer standalone
CloudflareDnsRecord resources instead.

- rule: proxied can only be true for A, AAAA, or CNAME records
- rule: priority is required for MX records

### spec.records[].name

`string` · required

The record name (or "@" for the zone apex).

- rule: {"required":true}

### spec.records[].type

`enum` · required

The DNS record type.

- rule: type must be specified (cannot be record_type_unspecified)
- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `record_type_unspecified` -- Unspecified record type (invalid).
- `A` -- IPv4 address record.
- `AAAA` -- IPv6 address record.
- `CNAME` -- Canonical name (alias) record.
- `MX` -- Mail exchange record.
- `TXT` -- Text record (SPF, DKIM, verification, etc.).
- `SRV` -- Service locator record.
- `NS` -- Nameserver record.
- `CAA` -- Certification Authority Authorization record.

### spec.records[].content

`string` · required

The record's value/target in presentation format. For A: an IPv4 address.
For AAAA: an IPv6 address. For CNAME/MX/NS: a hostname. For TXT: text.

- rule: {"required":true}

### spec.records[].proxied

`bool`

Whether the record is proxied through Cloudflare (orange cloud). Only valid
for A, AAAA, and CNAME records.

### spec.records[].ttl

`int32`

Time to live (TTL) in seconds. Leave 0 or set 1 for automatic; otherwise
30-86400.

- rule: ttl must be 0 or 1 (automatic) or between 30 and 86400 seconds

### spec.records[].priority

`int32`

Priority for MX records (lower is preferred). Required for MX.

- rule: priority must be between 0 and 65535

### spec.records[].comment

`string`

Optional comment/note for the record. Has no effect on DNS responses.

### spec.type

`enum`

The zone's deployment type. Defaults to "full".

Allowed values (use exactly as shown):

- `zone_type_unspecified`
- `full`
- `partial`
- `secondary`
- `internal`

### spec.vanityNameServers

`[]string`

Custom (vanity) name servers for the zone. Only available on Business and
Enterprise plans; leave empty to use Cloudflare's assigned name servers.

### spec.dnsSettings

`CloudflareDnsZoneDnsSettings`

Optional zone-wide DNS settings (CNAME flattening, zone mode, SOA, etc.).
Omit to leave Cloudflare's defaults in place.

### spec.dnsSettings.flattenAllCnames

`bool`

Flatten all CNAME records in the zone (a CNAME at the apex is always flattened).

### spec.dnsSettings.foundationDns

`bool`

Enable Foundation DNS Advanced Nameservers on the zone.

### spec.dnsSettings.multiProvider

`bool`

Enable multi-provider DNS (activates the zone even with non-Cloudflare NS
records and respects apex NS records during outbound transfers).

### spec.dnsSettings.secondaryOverrides

`bool`

Allow a secondary zone to use proxied override records and apex CNAME
flattening.

### spec.dnsSettings.nsTtl

`uint32`

Time to live (TTL), in seconds, for the zone's nameserver (NS) records.
Leave 0 to use the default; otherwise 30-86400.

- rule: ns_ttl must be 0 (default) or between 30 and 86400 seconds

### spec.dnsSettings.zoneMode

`enum`

Whether the zone is a standard zone or a CDN/DNS-only zone.

Allowed values (use exactly as shown):

- `zone_mode_unspecified`
- `standard`
- `cdn_only`
- `dns_only`

### spec.dnsSettings.soa

`CloudflareDnsZoneSoa`

Components of the zone's SOA record.

### spec.dnsSettings.soa.expire

`uint32`

Seconds a secondary will keep serving the zone after losing contact with the
primary (86400-2419200).

- rule: expire must be 0 (default) or between 86400 and 2419200 seconds

### spec.dnsSettings.soa.minTtl

`uint32`

TTL for negative caching of records within the zone (60-86400).

- rule: min_ttl must be 0 (default) or between 60 and 86400 seconds

### spec.dnsSettings.soa.mname

`string`

The primary nameserver (MNAME). Leave empty for a Cloudflare-assigned value.

### spec.dnsSettings.soa.refresh

`uint32`

Seconds after which a secondary re-checks the SOA for updates (600-86400).

- rule: refresh must be 0 (default) or between 600 and 86400 seconds

### spec.dnsSettings.soa.retry

`uint32`

Seconds a secondary waits before retrying after an unresponsive primary
(600-86400).

- rule: retry must be 0 (default) or between 600 and 86400 seconds

### spec.dnsSettings.soa.rname

`string`

The zone administrator's email (RNAME), first label being the local part.

### spec.dnsSettings.soa.ttl

`uint32`

TTL of the SOA record itself (300-86400).

- rule: ttl must be 0 (default) or between 300 and 86400 seconds

### spec.dnsSettings.nameservers

`CloudflareDnsZoneNameservers`

Settings determining the nameservers through which the zone is available.

### spec.dnsSettings.nameservers.nsSet

`uint32`

Configured nameserver set to use for this zone (1-5).

- rule: ns_set must be 0 (default) or between 1 and 5

### spec.dnsSettings.nameservers.type

`string`

Nameserver type: one of "cloudflare.standard", "custom.account",
"custom.tenant", "custom.zone".

- rule: type must be one of "cloudflare.standard", "custom.account", "custom.tenant", "custom.zone"

### spec.dnsSettings.internalDns

`CloudflareDnsZoneInternalDns`

Settings for an internal zone (only relevant when type is "internal").

### spec.dnsSettings.internalDns.referenceZoneId

`string | valueFrom`

The zone to fall back to for resolution. Can be a literal zone ID or a
reference to another CloudflareDnsZone.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.dnssec

`CloudflareDnsZoneDnssec`

Optional DNSSEC configuration. Enable to have Cloudflare sign the zone; the
DS record material to hand to your registrar is published as stack outputs.

### spec.dnssec.enabled

`bool`

Whether DNSSEC is active for the zone.

### spec.dnssec.multiSigner

`bool`

Enable multi-signer DNSSEC, allowing multiple providers to serve a
DNSSEC-signed zone simultaneously (required to add external DNSKEY records).

### spec.dnssec.presigned

`bool`

Allow transferring in an externally-signed (presigned) DNSSEC zone without
Cloudflare signing records on the fly (for Cloudflare-as-secondary setups).

### spec.dnssec.useNsec3

`bool`

Use NSEC3 rather than NSEC for authenticated denial of existence.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | The Cloudflare Zone ID of the created DNS zone. |
| `status.outputs.nameservers` | `[]string` | The list of nameserver addresses assigned to this DNS zone. |
| `status.outputs.status` | `string` | The zone status on Cloudflare (e.g., "initializing", "pending", "active", "moved"). |
| `status.outputs.dnssec_status` | `string` | DNSSEC status (e.g., "active", "disabled"). Empty when DNSSEC is not enabled. |
| `status.outputs.dnssec_ds` | `string` | The full DS record to enter at your registrar. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_digest` | `string` | The DS record digest. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_digest_type` | `string` | The DS digest type code (e.g., "2" for SHA-256). Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_digest_algorithm` | `string` | The DS digest algorithm name. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_algorithm` | `string` | The DNSKEY algorithm code. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_key_tag` | `string` | The DNSKEY key tag. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_public_key` | `string` | The DNSKEY public key. Empty unless DNSSEC is enabled. |
| `status.outputs.dnssec_flags` | `string` | The DNSKEY flags. Empty unless DNSSEC is enabled. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dnsSettings.internalDns.referenceZoneId` | CloudflareDnsZone | `status.outputs.zone_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| CloudflareCertificatePack | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareCustomHostname | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareCustomHostnameFallbackOrigin | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareDnsRecord | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareDnsZone | `spec.dnsSettings.internalDns.referenceZoneId` | `status.outputs.zone_id` |
| CloudflareEmailRoutingRule | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareEmailRoutingZone | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareLoadBalancer | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareR2Bucket | `spec.customDomains[].zoneId` | `status.outputs.zone_id` |
| CloudflareRuleset | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareWorker | `spec.customDomains[].zoneId` | `status.outputs.zone_id` |
| CloudflareWorker | `spec.routes[].zoneId` | `status.outputs.zone_id` |
| CloudflareZeroTrustAccessApplication | `spec.zoneId` | `status.outputs.zone_id` |
| CloudflareZeroTrustAccessGroup | `spec.zoneId` | `status.outputs.zone_id` |
| KubernetesExternalDns | `spec.cloudflare.zoneIdFilters` | `status.outputs.zone_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
