# AwsCloudTrail

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsCloudTrail example (hack/dev manifest and refgen Example
# source): a multi-region audit trail with log-file validation, advanced
# event selectors, Insights, and CloudWatch Logs mirroring. Literal
# ARNs/names stand in for composed references so the offline `tofu plan`
# renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsCloudTrail
metadata:
  name: org-audit-trail
  id: org-audit-trail
  org: test-org
  env: dev
spec:
  region: us-west-2
  s3BucketName:
    value: my-audit-trail-logs
  s3KeyPrefix: audit
  isMultiRegionTrail: true
  enableLogFileValidation: true
  kmsKeyId:
    value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
  snsTopicName:
    value: audit-trail-delivery
  cloudwatchLogs:
    logGroupArn:
      value: arn:aws:logs:us-west-2:123456789012:log-group:org-audit-trail
    roleArn:
      value: arn:aws:iam::123456789012:role/cloudtrail-to-cloudwatch
  advancedEventSelectors:
    - name: Management events
      fieldSelectors:
        - field: eventCategory
          equals: ["Management"]
    - name: S3 object writes in the data lake
      fieldSelectors:
        - field: eventCategory
          equals: ["Data"]
        - field: resources.type
          equals: ["AWS::S3::Object"]
        - field: resources.ARN
          startsWith: ["arn:aws:s3:::my-data-lake/"]
        - field: readOnly
          equals: ["false"]
  insightTypes:
    - ApiCallRateInsight
    - ApiErrorRateInsight
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.s3BucketName` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.s3KeyPrefix` | `string` |  |  |  |
| `spec.isMultiRegionTrail` | `bool` |  |  |  |
| `spec.includeGlobalServiceEvents` | `bool` |  |  |  |
| `spec.isOrganizationTrail` | `bool` |  |  |  |
| `spec.enableLogging` | `bool` |  |  |  |
| `spec.enableLogFileValidation` | `bool` |  |  |  |
| `spec.kmsKeyId` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.snsTopicName` | `string \| valueFrom` |  |  | AwsSnsTopic (`status.outputs.topic_name`) |
| `spec.cloudwatchLogs` | `AwsCloudTrailCloudwatchLogs` |  |  |  |
| `spec.cloudwatchLogs.logGroupArn` | `string \| valueFrom` | yes |  | AwsCloudwatchLogGroup (`status.outputs.log_group_arn`) |
| `spec.cloudwatchLogs.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.eventSelectors` | `[]AwsCloudTrailEventSelector` |  |  |  |
| `spec.eventSelectors[].readWriteType` | `string` |  |  |  |
| `spec.eventSelectors[].includeManagementEvents` | `bool` |  |  |  |
| `spec.eventSelectors[].excludeManagementEventSources` | `[]string` |  |  |  |
| `spec.eventSelectors[].dataResources` | `[]AwsCloudTrailDataResource` |  |  |  |
| `spec.eventSelectors[].dataResources[].type` | `string` |  |  |  |
| `spec.eventSelectors[].dataResources[].values` | `[]string` | yes |  |  |
| `spec.advancedEventSelectors` | `[]AwsCloudTrailAdvancedEventSelector` |  |  |  |
| `spec.advancedEventSelectors[].name` | `string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors` | `[]AwsCloudTrailFieldSelector` | yes |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].field` | `string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].equals` | `[]string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].notEquals` | `[]string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].startsWith` | `[]string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].notStartsWith` | `[]string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].endsWith` | `[]string` |  |  |  |
| `spec.advancedEventSelectors[].fieldSelectors[].notEndsWith` | `[]string` |  |  |  |
| `spec.insightTypes` | `[]string` |  |  |  |
| `spec.organizationDelegatedAdminAccountId` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.s3BucketName

`string | valueFrom` · required

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.s3KeyPrefix

`string`

- rule: {"string":{"maxLen":"2000"}}

### spec.isMultiRegionTrail

`bool`

### spec.includeGlobalServiceEvents

`bool` · optional (explicit presence)

### spec.isOrganizationTrail

`bool`

### spec.enableLogging

`bool` · optional (explicit presence)

### spec.enableLogFileValidation

`bool`

### spec.kmsKeyId

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.snsTopicName

`string | valueFrom`

- references: AwsSnsTopic (`status.outputs.topic_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_name}} -- a bare string does not parse

### spec.cloudwatchLogs

`AwsCloudTrailCloudwatchLogs`

### spec.cloudwatchLogs.logGroupArn

`string | valueFrom` · required

- references: AwsCloudwatchLogGroup (`status.outputs.log_group_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchLogGroup, name: <that resource's name>, fieldPath: status.outputs.log_group_arn}} -- a bare string does not parse

### spec.cloudwatchLogs.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.eventSelectors

`[]AwsCloudTrailEventSelector`

- rule: {"repeated":{"maxItems":"5"}}

### spec.eventSelectors[].readWriteType

`string`

- rule: {"string":{"in":["","ReadOnly","WriteOnly","All"]}}

### spec.eventSelectors[].includeManagementEvents

`bool` · optional (explicit presence)

### spec.eventSelectors[].excludeManagementEventSources

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["kms.amazonaws.com","rdsdata.amazonaws.com"]}}}}

