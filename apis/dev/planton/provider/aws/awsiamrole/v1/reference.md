# AwsIamRole

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsIamRoleSpec defines an IAM role: an assumable identity with temporary
credentials, the backbone of every service-to-service permission on AWS.

A role is two documents and a set of attachments. The trust policy answers
"who may assume this role" (a service principal like lambda.amazonaws.com,
another account, a federated identity). The permission side answers "what
can the role do once assumed": reusable managed policies attached by ARN
(see AwsIamPolicy), plus inline policies for permissions unique to this one
role. An optional permissions boundary caps the maximum permissions the
role can ever have, regardless of what its policies grant.

The role name comes from metadata.name. Name and path are create-only
(changing them replaces the role); the trust policy, description, session
duration, boundary, and policy attachments are all updatable in place.

Note that a role is assumed directly by most AWS services -- only EC2 needs
the AwsIamInstanceProfile wrapper to deliver a role to instances.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamRole
metadata:
  name: lambda-execution-role-demo
spec:
  region: us-west-2
  description: "Demo IAM role for Lambda function execution"
  path: "/service-roles/"
  trustPolicy:
    Version: "2012-10-17"
    Statement:
      - Effect: "Allow"
        Principal:
          Service: "lambda.amazonaws.com"
        Action: "sts:AssumeRole"
  managedPolicyArns:
    # Literal ARN (AWS-managed policy). Planton-defined policies attach via a
    # valueFrom reference to an AwsIamPolicy's policy_arn output.
    - value: "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  # Two intentionally differently-shaped inline policies. inline_policies is a free-form JSON map
  # (map<string, google.protobuf.Struct>); these two entries do not share a single object type
  # (different statement counts, one carries a Sid and the other does not), which is exactly the
  # shape that fails when the variable is typed map(any) instead of any. Keep BOTH so the fixture
  # reproduces the heterogeneous-map case a single-policy fixture would miss.
  inlinePolicies:
    extraLoggingPermissions:
      Version: "2012-10-17"
      Statement:
        - Sid: "CreateCloudWatchGroups"
          Effect: "Allow"
          Action:
            - "logs:CreateLogGroup"
          Resource: "*"
    customS3Access:
      Version: "2012-10-17"
      Statement:
        - Effect: "Allow"
          Action: "s3:ListBucket"
          Resource: "arn:aws:s3:::demo-bucket"
        - Effect: "Allow"
          Action:
            - "s3:GetObject"
            - "s3:PutObject"
          Resource: "arn:aws:s3:::demo-bucket/*"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.path` | `string` |  | `/` |  |
| `spec.trustPolicy` | `object` | yes |  |  |
| `spec.managedPolicyArns` | `[]string \| valueFrom` |  |  | AwsIamPolicy (`status.outputs.policy_arn`) |
| `spec.inlinePolicies` | `map<string, object>` |  |  |  |
| `spec.maxSessionDuration` | `int32` |  | `3600` |  |
| `spec.permissionsBoundary` | `string \| valueFrom` |  |  | AwsIamPolicy (`status.outputs.policy_arn`) |
| `spec.forceDetachPolicies` | `bool` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region used by the provider while managing this role.
IAM is a global service -- the role is assumable in every region -- but
every AWS API call is still made against a regional endpoint, so a region
is required.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.description

`string`

An optional human-readable description of the role's purpose, shown in
the IAM console. Updatable in place. Maximum 1000 characters; AWS rejects
typographic ("curly") quotes here, so stick to plain ASCII quoting.

- rule: {"string":{"maxLen":"1000"}}

### spec.path

`string`

The IAM path for the role, used to organize and match roles in IAM
policies (e.g. grant iam:PassRole only for
"arn:aws:iam::<acct>:role/service-roles/*"). Must begin and end with "/"
(e.g. "/service-roles/"). Defaults to "/" when omitted. Immutable:
changing the path replaces the role.

- default: `/`

### spec.trustPolicy

`object` · required

The trust policy as free-form JSON: the statement of WHO may assume this
role (service principals, AWS accounts, federated identities) and under
what conditions. This is the security-critical half of the role -- prefer
exact principals and add conditions (aws:SourceAccount, aws:SourceArn,
sts:ExternalId) to prevent confused-deputy access. Updatable in place.
Example:
  Version: "2012-10-17"
  Statement:
    - Effect: Allow
      Principal: { Service: lambda.amazonaws.com }
      Action: sts:AssumeRole

- rule: {"required":true}

### spec.managedPolicyArns

`[]string | valueFrom`

Managed policies to attach, each a reference to an AwsIamPolicy's
policy_arn output or a literal ARN (literals are how AWS-managed policies
like arn:aws:iam::aws:policy/ReadOnlyAccess attach). Attachments are
reconciled in place: adding or removing an entry attaches or detaches
without touching the role. Permissions unique to this role belong in
inline_policies instead.

- references: AwsIamPolicy (`status.outputs.policy_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamPolicy, name: <that resource's name>, fieldPath: status.outputs.policy_arn}} -- a bare string does not parse

