# AwsVpcPeering — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Peering is plumbing, not connectivity

An ACTIVE peering moves zero packets until both sides add routes toward the peer CIDR and open security groups. "Peering is up but nothing connects" is almost always a missing route, not a peering problem.

## The duplicate-connection trap

CreateVpcPeeringConnection on an already-peered VPC pair does not error — AWS returns the EXISTING connection's id. Declare each pair exactly once; a second instance for the same pair silently co-manages (and on destroy, deletes) the first one's connection.

## Deletes belong to the requester

The accept arm's destroy is a no-op at AWS: it abandons management, the peering stays ACTIVE, and the peer keeps routing. Only the request arm's destroy deletes the connection. Decommission cross-account peerings from the requester side first.

## Options need an ACTIVE connection

DNS-resolution options are a modification of an accepted peering — on a pending cross-account request AWS rejects them until the accepter accepts. Deploy the requester without the accepter-side option, let the peer accept, then set options from whichever side owns them.

## CIDR overlap is forever

AWS rejects peering between VPCs with overlapping CIDRs, and CIDRs are immutable — plan address space before you need the peering, not after. The non-transitive rule bites at scale: ten VPCs fully meshed is 45 peerings; consider Transit Gateway (already in the catalog) past a handful.
