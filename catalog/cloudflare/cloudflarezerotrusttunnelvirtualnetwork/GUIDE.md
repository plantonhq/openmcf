# CloudflareZeroTrustTunnelVirtualNetwork guide

Operational judgment for tunnel virtual networks. The README covers what each field is; this covers how the pieces interact.

## You need one when private ranges collide

An account starts with a built-in default virtual network, and plenty of deployments never need another. The moment two environments legitimately use the same private CIDR — staging and production both on 10.0.0.0/8, or an acquired company's network mirroring yours — each gets its own virtual network, and WARP device profiles select which one a user routes into. Creating them "just in case" adds routing surface for no benefit.

## The default flag is a lever, not a label

`is_default_network: true` does not just mark this network — it demotes the account's current default. Every route that was implicitly landing in the old default keeps pointing there, but every NEW route created without a `virtual_network_id` now lands here. Flip it in a maintenance window with the routes inventory in front of you.

## Names are identity for humans

Virtual network names are unique within the account and are how operators tell routing domains apart in the WARP client and dashboard. Name them after the network they represent (`prod-us-east`, `acquired-corp-lan`), not after the tool that made them.

## Empty it before deleting it

A virtual network with routes still in it will not delete — retire the routes first. A chart that owns both gets this ordering for free from the dependency graph.
