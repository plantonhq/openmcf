---
title: "Presets"
description: "Ready-to-deploy configuration presets for Redshift Serverless Namespace"
type: "preset-list"
componentSlug: "redshift-serverless-namespace"
componentTitle: "Redshift Serverless Namespace"
provider: "aws"
icon: "package"
order: 200
presets:
  - slug: "01-managed-password"
    rank: "01"
    title: "Managed-Password Namespace"
    excerpt: "This preset creates a Redshift Serverless namespace -- the data plane of the serverless warehouse -- with the AWS-managed admin password. AWS generates the password, stores it in Secrets Manager, and..."
  - slug: "02-customer-kms-and-roles"
    rank: "02"
    title: "Production Namespace with Customer KMS and Engine Roles"
    excerpt: "This preset creates a production-posture namespace: stored data encrypted with a customer-managed KMS key, engine IAM roles for COPY/UNLOAD and Redshift Spectrum composed from the resource graph, and..."
---

# Redshift Serverless Namespace Presets

Ready-to-deploy configuration presets for Redshift Serverless Namespace. Each preset is a complete manifest you can copy, customize, and deploy.
