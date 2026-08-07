# AwsDynamodb

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsDynamodbSpec defines an Amazon DynamoDB table: key schema and
indexes, capacity (on-demand, provisioned, or price-tuned with warm
throughput), streams, global-table replication, encryption, recovery,
and the table-scoped governance surface (resource policy, Kinesis
change-data destinations, CloudWatch contributor insights).

The table name comes from metadata.name (create-time immutable in
AWS). Everything the table composes with attaches by reference: the
KMS key that encrypts it, the Kinesis streams that receive its change
data, and the S3 bucket an import seeds it from. A table is a true
leaf -- nothing has to exist before it.

Capacity guidance: PAY_PER_REQUEST (on-demand) is the right default
for almost every new table -- zero capacity planning, scale-to-zero
cost, optional ceiling via on_demand_throughput. Choose PROVISIONED
only for sustained, predictable traffic where reserved capacity
pricing wins.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsDynamodb
metadata:
  name: awsdynamodb-demo
spec:
  region: us-west-2
  billingMode: PAY_PER_REQUEST
  attributeDefinitions:
    - name: pk
      type: S
    - name: sk
      type: S
    - name: gsi1pk
      type: S
  keySchema:
    - attributeName: pk
      keyType: HASH
    - attributeName: sk
      keyType: RANGE
  globalSecondaryIndexes:
    - name: gsi1
      keySchema:
        - attributeName: gsi1pk
          keyType: HASH
      projection:
        type: ALL
  ttl:
    enabled: true
    attributeName: expiresAt
  streamEnabled: true
  streamViewType: NEW_AND_OLD_IMAGES
  pointInTimeRecovery:
    enabled: true
  serverSideEncryption:
    enabled: true
  deletionProtectionEnabled: false
  contributorInsights:
    enabled: true
    gsiIndexNames:
      - gsi1
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.billingMode` | `string` |  |  |  |
| `spec.attributeDefinitions` | `[]AwsDynamodbAttribute` |  |  |  |
| `spec.attributeDefinitions[].name` | `string` | yes |  |  |
| `spec.attributeDefinitions[].type` | `string` | yes |  |  |
| `spec.keySchema` | `[]AwsDynamodbKeySchemaElement` |  |  |  |
| `spec.keySchema[].attributeName` | `string` | yes |  |  |
| `spec.keySchema[].keyType` | `string` | yes |  |  |
| `spec.provisionedThroughput` | `AwsDynamodbProvisionedThroughput` |  |  |  |
| `spec.provisionedThroughput.readCapacityUnits` | `int64` |  |  |  |
| `spec.provisionedThroughput.writeCapacityUnits` | `int64` |  |  |  |
| `spec.onDemandThroughput` | `AwsDynamodbOnDemandThroughput` |  |  |  |
| `spec.onDemandThroughput.maxReadRequestUnits` | `int64` |  |  |  |
| `spec.onDemandThroughput.maxWriteRequestUnits` | `int64` |  |  |  |
| `spec.warmThroughput` | `AwsDynamodbWarmThroughput` |  |  |  |
| `spec.warmThroughput.readUnitsPerSecond` | `int64` |  |  |  |
| `spec.warmThroughput.writeUnitsPerSecond` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes` | `[]AwsDynamodbGlobalSecondaryIndex` |  |  |  |
| `spec.globalSecondaryIndexes[].name` | `string` | yes |  |  |
| `spec.globalSecondaryIndexes[].keySchema` | `[]AwsDynamodbKeySchemaElement` | yes |  |  |
| `spec.globalSecondaryIndexes[].keySchema[].attributeName` | `string` | yes |  |  |
| `spec.globalSecondaryIndexes[].keySchema[].keyType` | `string` | yes |  |  |
| `spec.globalSecondaryIndexes[].projection` | `AwsDynamodbProjection` | yes |  |  |
| `spec.globalSecondaryIndexes[].projection.type` | `string` | yes |  |  |
| `spec.globalSecondaryIndexes[].projection.nonKeyAttributes` | `[]string` |  |  |  |
| `spec.globalSecondaryIndexes[].provisionedThroughput` | `AwsDynamodbProvisionedThroughput` |  |  |  |
| `spec.globalSecondaryIndexes[].provisionedThroughput.readCapacityUnits` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes[].provisionedThroughput.writeCapacityUnits` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes[].onDemandThroughput` | `AwsDynamodbOnDemandThroughput` |  |  |  |
| `spec.globalSecondaryIndexes[].onDemandThroughput.maxReadRequestUnits` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes[].onDemandThroughput.maxWriteRequestUnits` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes[].warmThroughput` | `AwsDynamodbWarmThroughput` |  |  |  |
| `spec.globalSecondaryIndexes[].warmThroughput.readUnitsPerSecond` | `int64` |  |  |  |
| `spec.globalSecondaryIndexes[].warmThroughput.writeUnitsPerSecond` | `int64` |  |  |  |
| `spec.localSecondaryIndexes` | `[]AwsDynamodbLocalSecondaryIndex` |  |  |  |
| `spec.localSecondaryIndexes[].name` | `string` | yes |  |  |
| `spec.localSecondaryIndexes[].rangeKey` | `string` | yes |  |  |
| `spec.localSecondaryIndexes[].projection` | `AwsDynamodbProjection` | yes |  |  |
| `spec.localSecondaryIndexes[].projection.type` | `string` | yes |  |  |
| `spec.localSecondaryIndexes[].projection.nonKeyAttributes` | `[]string` |  |  |  |
| `spec.ttl` | `AwsDynamodbTtl` |  |  |  |
| `spec.ttl.enabled` | `bool` |  |  |  |
| `spec.ttl.attributeName` | `string` |  |  |  |
| `spec.streamEnabled` | `bool` |  |  |  |
| `spec.streamViewType` | `string` |  |  |  |
| `spec.pointInTimeRecovery` | `AwsDynamodbPointInTimeRecovery` |  |  |  |
| `spec.pointInTimeRecovery.enabled` | `bool` |  |  |  |
| `spec.pointInTimeRecovery.recoveryPeriodInDays` | `int32` |  |  |  |
| `spec.serverSideEncryption` | `AwsDynamodbServerSideEncryption` |  |  |  |
| `spec.serverSideEncryption.enabled` | `bool` |  |  |  |
| `spec.serverSideEncryption.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.tableClass` | `string` |  |  |  |
| `spec.deletionProtectionEnabled` | `bool` |  |  |  |
| `spec.contributorInsights` | `AwsDynamodbContributorInsights` |  |  |  |
| `spec.contributorInsights.enabled` | `bool` |  |  |  |
| `spec.contributorInsights.mode` | `string` |  |  |  |
| `spec.contributorInsights.gsiIndexNames` | `[]string` |  |  |  |
| `spec.resourcePolicy` | `string` |  |  |  |
| `spec.kinesisStreamingDestination` | `AwsDynamodbKinesisStreamingDestination` |  |  |  |
| `spec.kinesisStreamingDestination.streamArn` | `string \| valueFrom` | yes |  | AwsKinesisStream (`status.outputs.stream_arn`) |
| `spec.kinesisStreamingDestination.approximateCreationDateTimePrecision` | `string` |  |  |  |
| `spec.replicas` | `[]AwsDynamodbReplica` |  |  |  |
| `spec.replicas[].regionName` | `string` | yes |  |  |
| `spec.replicas[].kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.replicas[].pointInTimeRecovery` | `bool` |  |  |  |
| `spec.replicas[].deletionProtectionEnabled` | `bool` |  |  |  |
| `spec.replicas[].propagateTags` | `bool` |  |  |  |
| `spec.replicas[].consistencyMode` | `string` |  |  |  |
| `spec.globalTableWitness` | `AwsDynamodbGlobalTableWitness` |  |  |  |
| `spec.globalTableWitness.regionName` | `string` | yes |  |  |
| `spec.restoreSourceName` | `string` |  |  |  |
| `spec.restoreSourceTableArn` | `string` |  |  |  |
| `spec.restoreDateTime` | `string` |  |  |  |
| `spec.restoreToLatestTime` | `bool` |  |  |  |
| `spec.restoreBackupArn` | `string` |  |  |  |
| `spec.importTable` | `AwsDynamodbImportTable` |  |  |  |
| `spec.importTable.s3Bucket` | `string \| valueFrom` | yes |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.importTable.s3BucketOwner` | `string` |  |  |  |
| `spec.importTable.s3KeyPrefix` | `string` |  |  |  |
| `spec.importTable.inputFormat` | `string` | yes |  |  |
| `spec.importTable.inputCompressionType` | `string` |  |  |  |
| `spec.importTable.csv` | `AwsDynamodbImportTableCsv` |  |  |  |
| `spec.importTable.csv.delimiter` | `string` |  |  |  |
| `spec.importTable.csv.headerList` | `[]string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the table is created in. Must match the region of
any KMS key or Kinesis stream it references. Replicas live in
OTHER regions and are listed under replicas.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.billingMode

