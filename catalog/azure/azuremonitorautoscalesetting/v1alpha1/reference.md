# AzureMonitorAutoscaleSetting

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Offline-plan test manifest. Exercises the DEEP shape the live smoke
# does not: three profiles (metric rules with dimensions, namespace, and
# per-instance division; a timezone'd recurrence; a fixed-date window),
# predictive autoscale, and both notification channels with webhook
# properties. (The live scenario proves the everyday shape: default
# profile + rule + recurrence + email.)
apiVersion: azure.planton.dev/v1alpha1
kind: AzureMonitorAutoscaleSetting
metadata:
  name: test-autoscale-setting
  org: test-org
  env: dev
spec:
  resourceGroup:
    value: test-rg
  name: web-pool-autoscale
  region: eastus
  targetResourceId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachineScaleSets/web-pool
  enabled: true
  predictive:
    scaleMode: ForecastOnly
    lookAheadTime: PT10M
  profiles:
    - name: default
      capacity:
        minimum: 2
        maximum: 10
        default: 3
      rules:
        - metricTrigger:
            metricName: Percentage CPU
            metricResourceId:
              value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachineScaleSets/web-pool
            timeGrain: PT1M
            statistic: Average
            timeWindow: PT10M
            timeAggregation: Average
            operator: GreaterThan
            threshold: 75
            metricNamespace: microsoft.compute/virtualmachinescalesets
            divideByInstanceCount: true
            dimensions:
              - name: VMName
                operator: NotEquals
                values:
                  - canary-0
          scaleAction:
            direction: Increase
            type: PercentChangeCount
            value: 20
            cooldown: PT5M
        - metricTrigger:
            metricName: Percentage CPU
            metricResourceId:
              value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Compute/virtualMachineScaleSets/web-pool
            timeGrain: PT1M
            statistic: Average
            timeWindow: PT20M
            timeAggregation: Average
            operator: LessThan
            threshold: 25
          scaleAction:
            direction: Decrease
            type: ChangeCount
            value: 1
            cooldown: PT15M
    - name: weekday-business-hours
      capacity:
        minimum: 4
        maximum: 12
        default: 4
      recurrence:
        timezone: Eastern Standard Time
        days: [Monday, Tuesday, Wednesday, Thursday, Friday]
        hour: 8
        minute: 30
    - name: launch-window
      capacity:
        minimum: 8
        maximum: 20
        default: 8
      fixedDate:
        timezone: Pacific Standard Time
        start: "2026-11-27T00:00:00Z"
        end: "2026-11-30T00:00:00Z"
  notification:
    email:
      # Explicit recipients are the only email wiring ARM still accepts --
      # the classic admin/co-admin flags were retired in April 2024.
      customEmails:
        - oncall@example.com
    webhooks:
      - serviceUri: https://hooks.example.com/scale-events
        properties:
          channel: ops
  tags:
    team: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.targetResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.predictive` | `AzureMonitorAutoscaleSettingPredictive` |  |  |  |
| `spec.predictive.scaleMode` | `string` | yes |  |  |
| `spec.predictive.lookAheadTime` | `string` |  |  |  |
| `spec.profiles` | `[]AzureMonitorAutoscaleSettingProfile` | yes |  |  |
| `spec.profiles[].name` | `string` | yes |  |  |
| `spec.profiles[].capacity` | `AzureMonitorAutoscaleSettingCapacity` | yes |  |  |
| `spec.profiles[].capacity.minimum` | `int32` | yes |  |  |
| `spec.profiles[].capacity.maximum` | `int32` | yes |  |  |
| `spec.profiles[].capacity.default` | `int32` | yes |  |  |
| `spec.profiles[].rules` | `[]AzureMonitorAutoscaleSettingRule` |  |  |  |
| `spec.profiles[].rules[].metricTrigger` | `AzureMonitorAutoscaleSettingMetricTrigger` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.metricName` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.metricResourceId` | `string \| valueFrom` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.timeGrain` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.statistic` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.timeWindow` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.timeAggregation` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.operator` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.threshold` | `double` |  |  |  |
| `spec.profiles[].rules[].metricTrigger.metricNamespace` | `string` |  |  |  |
| `spec.profiles[].rules[].metricTrigger.divideByInstanceCount` | `bool` |  |  |  |
| `spec.profiles[].rules[].metricTrigger.dimensions` | `[]AzureMonitorAutoscaleSettingDimension` |  |  |  |
| `spec.profiles[].rules[].metricTrigger.dimensions[].name` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.dimensions[].operator` | `string` | yes |  |  |
| `spec.profiles[].rules[].metricTrigger.dimensions[].values` | `[]string` | yes |  |  |
| `spec.profiles[].rules[].scaleAction` | `AzureMonitorAutoscaleSettingScaleAction` | yes |  |  |
| `spec.profiles[].rules[].scaleAction.direction` | `string` | yes |  |  |
| `spec.profiles[].rules[].scaleAction.type` | `string` | yes |  |  |
| `spec.profiles[].rules[].scaleAction.value` | `int32` | yes |  |  |
| `spec.profiles[].rules[].scaleAction.cooldown` | `string` | yes |  |  |
| `spec.profiles[].fixedDate` | `AzureMonitorAutoscaleSettingFixedDate` |  |  |  |
| `spec.profiles[].fixedDate.timezone` | `string` |  | `UTC` |  |
| `spec.profiles[].fixedDate.start` | `string` | yes |  |  |
| `spec.profiles[].fixedDate.end` | `string` | yes |  |  |
| `spec.profiles[].recurrence` | `AzureMonitorAutoscaleSettingRecurrence` |  |  |  |
| `spec.profiles[].recurrence.timezone` | `string` |  | `UTC` |  |
| `spec.profiles[].recurrence.days` | `[]string` | yes |  |  |
| `spec.profiles[].recurrence.hour` | `int32` | yes |  |  |
| `spec.profiles[].recurrence.minute` | `int32` | yes |  |  |
| `spec.notification` | `AzureMonitorAutoscaleSettingNotification` |  |  |  |
| `spec.notification.email` | `AzureMonitorAutoscaleSettingNotificationEmail` |  |  |  |
| `spec.notification.email.customEmails` | `[]string` |  |  |  |
| `spec.notification.webhooks` | `[]AzureMonitorAutoscaleSettingNotificationWebhook` |  |  |  |
| `spec.notification.webhooks[].serviceUri` | `string` | yes |  |  |
| `spec.notification.webhooks[].properties` | `map<string, string>` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.targetResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.enabled

