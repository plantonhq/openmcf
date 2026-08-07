# ScalewayDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayDnsRecordSpec defines the specification for a standalone
Scaleway DNS record.

This is the DAG-friendly record management kind. Each ScalewayDnsRecord
creates exactly one `scaleway_domain_record` Terraform resource with
an explicit dependency on the target DNS zone via `StringValueOrRef`.

**When to use standalone vs inline records:**
  - **Standalone ScalewayDnsRecord** (this kind): For records whose
    values come from other infrastructure resources (A records pointing
    to a Load Balancer IP, CNAMEs to a Kapsule cluster endpoint).
    Creates visible dependency edges in infra chart DAGs.
  - **Inline ScalewayDnsZone records**: For static records known at
    zone creation time (MX, SPF/TXT, CAA, domain verification).

**Scaleway's simpler record model:**
Unlike some other DNS providers, Scaleway does not have separate
`weight`, `port`, `flags`, or `tag` fields. These are embedded in
the `data` field using standard DNS record format:
  - SRV data: "weight port target" (e.g., "10 5060 sipserver.example.com.")
  - CAA data: 'flags tag value' (e.g., '0 issue "letsencrypt.org"')
  - TLSA data: "usage selector matching-type cert-data"
Only `priority` exists as a separate field (for MX/SRV ordering).

**Composition pattern:** This is a leaf resource (DAG Layer 3+).
Upstream: `zone_name` references ScalewayDnsZone, `data` can reference
any resource's output. Downstream: generally a terminal node; `fqdn`
output available for downstream consumers.

**Deferred features (not in v1):**
  - Dynamic record types (geo_ip, http_service, view, weighted)
  These can be added in future versions.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneName` | `string \| valueFrom` | yes |  | ScalewayDnsZone (`status.outputs.zone_name`) |
| `spec.name` | `string` |  |  |  |
| `spec.type` | `enum` | yes |  |  |
| `spec.data` | `string \| valueFrom` | yes |  |  |
| `spec.ttl` | `uint32` |  | `3600` |  |
| `spec.priority` | `uint32` |  |  |  |
| `spec.keepEmptyZone` | `bool` |  | `true` |  |

## Field Details

### spec.zoneName

`string | valueFrom` · required

The DNS zone name where this record will be created.

This can be a direct value (e.g., "example.com" or
"staging.example.com") or a reference to a ScalewayDnsZone
resource's output. The zone must already exist in Scaleway.

Examples:
  Direct: { value: "example.com" }
  Reference: { value_from: { kind: ScalewayDnsZone, name: "my-zone" } }

- references: ScalewayDnsZone (`status.outputs.zone_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_name}} -- a bare string does not parse

### spec.name

`string`

Record name relative to the zone.

Use empty string for the zone apex (root record). Scaleway
normalizes "@" and empty string to the zone root.

Cannot be changed after creation (ForceNew in Terraform). Changing
the name requires recreating the record resource.

Examples:
  "" -> zone apex (e.g., example.com)
  "www" -> www.example.com
  "api" -> api.example.com
  "_dmarc" -> _dmarc.example.com
  "_25._tcp.mail" -> for TLSA records

### spec.type

`enum` · required

DNS record type.

All Scaleway-supported record types and their data formats:
  A: IPv4 address (e.g., "192.0.2.1")
  AAAA: IPv6 address (e.g., "2001:db8::1")
  ALIAS: target hostname (e.g., "www.example.com.")
  CAA: flags tag value (e.g., '0 issue "letsencrypt.org"')
  CNAME: target hostname with trailing dot (e.g., "target.example.com.")
  DNAME: delegation target (e.g., "other.example.com.")
  MX: mail server with trailing dot (e.g., "mail.example.com.")
  NS: nameserver with trailing dot (e.g., "ns1.example.com.")
  PTR: pointer target (e.g., "host.example.com.")
  SOA: SOA parameters
  SRV: "weight port target" (e.g., "10 5060 sipserver.example.com.")
  TXT: text data (e.g., "v=spf1 include:_spf.google.com ~all")
  TLSA: "usage selector matching-type cert-data" (DANE)

Cannot be changed after creation (ForceNew in Terraform). Changing
the type requires recreating the record resource.

- rule: type must be specified (cannot be record_type_unspecified)
- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `record_type_unspecified` -- Unspecified record type (invalid).
- `A` -- IPv4 address record.
- `AAAA` -- IPv6 address record.
- `ALIAS` -- Auto-resolved alias record (Scaleway-native, like CNAME at zone apex).
- `CAA` -- Certificate Authority Authorization record.
- `CNAME` -- Canonical name (alias) record.
- `DNAME` -- Delegation name record. Redirects a subtree of the DNS name space.
- `MX` -- Mail exchange record.
- `NS` -- Nameserver record.
- `PTR` -- Pointer record (reverse DNS).
- `SOA` -- Start of authority record.
- `SRV` -- Service locator record.
- `TXT` -- Text record (SPF, DKIM, DMARC, domain verification, etc.).
- `TLSA` -- Transport Layer Security Association record (DANE).

### spec.data

`string | valueFrom` · required

Record data/value.

This can be a direct value or a reference to another resource's
output via StringValueOrRef. This is the primary mechanism for
cross-resource wiring in infra charts.

No default_kind is specified because record values can reference
many different resource types (Load Balancer IPs, Instance public
IPs, Kapsule cluster endpoints, etc.).

Format depends on record type -- see the `type` field documentation
for data format examples per type.

Examples:
  Literal A record: { value: "192.0.2.1" }
  Reference to LB IP: { value_from: { kind: ScalewayLoadBalancer, ... } }
  Reference to Kapsule DNS: { value_from: { kind: ScalewayKapsuleCluster, ... } }

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.ttl

`uint32`

Time to live in seconds.

Determines how long DNS resolvers cache this record before
querying Scaleway's nameservers again.

Common TTL values:
  300 (5 min) -- during migrations or cutover events
  3600 (1 hour) -- default for most records
  86400 (24 hours) -- static records that rarely change

Defaults to 3600 (1 hour) if not specified.

- default: `3600`

### spec.priority

`uint32`

Priority for MX and SRV records.

Lower values indicate higher priority. For MX records, this
determines mail delivery preference (e.g., 1 = primary, 10 = backup).

Ignored for record types other than MX and SRV.
Defaults to 0 if not specified.

### spec.keepEmptyZone

`bool`

Keep the DNS zone alive when this record is the last one being
destroyed.

When false (Scaleway default), destroying the last record in a zone
also deletes the zone itself. When true, the zone persists even
when empty.

Recommended to set to true when zones are managed by separate
ScalewayDnsZone resources (the common pattern), to prevent
accidental zone deletion.

- default: `true`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.record_id` | `string` | The unique identifier of the created DNS record in Scaleway. Format: "{dns_zone}/{record_id}" as assigned by the Scaleway API. Used for API operations and Terraform import. |
| `status.outputs.fqdn` | `string` | The fully qualified domain name of the DNS record. Computed by Scaleway from the record name and zone name. Examples: "www.example.com" (for name="www", zone="example.com") "example.com" (for name="", zone="example.com") "api.staging.example.com" (for name="api", zone="staging.example.com") This is the primary downstream output -- use it when other resources need to know the full record name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneName` | ScalewayDnsZone | `status.outputs.zone_name` |

## See Also

- [Overview](../README.md)
