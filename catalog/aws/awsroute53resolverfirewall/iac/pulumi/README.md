# AwsRoute53ResolverFirewall — Pulumi module

Manages one DNS Firewall rule group (`route53.ResolverFirewallRuleGroup`) with its domain lists (`route53.ResolverFirewallDomainList`), rules (`route53.ResolverFirewallRule`), and VPC associations (`route53.ResolverFirewallRuleGroupAssociation`), one resource per spec entry.

Module facts worth knowing before editing:

- **A rule's match source resolves in the module**: an owned list name becomes that list resource's ID output; an external list ID passes through; the threat arm sends DnsThreatProtection + ConfidenceThreshold. All ForceNew at the provider.
- **BlockOverrideDnsType is module-pinned to CNAME** whenever the response is OVERRIDE — the provider's one-value vocabulary, never surfaced as a spec knob.
- **`rule_match_ids` exports each rule's match identity** (domain list ID or the computed threat-protection ID) keyed by rule name — the second half of the rule's composite import ID.
- **Associations update priority/mutation_protection in place**; ENABLED mutation protection refuses the association's own destroy.

Outputs mirror the Terraform module key-for-key: `rule_group_id` (import ID), `rule_group_arn`, `share_status`, `domain_list_ids`, `association_ids`, `rule_match_ids`.
