# AwsSagemakerMlflowServer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerMlflowServer example (hack/dev manifest and
# refgen Example source): a Medium server exercising every arm -
# version pin, auto-registration, and the maintenance window. Literal
# ARNs stand in for composed references so the offline `tofu plan`
# renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerMlflowServer
metadata:
  name: ml-experiments
  id: ml-experiments
  org: test-org
  env: dev
spec:
  region: us-west-2
  artifactStoreUri: s3://my-mlflow/artifacts
  roleArn:
    value: arn:aws:iam::123456789012:role/mlflow-server
  size: Medium
  mlflowVersion: "3.0"
  automaticModelRegistration: true
  weeklyMaintenanceWindowStart: TUE:03:30
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.artifactStoreUri` | `string` | yes |  |  |
| `spec.roleArn` | `string \| valueFrom` | yes |  | AwsIamRole (`status.outputs.role_arn`) |
| `spec.size` | `string` |  |  |  |
| `spec.mlflowVersion` | `string` |  |  |  |
| `spec.automaticModelRegistration` | `bool` |  |  |  |
| `spec.weeklyMaintenanceWindowStart` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.artifactStoreUri

`string` · required

- rule: {"string":{"minLen":"1","maxLen":"1024","pattern":"^(https|s3)://([^/]+)/?(.*)$"}}

### spec.roleArn

`string | valueFrom` · required

- references: AwsIamRole (`status.outputs.role_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamRole, name: <that resource's name>, fieldPath: status.outputs.role_arn}} -- a bare string does not parse

### spec.size

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["Small","Medium","Large"]}}

### spec.mlflowVersion

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^\\d+\\.\\d+$"}}

### spec.automaticModelRegistration

`bool`

### spec.weeklyMaintenanceWindowStart

`string`

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"pattern":"^(MON|TUE|WED|THU|FRI|SAT|SUN):([01]\\d|2[0-3]):[0-5]\\d$"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerMlflowServer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.tracking_server_name` | `string` |  |
| `status.outputs.tracking_server_arn` | `string` |  |
| `status.outputs.tracking_server_url` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.roleArn` | AwsIamRole | `status.outputs.role_arn` |

## See Also

- [Overview](../README.md)
