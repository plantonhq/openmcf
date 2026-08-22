# DigitalOcean Droplet Autoscale Pool -- Operational Guide

What experience with this component teaches that the field reference cannot.

## Destroy destroys the droplets -- there is no other delete

DigitalOcean's only delete for an autoscale pool is the "dangerous" variant: the pool AND every member droplet it owns are terminated together. There is no adopt-the-members teardown, no orphan mode. Before destroying, drain traffic (load balancers, DNS) as if you were terminating that many droplets by hand -- because you are. The inverse hazard also holds: a FORGOTTEN pool keeps real droplets billing indefinitely, so pool hygiene is bill hygiene.

## Members are cattle -- keep state out of them

The pool creates and destroys members on its own schedule (scale events, health replacement, template rollouts). Anything written to a member's local disk is one scale-in from gone. Point members at managed databases, Spaces, or volumes owned elsewhere; use `userData` to bootstrap them identically on every boot.

## Dynamic scaling decides on agent metrics -- keep the agent on

CPU and memory targets are evaluated from the droplet monitoring agent's telemetry. Ship dynamic pools with `withDropletAgent: true` (the quick-start default here); without it, memory-based scaling has no data source at all. Static pools can skip the agent, though replacement health still benefits from it.

## Target the pool with tags, never with droplet IDs

Member droplet IDs churn with every scale event -- any firewall rule or load-balancer target list naming them goes stale immediately. The template's `tags` (plus the Planton labels both engines always apply) follow the membership automatically; tag-targeted firewall rules and load-balancer tag targets are the ONLY reliable way to address the fleet.

## Template changes roll the fleet

Editing the template (size, image, user data) applies in place on the pool, and DigitalOcean replaces members to converge on the new shape. Plan template edits like deployments: capacity dips while members roll. The image also has a read-back quirk -- DigitalOcean reports it as a numeric image ID even when you configured a slug; the modules keep your configured value, but a freshly IMPORTED pool will show an image diff on its first plan.

## Size the bounds for the bill, not just the load

A static pool bills `target_instances` droplets around the clock; a dynamic pool bills between `minInstances` and `maxInstances`. The max bound is your cost ceiling under sustained load (or under a runaway feedback loop) -- set it from budget, not optimism. Scale-ups respect `cooldownMinutes`, so short traffic spikes may ride on the existing members.

## What is deliberately NOT here

`public_networking` (schema-declared upstream but never sent in any create/update request at the pinned provider -- a field that cannot do anything; members always get public addresses, so restrict reachability with tag-targeted firewalls); per-member identity knobs (names are generated from the pool name -- members are interchangeable by design); and `current_utilization` / timestamps as outputs (volatile metrics, not identity -- read them from DigitalOcean monitoring).
