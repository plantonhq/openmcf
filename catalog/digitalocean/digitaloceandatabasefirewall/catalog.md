# Database Firewall on DigitalOcean

Declares the inbound trusted sources of a DigitalOcean managed database cluster: IP addresses and CIDR blocks, Droplets, Kubernetes clusters, App Platform apps, and Droplet tags -- each in its own typed list, with platform resources wired by reference instead of hand-copied ids. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for cluster, Droplet, DOKS, and App dependency wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Trusted-Sources Rule Set** -- the cluster's complete inbound allowlist, fanned out from the five typed lists to DigitalOcean's rule rows

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **A DigitalOceanDatabaseCluster** -- the cluster to protect, referenced by name (or an existing cluster's UUID as a literal).
- **Referenced sources** -- any Droplets, DOKS clusters, or Apps you trust by reference must exist (or deploy in the same chart).

### DigitalOcean Account

- Nothing beyond the cluster: trusted-source rules are free.

## After You Deploy

The cluster accepts connections only from the declared sources. Keep two facts in operational memory: declare ALL trusted sources in this ONE resource (a second rule set on the same cluster overwrites this one), and deleting this resource OPENS the database again -- the destroy clears the rule set rather than leaving the cluster protected.
