# AwsSagemakerMlflowApp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerMlflowApp example (hack/dev manifest and refgen
# Example source): an account-default serverless MLflow app exercising
# every arm - domain associations, auto-registration, and the
# maintenance window. Literal ARNs and domain IDs stand in for composed
# references so the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerMlflowApp
metadata:
  name: ml-experiments
  id: ml-experiments
  org: test-org
  env: dev
spec:
  region: us-west-2
  artifactStoreUri: s3://my-mlflow/artifacts
  roleArn:
    value: arn:aws:iam::123456789012:role/mlflow-app
  accountDefaultStatus: ENABLED
  defaultDomainIds:
    - value: d-abcdef123456
  modelRegistrationMode: AutoModelRegistrationEnabled
  weeklyMaintenanceWindowStart: SUN:03:00
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.artifactStoreUri` | `string` | yes |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.accountDefaultStatus` | `string` |  |  |  |
| `spec.defaultDomainIds` | `[]string \| valueFrom` |  |  | AwsSagemakerDomain (`status.outputs.domain_id`) |
| `spec.modelRegistrationMode` | `string` |  |  |  |
| `spec.weeklyMaintenanceWindowStart` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.artifactStoreUri

`string` · required

- rule: {"string":{"minLen":"1","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.accountDefaultStatus

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ENABLED","DISABLED"]}}

### spec.defaultDomainIds

`[]string | valueFrom`

- references: AwsSagemakerDomain (`status.outputs.domain_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSagemakerDomain, name: <that resource's name>, fieldPath: status.outputs.domain_id}} -- a bare string does not parse

### spec.modelRegistrationMode

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["AutoModelRegistrationEnabled","AutoModelRegistrationDisabled"]}}

### spec.weeklyMaintenanceWindowStart

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^(MON|TUE|WED|THU|FRI|SAT|SUN):([01]\\d|2[0-3]):[0-5]\\d$"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerMlflowApp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.app_arn` | `string` |  |
| `status.outputs.app_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |
| `spec.defaultDomainIds` | AwsSagemakerDomain | `status.outputs.domain_id` |

## See Also

- [Overview](../README.md)
