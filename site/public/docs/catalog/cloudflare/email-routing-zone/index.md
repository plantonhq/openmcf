---
title: "Email Routing Zone"
description: "Email Routing Zone deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareemailroutingzone"
---

# Email Routing Zone on Cloudflare

Enables Cloudflare Email Routing on a DNS zone -- the anchor of the Email Routing family. Enabling provisions the records inbound mail needs (MX, SPF, DKIM) automatically and configures the single per-zone catch-all rule that decides what happens to mail no other routing rule matched: leave Cloudflare's default (drop), forward to verified destinations, or hand it to an Email Worker. Individual routing rules and destination addresses are separate Cloud Resources that build on top of an enabled zone. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Email Routing enablement** -- turns Email Routing on for the zone
- **DNS records** -- the required MX, SPF, and DKIM records, created automatically (optionally locked against out-of-band edits)
- **Catch-all rule** (optional) -- the single per-zone rule for unmatched mail: drop, forward to addresses, or send to an Email Worker

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Email Routing edit access (and DNS edit access for the auto-provisioned records). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **An active zone** -- the domain must already exist as a Cloudflare zone (manage it with a `CloudflareDnsZone`, or reference an externally-managed zone ID).
- **Verified destinations** -- a forwarding catch-all can only deliver to verified `CloudflareEmailRoutingAddress` mailboxes.

## Deploy

### Console

Open the deployment store, find **Email Routing Zone on Cloudflare**, and click **Deploy**. The creation wizard captures the zone, whether the auto-created DNS records are locked, and the optional catch-all action. Leave the catch-all as **None** to keep Cloudflare's default (drop, disabled).

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareEmailRoutingZone
metadata:
  name: acme-com-email-routing
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-com
      fieldPath: status.outputs.zone_id
  catchAll:
    enabled: true
    type: forward
    forwardTo:
      - value: catch-all@example.com
```

```shell
planton apply -f cloudflare-email-routing-zone.yaml
```

This enables Email Routing on the `acme-com` zone and forwards all otherwise-unmatched mail to a catch-all mailbox. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when enabling Email Routing. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone (`zoneId`)** -- The zone to enable Email Routing on. Immutable -- changing it replaces the resource. Reference a `CloudflareDnsZone` to keep the dependency in the graph.

**Catch-all Action (`catchAll.type`)** -- What happens to unmatched mail. `drop` discards it silently (Cloudflare's default), `forward` delivers it to `forwardTo` destinations, `worker` hands it to an Email Worker. Omit `catchAll` to leave the default.

**Lock DNS Records (`lockDnsRecords`)** -- Manage the auto-created MX/SPF/DKIM records explicitly so they cannot be edited out-of-band.

## Outputs and Dependencies

### What This Component Consumes

The zone references a **CloudflareDnsZone** (via `zoneId`). A forwarding catch-all references **CloudflareEmailRoutingAddress** mailboxes (via `catchAll.forwardTo`); a worker catch-all references a **CloudflareWorker** (via `catchAll.worker`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `zone_id` | The zone Email Routing was enabled on | Referenced by `CloudflareEmailRoutingRule`s for the same zone |
| `enabled` | Whether Email Routing is enabled | Health checks |
| `status` | The Email Routing configuration status (ready / unconfigured / misconfigured) | Diagnostics |
| `name` | The zone's domain name | Dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Forward-all** -- a forwarding catch-all funnels every inbound message to one or more real mailboxes; add specific rules later for per-address routing.

**Strict routing** -- leave the catch-all at drop and rely on explicit `CloudflareEmailRoutingRule`s, so only addresses you route are accepted.

## Works With

- [**Email Routing Rule on Cloudflare**](/cloud-catalog/cloudflare-email-routing-rule) -- per-recipient rules that take precedence over the catch-all
- [**Email Routing Address on Cloudflare**](/cloud-catalog/cloudflare-email-routing-address) -- the verified mailboxes a forwarding catch-all delivers to
- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- the zone Email Routing is enabled on
