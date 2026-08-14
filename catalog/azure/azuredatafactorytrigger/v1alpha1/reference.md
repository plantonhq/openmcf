# AzureDataFactoryTrigger

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a schedule
# trigger exercising the full recurrence surface -- weekly frequency
# narrowed by the recurrence schedule (days, hours, minutes, monthly
# occurrences), start/end times, a time zone, and two pipelines with
# parameters. References are literal values so the manifest validates
# standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureDataFactoryTrigger
metadata:
  name: test-trigger
  id: test-trigger
  org: test-org
  env: test
spec:
  dataFactoryId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.DataFactory/factories/test-df
  name: weekly-loads
  description: Fires the weekly loads Monday and Friday at 02:00 and 14:30.
  annotations:
    - team:data
  activated: false
  schedule:
    frequency: Week
    interval: 1
    startTime: "2026-09-01T00:00:00Z"
    endTime: "2027-09-01T00:00:00Z"
    timeZone: UTC
    recurrenceSchedule:
      daysOfWeek:
        - Monday
        - Friday
      hours: [2, 14]
      minutes: [0, 30]
    pipelines:
      - name:
          value: ingest-weekly
        parameters:
          window: "@trigger().scheduledTime"
      - name:
          value: publish-weekly
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dataFactoryId` | `string \| valueFrom` | yes |  | AzureDataFactory (`status.outputs.data_factory_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.annotations` | `[]string` |  |  |  |
| `spec.activated` | `bool` |  | `true` |  |
| `spec.schedule` | `AzureDataFactoryTriggerSchedule` |  |  |  |
| `spec.schedule.frequency` | `string` |  | `Minute` |  |
| `spec.schedule.interval` | `int32` |  | `1` |  |
| `spec.schedule.startTime` | `string` |  |  |  |
| `spec.schedule.endTime` | `string` |  |  |  |
| `spec.schedule.timeZone` | `string` |  |  |  |
| `spec.schedule.recurrenceSchedule` | `AzureDataFactoryTriggerRecurrenceSchedule` |  |  |  |
| `spec.schedule.recurrenceSchedule.daysOfMonth` | `[]int32` |  |  |  |
| `spec.schedule.recurrenceSchedule.daysOfWeek` | `[]string` |  |  |  |
| `spec.schedule.recurrenceSchedule.hours` | `[]int32` |  |  |  |
| `spec.schedule.recurrenceSchedule.minutes` | `[]int32` |  |  |  |
| `spec.schedule.recurrenceSchedule.monthly` | `[]AzureDataFactoryTriggerMonthlyOccurrence` |  |  |  |
| `spec.schedule.recurrenceSchedule.monthly[].weekday` | `string` | yes |  |  |
| `spec.schedule.recurrenceSchedule.monthly[].week` | `int32` |  |  |  |
| `spec.schedule.pipelines` | `[]AzureDataFactoryTriggerPipelineReference` | yes |  |  |
| `spec.schedule.pipelines[].name` | `string \| valueFrom` | yes |  | AzureDataFactoryPipeline (`status.outputs.pipeline_name`) |
| `spec.schedule.pipelines[].parameters` | `map<string, string>` |  |  |  |
| `spec.tumblingWindow` | `AzureDataFactoryTriggerTumblingWindow` |  |  |  |
| `spec.tumblingWindow.frequency` | `string` | yes |  |  |
| `spec.tumblingWindow.interval` | `int32` |  |  |  |
| `spec.tumblingWindow.startTime` | `string` | yes |  |  |
| `spec.tumblingWindow.endTime` | `string` |  |  |  |
| `spec.tumblingWindow.delay` | `string` |  |  |  |
| `spec.tumblingWindow.maxConcurrency` | `int32` |  | `50` |  |
| `spec.tumblingWindow.retry` | `AzureDataFactoryTriggerRetryPolicy` |  |  |  |
| `spec.tumblingWindow.retry.count` | `int32` |  |  |  |
| `spec.tumblingWindow.retry.interval` | `int32` |  | `30` |  |
| `spec.tumblingWindow.dependencies` | `[]AzureDataFactoryTriggerDependency` |  |  |  |
| `spec.tumblingWindow.dependencies[].triggerName` | `string \| valueFrom` |  |  | AzureDataFactoryTrigger (`status.outputs.trigger_name`) |
| `spec.tumblingWindow.dependencies[].offset` | `string` |  |  |  |
| `spec.tumblingWindow.dependencies[].size` | `string` |  |  |  |
| `spec.tumblingWindow.additionalProperties` | `map<string, string>` |  |  |  |
| `spec.tumblingWindow.pipeline` | `AzureDataFactoryTriggerPipelineReference` | yes |  |  |
| `spec.tumblingWindow.pipeline.name` | `string \| valueFrom` | yes |  | AzureDataFactoryPipeline (`status.outputs.pipeline_name`) |
| `spec.tumblingWindow.pipeline.parameters` | `map<string, string>` |  |  |  |
| `spec.blobEvent` | `AzureDataFactoryTriggerBlobEvent` |  |  |  |
| `spec.blobEvent.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.blobEvent.events` | `[]string` | yes |  |  |
| `spec.blobEvent.blobPathBeginsWith` | `string` |  |  |  |
| `spec.blobEvent.blobPathEndsWith` | `string` |  |  |  |
| `spec.blobEvent.ignoreEmptyBlobs` | `bool` |  | `false` |  |
| `spec.blobEvent.additionalProperties` | `map<string, string>` |  |  |  |
| `spec.blobEvent.pipelines` | `[]AzureDataFactoryTriggerPipelineReference` | yes |  |  |
| `spec.blobEvent.pipelines[].name` | `string \| valueFrom` | yes |  | AzureDataFactoryPipeline (`status.outputs.pipeline_name`) |
| `spec.blobEvent.pipelines[].parameters` | `map<string, string>` |  |  |  |
| `spec.customEvent` | `AzureDataFactoryTriggerCustomEvent` |  |  |  |
| `spec.customEvent.eventgridTopicId` | `string \| valueFrom` | yes |  | AzureEventgridTopic (`status.outputs.topic_id`) |
| `spec.customEvent.events` | `[]string` | yes |  |  |
| `spec.customEvent.subjectBeginsWith` | `string` |  |  |  |
| `spec.customEvent.subjectEndsWith` | `string` |  |  |  |
| `spec.customEvent.additionalProperties` | `map<string, string>` |  |  |  |
| `spec.customEvent.pipelines` | `[]AzureDataFactoryTriggerPipelineReference` | yes |  |  |
| `spec.customEvent.pipelines[].name` | `string \| valueFrom` | yes |  | AzureDataFactoryPipeline (`status.outputs.pipeline_name`) |
| `spec.customEvent.pipelines[].parameters` | `map<string, string>` |  |  |  |

## Field Details

### spec.dataFactoryId

`string | valueFrom` · required

- references: AzureDataFactory (`status.outputs.data_factory_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactory, name: <that resource's name>, fieldPath: status.outputs.data_factory_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Trigger names must start with a letter, number, or underscore and must not contain < > * # . % & : \ + ? /
- rule: {"required":true}

