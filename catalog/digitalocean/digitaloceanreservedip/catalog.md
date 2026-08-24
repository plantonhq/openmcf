# Reserved IP on DigitalOcean

Reserves a static public IP address (IPv4 or IPv6) in a DigitalOcean region and optionally assigns it to a droplet. The address outlives the droplets behind it -- re-pointing it between droplets is the classic manual-failover building block. Integrates with Planton's Provider Connections for DigitalOcean API token management; the droplet is wired by reference.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Reserved IP** -- the static address in your chosen region (IPv4 or IPv6)
- **Droplet assignment** (when a droplet is referenced) -- the address pointed at your droplet

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Droplet** (optional) -- a DigitalOceanDroplet in the SAME region, if assigning at create time.

### DigitalOcean Account

- Billing awareness: an UNASSIGNED reserved IPv4 accrues a monthly charge (prorated hourly); assigned ones are free, and IPv6 is free either way.

## After You Deploy

Point DNS at the `reserved_ip_address` output instead of any droplet's own address -- replacing the droplet later means re-pointing the reservation, not re-publishing DNS. Destroy the reservation when it stops being needed; the idle state is the one that bills.
