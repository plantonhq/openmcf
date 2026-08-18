# AwsRdsProxy — Pulumi module

Manages one RDS Proxy (`rds.Proxy`) with its default target group (`rds.ProxyDefaultTargetGroup`), additional endpoints (`rds.ProxyEndpoint`), and database target (`rds.ProxyTarget`).

Module facts worth knowing before editing:

- **The default target group is always managed** — it exists on every proxy anyway (its provider delete is a no-op), managing it tunes the pool and anchors the target's ordering.
- **The auth scheme is pinned** — SECRETS is the only value AWS supports; the module never exposes it.
- **The target waits out CREATING databases** — the provider retries registration for 5 minutes against an instance still coming up.
- **Resource names key on endpoint names**, matching the Terraform module's for_each keys.

Outputs mirror the Terraform module key-for-key: `proxy_name` (import ID), `proxy_arn`, `endpoint`, `default_target_group_arn`/`_name`, `endpoint_addresses`, `endpoint_arns`, `target_type`, `target_rds_resource_id`.
