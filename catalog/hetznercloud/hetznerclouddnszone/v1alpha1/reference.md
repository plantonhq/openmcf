# HetznerCloudDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `hetzner-cloud.planton.dev/v1alpha1`

HetznerCloudDnsZoneSpec defines the specification for a Hetzner Cloud DNS
zone with its record sets.

A DNS zone represents a domain (e.g., "example.com") hosted on Hetzner
Cloud's authoritative nameservers. In primary mode, records are managed
directly through this component. In secondary mode, Hetzner Cloud
synchronizes records from an external primary nameserver via zone transfer
(AXFR/IXFR).

Record sets group all DNS records that share the same (name, type) pair.
Each record set is provisioned as a separate hcloud_zone_rrset resource,
which allows in-place value updates without destroying the record set.

Bundled provider resources:
  - hcloud_zone:       The DNS zone itself.
  - hcloud_zone_rrset: One per entry in the record_sets list.

Fields not exposed in this spec (hardcoded or derived in IaC modules):
  - name:   Set to domain_name (the DNS domain, not metadata.name).
  - labels: Derived from metadata (CG01 pattern). Standard labels take
            precedence over user-specified metadata.labels.

## Example

```yaml
apiVersion: hetzner-cloud.planton.dev/v1alpha1
kind: HetznerCloudDnsZone
metadata:
  name: hetznerclouddnszone-demo
  org: demo-org
  env: dev
  labels:
    team: platform
spec:
  domainName: example.com
  mode: primary
  ttl: 3600
  deleteProtection: false
  recordSets:
    - name: "@"
      type: A
      records:
        - value:
            value: "93.184.216.34"
          comment: "primary web server"
        - value:
            value: "93.184.216.35"
          comment: "secondary web server"
    - name: "@"
      type: AAAA
      records:
        - value:
            value: "2606:2800:220:1:248:1893:25c8:1946"
    - name: www
      type: CNAME
      ttl: 3600
      records:
        - value:
            value: "example.com."
    - name: "@"
      type: MX
      records:
        - value:
            value: "10 mail.example.com."
        - value:
            value: "20 mail2.example.com."
    - name: "@"
      type: TXT
      records:
        - value:
            value: "\"v=spf1 include:_spf.google.com ~all\""
    - name: _dmarc
      type: TXT
      records:
        - value:
            value: "\"v=DMARC1; p=reject; rua=mailto:dmarc@example.com\""
    - name: "@"
      type: CAA
      records:
        - value:
            value: "0 issue \"letsencrypt.org\""
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainName` | `string` | yes |  |  |
| `spec.mode` | `enum` | yes |  |  |
| `spec.ttl` | `int32` |  | `3600` |  |
| `spec.primaryNameservers` | `[]PrimaryNameserver` |  |  |  |
| `spec.primaryNameservers[].address` | `string` | yes |  |  |
| `spec.primaryNameservers[].port` | `int32` |  | `53` |  |
| `spec.primaryNameservers[].tsigAlgorithm` | `string` |  |  |  |
| `spec.primaryNameservers[].tsigKey` | `string` |  |  |  |
| `spec.recordSets` | `[]RecordSet` |  |  |  |
| `spec.recordSets[].name` | `string` | yes |  |  |
| `spec.recordSets[].type` | `string` | yes |  |  |
| `spec.recordSets[].ttl` | `int32` |  |  |  |
| `spec.recordSets[].records` | `[]RecordValue` | yes |  |  |
| `spec.recordSets[].records[].value` | `string \| valueFrom` | yes |  |  |
| `spec.recordSets[].records[].comment` | `string` |  |  |  |
| `spec.deleteProtection` | `bool` |  |  |  |

## Field Details

### spec.domainName

`string` · required

The DNS domain name for the zone (e.g., "example.com").

This is the actual domain, not the Planton resource identifier
(which is metadata.name). The zone's Hetzner Cloud "name" field
is set to this value.

Changing this value forces replacement of the zone.

- rule: {"string":{"minLen":"1"}}

### spec.mode

`enum` · required

Zone mode. Determines whether Hetzner Cloud is the primary authority
for this domain or synchronizes records from an external primary.

Changing this value forces replacement of the zone.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `mode_unspecified`
- `primary` -- Hetzner Cloud is the primary (authoritative) nameserver. Records are managed directly through the record_sets field.
- `secondary` -- Hetzner Cloud acts as a secondary nameserver, synchronizing records from an external primary via zone transfer (AXFR/IXFR). The primary_nameservers field is required; record_sets must be empty.