`bool` · optional (explicit presence)

- default: `true`

### spec.predictive

`AzureMonitorAutoscaleSettingPredictive`

### spec.predictive.scaleMode

`string` · required

- rule: {"required":true,"string":{"in":["Enabled","ForecastOnly"]}}

### spec.predictive.lookAheadTime

`string`

- rule: look_ahead_time must be an ISO 8601 duration between PT1M and PT1H, e.g. PT10M

### spec.profiles

`[]AzureMonitorAutoscaleSettingProfile` · required

- rule: {"repeated":{"minItems":"1","maxItems":"20"}}
- rule: a profile is default (no schedule), fixed-date, or recurring -- set at most one of fixed_date and recurrence

### spec.profiles[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.profiles[].capacity

`AzureMonitorAutoscaleSettingCapacity` · required

- rule: {"required":true}

### spec.profiles[].capacity.minimum

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":1000,"gte":0}}

### spec.profiles[].capacity.maximum

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":1000,"gte":0}}

### spec.profiles[].capacity.default

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":1000,"gte":0}}

### spec.profiles[].rules

`[]AzureMonitorAutoscaleSettingRule`

- rule: {"repeated":{"maxItems":"10"}}

### spec.profiles[].rules[].metricTrigger

`AzureMonitorAutoscaleSettingMetricTrigger` · required

- rule: {"required":true}

### spec.profiles[].rules[].metricTrigger.metricName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.profiles[].rules[].metricTrigger.metricResourceId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.profiles[].rules[].metricTrigger.timeGrain

`string` · required

- rule: time_grain must be an ISO 8601 duration, e.g. PT1M
- rule: {"required":true}

### spec.profiles[].rules[].metricTrigger.statistic

`string` · required

- rule: {"required":true,"string":{"in":["Average","Max","Min","Sum"]}}

### spec.profiles[].rules[].metricTrigger.timeWindow

`string` · required

- rule: time_window must be an ISO 8601 duration, e.g. PT10M
- rule: {"required":true}

### spec.profiles[].rules[].metricTrigger.timeAggregation

`string` · required

- rule: {"required":true,"string":{"in":["Average","Count","Maximum","Minimum","Total","Last"]}}

### spec.profiles[].rules[].metricTrigger.operator

`string` · required

- rule: {"required":true,"string":{"in":["Equals","NotEquals","GreaterThan","GreaterThanOrEqual","LessThan","LessThanOrEqual"]}}

### spec.profiles[].rules[].metricTrigger.threshold

`double`

### spec.profiles[].rules[].metricTrigger.metricNamespace

