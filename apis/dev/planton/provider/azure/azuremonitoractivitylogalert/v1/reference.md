# AzureMonitorActivityLogAlert

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureMonitorActivityLogAlertSpec** defines an Azure Monitor Activity Log
Alert: a rule that fires an action (notify an action group, post a
webhook) when a matching entry appears in the Azure Activity Log -- the
subscription-level record of control-plane operations, service-health
incidents, resource-health transitions, policy events, and Advisor
recommendations. It is the only way to alert on the events that never
show up as metrics: a VM someone deleted, a region-wide service incident,
a resource going Unavailable, a policy denial.

The alert watches one or more scopes (a subscription, a resource group,
or specific resources) and matches on a criteria block whose `category`
selects which slice of the Activity Log to watch; the remaining criteria
fields narrow within that category. When an entry matches, every action
fires.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureMonitorActivityLogAlert
metadata:
  name: test-activity-log-alert
spec:
  resourceGroup:
    value: test-rg
  name: vm-delete-alert
  scopes:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg
  criteria:
    category: ADMINISTRATIVE
    operationName: Microsoft.Compute/virtualMachines/delete
    levels:
      - ERROR
      - CRITICAL
    statuses:
      - Succeeded
  actions:
    - actionGroupId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Insights/actionGroups/ops
      webhookProperties:
        env: test
  description: Fires when a VM is deleted in the resource group
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.location` | `enum` |  |  |  |
| `spec.scopes` | `[]string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_id`) |
| `spec.criteria` | `AzureMonitorActivityLogAlertCriteria` | yes |  |  |
| `spec.criteria.category` | `enum` | yes |  |  |
| `spec.criteria.operationName` | `string` |  |  |  |
| `spec.criteria.caller` | `string` |  |  |  |
| `spec.criteria.levels` | `[]enum` |  |  |  |
| `spec.criteria.resourceProviders` | `[]string` |  |  |  |
| `spec.criteria.resourceTypes` | `[]string` |  |  |  |
| `spec.criteria.resourceGroups` | `[]string` |  |  |  |
| `spec.criteria.resourceIds` | `[]string` |  |  |  |
| `spec.criteria.statuses` | `[]string` |  |  |  |
| `spec.criteria.subStatuses` | `[]string` |  |  |  |
| `spec.criteria.recommendationCategory` | `enum` |  |  |  |
| `spec.criteria.recommendationImpact` | `enum` |  |  |  |
| `spec.criteria.recommendationType` | `string` |  |  |  |
| `spec.criteria.resourceHealth` | `AzureMonitorActivityLogAlertResourceHealth` |  |  |  |
| `spec.criteria.resourceHealth.current` | `[]enum` |  |  |  |
| `spec.criteria.resourceHealth.previous` | `[]enum` |  |  |  |
| `spec.criteria.resourceHealth.reason` | `[]enum` |  |  |  |
| `spec.criteria.serviceHealth` | `AzureMonitorActivityLogAlertServiceHealth` |  |  |  |
| `spec.criteria.serviceHealth.events` | `[]enum` |  |  |  |
| `spec.criteria.serviceHealth.locations` | `[]string` |  |  |  |
| `spec.criteria.serviceHealth.services` | `[]string` |  |  |  |
| `spec.actions` | `[]AzureMonitorActivityLogAlertAction` |  |  |  |
| `spec.actions[].actionGroupId` | `string \| valueFrom` | yes |  | AzureMonitorActionGroup (`status.outputs.action_group_id`) |
| `spec.actions[].webhookProperties` | `map<string, string>` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.enabled` | `bool` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the alert is created in. Can be a literal
resource-group name or a reference to an AzureResourceGroup's name
output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the alert, unique within the resource group. Fixed at
creation.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.location

`enum`

The region the alert RESOURCE lives in. Activity Log Alerts are only
supported in a handful of regions; unspecified applies "global" (the
correct choice for virtually every alert -- the alert evaluates the
subscription-global Activity Log regardless). The others exist for data
residency of the alert definition itself.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_location_unspecified` -- Not specified: "global".
- `GLOBAL` -- Global -- the alert definition is not pinned to a region (the norm).
- `WEST_EUROPE` -- West Europe (for data residency of the alert definition).
- `NORTH_EUROPE` -- North Europe (for data residency of the alert definition).
- `EAST_US_2_EUAP` -- East US 2 EUAP (for data residency of the alert definition).

### spec.scopes

`[]string | valueFrom` · required

The scopes the alert watches -- a subscription, resource group, or
specific resources, by ARM ID. At least one required. Events under a
scope are evaluated (a subscription scope watches everything in it).
Each entry defaults to referencing an AzureResourceGroup's ARM ID (the
common resource-group scope); the scope is polymorphic, so a
subscription or any resource id works via an explicit valueFrom.

- references: AzureResourceGroup (`status.outputs.resource_group_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_id}} -- a bare string does not parse

### spec.criteria

`AzureMonitorActivityLogAlertCriteria` · required

The matching criteria -- which Activity Log entries fire the alert.
Required: an alert with no criteria matches nothing.

- rule: {"required":true}
- rule: recommendation_type cannot be combined with recommendation_category or recommendation_impact
- rule: resource_health cannot be combined with caller or service_health
- rule: service_health cannot be combined with caller or resource_health

### spec.criteria.category

`enum` · required