### spec.ttl

`int32` · optional (explicit presence)

Default Time To Live (TTL) for records in the zone, in seconds.
Individual record sets can override this with their own TTL.

Default: 3600 (provider default when not specified).

- default: `3600`

### spec.primaryNameservers

`[]PrimaryNameserver`

Primary nameservers for zone transfer. Required when mode is
secondary; forbidden when mode is primary.

In secondary mode, Hetzner Cloud pulls DNS records from these
nameservers via AXFR/IXFR. At least one primary nameserver must
be specified.

### spec.primaryNameservers[].address

`string` · required

Public IPv4 or IPv6 address of the primary nameserver.

- rule: {"string":{"minLen":"1"}}

### spec.primaryNameservers[].port

`int32` · optional (explicit presence)

Port of the primary nameserver.

Default: 53

- default: `53`

### spec.primaryNameservers[].tsigAlgorithm

`string`

Transaction Signature (TSIG) algorithm for authenticating zone
transfers. Common values: "hmac-sha256", "hmac-sha512".
Leave empty if the primary does not require TSIG authentication.

### spec.primaryNameservers[].tsigKey

`string`

Transaction Signature (TSIG) key for authenticating zone transfers.
SENSITIVE: This is a shared secret between the primary and
secondary nameservers.
Leave empty if the primary does not require TSIG authentication.

### spec.recordSets

`[]RecordSet`

DNS record sets. Each entry manages all records for a unique
(name, type) pair using hcloud_zone_rrset.

Only valid when mode is primary. In secondary mode, records are
pulled from the primary nameserver and cannot be managed here.

### spec.recordSets[].name

`string` · required

Record name relative to the zone.

Use "@" for the zone apex (e.g., A record for "example.com"
itself), or a subdomain label (e.g., "www", "mail", "api").
Use "*" for wildcard records.

- rule: {"string":{"minLen":"1"}}

### spec.recordSets[].type

`string` · required

DNS record type. Standard types: "A", "AAAA", "CNAME", "MX",
"TXT", "NS", "SRV", "CAA", "PTR", "TLSA", "DS".

- rule: {"string":{"minLen":"1"}}

### spec.recordSets[].ttl

`int32` · optional (explicit presence)

TTL override for this record set, in seconds. When unset,
inherits the zone's default TTL.

### spec.recordSets[].records

`[]RecordValue` · required

Record values. At least one value is required per record set.

For most record types a single value suffices (CNAME, PTR).
Multiple values are used for round-robin A/AAAA records, multiple
MX servers, multiple TXT entries (e.g., SPF + DKIM), etc.

- rule: {"repeated":{"minItems":"1"}}

### spec.recordSets[].records[].value

`string | valueFrom` · required

The record value. Format depends on the record type:
  - A:     IPv4 address (e.g., "93.184.216.34")
  - AAAA:  IPv6 address (e.g., "2606:2800:220:1:248:1893:25c8:1946")
  - CNAME: Hostname with trailing dot (e.g., "example.com.")
  - MX:    Priority and hostname (e.g., "10 mail.example.com.")
  - TXT:   Quoted string (e.g., "\"v=spf1 include:_spf.google.com ~all\"")

Accepts a literal value or a reference to another component's
output via valueFrom for infra-chart composability. For example,
an A record can reference a HetznerCloudServer's IPv4 address:

  records:
    - value:
        valueFrom:
          kind: HetznerCloudServer
          name: web-1
          fieldPath: status.outputs.ipv4_address

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.recordSets[].records[].comment

`string`

Optional comment for this record.

### spec.deleteProtection

`bool`

Prevent accidental deletion of the zone via the Hetzner Cloud API.
When enabled, the zone cannot be deleted until protection is removed.

## Validation Rules

- `primary_forbids_primary_nameservers`: primary_nameservers must be empty when mode is primary
- `secondary_requires_primary_nameservers`: primary_nameservers is required when mode is secondary
- `secondary_forbids_record_sets`: record_sets must be empty when mode is secondary (records are pulled from the primary nameserver)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: HetznerCloudDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | The Hetzner Cloud numeric ID of the created zone (as a string). Can be referenced by other components via StringValueOrRef. |
| `status.outputs.nameservers` | `[]string` | The authoritative Hetzner nameservers assigned to the zone. Configure these NS records at your domain registrar to activate the zone. Example: ["helium.ns.hetzner.de", "hydrogen.ns.hetzner.com", "oxygen.ns.hetzner.com"] |

## See Also

- [Overview](../README.md)
