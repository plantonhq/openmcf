# AwsSagemakerEndpoint

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerEndpoint example (hack/dev manifest and refgen
# Example source): an instance-backed endpoint exercising every arm -
# weighted variants, managed instance scaling, routing, core dumps,
# data capture, async inference, KMS, and a canary blue/green
# deployment policy with alarm rollback. Literal ARNs stand in for
# composed references so the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerEndpoint
metadata:
  name: churn-scoring
  id: churn-scoring
  org: test-org
  env: dev
spec:
  region: us-west-2
  kmsKeyArn:
    value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
  productionVariants:
    - variantName: primary
      model:
        value: churn-scoring-model
      instanceType: ml.m5.large
      initialInstanceCount: 2
      initialVariantWeight: 0.9
      routingStrategy: LEAST_OUTSTANDING_REQUESTS
      volumeSizeGb: 64
      containerStartupHealthCheckTimeoutSeconds: 300
      modelDataDownloadTimeoutSeconds: 600
      enableSsmAccess: true
      managedInstanceScaling:
        status: ENABLED
        minInstanceCount: 1
        maxInstanceCount: 4
      coreDump:
        destinationS3Uri: s3://my-dumps/churn-scoring/
        kmsKeyArn:
          value: arn:aws:kms:us-west-2:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab
    - variantName: canary
      model:
        value: churn-scoring-model-next
      instanceType: ml.m5.large
      initialInstanceCount: 1
      initialVariantWeight: 0.1
  asyncInference:
    outputS3Path: s3://my-bucket/async-out/
    failureS3Path: s3://my-bucket/async-fail/
    maxConcurrentInvocationsPerInstance: 4
    successTopicArn:
      value: arn:aws:sns:us-west-2:123456789012:inference-ok
    errorTopicArn:
      value: arn:aws:sns:us-west-2:123456789012:inference-err
    includeInferenceResponseIn:
      - ERROR_NOTIFICATION_TOPIC
  dataCapture:
    destinationS3Uri: s3://my-bucket/capture/
    initialSamplingPercentage: 25
    captureModes:
      - Input
      - Output
    enableCapture: true
    jsonContentTypes:
      - application/json
  deployment:
    blueGreen:
      trafficRoutingType: CANARY
      waitIntervalSeconds: 300
      canarySize:
        type: CAPACITY_PERCENT
        value: 20
      terminationWaitSeconds: 120
      maximumExecutionTimeoutSeconds: 3600
    autoRollbackAlarmNames:
      - churn-scoring-5xx
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.productionVariants` | `[]AwsSagemakerEndpointVariant` | yes |  |  |
| `spec.productionVariants[].variantName` | `string` |  |  |  |
| `spec.productionVariants[].model` | `string \| valueFrom` |  |  | AwsSagemakerModel (`status.outputs.model_name`) |
| `spec.productionVariants[].instanceType` | `string` |  |  |  |
| `spec.productionVariants[].initialInstanceCount` | `int32` |  |  |  |
| `spec.productionVariants[].initialVariantWeight` | `float` |  |  |  |
| `spec.productionVariants[].serverless` | `AwsSagemakerEndpointServerlessConfig` |  |  |  |
| `spec.productionVariants[].serverless.maxConcurrency` | `int32` |  |  |  |
| `spec.productionVariants[].serverless.memorySizeMb` | `int32` |  |  |  |
| `spec.productionVariants[].serverless.provisionedConcurrency` | `int32` |  |  |  |
| `spec.productionVariants[].managedInstanceScaling` | `AwsSagemakerEndpointManagedInstanceScaling` |  |  |  |
| `spec.productionVariants[].managedInstanceScaling.status` | `string` |  |  |  |
| `spec.productionVariants[].managedInstanceScaling.minInstanceCount` | `int32` |  |  |  |
| `spec.productionVariants[].managedInstanceScaling.maxInstanceCount` | `int32` |  |  |  |
| `spec.productionVariants[].routingStrategy` | `string` |  |  |  |
| `spec.productionVariants[].volumeSizeGb` | `int32` |  |  |  |
| `spec.productionVariants[].containerStartupHealthCheckTimeoutSeconds` | `int32` |  |  |  |
| `spec.productionVariants[].modelDataDownloadTimeoutSeconds` | `int32` |  |  |  |
| `spec.productionVariants[].enableSsmAccess` | `bool` |  |  |  |
| `spec.productionVariants[].inferenceAmiVersion` | `string` |  |  |  |
| `spec.productionVariants[].acceleratorType` | `string` |  |  |  |
| `spec.productionVariants[].coreDump` | `AwsSagemakerEndpointCoreDump` |  |  |  |
| `spec.productionVariants[].coreDump.destinationS3Uri` | `string` | yes |  |  |
| `spec.productionVariants[].coreDump.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.productionVariants[].mlCapacityReservationArn` | `string` |  |  |  |
| `spec.shadowVariants` | `[]AwsSagemakerEndpointVariant` |  |  |  |
| `spec.shadowVariants[].variantName` | `string` |  |  |  |
| `spec.shadowVariants[].model` | `string \| valueFrom` |  |  | AwsSagemakerModel (`status.outputs.model_name`) |
| `spec.shadowVariants[].instanceType` | `string` |  |  |  |
| `spec.shadowVariants[].initialInstanceCount` | `int32` |  |  |  |
| `spec.shadowVariants[].initialVariantWeight` | `float` |  |  |  |
| `spec.shadowVariants[].serverless` | `AwsSagemakerEndpointServerlessConfig` |  |  |  |
| `spec.shadowVariants[].serverless.maxConcurrency` | `int32` |  |  |  |
| `spec.shadowVariants[].serverless.memorySizeMb` | `int32` |  |  |  |
| `spec.shadowVariants[].serverless.provisionedConcurrency` | `int32` |  |  |  |
| `spec.shadowVariants[].managedInstanceScaling` | `AwsSagemakerEndpointManagedInstanceScaling` |  |  |  |
| `spec.shadowVariants[].managedInstanceScaling.status` | `string` |  |  |  |
| `spec.shadowVariants[].managedInstanceScaling.minInstanceCount` | `int32` |  |  |  |
| `spec.shadowVariants[].managedInstanceScaling.maxInstanceCount` | `int32` |  |  |  |
| `spec.shadowVariants[].routingStrategy` | `string` |  |  |  |
| `spec.shadowVariants[].volumeSizeGb` | `int32` |  |  |  |
| `spec.shadowVariants[].containerStartupHealthCheckTimeoutSeconds` | `int32` |  |  |  |
| `spec.shadowVariants[].modelDataDownloadTimeoutSeconds` | `int32` |  |  |  |
| `spec.shadowVariants[].enableSsmAccess` | `bool` |  |  |  |
| `spec.shadowVariants[].inferenceAmiVersion` | `string` |  |  |  |
| `spec.shadowVariants[].acceleratorType` | `string` |  |  |  |
| `spec.shadowVariants[].coreDump` | `AwsSagemakerEndpointCoreDump` |  |  |  |
| `spec.shadowVariants[].coreDump.destinationS3Uri` | `string` | yes |  |  |
| `spec.shadowVariants[].coreDump.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.shadowVariants[].mlCapacityReservationArn` | `string` |  |  |  |
| `spec.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.executionRoleArn` | `string \| valueFrom` |  |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.asyncInference` | `AwsSagemakerEndpointAsyncInference` |  |  |  |
| `spec.asyncInference.outputS3Path` | `string` | yes |  |  |
| `spec.asyncInference.failureS3Path` | `string` |  |  |  |
| `spec.asyncInference.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.asyncInference.maxConcurrentInvocationsPerInstance` | `int32` |  |  |  |
| `spec.asyncInference.successTopicArn` | `string \| valueFrom` |  |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.asyncInference.errorTopicArn` | `string \| valueFrom` |  |  | AwsSnsTopic (`status.outputs.topic_arn`) |
| `spec.asyncInference.includeInferenceResponseIn` | `[]string` |  |  |  |
| `spec.dataCapture` | `AwsSagemakerEndpointDataCapture` |  |  |  |
| `spec.dataCapture.destinationS3Uri` | `string` | yes |  |  |
| `spec.dataCapture.initialSamplingPercentage` | `int32` |  |  |  |
| `spec.dataCapture.captureModes` | `[]string` | yes |  |  |
| `spec.dataCapture.enableCapture` | `bool` |  |  |  |
| `spec.dataCapture.csvContentTypes` | `[]string` |  |  |  |
| `spec.dataCapture.jsonContentTypes` | `[]string` |  |  |  |
| `spec.dataCapture.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.deployment` | `AwsSagemakerEndpointDeployment` |  |  |  |
| `spec.deployment.blueGreen` | `AwsSagemakerEndpointBlueGreenPolicy` |  |  |  |
| `spec.deployment.blueGreen.trafficRoutingType` | `string` |  |  |  |
| `spec.deployment.blueGreen.waitIntervalSeconds` | `int32` |  |  |  |
| `spec.deployment.blueGreen.canarySize` | `AwsSagemakerEndpointCapacitySize` |  |  |  |
| `spec.deployment.blueGreen.canarySize.type` | `string` |  |  |  |
| `spec.deployment.blueGreen.canarySize.value` | `int32` |  |  |  |
| `spec.deployment.blueGreen.linearStepSize` | `AwsSagemakerEndpointCapacitySize` |  |  |  |
| `spec.deployment.blueGreen.linearStepSize.type` | `string` |  |  |  |
| `spec.deployment.blueGreen.linearStepSize.value` | `int32` |  |  |  |
| `spec.deployment.blueGreen.terminationWaitSeconds` | `int32` |  |  |  |
| `spec.deployment.blueGreen.maximumExecutionTimeoutSeconds` | `int32` |  |  |  |
| `spec.deployment.rolling` | `AwsSagemakerEndpointRollingPolicy` |  |  |  |
| `spec.deployment.rolling.maximumBatchSize` | `AwsSagemakerEndpointCapacitySize` | yes |  |  |
| `spec.deployment.rolling.maximumBatchSize.type` | `string` |  |  |  |
| `spec.deployment.rolling.maximumBatchSize.value` | `int32` |  |  |  |
| `spec.deployment.rolling.waitIntervalSeconds` | `int32` |  |  |  |
| `spec.deployment.rolling.rollbackMaximumBatchSize` | `AwsSagemakerEndpointCapacitySize` |  |  |  |
| `spec.deployment.rolling.rollbackMaximumBatchSize.type` | `string` |  |  |  |
| `spec.deployment.rolling.rollbackMaximumBatchSize.value` | `int32` |  |  |  |
| `spec.deployment.rolling.maximumExecutionTimeoutSeconds` | `int32` |  |  |  |
| `spec.deployment.autoRollbackAlarmNames` | `[]string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.productionVariants

`[]AwsSagemakerEndpointVariant` · required

- rule: {"repeated":{"minItems":"1","maxItems":"10"}}
- rule: serverless variants cannot set instance_type, initial_instance_count, managed_instance_scaling, volume_size_gb, accelerator_type, inference_ami_version, core_dump, or ml_capacity_reservation_arn
- rule: either instance_type or serverless must be set

### spec.productionVariants[].variantName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.productionVariants[].model

`string | valueFrom`

- references: AwsSagemakerModel (`status.outputs.model_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSagemakerModel, name: <that resource's name>, fieldPath: status.outputs.model_name}} -- a bare string does not parse

### spec.productionVariants[].instanceType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^ml\\.[a-z0-9]+([.-][a-z0-9]+)*$"}}

### spec.productionVariants[].initialInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.productionVariants[].initialVariantWeight

`float` · optional (explicit presence)

- rule: {"float":{"gte":0}}

### spec.productionVariants[].serverless

`AwsSagemakerEndpointServerlessConfig`

- rule: provisioned_concurrency must not exceed max_concurrency

### spec.productionVariants[].serverless.maxConcurrency

`int32`

- rule: {"int32":{"lte":200,"gte":1}}

### spec.productionVariants[].serverless.memorySizeMb

`int32`

- rule: {"int32":{"in":[1024,2048,3072,4096,5120,6144]}}

### spec.productionVariants[].serverless.provisionedConcurrency

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":200,"gte":1}}

### spec.productionVariants[].managedInstanceScaling

`AwsSagemakerEndpointManagedInstanceScaling`

- rule: min_instance_count must not exceed max_instance_count

### spec.productionVariants[].managedInstanceScaling.status

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ENABLED","DISABLED"]}}

### spec.productionVariants[].managedInstanceScaling.minInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.productionVariants[].managedInstanceScaling.maxInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.productionVariants[].routingStrategy

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["LEAST_OUTSTANDING_REQUESTS","RANDOM"]}}

### spec.productionVariants[].volumeSizeGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":512,"gte":1}}

### spec.productionVariants[].containerStartupHealthCheckTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":60}}

### spec.productionVariants[].modelDataDownloadTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":60}}

### spec.productionVariants[].enableSsmAccess

`bool`

### spec.productionVariants[].inferenceAmiVersion

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["al2-ami-sagemaker-inference-gpu-2","al2-ami-sagemaker-inference-gpu-2-1","al2-ami-sagemaker-inference-gpu-3-1","al2-ami-sagemaker-inference-neuron-2","al2023-ami-sagemaker-inference-gpu-4-1"]}}

### spec.productionVariants[].acceleratorType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ml.eia1.medium","ml.eia1.large","ml.eia1.xlarge","ml.eia2.medium","ml.eia2.large","ml.eia2.xlarge"]}}

### spec.productionVariants[].coreDump

`AwsSagemakerEndpointCoreDump`

### spec.productionVariants[].coreDump.destinationS3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"512","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.productionVariants[].coreDump.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.productionVariants[].mlCapacityReservationArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.shadowVariants

`[]AwsSagemakerEndpointVariant`

- rule: {"repeated":{"maxItems":"10"}}
- rule: serverless variants cannot set instance_type, initial_instance_count, managed_instance_scaling, volume_size_gb, accelerator_type, inference_ami_version, core_dump, or ml_capacity_reservation_arn
- rule: either instance_type or serverless must be set

### spec.shadowVariants[].variantName

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"63","pattern":"^[0-9A-Za-z-]+$"}}

### spec.shadowVariants[].model

`string | valueFrom`

- references: AwsSagemakerModel (`status.outputs.model_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSagemakerModel, name: <that resource's name>, fieldPath: status.outputs.model_name}} -- a bare string does not parse

### spec.shadowVariants[].instanceType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^ml\\.[a-z0-9]+([.-][a-z0-9]+)*$"}}

### spec.shadowVariants[].initialInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.shadowVariants[].initialVariantWeight

`float` · optional (explicit presence)

- rule: {"float":{"gte":0}}

### spec.shadowVariants[].serverless

`AwsSagemakerEndpointServerlessConfig`

- rule: provisioned_concurrency must not exceed max_concurrency

### spec.shadowVariants[].serverless.maxConcurrency

`int32`

- rule: {"int32":{"lte":200,"gte":1}}

### spec.shadowVariants[].serverless.memorySizeMb

`int32`

- rule: {"int32":{"in":[1024,2048,3072,4096,5120,6144]}}

### spec.shadowVariants[].serverless.provisionedConcurrency

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":200,"gte":1}}

### spec.shadowVariants[].managedInstanceScaling

`AwsSagemakerEndpointManagedInstanceScaling`

- rule: min_instance_count must not exceed max_instance_count

### spec.shadowVariants[].managedInstanceScaling.status

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ENABLED","DISABLED"]}}

### spec.shadowVariants[].managedInstanceScaling.minInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":0}}

### spec.shadowVariants[].managedInstanceScaling.maxInstanceCount

`int32` · optional (explicit presence)

- rule: {"int32":{"gte":1}}

### spec.shadowVariants[].routingStrategy

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["LEAST_OUTSTANDING_REQUESTS","RANDOM"]}}

