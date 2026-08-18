# AWS Route 53 Resolver DNS Firewall

Block malware and phishing domains, alert on DNS tunneling, and sinkhole bad lookups — at the DNS layer, before a connection ever opens. One rule group carries the policy; associating it to VPCs turns it on.

## What Gets Managed

- The rule group and its owned domain lists (name-keyed, with their domains).
- Filtering rules in priority order: match an owned list, an external/AWS-managed list by ID, or a DNS threat class (DGA, dictionary DGA, DNS tunneling) with a confidence threshold; act with ALLOW, ALERT, or BLOCK — and shape BLOCK responses (NODATA, NXDOMAIN, or a CNAME OVERRIDE sinkhole with TTL).
- VPC associations with evaluation priority and mutation protection.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Route 53 Resolver permissions.

### AWS Prerequisites

- The VPCs to filter (nothing else — DNS Firewall needs no endpoints or agents).

## After You Deploy

- Queries from associated VPCs evaluate against the rules in ascending priority; the first match wins. Blocked lookups answer per the rule's block response.
- Pair with AwsRoute53ResolverQueryLog to see what fired — the firewall itself logs nothing.

## Common Changes

- Add domains to a list (in-place push), add rules (leave priority gaps — 100, 200 — for later inserts), associate more VPCs.
- A rule's match source is fixed for life (list-matched vs threat-matched replaces the rule); actions, priorities, and block responses update in place.
- Never enable mutation_protection on an association you intend to destroy declaratively — it refuses its own deletion until disabled.
