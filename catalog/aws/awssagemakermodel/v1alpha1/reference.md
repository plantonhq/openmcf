# AwsSagemakerModel

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerModel example (hack/dev manifest and refgen
# Example source): an inference pipeline exercising every arm - two
# containers with environments and hostnames, uncompressed S3 model
# data with EULA acceptance, an additional adapter channel, a private
# registry image config, multi-model caching, VPC attachment, and
# network isolation. Literal ARNs stand in for composed references so
# the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerModel
metadata:
  name: churn-scoring-pipeline
  id: churn-scoring-pipeline
  org: test-org
  env: dev
spec:
  region: us-west-2
  executionRoleArn:
    value: arn:aws:iam::123456789012:role/sagemaker-execution
  enableNetworkIsolation: true
  inferenceExecutionMode: Serial
  containers:
    - image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/churn-preprocess:2.1
      containerHostname: preprocess
      environment:
        LOG_LEVEL: info
      modelDataSource:
        s3Uri: s3://my-models/churn/preprocess/
        s3DataType: S3Prefix
        compressionType: None
      imageConfig:
        repositoryAccessMode: Vpc
        repositoryCredentialsProviderArn: arn:aws:lambda:us-west-2:123456789012:function:registry-creds
    - image: 123456789012.dkr.ecr.us-west-2.amazonaws.com/churn-score:2.1
      containerHostname: score
      mode: MultiModel
      multiModelCache: Disabled
      modelDataUrl: s3://my-models/churn/scoring/
      additionalModelDataSources:
        - channelName: adapters
          source:
            s3Uri: s3://my-models/churn/adapters/
            s3DataType: S3Prefix
            compressionType: None
            acceptEula: true
  vpcConfig:
    subnetIds:
      - value: subnet-0a1b2c3d4e5f60001
      - value: subnet-0a1b2c3d4e5f60002
    securityGroupIds:
      - value: sg-0a1b2c3d4e5f60001
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.executionRoleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.enableNetworkIsolation` | `bool` |  |  |  |
| `spec.primaryContainer` | `AwsSagemakerModelContainer` |  |  |  |
| `spec.primaryContainer.image` | `string` |  |  |  |
| `spec.primaryContainer.modelPackageArn` | `string` |  |  |  |
| `spec.primaryContainer.containerHostname` | `string` |  |  |  |
| `spec.primaryContainer.environment` | `map<string, string>` |  |  |  |
| `spec.primaryContainer.mode` | `string` |  |  |  |
| `spec.primaryContainer.modelDataUrl` | `string` |  |  |  |
| `spec.primaryContainer.modelDataSource` | `AwsSagemakerModelS3DataSource` |  |  |  |
| `spec.primaryContainer.modelDataSource.s3Uri` | `string` | yes |  |  |
| `spec.primaryContainer.modelDataSource.s3DataType` | `string` |  |  |  |
| `spec.primaryContainer.modelDataSource.compressionType` | `string` |  |  |  |
| `spec.primaryContainer.modelDataSource.acceptEula` | `bool` |  |  |  |
| `spec.primaryContainer.additionalModelDataSources` | `[]AwsSagemakerModelAdditionalDataSource` |  |  |  |
| `spec.primaryContainer.additionalModelDataSources[].channelName` | `string` | yes |  |  |
| `spec.primaryContainer.additionalModelDataSources[].source` | `AwsSagemakerModelS3DataSource` | yes |  |  |
| `spec.primaryContainer.additionalModelDataSources[].source.s3Uri` | `string` | yes |  |  |
| `spec.primaryContainer.additionalModelDataSources[].source.s3DataType` | `string` |  |  |  |
| `spec.primaryContainer.additionalModelDataSources[].source.compressionType` | `string` |  |  |  |
| `spec.primaryContainer.additionalModelDataSources[].source.acceptEula` | `bool` |  |  |  |
| `spec.primaryContainer.inferenceSpecificationName` | `string` |  |  |  |
| `spec.primaryContainer.multiModelCache` | `string` |  |  |  |
| `spec.primaryContainer.imageConfig` | `AwsSagemakerModelImageConfig` |  |  |  |
| `spec.primaryContainer.imageConfig.repositoryAccessMode` | `string` |  |  |  |
| `spec.primaryContainer.imageConfig.repositoryCredentialsProviderArn` | `string` |  |  |  |
| `spec.containers` | `[]AwsSagemakerModelContainer` |  |  |  |
| `spec.containers[].image` | `string` |  |  |  |
| `spec.containers[].modelPackageArn` | `string` |  |  |  |
| `spec.containers[].containerHostname` | `string` |  |  |  |
| `spec.containers[].environment` | `map<string, string>` |  |  |  |
| `spec.containers[].mode` | `string` |  |  |  |
| `spec.containers[].modelDataUrl` | `string` |  |  |  |
| `spec.containers[].modelDataSource` | `AwsSagemakerModelS3DataSource` |  |  |  |
| `spec.containers[].modelDataSource.s3Uri` | `string` | yes |  |  |
| `spec.containers[].modelDataSource.s3DataType` | `string` |  |  |  |
| `spec.containers[].modelDataSource.compressionType` | `string` |  |  |  |
| `spec.containers[].modelDataSource.acceptEula` | `bool` |  |  |  |
| `spec.containers[].additionalModelDataSources` | `[]AwsSagemakerModelAdditionalDataSource` |  |  |  |
| `spec.containers[].additionalModelDataSources[].channelName` | `string` | yes |  |  |
| `spec.containers[].additionalModelDataSources[].source` | `AwsSagemakerModelS3DataSource` | yes |  |  |
| `spec.containers[].additionalModelDataSources[].source.s3Uri` | `string` | yes |  |  |
| `spec.containers[].additionalModelDataSources[].source.s3DataType` | `string` |  |  |  |
| `spec.containers[].additionalModelDataSources[].source.compressionType` | `string` |  |  |  |
| `spec.containers[].additionalModelDataSources[].source.acceptEula` | `bool` |  |  |  |
| `spec.containers[].inferenceSpecificationName` | `string` |  |  |  |
| `spec.containers[].multiModelCache` | `string` |  |  |  |
| `spec.containers[].imageConfig` | `AwsSagemakerModelImageConfig` |  |  |  |
| `spec.containers[].imageConfig.repositoryAccessMode` | `string` |  |  |  |
| `spec.containers[].imageConfig.repositoryCredentialsProviderArn` | `string` |  |  |  |
| `spec.inferenceExecutionMode` | `string` |  |  |  |
| `spec.vpcConfig` | `AwsSagemakerModelVpcConfig` |  |  |  |
| `spec.vpcConfig.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.vpcConfig.securityGroupIds` | `[]string \| valueFrom` | yes |  | AwsSecurityGroup (`status.outputs.security_group_id`) |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.executionRoleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.enableNetworkIsolation

`bool`

### spec.primaryContainer

`AwsSagemakerModelContainer`

- rule: at least one of image and model_package_arn must be set
- rule: at most one of model_data_url and model_data_source may be set
- rule: environment keys must match ^[A-Za-z_][0-9A-Za-z_]*$ (max 1024 chars)
- rule: multi_model_cache requires mode MultiModel

### spec.primaryContainer.image

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"255","pattern":"^\\S+$"}}

### spec.primaryContainer.modelPackageArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.primaryContainer.containerHostname

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.primaryContainer.environment

`map<string, string>`

### spec.primaryContainer.mode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["SingleModel","MultiModel"]}}

### spec.primaryContainer.modelDataUrl

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.primaryContainer.modelDataSource

`AwsSagemakerModelS3DataSource`

### spec.primaryContainer.modelDataSource.s3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.primaryContainer.modelDataSource.s3DataType

`string`

- rule: {"string":{"in":["S3Prefix","S3Object"]}}

### spec.primaryContainer.modelDataSource.compressionType

`string`

- rule: {"string":{"in":["None","Gzip"]}}

### spec.primaryContainer.modelDataSource.acceptEula

`bool`

### spec.primaryContainer.additionalModelDataSources

`[]AwsSagemakerModelAdditionalDataSource`

### spec.primaryContainer.additionalModelDataSources[].channelName

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.primaryContainer.additionalModelDataSources[].source

`AwsSagemakerModelS3DataSource` · required

- rule: {"required":true}

### spec.primaryContainer.additionalModelDataSources[].source.s3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.primaryContainer.additionalModelDataSources[].source.s3DataType

`string`

- rule: {"string":{"in":["S3Prefix","S3Object"]}}

### spec.primaryContainer.additionalModelDataSources[].source.compressionType

`string`

- rule: {"string":{"in":["None","Gzip"]}}

### spec.primaryContainer.additionalModelDataSources[].source.acceptEula

`bool`

### spec.primaryContainer.inferenceSpecificationName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.primaryContainer.multiModelCache

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Enabled","Disabled"]}}

### spec.primaryContainer.imageConfig

`AwsSagemakerModelImageConfig`

- rule: repository_credentials_provider_arn requires repository_access_mode Vpc

### spec.primaryContainer.imageConfig.repositoryAccessMode

`string`

- rule: {"string":{"in":["Platform","Vpc"]}}

### spec.primaryContainer.imageConfig.repositoryCredentialsProviderArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.containers

`[]AwsSagemakerModelContainer`

- rule: {"repeated":{"maxItems":"15"}}
- rule: at least one of image and model_package_arn must be set
- rule: at most one of model_data_url and model_data_source may be set
- rule: environment keys must match ^[A-Za-z_][0-9A-Za-z_]*$ (max 1024 chars)
- rule: multi_model_cache requires mode MultiModel

### spec.containers[].image

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"255","pattern":"^\\S+$"}}

### spec.containers[].modelPackageArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.containers[].containerHostname

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.containers[].environment

`map<string, string>`

### spec.containers[].mode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["SingleModel","MultiModel"]}}

### spec.containers[].modelDataUrl

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.containers[].modelDataSource

`AwsSagemakerModelS3DataSource`

### spec.containers[].modelDataSource.s3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.containers[].modelDataSource.s3DataType

`string`

- rule: {"string":{"in":["S3Prefix","S3Object"]}}

### spec.containers[].modelDataSource.compressionType

`string`

- rule: {"string":{"in":["None","Gzip"]}}

### spec.containers[].modelDataSource.acceptEula

`bool`

### spec.containers[].additionalModelDataSources

`[]AwsSagemakerModelAdditionalDataSource`

### spec.containers[].additionalModelDataSources[].channelName

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.containers[].additionalModelDataSources[].source

`AwsSagemakerModelS3DataSource` · required

- rule: {"required":true}

### spec.containers[].additionalModelDataSources[].source.s3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.containers[].additionalModelDataSources[].source.s3DataType

`string`

- rule: {"string":{"in":["S3Prefix","S3Object"]}}

### spec.containers[].additionalModelDataSources[].source.compressionType

`string`

- rule: {"string":{"in":["None","Gzip"]}}

### spec.containers[].additionalModelDataSources[].source.acceptEula

`bool`

### spec.containers[].inferenceSpecificationName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.containers[].multiModelCache

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Enabled","Disabled"]}}

### spec.containers[].imageConfig

`AwsSagemakerModelImageConfig`

- rule: repository_credentials_provider_arn requires repository_access_mode Vpc

### spec.containers[].imageConfig.repositoryAccessMode

`string`

- rule: {"string":{"in":["Platform","Vpc"]}}

### spec.containers[].imageConfig.repositoryCredentialsProviderArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.inferenceExecutionMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Serial","Direct"]}}

### spec.vpcConfig

`AwsSagemakerModelVpcConfig`

### spec.vpcConfig.subnetIds

`[]string | valueFrom` · required

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"minItems":"1","maxItems":"16"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.vpcConfig.securityGroupIds

`[]string | valueFrom` · required

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"minItems":"1","maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

## Validation Rules

- `container_form_exactly_one`: exactly one of primary_container and containers must be set
- `execution_mode_requires_pipeline`: inference_execution_mode requires containers (an inference pipeline)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerModel, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.model_name` | `string` |  |
| `status.outputs.model_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.executionRoleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.vpcConfig.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.vpcConfig.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsSagemakerEndpoint | `spec.productionVariants[].model` | `status.outputs.model_name` |
| AwsSagemakerEndpoint | `spec.shadowVariants[].model` | `status.outputs.model_name` |

## See Also

- [Overview](../README.md)
