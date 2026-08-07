# AliCloudRamPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudRamPolicySpec defines the configuration for an Alibaba Cloud
Resource Access Management (RAM) custom policy.

A RAM policy is a JSON document that defines a set of permissions (actions
on resources with optional conditions). Custom policies are created when the
built-in system policies do not provide the exact permission boundaries you
need. Once created, a policy can be attached to RAM roles, users, or groups
via their respective components.

Alibaba Cloud supports up to 5 versions per policy. When a policy is updated,
a new version is created and set as the default. The rotate_strategy field
controls what happens when the version limit is reached.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudRamPolicy
metadata:
  name: alicloudrampolicy-demo
spec:
  region: cn-hangzhou
  policyName: planton-demo-policy
  policyDocument: |
    {
      "Version": "1",
      "Statement": [
        {
          "Effect": "Allow",
          "Action": [
            "oss:GetObject",
            "oss:ListObjects"
          ],
          "Resource": [
            "acs:oss:*:*:planton-demo-bucket/*"
          ]
        }
      ]
    }
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.policyName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.policyDocument` | `string` | yes |  |  |
| `spec.rotateStrategy` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.force` | `bool` |  | `false` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider endpoint configuration.
RAM is a global service, but the provider requires a region for API routing.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.policyName

`string` · required

RAM policy name. Must be unique within the Alibaba Cloud account.
1-128 characters: English letters, digits, and hyphens (-).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.description

`string`

Human-readable description of what this policy allows or denies.
Maximum 1024 characters.

- rule: {"string":{"maxLen":"1024"}}

### spec.policyDocument

`string` · required

JSON IAM policy document defining the permissions.
Maximum 6144 bytes. Must be valid JSON conforming to the Alibaba Cloud
RAM policy structure with Version, Statement, Effect, Action, and Resource.

Example granting read-only access to a specific OSS bucket:
  {"Version":"1","Statement":[{"Effect":"Allow",
   "Action":["oss:GetObject","oss:ListObjects"],
   "Resource":["acs:oss:*:*:my-bucket/*"]}]}

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.rotateStrategy

`string` · optional (explicit presence)

Strategy for handling policy versions when the 5-version limit is reached.
"None" (default): update fails if the limit is reached; you must manually
  delete old versions.
"DeleteOldestNonDefaultVersionWhenLimitExceeded": automatically deletes
  the oldest non-default version to make room.

- rule: {"string":{"in":["None","DeleteOldestNonDefaultVersionWhenLimitExceeded"]}}

### spec.tags

`map<string, string>`

Tags applied to the RAM policy for organizational and cost-tracking purposes.

### spec.force

`bool` · optional (explicit presence)

Force-delete the policy even if it is still attached to roles, users, or groups.
When true, the policy is detached from all entities and all non-default versions
are deleted before the policy itself is removed.
When false (default), deletion fails if the policy is still attached.
Default: false

- default: `false`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudRamPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.policy_name` | `string` | The RAM policy name as created. Used by AliCloudRamRole policy_attachments to reference this policy. |
| `status.outputs.policy_type` | `string` | The policy type. Always "Custom" for user-created policies. Included as an output because AliCloudRamRole policy_attachments require both policy_name and policy_type for attachment. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
