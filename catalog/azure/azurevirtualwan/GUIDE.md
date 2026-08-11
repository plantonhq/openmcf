# Azure Virtual WAN -- Operational Guide

Judgment that does not fit in field references.

## The WAN is free; the decisions are not

The WAN object costs nothing and provisions in minutes -- every real
cost and constraint lives in the hubs and gateways created under it.
Treat this component as the place where the GLOBAL decisions are made
(tier, transit policy, breakout policy), because changing some of them
later means touching everything underneath.

## Standard unless you can say why not

Basic exists for one narrow shape: site-to-site VPN only, on Basic
hubs, at a lower price. Everything else -- ExpressRoute, point-to-site,
hub-to-hub transit, most routing features -- needs Standard. The
ratchet only turns one way (Basic upgrades to Standard in place;
Standard never downgrades), so Basic is a safe starting point ONLY if
site-to-site-only is genuinely the end state.

## Deletion works bottom-up

ARM refuses to delete a WAN that still has hubs, and a hub that still
has gateways or connections. Tear down in reverse creation order:
connections → gateways → hubs → WAN. In charts this ordering falls out
of the reference graph; when operating manually, expect "in use"
errors to mean "you skipped a level".

## Branch-to-branch is a security decision

The default (on) makes the WAN a full transit network: any VPN site
can reach any other through the hubs. Topologies that treat branches
as untrusted islands (retail, OT networks) set
`allowBranchToBranchTraffic: false` and force everything through
hub-hosted inspection instead. Flipping it later is a routing-behavior
change across the whole estate -- decide deliberately at creation.
