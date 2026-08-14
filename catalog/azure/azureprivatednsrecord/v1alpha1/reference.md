# AzurePrivateDnsRecord

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Offline-plan test manifest. Exercises a structured payload at depth:
# an apex MX record set (name "@") with two entries at different
# preferences, an explicit non-default TTL, and user tags merged over
# the derived ones. The six other payload types are pure ARM property
# writes on the same create path -- the offline plans and the live lane
# cover them per the profile's record.
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateDnsRecord
metadata:
  name: test-private-dns-record
  org: test-org
  env: dev
spec:
  privateDnsZoneId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/platform-rg/providers/Microsoft.Network/privateDnsZones/internal.contoso.com
  name: "@"
  ttlSeconds: 3600
  mx:
    - preference: 10
      exchange: mail1.internal.contoso.com
    - preference: 20
      exchange: mail2.internal.contoso.com
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.privateDnsZoneId` | `string \| valueFrom` | yes |  | AzurePrivateDnsZone (`status.outputs.zone_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.ttlSeconds` | `int32` |  | `300` |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.a` | `[]string` |  |  |  |
| `spec.aaaa` | `[]string` |  |  |  |
| `spec.cname` | `string \| valueFrom` |  |  |  |
| `spec.mx` | `[]AzurePrivateDnsMxEntry` |  |  |  |
| `spec.mx[].preference` | `int32` | yes |  |  |
| `spec.mx[].exchange` | `string` | yes |  |  |
| `spec.ptr` | `[]string` |  |  |  |
| `spec.srv` | `[]AzurePrivateDnsSrvEntry` |  |  |  |
| `spec.srv[].priority` | `int32` | yes |  |  |
| `spec.srv[].weight` | `int32` | yes |  |  |
| `spec.srv[].port` | `int32` | yes |  |  |
| `spec.srv[].target` | `string` | yes |  |  |
| `spec.txt` | `[]string \| valueFrom` |  |  |  |

## Field Details

### spec.privateDnsZoneId

`string | valueFrom` · required

- references: AzurePrivateDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzurePrivateDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Record name must be '@' for the zone apex, '*' (optionally '*.<name>') for wildcards, or dot-separated labels of lowercase letters, digits, hyphens, and underscores -- e.g. db, api.v1, _sip._tcp
- rule: {"required":true}

### spec.ttlSeconds

`int32` · optional (explicit presence)

- default: `300`
- rule: {"int32":{"lte":2147483647,"gte":0}}

### spec.tags

`map<string, string>`

### spec.a

`[]string`

- rule: {"repeated":{"maxItems":"20","items":{"string":{"ipv4":true}}}}

### spec.aaaa

`[]string`

- rule: {"repeated":{"items":{"string":{"ipv6":true}}}}

### spec.cname

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.mx

`[]AzurePrivateDnsMxEntry`

### spec.mx[].preference

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":65535,"gte":0}}

### spec.mx[].exchange

`string` · required

- rule: {"required":true,"string":{"maxLen":"253"}}

### spec.ptr

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1","maxLen":"253"}}}}

### spec.srv

`[]AzurePrivateDnsSrvEntry`

### spec.srv[].priority

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":65535,"gte":0}}

### spec.srv[].weight

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":65535,"gte":0}}

### spec.srv[].port

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":65535,"gte":1}}

### spec.srv[].target

`string` · required

- rule: {"required":true,"string":{"maxLen":"253"}}

### spec.txt

`[]string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

## Validation Rules

- `azure_private_dns_record_exactly_one_payload`: Set exactly one record payload -- a, aaaa, cname, mx, ptr, srv, or txt -- the payload determines the record type

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePrivateDnsRecord, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.record_id` | `string` |  |
| `status.outputs.fqdn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.privateDnsZoneId` | AzurePrivateDnsZone | `status.outputs.zone_id` |

## See Also

- [Overview](../README.md)
