# Azure DNS Resolver Virtual Network Link -- Operational Guide

Judgment calls that matter when you run ruleset links in production.

## Link the hub's own network -- it is never implicit

The resolver's home network does not inherit its own ruleset: it needs a link like every spoke. Forgetting it produces the classic asymmetry -- spokes resolve on-premises names, the hub does not. Make the hub link part of the ruleset's own deployment story.

## One link per pair, named after the network

Azure allows exactly one link between a given ruleset and a given network; a duplicate create fails as already-exists. Name each link after the network it attaches ("spoke-payments") so the ruleset's link list reads as its audience roster, and use the metadata map to record the network's owner -- when a spoke team decommissions their network, the link to clean up is self-identifying.

## Replacement means a resolution blink for that network

Everything except metadata replaces the link, and while the link is gone the network stops forwarding -- captured domains fall back to Azure-internal resolution (usually NXDOMAIN) until the new link lands. Sequence link replacements outside change-sensitive windows for workloads that depend on on-premises names.

## The region wall is the ruleset's, not the resolver's

A link's network must live in the RULESET's region. Multi-region estates need a resolver + ruleset pair per region -- do not try to link a west VNet to the east rule book; build the west stack and copy the rules instead.

## Linking is the audit surface -- treat it like one

Which networks forward through which rule book is exactly what a security review asks. Because each link is a standalone resource, the answer is the resource list itself: keep links in the same repo/chart as the network they attach (spoke teams own their links), not centrally bundled with the ruleset -- that ownership split is why this attachment is its own component.