`string`

How the table is billed and how capacity is managed:
"PAY_PER_REQUEST" (on-demand -- pay per read/write, no capacity
planning, the recommended default for new tables) or "PROVISIONED"
(reserved read/write units, required for reserved-capacity
pricing). Empty keeps the AWS default (PROVISIONED) -- which then
requires provisioned_throughput, so most manifests set this
explicitly. Switching modes on a live table is an in-place update
(AWS allows one switch per 24 hours).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["PROVISIONED","PAY_PER_REQUEST"]}}

### spec.attributeDefinitions

`[]AwsDynamodbAttribute`

Attributes referenced by the table key schema and by every index
key schema. Only KEY attributes are declared here -- DynamoDB is
schemaless for everything else. Required unless the table is
created by restore (the source table carries the schema).

### spec.attributeDefinitions[].name

`string` · required

The attribute name.

- rule: {"string":{"minLen":"1"}}

### spec.attributeDefinitions[].type

`string` · required

The scalar type: "S" (string), "N" (number), or "B" (binary).

- rule: {"required":true,"string":{"in":["S","N","B"]}}

### spec.keySchema

`[]AwsDynamodbKeySchemaElement`

The table's primary key: exactly one HASH (partition) element and
at most one RANGE (sort) element, each naming a declared
attribute. Create-time immutable -- changing the key schema
replaces the table. Required unless the table is created by
restore (the source table carries the key schema).

