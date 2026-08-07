# AWS IAM Role

Deploys an IAM role with a configurable trust policy, managed policy attachments, and inline policies. The role integrates with Planton's Provider Connections for AWS credential management and provides `role_arn` and `role_name` outputs that downstream Cloud Resources consume via ValueFromRef.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **IAM Role** -- created with the specified trust policy (assume-role document), description, and IAM path
- **Managed Policy Attachments** -- one attachment per entry in `managedPolicyArns`, linking AWS-managed or customer-managed policies to the role
- **Inline Policies** -- one inline policy per entry in `inlinePolicies`, embedded directly on the role with the specified policy name and JSON document
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the role

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **IAM permissions** -- the credentials used by the Provider Connection must have `iam:CreateRole`, `iam:AttachRolePolicy`, `iam:PutRolePolicy`, and related IAM permissions.
- **Managed policy ARNs** -- any AWS-managed or customer-managed policy ARNs listed in `managedPolicyArns` must exist in the account before deployment.

## Deploy

### Console

Open the deployment store, find **AWS IAM Role**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Lambda Execution** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsIamRole
metadata:
  name: eks-service-role
  org: acme-corp
  env: prod
spec:
  region: us-west-2
  description: "Service role for EKS cluster management"
  path: "/"
  trustPolicy:
    Version: "2012-10-17"
    Statement:
      - Effect: Allow
        Principal:
          Service: eks.amazonaws.com
        Action: "sts:AssumeRole"
  managedPolicyArns:
    - "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
```

```shell
planton apply -f iam-role.yaml
```

This creates an IAM role that EKS can assume, with the `AmazonEKSClusterPolicy` managed policy attached. No inline policies are configured. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an IAM role. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Trust policy** -- The `trustPolicy` field defines which AWS services, accounts, or identity providers can assume this role. It is a required JSON structure following the IAM policy grammar. Common trust principals include `eks.amazonaws.com`, `lambda.amazonaws.com`, `ecs-tasks.amazonaws.com`, and cross-account root ARNs.

**Managed vs. inline policies** -- Use `managedPolicyArns` for reusable, centrally maintained policies (both AWS-managed and customer-managed). Use `inlinePolicies` for role-specific permissions that should not be shared. Managed policies are easier to audit and update across roles; inline policies are deleted when the role is deleted.

**IAM path** -- The `path` field defaults to `/`. Use IAM paths like `/service-roles/` or `/application/` to organize roles and enable path-based IAM conditions in permission boundaries. The path is create-only: it is part of the role's ARN, so changing it replaces the role.

**Session duration** -- `maxSessionDuration` caps how long AssumeRole sessions on this role may last. Unset keeps AWS's 1-hour default (right for service roles); raise it up to 43200 seconds (12 hours) for long-running human SSO or CI sessions. Updates in place.

**Permissions boundary** -- `permissionsBoundary` optionally caps the role's maximum permissions: effective permissions are the INTERSECTION of the boundary and the role's attached policies. Reference an AwsIamPolicy's `policy_arn` output or pass a literal policy ARN — the standard guardrail when teams mint their own roles.

**Deletion posture** -- `forceDetachPolicies` is off by default: deleting the role fails if out-of-band policy attachments exist, surfacing them instead of silently severing another owner's wiring. Turn it on for ephemeral or CI-owned roles where teardown must always succeed.

## Outputs and Dependencies

### What This Component Consumes

| Field | References | Via |
|-------|-----------|-----|
| `managedPolicyArns[]` | AwsIamPolicy | `status.outputs.policy_arn` (literal ARNs also accepted — how AWS-managed policies attach) |
| `permissionsBoundary` | AwsIamPolicy | `status.outputs.policy_arn` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `role_arn` | Amazon Resource Name of the IAM role | EKS cluster service role, Lambda execution role, ECS task role, S3 replication role |
| `role_name` | IAM role name | Policy attachment, instance profile association |
| `role_id` | The role's unique ID (never reused, unlike names) | Trust and resource policies that must survive role re-creation |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Lambda execution role** -- A role with `lambda.amazonaws.com` trust and the `AWSLambdaBasicExecutionRole` managed policy for CloudWatch Logs access. The standard starting point for Lambda functions. Start from the **Lambda Execution** preset.

**ECS task execution role** -- A role with `ecs-tasks.amazonaws.com` trust and the `AmazonECSTaskExecutionRolePolicy` managed policy for pulling container images and writing logs. Start from the **ECS Task Execution** preset.

**EC2 SSM role** -- A role with `ec2.amazonaws.com` trust and the `AmazonSSMManagedInstanceCore` managed policy for Systems Manager Session Manager access. Eliminates the need for SSH bastion hosts. Start from the **EC2 SSM** preset.

## Works With

- **AwsIamPolicy** -- customer-managed policies attached via `managedPolicyArns` or used as the `permissionsBoundary`.
- **AwsIamInstanceProfile** -- wraps this role for delivery to EC2 instances (the only service that needs the wrapper).
- **AwsEksCluster / AwsEksNodeGroup / AwsEksAccessEntry / AwsEksFargateProfile / AwsLambda / AwsEcsService** -- downstream components that consume the `role_arn` output.
