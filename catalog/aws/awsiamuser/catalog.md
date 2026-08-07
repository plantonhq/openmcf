# AWS IAM User

Deploys an IAM user with configurable managed policy attachments, inline policies, and optional access key generation. The component is designed for CI/CD pipelines and service integrations that require long-lived programmatic credentials. It integrates with Planton's Provider Connections for AWS credential management and provides `user_arn`, `access_key_id`, and `secret_access_key` outputs for downstream consumption.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **IAM User** -- created with the specified username matching the pattern `[a-zA-Z0-9+=,.@_-]{1,64}`, at the optional IAM `path` (defaults to `/`), with an optional `permissionsBoundary` capping its maximum permissions
- **Managed Policy Attachments** -- one attachment per entry in `managedPolicyArns`, linking AWS-managed or customer-managed policies to the user. Each entry is a literal ARN or a reference to an AwsIamPolicy's `policy_arn` output.
- **Inline Policies** -- one inline policy per entry in `inlinePolicies`, embedded directly on the user with the specified policy name and JSON document
- **Access Key** -- created only when `disableAccessKeys` is `false` (default); generates an active access key pair for programmatic AWS API access. The secret access key is base64-encoded in outputs.
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the user

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **IAM permissions** -- the credentials used by the Provider Connection must have `iam:CreateUser`, `iam:AttachUserPolicy`, `iam:PutUserPolicy`, `iam:CreateAccessKey`, and related IAM permissions.
- **Managed policy ARNs** -- any AWS-managed or customer-managed policy ARNs listed as literals in `managedPolicyArns` must exist in the account before deployment. Policies managed in Planton are referenced instead (see the consumes table below) and resolve at deploy time.

## Deploy

### Console

Open the deployment store, find **AWS IAM User**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **CI/CD Pipeline** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamUser
metadata:
  name: deployer
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  userName: github-actions-deployer
  managedPolicyArns:
    - "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
    - "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
```

```shell
planton apply -f iam-user.yaml
```

This creates an IAM user with ECR push access and S3 read-only access, plus an active access key pair. No inline policies are configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an IAM user. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Managed vs. inline policies** -- Use `managedPolicyArns` for reusable, centrally maintained policies (both AWS-managed and customer-managed). Use `inlinePolicies` for user-specific permissions that should not be shared. Managed policies are easier to audit across users; inline policies are deleted when the user is deleted.

**Access key generation** -- By default, an access key pair is created for programmatic access. Set `disableAccessKeys: true` for identity-only users (console access, federation) or monitoring integrations that use external credential delivery. The secret access key is base64-encoded in outputs and should be stored securely.

**Username conventions** -- The `userName` field accepts 1-64 characters matching `[a-zA-Z0-9+=,.@_-]`. Use descriptive names that identify the purpose (e.g., `github-actions-deployer`, `datadog-readonly`) rather than generic names. Unlike most IAM identifiers, the name is renameable in place -- AWS keeps the same underlying user and rewrites its ARN.

**IAM path** -- `path` organizes users for wildcard policy matching: a permission scoped to `user/ci/*` covers every user under the `/ci/` path. Must begin and end with `/` (e.g. `/ci/`); defaults to `/`. Users, unlike roles and policies, can move paths in place.

**Permissions boundary** -- `permissionsBoundary` names a managed policy that caps what the user can EVER do; effective permissions are the intersection of the boundary and the attached policies. Especially valuable on users, whose access keys are long-lived. Reference an AwsIamPolicy's `policy_arn` output or pass a literal ARN.

**Deletion posture** -- `forceDestroy` (off by default) controls what happens when credentials created outside this resource exist at delete time -- console login profiles, extra access keys, MFA devices, SSH keys, signing certificates. Off, deletion fails loudly and surfaces them; on, teardown always succeeds. Turn it on for ephemeral or CI-owned users.

**Least privilege** -- Scope policies to the minimum permissions required. Prefer specific managed policies over broad ones like `AdministratorAccess`. Use inline policies for fine-grained, resource-level permissions (e.g., allowing ECS updates only on specific clusters).

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AwsIamPolicy** (optional) | `managedPolicyArns[]` | `status.outputs.policy_arn` |
| **AwsIamPolicy** (optional) | `permissionsBoundary` | `status.outputs.policy_arn` |

Both fields also accept literal ARNs -- literals are how AWS-managed policies like `arn:aws:iam::aws:policy/ReadOnlyAccess` attach.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `user_arn` | Amazon Resource Name of the IAM user | IAM policies, cross-account access grants |
| `user_name` | IAM user name | Policy attachment, resource tagging |
| `user_id` | Stable unique identifier of the IAM user | Audit logging, resource ownership tracking |
| `access_key_id` | Access key ID (present when access keys are enabled) | CI/CD pipeline `AWS_ACCESS_KEY_ID` configuration |
| `secret_access_key` | Base64-encoded secret access key (sensitive) | CI/CD pipeline `AWS_SECRET_ACCESS_KEY` configuration |
| `console_url` | AWS console sign-in URL | Documentation, onboarding guides |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**CI/CD pipeline user** -- ECR power user, S3 read-only, and ECS deployment permissions with access keys enabled. The standard configuration for GitHub Actions, GitLab CI, or Jenkins pipelines that build container images and deploy to ECS. Start from the **CI/CD Pipeline** preset.

**Read-only service user** -- Broad read-only access with access keys disabled. Suitable for third-party monitoring tools (Datadog, New Relic) or audit integrations that only need to observe resources without making changes. Start from the **Read-Only Service** preset.

## Works With

- [**AWS IAM Policy**](/cloud-catalog/aws-iam-policy) -- provides customer-managed policies for attachment via `managedPolicyArns` and the permissions boundary via `permissionsBoundary`