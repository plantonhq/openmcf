# AwsIamPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsIamPolicySpec defines a customer-managed IAM policy: a standalone,
versioned permission document that can be attached to many roles, users, and
groups at once.

A managed policy is the reusable unit of AWS permissions. Unlike an inline
policy (which lives and dies with a single role or user), a managed policy
has its own lifecycle and its own ARN, so one definition -- "read-only access
to the analytics bucket", "the permissions boundary for CI jobs" -- can be
referenced everywhere it is needed and updated in exactly one place. Roles
and users attach it by ARN through their managed_policy_arns fields, and a
permissions_boundary field also takes a policy ARN, which makes this kind a
leaf that much of an AWS architecture composes onto.

The policy name comes from metadata.name. AWS treats the name, path, and
description as create-only (changing any of them replaces the policy);
only the document itself is updatable, and AWS stores each update as a new
policy version.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamPolicy
metadata:
  name: s3-read-only-demo
spec:
  region: us-west-2
  description: "Demo managed policy granting read-only access to one bucket"
  path: "/service-policies/"
  # A multi-statement document with heterogeneous statement shapes (one carries
  # a Sid, one does not; different Action arities). policy_document is free-form
  # JSON (google.protobuf.Struct), so the fixture deliberately exercises the
  # shape that would fail if the TF variable were typed map(any) instead of any.
  policyDocument:
    Version: "2012-10-17"
    Statement:
      - Sid: "ListBucket"
        Effect: "Allow"
        Action: "s3:ListBucket"
        Resource: "arn:aws:s3:::demo-bucket"
      - Effect: "Allow"
        Action:
          - "s3:GetObject"
          - "s3:GetObjectVersion"
        Resource: "arn:aws:s3:::demo-bucket/*"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.policyDocument` | `object` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.path` | `string` |  | `/` |  |

## Field Details

### spec.region

`string` · required

The AWS region used by the provider while managing this policy.
IAM is a global service -- the policy is visible and attachable in every
region -- but every AWS API call is still made against a regional
endpoint, so a region is required.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.policyDocument

`object` · required

The permission document as free-form JSON (an IAM policy document with
Version and Statement). This is the only part of a managed policy that can
change after creation: each update creates a new policy version and marks
it default. AWS keeps at most 5 versions, so both IaC engines prune the
oldest non-default version before saving a new one -- updates keep working
indefinitely without manual version cleanup.
Example:
  Version: "2012-10-17"
  Statement:
    - Effect: Allow
      Action: ["s3:GetObject"]
      Resource: "arn:aws:s3:::my-bucket/*"

- rule: {"required":true}

### spec.description

`string`

An optional human-readable description of what the policy grants and why.
Shown in the IAM console next to the policy. Immutable: AWS has no
update-description API for managed policies, so changing it replaces the
policy (attachments are re-created by the IaC engine). Maximum 1000
characters.

- rule: {"string":{"maxLen":"1000"}}

### spec.path

`string`

The IAM path for the policy, used to organize and match policies in other
IAM policies (e.g. grant iam:AttachRolePolicy only for
"arn:aws:iam::<acct>:policy/service-boundaries/*"). Must begin and end
with "/" (e.g. "/service-boundaries/"). Defaults to "/" when omitted.
Immutable: changing the path replaces the policy.

- default: `/`

## Validation Rules

- `path_format`: path must begin and end with '/' and contain no empty segments, e.g. '/' or '/service-boundaries/'
- `path_length`: path must be at most 512 characters

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsIamPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.policyArn` | `string` | The ARN of the managed policy (e.g. "arn:aws:iam::123456789012:policy/s3-read-only"). The primary handle other resources reference via status.outputs.policy_arn -- attachments and permissions boundaries both take this value. |
| `status.outputs.policyId` | `string` | The stable unique ID AWS assigns to the policy (e.g. "ANPA..."). Unlike the ARN it never encodes the name or path, so it survives as a stable identifier in audit trails. |
| `status.outputs.policyName` | `string` | The friendly name of the policy (mirrors metadata.name), for building IAM console URLs and CLI commands. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
