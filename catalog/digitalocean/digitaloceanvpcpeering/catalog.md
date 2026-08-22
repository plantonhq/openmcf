# VPC Peering on DigitalOcean

Connects two DigitalOcean VPCs over the private network fabric, so resources in both reach each other without touching the public internet -- free, and with free cross-VPC traffic. Integrates with Planton's Provider Connections for DigitalOcean API token management; both VPCs are wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **VPC Peering Connection** -- the named, symmetric private link between your two referenced VPCs

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Two VPCs** -- DigitalOceanVpc resources with NON-OVERLAPPING IP ranges (DigitalOcean rejects overlapping peers).

### DigitalOcean Account

- Nothing: peerings and the traffic crossing them are free.

## After You Deploy

Resources in each VPC reach the other network's private addresses immediately -- no route tables to edit. Remember the reachability is all-or-nothing between the two networks; per-host restriction stays with each droplet's firewall.
