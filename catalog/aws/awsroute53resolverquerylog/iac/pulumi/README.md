# AwsRoute53ResolverQueryLog — Pulumi module

Manages one resolver query logging configuration (`route53.ResolverQueryLogConfig`) with its VPC associations (`route53.ResolverQueryLogConfigAssociation`), one per spec entry.

Module facts worth knowing before editing:

- **The configuration is immutable except tags** — Name and DestinationArn are both ForceNew.
- **Associations are pure joins** (config, vpc) — every argument ForceNew, no update path; resource names key on the resolved VPC id, matching the Terraform module's for_each keys.
- **The association's failure is asynchronous** — an unwritable destination FAILS the association after a clean apply; the E2E verifier asserts ACTIVE on every association.

Outputs mirror the Terraform module key-for-key: `query_log_config_id` (import ID), `query_log_config_arn`, `share_status`, `association_ids`.
