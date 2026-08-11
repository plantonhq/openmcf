# Overview

The **AzurePrivateDnsResolverVirtualNetworkLink** component deploys the attachment that makes a DNS forwarding ruleset take effect in ONE virtual network. A ruleset without links steers nobody; each link adds one network to its audience -- and networks join and leave independently, without touching the ruleset or each other.

## Purpose

- **Per-network opt-in**: which networks forward through the rule book is an explicit, auditable attachment -- not an implicit blast radius.
- **Spokes without plumbing**: a linked network needs no resolver, no endpoints, and no peering to the resolver's network for DNS forwarding to work.
- **Independent lifecycles**: hundreds of spokes (up to 500 per ruleset) each own their link -- exactly the shape of the private DNS zone family's virtual network link.

## Key Features

- Full azurerm v5 surface: the link plus its free-form metadata map (ARM gives links metadata, not tags), modeled one-to-one.
- Chart-ready: references the ruleset and the network by typed outputs -- the references are the deploy-order edges; cross-subscription networks are supported.
- One link per ruleset-network pair; edits to metadata are in place, everything else replaces just this link.

## Use Cases

- **Hub-and-spoke DNS**: link every spoke to the hub's ruleset so corp domains resolve from anywhere -- including the hub's own network, which is never linked implicitly.
- **Selective rollout**: attach networks to a new rule book one at a time, watching resolution behavior per environment.
- **Cross-subscription estates**: link networks that live in other subscriptions to a centrally-owned ruleset.

## Future Enhancements

- The same-region service rule (a link's network must be in the ruleset's region) stays documentation until cross-resource facts can be introspected offline.
