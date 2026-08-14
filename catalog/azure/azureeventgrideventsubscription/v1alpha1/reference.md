# AzureEventgridEventSubscription

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Deep-shape example for docs and offline validation: a filtered
# webhook subscription on a custom topic, with delivery properties,
# tuned retry, and dead-lettering. References are literal values so
# the manifest validates standalone.
apiVersion: azure.planton.dev/v1alpha1
kind: AzureEventgridEventSubscription
metadata:
  name: test-eventgrid-event-subscription
  id: test-eventgrid-event-subscription
  org: test-org
  env: test
spec:
  scope:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.EventGrid/topics/orders-events
  name: big-orders-webhook
  destination:
    webhook:
      url: https://handler.example.com/events
      maxEventsPerBatch: 10
      preferredBatchSizeInKilobytes: 64
  deliveryProperties:
    - headerName: x-api-key
      type: Static
      value:
        value: test-api-key
      secret: true
    - headerName: x-event-subject
      type: Dynamic
      sourceField: subject
  deadLetter:
    storageAccountId:
      value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Storage/storageAccounts/deadletters
    storageBlobContainerName: eventgrid-dead-letters
  eventDeliverySchema: CloudEventSchemaV1_0
  includedEventTypes:
    - order.placed
    - order.cancelled
  subjectFilter:
    subjectBeginsWith: /orders/
    caseSensitive: false
  advancedFilter:
    numberGreaterThan:
      - key: data.amount
        value: 1000
    stringIn:
      - key: data.currency
        values:
          - USD
          - EUR
    numberInRange:
      - key: data.riskScore
        ranges:
          - from: 0
            to: 70
  advancedFilteringOnArraysEnabled: true
  labels:
    - orders
    - critical
  expirationTimeUtc: "2027-01-01T00:00:00Z"
  retryPolicy:
    maxDeliveryAttempts: 10
    eventTimeToLive: 240
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.scope` | `string \| valueFrom` |  |  | AzureEventgridTopic (`status.outputs.topic_id`) |
| `spec.systemTopicId` | `string \| valueFrom` |  |  | AzureEventgridSystemTopic (`status.outputs.system_topic_id`) |
| `spec.name` | `string` | yes |  |  |
| `spec.destination` | `AzureEventgridEventSubscriptionDestination` | yes |  |  |
| `spec.destination.azureFunction` | `AzureEventgridEventSubscriptionAzureFunctionDestination` |  |  |  |
| `spec.destination.azureFunction.functionId` | `string \| valueFrom` | yes |  |  |
| `spec.destination.azureFunction.maxEventsPerBatch` | `int32` |  |  |  |
| `spec.destination.azureFunction.preferredBatchSizeInKilobytes` | `int32` |  |  |  |
| `spec.destination.eventhubId` | `string \| valueFrom` |  |  | AzureEventHub (`status.outputs.event_hub_id`) |
| `spec.destination.hybridConnectionId` | `string \| valueFrom` |  |  |  |
| `spec.destination.serviceBusQueueId` | `string \| valueFrom` |  |  | AzureServiceBusQueue (`status.outputs.queue_id`) |
| `spec.destination.serviceBusTopicId` | `string \| valueFrom` |  |  | AzureServiceBusTopic (`status.outputs.topic_id`) |
| `spec.destination.storageQueue` | `AzureEventgridEventSubscriptionStorageQueueDestination` |  |  |  |
| `spec.destination.storageQueue.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.destination.storageQueue.queueName` | `string` | yes |  |  |
| `spec.destination.storageQueue.queueMessageTimeToLiveInSeconds` | `int32` |  |  |  |
| `spec.destination.webhook` | `AzureEventgridEventSubscriptionWebhookDestination` |  |  |  |
| `spec.destination.webhook.url` | `string` | yes |  |  |
| `spec.destination.webhook.maxEventsPerBatch` | `int32` |  |  |  |
| `spec.destination.webhook.preferredBatchSizeInKilobytes` | `int32` |  |  |  |
| `spec.destination.webhook.activeDirectoryTenantId` | `string` |  |  |  |
| `spec.destination.webhook.activeDirectoryAppIdOrUri` | `string` |  |  |  |
| `spec.deliveryIdentity` | `AzureEventgridEventSubscriptionIdentity` |  |  |  |
| `spec.deliveryIdentity.type` | `enum` | yes |  |  |
| `spec.deliveryIdentity.userAssignedIdentity` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.deliveryProperties` | `[]AzureEventgridEventSubscriptionDeliveryProperty` |  |  |  |
| `spec.deliveryProperties[].headerName` | `string` | yes |  |  |
| `spec.deliveryProperties[].type` | `string` | yes |  |  |
| `spec.deliveryProperties[].value` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.deliveryProperties[].sourceField` | `string` |  |  |  |
| `spec.deliveryProperties[].secret` | `bool` |  |  |  |
| `spec.deadLetter` | `AzureEventgridEventSubscriptionDeadLetter` |  |  |  |
| `spec.deadLetter.storageAccountId` | `string \| valueFrom` | yes |  | AzureStorageAccount (`status.outputs.storage_account_id`) |
| `spec.deadLetter.storageBlobContainerName` | `string` | yes |  |  |
| `spec.deadLetterIdentity` | `AzureEventgridEventSubscriptionIdentity` |  |  |  |
| `spec.deadLetterIdentity.type` | `enum` | yes |  |  |
| `spec.deadLetterIdentity.userAssignedIdentity` | `string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.eventDeliverySchema` | `string` |  | `EventGridSchema` |  |
| `spec.includedEventTypes` | `[]string` |  |  |  |
| `spec.subjectFilter` | `AzureEventgridEventSubscriptionSubjectFilter` |  |  |  |
| `spec.subjectFilter.subjectBeginsWith` | `string` |  |  |  |
| `spec.subjectFilter.subjectEndsWith` | `string` |  |  |  |
| `spec.subjectFilter.caseSensitive` | `bool` |  |  |  |
| `spec.advancedFilter` | `AzureEventgridEventSubscriptionAdvancedFilter` |  |  |  |
| `spec.advancedFilter.boolEquals` | `[]AzureEventgridEventSubscriptionBoolEqualsFilter` |  |  |  |
| `spec.advancedFilter.boolEquals[].key` | `string` | yes |  |  |
| `spec.advancedFilter.boolEquals[].value` | `bool` |  |  |  |
| `spec.advancedFilter.numberGreaterThan` | `[]AzureEventgridEventSubscriptionNumberFilter` |  |  |  |
| `spec.advancedFilter.numberGreaterThan[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberGreaterThan[].value` | `double` |  |  |  |
| `spec.advancedFilter.numberGreaterThanOrEquals` | `[]AzureEventgridEventSubscriptionNumberFilter` |  |  |  |
| `spec.advancedFilter.numberGreaterThanOrEquals[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberGreaterThanOrEquals[].value` | `double` |  |  |  |
| `spec.advancedFilter.numberLessThan` | `[]AzureEventgridEventSubscriptionNumberFilter` |  |  |  |
| `spec.advancedFilter.numberLessThan[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberLessThan[].value` | `double` |  |  |  |
| `spec.advancedFilter.numberLessThanOrEquals` | `[]AzureEventgridEventSubscriptionNumberFilter` |  |  |  |
| `spec.advancedFilter.numberLessThanOrEquals[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberLessThanOrEquals[].value` | `double` |  |  |  |
| `spec.advancedFilter.numberIn` | `[]AzureEventgridEventSubscriptionNumberListFilter` |  |  |  |
| `spec.advancedFilter.numberIn[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberIn[].values` | `[]double` | yes |  |  |
| `spec.advancedFilter.numberNotIn` | `[]AzureEventgridEventSubscriptionNumberListFilter` |  |  |  |
| `spec.advancedFilter.numberNotIn[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberNotIn[].values` | `[]double` | yes |  |  |
| `spec.advancedFilter.numberInRange` | `[]AzureEventgridEventSubscriptionNumberRangeFilter` |  |  |  |
| `spec.advancedFilter.numberInRange[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberInRange[].ranges` | `[]AzureEventgridEventSubscriptionNumberRange` | yes |  |  |
| `spec.advancedFilter.numberInRange[].ranges[].from` | `double` |  |  |  |
| `spec.advancedFilter.numberInRange[].ranges[].to` | `double` |  |  |  |
| `spec.advancedFilter.numberNotInRange` | `[]AzureEventgridEventSubscriptionNumberRangeFilter` |  |  |  |
| `spec.advancedFilter.numberNotInRange[].key` | `string` | yes |  |  |
| `spec.advancedFilter.numberNotInRange[].ranges` | `[]AzureEventgridEventSubscriptionNumberRange` | yes |  |  |
| `spec.advancedFilter.numberNotInRange[].ranges[].from` | `double` |  |  |  |
| `spec.advancedFilter.numberNotInRange[].ranges[].to` | `double` |  |  |  |
| `spec.advancedFilter.stringBeginsWith` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringBeginsWith[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringBeginsWith[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringNotBeginsWith` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringNotBeginsWith[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringNotBeginsWith[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringEndsWith` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringEndsWith[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringEndsWith[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringNotEndsWith` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringNotEndsWith[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringNotEndsWith[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringContains` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringContains[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringContains[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringNotContains` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringNotContains[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringNotContains[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringIn` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringIn[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringIn[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.stringNotIn` | `[]AzureEventgridEventSubscriptionStringListFilter` |  |  |  |
| `spec.advancedFilter.stringNotIn[].key` | `string` | yes |  |  |
| `spec.advancedFilter.stringNotIn[].values` | `[]string` | yes |  |  |
| `spec.advancedFilter.isNotNull` | `[]AzureEventgridEventSubscriptionKeyFilter` |  |  |  |
| `spec.advancedFilter.isNotNull[].key` | `string` | yes |  |  |
| `spec.advancedFilter.isNullOrUndefined` | `[]AzureEventgridEventSubscriptionKeyFilter` |  |  |  |
| `spec.advancedFilter.isNullOrUndefined[].key` | `string` | yes |  |  |
| `spec.advancedFilteringOnArraysEnabled` | `bool` |  | `false` |  |
| `spec.labels` | `[]string` |  |  |  |
| `spec.expirationTimeUtc` | `string` |  |  |  |
| `spec.retryPolicy` | `AzureEventgridEventSubscriptionRetryPolicy` |  |  |  |
| `spec.retryPolicy.maxDeliveryAttempts` | `int32` |  |  |  |
| `spec.retryPolicy.eventTimeToLive` | `int32` |  |  |  |

## Field Details

### spec.scope

`string | valueFrom`

- references: AzureEventgridTopic (`status.outputs.topic_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventgridTopic, name: <that resource's name>, fieldPath: status.outputs.topic_id}} -- a bare string does not parse

### spec.systemTopicId

`string | valueFrom`

- references: AzureEventgridSystemTopic (`status.outputs.system_topic_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventgridSystemTopic, name: <that resource's name>, fieldPath: status.outputs.system_topic_id}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Event subscription names must be 3-64 characters of letters, numbers, and hyphens
- rule: {"required":true}

### spec.destination

`AzureEventgridEventSubscriptionDestination` · required

- rule: {"required":true}
- rule: Set exactly one destination arm -- azure_function, eventhub_id, hybrid_connection_id, service_bus_queue_id, service_bus_topic_id, storage_queue, or webhook

### spec.destination.azureFunction

`AzureEventgridEventSubscriptionAzureFunctionDestination`

### spec.destination.azureFunction.functionId

`string | valueFrom` · required

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.destination.azureFunction.maxEventsPerBatch

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.destination.azureFunction.preferredBatchSizeInKilobytes

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.destination.eventhubId

`string | valueFrom`

- references: AzureEventHub (`status.outputs.event_hub_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureEventHub, name: <that resource's name>, fieldPath: status.outputs.event_hub_id}} -- a bare string does not parse

### spec.destination.hybridConnectionId

`string | valueFrom`

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.destination.serviceBusQueueId

`string | valueFrom`

- references: AzureServiceBusQueue (`status.outputs.queue_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureServiceBusQueue, name: <that resource's name>, fieldPath: status.outputs.queue_id}} -- a bare string does not parse

### spec.destination.serviceBusTopicId

`string | valueFrom`

- references: AzureServiceBusTopic (`status.outputs.topic_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureServiceBusTopic, name: <that resource's name>, fieldPath: status.outputs.topic_id}} -- a bare string does not parse

### spec.destination.storageQueue

`AzureEventgridEventSubscriptionStorageQueueDestination`

### spec.destination.storageQueue.storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.destination.storageQueue.queueName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.destination.storageQueue.queueMessageTimeToLiveInSeconds

`int32` · optional (explicit presence)

### spec.destination.webhook

`AzureEventgridEventSubscriptionWebhookDestination`

### spec.destination.webhook.url

`string` · required

- rule: The webhook url must be an https:// URL
- rule: {"required":true}

### spec.destination.webhook.maxEventsPerBatch

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":5000,"gte":1}}

### spec.destination.webhook.preferredBatchSizeInKilobytes

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":1024,"gte":1}}

### spec.destination.webhook.activeDirectoryTenantId

`string`

### spec.destination.webhook.activeDirectoryAppIdOrUri

`string`

### spec.deliveryIdentity

`AzureEventgridEventSubscriptionIdentity`

- rule: user_assigned_identity is required for USER_ASSIGNED and must be unset for SYSTEM_ASSIGNED

### spec.deliveryIdentity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_eventgrid_event_subscription_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`

### spec.deliveryIdentity.userAssignedIdentity

`string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.deliveryProperties

`[]AzureEventgridEventSubscriptionDeliveryProperty`

- rule: Static entries require value (and may set secret); Dynamic entries require source_field and must not set value or secret

### spec.deliveryProperties[].headerName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.deliveryProperties[].type

`string` · required

- rule: {"required":true,"string":{"in":["Static","Dynamic"]}}

### spec.deliveryProperties[].value

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.deliveryProperties[].sourceField

`string`

### spec.deliveryProperties[].secret

`bool`

### spec.deadLetter

`AzureEventgridEventSubscriptionDeadLetter`

### spec.deadLetter.storageAccountId

`string | valueFrom` · required

- references: AzureStorageAccount (`status.outputs.storage_account_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureStorageAccount, name: <that resource's name>, fieldPath: status.outputs.storage_account_id}} -- a bare string does not parse

### spec.deadLetter.storageBlobContainerName

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.deadLetterIdentity

`AzureEventgridEventSubscriptionIdentity`

- rule: user_assigned_identity is required for USER_ASSIGNED and must be unset for SYSTEM_ASSIGNED

### spec.deadLetterIdentity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_eventgrid_event_subscription_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`

### spec.deadLetterIdentity.userAssignedIdentity

`string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.eventDeliverySchema

`string` · optional (explicit presence)

- default: `EventGridSchema`
- rule: {"string":{"in":["EventGridSchema","CloudEventSchemaV1_0","CustomInputSchema"]}}

### spec.includedEventTypes

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.subjectFilter

`AzureEventgridEventSubscriptionSubjectFilter`

- rule: Set at least one of subject_begins_with, subject_ends_with, or case_sensitive

### spec.subjectFilter.subjectBeginsWith

`string`

### spec.subjectFilter.subjectEndsWith

`string`

### spec.subjectFilter.caseSensitive

`bool` · optional (explicit presence)

### spec.advancedFilter

`AzureEventgridEventSubscriptionAdvancedFilter`

- rule: An advanced_filter block must carry at least one condition

### spec.advancedFilter.boolEquals

`[]AzureEventgridEventSubscriptionBoolEqualsFilter`

### spec.advancedFilter.boolEquals[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.boolEquals[].value

`bool`

### spec.advancedFilter.numberGreaterThan

`[]AzureEventgridEventSubscriptionNumberFilter`

### spec.advancedFilter.numberGreaterThan[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberGreaterThan[].value

`double`

### spec.advancedFilter.numberGreaterThanOrEquals

`[]AzureEventgridEventSubscriptionNumberFilter`

### spec.advancedFilter.numberGreaterThanOrEquals[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberGreaterThanOrEquals[].value

`double`

### spec.advancedFilter.numberLessThan

`[]AzureEventgridEventSubscriptionNumberFilter`

### spec.advancedFilter.numberLessThan[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberLessThan[].value

`double`

### spec.advancedFilter.numberLessThanOrEquals

`[]AzureEventgridEventSubscriptionNumberFilter`

### spec.advancedFilter.numberLessThanOrEquals[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberLessThanOrEquals[].value

`double`

### spec.advancedFilter.numberIn

`[]AzureEventgridEventSubscriptionNumberListFilter`

### spec.advancedFilter.numberIn[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberIn[].values

`[]double` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25"}}

### spec.advancedFilter.numberNotIn

`[]AzureEventgridEventSubscriptionNumberListFilter`

### spec.advancedFilter.numberNotIn[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberNotIn[].values

`[]double` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25"}}

### spec.advancedFilter.numberInRange

`[]AzureEventgridEventSubscriptionNumberRangeFilter`

### spec.advancedFilter.numberInRange[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberInRange[].ranges

`[]AzureEventgridEventSubscriptionNumberRange` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25"}}

### spec.advancedFilter.numberInRange[].ranges[].from

`double`

### spec.advancedFilter.numberInRange[].ranges[].to

`double`

### spec.advancedFilter.numberNotInRange

`[]AzureEventgridEventSubscriptionNumberRangeFilter`

### spec.advancedFilter.numberNotInRange[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.numberNotInRange[].ranges

`[]AzureEventgridEventSubscriptionNumberRange` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25"}}

### spec.advancedFilter.numberNotInRange[].ranges[].from

`double`

### spec.advancedFilter.numberNotInRange[].ranges[].to

`double`

### spec.advancedFilter.stringBeginsWith

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringBeginsWith[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringBeginsWith[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringNotBeginsWith

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringNotBeginsWith[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringNotBeginsWith[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringEndsWith

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringEndsWith[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringEndsWith[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringNotEndsWith

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringNotEndsWith[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringNotEndsWith[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringContains

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringContains[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringContains[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringNotContains

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringNotContains[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringNotContains[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringIn

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringIn[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringIn[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.stringNotIn

`[]AzureEventgridEventSubscriptionStringListFilter`

### spec.advancedFilter.stringNotIn[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.stringNotIn[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"25","items":{"string":{"minLen":"1"}}}}

### spec.advancedFilter.isNotNull

`[]AzureEventgridEventSubscriptionKeyFilter`

### spec.advancedFilter.isNotNull[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilter.isNullOrUndefined

`[]AzureEventgridEventSubscriptionKeyFilter`

### spec.advancedFilter.isNullOrUndefined[].key

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.advancedFilteringOnArraysEnabled

`bool` · optional (explicit presence)

- default: `false`

### spec.labels

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.expirationTimeUtc

`string` · optional (explicit presence)

- rule: expiration_time_utc must be an RFC 3339 UTC timestamp like 2027-01-01T00:00:00Z

### spec.retryPolicy

`AzureEventgridEventSubscriptionRetryPolicy`

### spec.retryPolicy.maxDeliveryAttempts

`int32`

- rule: {"int32":{"lte":30,"gte":1}}

### spec.retryPolicy.eventTimeToLive

`int32`

- rule: {"int32":{"lte":1440,"gte":1}}

## Validation Rules

- `event_subscription_addressing_exactly_one`: Set exactly one of scope or system_topic_id -- the addressing choice determines the source
- `event_subscription_dead_letter_identity_requires_dead_letter`: dead_letter_identity requires dead_letter to be configured

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureEventgridEventSubscription, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.event_subscription_id` | `string` |  |
| `status.outputs.event_subscription_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.scope` | AzureEventgridTopic | `status.outputs.topic_id` |
| `spec.systemTopicId` | AzureEventgridSystemTopic | `status.outputs.system_topic_id` |
| `spec.destination.eventhubId` | AzureEventHub | `status.outputs.event_hub_id` |
| `spec.destination.serviceBusQueueId` | AzureServiceBusQueue | `status.outputs.queue_id` |
| `spec.destination.serviceBusTopicId` | AzureServiceBusTopic | `status.outputs.topic_id` |
| `spec.destination.storageQueue.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.deliveryIdentity.userAssignedIdentity` | AzureUserAssignedIdentity | `status.outputs.identity_id` |
| `spec.deadLetter.storageAccountId` | AzureStorageAccount | `status.outputs.storage_account_id` |
| `spec.deadLetterIdentity.userAssignedIdentity` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## See Also

- [Overview](../README.md)
