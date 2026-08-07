# AWS Redshift Serverless Namespace

Deploys an Amazon Redshift Serverless namespace -- the data plane of the serverless warehouse: the first database, admin credentials with optional Secrets Manager management, KMS data encryption, engine IAM roles for COPY/UNLOAD/Spectrum, and CloudWatch audit log exports. Compute attaches separately as `AwsRedshiftServerlessWorkgroup` nodes, so the same data can be served by multiple independently-sized compute planes. KMS keys and IAM roles compose by reference.

## What Gets Created

When you deploy an AwsRedshiftServerlessNamespace resource, Planton provisions:

- **Redshift Serverless Namespace** — a `redshiftserverless.Namespace` carrying the first database, admin credential strategy, data-encryption key, engine IAM roles, and audit log exports

No compute is created -- attach one or more `AwsRedshiftServerlessWorkgroup` nodes to query the data. RPU-hours accrue only while queries execute on a workgroup, so a namespace by itself costs only managed storage.

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A KMS key ARN** if encrypting stored data (or the managed password secret) with a customer-managed key
- **IAM role ARNs** if the engine needs to access S3, DynamoDB, Glue Data Catalog, or other AWS services during COPY/UNLOAD/Spectrum

## Quick Start

Create a file `redshift-serverless-namespace.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftServerlessNamespace
metadata:
  name: my-analytics-data
spec:
  region: us-west-2
  dbName: analytics
  adminUsername: admin
  manageAdminPassword: true
```

Deploy:

```shell
planton apply -f redshift-serverless-namespace.yaml
```

This creates a namespace whose admin password is generated, stored, and rotated by AWS Secrets Manager -- no secret in the manifest or the IaC state -- with the secret's ARN exported for applications to fetch credentials at runtime.

## Configuration Reference

### Required Fields

| Field | Type | Description |
| --- | --- | --- |
| `region` | string | AWS region; workgroups that attach must live in the same region |

### Credentials

| Field | Type | Description |
| --- | --- | --- |
| `manageAdminPassword` | bool | AWS generates, stores, and rotates the admin password in Secrets Manager (recommended) |
| `adminUserPassword` | string (sensitive) | Directly supplied admin password; mutually exclusive with the managed strategy |
| `adminPasswordSecretKmsKeyId` | ref → AwsKmsKey | KMS key for the managed secret; empty uses `aws/secretsmanager` |
| `adminUsername` | string | Empty keeps the AWS default (`admin`) |

### Data

| Field | Type | Description |
| --- | --- | --- |
| `dbName` | string | The first database (create-time only); empty keeps `dev` |
| `kmsKeyId` | ref → AwsKmsKey | Customer-managed data-encryption key; empty keeps the AWS-owned key |
| `iamRoles` | list of refs → AwsIamRole | Roles the engine assumes for COPY/UNLOAD/Spectrum |
| `defaultIamRoleArn` | ref → AwsIamRole | The role used when SQL says `IAM_ROLE default`; must also be in `iamRoles` |
| `logExports` | list | Any of `connectionlog`, `useractivitylog`, `userlog` |

## Stack Outputs

| Output | Description |
| --- | --- |
| `namespace_name` | The join key workgroups attach with |
| `namespace_id` | The unique identifier AWS assigns |
| `arn` | The namespace ARN, for IAM policies and usage limits |
| `db_name` | The first database's name |
| `admin_password_secret_arn` | The managed-password secret's ARN (managed mode only) |
