# AwsRoute53ResolverEndpoint — Pulumi module

Manages one resolver endpoint (`route53.ResolverEndpoint`) with its forwarding rules (`route53.ResolverRule`, one per spec entry) and their VPC associations (`route53.ResolverRuleAssociation`, one per rule-VPC pair).

Module facts worth knowing before editing:

- **Direction and SecurityGroupIds are ForceNew** on the endpoint; IpAddresses churn in place (added before removed, so the two-address floor never breaks).
- **SYSTEM rules carry no endpoint binding** — ResolverEndpointId is set only for FORWARD/DELEGATE rules; detaching a rule's endpoint forces the provider to replace the rule.
- **Associations are pure joins** with a deterministic cosmetic name (`{rule}-{vpc}`, capped at the provider's 64-character wall) — explicit on purpose, so the engines never diverge on Pulumi's URN auto-naming.
- **The tri-state metrics toggles** send only when the spec states them (nil leaves AWS's default; an explicit false is sent as false).

Outputs mirror the Terraform module key-for-key: `endpoint_id` (import ID), `endpoint_arn`, `host_vpc_id`, `ip_addresses`, `rule_ids`, `rule_association_ids`.