`string`

### spec.profiles[].rules[].metricTrigger.divideByInstanceCount

`bool`

### spec.profiles[].rules[].metricTrigger.dimensions

`[]AzureMonitorAutoscaleSettingDimension`

### spec.profiles[].rules[].metricTrigger.dimensions[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.profiles[].rules[].metricTrigger.dimensions[].operator

`string` · required

- rule: {"required":true,"string":{"in":["Equals","NotEquals"]}}

### spec.profiles[].rules[].metricTrigger.dimensions[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.profiles[].rules[].scaleAction

`AzureMonitorAutoscaleSettingScaleAction` · required

- rule: {"required":true}

### spec.profiles[].rules[].scaleAction.direction

`string` · required

- rule: {"required":true,"string":{"in":["Increase","Decrease"]}}

### spec.profiles[].rules[].scaleAction.type

`string` · required

- rule: {"required":true,"string":{"in":["ChangeCount","ExactCount","PercentChangeCount","ServiceAllowedNextValue"]}}

### spec.profiles[].rules[].scaleAction.value

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"gte":0}}

### spec.profiles[].rules[].scaleAction.cooldown

`string` · required

- rule: cooldown must be an ISO 8601 duration, e.g. PT5M
- rule: {"required":true}

### spec.profiles[].fixedDate

`AzureMonitorAutoscaleSettingFixedDate`

### spec.profiles[].fixedDate.timezone

`string` · optional (explicit presence)

- default: `UTC`
- rule: {"string":{"in":["Dateline Standard Time","UTC-11","Hawaiian Standard Time","Alaskan Standard Time","Pacific Standard Time (Mexico)","Pacific Standard Time","US Mountain Standard Time","Mountain Standard Time (Mexico)","Mountain Standard Time","Central America Standard Time","Central Standard Time","Central Standard Time (Mexico)","Canada Central Standard Time","SA Pacific Standard Time","Eastern Standard Time","US Eastern Standard Time","Venezuela Standard Time","Paraguay Standard Time","Atlantic Standard Time","Central Brazilian Standard Time","SA Western Standard Time","Pacific SA Standard Time","Newfoundland Standard Time","E. South America Standard Time","Argentina Standard Time","SA Eastern Standard Time","Greenland Standard Time","Montevideo Standard Time","Bahia Standard Time","UTC-02","Mid-Atlantic Standard Time","Azores Standard Time","Cape Verde Standard Time","Morocco Standard Time","UTC","GMT Standard Time","Greenwich Standard Time","W. Europe Standard Time","Central Europe Standard Time","Romance Standard Time","Central European Standard Time","W. Central Africa Standard Time","Namibia Standard Time","Jordan Standard Time","GTB Standard Time","Middle East Standard Time","Egypt Standard Time","Syria Standard Time","E. Europe Standard Time","South Africa Standard Time","FLE Standard Time","Turkey Standard Time","Israel Standard Time","Kaliningrad Standard Time","Libya Standard Time","Arabic Standard Time","Arab Standard Time","Belarus Standard Time","Russian Standard Time","E. Africa Standard Time","Iran Standard Time","Arabian Standard Time","Azerbaijan Standard Time","Russia Time Zone 3","Mauritius Standard Time","Georgian Standard Time","Caucasus Standard Time","Afghanistan Standard Time","West Asia Standard Time","Ekaterinburg Standard Time","Pakistan Standard Time","India Standard Time","Sri Lanka Standard Time","Nepal Standard Time","Central Asia Standard Time","Bangladesh Standard Time","N. Central Asia Standard Time","Myanmar Standard Time","SE Asia Standard Time","North Asia Standard Time","China Standard Time","North Asia East Standard Time","Singapore Standard Time","W. Australia Standard Time","Taipei Standard Time","Ulaanbaatar Standard Time","Tokyo Standard Time","Korea Standard Time","Yakutsk Standard Time","Cen. Australia Standard Time","AUS Central Standard Time","E. Australia Standard Time","AUS Eastern Standard Time","West Pacific Standard Time","Tasmania Standard Time","Magadan Standard Time","Vladivostok Standard Time","Russia Time Zone 10","Central Pacific Standard Time","Russia Time Zone 11","New Zealand Standard Time","UTC+12","Fiji Standard Time","Kamchatka Standard Time","Tonga Standard Time","Samoa Standard Time","Line Islands Standard Time"]}}

### spec.profiles[].fixedDate.start

`string` · required

- rule: start must be an RFC 3339 timestamp, e.g. 2026-11-27T00:00:00Z
- rule: {"required":true}