### spec.shadowVariants[].volumeSizeGb

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":512,"gte":1}}

### spec.shadowVariants[].containerStartupHealthCheckTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":60}}

### spec.shadowVariants[].modelDataDownloadTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":60}}

### spec.shadowVariants[].enableSsmAccess

`bool`

### spec.shadowVariants[].inferenceAmiVersion

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["al2-ami-sagemaker-inference-gpu-2","al2-ami-sagemaker-inference-gpu-2-1","al2-ami-sagemaker-inference-gpu-3-1","al2-ami-sagemaker-inference-neuron-2","al2023-ami-sagemaker-inference-gpu-4-1"]}}

### spec.shadowVariants[].acceleratorType

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ml.eia1.medium","ml.eia1.large","ml.eia1.xlarge","ml.eia2.medium","ml.eia2.large","ml.eia2.xlarge"]}}

### spec.shadowVariants[].coreDump

`AwsSagemakerEndpointCoreDump`

### spec.shadowVariants[].coreDump.destinationS3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"512","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.shadowVariants[].coreDump.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.shadowVariants[].mlCapacityReservationArn

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE"}

### spec.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.executionRoleArn

`string | valueFrom`

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.asyncInference

`AwsSagemakerEndpointAsyncInference`

### spec.asyncInference.outputS3Path

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"512","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.asyncInference.failureS3Path

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"512","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.asyncInference.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.asyncInference.maxConcurrentInvocationsPerInstance

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":1000,"gte":1}}

