---
title: "Production Namespace with Customer KMS and Engine Roles"
description: "This preset creates a production-posture namespace: stored data encrypted with a customer-managed KMS key, engine IAM roles for COPY/UNLOAD and Redshift Spectrum composed from the resource graph, and..."
type: "preset"
rank: "02"
presetSlug: "02-customer-kms-and-roles"
componentSlug: "redshift-serverless-namespace"
componentTitle: "Redshift Serverless Namespace"
provider: "aws"
icon: "package"
order: 2
---

# Production Namespace with Customer KMS and Engine Roles

This preset creates a production-posture namespace: stored data encrypted with a customer-managed KMS key, engine IAM roles for COPY/UNLOAD and Redshift Spectrum composed from the resource graph, and all three audit log types exported to CloudWatch Logs. The admin password stays AWS-managed.

## When to Use

- Production warehouses with a compliance requirement to own the data-encryption key
- Pipelines that load from S3 (COPY needs an IAM role) or query external tables through Redshift Spectrum
- Environments where connection and query audit trails must land in CloudWatch Logs

## Key Configuration Choices

- **Customer-managed data key** (`kmsKeyId` → `AwsKmsKey`) -- You own rotation and access policy for the key that encrypts stored data; switching keys later is an in-place but long-running re-encryption
- **Engine roles by reference** (`iamRoles` → `AwsIamRole`) -- The COPY and Spectrum roles come from the resource graph, so policy changes roll through without touching this manifest
- **Default role** (`defaultIamRoleArn`) -- SQL that says `IAM_ROLE default` uses the COPY role; the default must also appear in `iamRoles` (AWS rejects a default it has not been given)
- **Full audit exports** (`logExports`) -- Connection attempts, every executed query, and user changes stream to CloudWatch Logs
- **Managed password** (`manageAdminPassword: true`) -- The recommended posture, unchanged from the starter preset

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `my-data-key` | Name of your `AwsKmsKey` resource | Your resource graph |
| `my-redshift-copy-role` | Name of the `AwsIamRole` the engine assumes for COPY/UNLOAD | Your resource graph |
| `my-redshift-spectrum-role` | Name of the `AwsIamRole` for Spectrum external tables | Your resource graph |

## Related Presets

- **01-managed-password** -- The minimal starting point when AWS-owned encryption and no engine roles are enough
