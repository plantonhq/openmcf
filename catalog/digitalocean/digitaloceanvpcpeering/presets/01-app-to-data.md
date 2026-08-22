# App-to-Data Peering

This preset peers an application VPC with a data VPC, so app droplets reach databases over private addresses while the two networks keep separate lifecycles and separate firewall postures.

## When to Use

- The classic tier split: stateless app network peered to a locked-down data network
- Any two same-team VPCs that need private connectivity without collapsing into one network

## Key Configuration Choices

- **Both VPCs by reference** -- the peering follows the network resources; recreating a VPC re-wires the peering in the same apply.
- **Non-overlapping CIDRs are on you** -- DigitalOcean rejects overlapping peers, and VPC ranges are create-only; plan them before the first network exists.
- **Reachability is all-or-nothing** -- narrow access per host with droplet firewalls (e.g. the data VPC's databases trusting only app droplets).

## What You Get

A free private link between the two networks, active in minutes, with cross-VPC traffic staying on DigitalOcean's fabric at no charge.
