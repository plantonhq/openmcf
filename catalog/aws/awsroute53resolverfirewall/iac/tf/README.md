# AwsRoute53ResolverFirewall — Terraform/OpenTofu module

Manages one DNS Firewall rule group (`aws_route53_resolver_firewall_rule_group`) with its domain lists (`aws_route53_resolver_firewall_domain_list`, keyed by list name), rules (`aws_route53_resolver_firewall_rule`, keyed by rule name), and VPC associations (`aws_route53_resolver_firewall_rule_group_association`, keyed by association name).

Module facts worth knowing before editing:

- **The group is a name-and-tags container** — name ForceNew, update path tags-only.
- **A rule's match source resolves in the module**: an owned list name becomes that list resource's generated ID; an external list ID passes through; the threat arm sends `dns_threat_protection` + `confidence_threshold`. All three are ForceNew at the provider.
- **`block_override_dns_type` is module-pinned to CNAME** whenever the response is OVERRIDE — the provider's one-value vocabulary, never surfaced as a spec knob.
- **Domain list contents push post-create** (ADD on create, REPLACE on update) — a partially failed import is a retry error, not silent success.
- **Associations update priority/mutation_protection in place**; an ENABLED mutation protection refuses the association's own destroy.

Outputs mirror the Pulumi module key-for-key: `rule_group_id` (import ID), `rule_group_arn`, `share_status`, `domain_list_ids`, `association_ids`, `rule_match_ids`.
