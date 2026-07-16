# AWS DynamoDB: Table Rebuild to the Full Provider Surface

**Date**: July 4, 2026
**Type**: Feature (breaking, zero users)
**Components**: API Definitions, AWS Provider, Terraform Modules, Pulumi Modules, Tofu Tfvars Converter, Testing Framework

## Summary

`AwsDynamodb` is rebuilt from its early spec to the full `aws_dynamodb_table` v6 surface: Global Tables v2 replication up to Multi-Region Strong Consistency, all three capacity shapes (on-demand with spend ceilings, provisioned, pre-warmed), multi-attribute GSI keys, restore and S3-import create sources, and the folded table-scoped governance satellites (resource policy, Kinesis change-data destination, contributor insights). Both engines were rebuilt to one contract with the kind's first-ever E2E coverage — live dual-engine lanes green with a zero-orphan sweep. A live run also surfaced and fixed a framework-level bug in the proto-to-tfvars converter affecting every Terraform module.

## Problem Statement / Motivation

The spec was ~15 fields short of what DynamoDB actually offers, and the kind had never been deployable on Terraform at all.

### Pain Points

- **The Terraform module never worked against its own proto.** Its locals compared enum values to prefixed strings (`"KEY_TYPE_HASH"`, `"BILLING_MODE_PROVISIONED"`) that protojson never emits, referenced a `spec.auto_scale` field that did not exist, carried the legacy `{key,value}` label contract the tfvars pipeline cannot satisfy, and pinned the provider at `= 5.82.0`. The hack manifest carried the same phantom enum forms — invalid against its own spec.
- **The Pulumi module deployed differently-named tables.** It never set the table's `Name` argument, so Pulumi auto-naming appended a random suffix — the same manifest produced `<name>-<hex>` on Pulumi but `<name>` on Terraform.
- No global tables, no on-demand ceilings, no warm throughput, no PITR recovery window, no restores or imports, no resource policy, no Kinesis destination, plain-string KMS ARN with no `AwsKmsKey` composition.
- Six nested proto enums diverged from the settled family convention (provider-authentic strings validated by CEL), forcing enum-translation locals in both engines.

## Solution / What's New

- **Spec rebuilt to the full v6 surface** (27 top-level fields, 16 spec-level CEL rules): the capacity trio with billing-mode coupling CEL-enforced on the table and every GSI; multi-attribute GSI `key_schema` (1-4 HASH + 0-4 RANGE, the provider's own validator mirrored); LSIs modeled honestly as `range_key` + projection (their partition key is always the table's); folded `replicas` + `global_table_witness` with the MRSC topology rule (exactly two STRONG replicas, or one plus the witness) and the streams requirement (`NEW_AND_OLD_IMAGES`) failing at validate, not deploy; mutually-exclusive create sources (PITR restore by name/ARN, backup restore, S3 import); key schema required-unless-restored (restores inherit it from the source).
- **All six enums converted to provider strings** (`PAY_PER_REQUEST`, `S`, `HASH`, ...) — manifests read like AWS documentation, the `KEYS_ONLY` name-collision hack disappears, and both engines pass values through untranslated.
- **Composition**: `server_side_encryption.kms_key_arn` and per-replica `kms_key_arn` → `AwsKmsKey.key_arn` refs; `kinesis_streaming_destination.stream_arn` → `AwsKinesisStream.stream_arn`; `import_table.s3_bucket` → `AwsS3Bucket.bucket_id` — three new FK edges.
- **Folded satellites materialized per-name in both engines**: the resource policy, the Kinesis destination (one per table, AWS's own rule), and contributor insights (table + per-GSI, with the provider's `mode` argument).
- **Both engines rebuilt to one contract**: generator-owned `variables.tf` under the drift guard; provider floor `>= 6.37.0` (witness 6.22, GSI key_schema 6.29, the 6.37 GSI-removal fix; resolves v6.53); naming basis `metadata.name` on both engines with the cloud `Name` argument set explicitly; zero PARITY-EXCEPTIONs (every field verified on pulumi-aws v7.35.0).
- **Framework fix (live catch): the tfvars HCL writer now escapes template introducers.** A manifest string carrying `${...}` — an IAM policy condition variable like `${aws:ResourceAccount}`, shell user-data — rendered into `terraform.tfvars` unescaped and failed `tofu init` with "Extra characters after interpolation expression". The writer now emits HCL's own escapes (`$${`, `%%{`) for string values and array elements, with a regression test; policy-bearing string fields across every kind are now safe.
- **First-ever E2E**: a state-aware DescribeTable verifier (DELETING = absent, the RDS lifecycle class), profile + full-surface on-demand scenario (GSI, TTL, streams, PITR, SSE, resource policy, per-index insights), dual-engine entrypoints. No prerequisites — a true leaf.

## Validation

- Offline gate all green: spec/CEL tests (10 happy paths + 33 error paths), outputs conformance, TF drift guard, validate-refs, secret-coverage, tofu init+validate on v6.53, release-equivalent Pulumi build, `make build-go`, Bazel build of all touched targets, all 5 manifests CLI-validated, mechanical field-parity sweep clean across both engines, site catalog regenerated.
- **Live dual-engine E2E green**: Pulumi 3m28s, Terraform 1m30s (deploy → DescribeTable verify → outputs verify → destroy → verify-clean). Zero-orphan sweep: no tables; the only e2e-tagged account remnants are INACTIVE ECS artifacts AWS keeps describable after deletion.
- The Terraform lane's first run failed on the tfvars template-introducer bug (the scenario's resource policy carries `${aws:ResourceAccount}`); fixed at the root in the converter and re-run green.

## Breaking Changes

Zero users; no migration. For the record: enum values became bare provider strings, `point_in_time_recovery_enabled` became the `point_in_time_recovery` message, `contributor_insights_enabled` became the `contributor_insights` message, LSIs take `range_key` instead of a key-schema list, and `server_side_encryption.kms_key_arn` became a `StringValueOrRef`.

## Impact

DynamoDB is the most-adopted AWS data service; the catalog now covers it at the depth an advanced organization actually reaches — multi-region active-active with strong consistency, spend guardrails on on-demand tables, pre-warmed launch capacity, cross-account restores, and bulk S3 seeding — with composition into the KMS/Kinesis/S3 graph instead of loose ARN strings. The converter fix removes a whole failure class for every Terraform module that carries user policy or template content.
