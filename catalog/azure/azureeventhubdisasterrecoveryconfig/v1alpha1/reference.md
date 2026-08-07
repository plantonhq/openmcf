# AzureEventHubDisasterRecoveryConfig

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureEventHubDisasterRecoveryConfigSpec** defines a geo-disaster-
recovery pairing between two Event Hubs namespaces: metadata (hubs,
consumer groups, authorization rules -- not event data) continuously
replicates from the primary to the partner, and a failover-stable ALIAS
DNS name fronts whichever namespace is currently primary.

Clients connect through the alias
(`{alias_name}.servicebus.windows.net`) instead of either namespace, so
a failover needs no client reconfiguration. Alias-addressed connection
strings surface on the namespace's and the authorization rule kinds'
`*_connection_string_alias` outputs once a pairing exists.

**Azure's pairing contracts** (apply-time -- they involve both live
namespaces): the namespaces must be in DIFFERENT regions, on the same
tier (STANDARD or higher; geo-DR is not available on BASIC), and the
partner must be empty (no hubs) when pairing is created. Changing
partner_namespace_id breaks the existing pairing and re-pairs to the
new partner.

**Failover is an operational action, not a config change**: triggered
from the SECONDARY side (portal/CLI/SDK) during a regional incident;
after failover the alias points at the former partner. Deleting this
resource breaks the pairing gracefully -- both namespaces keep running
independently, and the alias name is released once Azure finishes the
break-pair/name-release choreography (the modules wait it out).

**ForceNew fields**: `alias_name`, `primary_namespace_id`.

## Example

```yaml
# Offline-plan manifest: an Event Hubs geo-DR pairing between two
# namespaces -- exercises the primary ARM-id parsing seam (resource
# group + namespace name) and the partner pass-through.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventHubDisasterRecoveryConfig
metadata:
  name: test-eh-geo-dr
spec:
  aliasName: hack-eventhub-alias
  primaryNamespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.EventHub/namespaces/my-primary-ehns
  partnerNamespaceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-rg/providers/Microsoft.EventHub/namespaces/my-partner-ehns
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.aliasName` | `string` | yes |  |  |
| `spec.primaryNamespaceId` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_id`) |
| `spec.partnerNamespaceId` | `string \| valueFrom` | yes |  | AzureEventHubNamespace (`status.outputs.namespace_id`) |

## Field Details

### spec.aliasName

`string` · required

The alias name -- becomes the failover-stable DNS name
`{alias_name}.servicebus.windows.net`, so it shares the namespace
name uniqueness scope (globally unique; it cannot collide with any
existing namespace name either). Up to 60 characters of letters,
numbers, periods, hyphens, and underscores, starting and ending with
a letter or number.

**ForceNew**: changing the alias replaces the pairing.

- rule: alias_name must be up to 60 characters of letters, numbers, periods, hyphens, and underscores, starting and ending with a letter or number
- rule: {"required":true,"string":{"minLen":"1","maxLen":"60"}}

### spec.primaryNamespaceId

`string | valueFrom` · required

The PRIMARY namespace -- the one clients actively use, by ARM ID.
References an AzureEventHubNamespace's namespace_id output. The
pairing lives under this namespace. Fixed at creation.

- references: AzureEventHubNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

### spec.partnerNamespaceId

`string | valueFrom` · required

The PARTNER (secondary) namespace metadata replicates to, by ARM ID.
Must be in a DIFFERENT region than the primary, on the same tier,
and empty (no hubs) at pairing time -- Azure validates all three
when the pairing is created. Changing it breaks the current pairing
and re-pairs to the new partner.

- references: AzureEventHubNamespace (`status.outputs.namespace_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHubNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventHubDisasterRecoveryConfig, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.disaster_recovery_config_id` | `string` | The Azure Resource Manager ID of the disaster-recovery config (under the primary namespace). Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.EventHub/namespaces/{ns}/disasterRecoveryConfigs/{alias} |
| `status.outputs.alias_name` | `string` | The alias name -- the failover-stable DNS identity `{alias_name}.servicebus.windows.net`. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.primaryNamespaceId` | AzureEventHubNamespace | `status.outputs.namespace_id` |
| `spec.partnerNamespaceId` | AzureEventHubNamespace | `status.outputs.namespace_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
