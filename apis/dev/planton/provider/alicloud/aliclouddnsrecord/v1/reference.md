# AliCloudDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudDnsRecordSpec defines the configuration for an Alibaba Cloud DNS
record managed by the Alidns service.

A DNS record maps a host record (subdomain) to a value within a parent
domain. The parent domain must already exist in Alidns -- either managed
by the AliCloudDnsZone component or added manually in the console.

Common record types: A (IPv4 address), AAAA (IPv6), CNAME (alias),
MX (mail exchange with priority), TXT (verification / SPF / DKIM),
NS (delegation), SRV (service locator), CAA (certificate authority).

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudDnsRecord
metadata:
  name: aliclouddnsrecord-demo
spec:
  region: cn-hangzhou
  domainName: demo.example.com
  rr: www
  type: A
  value: "203.0.113.10"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.domainName` | `string` | yes |  |  |
| `spec.rr` | `string` | yes |  |  |
| `spec.type` | `string` | yes |  |  |
| `spec.value` | `string` | yes |  |  |
| `spec.ttl` | `int32` |  |  |  |
| `spec.priority` | `int32` |  |  |  |
| `spec.line` | `string` |  |  |  |
| `spec.status` | `string` |  |  |  |
| `spec.remark` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider initialization.
Alidns is a global service, but the provider requires a region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.domainName

`string` · required

The parent domain name that this record belongs to.
Must already exist in the Alidns service (e.g., "example.com").
Cannot be changed after creation (ForceNew).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.rr

`string` · required

The host record (subdomain part) of the DNS record.
"rr" stands for "Resource Record" in the Alibaba Cloud API.
Examples: "www", "@" (apex), "*" (wildcard), "mail", "api.v2".
Maximum 253 characters; individual labels up to 63 characters.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.type

`string` · required

DNS record type.
Must be one of the types supported by the Alibaba Cloud Alidns API.
Note: "FORWORD_URL" is the Alibaba Cloud API's actual spelling.

- rule: type must be one of: A, AAAA, CNAME, MX, TXT, NS, SRV, CAA, REDIRECT_URL, FORWORD_URL
- rule: {"required":true}

### spec.value

`string` · required

The record value.
Interpretation depends on the record type:
  A       -> IPv4 address (e.g., "1.2.3.4")
  AAAA    -> IPv6 address (e.g., "2001:db8::1")
  CNAME   -> target domain (e.g., "cdn.example.com"), no trailing dot
  MX      -> mail server (e.g., "mx1.example.com"), no trailing dot
  TXT     -> text content (e.g., "v=spf1 include:example.com ~all")
  NS      -> nameserver (e.g., "ns1.example.com"), no trailing dot
  SRV     -> priority weight port target
  CAA     -> flags tag value (e.g., '0 issue "letsencrypt.org"')

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.ttl

`int32`

Time-to-live in seconds. Controls how long resolvers cache this record.
Range depends on the Alidns plan: Free [600-86400], Basic [120-86400],
Standard [60-86400], Ultimate [10-86400], Exclusive [1-86400].
Default: 600.

### spec.priority

`int32`

Priority for MX records. Required when type is "MX", ignored otherwise.
Range: 1 (highest priority) to 10 (lowest priority).

### spec.line

`string`

DNS resolution line. Controls which ISP or geographic line resolves
this record. Use "default" for standard resolution.
Advanced values include ISP-specific lines (e.g., "telecom", "unicom",
"mobile") and geographic lines for intelligent DNS routing.
Must be "default" when type is "FORWORD_URL".
Default: "default".

### spec.status

`string`

Record status. "ENABLE" activates the record; "DISABLE" keeps it
in Alidns but stops it from being served in DNS responses.
Useful for temporarily disabling a record without deleting it.
Default: "ENABLE".

- rule: status must be either ENABLE or DISABLE

### spec.remark

`string`

Remark or description for the record.
Visible in the Alidns console; useful for noting the record's purpose.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.record_id` | `string` | The record ID assigned by Alibaba Cloud. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
