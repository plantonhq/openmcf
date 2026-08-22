# AwsRdsProxy — Terraform/OpenTofu module

Manages one RDS Proxy (`aws_db_proxy`) with its default target group (`aws_db_proxy_default_target_group`), additional endpoints (`aws_db_proxy_endpoint`, keyed by name), and database target (`aws_db_proxy_target`).

Module facts worth knowing before editing:

- **The default target group is always managed** — it exists on every proxy anyway (its provider delete is a no-op), managing it tunes the pool and anchors the target's ordering.
- **The auth scheme is pinned** — SECRETS is the only value AWS supports; the module never exposes it.
- **The target waits out CREATING databases** — the provider retries registration for 5 minutes against an instance still coming up.
- **ForceNew edges**: engine_family, vpc_subnet_ids, and both network types on the proxy; everything on the endpoints except security groups; everything on the target.

Outputs mirror the Pulumi module key-for-key: `proxy_name` (import ID), `proxy_arn`, `endpoint`, `default_target_group_arn`/`_name`, `endpoint_addresses`, `endpoint_arns`, `target_type`, `target_rds_resource_id`.