### spec.keySchema[].attributeName

`string` · required

The declared attribute this key element uses.

- rule: {"string":{"minLen":"1"}}

### spec.keySchema[].keyType

`string` · required

"HASH" (partition key) or "RANGE" (sort key).

- rule: {"required":true,"string":{"in":["HASH","RANGE"]}}

### spec.provisionedThroughput

`AwsDynamodbProvisionedThroughput`

Reserved read/write capacity for the table. Required when the
effective billing mode is PROVISIONED; must stay unset for
PAY_PER_REQUEST. Pair with auto-scaling needs by sizing for the
sustained baseline -- on-demand is usually the better answer for
spiky traffic.

### spec.provisionedThroughput.readCapacityUnits

`int64`

Reserved read capacity units (one strongly consistent 4 KB read
per second each).

- rule: {"int64":{"gte":"0"}}

### spec.provisionedThroughput.writeCapacityUnits

`int64`

Reserved write capacity units (one 1 KB write per second each).

- rule: {"int64":{"gte":"0"}}

### spec.onDemandThroughput

`AwsDynamodbOnDemandThroughput`

Optional ceilings on on-demand consumption -- a spend guardrail
for PAY_PER_REQUEST tables. Requests beyond the ceiling are
throttled rather than billed. Only meaningful on
PAY_PER_REQUEST tables.

### spec.onDemandThroughput.maxReadRequestUnits

`int64`

Maximum read request units per second. 0 = no ceiling configured;
-1 removes a previously-set ceiling.

- rule: {"int64":{"gte":"-1"}}

### spec.onDemandThroughput.maxWriteRequestUnits

