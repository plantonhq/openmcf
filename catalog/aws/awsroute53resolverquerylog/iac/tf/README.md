# AwsRoute53ResolverQueryLog — Terraform/OpenTofu module

Manages one resolver query logging configuration (`aws_route53_resolver_query_log_config`) with its VPC associations (`aws_route53_resolver_query_log_config_association`, keyed by the resolved VPC id).

Module facts worth knowing before editing:

- **The configuration is immutable except tags** — name and destination_arn are both ForceNew; the provider's update path is tags-only.
- **Associations are pure joins** (config, vpc) — every argument ForceNew, no update path.
- **The association's failure is asynchronous** — the provider's waiter surfaces the association's error code when the resolver cannot write to the destination; the E2E verifier additionally asserts ACTIVE on every association.

Outputs mirror the Pulumi module key-for-key: `query_log_config_id` (import ID), `query_log_config_arn`, `share_status`, `association_ids`.
