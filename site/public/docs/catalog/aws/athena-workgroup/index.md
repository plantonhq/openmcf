---
title: "Athena Workgroup"
description: "Athena Workgroup deployment documentation"
icon: "package"
order: 100
componentName: "awsathenaworkgroup"
---

# AWS Athena Workgroup

Deploys an Amazon Athena workgroup with configurable query result storage (customer S3 or AWS-managed), server-side encryption, per-query cost controls, IAM Identity Center integration, log delivery, and Apache Spark execution support. The workgroup enforces governance settings so individual queries cannot override result locations or encryption policies.

## What Gets Created

When you deploy an AwsAthenaWorkgroup resource, Planton provisions:

- **Athena Workgroup** — an `aws_athena_workgroup` resource with the specified name, state, configuration enforcement, engine version, and cost controls
- **Result Storage** — either a customer-managed S3 result configuration (`resultConfiguration`) with optional encryption and ACL settings, or AWS-managed result storage (`managedQueryResults`) that needs no bucket; the two are mutually exclusive by AWS's own rule
- **Identity Integration** — created only when set: IAM Identity Center trusted identity propagation (`identityCenter`) and S3 Access Grants for per-user result access (`s3AccessGrants`)
- **Log Delivery** — created only when set: CloudWatch Logs, Athena-managed, and/or S3 logging arms (`monitoring`)
- **Spark Content Encryption** — created only when set: a customer KMS key for notebook and session content (`customerContentEncryptionKmsKey`)

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **An S3 bucket** for storing query results (only if setting `resultConfiguration.outputLocation`; `managedQueryResults` needs none)
- **An AWS Glue Data Catalog** with databases and tables defined over your S3 data sources
- **A KMS key ARN** if using SSE_KMS/CSE_KMS results, managed-results encryption, or Spark content encryption
- **An IAM execution role** for Spark or Identity Center workgroups (standard SQL workgroups do not need one)
- **An IAM Identity Center instance** only if setting `identityCenter`

## Quick Start

Create a file `athena-workgroup.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAthenaWorkgroup
metadata:
  name: analytics
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: my-project
    pulumi.planton.dev/stack.name: dev.AwsAthenaWorkgroup.analytics
spec:
  region: us-east-1
  resultConfiguration:
    outputLocation: "s3://my-athena-results/analytics/"
```

Deploy:

```shell
planton apply -f athena-workgroup.yaml
```

This creates an Athena workgroup named `analytics` with query results stored in S3, configuration enforcement enabled (default), and CloudWatch metrics published (default).

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | AWS region where the workgroup will be created (e.g., `us-east-1`, `eu-west-1`). | Required; non-empty |

