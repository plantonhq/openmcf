# AwsNetworkAcl — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Stateless bites twice

Every connection needs rules in BOTH directions: inbound 443 without an outbound 1024–65535 allow drops every response, and the failure reads like a hung service, not a firewall. When a NACL'd subnet "can't reach anything", check the reply direction first.

## Deny order is everything

First match wins by `rule_no`. A deny numbered above (higher than) the allow it should beat never fires. Put denies low (90 before the 100-series allows), leave gaps for insertions, and remember AWS's catch-all deny (32767) backstops everything unmatched.

## Association is replacement

A subnet always belongs to exactly one ACL. Listing it here replaces its previous association atomically; removing it hands it to the VPC's DEFAULT NACL (allow-all by AWS default) — removal can silently open a subnet up, not lock it down.

## Protocols are numbers underneath

AWS stores protocol numbers; the provider normalizes names ("tcp" → 6) when diffing, so both spellings are stable. On `-1`/`all` rules AWS stores no ports — the spec's CEL keeps them unset so plans stay clean.

## NACLs complement, not replace, security groups

The ACL is the subnet's coarse, stateless, deny-capable screen; the security group is the instance's fine, stateful allow-list. Tier boundaries and CIDR blocks belong here; service-to-service allows belong in security groups.
