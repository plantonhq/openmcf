---
title: "OpenSearch Serverless Collection"
description: "Amazon OpenSearch Serverless collection deployment documentation"
icon: "package"
order: 100
componentName: "awsopensearchserverlesscollection"
---

# AWS OpenSearch Serverless Collection

Deploys an Amazon OpenSearch Serverless collection — a fully managed, auto-scaling OpenSearch workspace with no domains or nodes to size (capacity is billed in OpenSearch Compute Units), for search, time-series, or vector workloads. VECTORSEARCH collections are the vector store Bedrock knowledge bases consume.

## One Manifest, One Usable Collection

OpenSearch Serverless separates the collection from three account-level POLICY objects that attach to collections by name-pattern matching: encryption security policies, network security policies, and data access policies. This component scopes all three (plus index retention) to exactly this collection — the modules render each policy with rules matching only this collection's name, so one manifest owns one collection and everything that makes it usable.

## What Gets Created

- **Encryption Security Policy** — always created (AWS rejects collection creation without one): AWS-owned key by default, or the customer-managed KMS key referenced in `encryption.kmsKeyArn`
- **Collection** — SEARCH, TIMESERIES (default), or VECTORSEARCH, with optional standby replicas (default on; disable for the half-cost dev floor) and optional collection-group membership
- **Network Security Policy** — public reachability by default (SigV4-authenticated; see below), or VPC-endpoint-restricted
- **Data Access Policy** — created from your `dataAccess` rules; without at least one rule nothing can read or write data
- **Lifecycle Policy** — created from your `retentionRules` (index retention; indefinite without rules)

## "Public" Network Access Is Not Public Data

Network `allowFromPublic` controls reachability only: every request must still be SigV4-signed and authorized by a data access rule. A collection with public network access and no data access rules is unreachable for data — like an S3 bucket reachable over the internet but private. IAM identity permissions alone grant nothing in OpenSearch Serverless; data access policies are the sole data-plane authorization.

## Quick Start

Create a file `collection.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsOpenSearchServerlessCollection
metadata:
  name: app-search
  annotations:
    planton.dev/provisioner: pulumi
spec:
  region: us-west-2
  type: SEARCH
  standbyReplicas: DISABLED
  dataAccess:
    - principals:
        - valueFrom:
            kind: AwsIamRole
            name: app-role
            fieldPath: status.outputs.role_arn
      indexPermissions:
        - aoss:ReadDocument
        - aoss:WriteDocument
        - aoss:CreateIndex
```

Deploy it:

```bash
planton apply -f collection.yaml
```

## Cost Model

Collections bill by OpenSearch Compute Units (OCUs). With standby replicas ENABLED (the AWS default, production posture) the floor is 2 indexing + 2 search OCUs; DISABLED halves it — the right choice for dev/test. The collection name, type, standby replicas, group membership, and encryption key are all fixed at create time (changing them replaces the collection).

## Spec Reference

See [reference](v1alpha1/reference.md) for the complete field reference, and the [presets](presets/) for ready-to-deploy configurations.