WHICH slice of the Activity Log to watch. Required.

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_category_unspecified`
- `ADMINISTRATIVE`
- `AUTOSCALE`
- `POLICY`
- `RECOMMENDATION`
- `RESOURCE_HEALTH`
- `SECURITY`
- `SERVICE_HEALTH`

### spec.criteria.operationName

`string`

The specific operation to match (e.g.
"Microsoft.Compute/virtualMachines/delete"). Empty matches any
operation in the category.

### spec.criteria.caller

`string`

The identity (user or service principal) that performed the operation.
Empty matches any caller. Cannot be combined with resource_health or
service_health (those categories have no caller).

### spec.criteria.levels

`[]enum`

The severity levels to match (Verbose, Informational, Warning, Error,
Critical). Empty matches any level.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_level_unspecified`
- `VERBOSE`
- `INFORMATIONAL`
- `WARNING`
- `ERROR`
- `CRITICAL`

### spec.criteria.resourceProviders

`[]string`

The resource providers to match (e.g. "Microsoft.Compute"). Empty
matches any.

### spec.criteria.resourceTypes

`[]string`

The resource types to match (e.g.
"Microsoft.Compute/virtualMachines"). Empty matches any.

### spec.criteria.resourceGroups

`[]string`

The resource groups (by name) to match. Empty matches any.

### spec.criteria.resourceIds

`[]string`

The specific resources (by ARM ID) to match. Empty matches any.

### spec.criteria.statuses

`[]string`

The operation statuses to match (e.g. "Succeeded", "Failed",
"Started"). Empty matches any.

### spec.criteria.subStatuses

`[]string`

The operation sub-statuses to match (e.g. "OK", "Created", HTTP-like
codes). Empty matches any.

### spec.criteria.recommendationCategory

`enum`

For the Recommendation category: the Advisor recommendation category to
match. Cannot be combined with recommendation_type.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_recommendation_category_unspecified`
- `COST`
- `RELIABILITY`
- `OPERATIONAL_EXCELLENCE`
- `PERFORMANCE`
- `HIGH_AVAILABILITY`
- `SECURITY_RECOMMENDATION`

### spec.criteria.recommendationImpact

`enum`

For the Recommendation category: the Advisor impact to match. Cannot be
combined with recommendation_type.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_recommendation_impact_unspecified`
- `HIGH`
- `MEDIUM`
- `LOW`

### spec.criteria.recommendationType

`string`

For the Recommendation category: match a specific recommendation type
id. Cannot be combined with recommendation_category or
recommendation_impact.

### spec.criteria.resourceHealth

`AzureMonitorActivityLogAlertResourceHealth`

For the ResourceHealth category: match resource-health transitions.
Cannot be combined with caller or service_health.

### spec.criteria.resourceHealth.current

`[]enum`

The current health states to match.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_health_status_unspecified`
- `AVAILABLE`
- `DEGRADED`
- `UNAVAILABLE`
- `UNKNOWN`

### spec.criteria.resourceHealth.previous

`[]enum`

The previous health states to match.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_health_status_unspecified`
- `AVAILABLE`
- `DEGRADED`
- `UNAVAILABLE`
- `UNKNOWN`

### spec.criteria.resourceHealth.reason

`[]enum`

The reasons for the health transition to match.

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_health_reason_unspecified`
- `PLATFORM_INITIATED`
- `USER_INITIATED`
- `REASON_UNKNOWN`

### spec.criteria.serviceHealth

`AzureMonitorActivityLogAlertServiceHealth`

For the ServiceHealth category: match Azure service-health incidents.
Cannot be combined with caller or resource_health.

### spec.criteria.serviceHealth.events

`[]enum`

The event types to match (Incident, Maintenance, Informational,
ActionRequired, Security).

Allowed values (use exactly as shown):

- `azure_monitor_activity_log_alert_service_health_event_unspecified`
- `INCIDENT`
- `MAINTENANCE`
- `EVENT_INFORMATIONAL`
- `ACTION_REQUIRED`
- `EVENT_SECURITY`

### spec.criteria.serviceHealth.locations

`[]string`

The Azure regions the incident affects (e.g. "East US"). Empty matches
any region.

### spec.criteria.serviceHealth.services

`[]string`

The Azure services the incident affects (e.g. "Virtual Machines").
Empty matches any service.

### spec.actions

`[]AzureMonitorActivityLogAlertAction`

The actions fired when an entry matches -- most commonly notifying an
action group. An alert with no actions still records matches but
notifies nobody.

### spec.actions[].actionGroupId

`string | valueFrom` · required

The action group to notify. Defaults to referencing an
AzureMonitorActionGroup's action_group_id output.

- references: AzureMonitorActionGroup (`status.outputs.action_group_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureMonitorActionGroup, name: <that resource's name>, fieldPath: status.outputs.action_group_id}} -- a bare string does not parse

### spec.actions[].webhookProperties

`map<string, string>`

Extra key-value properties passed to the action group's webhook
receivers when this alert fires.

### spec.description

`string`

A human-readable description of what the alert watches for and why.

### spec.enabled

`bool` · optional (explicit presence)

Whether the alert is active. Unspecified applies Azure's default
(enabled). Set false to keep the definition but stop it firing.

### spec.tags

`map<string, string>`

Free-form tags applied to the alert, merged over the Planton-derived
resource tags (organization, environment, resource id); a user tag with
the same key wins. Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMonitorActivityLogAlert, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.activity_log_alert_id` | `string` | The Azure Resource Manager ID of the activity log alert. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Insights/activityLogAlerts/{name} |
| `status.outputs.activity_log_alert_name` | `string` | The name of the activity log alert resource. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.scopes` | AzureResourceGroup | `status.outputs.resource_group_id` |
| `spec.actions[].actionGroupId` | AzureMonitorActionGroup | `status.outputs.action_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
