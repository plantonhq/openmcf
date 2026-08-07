# AWS Redshift Serverless Namespace — Architecture and Design

## Overview

Amazon Redshift Serverless is the serverless product surface of AWS's columnar
data warehouse. Instead of a provisioned cluster with a fixed node type and
count, it splits the warehouse into two independent resources:

- A **namespace** — the data plane: databases, admin credentials, the
  data-encryption key, engine IAM roles, and audit log configuration.
- A **workgroup** — the compute plane: Redshift Processing Unit (RPU)
  capacity, VPC placement, and endpoint configuration.

This component models the namespace. Compute is the separate
`AwsRedshiftServerlessWorkgroup` component.

## Why Two Kinds

The namespace/workgroup split is AWS's own resource model
(`CreateNamespace` / `CreateWorkgroup` in the `redshift-serverless` API),
and it meets the split test on every axis:

- **Independent lifecycle** — a workgroup is created, resized, and destroyed
  without touching the data; a namespace outlives every workgroup that ever
  served it.
- **Many-per-parent** — multiple workgroups can attach to one namespace: a
  capped development workgroup and an autoscaling production workgroup can
  serve the SAME data with separate spend guardrails.
- **FK-referenced** — the workgroup references the namespace by name; the
  credentials API, snapshots, and usage limits key off one or the other.

Bundling them into one kind would make "second compute plane over the same
data" — a first-class serverless capability — unrepresentable.

## Credential Model

A serverless namespace has three credential postures:

1. **AWS-managed admin password (recommended, the default posture here)** —
   `manage_admin_password` makes AWS generate the password, store it in
   Secrets Manager, and rotate it on schedule. No secret exists in the
   manifest or the IaC state; the secret's ARN is exported as
   `admin_password_secret_arn`.
2. **Direct password** — `admin_user_password` (annotated sensitive). The
   value lands in IaC state; prefer the managed path.
3. **No admin credentials at all** — unlike the provisioned cluster (whose
   CreateCluster hard-requires a master username), a serverless namespace can
   be created without admin credentials; IAM identities obtain temporary
   database credentials through the `GetCredentials` API scoped to a
   workgroup.

## Design Decisions

- **Name basis** — the namespace name is `metadata.name` (create-time
  immutable in AWS), the same basis every AWS kind uses so a manifest deploys
  identically on both engines.
- **`namespace_name` is a stack output** — downstream references resolve
  against stack outputs, so the workgroup's namespace reference points at
  `status.outputs.namespace_name`, never at metadata.
- **Composition by reference** — KMS keys (`kms_key_id`,
  `admin_password_secret_kms_key_id`) and IAM roles (`iam_roles`,
  `default_iam_role_arn`) are `StringValueOrRef` fields; the namespace never
  creates or mutates resources that deserve to be their own nodes.
- **No restore surface** — the provider's namespace resource does not model
  restore-from-snapshot; snapshot and recovery operations are separate API
  surfaces (see the deliberate-skips table).

## Deliberately Skipped Provider Surface

| Provider surface | Verdict | Reason |
| --- | --- | --- |
| `admin_user_password_wo` (+ `_wo_version`) | Skip | Write-only plaintext arm; the managed-password default makes it redundant surface, and the sensitive-annotated direct arm covers the rare direct-password need |
| `aws_redshiftserverless_snapshot` | Defer | Operational point-in-time snapshots of a namespace; scheduling plumbing, not namespace shape — joins via the exported namespace name on demand |
| `aws_redshiftserverless_resource_policy` | Defer | Cross-account snapshot-sharing governance; joins via the exported ARN on demand |
| `aws_redshiftserverless_usage_limit` | Defer | Cost-governance surface keyed by namespace/workgroup ARN; joins via the exported ARN on demand |

## Billing Note

A namespace by itself accrues only managed-storage charges. All compute
billing (RPU-hours) follows the workgroups that serve it, and only while
queries execute.

## References

- Redshift Serverless overview: https://docs.aws.amazon.com/redshift/latest/mgmt/serverless-whatis.html
- Namespaces and workgroups: https://docs.aws.amazon.com/redshift/latest/mgmt/serverless-workgroups-and-namespaces.html
- Temporary credentials (GetCredentials): https://docs.aws.amazon.com/redshift-serverless/latest/APIReference/API_GetCredentials.html
- Secrets Manager integration: https://docs.aws.amazon.com/redshift/latest/mgmt/redshift-secrets-manager-integration.html
