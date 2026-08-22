# Static Web Fleet

This preset holds exactly two identical web droplets in a VPC, with health-based replacement: a fixed-size fleet where DigitalOcean rebuilds any member that dies, without the operator paging anyone.

## When to Use

- A steady-traffic web tier that wants self-healing more than elasticity
- Any fixed-count worker fleet where droplets should be interchangeable cattle

## Key Configuration Choices

- **Static target of 2** -- capacity survives any single member's failure (and the count is the whole bill: two droplets around the clock).
- **SSH key and VPC by reference** -- the fleet composes with the account's key and network resources in one chart.
- **The `web` tag is the fleet's handle** -- point firewall rules and load-balancer targets at the tag, never at member IDs (members churn).
- **Agent enabled** -- replacement health and metrics visibility.

## What You Get

A self-healing fixed fleet -- and one warning to respect: destroying the pool destroys both member droplets; that is DigitalOcean's only delete.
