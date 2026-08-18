# AwsRoute53ResolverEndpoint

A Route 53 Resolver endpoint — the hybrid-DNS bridge between a VPC and networks outside it — with its forwarding rules and their VPC associations managed in-line. INBOUND endpoints receive queries from on-prem resolvers; OUTBOUND endpoints forward VPC queries out through the rules.

## Highlights

- **Rules ride the endpoint**: FORWARD rules bind to the owning endpoint and carry target name servers; SYSTEM rules punch recursive-resolution holes in broader forwards; each rule associates to the VPCs that honor it — one kind, the whole forwarding story.
- **The per-type walls are CELs**: FORWARD requires targets and an OUTBOUND endpoint, SYSTEM forbids targets, each target is IPv4 XOR IPv6 — AWS's server-side rejections front-loaded to validate-manifest.
- **Contracts taught in place**: detaching a rule's endpoint forces replacement (the provider's own CustomizeDiff), RECURSIVE rules are Resolver-owned (excluded from the vocabulary), RAM-shared rule tags never read back, and the metrics toggles are tri-state (unset leaves AWS's default; once set, revert by explicit false).

## Both Engines

Both modules render the endpoint, its rules, and their associations identically and export the same outputs: `endpoint_id` (import ID), `endpoint_arn`, `host_vpc_id`, `ip_addresses`, plus the `rule_ids` and `rule_association_ids` maps keyed like the spec entries.

## Chart Wiring

`ip_addresses.subnet_id` → AwsSubnet `subnet_id`; `security_group_ids` → AwsSecurityGroup `security_group_id`; `rules.vpc_ids` → AwsVpc `vpc_id`. The inbound endpoint's `ip_addresses` output is what on-prem resolvers target.
