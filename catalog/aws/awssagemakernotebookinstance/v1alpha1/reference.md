# AwsSagemakerNotebookInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerNotebookInstance example (hack/dev manifest and
# refgen Example source): a VPC-confined GPU notebook exercising every
# arm - private networking, KMS, code repositories, IMDSv2, root
# lockdown, and both lifecycle scripts. Literal ARNs stand in for
# composed references so the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerNotebookInstance
metadata:
  name: research-notebook
  id: research-notebook
  org: test-org
  env: dev
spec:
  region: us-west-2
  instanceType: ml.g4dn.xlarge
  roleArn:
    value: arn:aws:iam::123456789012:role/sagemaker-execution
  volumeSizeGb: 100
  subnetId:
    value: subnet-0a1b2c3d4e5f60001
  securityGroupIds:
    - value: sg-0a1b2c3d4e5f60001
  kmsKeyArn:
    value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
  directInternetAccess: Disabled
  rootAccess: Disabled
  platformIdentifier: notebook-al2023-v1
  defaultCodeRepository: https://github.com/example/research-notebooks.git
  additionalCodeRepositories:
    - https://github.com/example/shared-utils.git
  imdsMinimumVersion: "2"
  lifecycleConfig:
    onCreate: |
      #!/bin/bash
      set -e
      pip install --quiet pandas scikit-learn
    onStart: |
      #!/bin/bash
      echo "notebook started at $(date)" >> /home/ec2-user/SageMaker/.start-log
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.instanceType` | `string` |  |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.volumeSizeGb` | `int32` |  |  |  |
| `spec.subnetId` | `string \| valueFrom` |  |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.directInternetAccess` | `string` |  |  |  |
| `spec.rootAccess` | `string` |  |  |  |
| `spec.platformIdentifier` | `string` |  |  |  |
| `spec.defaultCodeRepository` | `string` |  |  |  |
| `spec.additionalCodeRepositories` | `[]string` |  |  |  |
| `spec.imdsMinimumVersion` | `string` |  |  |  |
| `spec.lifecycleConfig` | `AwsSagemakerNotebookInstanceLifecycleConfig` |  |  |  |
| `spec.lifecycleConfig.onCreate` | `string` |  |  |  |
| `spec.lifecycleConfig.onStart` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.instanceType

`string`

- rule: {"string":{"pattern":"^ml\\.[a-z0-9]+([.-][a-z0-9]+)*$"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.volumeSizeGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":16384,"gte":5}}

### spec.subnetId

`string | valueFrom`

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.directInternetAccess

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Enabled","Disabled"]}}

### spec.rootAccess

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Enabled","Disabled"]}}

### spec.platformIdentifier

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["notebook-al1-v1","notebook-al2-v1","notebook-al2-v2","notebook-al2-v3","notebook-al2023-v1"]}}

### spec.defaultCodeRepository

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.additionalCodeRepositories

`[]string`

- rule: {"repeated":{"maxItems":"3","unique":true,"items":{"string":{"minLen":"1"}}}}

### spec.imdsMinimumVersion

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["1","2"]}}

### spec.lifecycleConfig

`AwsSagemakerNotebookInstanceLifecycleConfig`

- rule: at least one of on_create and on_start must be set

### spec.lifecycleConfig.onCreate

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.lifecycleConfig.onStart

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

## Validation Rules

- `disabled_internet_requires_vpc`: direct_internet_access Disabled requires subnet_id and security_group_ids
- `security_groups_require_subnet`: security_group_ids requires subnet_id

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerNotebookInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.notebook_instance_name` | `string` |  |
| `status.outputs.notebook_instance_arn` | `string` |  |
| `status.outputs.url` | `string` |  |
| `status.outputs.network_interface_id` | `string` |  |
| `status.outputs.lifecycle_config_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.subnetId` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |

## See Also

- [Overview](../README.md)
