# Droplet Autoscale Pools on DigitalOcean

Runs a fleet of identical droplets that DigitalOcean keeps at your target size -- or grows and shrinks automatically on CPU/memory load between the bounds you set. Unhealthy members are replaced to hold the target. Integrates with Planton's Provider Connections for DigitalOcean API token management; SSH keys, the VPC, and the project are wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Droplet Autoscale Pool** -- the pool with your scaling mode (static count or dynamic bounds + utilization targets)
- **Member Droplets** -- provisioned and owned by the pool from your template (size, region, image, SSH keys, tags, networking)

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **An SSH key** -- a DigitalOceanSshKey resource; the API requires at least one key on every pool template.
- **A VPC and Project (optional)** -- referenced in the template when you want members placed explicitly.

### DigitalOcean Account

- Droplet quota for the pool's maximum size, and budget: every member is a real droplet billing hourly.

## After You Deploy

The pool reports `active` once every member is provisioned. Target members with tag-based firewall rules and load-balancer tags (tags follow the membership as it scales). And plan teardown deliberately: **destroying the pool destroys its member droplets** -- that is DigitalOcean's only delete for pools.