### spec.asyncInference.successTopicArn

`string | valueFrom`

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.asyncInference.errorTopicArn

`string | valueFrom`

- references: AwsSnsTopic (`status.outputs.topic_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSnsTopic, name: <that resource's name>, fieldPath: status.outputs.topic_arn}} -- a bare string does not parse

### spec.asyncInference.includeInferenceResponseIn

`[]string`

- rule: {"repeated":{"unique":true,"items":{"string":{"in":["SUCCESS_NOTIFICATION_TOPIC","ERROR_NOTIFICATION_TOPIC"]}}}}

### spec.dataCapture

`AwsSagemakerEndpointDataCapture`

### spec.dataCapture.destinationS3Uri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"512","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.dataCapture.initialSamplingPercentage

`int32`

- rule: {"int32":{"lte":100,"gte":0}}

### spec.dataCapture.captureModes

`[]string` · required

- rule: {"repeated":{"minItems":"1","maxItems":"2","unique":true,"items":{"string":{"in":["Input","Output","InputAndOutput"]}}}}

### spec.dataCapture.enableCapture

`bool`

### spec.dataCapture.csvContentTypes

`[]string`

- rule: {"repeated":{"maxItems":"10","items":{"string":{"minLen":"1","maxLen":"256","pattern":"^[0-9A-Za-z](-*[0-9A-Za-z])*\\/[0-9A-Za-z](-*[0-9A-Za-z.])*$"}}}}

### spec.dataCapture.jsonContentTypes

`[]string`

- rule: {"repeated":{"maxItems":"10","items":{"string":{"minLen":"1","maxLen":"256","pattern":"^[0-9A-Za-z](-*[0-9A-Za-z])*\\/[0-9A-Za-z](-*[0-9A-Za-z.])*$"}}}}

### spec.dataCapture.kmsKeyArn

`string | valueFrom`

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.deployment

`AwsSagemakerEndpointDeployment`

- rule: exactly one of blue_green and rolling must be set

### spec.deployment.blueGreen

`AwsSagemakerEndpointBlueGreenPolicy`

- rule: canary_size requires traffic_routing_type CANARY
- rule: linear_step_size requires traffic_routing_type LINEAR

### spec.deployment.blueGreen.trafficRoutingType

`string`

- rule: {"string":{"in":["ALL_AT_ONCE","CANARY","LINEAR"]}}

### spec.deployment.blueGreen.waitIntervalSeconds

`int32`

- rule: {"int32":{"lte":3600,"gte":0}}

### spec.deployment.blueGreen.canarySize

`AwsSagemakerEndpointCapacitySize`

### spec.deployment.blueGreen.canarySize.type

`string`

- rule: {"string":{"in":["INSTANCE_COUNT","CAPACITY_PERCENT"]}}

### spec.deployment.blueGreen.canarySize.value

`int32`

- rule: {"int32":{"gte":1}}

### spec.deployment.blueGreen.linearStepSize

`AwsSagemakerEndpointCapacitySize`

### spec.deployment.blueGreen.linearStepSize.type

`string`

- rule: {"string":{"in":["INSTANCE_COUNT","CAPACITY_PERCENT"]}}

### spec.deployment.blueGreen.linearStepSize.value

`int32`

- rule: {"int32":{"gte":1}}

### spec.deployment.blueGreen.terminationWaitSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":3600,"gte":0}}

### spec.deployment.blueGreen.maximumExecutionTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":14400,"gte":600}}

### spec.deployment.rolling

`AwsSagemakerEndpointRollingPolicy`

### spec.deployment.rolling.maximumBatchSize

`AwsSagemakerEndpointCapacitySize` · required

- rule: {"required":true}

### spec.deployment.rolling.maximumBatchSize.type

`string`

- rule: {"string":{"in":["INSTANCE_COUNT","CAPACITY_PERCENT"]}}

### spec.deployment.rolling.maximumBatchSize.value

`int32`

- rule: {"int32":{"gte":1}}

### spec.deployment.rolling.waitIntervalSeconds

`int32`

- rule: {"int32":{"lte":3600,"gte":0}}

### spec.deployment.rolling.rollbackMaximumBatchSize

`AwsSagemakerEndpointCapacitySize`

### spec.deployment.rolling.rollbackMaximumBatchSize.type

`string`

- rule: {"string":{"in":["INSTANCE_COUNT","CAPACITY_PERCENT"]}}

### spec.deployment.rolling.rollbackMaximumBatchSize.value

`int32`

- rule: {"int32":{"gte":1}}

### spec.deployment.rolling.maximumExecutionTimeoutSeconds

`int32` · optional (explicit presence)

- rule: {"int32":{"lte":14400,"gte":600}}

### spec.deployment.autoRollbackAlarmNames

`[]string`

- rule: {"repeated":{"maxItems":"10","unique":true,"items":{"string":{"minLen":"1"}}}}

## Validation Rules

- `variant_names_unique`: variant names must be unique across production_variants and shadow_variants
- `shadow_one_each_side`: shadow testing requires exactly one production variant and one shadow variant
- `role_required_without_models`: execution_role_arn is required when any variant omits model

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerEndpoint, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.endpoint_name` | `string` |  |
| `status.outputs.endpoint_arn` | `string` |  |
| `status.outputs.endpoint_config_name` | `string` |  |
| `status.outputs.endpoint_config_arn` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.productionVariants[].model` | AwsSagemakerModel | `status.outputs.model_name` |
| `spec.productionVariants[].coreDump.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.shadowVariants[].model` | AwsSagemakerModel | `status.outputs.model_name` |
| `spec.shadowVariants[].coreDump.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.executionRoleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.asyncInference.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.asyncInference.successTopicArn` | AwsSnsTopic | `status.outputs.topic_arn` |
| `spec.asyncInference.errorTopicArn` | AwsSnsTopic | `status.outputs.topic_arn` |
| `spec.dataCapture.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |

## See Also

- [Overview](../README.md)
