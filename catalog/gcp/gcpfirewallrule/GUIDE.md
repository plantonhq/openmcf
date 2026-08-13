# GcpFirewallRule Guide

The judgment this guide protects: firewall rules are the network's
policy layer, and their failure mode is silence — an over-broad rule
never pages anyone, and a deleted rule cuts traffic that three other
teams assumed was permanent. Write rules that say who, not just what.

## Target by identity, not by everything

A rule with no targets applies to EVERY instance in the network — almost
never the intent. Service-account targeting
(`targetServiceAccounts`) is the strongest form: it follows the workload
identity, survives re-tagging, and cannot be forged from inside a VM the
way tags can. Tag targeting is fine for coarse tiers (web, bastion).
The two systems cannot mix in one rule (GCP's constraint, spec-enforced)
— pick per rule.

## DENY wins ties; use priority deliberately

At equal priority DENY beats ALLOW; lower numbers win otherwise. A sane
layout: broad hygiene DENYs at high priority numbers (65000s), workload
ALLOWs in the middle (1000s), and emergency blocks at low numbers
reserved for incident response. `disabled: true` is the reversible way
to test removing a rule — flip it, watch, then delete.

## Egress rules are where data leaves

INGRESS rules guard the front door; EGRESS rules with
`destinationRanges` are the exfiltration control almost nobody writes.
A default-deny egress rule plus explicit allows to known ranges is the
posture; remember an EGRESS rule with no destination defaults to
0.0.0.0/0.

## Log before you tighten

`logConfig` turns on per-connection logs for the rule — enable it
(EXCLUDE_ALL_METADATA keeps volume sane) on any rule you are about to
tighten, watch a week of real traffic, then cut. Logging every rule on a
busy network is a real Logging bill; logging the rule under change is
cheap insurance.

## Tags at create time only

`resourceManagerTags` binds org tag values for IAM conditions and
org-policy scoping — and the provider sends them through a create-only
params block, so changing them REPLACES the rule (a brief enforcement
gap on the replaced rule). Plan tag changes with the same care as a rule
rewrite.

## Teardown discipline

`DELETE` cuts matched traffic over to the next matching rule the moment
it lands — on a permissive network that can mean "over to nothing".
`PREVENT` suits shared-plumbing rules (health-check allows, bastion
ingress) other teams silently depend on. `ABANDON` keeps the rule
enforcing while dropping management — the handoff path when a network
team takes ownership.
