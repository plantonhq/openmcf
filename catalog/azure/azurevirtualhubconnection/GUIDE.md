# Azure Virtual Hub Connection -- Operational Guide

Judgment that does not fit in field references.

## The connection is free; the routing is the product

Attaching a VNet costs nothing -- what you are really configuring is
WHICH routes the spoke's traffic follows (the associated table) and WHO
learns the spoke's prefixes (the propagation targets). Leave routing
unset and you get any-to-any through the hub's default table; every
deliberate topology overrides one or both halves.

## Isolation is two statements, not one

An isolated spoke (a) associates with an isolated table and (b)
propagates only to where it should be REACHED from (shared services,
default). Getting only one half right leaks reachability in one
direction -- the spoke can be unreachable yet still see everyone, or
vice versa. Test both directions after changing either half.

## Address overlap fails late -- check it early

ARM validates the spoke's address space against the hub and every
already-connected network AT ATTACH TIME. An overlap that existed for
months surfaces only when this connection deploys. Before attaching,
diff the VNet's space against the hub prefix and the WAN's connected
estate.

## The override criteria is a one-way door

`staticVnetLocalRouteOverrideCriteria` is fixed once the connection is
created (ARM replaces the connection to change it) -- decide CONTAINS
vs EQUAL before the first deploy if static routes toward an NVA are in
the design, even if the routes themselves come later.

## Deletion order matters twice

The hub cannot be deleted while this connection exists, and this
connection cannot be deleted while a hub BGP peering references it.
Reverse creation order: BGP peering → connection → hub.