`int64`

Maximum write request units per second. 0 = no ceiling configured;
-1 removes a previously-set ceiling.

- rule: {"int64":{"gte":"-1"}}

### spec.warmThroughput

`AwsDynamodbWarmThroughput`

Pre-warmed minimum throughput the table can serve instantly,
decoupled from billing mode -- for launch events and traffic
cliffs where waiting for organic scale-up is not acceptable. AWS
only allows warm throughput to INCREASE; lowering it replaces the
table. Leave unset to keep AWS's defaults (12,000 reads/s, 4,000
writes/s warm).

### spec.warmThroughput.readUnitsPerSecond

`int64`

Warm read units per second. 0 keeps the AWS default (12,000);
when set, AWS requires at least 12,000.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int64":{"gte":"12000"}}

### spec.warmThroughput.writeUnitsPerSecond

`int64`

Warm write units per second. 0 keeps the AWS default (4,000);
when set, AWS requires at least 4,000.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int64":{"gte":"4000"}}

### spec.globalSecondaryIndexes

`[]AwsDynamodbGlobalSecondaryIndex`

Global secondary indexes: alternate query shapes with their own
key schema (including multi-attribute keys), projection, and --
on PROVISIONED tables -- their own capacity. Added, modified, and
removed in place on a live table (one GSI mutation at a time, an
AWS serialization rule).

- rule: GSI key_schema must start with a HASH element and carry 1-4 HASH elements (all before any RANGE) and at most 4 RANGE elements

### spec.globalSecondaryIndexes[].name

`string` · required

The index name, unique within the table.

- rule: {"string":{"minLen":"1"}}

### spec.globalSecondaryIndexes[].keySchema

`[]AwsDynamodbKeySchemaElement` · required

The index key: 1-4 HASH elements first (multi-attribute partition
keys), then 0-4 RANGE elements (multi-attribute sort keys). Every
element names a declared attribute. The common case is one HASH
and at most one RANGE.

- rule: {"repeated":{"minItems":"1"}}

### spec.globalSecondaryIndexes[].keySchema[].attributeName

`string` · required

The declared attribute this key element uses.

- rule: {"string":{"minLen":"1"}}

### spec.globalSecondaryIndexes[].keySchema[].keyType

`string` · required

"HASH" (partition key) or "RANGE" (sort key).

- rule: {"required":true,"string":{"in":["HASH","RANGE"]}}

### spec.globalSecondaryIndexes[].projection

`AwsDynamodbProjection` · required

Which attributes the index carries.

- rule: {"required":true}
- rule: non_key_attributes must be set when projection type is INCLUDE and must stay empty otherwise

### spec.globalSecondaryIndexes[].projection.type

`string` · required

"ALL" (every attribute), "KEYS_ONLY" (table + index keys), or
"INCLUDE" (keys plus non_key_attributes). Projecting less makes
the index cheaper to store and write; projecting too little forces
costly fetch-backs to the table at query time.

- rule: {"required":true,"string":{"in":["ALL","KEYS_ONLY","INCLUDE"]}}

### spec.globalSecondaryIndexes[].projection.nonKeyAttributes

`[]string`

The non-key attributes projected when type is "INCLUDE". These do
not need to be declared in attribute_definitions.

- rule: {"repeated":{"unique":true}}

### spec.globalSecondaryIndexes[].provisionedThroughput

`AwsDynamodbProvisionedThroughput`

Per-index reserved capacity. Required when the table's effective
billing mode is PROVISIONED; must stay unset for PAY_PER_REQUEST.

### spec.globalSecondaryIndexes[].provisionedThroughput.readCapacityUnits

`int64`

Reserved read capacity units (one strongly consistent 4 KB read
per second each).

- rule: {"int64":{"gte":"0"}}

### spec.globalSecondaryIndexes[].provisionedThroughput.writeCapacityUnits

`int64`

Reserved write capacity units (one 1 KB write per second each).

- rule: {"int64":{"gte":"0"}}

### spec.globalSecondaryIndexes[].onDemandThroughput

