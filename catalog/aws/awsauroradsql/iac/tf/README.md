# AwsAuroraDsql — Terraform/OpenTofu module

Manages one Aurora DSQL cluster (`aws_dsql_cluster`) with its multi-region pairing (`aws_dsql_cluster_peering`).

Module facts worth knowing before editing:

- **The peering is a disguised UpdateCluster call** AWS accepts only while the cluster is PENDING_SETUP — created right after the cluster, the one valid order.
- **The peering has NO update path** (changes error at apply — the provider declares no replace triggers) and a no-op delete; peer changes mean recreating the cluster.
- **The endpoint output is module-derived** — AWS exposes no endpoint attribute; the documented `{identifier}.dsql.{region}.on.aws` shape is composed here (a recorded parity exclusion, identical in Pulumi).
- **Only the witness region replaces the cluster**; deletion protection and the KMS key update in place.

Outputs mirror the Pulumi module key-for-key: `identifier` (import ID), `cluster_arn`, `endpoint`, `vpc_endpoint_service_name`, `encryption_type`.
