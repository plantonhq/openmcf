# Azure DNS Forwarding Ruleset -- Operational Guide

Judgment calls that matter when you run forwarding rulesets in production.

## Write domains with the trailing dot, always

ARM stores rule domains as fully qualified names: `corp.contoso.com.` -- WITH the trailing dot. A rule captures the domain and everything under it; the most specific rule wins when domains nest. Write every domain in the spec exactly that way and spot-check the first deploy's plan.

## The ruleset steers nothing until networks are linked

Rules on an unlinked ruleset are inert configuration. Every virtual network that should forward through them needs its own AzurePrivateDnsResolverVirtualNetworkLink -- including the resolver's OWN network, which is not linked implicitly. Forgetting the hub's own link is the classic "works from spokes, fails from the hub" mystery.

## Domain edits replace the rule -- plan for the blink

Everything on a rule updates in place except `domain_name`: changing it deletes and recreates that rule, a brief window where the domain resolves inside Azure instead of forwarding. For a rename, ADD the new rule first, then remove the old one in a second apply -- rules are keyed by name, so both operations leave siblings untouched.

## Park rules with enabled: false, not by deleting them

A disabled rule keeps its configuration but forwards nothing (the domain resolves normally inside Azure). Staging a migration, testing a tunnel, or backing out a bad target list is a one-field flip rather than a delete-and-retype. Rule metadata is the right place to record WHY a rule is parked and who owns it.

## Target servers are tried in order -- put the healthy one first

Azure walks a rule's target list in order (up to 6 per rule). Order the primary datacenter's servers first and remember the targets must be reachable FROM THE OUTBOUND ENDPOINT'S SUBNET -- over VPN or ExpressRoute for on-premises targets. A rule whose targets are unreachable fails queries slowly rather than loudly; test resolution from a linked network after every target change.

## Do not forward what Azure already answers

A rule capturing a domain that also exists as a linked private DNS zone hijacks those queries to the on-premises servers -- which usually cannot answer for private endpoints. Keep rulesets to genuinely external namespaces and let private zones own theirs. (Azure refuses rules for its own service domains outright.)
