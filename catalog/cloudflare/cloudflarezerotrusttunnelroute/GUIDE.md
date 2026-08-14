# CloudflareZeroTrustTunnelRoute guide

Operational judgment for tunnel routes. The README covers what each field is; this covers how the pieces interact.

## A route is a promise WARP keeps

The route tells Cloudflare "this CIDR is reachable through that tunnel" — WARP clients in your Zero Trust organization then send matching traffic into the tunnel. The route creates fine even if the tunnel has no running connector; traffic just goes nowhere until cloudflared runs. Create-success is not end-to-end proof.

## CIDRs are unique per virtual network

Two routes for the same CIDR cannot coexist in one virtual network — the second create fails at the API. That is exactly what virtual networks are for: put overlapping ranges (two data centers both using 10.0.0.0/8) in separate virtual networks and let device enrollment pick the network. Omitting `virtual_network_id` lands the route in the account's default network, which is fine until the day it isn't.

## Prefer references over pasted IDs

`tunnel_id` and `virtual_network_id` accept references to the owning kinds' outputs. Wiring them as references makes the dependency explicit in the resource graph — the tunnel deploys first, the route follows, and a chart teardown retires them in the right order.

## Destroy is a soft delete

Like tunnels, a deleted route keeps its record with `deleted_at` set instead of 404ing. A recreated route for the same CIDR gets a new ID; nothing re-adopts the old one.
