# Uptime Check on DigitalOcean

Creates an availability/latency probe on any external endpoint -- your site, your API, anything reachable -- run from DigitalOcean's global vantage regions, with alert rules for downtime, latency, and certificate expiry delivered by email or Slack. The alert rules ride the check as composed rows with one lifecycle. Integrates with Planton's Provider Connections for DigitalOcean API token management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Uptime Check** -- the probe: target, protocol, vantage regions
- **Alert Rules** -- configured only when `alerts` rows are set; one alert object per row (`down`, `down_global`, `latency`, `ssl_expiry`), each with its own notification channels

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.

### DigitalOcean Account

- Nothing on the account -- the probed target is external. For Slack delivery: a Slack incoming-webhook URL. Alert emails may need to belong to the DigitalOcean team's verified members.

## After You Deploy

Probing starts immediately (unless `enabled: false`); results appear in the control panel's Monitoring -> Uptime section, and alert rules begin evaluating as soon as they exist. https probes also unlock `ssl_expiry` alert rules for certificate-renewal safety nets.
