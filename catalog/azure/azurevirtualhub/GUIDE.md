# Azure Virtual Hub -- Operational Guide

Judgment that does not fit in field references.

## The hub bills from creation, and it is slow

A Standard hub bills hourly (plus per-unit router capacity) the
moment ARM accepts the create -- before anything is connected to it.
Verified figures live in the generated estimate at
`catalog/_pricing/estimates/azurevirtualhub.yaml`, computed from the
pinned, source-dated price book.
Provisioning takes 15-30 minutes: ARM builds a managed router and
brings its routing state to Provisioned (measured live: ~17-20 minutes
for a plain hub). Composed children stretch the deploy further -- a
route map alone ran ~17 minutes live, because ARM re-polls the router
state around each child. Deletion runs 11-15 minutes. Plan maintenance
windows around that, and never create hubs speculatively.

## Size the address prefix once, correctly

The hub's router, gateways, and firewall all draw addresses from the
`addressPrefix`, it must not overlap ANY connected VNet or branch, and
it is fixed at creation -- resizing means replacing the hub and
everything attached to it. Microsoft's minimum is /24; use the
recommended /23 unless address space is genuinely scarce.

## Routing intent and custom routing do not mix

Routing intent takes over the hub's routing policy wholesale: ARM
rejects per-connection route-table customization while an intent is
active. Choose the model per hub -- EITHER label/table-based routing
(isolation, service chaining) OR intent-based firewall steering. The
intent's next hop must be an Azure Firewall deployed in THIS hub;
referencing one anywhere else fails at deploy.

## Custom route tables express isolation through connections

The hub's tables are only half the story: a spoke is isolated by what
its CONNECTION associates with and propagates to. The classic pattern
is one "isolated" table (spokes associate here, propagate to the
default table) so spokes reach shared services but not each other.
Labels exist to propagate to many tables in one statement -- label
production tables "prod" once and every connection can target the set.

## BGP peerings need the connection first

A `bgpConnections` entry peers the hub router with an NVA in a spoke --
but routes only flow when the NVA's VNet is attached through a hub
connection, and the peering references that connection's ID. Deploy
order in charts: hub → connection → hub update with the peering (or
reference the connection's output directly). The hub router's ASN is
always 65515; the NVA must use a different one.

## Deletion works bottom-up

ARM refuses to delete a hub that still has gateways or connections,
and refuses to delete the WAN while this hub exists. Tear down in
reverse creation order: connections → gateways → hub → WAN. "InUse"
errors during destroy almost always mean a level was skipped.
