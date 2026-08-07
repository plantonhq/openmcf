# Email Routing Rule on Cloudflare

Defines a single Cloudflare Email Routing rule for a zone: match inbound mail -- by a specific recipient or all messages -- and then drop it, forward it to verified destinations, or hand it to an Email Worker. Rules require Email Routing to be enabled on the zone (a `CloudflareEmailRoutingZone`) and are evaluated in priority order, so specific rules can take precedence over the zone's catch-all. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Routing Rule** -- a matcher set plus an action, evaluated at the configured priority
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Email Routing edit access (and Workers Scripts access when routing to an Email Worker). Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Email Routing enabled** -- the zone must already have Email Routing enabled (`CloudflareEmailRoutingZone`).
- **Verified destinations** -- a forwarding rule can only deliver to verified `CloudflareEmailRoutingAddress` mailboxes.

## Deploy

### Console

Open the deployment store, find **Email Routing Rule on Cloudflare**, and click **Deploy**. The creation wizard captures the zone and rule metadata (name, enabled, priority), at least one matcher (a specific recipient or all messages), and the action (drop / forward / Email Worker).

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareEmailRoutingRule
metadata:
  name: forward-support
  org: acme-corp
  env: prod
spec:
  zoneId:
    valueFrom:
      kind: CloudflareDnsZone
      name: acme-com
      fieldPath: status.outputs.zone_id
  name: forward-support
  matchers:
    - type: literal
      field: to
      value: support@example.com
  action:
    type: forward
    forwardTo:
      - value: ops@example.com
```

```shell
planton apply -f cloudflare-email-routing-rule.yaml
```

This forwards mail addressed to `support@example.com` to `ops@example.com`. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a routing rule. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Zone (`zoneId`)** -- The zone whose inbound mail this rule filters. Immutable -- changing it replaces the rule. Reference a `CloudflareDnsZone` to keep the dependency in the graph.

**Matchers (`matchers`)** -- At least one. A `literal` matcher (on the `to` field) matches a specific recipient; an `all` matcher matches every message. A message matches the rule if any matcher matches it.

**Action (`action.type`)** -- `drop` discards the message, `forward` delivers it to `forwardTo` destinations, `worker` hands it to an Email Worker referenced by `worker`.

**Priority (`priority`)** -- Lower numbers evaluate first. Put specific rules above broad ones so a catch-all-style rule does not shadow them.

## Outputs and Dependencies

### What This Component Consumes

The rule references a **CloudflareDnsZone** (via `zoneId`). A forwarding action references **CloudflareEmailRoutingAddress** mailboxes (via `action.forwardTo`); a worker action references a **CloudflareWorker** (via `action.worker`).

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `rule_id` | The Cloudflare-assigned identifier of the rule | Verification, dashboards |
| `zone_id` | The zone the rule belongs to | Auditing, grouping |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Role-address forwarding** -- a `literal` matcher on `support@` with a forwarding action routes a shared alias to a real inbox.

**Programmatic handling** -- a `worker` action hands matched mail to an Email Worker for parsing, ticket creation, or auto-responses.

## Works With

- [**Email Routing Zone on Cloudflare**](/cloud-catalog/cloudflare-email-routing-zone) -- enables Email Routing on the zone this rule belongs to
- [**Email Routing Address on Cloudflare**](/cloud-catalog/cloudflare-email-routing-address) -- the verified mailboxes a forwarding rule delivers to
- [**Worker on Cloudflare**](/cloud-catalog/cloudflare-worker) -- the Email Worker a worker rule hands matched mail to
