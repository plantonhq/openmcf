---
title: "DynamoDB"
description: "DynamoDB deployment documentation"
icon: "package"
order: 100
componentName: "awsdynamodb"
---

# AWS DynamoDB

Deploys an Amazon DynamoDB table through a single declarative manifest: key schema and indexes, capacity in any of AWS's three shapes (on-demand, provisioned, pre-warmed), streams, Global Tables v2 multi-region replication, encryption, recovery, and the table-scoped governance surface -- resource policy, Kinesis change-data destination, and CloudWatch contributor insights.

## What Gets Created

When you deploy an AwsDynamodb resource, Planton provisions:

- **DynamoDB Table** — keyed by `metadata.name`, with the specified key schema, billing mode, indexes, streams, replicas, encryption, and recovery settings
- **Global Secondary Indexes** — from `globalSecondaryIndexes`, each with its own key schema (including multi-attribute keys), projection, and per-index capacity
- **Local Secondary Indexes** — from `localSecondaryIndexes` (create-time only; they can never be added or removed later)
- **Global Table Replicas** — from `replicas`, one active read/write copy per listed region, up to Multi-Region Strong Consistency with an optional witness region
- **Resource Policy** — when `resourcePolicy` is set, a resource-based IAM policy attached to the table
- **Kinesis Streaming Destination** — when `kinesisStreamingDestination` is set, item-level change data flows into the referenced Kinesis Data Stream
- **Contributor Insights** — when `contributorInsights.enabled` is `true`, CloudWatch per-key access profiling on the table and each opted-in GSI

All resources are tagged with Planton metadata (organization, environment, resource kind, resource ID).

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A KMS key** (optional) — reference an `AwsKmsKey` for customer-managed encryption at rest
- **A Kinesis Data Stream** (optional) — reference an `AwsKinesisStream` to receive change data
- **An S3 bucket with source data** (optional) — reference an `AwsS3Bucket` to seed the table via `importTable`

## Quick Start

Create a file `dynamodb.yaml`:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: my-table
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: pk
      type: S
  keySchema:
    - attributeName: pk
      keyType: HASH
```

Deploy:

```shell
planton apply -f dynamodb.yaml
```

This creates an on-demand DynamoDB table with a single string partition key.

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `region` | `string` | The AWS region the table is created in. | Must be a valid AWS region |
| `attributeDefinitions` | `object[]` | Key attributes only (name + type `S`/`N`/`B`) -- DynamoDB is schemaless beyond the keys. | Required unless the table is created by restore |
| `keySchema` | `object[]` | The primary key: exactly one `HASH` element and at most one `RANGE` element, `HASH` first. Create-time immutable. | Required unless the table is created by restore |

### Capacity

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `billingMode` | `string` | `PROVISIONED` | `PAY_PER_REQUEST` (on-demand, the recommended default for new tables) or `PROVISIONED` (reserved units). Empty keeps the AWS default -- which then requires `provisionedThroughput`. |
| `provisionedThroughput` | `object` | — | `readCapacityUnits` + `writeCapacityUnits`. Required when the effective billing mode is `PROVISIONED`; must stay unset for `PAY_PER_REQUEST`. |
| `onDemandThroughput` | `object` | — | `maxReadRequestUnits` + `maxWriteRequestUnits` spend ceilings on a `PAY_PER_REQUEST` table; requests beyond the ceiling throttle instead of billing. `-1` removes a previously-set ceiling. |
| `warmThroughput` | `object` | — | `readUnitsPerSecond` (>= 12000) + `writeUnitsPerSecond` (>= 4000) pre-warmed instant capacity. Increase-only: lowering it replaces the table. |

### Indexes

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `globalSecondaryIndexes` | `object[]` | `[]` | Alternate query shapes, edited in place. Each carries `name`, `keySchema` (1-4 `HASH` elements first, then 0-4 `RANGE` -- multi-attribute keys), `projection`, and capacity matching the table's billing mode (`provisionedThroughput`, `onDemandThroughput`, `warmThroughput`). |
| `localSecondaryIndexes` | `object[]` | `[]` | Alternate sort orders over the table's partition key: `name`, `rangeKey`, `projection`. **Create-time only** and permanently caps item collections at 10 GB -- prefer a GSI unless you need strongly consistent reads. |
| `*.projection.type` | `string` | — | `ALL`, `KEYS_ONLY`, or `INCLUDE` (with `nonKeyAttributes`). |

### Streams, Recovery, and Safety

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `streamEnabled` | `bool` | `false` | Emit an ordered change stream of item modifications. Required (with `NEW_AND_OLD_IMAGES`) when the table has replicas. |
| `streamViewType` | `string` | — | `KEYS_ONLY`, `NEW_IMAGE`, `OLD_IMAGE`, or `NEW_AND_OLD_IMAGES`. Required when streams are enabled. |
| `pointInTimeRecovery` | `object` | — | `enabled` + `recoveryPeriodInDays` (1-35; 0 keeps the AWS default of 35). Continuous backups with per-second restore granularity. |
| `serverSideEncryption` | `object` | — | `enabled` switches from the AWS-owned key to the AWS-managed `aws/dynamodb` key; set `kmsKeyArn` (an `AwsKmsKey` reference or literal ARN) for a customer-managed key. |
| `tableClass` | `string` | `STANDARD` | `STANDARD` or `STANDARD_INFREQUENT_ACCESS` (~60% cheaper storage, ~25% costlier reads/writes). |
| `deletionProtectionEnabled` | `bool` | `false` | Refuse table deletion while `true`. |
| `ttl` | `object` | — | `enabled` + `attributeName` holding expiry as epoch seconds; expired items delete free of write cost. |

### Global Tables

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `replicas` | `object[]` | `[]` | Global Tables v2: one active read/write replica per listed region (`regionName`, per-replica `kmsKeyArn` reference, `pointInTimeRecovery`, `deletionProtectionEnabled`, `propagateTags`, `consistencyMode`). Requires streams with `NEW_AND_OLD_IMAGES`. |
| `replicas[].consistencyMode` | `string` | `EVENTUAL` | `EVENTUAL` (classic global tables) or `STRONG` (Multi-Region Strong Consistency: exactly two STRONG replicas, or one plus the witness). |
| `globalTableWitness` | `object` | — | The MRSC witness region -- stores replicated writes for quorum but serves no reads or writes. Must accompany exactly one `STRONG` replica. |

### Create Sources (mutually exclusive)

| Field | Type | Description |
|-------|------|-------------|
| `restoreSourceName` | `string` | Create by point-in-time restore of a same-account, same-region table (with `restoreDateTime` XOR `restoreToLatestTime`). Key schema is inherited from the source. |
| `restoreSourceTableArn` | `string` | The cross-region / cross-account point-in-time restore form. |
| `restoreBackupArn` | `string` | Create by restoring an on-demand or AWS Backup backup. |
| `importTable` | `object` | Seed a brand-new table from S3 (`s3Bucket` reference, `inputFormat` `CSV`/`DYNAMODB_JSON`/`ION`, optional compression and CSV options) -- billed as a one-time import. |

### Governance

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `resourcePolicy` | `string` | — | A resource-based IAM policy document (JSON) attached to the table -- cross-account access without role assumption. |
| `kinesisStreamingDestination` | `object` | — | Streams item-level change data into a Kinesis Data Stream (`streamArn` reference; optional timestamp precision). One destination per table. |
| `contributorInsights` | `object` | — | `enabled`, optional `mode` (`ACCESSED_AND_THROTTLED_KEYS` or `THROTTLED_KEYS`), and `gsiIndexNames` that also get insights. |

## Examples

### On-Demand Table with Composite Key and Spend Ceiling

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: orders-table
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: customerId
      type: S
    - name: orderId
      type: S
  keySchema:
    - attributeName: customerId
      keyType: HASH
    - attributeName: orderId
      keyType: RANGE
  onDemandThroughput:
    maxReadRequestUnits: 10000
    maxWriteRequestUnits: 5000
```

