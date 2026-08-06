# ScalewayDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayDnsZoneSpec defines the specification for a Scaleway DNS zone
with optional inline DNS records.

Scaleway DNS zones are managed through the Scaleway Domains and DNS
product. A zone represents a delegated portion of the DNS namespace
for a domain you own or have been delegated control of. Zones are
created by specifying a parent domain and an optional subdomain prefix.

This is a **composite resource** wrapping:
  - `scaleway_domain_zone` (1x) -- the zone itself
  - `scaleway_domain_record` (0..Nx) -- one per inline record entry

**Zone naming:**
  - Root zone: `domain = "example.com"`, `subdomain = ""` -> zone name is `example.com`
  - Subdomain zone: `domain = "example.com"`, `subdomain = "staging"` -> zone name is `staging.example.com`

**Two ways to manage records:**
  - **Inline records** (this spec's `records` field): Convenience for
    static records known at zone creation time (MX, SPF/TXT, CAA).
    Records are managed as part of the zone's lifecycle.
  - **Standalone ScalewayDnsRecord** (R16): DAG-friendly for records
    that reference other resources' outputs (A records pointing to a
    Load Balancer IP, CNAMEs to a Kapsule cluster endpoint). These
    create explicit dependency edges in infra charts.

Both patterns coexist. Use inline records for convenience, standalone
records for cross-resource wiring.

**DNS delegation:** After creating a zone, configure the nameservers
(from `status.outputs.name_servers`) at your domain registrar to
delegate DNS resolution to Scaleway.

**Scaleway DNS limitations:**
  - No DNSSEC support
  - No traffic routing policies (geo, weighted, latency, failover)

**Composition pattern:** This is a foundation resource (DAG Layer 0)
with no upstream `StringValueOrRef` dependencies. The primary output
`zone_name` is referenced by downstream ScalewayDnsRecord resources.
Inline record values use `StringValueOrRef` so they can reference
other resources' outputs (e.g., a Load Balancer's IP address).

**Deferred features (not in v1):**
  - DNSSEC configuration
  - Dynamic record types (geo_ip, http_service, view, weighted)
  These can be added in future versions.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domain` | `string` | yes |  |  |
| `spec.subdomain` | `string` |  |  |  |
| `spec.records` | `[]ScalewayDnsZoneRecord` |  |  |  |
| `spec.records[].name` | `string` |  |  |  |
| `spec.records[].type` | `enum` | yes |  |  |
| `spec.records[].data` | `string \| valueFrom` | yes |  |  |
| `spec.records[].ttl` | `uint32` |  | `3600` |  |
| `spec.records[].priority` | `uint32` |  |  |  |

## Field Details

### spec.domain

`string` · required

The registered parent domain name (e.g., "example.com").

This must be a domain you own or have been delegated control of.
Scaleway does not perform domain registration -- the domain must
already exist at a registrar (Namecheap, Google Domains, etc.).

Cannot be changed after creation (ForceNew in Terraform). Changing
the domain requires recreating the zone resource.

- rule: {"required":true}

### spec.subdomain

`string`

Subdomain prefix for this zone.

Leave empty ("") for the root zone of the domain.
Set to a value like "staging" to create a zone for
"staging.example.com", which enables subdomain delegation
(a separate set of nameservers for a portion of the namespace).

Can be updated after creation without recreating the zone.

Common subdomain zones:
  - "staging" -> staging.example.com
  - "dev" -> dev.example.com
  - "internal" -> internal.example.com

### spec.records

`[]ScalewayDnsZoneRecord`

DNS records to create within this zone.

Each entry creates one `scaleway_domain_record` resource.
Suitable for static records known at zone creation time:
  - MX records for email routing
  - TXT records for SPF, DKIM, DMARC, domain verification
  - CAA records for certificate authority authorization
  - NS records for subdomain delegation

For records whose values come from other infrastructure resources
(A records pointing to a Load Balancer IP, CNAMEs to a Kapsule
wildcard DNS), prefer the standalone ScalewayDnsRecord kind for
explicit DAG dependency visibility.

If no records are provided, the zone is created empty (common
when all records are managed as standalone ScalewayDnsRecord
resources or by external systems).

### spec.records[].name

`string`

Record name relative to the zone.

Use empty string or omit for the zone apex (root record).
Scaleway normalizes "@" and empty string to the zone root.

Examples:
  "" or "@" -> zone apex (e.g., example.com)
  "www" -> www.example.com
  "api" -> api.example.com
  "_dmarc" -> _dmarc.example.com
  "mail" -> mail.example.com

### spec.records[].type

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

- rule: {"required":true}

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

### spec.records[].data

`string | valueFrom` · required

Record data/value.

Can be a literal string or a reference to another resource's
output via StringValueOrRef. This enables inline records to
reference dynamic values from other infrastructure:
  - A Load Balancer's IP address
  - A Kapsule cluster's wildcard DNS endpoint
  - An Instance's public IP

No default_kind is specified because record values can reference
many different resource types.

Examples:
  Literal A record: { value: "192.0.2.1" }
  Reference to LB IP: { value_from: { kind: ScalewayLoadBalancer, ... } }

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.records[].ttl

`uint32`

Time to live in seconds.

Determines how long DNS resolvers cache this record before
querying Scaleway's nameservers again.

Common TTL values:
  300 (5 min) -- during migrations or cutover events
  3600 (1 hour) -- default for most records
  86400 (24 hours) -- static records that rarely change

Valid range: 60-2592000 seconds (1 minute to 30 days).
Defaults to 3600 (1 hour) if not specified.

- default: `3600`

### spec.records[].priority

`uint32`

Priority for MX and SRV records.

Lower values indicate higher priority. For MX records, this
determines mail delivery preference (e.g., 1 = primary, 10 = backup).

Ignored for record types other than MX and SRV.
Defaults to 0 if not specified.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_name` | `string` | The computed zone name. Format: "{subdomain}.{domain}" for subdomain zones, or just "{domain}" for root zones. Examples: Root zone: "example.com" Subdomain zone: "staging.example.com" This is the primary output for downstream cross-resource references. ScalewayDnsRecord resources reference this value to identify which zone their records belong to. |
| `status.outputs.name_servers` | `[]string` | Nameservers assigned by Scaleway for this zone. These are the NS records that must be configured at the domain registrar for DNS delegation. Until delegation is complete, DNS queries for this zone will not resolve through Scaleway. Typically contains 2-4 Scaleway nameserver addresses. |
| `status.outputs.name_servers_default` | `[]string` | Scaleway's default nameservers for this zone. These are the standard NS records assigned by Scaleway's DNS infrastructure. Usually identical to `name_servers` unless custom nameservers have been configured. |
| `status.outputs.name_servers_master` | `[]string` | Master nameservers for this zone. For standard zones, this is typically the same as the default nameservers. For zones with secondary/slave configurations, this identifies the primary nameserver(s). |
| `status.outputs.status` | `string` | Zone status. Indicates the current state of the DNS zone in Scaleway's infrastructure. Common values: "active" -- zone is operational and serving DNS queries "pending" -- zone is being provisioned "error" -- zone has configuration issues Exported for observability and health monitoring. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| ScalewayDnsRecord | `spec.zoneName` | `status.outputs.zone_name` |

## See Also

- [Overview](./README.md)