### spec.profiles[].fixedDate.end

`string` · required

- rule: end must be an RFC 3339 timestamp, e.g. 2026-11-30T00:00:00Z
- rule: {"required":true}

### spec.profiles[].recurrence

`AzureMonitorAutoscaleSettingRecurrence`

### spec.profiles[].recurrence.timezone

`string` · optional (explicit presence)

- default: `UTC`
- rule: {"string":{"in":["Dateline Standard Time","UTC-11","Hawaiian Standard Time","Alaskan Standard Time","Pacific Standard Time (Mexico)","Pacific Standard Time","US Mountain Standard Time","Mountain Standard Time (Mexico)","Mountain Standard Time","Central America Standard Time","Central Standard Time","Central Standard Time (Mexico)","Canada Central Standard Time","SA Pacific Standard Time","Eastern Standard Time","US Eastern Standard Time","Venezuela Standard Time","Paraguay Standard Time","Atlantic Standard Time","Central Brazilian Standard Time","SA Western Standard Time","Pacific SA Standard Time","Newfoundland Standard Time","E. South America Standard Time","Argentina Standard Time","SA Eastern Standard Time","Greenland Standard Time","Montevideo Standard Time","Bahia Standard Time","UTC-02","Mid-Atlantic Standard Time","Azores Standard Time","Cape Verde Standard Time","Morocco Standard Time","UTC","GMT Standard Time","Greenwich Standard Time","W. Europe Standard Time","Central Europe Standard Time","Romance Standard Time","Central European Standard Time","W. Central Africa Standard Time","Namibia Standard Time","Jordan Standard Time","GTB Standard Time","Middle East Standard Time","Egypt Standard Time","Syria Standard Time","E. Europe Standard Time","South Africa Standard Time","FLE Standard Time","Turkey Standard Time","Israel Standard Time","Kaliningrad Standard Time","Libya Standard Time","Arabic Standard Time","Arab Standard Time","Belarus Standard Time","Russian Standard Time","E. Africa Standard Time","Iran Standard Time","Arabian Standard Time","Azerbaijan Standard Time","Russia Time Zone 3","Mauritius Standard Time","Georgian Standard Time","Caucasus Standard Time","Afghanistan Standard Time","West Asia Standard Time","Ekaterinburg Standard Time","Pakistan Standard Time","India Standard Time","Sri Lanka Standard Time","Nepal Standard Time","Central Asia Standard Time","Bangladesh Standard Time","N. Central Asia Standard Time","Myanmar Standard Time","SE Asia Standard Time","North Asia Standard Time","China Standard Time","North Asia East Standard Time","Singapore Standard Time","W. Australia Standard Time","Taipei Standard Time","Ulaanbaatar Standard Time","Tokyo Standard Time","Korea Standard Time","Yakutsk Standard Time","Cen. Australia Standard Time","AUS Central Standard Time","E. Australia Standard Time","AUS Eastern Standard Time","West Pacific Standard Time","Tasmania Standard Time","Magadan Standard Time","Vladivostok Standard Time","Russia Time Zone 10","Central Pacific Standard Time","Russia Time Zone 11","New Zealand Standard Time","UTC+12","Fiji Standard Time","Kamchatka Standard Time","Tonga Standard Time","Samoa Standard Time","Line Islands Standard Time"]}}

### spec.profiles[].recurrence.days

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"in":["Monday","Tuesday","Wednesday","Thursday","Friday","Saturday","Sunday"]}}}}

### spec.profiles[].recurrence.hour

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":23,"gte":0}}

### spec.profiles[].recurrence.minute

`int32` · required · optional (explicit presence)

- rule: {"required":true,"int32":{"lte":59,"gte":0}}

### spec.notification

`AzureMonitorAutoscaleSettingNotification`

- rule: configure at least one notification channel: email and/or webhooks

### spec.notification.email

`AzureMonitorAutoscaleSettingNotificationEmail`

- rule: email notifications need at least one address in custom_emails (Azure retired the subscription-administrator email flags in April 2024)

### spec.notification.email.customEmails

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.notification.webhooks

`[]AzureMonitorAutoscaleSettingNotificationWebhook`

### spec.notification.webhooks[].serviceUri

`string` · required

- rule: the webhook service_uri must be an http:// or https:// URL
- rule: {"required":true}

### spec.notification.webhooks[].properties

`map<string, string>`

### spec.tags

`map<string, string>`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureMonitorAutoscaleSetting, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.autoscale_setting_id` | `string` |  |
| `status.outputs.autoscale_setting_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## See Also

- [Overview](../README.md)
