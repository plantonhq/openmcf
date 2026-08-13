# Azure DNS Private Resolver -- Operational Guide

Judgment calls that matter when you run the resolver in production.

## One resolver per network is a hard wall -- put it in the hub

Azure allows exactly one resolver per virtual network. In hub-and-spoke, that means the hub owns THE resolver and every spoke consumes it through ruleset links (spokes need no resolver, no endpoints, not even peering for DNS to flow). Placing a resolver in a spoke buys nothing and burns the spoke's one slot.

## Carve the delegated subnets when you design the network

Every endpoint needs its own subnet -- delegated to `Microsoft.Network/dnsResolvers`, /28 to /24, carrying nothing else, one endpoint per subnet. ARM enforces all of it at deploy time, and the delegation lives on the subnet resource: a missing delegation surfaces as a deploy-time failure, not a manifest error. A /28 costs sixteen addresses and is all an endpoint ever needs; carve inbound and outbound subnets up front even if you deploy only one endpoint today.

## Pin the inbound IP if on-premises config is expensive to change

Dynamic allocation picks a free address at create; that address survives the endpoint's lifetime but NOT its replacement -- and every endpoint field except tags replaces the endpoint. If the inbound IP is fanned out across datacenter forwarder configs, use STATIC allocation and choose the address deliberately. Dynamic is fine when a config-management system owns the forwarder fleet.

## Endpoints are the billing meter and the capacity dial

The resolver object is free; each endpoint bills hourly and serves ~10,000 queries/second. One endpoint each way is the right day-one shape -- add endpoints (up to 5 each way) for throughput, not for redundancy: a single endpoint is already zone-resilient where the region has zones.

## The deploy-order edge runs through the outputs

Forwarding rulesets bind the resolver's outbound endpoint by ARM id -- reference this component's `outbound_endpoint_id` output rather than composing the id by hand, and deploy order takes care of itself (the reference IS the ordering edge). The same holds for on-premises pointing: read `inbound_endpoint_ip` from the outputs, never guess the subnet's .4 address.

## Deletes are slower than they look

Endpoint deletes poll past Azure's first "deleted" answer until the endpoint is verifiably gone (minutes each, sequential with the resolver's own delete). Budget teardown time accordingly in pipelines that create and destroy resolvers routinely.
