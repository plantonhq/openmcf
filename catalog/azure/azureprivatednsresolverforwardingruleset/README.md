# Overview

The **AzurePrivateDnsResolverForwardingRuleset** component deploys a DNS forwarding ruleset -- the rule book that decides which DNS names leave Azure and which servers answer them. Each rule pairs a domain ("corp.contoso.com.") with the on-premises (or other custom) DNS servers that own it; everything else keeps resolving inside Azure.

## Purpose

- **Conditional forwarding as configuration**: the corp-domain-to-datacenter routing that used to live in forwarder VM config files becomes declarative rules on a managed object.
- **One rule book, many networks**: rules are written once on the ruleset; every virtual network linked to it (AzurePrivateDnsResolverVirtualNetworkLink) inherits them.
- **Steers the resolver's outbound pipe**: the ruleset binds a DNS Private Resolver's outbound endpoint -- that is the path captured queries take out of Azure.

## Key Features

- Full azurerm v5 surface: the ruleset plus its forwarding rules (domain, ordered target servers with ports, enabled flag, per-rule metadata), composed as one component.
- Rules are keyed by name -- adding, removing, or editing one rule never touches its siblings; everything but a rule's domain updates in place.
- Chart-ready: references the resolver's outbound endpoint by typed output (the reference is the deploy-order edge); publishes the ruleset id that links bind.

## Use Cases

- **Active Directory resolution from Azure**: forward "corp.contoso.com." to the domain controllers over VPN/ExpressRoute.
- **Multi-domain estates**: one ruleset carrying rules for each acquired company's namespace, each pointing at its own DNS servers.
- **Staged cutovers**: park a rule with `enabled: false` and flip it on when the tunnel goes live -- no redeploys.

## Future Enhancements

- The at-most-2-outbound-endpoints (same resolver) service rule stays documentation until cross-resource facts can be introspected offline.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
