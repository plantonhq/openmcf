# AwsSagemakerModelRegistry

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
# Canonical AwsSagemakerModelRegistry example (hack/dev manifest and
# refgen Example source): a described group with a cross-account
# sharing policy. Literal account IDs stand in for real principals so
# the offline `tofu plan` renders every arm.
apiVersion: aws.planton.dev/v1alpha1
kind: AwsSagemakerModelRegistry
metadata:
  name: churn-models
  id: churn-models
  org: test-org
  env: dev
spec:
  region: us-west-2
  description: Registered churn-scoring model versions
  resourcePolicy:
    Version: "2012-10-17"
    Statement:
      - Sid: AllowCrossAccountDescribe
        Effect: Allow
        Principal:
          AWS: arn:aws:iam::210987654321:root
        Action:
          - sagemaker:DescribeModelPackage
          - sagemaker:ListModelPackages
        Resource: "*"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.resourcePolicy` | `object` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

- rule: {"string":{"maxLen":"1024"}}

### spec.resourcePolicy

`object`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsSagemakerModelRegistry, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.model_package_group_name` | `string` |  |
| `status.outputs.model_package_group_arn` | `string` |  |

## See Also

- [Overview](../README.md)