### Production Table with Customer-Managed Encryption and Insights

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: audit-log
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: pk
      type: S
    - name: sk
      type: S
    - name: eventType
      type: S
  keySchema:
    - attributeName: pk
      keyType: HASH
    - attributeName: sk
      keyType: RANGE
  globalSecondaryIndexes:
    - name: event-type-index
      keySchema:
        - attributeName: eventType
          keyType: HASH
      projection:
        type: INCLUDE
        nonKeyAttributes:
          - userId
          - action
  streamEnabled: true
  streamViewType: NEW_AND_OLD_IMAGES
  pointInTimeRecovery:
    enabled: true
  serverSideEncryption:
    enabled: true
    kmsKeyArn:
      valueFrom:
        kind: AwsKmsKey
        name: audit-log-key
        fieldPath: status.outputs.key_arn
  deletionProtectionEnabled: true
  contributorInsights:
    enabled: true
    gsiIndexNames:
      - event-type-index
```

### Global Table with Strong Consistency

Two STRONG replicas give synchronous multi-region writes; swap one replica for `globalTableWitness` for the cheaper witness topology:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: payments-ledger
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: pk
      type: S
  keySchema:
    - attributeName: pk
      keyType: HASH
  streamEnabled: true
  streamViewType: NEW_AND_OLD_IMAGES
  replicas:
    - regionName: us-east-2
      consistencyMode: STRONG
    - regionName: us-west-2
      consistencyMode: STRONG
```

### Change-Data Fan-Out to Kinesis

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: events-table
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: pk
      type: S
  keySchema:
    - attributeName: pk
      keyType: HASH
  kinesisStreamingDestination:
    streamArn:
      valueFrom:
        kind: AwsKinesisStream
        name: events-cdc
        fieldPath: status.outputs.stream_arn
```

### Seed a Table from S3

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsDynamodb
metadata:
  name: catalog-table
spec:
  region: us-east-1
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: sku
      type: S
  keySchema:
    - attributeName: sku
      keyType: HASH
  importTable:
    s3Bucket:
      valueFrom:
        kind: AwsS3Bucket
        name: catalog-seed-data
        fieldPath: status.outputs.bucket_id
    s3KeyPrefix: exports/catalog/
    inputFormat: DYNAMODB_JSON
    inputCompressionType: GZIP
```

## Stack Outputs

After deployment, the following outputs are available in `status.outputs`:

| Output | Type | Description |
|--------|------|-------------|
| `table_name` | `string` | Name of the DynamoDB table -- what SDK calls and IAM policies reference |
| `table_arn` | `string` | ARN of the table -- the join key for policies and cross-service integrations |
| `table_id` | `string` | Provider-assigned table ID |
| `stream_arn` | `string` | Stream ARN, populated when `streamEnabled` is `true` -- what Lambda event-source mappings attach to |
| `stream_label` | `string` | Stream label, populated when `streamEnabled` is `true` |

## Related Components

- [AwsKmsKey](/docs/catalog/aws/kms-key) — customer-managed encryption keys for the table and its replicas
- [AwsKinesisStream](/docs/catalog/aws/kinesis-data-stream) — receives the table's item-level change data
- [AwsS3Bucket](/docs/catalog/aws/s3-bucket) — holds source data for table imports
- [AwsIamRole](/docs/catalog/aws/iam-role) — IAM roles with policies for table access
- [AwsLambda](/docs/catalog/aws/lambda) — triggered by DynamoDB Streams events