### spec.inlinePolicies

`map<string, object>`

Inline policies embedded in this role: a map of policy name to a
free-form JSON permission document. An inline policy lives and dies with
the role, so use it for permissions that make no sense anywhere else
(e.g. access to this service's own queue); anything reused across
principals belongs in a first-class AwsIamPolicy attached via
managed_policy_arns.

- rule: {"map":{"keys":{"string":{"maxLen":"128"}}}}

### spec.maxSessionDuration

`int32`

The maximum duration, in seconds, of sessions assumed on this role
(the ceiling for the AssumeRole DurationSeconds parameter). Between 3600
(1 hour, the AWS default when unset) and 43200 (12 hours). Raise it for
long-running human or CI sessions; keep the default for service roles.
Updatable in place.

- default: `3600`

### spec.permissionsBoundary

`string | valueFrom`

An optional permissions boundary: a managed policy whose grants cap the
maximum permissions this role can ever have -- effective permissions are
the INTERSECTION of the boundary and the role's permission policies.
Reference an AwsIamPolicy's policy_arn output or pass a literal policy
ARN. Setting or changing the boundary is in-place; clearing it removes
the ceiling.

- references: AwsIamPolicy (`status.outputs.policy_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsIamPolicy, name: <that resource's name>, fieldPath: status.outputs.policy_arn}} -- a bare string does not parse

### spec.forceDetachPolicies

`bool`

Whether deleting the role force-detaches any policies still attached to
it (including attachments made outside this resource). Off by default:
deletion fails if out-of-band attachments exist, surfacing them instead
of silently severing another owner's wiring. Turn on for ephemeral or
CI-owned roles where teardown must always succeed.

## Validation Rules

- `path_format`: path must begin and end with '/' and contain no empty segments, e.g. '/' or '/service-roles/'
- `path_length`: path must be at most 512 characters
- `max_session_duration_range`: max_session_duration must be between 3600 (1h) and 43200 (12h) seconds

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsIamRole, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.roleArn` | `string` | The ARN of the IAM role (e.g. "arn:aws:iam::123456789012:role/my-role"). What most service integrations reference via status.outputs.role_arn -- Lambda's role, an ECS task role, a Step Functions execution role, and so on. To deliver this role to EC2 instances, wrap it in an AwsIamInstanceProfile instead of referencing the role directly. |
| `status.outputs.roleName` | `string` | The friendly name of the IAM role (mirrors metadata.name). What an AwsIamInstanceProfile's role field references, and what the AWS CLI and console use. |
| `status.outputs.roleId` | `string` | The stable unique ID AWS assigns to the role (e.g. "AROA..."). Unlike the ARN it never encodes the name or path, so it is what appears in policy aws:userid conditions and audit trails. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.managedPolicyArns` | AwsIamPolicy | `status.outputs.policy_arn` |
| `spec.permissionsBoundary` | AwsIamPolicy | `status.outputs.policy_arn` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