`AwsDynamodbOnDemandThroughput`

Per-index on-demand ceilings (PAY_PER_REQUEST tables only).

### spec.globalSecondaryIndexes[].onDemandThroughput.maxReadRequestUnits

`int64`

Maximum read request units per second. 0 = no ceiling configured;
-1 removes a previously-set ceiling.

- rule: {"int64":{"gte":"-1"}}

### spec.globalSecondaryIndexes[].onDemandThroughput.maxWriteRequestUnits

`int64`

Maximum write request units per second. 0 = no ceiling configured;
-1 removes a previously-set ceiling.

- rule: {"int64":{"gte":"-1"}}

### spec.globalSecondaryIndexes[].warmThroughput

`AwsDynamodbWarmThroughput`

Per-index pre-warmed throughput; increase-only, like the table's.

### spec.globalSecondaryIndexes[].warmThroughput.readUnitsPerSecond

`int64`

Warm read units per second. 0 keeps the AWS default (12,000);
when set, AWS requires at least 12,000.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int64":{"gte":"12000"}}

### spec.globalSecondaryIndexes[].warmThroughput.writeUnitsPerSecond

`int64`

Warm write units per second. 0 keeps the AWS default (4,000);
when set, AWS requires at least 4,000.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int64":{"gte":"4000"}}

### spec.localSecondaryIndexes

`[]AwsDynamodbLocalSecondaryIndex`

Local secondary indexes: alternate sort orders sharing the table's
partition key. CREATE-TIME ONLY -- LSIs can never be added or
removed after the table exists, and their presence permanently
caps each item collection at 10 GB. Prefer a GSI unless you need
strongly consistent reads on the alternate sort order.

### spec.localSecondaryIndexes[].name

`string` · required

The index name, unique within the table.

- rule: {"string":{"minLen":"1"}}

### spec.localSecondaryIndexes[].rangeKey

`string` · required

The alternate sort-key attribute (the partition key is always the
table's own HASH key). Must be declared in attribute_definitions.

- rule: {"string":{"minLen":"1"}}

### spec.localSecondaryIndexes[].projection

`AwsDynamodbProjection` · required

Which attributes the index carries.

- rule: {"required":true}
- rule: non_key_attributes must be set when projection type is INCLUDE and must stay empty otherwise

### spec.localSecondaryIndexes[].projection.type

`string` · required

"ALL" (every attribute), "KEYS_ONLY" (table + index keys), or
"INCLUDE" (keys plus non_key_attributes). Projecting less makes
the index cheaper to store and write; projecting too little forces
costly fetch-backs to the table at query time.

- rule: {"required":true,"string":{"in":["ALL","KEYS_ONLY","INCLUDE"]}}

### spec.localSecondaryIndexes[].projection.nonKeyAttributes

`[]string`

The non-key attributes projected when type is "INCLUDE". These do
not need to be declared in attribute_definitions.

- rule: {"repeated":{"unique":true}}

### spec.ttl

`AwsDynamodbTtl`

Time-to-live: DynamoDB deletes items (within ~48h, free of write
cost) once the epoch-seconds value in the named attribute passes.

- rule: attribute_name is required when TTL is enabled and must stay empty when disabled

### spec.ttl.enabled

`bool`

Turn TTL on.

### spec.ttl.attributeName

`string`

The attribute holding the expiry time as epoch seconds. Required
when enabled; must stay empty when disabled.

### spec.streamEnabled

`bool`

Emit an ordered change stream of item modifications, consumable by
Lambda event sources and Kinesis. Required (with view type
NEW_AND_OLD_IMAGES) when the table has replicas.

### spec.streamViewType

`string`

What each stream record carries: "KEYS_ONLY", "NEW_IMAGE",
"OLD_IMAGE", or "NEW_AND_OLD_IMAGES". Required when stream_enabled
is true; must stay empty when disabled. Global tables require
NEW_AND_OLD_IMAGES.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["KEYS_ONLY","NEW_IMAGE","OLD_IMAGE","NEW_AND_OLD_IMAGES"]}}

