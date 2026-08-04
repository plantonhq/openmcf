# CloudflareEmailRoutingRule

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `cloudflare.planton.dev/v1`

CloudflareEmailRoutingRuleSpec declares a single Email Routing rule for a zone:
match incoming mail (by recipient or all) and drop it, forward it to verified
destinations, or hand it to an Email Worker. Requires Email Routing to be
enabled on the zone (CloudflareEmailRoutingZone).

## Example

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareEmailRoutingRule
metadata:
  name: test-email-routing-rule
spec:
  zoneId:
    value: "023e105f4ecef8ad9ca31a8372d0c353"
  name: support-to-ops
  enabled: true
  priority: 0
  matchers:
    - type: literal
      field: to
      value: support@example.com
  action:
    type: forward
    forwardTo:
      - value: ops@example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.zoneId` | `string \| valueFrom` | yes |  | CloudflareDnsZone (`status.outputs.zone_id`) |
| `spec.name` | `string` |  |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.priority` | `int64` |  |  |  |
| `spec.matchers` | `[]CloudflareEmailRoutingRuleMatcher` | yes |  |  |
| `spec.matchers[].type` | `enum` |  |  |  |
| `spec.matchers[].field` | `string` |  |  |  |
| `spec.matchers[].value` | `string` |  |  |  |
| `spec.action` | `CloudflareEmailRoutingRuleAction` | yes |  |  |
| `spec.action.type` | `enum` |  |  |  |
| `spec.action.forwardTo` | `[]string \| valueFrom` |  |  | CloudflareEmailRoutingAddress (`status.outputs.email`) |
| `spec.action.worker` | `string \| valueFrom` |  |  | CloudflareWorker (`status.outputs.script_name`) |

## Field Details

### spec.zoneId

`string | valueFrom` · required

The zone this rule belongs to. A literal zone ID or a reference to a
CloudflareDnsZone resource.

- references: CloudflareDnsZone (`status.outputs.zone_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareDnsZone, name: <that resource's name>, fieldPath: status.outputs.zone_id}} -- a bare string does not parse

### spec.name

`string`

An optional descriptive name for the rule.

### spec.enabled

`bool` · optional (explicit presence)

Whether the rule is active. Defaults to true.

- default: `true`

### spec.priority

`int64`

Evaluation priority; lower numbers are evaluated first. Leave 0 for the
default.

- rule: priority must be zero or positive

### spec.matchers

`[]CloudflareEmailRoutingRuleMatcher` · required

The patterns that select which messages this rule applies to (at least one).

- rule: {"repeated":{"minItems":"1"}}
- rule: a literal matcher requires both field and value
- rule: an all matcher must not set field or value

### spec.matchers[].type

`enum`

The matcher type.

- rule: matcher type must be one of all, literal
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `matcher_type_unspecified`
- `all`
- `literal`

### spec.matchers[].field

`string`

The message field to match. Only "to" is supported. Required for literal
matchers; must be empty for all matchers.

- rule: field must be empty or 'to'

### spec.matchers[].value

`string`

The value to match (the recipient address for a literal "to" matcher).
Required for literal matchers; must be empty for all matchers.

### spec.action

`CloudflareEmailRoutingRuleAction` · required

The action taken on matched messages.

- rule: {"required":true}
- rule: forward_to is required when action type is forward
- rule: worker is required when action type is worker

### spec.action.type

`enum`

The action type.

- rule: action type must be one of drop, forward, worker
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `rule_action_type_unspecified`
- `drop`
- `forward`
- `worker`

### spec.action.forwardTo

`[]string | valueFrom`

Destination addresses to forward to (required when type is forward). Each is
a verified destination email or a reference to a CloudflareEmailRoutingAddress.

- references: CloudflareEmailRoutingAddress (`status.outputs.email`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareEmailRoutingAddress, name: <that resource's name>, fieldPath: status.outputs.email}} -- a bare string does not parse

### spec.action.worker

`string | valueFrom`

The Email Worker script to invoke (required when type is worker). A script
name or a reference to a CloudflareWorker.

- references: CloudflareWorker (`status.outputs.script_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: CloudflareWorker, name: <that resource's name>, fieldPath: status.outputs.script_name}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: CloudflareEmailRoutingRule, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.rule_id` | `string` | The Cloudflare-assigned identifier of the routing rule. |
| `status.outputs.zone_id` | `string` | The zone ID the rule belongs to. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.zoneId` | CloudflareDnsZone | `status.outputs.zone_id` |
| `spec.action.forwardTo` | CloudflareEmailRoutingAddress | `status.outputs.email` |
| `spec.action.worker` | CloudflareWorker | `status.outputs.script_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
