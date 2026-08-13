# AwsSagemakerImage

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerImage example (hack/dev manifest and refgen
# Example source): a registry entry with one fully-annotated version.
# Literal ECR paths stand in for real images so the offline `tofu plan`
# renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerImage
metadata:
  name: team-pytorch-kernels
  id: team-pytorch-kernels
  org: test-org
  env: dev
spec:
  region: us-west-2
  roleArn:
    value: arn:aws:iam::123456789012:role/sagemaker-execution
  displayName: PyTorch Kernels
  description: Team PyTorch kernel images for Studio
  versions:
    - baseImage: 123456789012.dkr.ecr.us-west-2.amazonaws.com/team-kernels:pytorch-2.4
      aliases:
        - latest
        - stable
      horovod: false
      jobType: NOTEBOOK_KERNEL
      mlFramework: PyTorch 2.4
      processor: GPU
      programmingLang: python 3.12
      releaseNotes: CUDA 12.4 base image
      vendorGuidance: STABLE
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.versions` | `[]AwsSagemakerImageVersion` |  |  |  |
| `spec.versions[].baseImage` | `string` | yes |  |  |
| `spec.versions[].aliases` | `[]string` |  |  |  |
| `spec.versions[].horovod` | `bool` |  |  |  |
| `spec.versions[].jobType` | `string` |  |  |  |
| `spec.versions[].mlFramework` | `string` |  |  |  |
| `spec.versions[].processor` | `string` |  |  |  |
| `spec.versions[].programmingLang` | `string` |  |  |  |
| `spec.versions[].releaseNotes` | `string` |  |  |  |
| `spec.versions[].vendorGuidance` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.displayName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"128"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"512"}}

### spec.versions

`[]AwsSagemakerImageVersion`

### spec.versions[].baseImage

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"255","pattern":"^\\S+$"}}

### spec.versions[].aliases

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"minLen":"1"}}}}

### spec.versions[].horovod

`bool`

### spec.versions[].jobType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["TRAINING","INFERENCE","NOTEBOOK_KERNEL"]}}

### spec.versions[].mlFramework

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[a-zA-Z]+ ?\\d+\\.\\d+(\\.\\d+)?$"}}

### spec.versions[].processor

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["CPU","GPU"]}}

### spec.versions[].programmingLang

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^[a-zA-Z]+ ?\\d+\\.\\d+(\\.\\d+)?$"}}

### spec.versions[].releaseNotes

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"255"}}

### spec.versions[].vendorGuidance

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["NOT_PROVIDED","STABLE","TO_BE_ARCHIVED","ARCHIVED"]}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerImage, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.image_name` | `string` |  |
| `status.outputs.image_arn` | `string` |  |
| `status.outputs.version_numbers` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](../README.md)
