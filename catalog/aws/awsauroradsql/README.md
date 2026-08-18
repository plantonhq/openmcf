# AwsAuroraDsql

An Aurora DSQL cluster — AWS's serverless, PostgreSQL-compatible distributed SQL database with no instances, no capacity dials, and active-active multi-region writes — with its multi-region pairing managed in-line.

## Highlights

- **Nothing to size**: no instances, no ACUs, no storage allocation — DSQL scales to zero when idle and bills per request and per byte stored. A single-region cluster deploys from defaults alone.
- **The endpoint is derived honestly**: AWS exposes no endpoint attribute, so both modules derive the documented `{identifier}.dsql.{region}.on.aws` connection host as the `endpoint` output — the chart-ready database host.
- **Multi-region is a one-shot pairing, taught up front**: peering happens while a fresh cluster is still in PENDING_SETUP (the modules order it correctly); the peering has no update path and a no-op delete — changing peers means recreating the cluster.

## Both Engines

Both modules render the cluster and its pairing identically and export the same outputs: `identifier` (import ID), `cluster_arn` (what a peer's multi_region references), `endpoint`, `vpc_endpoint_service_name` (PrivateLink), and `encryption_type`.

## Chart Wiring

`kms_encryption_key` → AwsKmsKey `key_arn`; `multi_region.peer_cluster_arns` → another AwsAuroraDsql's `cluster_arn`. Applications take the `endpoint` output as their PostgreSQL host and sign in with IAM auth tokens — DSQL has no native passwords.
