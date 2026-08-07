# AWS Account Provider

This preset trusts workloads running in one AWS account: EC2 instances, Lambda functions, and ECS tasks federate into Google Cloud using their native AWS credentials — no tokens to mint, no keys to sync between clouds. The standard bridge for cross-cloud data flows and migrations.

## When to Use

- AWS-hosted workloads need to read/write GCS, BigQuery, or Pub/Sub
- A migration is running workloads on both clouds and syncing data between them
- You are replacing GCP service-account keys stored in AWS Secrets Manager

## Key Configuration Choices

- **Default attribute mapping (omitted)** — AWS providers map `google.subject` from the caller's ARN and derive `attribute.aws_role` automatically; only add a custom mapping when you need additional attributes.
- **Condition on `attribute.aws_role`** — the account-level trust is broad (every role in the account); the condition narrows federation to the one IAM role your workloads assume. Drop it only if genuinely every workload in the account should federate.
- **Per-role grants** — bind GCP roles to `principalSet://iam.googleapis.com/<pool name>/attribute.aws_role/<role arn>` so each AWS role maps to its own GCP access.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `ci-federation-pool` | Your GcpWorkloadIdentityPool resource name | The pool manifest's `metadata.name` |
| `123456789012` | Your 12-digit AWS account ID (in `accountId` and the condition's ARN) | AWS Console account menu |
| `<role-name>` | The IAM role your workloads assume | AWS IAM console |

## Related Presets

- **01-github-actions-oidc** — Trust GitHub Actions workflows
- **03-gitlab-ci-oidc** — Trust GitLab CI pipelines