### spec.description

`string`

### spec.annotations

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.activated

`bool` · optional (explicit presence)

- default: `true`

### spec.schedule

`AzureDataFactoryTriggerSchedule`

### spec.schedule.frequency

`string` · optional (explicit presence)

- default: `Minute`
- rule: {"string":{"in":["Minute","Hour","Day","Week","Month"]}}

### spec.schedule.interval

`int32` · optional (explicit presence)

- default: `1`
- rule: {"int32":{"gte":1}}

### spec.schedule.startTime

`string`

- rule: start_time must be an RFC 3339 timestamp, e.g. 2026-09-01T00:00:00Z

### spec.schedule.endTime

`string`

- rule: end_time must be an RFC 3339 timestamp, e.g. 2026-12-31T00:00:00Z

### spec.schedule.timeZone

`string`

### spec.schedule.recurrenceSchedule

`AzureDataFactoryTriggerRecurrenceSchedule`

### spec.schedule.recurrenceSchedule.daysOfMonth

`[]int32`

- rule: {"repeated":{"items":{"cel":[{"id":"data_factory_trigger_days_of_month_range","message":"days_of_month entries must be 1-31 (from the start) or -1..-31 (from the end)","expression":"this != 0 && this >= -31 && this <= 31"}]}}}

### spec.schedule.recurrenceSchedule.daysOfWeek

