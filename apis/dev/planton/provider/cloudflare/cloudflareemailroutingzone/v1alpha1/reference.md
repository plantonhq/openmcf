# CloudflareEmailRoutingZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1alpha1`

CloudflareEmailRoutingZoneSpec enables Email Routing on a zone — the anchor for
the Email Routing family. Enabling provisions the zone's required MX/SPF/DKIM
DNS records automatically. The single per-zone catch-all rule is folded in
(spec.catch_all) since it is 1:1 with the zone; individual routing rules and
destination addresses are separate kinds (CloudflareEmailRoutingRule,
CloudflareEmailRoutingAddress).

## Example

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareEmailRoutingZone
metadata:
  name: test-email-routing-zone
spec:
  zoneId:
    value: "023e105f4ecef8ad9ca31a8372d0c353"
  catchAll:
    enabled: true
    type: drop
  lockDnsRecords: false
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneId` | `string \| valueFrom` | yes |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.catchAll` | `CloudflareEmailRoutingZoneCatchAll` |  |  |  |
| `spec.catchAll.enabled` | `bool` |  |  |  |
| `spec.catchAll.type` | `enum` |  |  |  |
| `spec.catchAll.forwardTo` | `[]string \| valueFrom` |  |  | CloudflareEmailRoutingAddress (`status.outputs.email`) |
| `spec.catchAll.worker` | `string \| valueFrom` |  |  | CloudflareWorker (`status.outputs.script_name`) |
| `spec.lockDnsRecords` | `bool` |  |  |  |

## Field Details

### spec.zoneId

`string | valueFrom` · required

The zone to enable Email Routing on. A literal zone ID or a reference to a
CloudflareDnsZone resource.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.catchAll

`CloudflareEmailRoutingZoneCatchAll`

Optional catch-all rule for mail no other rule matched. Omit to leave the
zone's catch-all at Cloudflare's default (drop, disabled).

- rule: forward_to is required when catch-all type is forward
- rule: worker is required when catch-all type is worker

### spec.catchAll.enabled

`bool`

Whether the catch-all rule is active.

### spec.catchAll.type

`enum`

What to do with unmatched mail.

- rule: catch-all type must be one of drop, forward, worker
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `catch_all_action_type_unspecified`
- `drop`
- `forward`
- `worker`

### spec.catchAll.forwardTo

`[]string | valueFrom`

Destination addresses to forward to (required when type is forward). Each is
a verified destination email, or a reference to a CloudflareEmailRoutingAddress.

- references: CloudflareEmailRoutingAddress (`status.outputs.email`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareEmailRoutingAddress, name: <that resource's name>, fieldPath: status.outputs.email}} -- a bare string does not parse

### spec.catchAll.worker

`string | valueFrom`

The Email Worker script to invoke (required when type is worker). A script
name or a reference to a CloudflareWorker.

- references: CloudflareWorker (`status.outputs.script_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareWorker, name: <that resource's name>, fieldPath: status.outputs.script_name}} -- a bare string does not parse

### spec.lockDnsRecords

`bool`

Lock the Email Routing DNS records so they cannot be modified out-of-band.
The records are created automatically on enable regardless; locking manages
them explicitly. Defaults to false.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareEmailRoutingZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.zone_id` | `string` | The zone ID Email Routing was enabled on. |
| `status.outputs.enabled` | `bool` | Whether Email Routing is enabled on the zone. |
| `status.outputs.status` | `string` | The Email Routing configuration status (e.g. ready, unconfigured, misconfigured). |
| `status.outputs.name` | `string` | The zone's domain name. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneId` | CloudflareDnsZone | `status.outputs.zone_id` |
| `spec.catchAll.forwardTo` | CloudflareEmailRoutingAddress | `status.outputs.email` |
| `spec.catchAll.worker` | CloudflareWorker | `status.outputs.script_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