### spec.eventSelectors[].dataResources

`[]AwsCloudTrailDataResource`

### spec.eventSelectors[].dataResources[].type

`string`

- rule: {"string":{"in":["AWS::DynamoDB::Table","AWS::Lambda::Function","AWS::S3::Object"]}}

### spec.eventSelectors[].dataResources[].values

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"250"}}

### spec.advancedEventSelectors

`[]AwsCloudTrailAdvancedEventSelector`

### spec.advancedEventSelectors[].name

`string`

- rule: {"string":{"maxLen":"1000"}}

### spec.advancedEventSelectors[].fieldSelectors

`[]AwsCloudTrailFieldSelector` · required

- rule: {"repeated":{"minItems":"1"}}
- rule: set at least one of equals, not_equals, starts_with, not_starts_with, ends_with, not_ends_with

### spec.advancedEventSelectors[].fieldSelectors[].field

`string`

- rule: {"string":{"in":["errorCode","eventCategory","eventName","eventSource","eventType","readOnly","resources.ARN","resources.type","sessionCredentialFromConsole","userIdentity.arn","vpcEndpointId"]}}

### spec.advancedEventSelectors[].fieldSelectors[].equals

`[]string`

### spec.advancedEventSelectors[].fieldSelectors[].notEquals

`[]string`

### spec.advancedEventSelectors[].fieldSelectors[].startsWith

`[]string`

### spec.advancedEventSelectors[].fieldSelectors[].notStartsWith

`[]string`

### spec.advancedEventSelectors[].fieldSelectors[].endsWith

`[]string`

### spec.advancedEventSelectors[].fieldSelectors[].notEndsWith

`[]string`

### spec.insightTypes

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["ApiCallRateInsight","ApiErrorRateInsight"]}}}}

### spec.organizationDelegatedAdminAccountId

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[0-9]{12}$"}}

## Validation Rules

- `one_selector_style`: event_selectors and advanced_event_selectors are mutually exclusive - pick one style
- `delegated_admin_requires_org_trail`: organization_delegated_admin_account_id requires is_organization_trail = true

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsCloudTrail, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.trail_arn` | `string` |  |
| `status.outputs.home_region` | `string` |  |
| `status.outputs.sns_topic_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.s3BucketName` | AwsS3Bucket | `status.outputs.bucket_id` |
| `spec.kmsKeyId` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.snsTopicName` | AwsSnsTopic | `status.outputs.topic_name` |
| `spec.cloudwatchLogs.logGroupArn` | AwsCloudwatchLogGroup | `status.outputs.log_group_arn` |
| `spec.cloudwatchLogs.roleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](../README.md)
