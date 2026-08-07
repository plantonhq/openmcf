# DigitalOceanDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanDnsZoneSpec defines the specification required to create a DNS zone (domain) on DigitalOcean.
This allows you to manage DNS records for a given domain via DigitalOcean's DNS service, focusing on the essential parameters (80/20 principle).

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDnsZone
metadata:
  name: first-dns-zone
spec:
  domainName: planton.com                # fully‑qualified domain
  records:
    # Root A record – points the zone apex to an IP address
    - name: "@"
      type: A
      values:
        - value: "203.0.113.10"
      ttlSeconds: 3600                   # optional; defaults to 3600

    # www CNAME record – redirects www.example.com to the apex
    - name: www
      type: CNAME
      values:
        - value: "@"

    # TXT record – typical SPF configuration (quotes required in YAML)
    - name: "@"
      type: TXT
      values:
        - value: "v=spf1 include:mail.example.com ~all"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.domainName` | `string` | yes |  |  |
| `spec.records` | `[]DigitalOceanDnsZoneRecord` |  |  |  |
| `spec.records[].name` | `string` | yes |  |  |
| `spec.records[].values` | `[]string \| valueFrom` | yes |  |  |
| `spec.records[].ttlSeconds` | `uint32` |  | `3600` |  |
| `spec.records[].type` | `enum` | yes |  |  |
| `spec.records[].priority` | `uint32` |  | `0` |  |
| `spec.records[].weight` | `uint32` |  | `0` |  |
| `spec.records[].port` | `uint32` |  | `0` |  |
| `spec.records[].flags` | `uint32` |  | `0` |  |
| `spec.records[].tag` | `string` |  |  |  |

## Field Details

### spec.domainName

`string` · required

The domain name for the DNS zone.
Must be a valid fully-qualified domain name (e.g., "example.com").

- rule: {"required":true,"string":{"pattern":"^(?:[A-Za-z0-9-]+\\.)+[A-Za-z]{2,}$"}}

### spec.records

`[]DigitalOceanDnsZoneRecord`

A list of DNS records to create within the zone (optional).
Each record includes its type, name, value(s), and TTL.

### spec.records[].name

`string` · required

The host/name for the DNS record, relative to the zone.
For root records, use "@" to denote the zone itself.

- rule: {"required":true}

### spec.records[].values

`[]string | valueFrom` · required

The value or values for the DNS record.
- For A/AAAA: one or more IP address(es).
- For CNAME: the target domain name.
- For TXT: the text data (if multiple strings, they will be concatenated by DNS).
- For MX: the mail server domain name (priority is specified separately in the priority field).
- For SRV: the target server (priority, weight, and port are specified separately).
- For CAA: the certificate authority domain (flags and tag are specified separately).
Each value can be a literal or a reference to another resource's output.

- rule: {"required":true,"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.records[].ttlSeconds

`uint32`

The time-to-live for this DNS record, in seconds.
Determines how long resolvers cache the record. Defaults to 3600 seconds (1 hour) if not set.

- default: `3600`

### spec.records[].type

`enum` · required

The type of the DNS record.
This field is required and must be one of the supported record types.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `unspecified`
- `A` -- host address
- `AAAA` -- ipv6 host address
- `ALIAS` -- auto resolved alias
- `CNAME` -- canonical name for an alias
- `MX` -- mail exchange
- `NS` -- name server
- `PTR` -- pointer
- `SOA` -- start of authority
- `SRV` -- location of service
- `TXT` -- descriptive text
- `CAA` -- certificate authority authorization

### spec.records[].priority

`uint32`

Priority for MX and SRV records.
- For MX records: Lower values indicate higher priority (e.g., 1 is higher priority than 10).
- For SRV records: Used in conjunction with weight for load distribution.
- Ignored for other record types.
Defaults to 0 if not specified.

- default: `0`

### spec.records[].weight

`uint32`

Weight for SRV records.
Specifies the relative weight for records with the same priority.
Higher weights are chosen more often.
Ignored for non-SRV record types.
Defaults to 0 if not specified.

- default: `0`

### spec.records[].port

`uint32`

Port for SRV records.
Specifies the TCP or UDP port on which the service is available.
Required for SRV records, ignored for other types.
Defaults to 0 if not specified.

- default: `0`

### spec.records[].flags

`uint32`

Flags for CAA records.
- 0: Non-critical (default) - If a CA doesn't understand the tag, it can ignore it.
- 128: Critical - If a CA doesn't understand the tag, it must refuse to issue.
Ignored for non-CAA record types.
Defaults to 0 if not specified.

- default: `0`

### spec.records[].tag

`string`

Tag for CAA records.
Specifies the property being authorized:
- "issue": Authorizes a CA to issue certificates for the domain.
- "issuewild": Authorizes a CA to issue wildcard certificates.
- "iodef": Specifies a URL to which a CA may report policy violations.
Required for CAA records, ignored for other types.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_name` | `string` | name of the DNS zone (domain) created on DigitalOcean. |
| `status.outputs.zone_id` | `string` | The unique identifier of the created DNS zone (typically the domain name or ID assigned by DigitalOcean). |
| `status.outputs.name_servers` | `[]string` | The list of nameserver addresses for the DNS zone. These are the nameservers that need to be set at the domain registrar for this zone. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanAppPlatformService | `spec.customDomain` | `spec.domain_name` |
| DigitalOceanDnsRecord | `spec.domain` | `status.outputs.zone_name` |

## See Also

- [Overview](../README.md)
