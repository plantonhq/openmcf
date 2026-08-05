# AliCloudRamRole

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudRamRoleSpec defines the configuration for an Alibaba Cloud
Resource Access Management (RAM) role with optional policy attachments.

A RAM role is a virtual identity that does not have permanent credentials.
Trusted entities (services, accounts, or federated identities) assume the
role via STS to obtain temporary security tokens. This is the standard
mechanism for granting Alibaba Cloud services (ECS, ACK, FC, SAE) the
permissions they need to operate on your behalf.

Policy attachments are bundled (per DD07) because a role without policies
grants zero permissions and is non-functional.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudRamRole
metadata:
  name: alicloudramrole-demo
spec:
  region: cn-hangzhou
  roleName: planton-demo-role
  assumeRolePolicyDocument: |
    {
      "Statement": [{
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {"Service": ["ecs.aliyuncs.com"]}
      }],
      "Version": "1"
    }
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.roleName` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.assumeRolePolicyDocument` | `string` | yes |  |  |
| `spec.maxSessionDuration` | `int32` |  | `3600` |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.force` | `bool` |  | `false` |  |
| `spec.policyAttachments` | `[]AliCloudRamRolePolicyAttachment` |  |  |  |
| `spec.policyAttachments[].policyName` | `string` | yes |  |  |
| `spec.policyAttachments[].policyType` | `string` |  | `System` |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider endpoint configuration.
RAM is a global service, but the provider requires a region for API routing.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.roleName

`string` · required

RAM role name. Must be unique within the Alibaba Cloud account.
1-64 characters: letters, digits, periods (.), hyphens (-), and underscores (_).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"64"}}

### spec.description

`string`

Human-readable description of the role's purpose.

### spec.assumeRolePolicyDocument

`string` · required

JSON trust policy document defining which principals can assume this role.
This controls who can call sts:AssumeRole to obtain temporary credentials.

Example for allowing ECS service to assume the role:
  {"Statement":[{"Action":"sts:AssumeRole","Effect":"Allow",
   "Principal":{"Service":["ecs.aliyuncs.com"]}}],"Version":"1"}

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.maxSessionDuration

`int32` · optional (explicit presence)

Maximum session duration in seconds when assuming this role via STS.
Longer durations are useful for CI/CD pipelines and batch workloads.
Range: 3600-43200 (1 hour to 12 hours).
Default: 3600

- default: `3600`
- rule: {"int32":{"lte":43200,"gte":3600}}

### spec.tags

`map<string, string>`

Tags applied to the RAM role for organizational and cost-tracking purposes.

### spec.force

`bool` · optional (explicit presence)

Force-detach all attached policies before deleting the role.
When false (default), deletion fails if policies are still attached.
Default: false

- default: `false`

### spec.policyAttachments

`[]AliCloudRamRolePolicyAttachment`

Policies to attach to this role. Each attachment grants the role a set
of permissions defined by either a system-managed or custom policy.

### spec.policyAttachments[].policyName

`string` · required

Name of the policy to attach.
System policies: "AliyunECSFullAccess", "AliyunOSSReadOnlyAccess", etc.
Custom policies: the name you gave when creating the policy.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.policyAttachments[].policyType

`string` · optional (explicit presence)

Policy type: "System" for Alibaba Cloud managed policies,
"Custom" for user-created policies.
Default: "System"

- default: `System`
- rule: {"string":{"in":["System","Custom"]}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudRamRole, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.role_id` | `string` | The RAM role ID assigned by Alibaba Cloud. |
| `status.outputs.role_name` | `string` | The RAM role name as created. |
| `status.outputs.arn` | `string` | The Alibaba Cloud Resource Name (ARN) for the role. Format: acs:ram::<account-id>:role/<role-name> Used by other components (ACK, FC, ECS) to reference this role. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudFunction | `spec.role` | `status.outputs.arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
