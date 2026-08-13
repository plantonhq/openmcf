# AwsSagemakerPipeline

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerPipeline example (hack/dev manifest and refgen
# Example source): an S3-located definition with parallelism and
# display metadata. Literal ARNs stand in for composed references so
# the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerPipeline
metadata:
  name: churn-training
  id: churn-training
  org: test-org
  env: dev
spec:
  region: us-west-2
  displayName: churn-training-nightly
  description: Nightly churn-model retraining
  roleArn:
    value: arn:aws:iam::123456789012:role/sagemaker-pipeline-execution
  definitionS3Location:
    bucket:
      value: my-pipeline-definitions
    objectKey: pipelines/churn-training.json
    versionId: 3sL4kqtJlcpXroDTDmJ+rmSpXd3dIbrHY+MTRCxf3vjVBH40Nr8X8gdRQBpUMLUo
  parallelismMaxSteps: 4
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.definition` | `object` |  |  |  |
| `spec.definitionS3Location` | `AwsSagemakerPipelineDefinitionS3Location` |  |  |  |
| `spec.definitionS3Location.bucket` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.definitionS3Location.objectKey` | `string` | yes |  |  |
| `spec.definitionS3Location.versionId` | `string` |  |  |  |
| `spec.parallelismMaxSteps` | `int32` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.displayName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"256","pattern":"^[0-9A-Za-z]([0-9A-Za-z-])*$"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"3072"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.definition

`object`

### spec.definitionS3Location

`AwsSagemakerPipelineDefinitionS3Location`

### spec.definitionS3Location.bucket

`string | valueFrom` · required

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.definitionS3Location.objectKey

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.definitionS3Location.versionId

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.parallelismMaxSteps

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

## Validation Rules

- `definition_exactly_one`: exactly one of definition and definition_s3_location must be set

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerPipeline, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.pipeline_name` | `string` |  |
| `status.outputs.pipeline_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.definitionS3Location.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |

## See Also

- [Overview](../README.md)