### spec.pointInTimeRecovery

`AwsDynamodbPointInTimeRecovery`

Continuous backups with per-second restore granularity over the
recovery window. The insurance policy for fat-finger deletes and
bad deploys -- production tables should enable it.

- rule: recovery_period_in_days only applies when point-in-time recovery is enabled

### spec.pointInTimeRecovery.enabled

`bool`

Turn point-in-time recovery on.

### spec.pointInTimeRecovery.recoveryPeriodInDays

`int32`

How many days of history to retain, 1-35. 0 keeps the AWS default
(35 days).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":35,"gte":1}}

### spec.serverSideEncryption

`AwsDynamodbServerSideEncryption`

Encryption at rest. DynamoDB always encrypts with an AWS-owned key
even when this is unset; enable this to switch to the AWS-managed
aws/dynamodb key or a customer-managed KMS key you control
(required for cross-account access patterns and key-rotation
policies of your own).

- rule: kms_key_arn only applies when server-side encryption is enabled

### spec.serverSideEncryption.enabled

`bool`

Switch from the AWS-owned key to an AWS-managed or
customer-managed KMS key.

### spec.serverSideEncryption.kmsKeyArn

`string | valueFrom`

The customer-managed KMS key that encrypts the table. Empty (with
enabled true) uses the AWS-managed aws/dynamodb key. Reference an
AwsKmsKey key_arn output or pass a literal key ARN. Required to be
configured when the table is created by cross-region/cross-account
restore.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.tableClass

`string`

Storage class: "STANDARD" (the default) or
"STANDARD_INFREQUENT_ACCESS" (~60% cheaper storage, ~25% costlier
reads/writes -- for large, rarely-read tables like audit logs).
Empty keeps the AWS default. Switchable in place (twice per 30
days per AWS).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["STANDARD","STANDARD_INFREQUENT_ACCESS"]}}

### spec.deletionProtectionEnabled

`bool`

Refuse deletion of the table while true. Flip to false (an
in-place update) before a genuine teardown.

### spec.contributorInsights

`AwsDynamodbContributorInsights`

CloudWatch Contributor Insights: per-key access profiling that
answers "which partition keys are hot / throttled". Enables at the
table level and, optionally, per GSI.

- rule: mode and gsi_index_names only apply when contributor insights is enabled

### spec.contributorInsights.enabled

`bool`

Enable insights on the table itself.

### spec.contributorInsights.mode

`string`

What to profile: "ACCESSED_AND_THROTTLED_KEYS" (the AWS default)
or "THROTTLED_KEYS" (cheaper -- only throttle diagnostics). Empty
keeps the AWS default. Applies to the table and every listed
index.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ACCESSED_AND_THROTTLED_KEYS","THROTTLED_KEYS"]}}

### spec.contributorInsights.gsiIndexNames

`[]string`

Global secondary indexes (by name) that also get insights, in
addition to the table.

- rule: {"repeated":{"unique":true}}

### spec.resourcePolicy

`string`

A resource-based IAM policy attached to the table, as a JSON
document -- cross-account access grants without assuming roles.
Table-scoped only (stream policies are a separate niche surface).

### spec.kinesisStreamingDestination

`AwsDynamodbKinesisStreamingDestination`

A Kinesis Data Stream that receives the table's item-level change
data -- the fan-out path for analytics and search-indexing
pipelines (independent of DynamoDB Streams and usable alongside
it). AWS allows exactly one Kinesis destination per table.

### spec.kinesisStreamingDestination.streamArn

`string | valueFrom` · required

The destination Kinesis Data Stream. Reference an AwsKinesisStream
stream_arn output or pass a literal stream ARN.

- references: AwsKinesisStream (`status.outputs.stream_arn`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKinesisStream, name: <that resource's name>, fieldPath: status.outputs.stream_arn}} -- a bare string does not parse

### spec.kinesisStreamingDestination.approximateCreationDateTimePrecision

`string`

Timestamp precision on emitted records: "MILLISECOND" or
"MICROSECOND". Empty keeps the AWS default (MICROSECOND).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["MILLISECOND","MICROSECOND"]}}

