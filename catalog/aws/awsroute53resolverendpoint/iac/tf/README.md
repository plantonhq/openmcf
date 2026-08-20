# AwsRoute53ResolverEndpoint — Terraform/OpenTofu module

Manages one resolver endpoint (`aws_route53_resolver_endpoint`) with its forwarding rules (`aws_route53_resolver_rule`, keyed by rule name) and their VPC associations (`aws_route53_resolver_rule_association`, keyed `rule//vpc`).

Module facts worth knowing before editing:

- **Direction and security groups are ForceNew** on the endpoint; ip_addresses churn in place (the provider adds before it removes, one waiter round-trip per change — large IP churns are slow by design).
- **SYSTEM rules carry no endpoint binding** — the module sends `resolver_endpoint_id` only for FORWARD/DELEGATE rules; detaching a rule's endpoint (id → empty) is the provider's one forced replacement on the rule.
- **Associations are pure joins** — every argument ForceNew, no update path; their cosmetic name is derived deterministically (`{rule}-{vpc}`, capped at 64) so both engines send the identical value instead of Pulumi auto-naming.
- **Trailing dots on rule domains normalize both ways** (the provider's StateFunc) — no drift from `corp.example.com.`.
- **RAM-shared rule tags never read back** (the provider skips them for SHARED_WITH_ME rules).

Outputs mirror the Pulumi module key-for-key: `endpoint_id` (import ID), `endpoint_arn`, `host_vpc_id`, `ip_addresses`, `rule_ids`, `rule_association_ids`.
