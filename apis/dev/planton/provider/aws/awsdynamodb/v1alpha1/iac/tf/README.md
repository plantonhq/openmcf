# Terraform Module to Deploy AwsDynamodb

This module provisions an Amazon DynamoDB table -- key schema and
indexes, capacity in any of AWS's three shapes, streams, Global Tables
v2 replication, encryption, recovery, and the table-scoped governance
satellites (resource policy, Kinesis change-data destination,
contributor insights) -- aligned with the Planton API. The KMS key, the
Kinesis stream, and the S3 import bucket all attach by reference; the
module never creates a resource that deserves to be its own node.

## CLI (local backend)

```shell
planton tofu init --manifest ../hack/manifest.yaml
planton tofu plan --manifest ../hack/manifest.yaml
planton tofu apply --manifest ../hack/manifest.yaml --auto-approve
planton tofu destroy --manifest ../hack/manifest.yaml --auto-approve
```

Credentials are passed via the stack input through the CLI, not in `spec`.

## Files

- `variables.tf` (generated; do not edit)
- `provider.tf` — provider setup (`hashicorp/aws >= 6.37.0`; the
  global-table witness, multi-attribute GSI keys, and the GSI
  key_schema removal fix ride the v6 line)
- `locals.tf` — naming basis, identity tags, key-schema lowering
- `table.tf` — the `aws_dynamodb_table` resource
- `satellites.tf` — resource policy, Kinesis streaming destination, and
  contributor insights (table + per-GSI)
- `outputs.tf` — outputs matching `AwsDynamodbStackOutputs`

## Outputs

`table_name`, `table_arn`, `table_id`, `stream_arn`, `stream_label` --
the stream pair is empty when streams are disabled, matching the Pulumi
module exactly.
