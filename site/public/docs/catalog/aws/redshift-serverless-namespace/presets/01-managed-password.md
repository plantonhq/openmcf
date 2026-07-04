---
title: "Managed-Password Namespace"
description: "This preset creates a Redshift Serverless namespace -- the data plane of the serverless warehouse -- with the AWS-managed admin password. AWS generates the password, stores it in Secrets Manager, and..."
type: "preset"
rank: "01"
presetSlug: "01-managed-password"
componentSlug: "redshift-serverless-namespace"
componentTitle: "Redshift Serverless Namespace"
provider: "aws"
icon: "package"
order: 1
---

# Managed-Password Namespace

This preset creates a Redshift Serverless namespace -- the data plane of the serverless warehouse -- with the AWS-managed admin password. AWS generates the password, stores it in Secrets Manager, and rotates it on schedule; no secret ever touches the manifest or the IaC state. The secret's ARN surfaces as the `admin_password_secret_arn` output, which applications use to fetch credentials at runtime.

## When to Use

- The starting point for nearly every serverless warehouse -- pair it with an `AwsRedshiftServerlessWorkgroup` to get compute
- Teams that want zero secret handling: the manifest carries no password and the IaC state stores none
- Environments where credential rotation must happen without redeploys

## Key Configuration Choices

- **Managed password** (`manageAdminPassword: true`) -- AWS Secrets Manager owns the credential lifecycle; the recommended posture
- **Named first database** (`dbName: analytics`) -- Create-time only; additional databases are created with SQL
- **Explicit admin** (`adminUsername: admin`) -- Matches the AWS default; set it explicitly so the value is visible in review
- **AWS-owned data encryption** -- No `kmsKeyId` pin; stored data is encrypted with the AWS-owned Redshift service key

## Related Presets

- **02-customer-kms-and-roles** -- Use when data must be encrypted with a customer-managed KMS key and the engine needs IAM roles for COPY/UNLOAD/Spectrum
