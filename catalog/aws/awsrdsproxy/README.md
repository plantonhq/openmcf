# AwsRdsProxy

An RDS Proxy — the managed connection pool between connection-hungry applications (Lambda above all) and a database — with its pool tuning, additional endpoints, and database target managed in-line.

## Highlights

- **One kind, the whole proxy**: the proxy, its default target group's pool tuning (a PATCH satellite — it always exists on the proxy, delete is a no-op), name-keyed additional endpoints, and the instance-XOR-cluster target registration.
- **Secrets-based sign-in, IAM-optional client auth**: the proxy reads database credentials from Secrets Manager under your IAM role; clients connect with native credentials or IAM tokens per auth entry (`iam_auth`), enforced over TLS with `require_tls`.
- **Fixed-for-life dials called out**: engine family, the proxy's subnets, and both network types replace the proxy (its endpoint DNS changes with it) — everything else updates in place.

## Both Engines

Both modules render the proxy, pool, endpoints, and target identically and export the same outputs: `proxy_name` (import ID), `proxy_arn`, `endpoint` (the default DNS name applications connect to), `default_target_group_arn`/`_name`, the `endpoint_addresses`/`endpoint_arns` maps keyed by endpoint name, and the target's `target_type`/`target_rds_resource_id`.

## Chart Wiring

`role_arn` → AwsIamRole `role_arn`; `auth[].secret_arn` → AwsSecretsManagerSecret `secret_arn`; `vpc_subnet_ids` → AwsSubnet `subnet_id`; `vpc_security_group_ids` → AwsSecurityGroup `security_group_id`; `target.db_instance_identifier` → AwsRdsInstance `instance_identifier` (or `target.db_cluster_identifier` → AwsRdsCluster `cluster_identifier`). Applications take the `endpoint` output as their database host.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
