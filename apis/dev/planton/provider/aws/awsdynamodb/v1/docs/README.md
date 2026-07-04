# AWS DynamoDB — Architecture and Design

## Overview

`AwsDynamodb` models one Amazon DynamoDB table at the full provider
surface: key schema and indexes, all three capacity shapes (on-demand
with spend ceilings, provisioned, pre-warmed), streams, Global Tables v2
replication up to Multi-Region Strong Consistency, encryption at rest,
point-in-time recovery, create-from sources (restores and S3 import),
and the table-scoped governance surface (resource policy, Kinesis
change-data destination, CloudWatch contributor insights).

A table is a true leaf in the resource graph -- nothing has to exist
before it -- and everything it composes with attaches by reference: the
KMS keys that encrypt the table and each replica (`AwsKmsKey`), the
Kinesis Data Stream that receives change data (`AwsKinesisStream`), and
the S3 bucket an import seeds it from (`AwsS3Bucket`).

## Design Decisions

- **Provider strings, not enums.** Every categorical field carries the
  exact string the AWS API accepts (`PAY_PER_REQUEST`, `S`, `HASH`,
  `KEYS_ONLY`, ...), validated by CEL sets. Manifests read like AWS
  documentation and both engines pass values through untranslated.
- **Key schema honesty.** The table's primary key is a `key_schema`
  list (one `HASH`, optional `RANGE` -- create-time immutable), and
  every key element must name a declared attribute; validation enforces
  both before anything deploys. A local secondary index is modeled as
  `range_key` + projection because its partition key is always the
  table's own -- there is no independent choice to express.
- **Restore-created tables relax the schema requirement.** Key schema
  and attribute definitions are required unless the table is created by
  restore (`restore_source_name`, `restore_source_table_arn`,
  `restore_backup_arn`), where AWS inherits them from the source --
  the create-time-derivation pattern.
- **Global Tables v2 folded as `replicas`.** Each replica is a setting
  of exactly one table (the provider models it as a nested block),
  referenced by nothing else -- folding is honest composition. The
  MRSC topology rule (exactly two STRONG replicas, or one STRONG
  replica plus the witness region) and the streams requirement
  (NEW_AND_OLD_IMAGES) are enforced at validation time, so an invalid
  topology never reaches AWS.
- **Table-scoped satellites folded, materialized per-name.** The
  resource policy, the Kinesis streaming destination, and contributor
  insights are standalone provider resources, but each is keyed by the
  table with replace-on-change semantics and owned by exactly one
  table -- a table setting, not a graph node. Both engines materialize
  each as its own provider resource, so edits are in place.
- **Capacity modeled in all three AWS shapes.** On-demand ceilings
  (`on_demand_throughput`), reserved units (`provisioned_throughput`,
  table and per-GSI), and pre-warmed floors (`warm_throughput`,
  increase-only -- lowering it replaces the table, stated honestly).
  Billing-mode coupling is CEL-enforced on the table and every GSI.
- **Multi-attribute GSI keys.** GSI `key_schema` accepts 1-4 HASH
  elements followed by up to 4 RANGE elements -- the Multi-Attribute
  Keys design pattern -- with the provider's own ordering and count
  rules as validation.

## Deliberately Skipped Provider Surface

| Provider surface | Verdict | Reason |
|---|---|---|
| `aws_dynamodb_table_item` | DEFER | Data-plane seed items, not infrastructure shape; `import_table` covers bulk seeding. Revisit on concrete pull. |
| `aws_dynamodb_table_export` | DEFER | A one-shot operational export job, not table shape. |
| `aws_dynamodb_global_table` (v1, 2017.11.29) | SKIP | Superseded by Global Tables v2, which the folded `replicas` blocks model. |
| `aws_dynamodb_table_replica` | SKIP | An alternate per-replica resource for split-provider-config setups; the folded `replicas` blocks cover the composable case. |
| `aws_dynamodb_global_secondary_index` (standalone) | SKIP | An alternate way to manage a GSI outside its table; the folded `global_secondary_indexes` list is the honest single-owner shape. |
| `aws_dynamodb_tag` | SKIP | Per-tag glue; identity tags derive from metadata, and `replicas[].propagate_tags` covers replica tagging. |
| Stream-scoped resource policies | SKIP | `resource_policy` targets the table; policy on the stream ARN is a niche surface with the same JSON escape hatch available later. |
| DAX (`aws_dax_*`) | DEFER | A separate caching product with its own cluster/parameter/subnet model -- a candidate kind of its own on demand. |

## Billing Note

An on-demand table with zero traffic costs storage only; the E2E
scenario writes no items, so a full create-verify-destroy lane accrues
effectively nothing. Warm throughput and provisioned capacity bill from
the moment they are set -- presets default to on-demand.

## References

- [DynamoDB Developer Guide](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html)
- [Global Tables version 2019.11.21](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/globaltables.V2.html)
- [Multi-Region Strong Consistency](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/multi-region-strong-consistency-gt.html)
- [Multi-Attribute Keys design pattern](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/GSI.DesignPattern.MultiAttributeKeys.html)
- [Warm throughput](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/warm-throughput.html)
