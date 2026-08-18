# AwsRoute53ResolverFirewall

A Route 53 Resolver DNS Firewall rule group — the DNS-layer block/allow policy for queries leaving VPCs — with its domain lists, filtering rules, and VPC associations managed in-line. Rules match owned domain lists, external/AWS-managed lists by ID, or DNS threat classes (Advanced protection), and apply ALLOW/ALERT/BLOCK.

## Highlights

- **The rule group is the pivot**: rules cannot exist without one, and associating the group to a VPC is what turns the policy on there — one kind, the whole DNS-firewall story.
- **The provider's plan-time union is CELs**: each rule matches exactly one of an owned list, an external list ID, or a threat class (paired with its confidence threshold); BLOCK-response tuning gates on the BLOCK action; OVERRIDE demands its record and TTL.
- **The one-value knob is module-owned**: `block_override_dns_type` has exactly one legal value (CNAME) — the modules pin it whenever the response is OVERRIDE, so no dead knob ships.
- **Per-VPC fail-open deliberately lives elsewhere**: it is a VPC setting (delete merely resets it), recorded with the VPC's other resolver settings — a rule group associable to many VPCs cannot own per-VPC state.

## Both Engines

Both modules render the group, lists, rules, and associations identically and export the same outputs: `rule_group_id` (import ID), `rule_group_arn`, `share_status`, plus the `domain_list_ids`, `association_ids`, and `rule_match_ids` maps keyed like the spec entries.

## Chart Wiring

`vpc_associations.vpc_id` → AwsVpc `vpc_id`. Rules reference owned lists by name (uniqueness and existence are CELs) or any external list by literal ID.
