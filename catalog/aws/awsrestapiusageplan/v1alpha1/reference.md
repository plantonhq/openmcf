# AwsRestApiUsagePlan

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsRestApiUsagePlan example (hack/dev manifest and refgen
# Example source): a plan covering one REST API stage with a quota, a
# throttle, a per-method throttle, and one API key. Literal ids stand
# in for composed references so the offline `tofu plan` renders every
# arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRestApiUsagePlan
metadata:
  name: orders-free-tier
  id: orders-free-tier
  org: test-org
  env: dev
spec:
  region: us-west-2
  description: Free tier -- 1000 requests/day
  apiStages:
    - restApiId:
        value: abcdef1234
      stageName:
        value: prod
      methodThrottles:
        - path: /search/GET
          burstLimit: 10
          rateLimit: 5
  quota:
    limit: 1000
    period: DAY
    offset: 0
  throttle:
    burstLimit: 100
    rateLimit: 50
  apiKeys:
    - name: acme-mobile
      description: acme-corp mobile app
      enabled: true
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.productCode` | `string` |  |  |  |
| `spec.apiStages` | `[]AwsRestApiUsagePlanApiStage` |  |  |  |
| `spec.apiStages[].restApiId` | `string \| valueFrom` | yes |  | AwsRestApiGateway (`status.outputs.rest_api_id`) |
| `spec.apiStages[].stageName` | `string \| valueFrom` | yes |  | AwsRestApiGateway (`status.outputs.stage_name`) |
| `spec.apiStages[].methodThrottles` | `[]AwsRestApiUsagePlanMethodThrottle` |  |  |  |
| `spec.apiStages[].methodThrottles[].path` | `string` | yes |  |  |
| `spec.apiStages[].methodThrottles[].burstLimit` | `int32` |  |  |  |
| `spec.apiStages[].methodThrottles[].rateLimit` | `double` |  |  |  |
| `spec.quota` | `AwsRestApiUsagePlanQuota` |  |  |  |
| `spec.quota.limit` | `int32` |  |  |  |
| `spec.quota.period` | `string` |  |  |  |
| `spec.quota.offset` | `int32` |  |  |  |
| `spec.throttle` | `AwsRestApiUsagePlanThrottle` |  |  |  |
| `spec.throttle.burstLimit` | `int32` |  |  |  |
| `spec.throttle.rateLimit` | `double` |  |  |  |
| `spec.apiKeys` | `[]AwsRestApiUsagePlanApiKey` |  |  |  |
| `spec.apiKeys[].name` | `string` | yes |  |  |
| `spec.apiKeys[].description` | `string` |  |  |  |
| `spec.apiKeys[].enabled` | `bool` |  |  |  |
| `spec.apiKeys[].customerId` | `string` |  |  |  |
| `spec.apiKeys[].value` | `string` (sensitive) | yes |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.productCode

`string`

### spec.apiStages

`[]AwsRestApiUsagePlanApiStage`

- rule: method_throttles entries must have unique path values

### spec.apiStages[].restApiId

`string | valueFrom` · required

- references: AwsRestApiGateway (`status.outputs.rest_api_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRestApiGateway, name: <that resource's name>, fieldPath: status.outputs.rest_api_id}} -- a bare string does not parse

### spec.apiStages[].stageName

`string | valueFrom` · required

- references: AwsRestApiGateway (`status.outputs.stage_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRestApiGateway, name: <that resource's name>, fieldPath: status.outputs.stage_name}} -- a bare string does not parse

### spec.apiStages[].methodThrottles

`[]AwsRestApiUsagePlanMethodThrottle`

### spec.apiStages[].methodThrottles[].path

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.apiStages[].methodThrottles[].burstLimit

`int32`

- rule: {"int32":{"gte":0}}

### spec.apiStages[].methodThrottles[].rateLimit

`double`

- rule: {"double":{"gte":0}}

### spec.quota

`AwsRestApiUsagePlanQuota`

- rule: quota offset must be 0 for DAY, 0-6 for WEEK, and 0-27 for MONTH

### spec.quota.limit

`int32`

- rule: {"int32":{"gte":1}}

### spec.quota.period

`string`

- rule: {"string":{"in":["DAY","WEEK","MONTH"]}}

### spec.quota.offset

`int32`

- rule: {"int32":{"gte":0}}

### spec.throttle

`AwsRestApiUsagePlanThrottle`

- rule: set at least one of burst_limit or rate_limit (omit throttle entirely for account defaults)

### spec.throttle.burstLimit

`int32`

- rule: {"int32":{"gte":0}}

### spec.throttle.rateLimit

`double`

- rule: {"double":{"gte":0}}

### spec.apiKeys

`[]AwsRestApiUsagePlanApiKey`

### spec.apiKeys[].name

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.apiKeys[].description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.apiKeys[].enabled

`bool` · optional (explicit presence)

### spec.apiKeys[].customerId

`string`

### spec.apiKeys[].value

`string` · required · sensitive

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"minLen":"20","maxLen":"128"}}

## Validation Rules

- `api_key_names_unique`: api_keys entries must have unique names

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRestApiUsagePlan, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.usage_plan_id` | `string` |  |
| `status.outputs.usage_plan_arn` | `string` |  |
| `status.outputs.api_key_ids` | `map<string, string>` |  |
| `status.outputs.api_key_arns` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.apiStages[].restApiId` | AwsRestApiGateway | `status.outputs.rest_api_id` |
| `spec.apiStages[].stageName` | AwsRestApiGateway | `status.outputs.stage_name` |

## See Also

- [Overview](../README.md)
