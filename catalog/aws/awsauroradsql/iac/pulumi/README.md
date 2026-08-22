# AwsAuroraDsql — Pulumi module

Manages one Aurora DSQL cluster (`dsql.Cluster`) with its multi-region pairing (`dsql.ClusterPeering`).

Module facts worth knowing before editing:

- **The peering is a disguised UpdateCluster call** AWS accepts only while the cluster is PENDING_SETUP — created right after the cluster with an explicit DependsOn, the one valid order.
- **The peering has NO update path** (changes error at apply — the provider declares no replace triggers) and a no-op delete; peer changes mean recreating the cluster.
- **The endpoint output is module-derived** — AWS exposes no endpoint attribute; the documented `{identifier}.dsql.{region}.on.aws` shape is composed here (a recorded parity exclusion, identical in Terraform).
- **Only the witness region replaces the cluster**; deletion protection and the KMS key update in place.

Outputs mirror the Terraform module key-for-key: `identifier` (import ID), `cluster_arn`, `endpoint`, `vpc_endpoint_service_name`, `encryption_type`.
