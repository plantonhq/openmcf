# AwsSagemakerFeatureGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerFeatureGroup example (hack/dev manifest and
# refgen Example source): a dual-store group exercising every arm -
# online store with InMemory storage, KMS, and TTL; offline store with
# Iceberg, KMS, and a named Glue catalog entry; a vector feature;
# provisioned throughput. Literal ARNs stand in for composed references
# so the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerFeatureGroup
metadata:
  name: customer-features
  id: customer-features
  org: test-org
  env: dev
spec:
  region: us-west-2
  recordIdentifierFeatureName: customer_id
  eventTimeFeatureName: event_time
  description: Customer features for churn scoring
  roleArn:
    value: arn:aws:iam::123456789012:role/sagemaker-execution
  featureDefinitions:
    - name: customer_id
      type: String
    - name: event_time
      type: Fractional
    - name: lifetime_value
      type: Fractional
    - name: embedding
      type: Fractional
      collectionType: Vector
      vectorDimension: 256
  onlineStore:
    enabled: true
    kmsKeyArn:
      value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
    storageType: InMemory
    ttl:
      unit: Days
      value: 30
  offlineStore:
    s3Uri: s3://my-feature-store/customers/
    kmsKeyArn:
      value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
    tableFormat: Iceberg
    disableGlueTableCreation: false
    dataCatalog:
      catalog: AwsDataCatalog
      database: feature_store
      tableName: customer_features
  throughput:
    mode: Provisioned
    provisionedReadCapacityUnits: 100
    provisionedWriteCapacityUnits: 50
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.recordIdentifierFeatureName` | `string` | yes |  |  |
| `spec.eventTimeFeatureName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.featureDefinitions` | `[]AwsSagemakerFeatureGroupFeature` | yes |  |  |
| `spec.featureDefinitions[].name` | `string` | yes |  |  |
| `spec.featureDefinitions[].type` | `string` |  |  |  |
| `spec.featureDefinitions[].collectionType` | `string` |  |  |  |
| `spec.featureDefinitions[].vectorDimension` | `int32` |  |  |  |
| `spec.onlineStore` | `AwsSagemakerFeatureGroupOnlineStore` |  |  |  |
| `spec.onlineStore.enabled` | `bool` |  |  |  |
| `spec.onlineStore.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.onlineStore.storageType` | `string` |  |  |  |
| `spec.onlineStore.ttl` | `AwsSagemakerFeatureGroupTtl` |  |  |  |
| `spec.onlineStore.ttl.unit` | `string` |  |  |  |
| `spec.onlineStore.ttl.value` | `int32` |  |  |  |
| `spec.offlineStore` | `AwsSagemakerFeatureGroupOfflineStore` |  |  |  |
| `spec.offlineStore.s3Uri` | `string` | yes |  |  |
| `spec.offlineStore.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.offlineStore.disableGlueTableCreation` | `bool` |  |  |  |
| `spec.offlineStore.tableFormat` | `string` |  |  |  |
| `spec.offlineStore.dataCatalog` | `AwsSagemakerFeatureGroupDataCatalog` |  |  |  |
| `spec.offlineStore.dataCatalog.catalog` | `string` | yes |  |  |
| `spec.offlineStore.dataCatalog.database` | `string` | yes |  |  |
| `spec.offlineStore.dataCatalog.tableName` | `string` | yes |  |  |
| `spec.throughput` | `AwsSagemakerFeatureGroupThroughput` |  |  |  |
| `spec.throughput.mode` | `string` |  |  |  |
| `spec.throughput.provisionedReadCapacityUnits` | `int32` |  |  |  |
| `spec.throughput.provisionedWriteCapacityUnits` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.recordIdentifierFeatureName

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"64","pattern":"^[0-9A-Za-z]([-_]*[0-9A-Za-z])*$"}}

### spec.eventTimeFeatureName

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"64","pattern":"^[0-9A-Za-z]([-_]*[0-9A-Za-z])*$"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"128"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.featureDefinitions

`[]AwsSagemakerFeatureGroupFeature` · required

- rule: {"repeated":{"minItems":"1","maxItems":"2500"}}
- rule: vector_dimension is required when collection_type is Vector and forbidden otherwise

### spec.featureDefinitions[].name

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"64","pattern":"^[0-9A-Za-z]([-_]*[0-9A-Za-z])*$","notIn":["is_deleted","write_time","api_invocation_time"]}}

### spec.featureDefinitions[].type

`string`

- rule: {"string":{"in":["Integral","Fractional","String"]}}

### spec.featureDefinitions[].collectionType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["List","Set","Vector"]}}

### spec.featureDefinitions[].vectorDimension

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":8192,"gte":1}}

### spec.onlineStore

`AwsSagemakerFeatureGroupOnlineStore`

### spec.onlineStore.enabled

`bool`

### spec.onlineStore.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.onlineStore.storageType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Standard","InMemory"]}}

### spec.onlineStore.ttl

`AwsSagemakerFeatureGroupTtl`

### spec.onlineStore.ttl.unit

`string`

- rule: {"string":{"in":["Seconds","Minutes","Hours","Days","Weeks"]}}

### spec.onlineStore.ttl.value

`int32`

- rule: {"int32":{"gte":1}}

### spec.offlineStore

`AwsSagemakerFeatureGroupOfflineStore`

### spec.offlineStore.s3Uri

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.offlineStore.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.offlineStore.disableGlueTableCreation

`bool`

### spec.offlineStore.tableFormat

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Glue","Iceberg"]}}

### spec.offlineStore.dataCatalog

`AwsSagemakerFeatureGroupDataCatalog`

### spec.offlineStore.dataCatalog.catalog

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.offlineStore.dataCatalog.database

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.offlineStore.dataCatalog.tableName

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.throughput

`AwsSagemakerFeatureGroupThroughput`

- rule: provisioned capacity units require mode Provisioned

### spec.throughput.mode

`string`

- rule: {"string":{"in":["OnDemand","Provisioned"]}}

### spec.throughput.provisionedReadCapacityUnits

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":10000000,"gte":0}}

### spec.throughput.provisionedWriteCapacityUnits

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":10000000,"gte":0}}

## Validation Rules

- `at_least_one_store`: at least one of online_store (enabled) and offline_store must be configured
- `feature_names_unique`: feature_definitions entries must have unique names
- `record_identifier_defined`: record_identifier_feature_name must be one of feature_definitions
- `event_time_defined`: event_time_feature_name must be one of feature_definitions

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerFeatureGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.feature_group_name` | `string` |  |
| `status.outputs.feature_group_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.onlineStore.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.offlineStore.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |

## See Also

- [Overview](../README.md)
