# AzureServiceBusSubscription -- Design Research

## The Resource

A Service Bus subscription
(`Microsoft.ServiceBus/namespaces/topics/subscriptions`) is a consuming
view of a topic; its rules
(`…/subscriptions/rules`) decide which messages it admits. The
component maps onto `azurerm_servicebus_subscription` +
`azurerm_servicebus_subscription_rule` (azurerm v4.x,
`internal/services/servicebus/`), parity-verified against pulumi-azure
v6 (`servicebus.Subscription` + `servicebus.SubscriptionRule`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `topic_id` | `topic_id` | The v4 typed parent id; single authoritative parent FK, ForceNew |
| `name` | `subscription_name` | Required, ForceNew, provider regex (≤50, `[-._a-zA-Z0-9]`, underscore-friendly ends) |
| `max_delivery_count` | same | REQUIRED in azurerm (no default) -- kept required |
| `lock_duration` / `default_message_ttl` / `auto_delete_on_idle` | same | Unset leaves Azure's defaults |
| `dead_lettering_on_message_expiration` | same | Default false |
| `dead_lettering_on_filter_evaluation_error` | same | optional bool, default true (Azure's default) |
| `requires_session` | same | ForceNew |
| `batched_operations_enabled` | same | Subscription-side default false (differs from queue/topic -- azurerm's own shape) |
| `forward_to` / `forward_dead_lettered_messages_to` | same | Entity NAMES in the same namespace |
| `status` | enum | ACTIVE/DISABLED/RECEIVE_DISABLED (the provider's set for subscriptions) |
| `client_scoped_subscription_enabled` + block | `client_scoped_subscription` | Presence of the block IS the toggle; ForceNew identity fields; the internal `{name}$client$D` entity naming round-trips in the provider |
| rule `filter_type` | enum | SqlFilter/CorrelationFilter, XOR contracts as CELs |
| rule `sql_filter` | same | ≤1024 (provider validator) |
| rule `correlation_filter` | message | Nine matchers + properties map; the provider's at-least-one contract as a message CEL |
| rule `action` | same | SQL annotation on matched messages |

Outputs: `subscription_id`, `subscription_name`, `topic_name`,
`namespace_name`.

## Decomposition Decisions

- **Rules FOLD into the subscription.** A rule has no life outside its
  subscription, nothing FK-references one, and the rule set is one
  filtering policy configured as a unit -- the fold test passes on
  every axis. azurerm's standalone rule resource is realized inside
  both modules.
- **The `$Default` contract (live-verified):** Azure auto-creates a
  catch-all TrueFilter rule named `$Default` on every new subscription,
  and rules combine with OR semantics -- declared rules are ADDITIVE
  alongside it. The catch-all CANNOT be declared: both providers'
  create paths run an import-existence check and refuse to adopt the
  service-created rule (confirmed live on both engines -- "a resource
  with the ID .../rules/$Default already exists ... needs to be
  imported"). The spec reserves the name with a CEL and teaches the
  one-time out-of-band removal
  (`az servicebus topic subscription rule delete --name '$Default'`)
  for restrictive delivery.

## Recorded Skips (with reasons)

- **`sql_filter_compatibility_level`** -- a Computed, reserved output
  in azurerm (hard-coded 20 by the service); not a knob.
- **Rule-name uniqueness within the list** is not CEL-enforceable
  (list-comprehension uniqueness is outside protovalidate's CEL
  subset); duplicate names collapse into one upserted rule at ARM, and
  the modules key resources by rule name so a duplicate fails the plan
  on both engines.

## Operational Behavior Worth Knowing

- **Renaming a subscription resets its read position** -- the new
  subscription starts from the stream's tail; undelivered messages in
  the old one are lost with it.
- **A filter that throws on a malformed message dead-letters it** when
  `dead_lettering_on_filter_evaluation_error` is true (Azure's
  default) -- watch the dead-letter sub-queue when producers evolve
  schemas.

## Composition

- `topic_id` → `AzureServiceBusTopic.status.outputs.topic_id`
- `forward_to` → a queue/topic NAME in the same namespace (compose with
  the target's `queue_name`/`topic_name` output via explicit valueFrom)