`[]string`

- rule: {"repeated":{"maxItems":"7","items":{"string":{"in":["Monday","Tuesday","Wednesday","Thursday","Friday","Saturday","Sunday"]}}}}

### spec.schedule.recurrenceSchedule.hours

`[]int32`

- rule: {"repeated":{"items":{"int32":{"lte":24,"gte":0}}}}

### spec.schedule.recurrenceSchedule.minutes

`[]int32`

- rule: {"repeated":{"items":{"int32":{"lte":60,"gte":0}}}}

### spec.schedule.recurrenceSchedule.monthly

`[]AzureDataFactoryTriggerMonthlyOccurrence`

### spec.schedule.recurrenceSchedule.monthly[].weekday

`string` · required

- rule: {"required":true,"string":{"in":["Monday","Tuesday","Wednesday","Thursday","Friday","Saturday","Sunday"]}}

### spec.schedule.recurrenceSchedule.monthly[].week

`int32` · optional (explicit presence)

- rule: week must be 1-5 (from the start) or -1..-5 (from the end)

### spec.schedule.pipelines

`[]AzureDataFactoryTriggerPipelineReference` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.schedule.pipelines[].name

`string | valueFrom` · required

- references: AzureDataFactoryPipeline (`status.outputs.pipeline_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryPipeline, name: <that resource's name>, fieldPath: status.outputs.pipeline_name}} -- a bare string does not parse

### spec.schedule.pipelines[].parameters

`map<string, string>`

### spec.tumblingWindow

`AzureDataFactoryTriggerTumblingWindow`

### spec.tumblingWindow.frequency

`string` · required

- rule: {"required":true,"string":{"in":["Minute","Hour","Month"]}}

### spec.tumblingWindow.interval

`int32`

- rule: {"int32":{"gte":1}}

### spec.tumblingWindow.startTime

`string` · required

- rule: start_time must be an RFC 3339 timestamp, e.g. 2026-09-01T00:00:00Z
- rule: {"required":true}

### spec.tumblingWindow.endTime

`string`

- rule: end_time must be an RFC 3339 timestamp, e.g. 2026-12-31T00:00:00Z

### spec.tumblingWindow.delay

`string`

- rule: delay must be a TimeSpan like 00:10:00 or 1.06:00:00

### spec.tumblingWindow.maxConcurrency

`int32` · optional (explicit presence)

- default: `50`
- rule: {"int32":{"lte":50,"gte":1}}

### spec.tumblingWindow.retry

`AzureDataFactoryTriggerRetryPolicy`

### spec.tumblingWindow.retry.count

`int32`

- rule: {"int32":{"gte":1}}

### spec.tumblingWindow.retry.interval

`int32` · optional (explicit presence)

- default: `30`
- rule: {"int32":{"gte":0}}

### spec.tumblingWindow.dependencies

`[]AzureDataFactoryTriggerDependency`

### spec.tumblingWindow.dependencies[].triggerName

`string | valueFrom`

- references: AzureDataFactoryTrigger (`status.outputs.trigger_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryTrigger, name: <that resource's name>, fieldPath: status.outputs.trigger_name}} -- a bare string does not parse

### spec.tumblingWindow.dependencies[].offset

`string`

- rule: offset must be a TimeSpan like 24:00:00 or -24:00:00

### spec.tumblingWindow.dependencies[].size

`string`

- rule: size must be a TimeSpan like 06:00:00

### spec.tumblingWindow.additionalProperties

`map<string, string>`

### spec.tumblingWindow.pipeline

`AzureDataFactoryTriggerPipelineReference` · required

- rule: {"required":true}

### spec.tumblingWindow.pipeline.name

`string | valueFrom` · required

- references: AzureDataFactoryPipeline (`status.outputs.pipeline_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryPipeline, name: <that resource's name>, fieldPath: status.outputs.pipeline_name}} -- a bare string does not parse

### spec.tumblingWindow.pipeline.parameters

`map<string, string>`

### spec.blobEvent

`AzureDataFactoryTriggerBlobEvent`

- rule: Set blob_path_begins_with, blob_path_ends_with, or both

### spec.blobEvent.storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.blobEvent.events

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"in":["Microsoft.Storage.BlobCreated","Microsoft.Storage.BlobDeleted"]}}}}

### spec.blobEvent.blobPathBeginsWith

`string`

### spec.blobEvent.blobPathEndsWith

`string`

### spec.blobEvent.ignoreEmptyBlobs

`bool` · optional (explicit presence)

- default: `false`

### spec.blobEvent.additionalProperties

`map<string, string>`

### spec.blobEvent.pipelines

`[]AzureDataFactoryTriggerPipelineReference` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.blobEvent.pipelines[].name

`string | valueFrom` · required

- references: AzureDataFactoryPipeline (`status.outputs.pipeline_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryPipeline, name: <that resource's name>, fieldPath: status.outputs.pipeline_name}} -- a bare string does not parse

### spec.blobEvent.pipelines[].parameters

`map<string, string>`

### spec.customEvent

`AzureDataFactoryTriggerCustomEvent`

### spec.customEvent.eventgridTopicId

`string | valueFrom` · required

- references: AzureEventgridTopic (`status.outputs.topic_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventgridTopic, name: <that resource's name>, fieldPath: status.outputs.topic_id}} -- a bare string does not parse

### spec.customEvent.events

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.customEvent.subjectBeginsWith

`string`

### spec.customEvent.subjectEndsWith

`string`

### spec.customEvent.additionalProperties

`map<string, string>`

### spec.customEvent.pipelines

`[]AzureDataFactoryTriggerPipelineReference` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.customEvent.pipelines[].name

`string | valueFrom` · required

- references: AzureDataFactoryPipeline (`status.outputs.pipeline_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureDataFactoryPipeline, name: <that resource's name>, fieldPath: status.outputs.pipeline_name}} -- a bare string does not parse

### spec.customEvent.pipelines[].parameters

`map<string, string>`

## Validation Rules

- `azure_data_factory_trigger_exactly_one_variant`: Set exactly one trigger variant -- schedule, tumbling_window, blob_event, or custom_event -- the variant determines the trigger type

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureDataFactoryTrigger, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.trigger_id` | `string` |  |
| `status.outputs.trigger_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.dataFactoryId` | AzureDataFactory | `status.outputs.data_factory_id` |
| `spec.schedule.pipelines[].name` | AzureDataFactoryPipeline | `status.outputs.pipeline_name` |
| `spec.tumblingWindow.dependencies[].triggerName` | AzureDataFactoryTrigger | `status.outputs.trigger_name` |
| `spec.tumblingWindow.pipeline.name` | AzureDataFactoryPipeline | `status.outputs.pipeline_name` |
| `spec.blobEvent.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.blobEvent.pipelines[].name` | AzureDataFactoryPipeline | `status.outputs.pipeline_name` |
| `spec.customEvent.eventgridTopicId` | AzureEventgridTopic | `status.outputs.topic_id` |
| `spec.customEvent.pipelines[].name` | AzureDataFactoryPipeline | `status.outputs.pipeline_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureDataFactoryTrigger | `spec.tumblingWindow.dependencies[].triggerName` | `status.outputs.trigger_name` |

## See Also

- [Overview](../README.md)
