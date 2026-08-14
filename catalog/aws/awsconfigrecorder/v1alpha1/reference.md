# AwsConfigRecorder

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsConfigRecorder example (hack/dev manifest and refgen
# Example source): the region's recording posture -- a scoped inclusion
# recorder delivering to S3 with a retention window. Literal
# ARNs/names stand in for composed references so the offline `tofu
# plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsConfigRecorder
metadata:
  name: config-recording-us-west-2
  id: config-recording-us-west-2
  org: test-org
  env: dev
spec:
  region: us-west-2
  roleArn:
    value: arn:aws:iam::123456789012:role/config-recorder
  recordingGroup:
    allSupported: false
    resourceTypes:
      - AWS::S3::Bucket
      - AWS::EC2::SecurityGroup
    recordingStrategy: INCLUSION_BY_RESOURCE_TYPES
  recordingMode:
    recordingFrequency: CONTINUOUS
    override:
      description: Daily snapshots for noisy compute types
      recordingFrequency: DAILY
      resourceTypes:
        - AWS::EC2::Instance
  deliveryChannel:
    s3BucketName:
      value: my-config-history
    s3KeyPrefix: config
    snapshotDeliveryFrequency: TwentyFour_Hours
  retentionPeriodInDays: 365
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.recordingEnabled` | `bool` |  |  |  |
| `spec.recordingGroup` | `AwsConfigRecorderRecordingGroup` |  |  |  |
| `spec.recordingGroup.allSupported` | `bool` |  |  |  |
| `spec.recordingGroup.includeGlobalResourceTypes` | `bool` |  |  |  |
| `spec.recordingGroup.resourceTypes` | `[]string` |  |  |  |
| `spec.recordingGroup.exclusionByResourceTypes` | `[]string` |  |  |  |
| `spec.recordingGroup.recordingStrategy` | `string` |  |  |  |
| `spec.recordingMode` | `AwsConfigRecorderRecordingMode` |  |  |  |
| `spec.recordingMode.recordingFrequency` | `string` |  |  |  |
| `spec.recordingMode.override` | `AwsConfigRecorderRecordingModeOverride` |  |  |  |
| `spec.recordingMode.override.description` | `string` |  |  |  |
| `spec.recordingMode.override.recordingFrequency` | `string` |  |  |  |
| `spec.recordingMode.override.resourceTypes` | `[]string` | yes |  |  |
| `spec.deliveryChannel` | `AwsConfigRecorderDeliveryChannel` |  |  |  |
| `spec.deliveryChannel.s3BucketName` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.deliveryChannel.s3KeyPrefix` | `string` |  |  |  |
| `spec.deliveryChannel.s3KmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.deliveryChannel.snsTopicArn` | `string \| valueFrom` |  |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.deliveryChannel.snapshotDeliveryFrequency` | `string` |  |  |  |
| `spec.retentionPeriodInDays` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.recordingEnabled

`bool` · optional (explicit presence)

### spec.recordingGroup

`AwsConfigRecorderRecordingGroup`

- rule: recording_strategy ALL_SUPPORTED_RESOURCE_TYPES requires all_supported (and vice versa when a strategy is set)
- rule: exclusion_by_resource_types requires all_supported = false, an empty resource_types, and the EXCLUSION_BY_RESOURCE_TYPES strategy
- rule: resource_types requires all_supported = false, an empty exclusion list, and the INCLUSION_BY_RESOURCE_TYPES strategy

### spec.recordingGroup.allSupported

`bool` · optional (explicit presence)

### spec.recordingGroup.includeGlobalResourceTypes

`bool` · optional (explicit presence)

### spec.recordingGroup.resourceTypes

`[]string`

- rule: {"repeated":{"unique":true}}

### spec.recordingGroup.exclusionByResourceTypes

`[]string`

- rule: {"repeated":{"unique":true}}

### spec.recordingGroup.recordingStrategy

`string`

- rule: {"string":{"in":["","ALL_SUPPORTED_RESOURCE_TYPES","INCLUSION_BY_RESOURCE_TYPES","EXCLUSION_BY_RESOURCE_TYPES"]}}

### spec.recordingMode

`AwsConfigRecorderRecordingMode`

### spec.recordingMode.recordingFrequency

`string`

- rule: {"string":{"in":["","CONTINUOUS","DAILY"]}}

### spec.recordingMode.override

`AwsConfigRecorderRecordingModeOverride`

### spec.recordingMode.override.description

`string`

### spec.recordingMode.override.recordingFrequency

`string`

- rule: {"string":{"in":["CONTINUOUS","DAILY"]}}

### spec.recordingMode.override.resourceTypes

`[]string` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.deliveryChannel

`AwsConfigRecorderDeliveryChannel`

### spec.deliveryChannel.s3BucketName

`string | valueFrom` · required

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.deliveryChannel.s3KeyPrefix

`string`

### spec.deliveryChannel.s3KmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.deliveryChannel.snsTopicArn

`string | valueFrom`

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.deliveryChannel.snapshotDeliveryFrequency

`string`

- rule: {"string":{"in":["","One_Hour","Three_Hours","Six_Hours","Twelve_Hours","TwentyFour_Hours"]}}

### spec.retentionPeriodInDays

`int32`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":2557,"gte":30}}

## Validation Rules

- `running_recorder_needs_channel`: delivery_channel is required unless recording_enabled is explicitly false - AWS refuses to start a recorder without one

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsConfigRecorder, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.recorder_name` | `string` |  |
| `status.outputs.delivery_channel_name` | `string` |  |
| `status.outputs.recording_enabled` | `bool` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.deliveryChannel.s3BucketName` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.deliveryChannel.s3KmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.deliveryChannel.snsTopicArn` | AwsSnsTopic | `status.outputs.topic_arn` |

## See Also

- [Overview](../README.md)