However, most practical deployments set at least `resultConfiguration.outputLocation`.

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `description` | `string` | — | Console description of the workgroup (max 1024 chars) |
| `state` | `string` | `ENABLED` | `ENABLED` or `DISABLED`; disabling rejects new queries but keeps configuration and history |
| `resultConfiguration` | `object` | — | Customer-managed S3 result storage and encryption settings |
| `resultConfiguration.outputLocation` | `string` | — | S3 URI where query results are stored (e.g., `s3://bucket/prefix/`) |
| `resultConfiguration.encryptionOption` | `string` | — | `SSE_S3`, `SSE_KMS`, or `CSE_KMS` |
| `resultConfiguration.kmsKeyArn` | `string` | — | KMS key ARN for SSE_KMS/CSE_KMS. Can reference AwsKmsKey resource via `valueFrom` |
| `resultConfiguration.expectedBucketOwner` | `string` | — | AWS account ID for cross-account S3 buckets |
| `resultConfiguration.s3AclOption` | `string` | — | `BUCKET_OWNER_FULL_CONTROL` for cross-account result ownership |
| `managedQueryResults` | `object` | — | AWS-managed result storage (24h retention, retrieved via Athena APIs); presence enables it; mutually exclusive with `outputLocation` |
| `managedQueryResults.kmsKey` | `string` | AWS-owned key | KMS key for managed-results encryption. Can reference AwsKmsKey via `valueFrom` |
| `bytesScannedCutoffPerQuery` | `int64` | `0` (no limit) | Max bytes a query can scan. Must be 0 or >= 10485760 (10 MB) |
| `enforceWorkgroupConfiguration` | `bool` | `true` | Lock settings so queries cannot override them |
| `publishCloudwatchMetricsEnabled` | `bool` | `true` | Publish query metrics to CloudWatch |
| `requesterPaysEnabled` | `bool` | `false` | Requester pays for S3 data access |
| `enableMinimumEncryptionConfiguration` | `bool` | `false` | Require at least SSE_S3 for all query results |
| `selectedEngineVersion` | `string` | `AUTO` | Athena engine version (`Athena engine version 3`, `PySpark engine version 3`, or `AUTO`) |
| `executionRole` | `string` | — | IAM role for Spark / Identity Center workgroups. Can reference AwsIamRole via `valueFrom` |
| `customerContentEncryptionKmsKey` | `string` | AWS-owned key | KMS key encrypting Spark notebook/session content. Can reference AwsKmsKey via `valueFrom` |
| `identityCenter.enableIdentityCenter` | `bool` | `false` | IAM Identity Center trusted identity propagation (create-time) |
| `identityCenter.identityCenterInstanceArn` | `string` | — | Identity Center instance ARN (create-time) |
| `s3AccessGrants.enableS3AccessGrants` | `bool` | `false` | Obtain result-bucket credentials from S3 Access Grants |
| `s3AccessGrants.authenticationType` | `string` | — | Only `DIRECTORY_IDENTITY` is supported by AWS |
| `s3AccessGrants.createUserLevelPrefix` | `bool` | `false` | Per-user prefix under the result location |
| `monitoring.cloudWatchLogging` | `object` | — | CloudWatch Logs delivery: `logGroup`, `logStreamNamePrefix`, `logTypes` (worker type → log streams) |
| `monitoring.managedLogging` | `object` | — | Athena-managed log storage; optional `kmsKey` |
| `monitoring.s3Logging` | `object` | — | S3 log delivery: `logLocation` (s3:// URI), optional `kmsKey` |
| `forceDestroy` | `bool` | `false` | Delete named queries and prepared statements on workgroup destroy |

## Examples

### Basic SQL Workgroup

A minimal workgroup directing query results to S3 with all governance defaults.

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAthenaWorkgroup
metadata:
  name: analytics-team
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: analytics
    pulumi.planton.dev/stack.name: dev.AwsAthenaWorkgroup.analytics-team
spec:
  region: us-east-1
  resultConfiguration:
    outputLocation: "s3://my-athena-results/analytics-team/"
```

### Cost-Controlled with SSE_S3

Workgroup with a 10 GB per-query scan limit and enforced minimum encryption.

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAthenaWorkgroup
metadata:
  name: data-science
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: my-org
    pulumi.planton.dev/project: data
    pulumi.planton.dev/stack.name: prod.AwsAthenaWorkgroup.data-science
spec:
  region: us-east-1
  resultConfiguration:
    outputLocation: "s3://data-science-results/queries/"
    encryptionOption: SSE_S3
  bytesScannedCutoffPerQuery: 10737418240
  enableMinimumEncryptionConfiguration: true
```

### Production KMS-Encrypted with valueFrom

Production workgroup with SSE_KMS encryption referencing a KMS key from another Planton resource.

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAthenaWorkgroup
metadata:
  name: prod-analytics
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: analytics
    pulumi.planton.dev/stack.name: prod.AwsAthenaWorkgroup.prod-analytics
spec:
  region: us-east-1
  resultConfiguration:
    outputLocation: "s3://prod-athena-results/queries/"
    encryptionOption: SSE_KMS
    kmsKeyArn:
      valueFrom:
        kind: AwsKmsKey
        name: analytics-encryption-key
        fieldPath: status.outputs.key_arn
  bytesScannedCutoffPerQuery: 53687091200
  enforceWorkgroupConfiguration: true
  publishCloudwatchMetricsEnabled: true
  enableMinimumEncryptionConfiguration: true
  selectedEngineVersion: "Athena engine version 3"
```

### Managed Results with No Bucket

BI-facing workgroup on AWS-managed result storage — no S3 bucket to create,
secure, or lifecycle. Results are retained 24 hours and retrieved through the
Athena APIs.

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsAthenaWorkgroup
metadata:
  name: bi-queries
  labels:
    planton.dev/provisioner: pulumi
    pulumi.planton.dev/organization: acme
    pulumi.planton.dev/project: analytics
    pulumi.planton.dev/stack.name: prod.AwsAthenaWorkgroup.bi-queries
spec:
  region: us-east-1
  managedQueryResults: {}
  bytesScannedCutoffPerQuery: 10737418240
```

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `workgroup_arn` | `string` | ARN of the Athena workgroup, used for IAM policies and cross-service references |
| `workgroup_name` | `string` | Name of the workgroup, used in Athena API calls (`StartQueryExecution`, etc.) |
| `effective_engine_version` | `string` | Actual engine version in use (resolved from `selectedEngineVersion` or `AUTO`) |

## Related Components

- [AWS S3 Bucket](/docs/catalog/aws/s3-bucket) — S3 bucket for query result storage
- [AWS KMS Key](/docs/catalog/aws/kms-key) — Customer-managed encryption for query results
- [AWS IAM Role](/docs/catalog/aws/iam-role) — Execution role for Spark workgroups
- [AWS CloudWatch Log Group](/docs/catalog/aws/cloudwatch-log-group) — Query execution logging