### spec.replicas

`[]AwsDynamodbReplica`

Global Tables v2: multi-region replicas of this table, each an
active read/write endpoint. Requires streams with
NEW_AND_OLD_IMAGES. Adding/removing an entry adds/removes that
region's replica in place. For Multi-Region Strong Consistency,
set consistency_mode STRONG on every replica (exactly two
replicas, or one replica plus global_table_witness).

### spec.replicas[].regionName

`string` · required

The region the replica lives in. Must differ from the table's own
region.

- rule: {"string":{"minLen":"1"}}

### spec.replicas[].kmsKeyArn

`string | valueFrom`

The customer-managed KMS key encrypting THIS replica (each region
encrypts independently). Empty uses the AWS-managed aws/dynamodb
key in the replica region. Changing it recreates the replica, not
the table. Reference an AwsKmsKey key_arn output or pass a literal
key ARN.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.replicas[].pointInTimeRecovery

`bool`

Enable point-in-time recovery on the replica (independent of the
source table's setting).

### spec.replicas[].deletionProtectionEnabled

`bool`

Refuse deletion of the replica while true.

### spec.replicas[].propagateTags

`bool`

Propagate the table's tags to the replica. One-way: tag changes
flow from the table to the replica; replica-side drift is left
alone.

### spec.replicas[].consistencyMode

`string`

"EVENTUAL" (the default -- classic global tables, last-writer-wins)
or "STRONG" (Multi-Region Strong Consistency -- synchronous quorum
writes; requires exactly two STRONG replicas, or one STRONG
replica plus the witness region).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["EVENTUAL","STRONG"]}}

### spec.globalTableWitness

`AwsDynamodbGlobalTableWitness`

The witness region of a Multi-Region Strong Consistency global
table -- it stores replicated writes to support quorum but serves
no reads or writes. Must accompany exactly one replica with
consistency_mode STRONG (the cheaper MRSC topology; the
alternative is two STRONG replicas and no witness).

### spec.globalTableWitness.regionName

`string` · required

The witness region. Must differ from the table's region and the
replica's region.

- rule: {"string":{"minLen":"1"}}

### spec.restoreSourceName

`string`

Create this table by restoring another table's point-in-time
state: the name of the source table in this account and region.
Key schema and attributes are inherited from the source. Mutually
exclusive with the other restore/import sources.

### spec.restoreSourceTableArn

`string`

Create this table by restoring from a source table ARN -- the
cross-region / cross-account form of point-in-time restore.
Server-side encryption must be configured on the restored table.
Mutually exclusive with the other restore/import sources.

### spec.restoreDateTime

`string`

The point in time to restore, as a UTC RFC3339 timestamp (e.g.
"2026-07-04T06:00:00Z"). Exactly one of restore_date_time or
restore_to_latest_time accompanies a point-in-time restore source.

### spec.restoreToLatestTime

`bool`

Restore the most recent recoverable state of the source table
instead of a specific timestamp.

### spec.restoreBackupArn

`string`

Create this table by restoring an on-demand or AWS Backup backup,
by ARN. Key schema and attributes are inherited from the backup.
Mutually exclusive with the other restore/import sources.

### spec.importTable

`AwsDynamodbImportTable`

Create this table pre-loaded with data imported from S3 (CSV,
DynamoDB JSON, or Amazon Ion) -- billed as a one-time import, far
cheaper than writing items individually. Mutually exclusive with
the restore sources; the key schema and attributes above are
required (imports define a brand-new table).

- rule: csv options only apply when input_format is CSV

### spec.importTable.s3Bucket

`string | valueFrom` · required

The S3 bucket holding the source data. Reference an AwsS3Bucket
bucket_id output or pass a literal bucket name.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.importTable.s3BucketOwner

`string`

The bucket owner's account ID -- set when importing from a bucket
in another account.

### spec.importTable.s3KeyPrefix

`string`

Only objects under this key prefix are imported. Empty imports the
whole bucket.

