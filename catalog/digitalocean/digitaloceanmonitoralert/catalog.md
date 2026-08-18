# Monitor Alert on DigitalOcean

Creates an alert policy on DigitalOcean's built-in metrics -- droplet CPU/memory/disk/bandwidth, load-balancer health and error rates, or managed-database utilization -- with email and Slack delivery. Targets are wired by typed references (droplets, load balancers, database clusters) or droplet tags, validated against the metric family. Integrates with Planton's Provider Connections for DigitalOcean API token management and ValueFromRef for target wiring.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Alert Policy** -- the metric, comparison, threshold, and sampling window
- **Targets** -- the referenced droplets/load balancers/database clusters, or tag-targeted droplet fleets
- **Notification Channels** -- email addresses and/or Slack webhooks (webhook URLs are treated as secrets)

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Targets** -- the droplets, load balancers, or database clusters to watch (referenced by name), or a droplet tag.

### DigitalOcean Account

- For Slack delivery: a Slack incoming-webhook URL. Alert emails may need to belong to the DigitalOcean team's verified members.

## After You Deploy

The policy starts evaluating immediately (unless `enabled: false`). Tag-targeted policies track fleet membership automatically -- droplets gaining or losing the tag join and leave the alert with no manifest change.