### spec.importTable.inputFormat

`string` · required

The source data format: "CSV", "DYNAMODB_JSON", or "ION".

- rule: {"required":true,"string":{"in":["CSV","DYNAMODB_JSON","ION"]}}

### spec.importTable.inputCompressionType

`string`

How the source objects are compressed: "GZIP", "ZSTD", or "NONE".
Empty keeps the AWS default (NONE).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["GZIP","ZSTD","NONE"]}}

### spec.importTable.csv

`AwsDynamodbImportTableCsv`

CSV parsing options; only meaningful when input_format is "CSV".

### spec.importTable.csv.delimiter

`string`

The field delimiter. Empty keeps the AWS default (",").

### spec.importTable.csv.headerList

`[]string`

Column names, when the files carry no header row. Empty treats the
first row of each file as the header.

## Validation Rules

- `key_schema_required_unless_restored`: key_schema and attribute_definitions are required -- omit them only when the table is created by restore (restore_source_name, restore_source_table_arn, or restore_backup_arn), which inherits them from the source
- `table_key_schema_shape`: key_schema must have exactly one HASH element and at most one RANGE element, HASH first
- `key_attributes_declared`: every key_schema element must name an attribute declared in attribute_definitions
- `provisioned_requires_throughput`: the effective billing mode is PROVISIONED (the AWS default when billing_mode is empty), which requires provisioned_throughput with read and write capacity > 0
- `on_demand_forbids_provisioned_throughput`: provisioned_throughput must stay unset on a PAY_PER_REQUEST table
- `on_demand_throughput_needs_on_demand`: on_demand_throughput ceilings only apply to PAY_PER_REQUEST tables
- `gsi_throughput_matches_billing`: each global secondary index must carry provisioned_throughput (read and write > 0) when the effective billing mode is PROVISIONED, and must not carry it on a PAY_PER_REQUEST table; per-index on_demand_throughput only applies to PAY_PER_REQUEST tables
- `index_key_attributes_declared`: every index key element (GSI key_schema, LSI range_key) must name an attribute declared in attribute_definitions
- `stream_view_requires_enabled`: stream_view_type is required when stream_enabled is true and must stay empty when streams are disabled
- `replicas_require_streams`: global-table replicas require stream_enabled with stream_view_type NEW_AND_OLD_IMAGES
- `replica_consistency_modes_agree`: all replicas must use the same consistency_mode -- mixing STRONG and EVENTUAL replicas is not a valid global-table topology
- `mrsc_topology`: Multi-Region Strong Consistency requires exactly two STRONG replicas, or exactly one STRONG replica plus global_table_witness
- `restore_import_sources_exclusive`: at most one create source may be set: restore_source_name, restore_source_table_arn, restore_backup_arn, or import_table
- `restore_point_requires_pitr_source`: restore_date_time / restore_to_latest_time require a point-in-time restore source (restore_source_name or restore_source_table_arn), and exactly one of the two must be chosen
- `insights_indexes_exist`: contributor_insights.gsi_index_names must name global secondary indexes defined on this table

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsDynamodb, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.table_name` | `string` | The table name -- what SDK calls, IAM policy resources, and application configuration reference. |
| `status.outputs.table_arn` | `string` | The table ARN -- the join key for IAM policies, resource policies, and cross-service integrations. |
| `status.outputs.table_id` | `string` | The provider-assigned table identifier. |
| `status.outputs.stream_arn` | `string` | The DynamoDB Streams ARN -- what Lambda event-source mappings and other stream consumers attach to. Empty when streams are disabled. |
| `status.outputs.stream_label` | `string` | The stream label (a per-stream timestamp qualifier); combined with the account and table name it uniquely identifies the stream. Empty when streams are disabled. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.serverSideEncryption.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.kinesisStreamingDestination.streamArn` | AwsKinesisStream | `status.outputs.stream_arn` |
| `spec.replicas[].kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.importTable.s3Bucket` | AwsS3Bucket | `status.outputs.bucket_id` |

## See Also

- [Overview](../README.md)
